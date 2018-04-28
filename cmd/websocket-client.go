// 脚本执行
// go run websocket-client.go --addr local.cc.ndmooc.com --unitid A16 --token 53ba4bf6af4e4ebd8bdfd09cb99c5c9079485fcb

package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

// 172.17.0.2:5201/v2/ngx/center/units/A1001/?token=95bc6aff422d12777df02ba42d6167ff4c63b5bf
var addr = flag.String("addr", "", "http service address")
var unitid = flag.String("unitid", "", "unit id")
var token = flag.String("token", "", "accesstoken or devicetoken")

type Msg struct {
	Typ byte
	Uid int32
	Act byte
	//Data string
}

func main() {
	flag.Parse()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	v := url.Values{}
	v.Set("token", *token)
	rawquery := v.Encode()
	path := fmt.Sprintf("v2/ngx/center/units/%s/", *unitid)
	u := url.URL{Scheme: "ws", Host: *addr, Path: path, RawQuery: rawquery}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial: ", err)
	}
	defer c.Close()

	// 注册
	reg := map[string]string{
		"act": "1",
		"os":  "1",
		"vi":  "1",
		"hw":  "0",
	}
	regjson, err := json.Marshal(reg)
	if err != nil {
		log.Fatal(err)
	}
	err = c.WriteMessage(websocket.TextMessage, []byte(regjson))
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			fmt.Println("Bye!")
			return
			//case t := <-ticker.C:
		case <-ticker.C:
			//err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			//if err != nil {
			//	log.Println("write:", err)
			//	return
			//}

			// TODO 为什么该数据类型不能以二进制流发放
			data := &Msg{1, 123, 2}
			//data := &Msg{1, 123, 2}
			buf := new(bytes.Buffer)
			binary.Write(buf, binary.BigEndian, data)
			bytes := buf.Bytes()
			bytes = append(bytes, []byte("BBB")...)
			fmt.Println(bytes)
			err := c.WriteMessage(websocket.BinaryMessage, bytes)
			if err != nil {
				log.Println("write: ", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message
			// and then waiting (with timeout) for the server to close
			// the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
				fmt.Println("Bye bye!")
			}
			return
		}
	}
}
