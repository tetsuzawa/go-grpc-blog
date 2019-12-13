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
	"log"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	blogpb "github.com/tetsuzawa/go-grpc-blog/protocols/blog"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read blog",
	Long:  `Read blog reads a blog content`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("read called")
		id, err := cmd.Flags().GetString("id")
		check(err)

		req := &blogpb.ReadBlogReq{
			Id: id,
		}

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

		conn, err := grpc.Dial("127.0.0.1:31060", grpc.WithInsecure())
		if err != nil {
			log.Fatal("client connection error:", err)
		}
		defer conn.Close()

		blogClient := blogpb.NewBlogDataClient(conn)
		res, err := blogClient.Read(ctx, req)
		check(err)
		fmt.Printf("reqest=%v\n", req)
		fmt.Printf("response=%v\n", res)
	},
}

func init() {
	rootCmd.AddCommand(readCmd)

	readCmd.Flags().StringP("id", "i", "ididid", "id option")

	var err error
	err = readCmd.MarkFlagRequired("id")
	check(err)
}
