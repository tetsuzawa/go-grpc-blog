package main

import (
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	"github.com/tetsuzawa/go-grpc-blog/interfaces/models"
	"github.com/tetsuzawa/go-grpc-blog/protocols/blog"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("Starting server...")
	address := "0.0.0.0:31090"
	listenPort, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalln(err)
	}
	server := grpc.NewServer()
	blogService := &models.BlogServicer{}
	// 実行したい実処理をseverに登録する
	blogpb.RegisterBlogDataServer(server, blogService)

	go func() {
		err = server.Serve(listenPort)
		if err != nil {
			log.Fatalln(err)
		}
	}()
	log.Println("Server successfully started on port:", address)

	// shutdown signal
	c := make(chan os.Signal)

	signal.Notify(c, os.Interrupt)

	<-c

	log.Println("Stopping the server...")
	server.Stop()
	log.Printf("Closing MongoDB connection...")
	err = models.Disconnect()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Done.")
}
