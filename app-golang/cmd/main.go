package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ViniciusCampos12/businessHub/app-golang/internal/infra/adapters"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/infra/http/routes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "github.com/ViniciusCampos12/businessHub/app-golang/docs"
	log "github.com/sirupsen/logrus"
)

var mongoDB *mongo.Database
var rbmqPub *adapters.RabbitMqAdapter

func init() {
	mongoDB = mongoStart()
	rbmqPub = rabbitmqStart()
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	r := gin.New()
	r.Use(LogrusMiddleware(), gin.Recovery())

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	r.Use(ErrorMiddleware())
	routes.SetupBaseRouter(r, mongoDB, rbmqPub)

	go shutdown(server)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe: %v", err)
	}
}

func LogrusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)

		log.WithFields(log.Fields{
			"status":   c.Writer.Status(),
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"latency":  latency.String(),
			"clientIP": c.ClientIP(),
		}).Info("HTTP request")
	}
}

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 && !c.Writer.Written() {
			for _, e := range c.Errors {
				log.Errorf("[GIN-ERROR] %v\n", e.Err)
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
	}
}

func shutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Infof("Server forced to shutdown: %v", err)
	}

	cleanup(ctx, mongoDB, rbmqPub)
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
