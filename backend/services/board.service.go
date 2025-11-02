package services

import (
	"backend/models"
	"backend/repositories"
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type BoardService struct {
	boardRepository repositories.BoardRepository
}

func NewBoardService(boardRepository *repositories.BoardRepository) *BoardService {
	return &BoardService{boardRepository: *boardRepository}
}

func (b *BoardService) GetAll(ctx context.Context) ([]models.Board, error) {
	boards, err := b.boardRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return boards, nil
}

func (b *BoardService) GetById(ctx context.Context, id string) (*models.Board, error) {
	board, err := b.boardRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return board, nil
}

func (b *BoardService) GetByAlias(ctx context.Context, alias string) (*models.Board, error) {
	board, err := b.boardRepository.GetByAlias(ctx, alias)
	if err != nil {
		return nil, err
	}
	return board, nil
}

func (b *BoardService) CreateBoard(ctx context.Context, name string) (*models.Board, error) {
	newBoard := &models.Board{
		Name:      name,
		CardIDs:   []bson.ObjectID{},
		Alias:     strings.ToLower(name),
		CreatedAt: time.Now().UTC(),
	}

	return b.boardRepository.CreateBoard(ctx, newBoard)
}

func (b *BoardService) AddCard(ctx context.Context, boardId string, cardId string) (*models.Board, error) {
	boardObjectID, err := bson.ObjectIDFromHex(boardId)
	if err != nil {
		return nil, err
	}
	cardObjectID, err := bson.ObjectIDFromHex(cardId)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$addToSet": bson.M{"cardIds": cardObjectID},
	}

	updatedBoard, err := b.boardRepository.UpdateBoard(ctx, boardObjectID, update)
	if err != nil {
		return nil, err
	}

	return updatedBoard, nil
}
