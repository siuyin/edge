package main

import (
	"encoding/json"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"net/http"
	"net/url"
	"os"
)

func main() {
	type SMS struct {
		Dest string `json:"dest"`
		Msg  string `json:"msg"`
	}
	socket, _ := zmq.NewSocket(zmq.REP)
	defer socket.Close()
	endpoint := "ipc:///var/www/socks/sms.ipc"
	socket.Bind(endpoint)

	hoiioURL := "https://secure.hoiio.com/open/sms/send"
	appID := os.Getenv("HOIIO_APP_ID")
	accessToken := os.Getenv("HOIIO_ACCESS_TOKEN")
	v := url.Values{}
	v.Add("app_id", appID)
	v.Add("access_token", accessToken)
	fmt.Println("Starting hoiiosms")
	sms := new(SMS)
	for {
		msg, _ = socket.Recv(0)
		fmt.Println("Received ", msg)
		json.Unmarshal([]byte(msg), sms)
	}
	v.Add("dest", sms.Dest)
	v.Add("msg", sms.Msg)
	str := hoiioURL + "?" + v.Encode()
	fmt.Println(str)
	/*
		res, err := http.Get(str)
		if err != nil {
			fmt.Println("get error:", err)
		}
		fmt.Println(res)
	*/
}
