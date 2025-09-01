package entities

import (
	"fmt"
	"math"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// swagger:model Company
type Company struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"  swaggerignore:"true"`
	Document          string             `json:"document" bson:"document" validate:"required" example:"19.862.056/0002-23"`
	FantasyName       string             `json:"fantasy_name" bson:"fantasy_name" validate:"required" example:"My Company"`
	SocialReason      string             `json:"social_reason" bson:"social_reason" validate:"required" example:"My Company LTDA"`
	Address           Address            `json:"address" bson:"address" validate:"required,dive"`
	TotalEmployees    int                `json:"total_employees" bson:"total_employees" validate:"gte=0" example:"100"`
	TotalEmployeesPwd int                `json:"total_employees_pwd" bson:"total_employees_pwd" validate:"gte=0" example:"1"`
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at" swaggerignore:"true"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at" swaggerignore:"true"`
}

func (c *Company) pwdPercentage(totalEmployees int) float64 {
	switch {
	case totalEmployees >= 100 && totalEmployees <= 200:
		return 0.02
	case totalEmployees >= 201 && totalEmployees <= 500:
		return 0.03
	case totalEmployees >= 501 && totalEmployees <= 1000:
		return 0.04
	case totalEmployees > 1000:
		return 0.05
	default:
		return 0
	}
}

func (c *Company) CheckPWDQuota(totalEmployees int, currentPWDs int) error {
	percentage := c.pwdPercentage(totalEmployees)
	if percentage == 0 {
		return nil
	}

	minPWDs := int(math.Ceil(float64(totalEmployees) * percentage))

	if currentPWDs < minPWDs {
		return fmt.Errorf("Insufficient quota: company must have %d PWD(s), but has %d (required by Brazilian Law nº 8.213/91, art. 93)", minPWDs, currentPWDs)
	}

	return nil
}

func (c *Company) UnsmaskDocument() {
	document := c.Document
	document = strings.ReplaceAll(document, ".", "")
	document = strings.ReplaceAll(document, "/", "")
	document = strings.ReplaceAll(document, "-", "")

	c.Document = document
}
