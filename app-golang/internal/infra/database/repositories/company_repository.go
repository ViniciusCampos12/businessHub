package repositories

import (
	"context"
	"errors"

	"github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CompanyRepository struct {
	Mongo *mongo.Database
}

func (r *CompanyRepository) Create(c *entities.Company) (*entities.Company, error) {
	_, err := r.Mongo.Collection("companies").InsertOne(context.Background(), c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *CompanyRepository) FindByDocument(document string) (*entities.Company, error) {
	var company entities.Company
	err := r.Mongo.Collection("companies").FindOne(context.Background(), map[string]interface{}{
		"document": document,
	},
	).Decode(&company)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &company, nil
}

func (r *CompanyRepository) FindMany() ([]*entities.Company, error) {
	cursor, err := r.Mongo.Collection("companies").Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var companies []*entities.Company
	if err = cursor.All(context.TODO(), &companies); err != nil {
		return nil, err
	}

	return companies, nil
}

func (r *CompanyRepository) FindById(id string) (*entities.Company, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var company entities.Company
	err = r.Mongo.Collection("companies").FindOne(context.TODO(), bson.M{
		"_id": objID,
	},
	).Decode(&company)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &company, nil
}

func (r *CompanyRepository) Save(id string, c *entities.Company) (bool, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	update := bson.M{
		"$set": bson.M{
			"document":            c.Document,
			"fantasy_name":        c.FantasyName,
			"social_reason":       c.SocialReason,
			"address":             c.Address,
			"total_employees":     c.TotalEmployees,
			"total_employees_pwd": c.TotalEmployeesPwd,
			"updated_at":          c.UpdatedAt,
		},
	}

	result, err := r.Mongo.Collection("companies").UpdateOne(context.TODO(), bson.M{
		"_id": objID,
	}, update)
	if err != nil {
		return false, err
	}

	if result.MatchedCount == 0 {
		return false, nil
	}

	return true, nil
}

func (r *CompanyRepository) Delete(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ObjectID")
	}

	res, err := r.Mongo.Collection("companies").DeleteOne(context.TODO(), bson.M{
		"_id": objID,
	},
	)

	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("delete failed")
	}

	return nil
}
