package handlers

import (
	"net/http"

	usecases "github.com/ViniciusCampos12/businessHub/app-golang/internal/application/useCases"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/helpers"
	"github.com/gin-gonic/gin"
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
	id := c.Param("id")
	err := dc.UseCase.Handle(id)

	if err != nil {
		if err.Error() == "company not found" {
			helpers.ResponseError(c, err, http.StatusNotFound)
			return
		}
		helpers.ResponseError(c, err, http.StatusInternalServerError)
		return
	}

	helpers.ResponseSuccess(c, nil, http.StatusNoContent)
}
