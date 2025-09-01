package routes

import (
	usecases "github.com/ViniciusCampos12/businessHub/app-golang/internal/application/useCases"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/infra/adapters"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/infra/database/repositories"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/infra/http/handlers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupCompanyRouter(rg *gin.RouterGroup, mongo *mongo.Database, rbmqPub *adapters.RabbitMqAdapter) {
	createCompanyHandler := createCompanyHandlerFactory(mongo, rbmqPub)
	listCompaniesHandler := listCompaniesHandlerFactory(mongo, rbmqPub)
	editCompanyHandler := editCompanyHandlerFactory(mongo, rbmqPub)
	deleteCompanyHandler := deleteCompanyHandlerFactory(mongo, rbmqPub)

	companies := rg.Group("/companies")
	{
		companies.GET("/", listCompaniesHandler.Execute)
		companies.POST("/", createCompanyHandler.Execute)
		companies.PUT("/:id", editCompanyHandler.Execute)
		companies.DELETE("/:id", deleteCompanyHandler.Execute)
	}
}

func createCompanyHandlerFactory(mongo *mongo.Database, rbmqPub *adapters.RabbitMqAdapter) *handlers.CreateCompany {
	companyRepository := &repositories.CompanyRepository{
		Mongo: mongo,
	}
	createCompanyUseCase := &usecases.CreateCompany{
		Repo:   companyRepository,
		Broker: rbmqPub,
	}
	createCompanyHandler := &handlers.CreateCompany{
		UseCase: createCompanyUseCase,
	}

	return createCompanyHandler
}

func listCompaniesHandlerFactory(mongo *mongo.Database, rbmqPub *adapters.RabbitMqAdapter) *handlers.ListCompanies {
	companyRepository := &repositories.CompanyRepository{
		Mongo: mongo,
	}
	listCompaniesUseCase := &usecases.ListCompanies{
		Repo:   companyRepository,
		Broker: rbmqPub,
	}
	listCompaniesHandler := &handlers.ListCompanies{
		UseCase: listCompaniesUseCase,
	}

	return listCompaniesHandler
}

func editCompanyHandlerFactory(mongo *mongo.Database, rbmqPub *adapters.RabbitMqAdapter) *handlers.EditCompany {
	companyRepository := &repositories.CompanyRepository{
		Mongo: mongo,
	}
	editCompanyUseCase := &usecases.EditCompany{
		Repo:   companyRepository,
		Broker: rbmqPub,
	}
	editCompanyHandler := &handlers.EditCompany{
		UseCase: editCompanyUseCase,
	}

	return editCompanyHandler
}

func deleteCompanyHandlerFactory(mongo *mongo.Database, rbmqPub *adapters.RabbitMqAdapter) *handlers.DeleteCompany {
	companyRepository := &repositories.CompanyRepository{
		Mongo: mongo,
	}
	deleteCompanyUseCase := &usecases.DeleteCompany{
		Repo:   companyRepository,
		Broker: rbmqPub,
	}
	deleteCompanyHandler := &handlers.DeleteCompany{
		UseCase: deleteCompanyUseCase,
	}

	return deleteCompanyHandler
}
