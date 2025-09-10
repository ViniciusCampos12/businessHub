package routes

import (
	"fmt"
	"log"

	usecases "github.com/ViniciusCampos12/businessHub/app-golang/internal/application/useCases"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/infra/adapters"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/infra/database/repositories"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/infra/http/handlers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupCompanyRouter(companies *gin.RouterGroup, mongo *mongo.Database, rbmqPub *adapters.RabbitMqAdapter) {
	companyRepo, err := loadCompanyRepo(mongo, repositories.NewCompanyRepository)

	if err != nil {
		log.Fatalf("failed to load repository: %v", err)
	}

	createCompanyHandler := createCompanyHandlerFactory(companyRepo, rbmqPub)
	listCompaniesHandler := listCompaniesHandlerFactory(companyRepo, rbmqPub)
	editCompanyHandler := editCompanyHandlerFactory(companyRepo, rbmqPub)
	deleteCompanyHandler := deleteCompanyHandlerFactory(companyRepo, rbmqPub)

	companies.GET("/", listCompaniesHandler.Execute)
	companies.POST("/", createCompanyHandler.Execute)
	companies.PUT("/:id", editCompanyHandler.Execute)
	companies.DELETE("/:id", deleteCompanyHandler.Execute)
}

func createCompanyHandlerFactory(repo *repositories.CompanyRepository, rbmqPub *adapters.RabbitMqAdapter) *handlers.CreateCompany {
	createCompanyUseCase := &usecases.CreateCompany{
		Repo:   repo,
		Broker: rbmqPub,
	}

	createCompanyHandler := &handlers.CreateCompany{
		UseCase: createCompanyUseCase,
	}

	return createCompanyHandler
}

func listCompaniesHandlerFactory(repo *repositories.CompanyRepository, rbmqPub *adapters.RabbitMqAdapter) *handlers.ListCompanies {
	listCompaniesUseCase := &usecases.ListCompanies{
		Repo:   repo,
		Broker: rbmqPub,
	}
	listCompaniesHandler := &handlers.ListCompanies{
		UseCase: listCompaniesUseCase,
	}

	return listCompaniesHandler
}

func editCompanyHandlerFactory(repo *repositories.CompanyRepository, rbmqPub *adapters.RabbitMqAdapter) *handlers.EditCompany {
	editCompanyUseCase := &usecases.EditCompany{
		Repo:   repo,
		Broker: rbmqPub,
	}
	editCompanyHandler := &handlers.EditCompany{
		UseCase: editCompanyUseCase,
	}

	return editCompanyHandler
}

func deleteCompanyHandlerFactory(repo *repositories.CompanyRepository, rbmqPub *adapters.RabbitMqAdapter) *handlers.DeleteCompany {
	deleteCompanyUseCase := &usecases.DeleteCompany{
		Repo:   repo,
		Broker: rbmqPub,
	}
	deleteCompanyHandler := &handlers.DeleteCompany{
		UseCase: deleteCompanyUseCase,
	}

	return deleteCompanyHandler
}

func loadCompanyRepo(mongo *mongo.Database, newRepo func(db *mongo.Database) (*repositories.CompanyRepository, error)) (*repositories.CompanyRepository, error) {
	repo, err := newRepo(mongo)

	if err != nil {
		return nil, fmt.Errorf("fail to load repo in handlers: %w", err)
	}

	return repo, nil
}
