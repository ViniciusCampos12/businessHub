package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	usecases "github.com/ViniciusCampos12/businessHub/app-golang/internal/application/useCases"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/entities"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/fails"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/helpers"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EditCompany struct {
	UseCase *usecases.EditCompany
}

// EditCompany godoc
// @Summary edit a existing company
// @Description edit a existing company in businessHub
// @Tags Companies
// @Accept json
// @Produce json
// @Param company body entities.Company true "Company Data"
// @Param id path string true "Company ID"
// @Success 200 {object} viewmodels.CompanyUpdatedResponse
// @Failure 400 {object} viewmodels.CompanyBadRequestResponse
// @Failure 404 {object} viewmodels.CompanyNotFoundResponse
// @Failure 409 {object} viewmodels.CompanyConflitResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /api/companies/{id} [put]
func (e *EditCompany) Execute(c *gin.Context) {
	idParam := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		err = fmt.Errorf("%w: %v", fails.ErrValidation, err)
		editHandlerError(err, c)
		return
	}

	var input entities.Company
	if err := c.ShouldBindJSON(&input); err != nil {
		err = fmt.Errorf("%w: %v", fails.ErrValidation, err)
		editHandlerError(err, c)
		return
	}

	hasUpdated, err := e.UseCase.Handle(objID, &input, c.Request.Context())

	if err != nil {
		editHandlerError(err, c)
		return
	}

	helpers.ResponseSuccess(c, hasUpdated, http.StatusOK)
}

func editHandlerError(err error, c *gin.Context) {
	switch {
	case errors.Is(err, fails.ErrValidation):
		log.Errorf("validation error: %v", err)
		helpers.ResponseError(c, err, http.StatusBadRequest)
	case errors.Is(err, fails.ErrCompanyNotFound):
		helpers.ResponseError(c, err, http.StatusNotFound)
	case errors.Is(err, fails.ErrDbUpdateFailed):
		helpers.ResponseError(c, err, http.StatusNotFound)
	case errors.Is(err, fails.ErrCompanyDocumentIsAlreadyInUse):
		helpers.ResponseError(c, err, http.StatusConflict)
	case errors.Is(err, fails.ErrInsufficientPWDQuota):
		helpers.ResponseError(c, err, http.StatusBadRequest)
	case errors.Is(err, context.Canceled):
		log.Error("request cancelled by client")
		return
	default:
		log.Errorf("unexpected error: %v", err)
		helpers.ResponseError(c, err, http.StatusInternalServerError)
	}
}
