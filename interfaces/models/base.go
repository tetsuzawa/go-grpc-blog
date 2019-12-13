package models

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tetsuzawa/go-3daudio/web-app/config"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	tableNameBlogData = "blog"
)

//const tFormat = "2006-01-02 15:04:05"
const tFormat = "2006-01-02T15:04:05.000Z"

func GetTableName(name string) string {
	return fmt.Sprintf("%s", name)
}

var client *mongo.Client
var db *mongo.Database

func init() {
	var err error
	err = godotenv.Load() //Load env.file
	if err != nil {
		log.Fatalln(errors.Wrap(err, "failed to load .env file at godotenv.Load()"))
	}

	dbName := fmt.Sprintf("mongodb://%s:%s@%s:%d",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		config.Cfg.DB.Host,
		config.Cfg.DB.Port,
		//config.Cfg.DB.Name,
		//config.Cfg.DB.ETC,
	)

	clientOptions := options.Client().ApplyURI(dbName)
	client, err = mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "failed to make a instance of client at mongo.NewClient()"))
	}

	ctx, _ := context.WithTimeout(context.TODO(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "failed to connect to DB at mongo.NewClient()"))
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to connect to DB at mongo.NewClient()"))
	}

	db = client.Database(config.Cfg.DB.Name)
}
