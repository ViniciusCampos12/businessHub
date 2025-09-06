package routes

import (
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/infra/adapters"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupBaseRouter(r *gin.Engine, mongo *mongo.Database, rbmqPub *adapters.RabbitMqAdapter) {
	api := r.Group("/api")
	{
		api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		companies := api.Group("/companies")
		{
			SetupCompanyRouter(companies, mongo, rbmqPub)
		}
	}
}
