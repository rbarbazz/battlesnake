package main

// This file can be a nice home for your Battlesnake logic and related helper functions.
//
// We have started this for you, with a function to help remove the 'neck' direction
// from the list of possible moves!

import (
	"log"
	"math"
)

// This function is called when you register your Battlesnake on play.battlesnake.com
// See https://docs.battlesnake.com/guides/getting-started#step-4-register-your-battlesnake
// It controls your Battlesnake appearance and author permissions.
// For customization options, see https://docs.battlesnake.com/references/personalization
// TIP: If you open your Battlesnake URL in browser you should see this data.
func info() BattlesnakeInfoResponse {
	log.Println("INFO")
	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "rbarbazz",
		Color:      "#cc241d",
		Head:       "default", // TODO: Personalize
		Tail:       "default", // TODO: Personalize
	}
}

// This function is called everytime your Battlesnake is entered into a game.
// The provided GameState contains information about the game that's about to be played.
// It's purely for informational purposes, you don't have to make any decisions here.
func start(state GameState) {
	log.Printf("%s START\n", state.Game.ID)
}

// This function is called when a game your Battlesnake was in has ended.
// It's purely for informational purposes, you don't have to make any decisions here.
func end(state GameState) {
	log.Printf("%s END\n\n", state.Game.ID)
}

func getPosition(coord Coord, width int) int {
	return coord.X + coord.Y*width
}

type PositionArray struct {
	positions [][]bool // true means the position is already taken, false means it's free
	height    int
	width     int
}

// This function initializes a new position array [][]bool where the default values are false
func newPositionArray(width int, height int) *PositionArray {
	positionArray := PositionArray{height: height, width: width}
	positions := make([][]bool, height)
	for i, _ := range positions {
		positions[i] = make([]bool, width)
	}

	positionArray.positions = positions

	return &positionArray
}

// Marks the positions where a snake is located as true in the position array
func (p PositionArray) processPositions(coords []Coord) {
	for _, coord := range coords {
		// The origin of the map is 0 0, starting from the bottom left
		// Our position array starts from the top left
		p.positions[p.height-1-coord.Y][coord.X] = true
	}
}

type PotentialMove struct {
	Move        string
	IsAvailable bool
}

// Returns a possible move based on the position array
func (p PositionArray) findNextMove(head Coord) [4]PotentialMove {
	flippedY := p.height - 1 - head.Y
	potentialMoves := [4]PotentialMove{
		PotentialMove{Move: "up", IsAvailable: false},
		PotentialMove{Move: "down", IsAvailable: false},
		PotentialMove{Move: "left", IsAvailable: false},
		PotentialMove{Move: "right", IsAvailable: false},
	}

	// Up
	if flippedY-1 > 0 && !p.positions[flippedY-1][head.X] {
		potentialMoves[0].IsAvailable = true
	}
	// Down
	if flippedY+1 < p.height && !p.positions[flippedY+1][head.X] {
		potentialMoves[1].IsAvailable = true
	}
	// Left
	if head.X-1 > 0 && !p.positions[flippedY][head.X-1] {
		potentialMoves[2].IsAvailable = true
	}
	// Right
	if head.X+1 < p.width && !p.positions[flippedY][head.X+1] {
		potentialMoves[3].IsAvailable = true
	}

	return potentialMoves
}

func identifyNearestFood(head Coord, foodList []Coord) Coord {
	type NearestFood struct {
		TotalDiff float64
		FoodItem  Coord
	}
	nearestFood := NearestFood{}

	for i, foodItem := range foodList {
		xDiff := math.Abs(float64(foodItem.X - head.X))
		yDiff := math.Abs(float64(foodItem.Y - head.Y))
		totalDiff := xDiff + yDiff

		if i == 0 || totalDiff < nearestFood.TotalDiff {
			nearestFood = NearestFood{TotalDiff: totalDiff, FoodItem: foodItem}
		}
	}

	return nearestFood.FoodItem
}

func getDirection(potentialMoves [4]PotentialMove, head Coord, foodItem Coord) string {
	if foodItem.X > head.X && potentialMoves[2].IsAvailable {
		return "left"
	} else if foodItem.X < head.X && potentialMoves[3].IsAvailable {
		return "right"
	} else if foodItem.Y > head.Y && potentialMoves[0].IsAvailable {
		return "up"
	}
	return "down"
}

// Todo for next time:
// Figure out if we should call identifyNearestFood before or after findNextMove
// Possibility of having a priority list of moves that would be passed around functions that help determine the next move

// This function is called on every turn of a game. Use the provided GameState to decide
// where to move -- valid moves are "up", "down", "left", or "right".
// We've provided some code and comments to get you started.
func move(state GameState) BattlesnakeMoveResponse {
	// Use information in GameState to prevent your Battlesnake from moving beyond the boundaries of the board.
	boardWidth := state.Board.Width
	boardHeight := state.Board.Height

	positionArray := newPositionArray(boardWidth, boardHeight)

	// Step 0: Don't let your Battlesnake move back in on it's own neck
	mybody := state.You.Body
	myHead := mybody[0] // Coordinates of your head

	// Call processPositions on all the snakes including myself
	positionArray.processPositions(state.You.Body)

	for _, snake := range state.Board.Snakes {
		positionArray.processPositions(snake.Body)
	}

	// TODO: Step 4 - Find food.
	// Use information in GameState to seek out and find food.

	return BattlesnakeMoveResponse{
		Move: positionArray.findNextMove(myHead),
	}
}
