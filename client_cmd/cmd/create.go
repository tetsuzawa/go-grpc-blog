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

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create blog",
	Long:  `Create blog creates a blog content`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
		id, err := cmd.Flags().GetString("id")
		check(err)
		author_id, err := cmd.Flags().GetString("author_id")
		check(err)
		title, err := cmd.Flags().GetString("title")
		check(err)
		content, err := cmd.Flags().GetString("content")
		check(err)

		blog := &blogpb.Blog{
			Id:       id,
			AutherId: author_id,
			Title:    title,
			Content:  content,
		}

		req := &blogpb.CreateBlogReq{
			Blog: blog,
		}

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

		conn, err := grpc.Dial("127.0.0.1:31060", grpc.WithInsecure())
		if err != nil {
			log.Fatal("client connection error:", err)
		}
		defer conn.Close()

		blogClient := blogpb.NewBlogDataClient(conn)
		res, err := blogClient.Create(ctx, req)
		check(err)
		fmt.Printf("reqest=%v\n", req)
		fmt.Printf("response=%v\n", res)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("id", "i", "ididid", "id option")
	createCmd.Flags().StringP("author_id", "a", "authoridid", "author id option")
	createCmd.Flags().StringP("title", "t", "Sample title", "title option")
	createCmd.Flags().StringP("content", "c", "sample content content", "content option")

	var err error
	err = createCmd.MarkFlagRequired("id")
	check(err)
	err = createCmd.MarkFlagRequired("author_id")
	check(err)
	err = createCmd.MarkFlagRequired("title")
	check(err)
	err = createCmd.MarkFlagRequired("content")
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
