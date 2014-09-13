package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	zmq "github.com/pebbe/zmq4"
	"log"
)

func main() {
	type Item struct {
		ID       int     `json:"id"`
		Name     string  `json:"name"`
		ImageURL string  `json:"image_url"`
		Price    float64 `json:"price"`
		Points   int     `json:"points"`
	}

	db, err := sql.Open("postgres", "host=/var/run/postgresql user=www-data dbname=edge_development sslmode=disable")
	if err != nil {
		log.Fatal("dbopen", err)
	}
	//  Socket to talk to clients
	responder, _ := zmq.NewSocket(zmq.REP)
	defer responder.Close()
	endpoint := "ipc:///var/www/socks/itemlist.ipc"
	responder.Bind(endpoint)

	fmt.Println("Starting itemlistsvr")
	var recs []*Item
	for {
		recs = make([]*Item, 0)
		rec := new(Item)
		msg, _ := responder.Recv(0)
		fmt.Println("Received ", msg)

		rows, err := db.Query("select id, name, image_url,price::money::numeric::float8,points from items")
		if err != nil {
			log.Fatal("query", err)
		}
		for rows.Next() {
			if err := rows.Scan(&rec.ID, &rec.Name, &rec.ImageURL, &rec.Price, &rec.Points); err != nil {
				log.Fatal("scan:", err)
			}
			fmt.Println(rec.ID, rec.Name)
			recs = append(recs, rec)
		}
		rows.Close()
		//  Send reply back to client
		reply, _ := json.Marshal(recs)
		responder.Send(string(reply), 0)
	}
}
