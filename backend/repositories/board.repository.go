package repositories

import (
	"backend/models"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BoardRepository struct {
	db              *mongo.Database
	boardCollection *mongo.Collection
}

func NewBoardRepository(db *mongo.Database) *BoardRepository {
	return &BoardRepository{
		db:              db,
		boardCollection: db.Collection("boards"),
	}
}

func (b *BoardRepository) GetAll(ctx context.Context) ([]models.Board, error) {
	var boards []models.Board

	cursor, err := b.boardCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error finding boards: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &boards); err != nil {
		log.Printf("Error decoding boards: %v", err)
		return nil, err
	}

	return boards, nil
}

func (b *BoardRepository) GetById(ctx context.Context, id string) (*models.Board, error) {
	var board models.Board

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid ObjectID: %v", err)
		return nil, err
	}

	err = b.boardCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&board)
	if err != nil {
		log.Printf("Error finding board with ID %s: %v", id, err)
		return nil, err
	}

	return &board, nil
}

func (b *BoardRepository) GetByAlias(ctx context.Context, alias string) (*models.Board, error) {
	var board models.Board

	err := b.boardCollection.FindOne(ctx, bson.M{"alias": alias}).Decode(&board)
	if err != nil {
		log.Printf("Error finding board with alias %s: %v", alias, err)
		return nil, err
	}

	return &board, nil
}

func (b *BoardRepository) CreateBoard(ctx context.Context, board *models.Board) (*models.Board, error) {
	result, err := b.boardCollection.InsertOne(ctx, board)
	if err != nil {
		return nil, err
	}

	board.ID = result.InsertedID.(bson.ObjectID)
	return board, nil
}

func (b *BoardRepository) UpdateBoard(ctx context.Context, boardID bson.ObjectID, update bson.M) (*models.Board, error) {
	_, err := b.boardCollection.UpdateOne(ctx, bson.M{"_id": boardID}, update)
	if err != nil {
		return nil, err
	}

	var updated models.Board
	if err := b.boardCollection.FindOne(ctx, bson.M{"_id": boardID}).Decode(&updated); err != nil {
		return nil, err
	}

	return &updated, nil
}
