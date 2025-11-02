package repositories

import (
	"backend/models"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CardRepository struct {
	db             *mongo.Database
	cardCollection *mongo.Collection
}

func NewCardRepository(db *mongo.Database) *CardRepository {
	return &CardRepository{
		db:             db,
		cardCollection: db.Collection("cards"),
	}
}

func (c *CardRepository) GetAll(ctx context.Context) ([]models.Card, error) {
	var cards []models.Card

	cursor, err := c.cardCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error finding cards: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &cards); err != nil {
		log.Printf("Error decoding cards: %v", err)
		return nil, err
	}

	return cards, nil
}

func (c *CardRepository) GetById(ctx context.Context, id string) (*models.Card, error) {
	var card models.Card

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid ObjectID: %v", err)
		return nil, err
	}

	err = c.cardCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&card)
	if err != nil {
		log.Printf("Error finding card with ID %s: %v", id, err)
		return nil, err
	}

	return &card, nil
}

func (c *CardRepository) Create(ctx context.Context, card *models.Card) (*models.Card, error) {
	result, err := c.cardCollection.InsertOne(ctx, card)
	if err != nil {
		return nil, err
	}

	card.Id = result.InsertedID.(bson.ObjectID)
	return card, nil
}
