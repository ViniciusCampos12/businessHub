package handlers

import (
	"context"
	"errors"
	"net/http"

	usecases "github.com/ViniciusCampos12/businessHub/app-golang/internal/application/useCases"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/helpers"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type ListCompanies struct {
	UseCase *usecases.ListCompanies
}

// ListCompanies godoc
// @Summary list all companies
// @Description list all companies
// @Tags Companies
// @Produce json
// @Success 200 {object} viewmodels.CompaniesListResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /api/companies [get]
func (lc *ListCompanies) Execute(c *gin.Context) {
	companies, err := lc.UseCase.Handle(c.Request.Context())

	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Error("request cancelled by client")
			return
		}

		log.Errorf("unexpected error: %v", err)
		helpers.ResponseError(c, err, http.StatusInternalServerError)
		return
	}

	helpers.ResponseSuccess(c, companies, http.StatusOK)
}
