package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ViniciusCampos12/businessHub/app-golang/internal/infra/adapters"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/infra/http/routes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "github.com/ViniciusCampos12/businessHub/app-golang/docs"
)

var mongoDB *mongo.Database
var ctx context.Context
var rbmqPub *adapters.RabbitMqAdapter

func init() {
	mongoDB, ctx = mongoStart()
	rbmqPub = rabbitmqStart()
}

func main() {
	defer rbmqPub.Close()
	defer mongoDB.Client().Disconnect(ctx)
	r := gin.Default()

	r.Use(errorMiddleware())
	routes.SetupBaseRouter(r, mongoDB, rbmqPub)

	r.Run("0.0.0.0:8080")
}

func mongoStart() (*mongo.Database, context.Context) {
	mongoURI, MongoDbDatabase := os.Getenv("MONGO_URI"), os.Getenv("DB_DATABASE")

	if MongoDbDatabase == "" || mongoURI == "" {
		panic("Environment variables MONGO_URI are not defined")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		panic(err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		panic(err)
	}

	log.Println("MongoDB connected successfully")

	return client.Database(MongoDbDatabase), context.Background()
}

func rabbitmqStart() *adapters.RabbitMqAdapter {
	rabbitmqUrl := os.Getenv("RABBITMQ_URL")

	if rabbitmqUrl == "" {
		panic("Environment variables RABBITMQ_URL are not defined")
	}

	return adapters.NewRabbitMqAdapter(rabbitmqUrl)
}

func errorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 && !c.Writer.Written() {
			for _, e := range c.Errors {
				log.Printf("[GIN-ERROR] %v\n", e.Err)
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
	}
}
