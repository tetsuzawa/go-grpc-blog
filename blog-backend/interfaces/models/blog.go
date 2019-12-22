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

type Blog struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id,omitepty"`
	AuthorID string             `json:"author_id" bson:"author_id"`
	Title    string             `json:"title" bson:"title"`
	Content  string             `json:"content" bson:"content"`
}

type BlogServicer struct{}

func DocumentToBlogpb(data *Blog) *blogpb.Blog {
	return &blogpb.Blog{
		Id:       data.ID.Hex(),
		AuthorId: data.AuthorID,
		Title:    data.Title,
		Content:  data.Content,
	}

}

func NewBlog(id, authorId, title, content string) *blogpb.Blog {
	return &blogpb.Blog{
		Id:       id,
		AuthorId: authorId,
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

	//data := blogpb.Blog{
	//	//ID:        blog.Id,  //empty, Mongodb generates a unique object ID
	//	AuthorId: ,
	//	Title:
	//	Content:
	//}
	data := Blog{
		ID:       primitive.NewObjectID(),
		AuthorID: blog.GetAuthorId(),
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
	old , ok := result.InsertedID.(primitive.ObjectID)
	if !ok{
		log.Printf("failed to convert type: %v\n", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}
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

	//var blog = new(blogpb.Blog)
	var data Blog
	err = BlogCollection.FindOne(ctxMongo, filter).Decode(&data)
	if err != nil {
		log.Printf("failed to decode document at Decode: %v\n", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}
	return &blogpb.ReadBlogRes{Blog: DocumentToBlogpb(&data)}, nil
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
	//b := blogpb.Blog{
	//	AuthorId: blog.AuthorId,
	//	Title:    blog.Title,
	//	Content:  blog.Content,
	//}
	b := Blog{
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetTitle(),
	}
	update := bson.M{"$set": b}

	_, err = BlogCollection.UpdateOne(ctxMongo, filter, update)
	if err != nil {
		log.Printf("failed to update document at UpdateOne: %v\n", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}
	var data Blog
	err = BlogCollection.FindOne(ctxMongo, filter).Decode(&data)
	if err != nil {
		log.Printf("failed to decode document at Decode: %v\n", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}

	return &blogpb.UpdateBlogRes{Blog: DocumentToBlogpb(&data)}, nil
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
	if result.DeletedCount == 0 {
		err = fmt.Errorf("not found error: could not find document")
		log.Printf("failed to delete document at DeleteOne (document count issue): %v", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}
	return &blogpb.DeleteBlogRes{IsSuccess: true}, nil
}

func (u *BlogServicer) List(req *blogpb.ListBlogReq, stream blogpb.BlogData_ListServer) error {
	BlogCollection := db.Collection(u.TableName())
	log.Printf("List blog invoked\n")

	//list, err := BlogCollection.Find(ctxMongo, primitive.D{{}})
	list, err := BlogCollection.Find(ctxMongo, bson.D{})
	if err != nil {
		log.Printf("failed to list document at Find: %v\n", err)
		return status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}

	for list.Next(ctxMongo) {
		var data Blog
		err = list.Decode(&data)
		if err != nil {
			log.Printf("failed to decode document at Decode: %v\n", err)
			return status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
		}
		err = stream.Send(&blogpb.ListBlogRes{Blog: DocumentToBlogpb(&data)})
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
