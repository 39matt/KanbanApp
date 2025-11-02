package main

import (
	"backend/models"
	"backend/repositories"
	"backend/services"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Server struct {
	db           *mongo.Database
	boardService *services.BoardService
	cardService  *services.CardService
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := models.InitServer(ctx)

	boardRepo := repositories.NewBoardRepository(db)
	boardService := services.NewBoardService(boardRepo)
	cardRepo := repositories.NewCardRepository(db)
	cardService := services.NewCardService(cardRepo)
	server := &Server{db, boardService, cardService}

	router := gin.Default()
	router.Use(corsMiddleware())
	router.GET("/boards/get-all", server.GetBoardsHandler)
	router.POST("/boards/get-by-id", server.GetBoardByIdHandler)
	router.POST("/boards/:alias", server.GetBoardByAliasHandler)
	router.POST("/boards/add", server.AddBoardHandler)
	router.POST("/boards/add-card-to-board", server.AddCardToBoardHandler)

	router.GET("/cards/get-all", server.GetCardsHandler)
	router.POST("/cards/get-by-id", server.GetCardByIdHandler)
	router.POST("/cards/add", server.AddCardHandler)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func (s *Server) GetBoardsHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	boards, err := s.boardService.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"boards": boards})
}

type GetBoardRequest struct {
	Id string `json:"id"`
}

func (s *Server) GetBoardByIdHandler(c *gin.Context) {
	var req GetBoardRequest
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		http.Error(c.Writer, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	board, err := s.boardService.GetById(ctx, req.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"board": board})
}

func (s *Server) GetBoardByAliasHandler(c *gin.Context) {
	alias := c.Param("alias")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	board, err := s.boardService.GetByAlias(ctx, alias)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"board": board})
}

type AddBoardRequest struct {
	Name string `json:"name"`
}

func (s *Server) AddBoardHandler(c *gin.Context) {
	var req AddBoardRequest
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		http.Error(c.Writer, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	board, err := s.boardService.CreateBoard(ctx, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"board": board})
}

func (s *Server) GetCardsHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	cards, err := s.cardService.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cards": cards})
}

type GetCardRequest struct {
	Id string `json:"id"`
}

func (s *Server) GetCardByIdHandler(c *gin.Context) {
	var req GetCardRequest
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		http.Error(c.Writer, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	card, err := s.cardService.GetById(ctx, req.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"card": card})
}

type AddCardRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Section     string `json:"section"`
}

func (s *Server) AddCardHandler(c *gin.Context) {
	var req AddCardRequest
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		http.Error(c.Writer, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	card, err := s.cardService.Create(ctx, req.Title, req.Description, req.Section)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"card": card})
}

type AddCardToBoardRequest struct {
	BoardId string `json:"boardId"`
	CardId  string `json:"CardId"`
}

func (s *Server) AddCardToBoardHandler(c *gin.Context) {
	var req AddCardToBoardRequest
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		http.Error(c.Writer, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	board, err := s.boardService.AddCard(ctx, req.BoardId, req.CardId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"board": board})
}
