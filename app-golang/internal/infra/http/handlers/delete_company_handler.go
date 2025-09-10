package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	usecases "github.com/ViniciusCampos12/businessHub/app-golang/internal/application/useCases"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/fails"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/helpers"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeleteCompany struct {
	UseCase *usecases.DeleteCompany
}

// DeleteComapny godoc
// @Summary delete a existing company
// @Description delete a existing company in businessHub
// @Tags Companies
// @Param id path string true "Company ID"
// @Success 204
// @Failure 404 {object} viewmodels.CompanyNotFoundResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /api/companies/{id} [delete]
func (dc *DeleteCompany) Execute(c *gin.Context) {
	idParam := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		err = fmt.Errorf("%w: %v", fails.ErrValidation, err)
		deleteHandlerError(err, c)
		return
	}

	err = dc.UseCase.Handle(objID, c.Request.Context())

	if err != nil {
		deleteHandlerError(err, c)
		return
	}

	helpers.ResponseSuccess(c, nil, http.StatusNoContent)
}

func deleteHandlerError(err error, c *gin.Context) {
	switch {
	case errors.Is(err, fails.ErrValidation):
		log.Errorf("validation error: %v", err)
		helpers.ResponseError(c, err, http.StatusBadRequest)
	case errors.Is(err, fails.ErrCompanyNotFound):
		helpers.ResponseError(c, err, http.StatusNotFound)
	case errors.Is(err, fails.ErrDbDeleteFailed):
		helpers.ResponseError(c, err, http.StatusNotFound)
	case errors.Is(err, context.Canceled):
		log.Error("request cancelled by client")
		return
	default:
		log.Errorf("unexpected error: %v", err)
		helpers.ResponseError(c, err, http.StatusInternalServerError)
	}
}
