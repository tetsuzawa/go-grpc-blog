package models

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"

	"github.com/tetsuzawa/go-grpc-blog/config"
)

const (
	tableNameBlogData = "blog"
)

//const tFormat = "2006-01-02 15:04:05"
const tFormat = "2006-01-02T15:04:05.000Z"

func GetTableName(name string) string {
	return fmt.Sprintf("%s", name)
}

var ctxMongo context.Context
var client *mongo.Client
var db *mongo.Database

func init() {
	log.Println("Connecting to MongoDB...")

	var err error
	err = godotenv.Load() //Load env.file
	if err != nil {
		log.Fatalln(errors.Wrap(err, "failed to load .env file at godotenv.Load()"))
	}

	/*
		dbName := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			config.Cfg.DB.Host,
			config.Cfg.DB.Port,
			config.Cfg.DB.Name,
			//config.Cfg.DB.ETC,
		)
	*/
	fmt.Println(config.Cfg.DB.Host)
	fmt.Println(config.Cfg.DB.Port)
	fmt.Println(config.Cfg.DB.Name)
	dbName := fmt.Sprintf("mongodb://%s:%d/%s",
		config.Cfg.DB.Host,
		config.Cfg.DB.Port,
		config.Cfg.DB.Name,
		//config.Cfg.DB.ETC,
	)


	log.Println("Connecting to MongoDB...")
	clientOptions := options.Client().ApplyURI(dbName)

	//sha1Pass := sha256.Sum256([]byte(os.Getenv("DB_PASSWORD")))
	/*
	credential := options.Credential{
		AuthMechanism: "SCRAM-SHA-256",
		Username:      os.Getenv("DB_USER"),
		Password:      hex.EncodeToString(sha1Pass[:]),
	}
	*/
	credential := options.Credential{
		Username:      os.Getenv("DB_USER"),
		Password:      os.Getenv("DB_PASSWORD"),
	}

	clientOptions.SetAuth(credential)
	err = clientOptions.Validate()
	if err != nil {
		log.Fatalln(errors.Wrap(err, "failed to validate URI at mongo.Validate()"))
	}
	client, err = mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "failed to make a instance of client at mongo.NewClient()"))
	}

	//ctxMongo, _ = context.WithTimeout(context.Background(), 10*time.Second)
	ctxMongo = context.Background()
	err = client.Connect(ctxMongo)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "failed to connect to DB at mongo.NewClient()"))
	}
	// Check the connection
	err = client.Ping(ctxMongo, nil)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to connect to DB at mongo.NewClient()"))
	}

	db = client.Database(config.Cfg.DB.Name)
	log.Println("Successfully connected to MongoDB!")
}

func Disconnect() error {
	err := client.Disconnect(ctxMongo)
	if err != nil {
		return errors.Wrap(err, "failed to disconnect from DB at Disconnect()")
	}
	return nil
}
