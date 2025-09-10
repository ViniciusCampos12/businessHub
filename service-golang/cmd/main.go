package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ViniciusCampos12/businessHub/service-golang/internal/infra/adapters"
	"github.com/ViniciusCampos12/businessHub/service-golang/internal/infra/jobs"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoDB *mongo.Database
var rbmqAdpt *adapters.RabbitMqAdapter

func init() {
	mongoDB = mongoStart()
	rbmqAdpt = rabbitmqStart()
	startJobs(mongoDB)
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.Info("start app")
	shutdown()
}

func shutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cleanup(ctx, mongoDB, rbmqAdpt)
}

func mongoStart() *mongo.Database {
	mongoURI, MongoDbDatabase := os.Getenv("MONGO_URI"), os.Getenv("MONGO_DATABASE")

	if MongoDbDatabase == "" || mongoURI == "" {
		log.Fatal("Environment variables MONGO_URI are not defined")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	log.Info("MongoDB connected successfully")

	return client.Database(MongoDbDatabase)
}

func rabbitmqStart() *adapters.RabbitMqAdapter {
	rabbitmqUrl := os.Getenv("RABBITMQ_URL")

	if rabbitmqUrl == "" {
		log.Fatal("Environment variables RABBITMQ_URL are not defined")
	}

	return adapters.NewRabbitMqAdapter(rabbitmqUrl)
}

func cleanup(ctx context.Context, mongoDB *mongo.Database, rbmqPub *adapters.RabbitMqAdapter) {
	if err := mongoDB.Client().Disconnect(ctx); err != nil {
		log.Errorf("Error disconnecting Mongo: %v", err)
	}

	log.Info("Mongo successfully disconnected")

	if err := rbmqPub.Close(); err != nil {
		log.Errorf("Error disconnecting RabbitMQ: %v", err)
	}

	log.Info("RabbitMQ successfully disconnected")
}

func startJobs(mongoDB *mongo.Database) {
	rbQueue := os.Getenv("RABBITMQ_COMPANY_QUEUE")

	if rbQueue == "" {
		log.Fatal("Environment variables RABBITMQ_COMPANY_QUEUE are not defined")
	}

	jobs.StartCompanyConsumer(rbQueue, rbmqAdpt, mongoDB)
}
