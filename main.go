/*
* magicedit v0.1.0
*
* Created by: Berkay Çubuk<berkay@berkaycubuk.com>
 */
package main

import (
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := 800
	screenHeight := 450

	rl.InitWindow(int32(screenWidth), int32(screenHeight), "magicedit v0.1.0")
	defer rl.CloseWindow()

	var textBuffer strings.Builder
	var commandBuffer strings.Builder
	commandBufferCursorPos := 0

	currentMode := "NORMAL"

	maxTextLength := 1000

	cursorPos := 0
	cursorLine := 0
    cursorCol := 0

	isHoldingBackspace := false
	isHoldingEnter := false

	backspaceTimer := float32(0)
	enterTimer := float32(0)

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

	for !closeWindow {
		if wantToClose || rl.WindowShouldClose() {
			closeWindow = true
		}

		deltaTime := rl.GetFrameTime()

		pressedChar := rl.GetCharPressed()

		for pressedChar > 0 {
			if currentMode == "INSERT" {
				if textBuffer.Len() < maxTextLength {
					textBuffer.WriteRune(rune(pressedChar))
					cursorPos++
					cursorLine, cursorCol = getCursorLineCol(textBuffer.String(), cursorPos)
				}
			} else if currentMode == "NORMAL" {
				if pressedChar == 105 { // i
					currentMode = "INSERT"
				} else if pressedChar == 58 { // :
					currentMode = "COMMAND"
					commandBuffer.WriteString(":")
					commandBufferCursorPos++
				}
			} else if currentMode == "COMMAND" {
				commandBuffer.WriteRune(rune(pressedChar))
				commandBufferCursorPos++
			}

			pressedChar = rl.GetCharPressed()
		}

		if rl.IsKeyPressed(rl.KeyEscape) {
			if currentMode == "INSERT" || currentMode == "COMMAND" {
				currentMode = "NORMAL"
				commandBuffer.Reset()
				commandBufferCursorPos = 0
			}
		}

		if rl.IsKeyDown(rl.KeyBackspace) {
			if !isHoldingBackspace {
				if cursorPos > 0 {
					textStr := textBuffer.String()
					textStr = textStr[:cursorPos - 1] + textStr[cursorPos:]
					textBuffer.Reset()
					textBuffer.WriteString(textStr)
					cursorPos--
					cursorLine, cursorCol = getCursorLineCol(textBuffer.String(), cursorPos)
				}

				isHoldingBackspace = true
				backspaceTimer = 0
			} else {
				backspaceTimer += deltaTime
				if backspaceTimer > keyRepeatDelay {
					if backspaceTimer - keyRepeatDelay > keyRepeatRate {
						if cursorPos > 0 {
							textStr := textBuffer.String()
							textStr = textStr[:cursorPos - 1] + textStr[cursorPos:]
							textBuffer.Reset()
							textBuffer.WriteString(textStr)
							cursorPos--
							cursorLine, cursorCol = getCursorLineCol(textBuffer.String(), cursorPos)
						}
						backspaceTimer = keyRepeatDelay
					}
				}
			}
		} else {
			isHoldingBackspace = false
		}

		if rl.IsKeyDown(rl.KeyEnter) {
			if !isHoldingEnter {
				if currentMode == "COMMAND" {
					if commandBuffer.String() == ":q" {
						wantToClose = true
					}
				}

				textBuffer.WriteString("\n")
				cursorPos++
				cursorLine, cursorCol = getCursorLineCol(textBuffer.String(), cursorPos)

				isHoldingEnter = true
				enterTimer = 0
			} else {
				enterTimer += deltaTime
				if enterTimer > keyRepeatDelay {
					if enterTimer - keyRepeatDelay > keyRepeatRate {
						textBuffer.WriteString("\n")
						cursorPos++
						cursorLine, cursorCol = getCursorLineCol(textBuffer.String(), cursorPos)
						enterTimer = keyRepeatDelay
					}
				}
			}
		} else {
			isHoldingEnter = false
		}

		/* Draw Area */
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		textStr := textBuffer.String()
		lines := strings.Split(textStr, "\n")
		yOffset := float32(10) // Y offset for the text
		
		for i, line := range lines {
			if i == cursorLine {
                cursorStartX := rl.MeasureTextEx(mainFont, line[:cursorCol], float32(fontSize), 1).X
                cursorStartY := yOffset + float32(i)*float32(lineHeight)

				if currentMode == "INSERT" {
					rl.DrawRectangle(int32(cursorStartX+10), int32(cursorStartY), 2, int32(lineHeight), rl.White)
				} else if currentMode == "NORMAL" {
					rl.DrawRectangle(int32(cursorStartX+10), int32(cursorStartY), 10, int32(lineHeight), rl.NewColor(255,255,255,50))
				}

			}

			rl.DrawTextEx(mainFont, line, rl.NewVector2(10, yOffset+float32(i)*float32(lineHeight)), float32(fontSize), 1, rl.White)
		}

		// Bottom bar
		rl.DrawTextEx(mainFont, currentMode, rl.NewVector2(10, float32(screenHeight - (2 * lineHeight) - 10)), float32(fontSize), 1, rl.White)

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
