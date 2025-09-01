package handlers

import (
	"net/http"

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
// @Success 201 {object} helpers.SuccessResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 409 {object} helpers.ErrorResponse
// @Router /api/companies [post]
func (h *CreateCompany) Execute(c *gin.Context) {
	var input *entities.Company
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.ResponseError(c, err, http.StatusBadRequest)
		return
	}

	newCompany, err := h.UseCase.Handle(input)
	
	if err != nil {
		helpers.ResponseError(c, err, http.StatusConflict)
		return
	}

	helpers.ResponseSuccess(c, newCompany, http.StatusCreated)
}
