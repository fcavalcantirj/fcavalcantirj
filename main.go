package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	currentGame *GameState
	logger      *zap.Logger
)

func init() {
	// Initialize logger
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	logger, err = config.Build()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from the frontend
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Create Chi router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))

	// Routes
	r.Post("/api/new-game", handleNewGame)
	r.Put("/api/click", handleClick)
	r.Get("/api/game-state", handleGameState)

	// Add CORS middleware
	handler := corsMiddleware(r)

	// Create server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				logger.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			logger.Fatal("server shutdown failed", zap.Error(err))
		}
		serverStopCtx()
	}()

	// Run the server
	logger.Info("server starting on :8080")
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Fatal("server failed to start", zap.Error(err))
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}

func handleNewGame(w http.ResponseWriter, r *http.Request) {
	selectedStory := stories[rand.Intn(len(stories))]

	currentGame = &GameState{
		Story:       selectedStory.Description,
		StartTime:   time.Now(),
		ClicksLeft:  10,
		Hints:       selectedStory.Hints,
		Solved:      false,
		Solution:    selectedStory.Solution,
		FoundPieces: 0,
	}

	logger.Info("new game started",
		zap.String("story", selectedStory.Description),
		zap.Time("startTime", currentGame.StartTime))

	json.NewEncoder(w).Encode(currentGame)
}

func handleClick(w http.ResponseWriter, r *http.Request) {
	if currentGame == nil {
		http.Error(w, "No active game", http.StatusBadRequest)
		return
	}

	var clickReq ClickRequest
	if err := json.NewDecoder(r.Body).Decode(&clickReq); err != nil {
		logger.Error("failed to decode click request", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	logger.Info("click received",
		zap.String("shapeName", clickReq.ShapeName),
		zap.Float64("cameraX", clickReq.CameraX),
		zap.Float64("cameraY", clickReq.CameraY),
		zap.Float64("cameraZ", clickReq.CameraZ))

	response := processClick(clickReq)
	json.NewEncoder(w).Encode(response)
}

func handleGameState(w http.ResponseWriter, r *http.Request) {
	if currentGame == nil {
		http.Error(w, "No active game", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(currentGame)
}
