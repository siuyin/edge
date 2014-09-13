package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	zmq "github.com/pebbe/zmq4"
	"log"
	"os"
)

func genNewUserLogin() string {
	return "http://localhost/session/342423443"
}

func checkUserExists(db *sql.DB, email string) (bool, string) {
	rows, err := db.Query("select id from users where email = $1", email)
	defer rows.Close()
	if err != nil {
		log.Fatal("query", err)
	}
	if rows.Next() {
		return true, ""
	} else {
		return false, genNewUserLogin()
	}
}
func promoImages(db *sql.DB) []string {
	rows, err := db.Query("select image_url from items limit 5")
	if err != nil {
		log.Fatal("query", err)
	}
	imgs := make([]string, 0)
	for rows.Next() {
		var img string
		if err := rows.Scan(&img); err != nil {
			log.Fatal(err)
		}
		imgs = append(imgs, img)
		fmt.Println(img)
	}
	return imgs
}
func genPointsCode() string {
	return "3449080293840"
}

func main() {
	type InRec struct {
		Email string `json:"email"`
	}
	type OutRec struct {
		Email      string   `json:"email"`
		Existing   bool     `json:"existing"`
		LoginLink  string   `json:"login_link"`
		Images     []string `json:"images"`
		PointsCode string   `json:"pointscode"`
	}
	db, err := sql.Open("postgres", "host=/var/run/postgresql user=www-data dbname=edge_development sslmode=disable")
	if err != nil {
		log.Fatal("dbopen", err)
	}
	//  Socket to talk to clients
	responder, _ := zmq.NewSocket(zmq.REP)
	defer responder.Close()
	endpoint := "ipc:///var/www/socks/email.ipc"
	responder.Bind(endpoint)
	os.Chmod(endpoint, 0777)

	fmt.Println("Starting emailsvr")
	inrec := new(InRec)
	outrec := new(OutRec)

	for {
		//  Wait for next request from client
		msg, _ := responder.Recv(0)
		fmt.Println("Received ", msg)
		json.Unmarshal([]byte(msg), inrec)
		outrec.Email = inrec.Email
		outrec.Existing, outrec.LoginLink = checkUserExists(db, inrec.Email)
		outrec.Images = promoImages(db)
		outrec.PointsCode = genPointsCode()

		//  Send reply back to client
		reply, _ := json.Marshal(outrec)
		responder.Send(string(reply), 0)
	}
}
