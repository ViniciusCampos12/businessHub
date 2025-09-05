package handlers

import (
	"net/http"

	usecases "github.com/ViniciusCampos12/businessHub/app-golang/internal/application/useCases"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/helpers"
	"github.com/gin-gonic/gin"
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
	companies, err := lc.UseCase.Handle()

	if err != nil {
		helpers.ResponseError(c, err, http.StatusInternalServerError)
		return
	}

	helpers.ResponseSuccess(c, companies, http.StatusOK)
}
