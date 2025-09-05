package handlers

import (
	"net/http"
	"strings"

	usecases "github.com/ViniciusCampos12/businessHub/app-golang/internal/application/useCases"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/entities"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/helpers"
	"github.com/gin-gonic/gin"
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
// @Failure 500 {object} helpers.ErrorResponse
// @Router /api/companies/{id} [put]
func (e *EditCompany) Execute(c *gin.Context) {
	id := c.Param("id")

	var input entities.Company
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.ResponseError(c, err, http.StatusBadRequest)
		return
	}

	hasUpdated, err := e.UseCase.Handle(id, &input)
	if err != nil {
		if err.Error() == "company not found" {
			helpers.ResponseError(c, err, http.StatusNotFound)
			return
		}

		if strings.HasPrefix(err.Error(), "Insufficient quota: company must have") {
			helpers.ResponseError(c, err, http.StatusBadRequest)
			return
		}

		helpers.ResponseError(c, err, http.StatusInternalServerError)
		return
	}

	helpers.ResponseSuccess(c, hasUpdated, http.StatusOK)
}
