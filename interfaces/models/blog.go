package models

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	blogpb "github.com/tetsuzawa/go-grpc-blog/protocols/blog"
)

//type Blog struct {
//	ID        string `json:"id" bson:"id"`
//	BlogName  string `json:"Blog_name" bson:"Blog_name"`
//	Password  string `json:"password" bson:"password"`
//	FirstName string `json:"first_name" bson:"first_name"`
//	LastName  string `json:"last_name" bson:"last_name"`
//	Role      string `json:"role" bson:"role"`
//}

type BlogServicer struct{}

func NewBlog(id, authorId, title, content string) *blogpb.Blog {
	return &blogpb.Blog{
		Id:       id,
		AutherId: authorId,
		Title:    title,
		Content:  content,
	}
}

func (u *BlogServicer) TableName() string {
	return GetTableName(tableNameBlogData)
}

func (u *BlogServicer) Create(ctx context.Context, req *blogpb.CreateBlogReq) (*blogpb.CreateBlogRes, error) {
	BlogCollection := db.Collection(u.TableName())

	blog := req.GetBlog()

	log.Printf("Create blog invoked. Blog:%v\n", blog)

	data := blogpb.Blog{
		//ID:        blog.Id,  //empty, Mongodb generates a unique object ID
		AuthorId: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}

	result, err := BlogCollection.InsertOne(ctxMongo, data)
	if err != nil {
		// return internal gRPC error to be handled later
		log.Printf("failed to insert document at InsertOne: %v\n", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}

	// add the id to blog, first cast the "generic type" (go doesn't have real generics yet) to an Object ID.
	old := result.InsertedID.(primitive.ObjectID)
	//	// Convert the object id to it's string counterpart
	blog.Id = old.Hex()
	return &blogpb.CreateBlogRes{Blog: blog}, nil
}

func (u *BlogServicer) Read(ctx context.Context, req *blogpb.ReadBlogReq) (*blogpb.ReadBlogRes, error) {
	BlogCollection := db.Collection(u.TableName())

	id := req.GetId()

	log.Printf("Read blog invoked. ID:%v\n", id)

	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("failed to decode id at ObjectIdFromHex: %v\n", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}
	filter := bson.M{"_id": oID}

	var blog = new(blogpb.Blog)
	err = BlogCollection.FindOne(ctxMongo, filter).Decode(blog)
	if err != nil {
		log.Printf("failed to decode document at Decode: %v\n", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}
	return &blogpb.ReadBlogRes{Blog: blog,}, nil
}

func (u *BlogServicer) Update(ctx context.Context, req *blogpb.UpdateBlogReq) (*blogpb.UpdateBlogRes, error) {
	BlogCollection := db.Collection(u.TableName())

	blog := req.GetBlog()

	log.Printf("Read blog invoked. Blog:%v\n", blog)

	oID, err := primitive.ObjectIDFromHex(blog.Id)
	if err != nil {
		log.Printf("failed to decode id at ObjectIdFromHex: %v\n", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}
	filter := bson.M{"_id": oID}
	b := blogpb.Blog{
		AuthorId: blog.AuthorId,
		Title:    blog.Title,
		Content:  blog.Content,
	}
	update := bson.M{"$set": b}

	_, err = BlogCollection.UpdateOne(ctxMongo, filter, update)
	if err != nil {
		log.Printf("failed to update document at UpdateOne: %v\n", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}
	err = BlogCollection.FindOne(ctxMongo, filter).Decode(blog)
	if err != nil {
		log.Printf("failed to decode document at Decode: %v\n", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}

	return &blogpb.UpdateBlogRes{Blog: blog}, nil
}

func (u *BlogServicer) Delete(ctx context.Context, req *blogpb.DeleteBlogReq) (*blogpb.DeleteBlogRes, error) {
	BlogCollection := db.Collection(u.TableName())

	id := req.GetId()

	log.Printf("Delete blog invoked. ID:%v\n", id)

	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("failed to decode id at ObjectIdFromHex: %v\n", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}
	filter := bson.M{"_id": oID}

	result, err := BlogCollection.DeleteOne(ctxMongo, filter)
	if err != nil {
		log.Printf("failed to delete document at DeleteOne: %v\n", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}
	if result.DeletedCount == 1 {
		err = fmt.Errorf("not found error: could not find document")
		log.Printf("failed to delete document at DeleteOne (document count issue): %v", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}
	return &blogpb.DeleteBlogRes{IsSuccess: true}, nil
}

func (u *BlogServicer) List(ctx context.Context, req *blogpb.ListBlogReq, stream blogpb.BlogData_ListServer) error {
	BlogCollection := db.Collection(u.TableName())
	log.Printf("List blog invoked\n")

	list, err := BlogCollection.Find(ctxMongo, primitive.D{{}})
	if err != nil {
		log.Printf("failed to list document at Find: %v\n", err)
		return status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}

	for list.Next(ctxMongo) {
		var blog = new(blogpb.Blog)
		err = list.Decode(blog)
		if err != nil {
			log.Printf("failed to decode document at Decode: %v\n", err)
			return status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
		}
		err = stream.Send(&blogpb.ListBlogRes{Blog: blog})
		if err != nil {
			log.Printf("failed to send document at Send: %v\n", err)
			return status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
		}
	}
	err = list.Close(ctxMongo)
	if err != nil {
		log.Printf("failed to close mongo list at Close: %v\n", err)
		return status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}

	return nil
}

//func GetBlog(id string) (*BlogServicer, error) {
//	BlogCollection := db.Collection(GetTableName(tableNameBlogData))
//
//	filter := bson.D{{"id", id}}
//
//	var u BlogServicer
//	err := BlogCollection.FindOne(ctxMongo, filter).Decode(&u)
//	if err != nil {
//		return nil, errors.Wrap(err, "failed to find data at FindOne()")
//	}
//	return NewBlog(u.ID, u.BlogName, u.Password, u.FirstName, u.LastName, u.Role), nil
//}
//
//func GetBlogByBlogName(un string) (*BlogServicer, error) {
//	BlogCollection := db.Collection(GetTableName(tableNameBlogData))
//
//	filter := bson.D{{"Blog_name", un}}
//
//	var u BlogServicer
//	err := BlogCollection.FindOne(ctxMongo, filter).Decode(&u)
//	if err != nil {
//		return nil, errors.Wrap(err, "failed to find data at FindOne()")
//	}
//	return NewBlog(u.ID, u.BlogName, u.Password, u.FirstName, u.LastName, u.Role), nil
//}
