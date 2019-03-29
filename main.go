package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"gopkg.in/AlecAivazis/survey.v1"
)

type coordinates struct {
	x int
	y int
}

type tile struct {
	content string
	state   int // 0 => hide, 1 => visible, 2 => flag, 3 => unsafe
}

type mode struct {
	name   string
	width  int
	height int
	mines  int
}

var board []tile
var gameMode mode
var clear string

func getIndex(coord coordinates) int {
	return coord.y*gameMode.width + coord.x
}

func selectMode() {
	prompt := &survey.Select{
		Message: "Chose a difficulty:",
		Options: []string{"Easy", "Medium", "Hard"},
	}
	err := survey.AskOne(prompt, &gameMode.name, nil)
	if err != nil {
		os.Exit(1)
	}
}

func generateCoord() coordinates {
	var coord coordinates
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	coord.x = rng.Intn(gameMode.width)
	coord.y = rng.Intn(gameMode.height)

	return coord
}

func hasMine(coord coordinates) bool {
	if board[coord.y*10+coord.x].content == "💥" {
		return true
	}
	return false
}

func generateBoard() {
	switch gameMode.name {
	case "Easy":
		gameMode.width = 10
		gameMode.height = 10
		gameMode.mines = 10
		break

	case "Medium":
		gameMode.width = 16
		gameMode.height = 16
		gameMode.mines = 40
		break

	case "Hard":
		gameMode.width = 30
		gameMode.height = 16
		gameMode.mines = 99
		break
	}

	// create the board
	board = make([]tile, gameMode.width*gameMode.height)

	// add mines
	for i := 0; i < gameMode.mines; i++ {
		coord := generateCoord()
		if hasMine(coord) {
			i--
		} else {
			board[getIndex(coord)].content = "💥"
		}
	}

	// set numbers. I loved doing this part 💕
	for y := 0; y < gameMode.height; y++ {
		for x := 0; x < gameMode.width; x++ {
			i := 0
			if board[getIndex(coordinates{x, y})].content != "💥" {
				if y-1 >= 0 && board[getIndex(coordinates{x, y - 1})].content == "💥" {
					i++
				}
				if y-1 >= 0 && x+1 < gameMode.width && board[getIndex(coordinates{x + 1, y - 1})].content == "💥" {
					i++
				}
				if x+1 < gameMode.width && board[getIndex(coordinates{x + 1, y})].content == "💥" {
					i++
				}
				if y+1 < gameMode.height && x+1 < gameMode.width && board[getIndex(coordinates{x + 1, y + 1})].content == "💥" {
					i++
				}
				if y+1 < gameMode.height && board[getIndex(coordinates{x, y + 1})].content == "💥" {
					i++
				}
				if y+1 < gameMode.height && x-1 >= 0 && board[getIndex(coordinates{x - 1, y + 1})].content == "💥" {
					i++
				}
				if x-1 >= 0 && board[getIndex(coordinates{x - 1, y})].content == "💥" {
					i++
				}
				if y-1 >= 0 && x-1 >= 0 && board[getIndex(coordinates{x - 1, y - 1})].content == "💥" {
					i++
				}
				switch i {
				case 0:
					board[getIndex(coordinates{x, y})].content = "0️⃣ "
					break

				case 1:
					board[getIndex(coordinates{x, y})].content = "1️⃣ "
					break

				case 2:
					board[getIndex(coordinates{x, y})].content = "2️⃣ "
					break

				case 3:
					board[getIndex(coordinates{x, y})].content = "3️⃣ "
					break

				case 4:
					board[getIndex(coordinates{x, y})].content = "4️⃣ "
					break

				case 5:
					board[getIndex(coordinates{x, y})].content = "5️⃣ "
					break

				case 6:
					board[getIndex(coordinates{x, y})].content = "6️⃣ "
					break

				case 7:
					board[getIndex(coordinates{x, y})].content = "7️⃣ "
					break

				case 8:
					board[getIndex(coordinates{x, y})].content = "8️⃣ "
					break
				}
			}
		}
	}
}

func drawBoard(coord coordinates) {
	cmd := exec.Command(clear)
	cmd.Stdout = os.Stdout
	cmd.Run()

	blinkOn := "\033[5m"
	blinkOff := "\033[0m"

	fmt.Print("    ")
	for i := 0; i < gameMode.width; i++ {
		fmt.Print("", i)
		if i < 10 {
			fmt.Print(" ")
		}
	}
	fmt.Print("\n    ")
	for i := 0; i < gameMode.width; i++ {
		fmt.Print("⬇️ ")
	}
	fmt.Print("\n")

	for y := 0; y < gameMode.height; y++ {
		fmt.Print(y, " ➡️ ")
		if coord.y == y && coord.x == -1 {
			fmt.Print(blinkOn)
		}
		for x := 0; x < gameMode.width; x++ {
			if coord.y == y && coord.x == x {
				fmt.Print(blinkOn)
			}
			if board[getIndex(coordinates{x, y})].state == 0 {
				fmt.Print("🔳")
			} else if board[getIndex(coordinates{x, y})].state == 1 {
				fmt.Print(board[getIndex(coordinates{x, y})].content)
			} else if board[getIndex(coordinates{x, y})].state == 2 {
				fmt.Print("⛔")
			} else if board[getIndex(coordinates{x, y})].state == 3 {
				fmt.Print("❓")
			}
			if coord.y == y && coord.x == x {
				fmt.Print(blinkOff)
			}
			// fmt.Print(board[getIndex(coordinates{x, y})].content, "\t")
		}
		if coord.y == y && coord.x == -1 {
			fmt.Print(blinkOff)
		}
		fmt.Print("\n")
	}
}

func selectCoord() coordinates {
	coord := coordinates{-1, -1}
	array := make([]string, gameMode.height)

	for i := 0; i < gameMode.height; i++ {
		array[i] = strconv.Itoa(i)
	}

	prompt := &survey.Select{
		Message:  "Choose a line:",
		Options:  array,
		PageSize: 20,
	}

	err := survey.AskOne(prompt, &coord.y, nil)
	if err != nil {
		os.Exit(1)
	}

	drawBoard(coord)

	array = make([]string, gameMode.width)

	for i := 0; i < gameMode.width; i++ {
		array[i] = strconv.Itoa(i)
	}

	prompt = &survey.Select{
		Message:  "Choose a column:",
		Options:  array,
		PageSize: 20,
	}
	err = survey.AskOne(prompt, &coord.x, nil)
	if err != nil {
		os.Exit(1)
	}

	drawBoard(coord)

	return coord
}

func selectAction(coord coordinates) string {
	var action string
	var array []string

	if board[getIndex(coord)].state == 0 {
		array = append(array, "Discover", "Signal a mine", "Set unsafe")
	} else if board[getIndex(coord)].state == 1 {
		return "None"
	} else if board[getIndex(coord)].state == 2 {
		array = append(array, "Set unsafe", "Remove tag")
	} else if board[getIndex(coord)].state == 3 {
		array = append(array, "Signal a mine", "Remove tag")
	}

	prompt := &survey.Select{
		Message:  "Choose an action:",
		Options:  array,
		PageSize: 20,
	}
	err := survey.AskOne(prompt, &action, nil)
	if err != nil {
		os.Exit(1)
	}

	return action
}

func revealMore(coord coordinates) {
	if coord.y-1 >= 0 && board[getIndex(coordinates{coord.x, coord.y - 1})].state == 0 {
		board[getIndex(coordinates{coord.x, coord.y - 1})].state = 1
		if board[getIndex(coordinates{coord.x, coord.y - 1})].content == "0️⃣ " {
			revealMore(coordinates{coord.x, coord.y - 1})
		}
	}
	if coord.y-1 >= 0 && coord.x+1 < gameMode.width && board[getIndex(coordinates{coord.x + 1, coord.y - 1})].state == 0 {
		board[getIndex(coordinates{coord.x + 1, coord.y - 1})].state = 1
		if board[getIndex(coordinates{coord.x + 1, coord.y - 1})].content == "0️⃣ " {
			revealMore(coordinates{coord.x + 1, coord.y - 1})
		}
	}
	if coord.x+1 < gameMode.width && board[getIndex(coordinates{coord.x + 1, coord.y})].state == 0 {
		board[getIndex(coordinates{coord.x + 1, coord.y})].state = 1
		if board[getIndex(coordinates{coord.x + 1, coord.y})].content == "0️⃣ " {
			revealMore(coordinates{coord.x + 1, coord.y})
		}
	}
	if coord.y+1 < gameMode.height && coord.x+1 < gameMode.width && board[getIndex(coordinates{coord.x + 1, coord.y + 1})].state == 0 {
		board[getIndex(coordinates{coord.x + 1, coord.y + 1})].state = 1
		if board[getIndex(coordinates{coord.x + 1, coord.y + 1})].content == "0️⃣ " {
			revealMore(coordinates{coord.x + 1, coord.y + 1})
		}
	}
	if coord.y+1 < gameMode.height && board[getIndex(coordinates{coord.x, coord.y + 1})].state == 0 {
		board[getIndex(coordinates{coord.x, coord.y + 1})].state = 1
		if board[getIndex(coordinates{coord.x, coord.y + 1})].content == "0️⃣ " {
			revealMore(coordinates{coord.x, coord.y + 1})
		}
	}
	if coord.y+1 < gameMode.height && coord.x-1 >= 0 && board[getIndex(coordinates{coord.x - 1, coord.y + 1})].state == 0 {
		board[getIndex(coordinates{coord.x - 1, coord.y + 1})].state = 1
		if board[getIndex(coordinates{coord.x - 1, coord.y + 1})].content == "0️⃣ " {
			revealMore(coordinates{coord.x - 1, coord.y + 1})
		}
	}
	if coord.x-1 >= 0 && board[getIndex(coordinates{coord.x - 1, coord.y})].state == 0 {
		board[getIndex(coordinates{coord.x - 1, coord.y})].state = 1
		if board[getIndex(coordinates{coord.x - 1, coord.y})].content == "0️⃣ " {
			revealMore(coordinates{coord.x - 1, coord.y})
		}
	}
	if coord.y-1 >= 0 && coord.x-1 >= 0 && board[getIndex(coordinates{coord.x - 1, coord.y - 1})].state == 0 {
		board[getIndex(coordinates{coord.x - 1, coord.y - 1})].state = 1
		if board[getIndex(coordinates{coord.x - 1, coord.y - 1})].content == "0️⃣ " {
			revealMore(coordinates{coord.x - 1, coord.y - 1})
		}
	}
}

func manageTile(coord coordinates, action string) {
	switch action {
	case "Discover":
		board[getIndex(coord)].state = 1
		if board[getIndex(coord)].content == "💥" {
			drawBoard(coord)
			fmt.Println("Ho you lost 😞. Try again !")
			os.Exit(1)
		}
		if board[getIndex(coord)].content == "0️⃣ " {
			revealMore(coord)
		}
		break

	case "Signal a mine":
		board[getIndex(coord)].state = 2
		break

	case "Set unsafe":
		board[getIndex(coord)].state = 3
		break

	case "Remove tag":
		board[getIndex(coord)].state = 0
		break

	case "None":
		break
	}
}

func checkWin() {
	count := 0
	for i := len(board); i < len(board); i++ {
		if board[i].state == 0 {
			return
		} else if board[i].state == 2 {
			count++
		}
	}
	if count == gameMode.mines-1 {
		fmt.Println("Congrats you win 🎉🎉🎉")
		os.Exit(0)
	}
}

func main() {
	if runtime.GOOS == "windows" {
		clear = "cls"
	} else {
		clear = "clear"
	}

	selectMode()
	generateBoard()
	for true {
		drawBoard(coordinates{-1, -1})
		checkWin()
		coord := selectCoord()
		action := selectAction(coord)
		manageTile(coord, action)
	}
}
