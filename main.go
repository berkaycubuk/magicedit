package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func clearScreen() {
	fmt.Print("\033[2J") // clear screen
	fmt.Print("\033[H") // move cursor to top left
}

func printInfoBar(height int, width int, mode string, row int, col int) {
	fmt.Printf("\033[%d;1H", height) // move the cursor to the bottom left

	fmt.Printf("%s	%d:%d", mode, row, col)
}

func moveCursor(row int, col int) {
	fmt.Printf("\033[%d;%dH", row, col)
}

func main() {
	clearScreen()

	oldState, err := term.MakeRaw(int(os.Stderr.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}

	row := 0
	col := 0

	lineCount := 0

	fmt.Println("MagicEdit v0.1.0")

	mode := "NORMAL"

	isCleared := false

	for {
		var buf [1]byte
		os.Stdin.Read(buf[:])

		if !isCleared {
			isCleared = true
			clearScreen()
		}

		if mode == "NORMAL" {
		} else if mode == "INSERT" {
			if buf[0] == 13 {
				row += 1
				col = 0
			} else {
				fmt.Printf("%c", buf[0])
				col += 1
			}
		}

		if mode == "NORMAL" && buf[0] == 'j' {
			row += 1
			if row > height {
				row = height
			}
		}

		if mode == "NORMAL" && buf[0] == 'k' {
			row -= 1
			if row < 0 {
				row = 0
			}
		}

		if mode == "NORMAL" && buf[0] == 'l' {
			col += 1
			if col > width {
				col = width
			}
		}

		if mode == "NORMAL" && buf[0] == 'h' {
			col -= 1
			if col < 0 {
				col = 0
			}
		}

		if buf[0] == '\r' {
			fmt.Printf("%c", '\n')
		}

		if buf[0] == 'i' {
			mode = "INSERT"
		}

		if buf[0] == 27 { // esc
			mode = "NORMAL"
		}

		if mode == "NORMAL" && buf[0] == 'q' {
			break
		}

		printInfoBar(height, width, mode, row, col)
		moveCursor(row, col)
	}

	clearScreen()
}
