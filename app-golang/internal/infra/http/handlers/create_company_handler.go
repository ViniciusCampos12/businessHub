package handlers

import (
	"log"
	"net/http"
	"strings"

	usecases "github.com/ViniciusCampos12/businessHub/app-golang/internal/application/useCases"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/entities"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/helpers"
	"github.com/gin-gonic/gin"
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
		helpers.ResponseError(c, err, http.StatusBadRequest)
		return
	}

	newCompany, err := h.UseCase.Handle(&input)

	if err != nil {
		if err.Error() == "company already exists" {
			helpers.ResponseError(c, err, http.StatusConflict)
			return
		}

		if strings.HasPrefix(err.Error(), "Insufficient quota: company must have") {
			helpers.ResponseError(c, err, http.StatusBadRequest)
			return
		}

		log.Printf("%v", err.Error())
		helpers.ResponseError(c, nil, http.StatusInternalServerError)
		return
	}

	helpers.ResponseSuccess(c, newCompany, http.StatusCreated)
}
