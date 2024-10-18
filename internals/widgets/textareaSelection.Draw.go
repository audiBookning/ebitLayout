package widgets

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func (t *TextAreaSelection) Draw(screen *ebiten.Image) {

	// Split lines only if text has changed

	// Draw the background of the text area
	t.drawBackground(screen)

	yOffset := float64(t.y + t.paddingTop)
	// Apply scroll offset
	startLine := t.scrollOffset
	endLine := clamp(startLine+t.maxLines, 0, len(t.cachedLines))

	// Retrieve normalized selection bounds
	minPos, maxPos := t.getSelectionBounds()
	//fmt.Printf("Drawing selection from byte %d to byte %d\n", minPos, maxPos)

	for i := startLine; i < endLine; i++ {
		line := t.cachedLines[i]

		lineText := line
		lineX := t.x + t.paddingLeft
		lineY := int(yOffset)

		// Draw selection if active and within this line
		if minPos != maxPos {
			t.drawSelection(screen, minPos, maxPos, i, line, yOffset)
		}

		t.textWrapper.DrawText(screen, lineText, float64(lineX), float64(lineY))

		yOffset += t.lineHeight
	}

	// Draw the scrollbar if content exceeds maxLines
	if len(t.cachedLines) > t.maxLines {
		t.drawScrollbar(screen, len(t.cachedLines))
	}

	// Draw the cursor if the text area has focus and the cursor is within the visible lines
	if t.hasFocus {
		t.drawCursor(screen)
	}
}
