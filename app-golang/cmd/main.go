package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ViniciusCampos12/businessHub/app-golang/internal/infra/adapters"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/infra/http/routes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "github.com/ViniciusCampos12/businessHub/app-golang/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func initMongo() *mongo.Database {
	mongoURI, MongoDbDatabase := os.Getenv("MONGO_URI"), os.Getenv("DB_DATABASE")

	if MongoDbDatabase == "" || mongoURI == "" {
		log.Fatal("Environment variables MONGO_URI are not defined")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB is not responding:", err)
	}

	log.Println("MongoDB connected successfully")

	return client.Database(MongoDbDatabase)
}

func initRabbitMq() *adapters.RabbitMqAdapter {
	rabbitmqUrl := os.Getenv("RABBITMQ_URL")

	if rabbitmqUrl == "" {
		log.Fatal("Environment variables RABBITMQ_URL are not defined")
	}

	publisher := adapters.NewRabbitMqAdapter(rabbitmqUrl)

	return publisher
}

func main() {
	mongo := initMongo()
	rbmqPub := initRabbitMq()
	defer rbmqPub.Close()
	r := gin.Default()

	api := r.Group("/api")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.SetupCompanyRouter(api, mongo, rbmqPub)

	r.Run("0.0.0.0:8080")
}
