package cmd

import (
	"fmt"
	blogpb "github.com/tetsuzawa/go-grpc-blog/protocols/blog"
	"github.com/tetsuzawa/go-grpc-test/proto/cat"
	"github.com/tetsuzawa/go-grpc-test/proto/ping"
	"google.golang.org/grpc"
	"log"
	"os"
)

var blogClient  blogpb.BlogDataClient

func init() {
	conn, err := grpc.Dial("127.0.0.1:31060", grpc.WithInsecure())
	if err != nil {
		log.Fatal("client connection error:", err)
	}
	defer conn.Close()

	blogClient = blogpb.NewBlogDataClient(conn)

	fmt.Printf("result:%#v \n", catRes)
	fmt.Printf("error::%#v \n", err)

}

