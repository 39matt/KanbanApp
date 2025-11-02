package services

import (
	"backend/models"
	"backend/repositories"
	"context"
	"time"
)

type CardService struct {
	cardRepository repositories.CardRepository
}

func NewCardService(cardRepository *repositories.CardRepository) *CardService {
	return &CardService{cardRepository: *cardRepository}
}

func (c *CardService) GetAll(ctx context.Context) ([]models.Card, error) {
	cards, err := c.cardRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (c *CardService) GetById(ctx context.Context, id string) (*models.Card, error) {
	card, err := c.cardRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return card, nil
}

func (c *CardService) Create(ctx context.Context, title string, description string, section string) (*models.Card, error) {
	newCard := &models.Card{
		Title:       title,
		Description: description,
		Section:     section,
		CreatedAt:   time.Now().UTC(),
	}

	return c.cardRepository.Create(ctx, newCard)
}
