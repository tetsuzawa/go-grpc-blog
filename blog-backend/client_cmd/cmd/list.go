/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"

	blogpb "github.com/tetsuzawa/go-grpc-blog/protocols/blog"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List blog",
	Long:  `List blog lists a blog content`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")

		conn, err := grpc.Dial("127.0.0.1:31060", grpc.WithInsecure())
		if err != nil {
			log.Fatal("client connection error:", err)
		}
		defer conn.Close()

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

		blogClient := blogpb.NewBlogDataClient(conn)
		req := &blogpb.ListBlogReq{}
		fmt.Printf("reqest=%v\n", req)
		stream, err := blogClient.List(ctx, req)
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			check(err)
			fmt.Printf("response=%v\n", res)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
