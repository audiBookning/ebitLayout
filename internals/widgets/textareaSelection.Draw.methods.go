package widgets

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (t *TextAreaSelection) drawBackground(screen *ebiten.Image) {
	// Draw the background of the text area
	vector.DrawFilledRect(screen, float32(t.x), float32(t.y), float32(t.w), float32(t.h), color.RGBA{200, 200, 200, 255}, true)
	t.drawGrid(screen)
}

func (t *TextAreaSelection) drawSelection(screen *ebiten.Image, minPos, maxPos, lineIndex int, line string, yOffset float64) {
	startLineSel, startCol := t.getCursorLineAndColForPos(minPos)
	endLineSel, endCol := t.getCursorLineAndColForPos(maxPos)

	if lineIndex >= startLineSel && lineIndex <= endLineSel {
		// Determine selection bounds within the current line
		var selStart, selEnd int
		if lineIndex == startLineSel {
			selStart = startCol
		} else {
			selStart = 0
		}
		if lineIndex == endLineSel {
			selEnd = endCol
		} else {
			selEnd = len(line)
		}

		// Calculate x positions based on byte indices
		selectionXStart := t.x + t.paddingLeft + int(t.textWidth(line[:selStart]))
		selectionXEnd := t.x + t.paddingLeft + int(t.textWidth(line[:selEnd]))

		// Clamp the selection rectangle within textarea bounds
		selectionXStart = clamp(selectionXStart, t.x, t.x+t.w)
		selectionXEnd = clamp(selectionXEnd, t.x, t.x+t.w)

		// Clamp the yOffset within textarea bounds
		clampedYOffset := clampFloat(yOffset, float64(t.y+t.paddingTop), float64(t.y+t.h+t.paddingTop))

		// Draw the selection rectangle
		vector.DrawFilledRect(screen,
			float32(selectionXStart),
			float32(clampedYOffset),
			float32(selectionXEnd-selectionXStart),
			float32(t.lineHeight),
			color.RGBA{0, 0, 255, 128},
			true)
	}
}

func (t *TextAreaSelection) drawScrollbar(screen *ebiten.Image, totalLines int) {
	// Initialize scrollbar properties
	t.scrollbarWidth = 10
	t.scrollbarX = t.x + t.w - t.scrollbarWidth
	t.scrollbarY = t.y
	t.scrollbarHeight = t.h

	// Calculate the height of the scrollbar thumb
	visibleRatio := float64(t.maxLines) / float64(totalLines)
	t.scrollbarThumbH = float64(t.scrollbarHeight) * visibleRatio
	if t.scrollbarThumbH < 20 {
		t.scrollbarThumbH = 20 // Minimum thumb height
	}

	// Calculate the Y position of the scrollbar thumb relative to the scrollbar track
	maxScrollOffset := totalLines - t.maxLines
	if maxScrollOffset < 1 {
		maxScrollOffset = 1
	}
	thumbMaxY := float64(t.scrollbarHeight - int(t.scrollbarThumbH))
	t.scrollbarThumbY = float64(t.scrollOffset) / float64(maxScrollOffset) * thumbMaxY

	// Draw the scrollbar track
	vector.DrawFilledRect(screen, float32(t.scrollbarX), float32(t.scrollbarY), float32(t.scrollbarWidth), float32(t.scrollbarHeight), color.RGBA{220, 220, 220, 255}, true)

	// Draw the scrollbar thumb
	vector.DrawFilledRect(
		screen,
		float32(t.scrollbarX),
		float32(t.scrollbarY)+float32(t.scrollbarThumbY),
		float32(t.scrollbarWidth),
		float32(t.scrollbarThumbH),
		color.RGBA{160, 160, 160, 255},
		true)
}

func (t *TextAreaSelection) drawCursor(screen *ebiten.Image) {
	cursorLine, cursorCol := t.getCursorLineAndCol()
	if cursorLine >= t.scrollOffset && cursorLine < t.scrollOffset+t.maxLines {

		cursorX := t.x + t.paddingLeft + int(t.textWidth(strings.Split(t.text, "\n")[cursorLine][:cursorCol]))

		cursorY := float64(t.y+t.paddingTop) + float64(cursorLine-t.scrollOffset)*t.lineHeight

		// Clamp the cursor position within textarea bounds
		cursorX = clamp(cursorX, t.x, t.x+t.w)
		cursorY = clampFloat(cursorY, float64(t.y), float64(t.y+t.h)-t.lineHeight)

		// Render the blinking cursor
		if t.counter%(t.cursorBlinkRate*2) < t.cursorBlinkRate {
			vector.DrawFilledRect(screen,
				float32(cursorX),
				float32(cursorY),
				2,
				float32(t.lineHeight),
				color.RGBA{0, 0, 0, 255},
				true)
		}
	}
}

func (t *TextAreaSelection) drawGrid(screen *ebiten.Image) {
	gridColor := color.RGBA{255, 0, 0, 255} // Red color
	strokeWidth := float32(1)               // Thickness of grid lines

	// Calculate the drawable area considering padding
	startX := float32(t.x + t.paddingLeft)
	endX := float32(t.x + t.w - t.paddingLeft)
	startY := float32(t.y + t.paddingTop)
	endY := float32(t.y + t.h - t.paddingTop)

	// Draw vertical lines
	for x := float64(startX); x <= float64(endX); x += t.stepX {
		vector.StrokeLine(
			screen,
			float32(x),
			startY,
			float32(x),
			endY,
			strokeWidth,
			gridColor,
			false, // antialiasing
		)
	}

	// Draw horizontal lines
	for y := float64(startY); y <= float64(endY); y += t.stepY {
		vector.StrokeLine(
			screen,
			startX,
			float32(y),
			endX,
			float32(y),
			strokeWidth,
			gridColor,
			false, // antialiasing
		)
	}
}

func (t *TextAreaSelection) getCursorLineAndCol() (int, int) {
	return t.getCursorLineAndColForPos(t.cursorPos)
}
