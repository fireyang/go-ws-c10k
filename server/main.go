package main

import (
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	// conn.SetPongHandler(func(string) error {
	//     log.Printf("pong")
	//     conn.SetReadDeadline(time.Now().Add(pongWait))
	//     return nil
	// })
	conn.SetPingHandler(func(string) error {
		// log.Printf("ping")
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	cur := atomic.LoadInt32(&conn_num)

	for {
		// time.Sleep(1 * time.Second)
		// log.Printf("start read")
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			log.Printf("err2:id,%d, err,%v", cur, err)
			break
		}
		log.Printf("msg:%v", msg)
		// msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
		// w, err := conn.NextWriter(websocket.TextMessage)
		// if err != nil {
		//     return
		// }
		// log.Printf("write msg:%v", string(msg))
		// w.Write(msg)
		// if err := w.Close(); err != nil {
		//     return
		// }
	}
	i := atomic.AddInt32(&conn_num, -1)
	log.Printf("close conn:%d", i)
}

var conn_num = int32(0)

func main() {
	log.Println("start server")
	addr := ":8080"
	// i := 0
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		i := atomic.AddInt32(&conn_num, 1)
		log.Printf("new conn:%d, %v", i, r.RemoteAddr)
		serveWs(w, r)
		log.Printf("close ws:%d", i)
	})
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
