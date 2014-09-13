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
	type User struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		MobileNo string `json:"mobile_no"`
	}
	type InRec struct {
		UserID int `json:"user_id"`
		ItemID int `json:"item_id"`
	}
	db, err := sql.Open("postgres", "host=/var/run/postgresql user=www-data dbname=edge_development sslmode=disable")
	if err != nil {
		log.Fatal("dbopen", err)
	}
	//  Socket to talk to clients
	responder, _ := zmq.NewSocket(zmq.REP)
	defer responder.Close()
	endpoint := "ipc:///var/www/socks/buy.ipc"
	responder.Bind(endpoint)

	fmt.Println("Starting buysvr")
	rec := new(User)
	inrec := new(InRec)
	for {
		msg, _ := responder.Recv(0)
		fmt.Println("Received ", msg)
		json.Unmarshal([]byte(msg), inrec)
fmt.Println(inrec.UserID)

		rows, err := db.Query("select id, name, email,mobile_no from users where id = $1", inrec.UserID)
		if err != nil {
			log.Fatal("query", err)
		}
		for rows.Next() {
			if err := rows.Scan(&rec.ID, &rec.Name, &rec.Email, &rec.MobileNo); err != nil {
				log.Fatal(err)
			}
			fmt.Println(rec.ID, rec.Name)
		}
		rows.Close()
		//  Send reply back to client
		reply, _ := json.Marshal(rec)
		responder.Send(string(reply), 0)
	}
}
