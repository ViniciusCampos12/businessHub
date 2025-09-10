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
)

type CreateCompany struct {
	UseCase *usecases.CreateCompany
}

// CreateCompany godoc
// @Summary create a new company
// @Description Create a new company in businessHub
// @Tags Companies
// @Accept json
// @Produce json
// @Param company body entities.Company true "Company Data"
// @Success 201 {object} viewmodels.CompanyCreatedResponse
// @Failure 400 {object} viewmodels.CompanyBadRequestResponse
// @Failure 409 {object} viewmodels.CompanyConflitResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /api/companies [post]
func (h *CreateCompany) Execute(c *gin.Context) {
	var input entities.Company
	if err := c.ShouldBindJSON(&input); err != nil {
		err = fmt.Errorf("%w: %v", fails.ErrValidation, err)
		createHandlerError(err, c)
		return
	}

	newCompany, err := h.UseCase.Handle(&input, c.Request.Context())

	if err != nil {
		createHandlerError(err, c)
		return
	}

	helpers.ResponseSuccess(c, newCompany, http.StatusCreated)
}

func createHandlerError(err error, c *gin.Context) {
	switch {
	case errors.Is(err, fails.ErrValidation):
		log.Errorf("validation error: %v", err)
		helpers.ResponseError(c, err, http.StatusBadRequest)
	case errors.Is(err, fails.ErrCompanyAlreadyExists):
		helpers.ResponseError(c, err, http.StatusConflict)
	case errors.Is(err, fails.ErrInsufficientPWDQuota):
		helpers.ResponseError(c, err, http.StatusBadRequest)
	case errors.Is(err, context.Canceled):
		log.Errorf("request cancelled by client")
		return
	default:
		log.Errorf("unexpected error: %v", err)
		helpers.ResponseError(c, err, http.StatusInternalServerError)
	}
}
