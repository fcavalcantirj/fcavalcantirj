package main

import "time"

// GameState represents the current state of the game
type GameState struct {
	Story       string    `json:"story"`
	StartTime   time.Time `json:"startTime"`
	ClicksLeft  int       `json:"clicksLeft"`
	Hints       []string  `json:"hints"`
	Solved      bool      `json:"solved"`
	Solution    []string  `json:"solution"`
	FoundPieces int       `json:"foundPieces"`
}

// ClickRequest represents a click event from the frontend
type ClickRequest struct {
	ShapeName string  `json:"shapeName"`
	CameraX   float64 `json:"cameraX"`
	CameraY   float64 `json:"cameraY"`
	CameraZ   float64 `json:"cameraZ"`
}

// Story represents a game story with its solution and hints
type Story struct {
	Title       string
	Description string
	Solution    []string
	Hints       []string
}

// ShapeInfo contains information about a clickable shape
type ShapeInfo struct {
	Name        string
	Description string
	Location    string
	AngleRange  struct {
		MinX float64
		MaxX float64
		MinY float64
		MaxY float64
	}
}
