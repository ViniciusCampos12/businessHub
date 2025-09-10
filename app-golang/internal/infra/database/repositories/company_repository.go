package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/ViniciusCampos12/businessHub/app-golang/internal/domain/entities"
	"github.com/ViniciusCampos12/businessHub/app-golang/internal/fails"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CompanyRepository struct {
	Mongo *mongo.Database
	db    *mongo.Collection
}

func NewCompanyRepository(db *mongo.Database) (*CompanyRepository, error) {
	repo := &CompanyRepository{
		Mongo: db,
	}

	repo.db = repo.Mongo.Collection("companies")

	if err := repo.ensureIndexes(); err != nil {
		return nil, fmt.Errorf("mongodb fail to ensure indexes: %w", err)
	}

	return repo, nil
}

func (r *CompanyRepository) ensureIndexes() error {
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"document": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := r.db.Indexes().CreateOne(context.Background(), indexModel)
	return err
}

func (r *CompanyRepository) Create(c *entities.Company, ctx context.Context) (*entities.Company, error) {
	_, err := r.db.InsertOne(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("mongodb fail to create: %w", err)
	}
	return c, nil
}

func (r *CompanyRepository) FindByDocument(document string, ctx context.Context) (*entities.Company, error) {
	var company entities.Company
	err := r.db.FindOne(ctx, bson.M{
		"document": document,
	},
	).Decode(&company)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("mongodb fail to find by document: %w", err)
	}

	return &company, nil
}

func (r *CompanyRepository) FindMany(ctx context.Context) ([]*entities.Company, error) {
	cursor, err := r.db.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("mongodb fail to find: %w", err)
	}
	defer cursor.Close(ctx)

	var companies []*entities.Company
	if err = cursor.All(ctx, &companies); err != nil {
		return nil, fmt.Errorf("mongodb fail to iteration with cursor all: %w", err)
	}

	return companies, nil
}

func (r *CompanyRepository) FindById(objID primitive.ObjectID, ctx context.Context) (*entities.Company, error) {
	var company entities.Company
	err := r.db.FindOne(ctx, bson.M{
		"_id": objID,
	},
	).Decode(&company)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("mongodb fail to find by id: %w", err)
	}

	return &company, nil
}

func (r *CompanyRepository) Save(objID primitive.ObjectID, c *entities.Company, ctx context.Context) error {
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

	result, err := r.db.UpdateOne(ctx, bson.M{
		"_id": objID,
	}, update)
	if err != nil {
		return fmt.Errorf("mongodb fail to update one: %w", err)
	}

	if result.MatchedCount == 0 {
		return fails.ErrDbUpdateFailed
	}

	return nil
}

func (r *CompanyRepository) Delete(objID primitive.ObjectID, ctx context.Context) error {
	res, err := r.db.DeleteOne(ctx, bson.M{
		"_id": objID,
	},
	)

	if err != nil {
		return fmt.Errorf("mongodb fail to delete one: %w", err)
	}

	if res.DeletedCount == 0 {
		return fails.ErrDbDeleteFailed
	}

	return nil
}
