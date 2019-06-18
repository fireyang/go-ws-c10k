package main

import (
	"context"
	"log"
	"net"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var wg sync.WaitGroup

func main() {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	localAddr := os.Args[1]
	k := 0
	for i := 0; i < 20; i++ {
		for j := 0; j < 1000; j++ {
			k++
			wg.Add(1)
			go newConn(k, u.String(), localAddr)
		}
		time.Sleep(1 * time.Second)
	}
	// interrupt := make(chan os.Signal, 1)
	// signal.Notify(interrupt, os.Interrupt)
	// for {
	//     select {
	//     case <-interrupt:
	//         return
	//     }
	// }
	wg.Wait()
	// newConn(1, u.String(), localAddr)
}

func newConn(id int, url, localAddr string) {
	defer wg.Done()
	dialer := websocket.DefaultDialer
	// netAddr := &net.TCPAddr{IP: []byte("192.168.190.154")}
	// netAddr, err := net.ResolveTCPAddr("tcp4", "192.168.190.154:0")
	// log.Printf("33333333333")
	netAddr, err := net.ResolveTCPAddr("tcp4", localAddr)
	netAddr.Port = 0
	// netAddr.Port = id + netAddr.Port
	// netAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:10000")
	if err != nil {
		log.Printf("e,%v", err)
		return
	}
	log.Printf("11111111111,%d", netAddr.Port)
	dialer.NetDialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		d := net.Dialer{LocalAddr: netAddr}
		return d.DialContext(ctx, network, addr)
	}
	// log.Printf("22222222222")
	c, _, err := dialer.Dial(url, nil)
	if err != nil {
		log.Printf("dial:%v", err)
		return
	}
	defer c.Close()

	log.Printf("new client, %d", id)
	// interrupt := make(chan os.Signal, 1)
	// signal.Notify(interrupt, os.Interrupt)
	go func() {
		for {
			if err := c.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
			time.Sleep(20 * time.Second)
		}
	}()
	for {
		// log.Printf("ping")
		_, _, err := c.ReadMessage()
		if err != nil {
			//     if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			//         log.Printf("error: %v", err)
			//     }
			log.Printf("err2:id,%d, err,%v", id, err)
			break
		}
		// select {
		// case <-interrupt:
		//     log.Println("interrupt")
		//     err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		//     if err != nil {
		//         log.Println("write close:", err)
		//     }
		//     return
		// }
	}
}
