package main

// Welcome to
// __________         __    __  .__                               __
// \______   \_____ _/  |__/  |_|  |   ____   ______ ____ _____  |  | __ ____
//  |    |  _/\__  \\   __\   __\  | _/ __ \ /  ___//    \\__  \ |  |/ // __ \
//  |    |   \ / __ \|  |  |  | |  |_\  ___/ \___ \|   |  \/ __ \|    <\  ___/
//  |________/(______/__|  |__| |____/\_____>______>___|__(______/__|__\\_____>
//
// This file can be a nice home for your Battlesnake logic and helper functions.
//
// To get you started we've included code to prevent your Battlesnake from moving backwards.
// For more info see docs.battlesnake.com

import (
  "log"
)

// info is called when you create your Battlesnake on play.battlesnake.com
// and controls your Battlesnake's appearance
// TIP: If you open your Battlesnake URL in a browser you should see this data

func floodFill(state GameState, start Coord) int {
    width, height := state.Board.Width, state.Board.Height
    visited := make([][]bool, height)
    for i := range visited {
        visited[i] = make([]bool, width)
    }

    // Mark snake bodies as visited
    for _, snake := range state.Board.Snakes {
        for _, bodyPart := range snake.Body {
            visited[bodyPart.Y][bodyPart.X] = true
        }
    }

    return floodFillRecursive(start, visited, width, height)
}
func floodFillRecursive(coord Coord, visited [][]bool, width, height int) int {
    if coord.X < 0 || coord.Y < 0 || coord.X >= width || coord.Y >= height || visited[coord.Y][coord.X] {
        return 0
    }

    visited[coord.Y][coord.X] = true
    return 1 +
        floodFillRecursive(Coord{X: coord.X - 1, Y: coord.Y}, visited, width, height) +
        floodFillRecursive(Coord{X: coord.X + 1, Y: coord.Y}, visited, width, height) +
        floodFillRecursive(Coord{X: coord.X, Y: coord.Y - 1}, visited, width, height) +
        floodFillRecursive(Coord{X: coord.X, Y: coord.Y + 1}, visited, width, height)
}







func info() BattlesnakeInfoResponse {
  log.Println("INFO")

  return BattlesnakeInfoResponse{
    APIVersion: "1",
    Author:     "",        // TODO: Your Battlesnake username
    Color:      "#FF0000", // TODO: Choose color
    Head:       "default", // TODO: Choose head
    Tail:       "default", // TODO: Choose tail
  }
}

// start is called when your Battlesnake begins a game
func start(state GameState) {
  log.Println("GAME START")
}

// end is called when your Battlesnake finishes a game
func end(state GameState) {
  log.Printf("GAME OVER\n\n")
}

// move is called on every turn and returns your next move
// Valid moves are "up", "down", "left", or "right"
// See https://docs.battlesnake.com/api/example-move for available data
func move(state GameState) BattlesnakeMoveResponse {

  isMoveSafe := map[string]bool{
    "up":    true,
    "down":  true,
    "left":  true,
    "right": true,
  }



  // We've included code to prevent your Battlesnake from moving backwards
  myHead := state.You.Body[0] // Coordinates of your head
  myNeck := state.You.Body[1] // Coordinates of your "neck"
  boardWidth := state.Board.Width
  boardHeight := state.Board.Height
  if myNeck.X < myHead.X { // Neck is left of head, don't move left
    isMoveSafe["left"] = false

  } else if myNeck.X > myHead.X { // Neck is right of head, don't move right
    isMoveSafe["right"] = false

  } else if myNeck.Y < myHead.Y { // Neck is below head, don't move down
    isMoveSafe["down"] = false

  } else if myNeck.Y > myHead.Y { // Neck is above head, don't move up
    isMoveSafe["up"] = false
  } 
  if myHead.X == 0{
     isMoveSafe["left"] = false
  }
  if myHead.X == boardWidth-1{
    isMoveSafe["right"] = false
  }
  if myHead.Y == boardHeight-1{
    isMoveSafe["up"] = false
  }
  if myHead.Y == 0 {
    isMoveSafe["down"] = false
  }


  // TODO: Step 1 - Prevent your Battlesnake from moving out of bounds - done

  // TODO: Step 2 - Prevent your Battlesnake from colliding with itself
  // mybody := state.You.Body
  myBody := state.You.Body
  for _, bodyPart := range myBody[1:] { // Exclude the head
      if myHead.X+1 == bodyPart.X && myHead.Y == bodyPart.Y {
          isMoveSafe["right"] = false
      } else if myHead.X-1 == bodyPart.X && myHead.Y == bodyPart.Y {
          isMoveSafe["left"] = false
      } else if myHead.Y+1 == bodyPart.Y && myHead.X == bodyPart.X {
          isMoveSafe["up"] = false
      } else if myHead.Y-1 == bodyPart.Y && myHead.X == bodyPart.X {
          isMoveSafe["down"] = false
      }
  }
  // TODO: Step 3 - Prevent your Battlesnake from colliding with other Battlesnakes
  // opponents := state.Board.Snakes

  // Are there any safe moves left?
  bestMove := ""
  maxArea := -1

  for move, isSafe := range isMoveSafe {
      if isSafe {
          nextHead := getNextHead(state.You.Head, move)
          area := floodFill(state, nextHead)
          if area > maxArea {
              maxArea = area
              bestMove = move
          }
      }
  }

  if bestMove == "" {
      bestMove = "down" // Fallback move
  }

  log.Printf("MOVE %d: %s\n", state.Turn, bestMove)
  return BattlesnakeMoveResponse{Move: bestMove}
}

func main() {
  RunServer()
}

func getNextHead(head Coord, move string) Coord {
    nextHead := head
    switch move {
    case "up":
        nextHead.Y++
    case "down":
        nextHead.Y--
    case "left":
        nextHead.X--
    case "right":
        nextHead.X++
    }
    return nextHead
}