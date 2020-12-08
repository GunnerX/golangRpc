package main

import (
	"encoding/json"
	"fmt"
	"golangRpc"
	"golangRpc/codec"
	"log"
	"net"
	"time"
)

func startServer(addr chan string) {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	golangRpc.Accept(l)
}

func main() {
	addr := make(chan string)

	go startServer(addr)

	conn, _ := net.Dial("tcp", <-addr)
	defer func() {
		conn.Close()
	}()

	time.Sleep(time.Second)
	// send options
	json.NewEncoder(conn).Encode(golangRpc.DefaultOption)
	cc := codec.NewGobCodec(conn)

	// send req & receive rsp
	for i := 0; i < 5; i++ {
		log.Printf("in for i = %d", i)
		h := &codec.Header{
			ServiceMethod: "Foo.Sum",
			Seq:           uint64(i),
		}
		cc.Write(h, fmt.Sprintf("golangRpc req %d", h.Seq))
		cc.ReadHeader(h)
		var reply string
		cc.ReadBody(&reply)
		log.Println("reply: ", reply)
	}

}
