/*
* magicedit v0.1.0
*
* Created by: Berkay Çubuk<berkay@berkaycubuk.com>
*/

package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	var textBuffer strings.Builder
	var filename string

	args := os.Args[1:]

	if len(args) > 0 {
		f, err := os.Open(args[0])
		if err != nil {
			panic(err)
		}
		defer f.Close()
		
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			textBuffer.WriteString(scanner.Text() + "\n")
		}

		filename = args[0]
	}

	screenWidth := 800
	screenHeight := 450

	viewOffsetY := 0
	viewOffsetX := 0

	tabSpace := 3
	rl.SetConfigFlags(rl.FlagWindowResizable)

	rl.InitWindow(int32(screenWidth), int32(screenHeight), "magicedit v0.1.0")
	defer rl.CloseWindow()

	var commandBuffer strings.Builder
	commandBufferCursorPos := 0

	currentMode := "NORMAL"

	cursorPos := 0
	cursorLine := 0
	cursorCol := 0

	holdingKey := ""
	keyTimer := float32(0)
	keyRepeatDelay := float32(0.5)
	keyRepeatRate := float32(0.05)

	wantToClose := false
	closeWindow := false

	fontSize := 20
	lineHeight := fontSize
	mainFont := rl.LoadFontEx("fonts/jetbrainsmono-regular.ttf", int32(fontSize), nil, 250)
	defer rl.UnloadFont(mainFont)

	rl.SetExitKey(rl.KeyNull)

	rl.SetTargetFPS(60)

	var lastPressedChar int32

	for !closeWindow {
		if wantToClose || rl.WindowShouldClose() {
			closeWindow = true
		}

		viewOffsetY += int(rl.GetMouseWheelMove()) * 20
		if viewOffsetY > 0 {
			viewOffsetY = 0
		}

		// Update window size
		screenWidth = rl.GetScreenWidth()
		screenHeight = rl.GetScreenHeight()

		deltaTime := rl.GetFrameTime()

		pressedChar := rl.GetCharPressed()

		for pressedChar > 0 {
			if currentMode == "INSERT" {
				textStr := textBuffer.String()
				textStr = textStr[:cursorPos] + string(pressedChar) + textStr[cursorPos:]
				textBuffer.Reset()
				textBuffer.WriteString(textStr)
				cursorPos++
				cursorLine, cursorCol = getCursorLineCol(textBuffer.String(), cursorPos)
			} else if currentMode == "NORMAL" {
				if pressedChar == 105 { // i
					currentMode = "INSERT"
				} else if pressedChar == 97 { // a
					// TODO: Slice out of bounds error
					currentMode = "INSERT"
					//cursorCol++
					//cursorPos++
					//cursorLine, cursorCol = getCursorLineCol(textBuffer.String(), cursorPos)
				} else if pressedChar == 58 { // :
					currentMode = "COMMAND"
					commandBuffer.WriteString(":")
					commandBufferCursorPos++
				} else if pressedChar == 71 { // G
					cursorPos = textBuffer.Len() - 1
					cursorLine, cursorCol = getCursorLineCol(textBuffer.String(), cursorPos)
				} else if pressedChar == 111 { // o
					// TODO: adds the new line to the top when used at the bottom
					currentMode = "INSERT"
					cursorLine++
					cursorCol = 0
					cursorPos = getLineStartPos(textBuffer.String(), cursorLine)

					textStr := textBuffer.String()
					textStr = textStr[:cursorPos] + "\n" + textStr[cursorPos:]
					textBuffer.Reset()
					textBuffer.WriteString(textStr)
				} else if pressedChar == 79 { // O
					currentMode = "INSERT"
					cursorCol = 0
					cursorPos = getLineStartPos(textBuffer.String(), cursorLine)

					textStr := textBuffer.String()
					textStr = textStr[:cursorPos] + "\n" + textStr[cursorPos:]
					textBuffer.Reset()
					textBuffer.WriteString(textStr)
				} else if lastPressedChar == pressedChar {
					if pressedChar == 103 { // g
						cursorPos = 0
						cursorLine, cursorCol = getCursorLineCol(textBuffer.String(), cursorPos)
					} else if pressedChar == 100 { // d
						cursorPos = getLineStartPos(textBuffer.String(), cursorLine)
						nextLineStartPos := getLineStartPos(textBuffer.String(), cursorLine + 1)

						textStr := textBuffer.String()
						textStr = textStr[:cursorPos] + textStr[nextLineStartPos:]
						textBuffer.Reset()
						textBuffer.WriteString(textStr)
					}
				}
			} else if currentMode == "COMMAND" {
				commandBuffer.WriteRune(rune(pressedChar))
				commandBufferCursorPos++
			}

			lastPressedChar = pressedChar
			pressedChar = rl.GetCharPressed()
		}

		if rl.IsKeyPressed(rl.KeyEscape) {
			if currentMode == "INSERT" || currentMode == "COMMAND" {
				currentMode = "NORMAL"
				commandBuffer.Reset()
				commandBufferCursorPos = 0
			}
		}

		if currentMode == "NORMAL" {
			if rl.IsKeyDown(rl.KeyH) && cursorPos > 0 {
				if holdingKey != "h" {
					cursorPos--
					cursorLine, cursorCol = getCursorLineCol(textBuffer.String(), cursorPos)
					holdingKey = "h"
					keyTimer = 0
				} else {
					keyTimer += deltaTime
					if keyTimer > keyRepeatDelay {
						if keyTimer - keyRepeatDelay > keyRepeatRate {
							cursorPos--
							cursorLine, cursorCol = getCursorLineCol(textBuffer.String(), cursorPos)
							keyTimer = keyRepeatDelay
						}
					}
				}
			} else if rl.IsKeyDown(rl.KeyL) && cursorPos < textBuffer.Len() {
				if holdingKey != "l" {
					cursorPos++
					cursorLine, cursorCol = getCursorLineCol(textBuffer.String(), cursorPos)
					holdingKey = "l"
					keyTimer = 0
				} else {
					keyTimer += deltaTime
					if keyTimer > keyRepeatDelay {
						if keyTimer - keyRepeatDelay > keyRepeatRate {
							cursorPos++
							cursorLine, cursorCol = getCursorLineCol(textBuffer.String(), cursorPos)
							keyTimer = keyRepeatDelay
						}
					}
				}
			} else if rl.IsKeyDown(rl.KeyK) && cursorLine > 0 {
				if holdingKey != "k" {
					cursorLine--
					cursorCol = min(cursorCol, len(getLine(textBuffer.String(), cursorLine)))
					cursorPos = getLineStartPos(textBuffer.String(), cursorLine) + cursorCol
					holdingKey = "k"
					keyTimer = 0
				} else {
					keyTimer += deltaTime
					if keyTimer > keyRepeatDelay {
						if keyTimer - keyRepeatDelay > keyRepeatRate {
							cursorLine--
							cursorCol = min(cursorCol, len(getLine(textBuffer.String(), cursorLine)))
							cursorPos = getLineStartPos(textBuffer.String(), cursorLine) + cursorCol
							keyTimer = keyRepeatDelay
						}
					}
				}
			} else if rl.IsKeyDown(rl.KeyJ) && cursorLine < getLineCount(textBuffer.String()) - 1 {
				if holdingKey != "j" {
					cursorLine++
					cursorCol = min(cursorCol, len(getLine(textBuffer.String(), cursorLine)))
					cursorPos = getLineStartPos(textBuffer.String(), cursorLine) + cursorCol
					holdingKey = "j"
					keyTimer = 0
				} else {
					keyTimer += deltaTime
					if keyTimer > keyRepeatDelay {
						if keyTimer - keyRepeatDelay > keyRepeatRate {
							cursorLine++
							cursorCol = min(cursorCol, len(getLine(textBuffer.String(), cursorLine)))
							cursorPos = getLineStartPos(textBuffer.String(), cursorLine) + cursorCol
							keyTimer = keyRepeatDelay
						}
					}
				}
			} else {
				holdingKey = ""
			}
		} else if currentMode == "COMMAND" {
			if rl.IsKeyPressed(rl.KeyEnter) {
				if commandBuffer.String() == ":q" {
					wantToClose = true
				} else if commandBuffer.String() == ":wq" {
					// Save file
					if filename != "" {
						err := os.WriteFile(filename, []byte(textBuffer.String()), 0755)
						if err != nil {
							log.Fatal(err)
						}
					}

					wantToClose = true
				}
			}
		} else if currentMode == "INSERT" {
			if rl.IsKeyDown(rl.KeyBackspace) {
				if holdingKey != "backspace" {
					if cursorPos > 0 {
						textStr := textBuffer.String()
						textStr = textStr[:cursorPos - 1] + textStr[cursorPos:]
						textBuffer.Reset()
						textBuffer.WriteString(textStr)
						cursorPos--
						cursorLine, cursorCol = getCursorLineCol(textBuffer.String(), cursorPos)
					}

					holdingKey = "backspace"
					keyTimer = 0
				} else {
					keyTimer += deltaTime
					if keyTimer > keyRepeatDelay {
						if keyTimer - keyRepeatDelay > keyRepeatRate {
							if cursorPos > 0 {
								textStr := textBuffer.String()
								textStr = textStr[:cursorPos - 1] + textStr[cursorPos:]
								textBuffer.Reset()
								textBuffer.WriteString(textStr)
								cursorPos--
								cursorLine, cursorCol = getCursorLineCol(textBuffer.String(), cursorPos)
							}
							keyTimer = keyRepeatDelay
						}
					}
				}
			} else if rl.IsKeyDown(rl.KeyEnter) {
				if holdingKey != "enter" {
					textStr := textBuffer.String()
					textStr = textStr[:cursorPos] + string("\n") + textStr[cursorPos:]
					textBuffer.Reset()
					textBuffer.WriteString(textStr)
					cursorPos++
					cursorLine, cursorCol = getCursorLineCol(textBuffer.String(), cursorPos)

					holdingKey = "enter"
					keyTimer = 0
				} else {
					keyTimer += deltaTime
					if keyTimer > keyRepeatDelay {
						if keyTimer - keyRepeatDelay > keyRepeatRate {
							textStr := textBuffer.String()
							textStr = textStr[:cursorPos] + string("\n") + textStr[cursorPos:]
							textBuffer.Reset()
							textBuffer.WriteString(textStr)
							cursorPos++
							cursorLine, cursorCol = getCursorLineCol(textBuffer.String(), cursorPos)

							keyTimer = keyRepeatDelay
						}
					}
				}
			} else if rl.IsKeyDown(rl.KeyTab) {
				if holdingKey != "tab" {
					textStr := textBuffer.String()
					textStr = textStr[:cursorPos] + string("\t") + textStr[cursorPos:]
					textBuffer.Reset()
					textBuffer.WriteString(textStr)
					cursorPos++
					cursorLine, cursorCol = getCursorLineCol(textBuffer.String(), cursorPos)

					holdingKey = "tab"
					keyTimer = 0
				} else {
					keyTimer += deltaTime
					if keyTimer > keyRepeatDelay {
						if keyTimer - keyRepeatDelay > keyRepeatRate {
							textStr := textBuffer.String()
							textStr = textStr[:cursorPos] + string("\t") + textStr[cursorPos:]
							textBuffer.Reset()
							textBuffer.WriteString(textStr)
							cursorPos++
							cursorLine, cursorCol = getCursorLineCol(textBuffer.String(), cursorPos)

							keyTimer = keyRepeatDelay
						}
					}
				}
			} else {
				holdingKey = ""
			}
		}

		lineNumberOffset := 4 * 10 // 4 for length

		// Scrolling with cursor position
		if cursorLine * lineHeight + 10 + (lineHeight * 5) >= screenHeight {
			viewOffsetY = screenHeight - (cursorLine * lineHeight + 10 + (lineHeight * 5))
		} else {
			viewOffsetY = 0
		}

		// TODO: this does not count the tab size
		if (cursorCol + 2) * 10 + 10 >= screenWidth - lineNumberOffset {
			viewOffsetX = (screenWidth - lineNumberOffset) - ((cursorCol + 2) * 10 + 10)
		} else {
			viewOffsetX = 0
		}

		/* Draw Area */
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		textStr := textBuffer.String()
		lines := strings.Split(textStr, "\n")
		yOffset := float32(10) // Y offset for the text
		
		for i, line := range lines {
			xOffset := 0
			tabOffset := 0
			for index, char := range line {
				rl.DrawTextEx(mainFont, string(char), rl.NewVector2(10 + float32(xOffset) + float32(viewOffsetX) + float32(lineNumberOffset), float32(viewOffsetY) + yOffset+float32(i)*float32(lineHeight)), float32(fontSize), 1, rl.White)
				
				if string(char) == "\t" {
					rl.DrawTextEx(mainFont, "»", rl.NewVector2(10 + float32(xOffset) + float32(viewOffsetX) + float32(lineNumberOffset), float32(viewOffsetY) + yOffset+float32(i)*float32(lineHeight)), float32(fontSize), 1, rl.NewColor(255,255,255,60))

					xOffset += tabSpace * 10
					if index < cursorCol {
						tabOffset += (tabSpace - 1) * 10
					}
				} else {
					xOffset += 10
				}
			}

			// Line numbers
			rl.DrawTextEx(mainFont, strconv.Itoa(i + 1), rl.NewVector2(10, float32(viewOffsetY) + yOffset+float32(i)*float32(lineHeight)), float32(fontSize), 1, rl.NewColor(255,255,255,60))

			if i == cursorLine {
                cursorStartX := rl.MeasureTextEx(mainFont, line[:cursorCol], float32(fontSize), 1).X + float32(viewOffsetX) + float32(lineNumberOffset)
                cursorStartY := yOffset + float32(i)*float32(lineHeight) + float32(viewOffsetY)

				if currentMode == "INSERT" {
					rl.DrawRectangle(int32(cursorStartX+10) + int32(tabOffset), int32(cursorStartY), 2, int32(lineHeight), rl.White)
				} else if currentMode == "NORMAL" {
					rl.DrawRectangle(int32(cursorStartX+10) + int32(tabOffset), int32(cursorStartY), 10, int32(lineHeight), rl.NewColor(255,255,255,50))
				}
			}
		}

		// Bottom bar
		rl.DrawRectangle(0, int32(screenHeight - 2 *lineHeight), int32(screenWidth), int32(3 *lineHeight), rl.NewColor(0, 117, 44, 150))
		rl.DrawTextEx(mainFont, currentMode, rl.NewVector2(float32(screenWidth- 10 - (10 * len(currentMode))), float32(screenHeight - lineHeight - 10)), float32(fontSize), 1, rl.White)

		// Command bar
		rl.DrawTextEx(mainFont, commandBuffer.String(), rl.NewVector2(10, float32(screenHeight - lineHeight - 10)), float32(fontSize), 1, rl.White)

		rl.EndDrawing()
		/* Draw Area */
	}
}

func getCursorLineCol(text string, pos int) (line, col int) {
	lines := strings.Split(text[:pos], "\n")
	return len(lines) - 1, len(lines[len(lines)-1])
}

func getLine(text string, line int) string {
	lines := strings.Split(text, "\n")
	if line < len(lines) {
		return lines[line]
	}
	return ""
}

func getLineStartPos(text string, line int) int {
	lines := strings.Split(text, "\n")
	if line < len(lines) {
		start := 0
		for i := 0; i < line; i++ {
			start += len(lines[i]) + 1 // +1 for newline character
		}
		return start
	}
	return 0
}

func getLineCount(text string) int {
	return len(strings.Split(text, "\n"))
}
