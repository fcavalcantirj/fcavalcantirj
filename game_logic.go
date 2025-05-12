package main

import (
	"math/rand"
	"time"

	"go.uber.org/zap"
)

var stories = []Story{
	{
		Title:       "The Missing Key",
		Description: "A valuable key has gone missing from the Little Tokyo bank. The last person who saw it was the bank manager, who claims to have left it on his desk. However, security footage shows someone entering the building after hours. Find the key and discover who took it.",
		Solution:    []string{"Object649_alpha_glass_0", "Mesh_0_Body_0"},
		Hints: []string{
			"Look for a window that might have been used as an entry point",
			"The thief might have used a vehicle to escape",
			"Check the building's exterior for any signs of forced entry",
			"The key might be hidden near where the thief entered",
			"Look for any unusual objects near the building",
		},
	},
	{
		Title:       "The Mysterious Package",
		Description: "A suspicious package was delivered to the main building. Security is concerned it might contain something dangerous. The delivery person was seen acting strangely. Investigate the package and find out what's inside.",
		Solution:    []string{"Mesh_1_Windows_0", "Object649_alpha_glass_0"},
		Hints: []string{
			"The package was left near a window",
			"Look for any unusual objects near the building",
			"Check the windows for any signs of tampering",
			"The delivery person might have left something behind",
			"Look for anything that doesn't belong in the scene",
		},
	},
}

var shapeInfo = map[string]ShapeInfo{
	"Object649_alpha_glass_0": {
		Name:        "Building Window",
		Description: "A large window on the main building. There are signs of forced entry!",
		Location:    "Main building facade",
		AngleRange: struct {
			MinX float64
			MaxX float64
			MinY float64
			MaxY float64
		}{
			MinX: 0,
			MaxX: 10,
			MinY: 0,
			MaxY: 5,
		},
	},
	"Object649_paintmat_0": {
		Name:        "Building Wall",
		Description: "The main wall of the building. Look for any unusual marks or damage.",
		Location:    "Main building facade",
		AngleRange: struct {
			MinX float64
			MaxX float64
			MinY float64
			MaxY float64
		}{
			MinX: 0,
			MaxX: 10,
			MinY: 0,
			MaxY: 5,
		},
	},
	"Mesh_0_Body_0": {
		Name:        "Suspicious Vehicle",
		Description: "A vehicle parked near the building. The driver seems to be watching the entrance.",
		Location:    "Street level",
		AngleRange: struct {
			MinX float64
			MaxX float64
			MinY float64
			MaxY float64
		}{
			MinX: 0,
			MaxX: 10,
			MinY: 0,
			MaxY: 5,
		},
	},
	"Mesh_1_Windows_0": {
		Name:        "Vehicle Windows",
		Description: "The windows of a suspicious vehicle. Something seems off...",
		Location:    "Street level",
		AngleRange: struct {
			MinX float64
			MaxX float64
			MinY float64
			MaxY float64
		}{
			MinX: 0,
			MaxX: 10,
			MinY: 0,
			MaxY: 5,
		},
	},
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateRandomStory() string {
	story := stories[rand.Intn(len(stories))]
	logger.Info("generating new story",
		zap.String("title", story.Title),
		zap.Strings("solution", story.Solution))
	return story.Description
}

func processClick(req ClickRequest) map[string]interface{} {
	if currentGame == nil {
		return map[string]interface{}{
			"clicksLeft": 0,
			"hint":       "No active game. Please start a new game.",
			"solved":     false,
		}
	}

	// Check if we have clicks left
	if currentGame.ClicksLeft <= 0 {
		return map[string]interface{}{
			"clicksLeft": 0,
			"hint":       "You've used all your clicks! Start a new game to continue investigating.",
			"solved":     false,
		}
	}

	// Decrement clicks first
	currentGame.ClicksLeft--

	shape, exists := shapeInfo[req.ShapeName]
	if !exists {
		logger.Info("clicked unknown shape",
			zap.String("shapeName", req.ShapeName),
			zap.Float64("cameraX", req.CameraX),
			zap.Float64("cameraY", req.CameraY),
			zap.Float64("cameraZ", req.CameraZ))

		// Generate a random message from a list of helpful hints
		messages := []string{
			"Focus on the main building and any suspicious vehicles nearby.",
			"Look for signs of forced entry or unusual activity around the building.",
			"Check the windows and walls of the main building for clues.",
			"Pay attention to any vehicles that might have been used in the crime.",
			"Examine the building's exterior for any signs of tampering.",
			"Look for anything that seems out of place in the scene.",
		}
		randomMessage := messages[rand.Intn(len(messages))]

		return map[string]interface{}{
			"clicksLeft": currentGame.ClicksLeft,
			"hint":       randomMessage,
			"solved":     false,
		}
	}

	// Check if the camera angle is within the correct range for this shape
	if isCorrectAngle(req.CameraX, req.CameraY, shape.AngleRange) {
		logger.Info("clicked shape from correct angle",
			zap.String("shapeName", req.ShapeName),
			zap.String("shapeDescription", shape.Description),
			zap.Float64("cameraX", req.CameraX),
			zap.Float64("cameraY", req.CameraY),
			zap.Float64("cameraZ", req.CameraZ),
			zap.Float64("angleRangeMinX", shape.AngleRange.MinX),
			zap.Float64("angleRangeMaxX", shape.AngleRange.MaxX),
			zap.Float64("angleRangeMinY", shape.AngleRange.MinY),
			zap.Float64("angleRangeMaxY", shape.AngleRange.MaxY))

		// Check if this is part of the solution
		for _, solutionShape := range currentGame.Solution {
			if req.ShapeName == solutionShape {
				// Found a solution piece!
				currentGame.FoundPieces++
				if currentGame.FoundPieces == len(currentGame.Solution) {
					// All pieces found!
					return map[string]interface{}{
						"clicksLeft": currentGame.ClicksLeft,
						"hint":       "You found a crucial clue! " + shape.Description,
						"solved":     true,
					}
				}
				return map[string]interface{}{
					"clicksLeft": currentGame.ClicksLeft,
					"hint":       "You found a crucial clue! " + shape.Description + "\n\nYou're getting closer to solving the mystery!",
					"solved":     false,
				}
			}
		}

		return map[string]interface{}{
			"clicksLeft": currentGame.ClicksLeft,
			"hint":       "You found something interesting! " + shape.Description,
			"solved":     false,
		}
	}

	logger.Info("clicked shape from wrong angle",
		zap.String("shapeName", req.ShapeName),
		zap.String("shapeDescription", shape.Description),
		zap.Float64("cameraX", req.CameraX),
		zap.Float64("cameraY", req.CameraY),
		zap.Float64("cameraZ", req.CameraZ),
		zap.Float64("angleRangeMinX", shape.AngleRange.MinX),
		zap.Float64("angleRangeMaxX", shape.AngleRange.MaxX),
		zap.Float64("angleRangeMinY", shape.AngleRange.MinY),
		zap.Float64("angleRangeMaxY", shape.AngleRange.MaxY))

	return map[string]interface{}{
		"clicksLeft": currentGame.ClicksLeft,
		"hint":       "From this angle, you can't see anything relevant. Try moving around to get a better view.",
		"solved":     false,
	}
}

func isCorrectAngle(x, y float64, angleRange struct {
	MinX float64
	MaxX float64
	MinY float64
	MaxY float64
}) bool {
	return x >= angleRange.MinX && x <= angleRange.MaxX && y >= angleRange.MinY && y <= angleRange.MaxY
}
