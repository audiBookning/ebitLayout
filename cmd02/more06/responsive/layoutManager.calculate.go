package responsive

import (
	"example.com/menu/cmd02/more06/types"
)

func (lm *LayoutManager) calcVertical(screenWidth, screenHeight int, elementIDs []string, alignment Alignment, elements []types.Element) map[string]types.Position {
	numElements := len(elements)
	if numElements == 0 {
		return nil
	}

	spacing := 20

	totalHeight := 0
	for _, elem := range elements {
		_, height := elem.GetSize()
		totalHeight += height
	}
	totalHeight += (numElements - 1) * spacing

	startY := (screenHeight - totalHeight) / 2

	var startX int
	switch alignment {
	case AlignLeft:
		startX = 0
	case AlignCenter:
		startX = (screenWidth) / 2
	case AlignRight:
		startX = screenWidth
	}

	positions := make(map[string]types.Position)
	y := startY
	for i, elem := range elements {
		width, height := elem.GetSize()
		switch alignment {
		case AlignLeft:

		case AlignCenter:
			startX = (screenWidth - width) / 2
		case AlignRight:
			startX = screenWidth - width
		}
		positions[elementIDs[i]] = types.Position{
			X:      startX,
			Y:      y,
			Width:  width,
			Height: height,
		}
		y += height + spacing
	}

	return positions
}

func (lm *LayoutManager) calcGrid(screenWidth, screenHeight int, elementIDs []string, alignment Alignment, elements []types.Element) map[string]types.Position {
	numElements := len(elements)
	if numElements == 0 {
		return nil
	}

	columns := 2
	rows := (numElements + columns - 1) / columns

	spacingX := 30
	spacingY := 30

	colWidths := make([]int, columns)
	rowHeights := make([]int, rows)

	for i, elem := range elements {
		col := i % columns
		row := i / columns
		width, height := elem.GetSize()
		if width > colWidths[col] {
			colWidths[col] = width
		}
		if height > rowHeights[row] {
			rowHeights[row] = height
		}
	}

	totalWidth := 0
	for _, w := range colWidths {
		totalWidth += w
	}
	totalWidth += (columns - 1) * spacingX

	totalHeight := 0
	for _, h := range rowHeights {
		totalHeight += h
	}
	totalHeight += (rows - 1) * spacingY

	var startX, startY int

	switch alignment {
	case AlignLeft:
		startX = 0
	case AlignCenter:
		startX = (screenWidth - totalWidth) / 2
	case AlignRight:
		startX = screenWidth - totalWidth
	}

	startY = (screenHeight - totalHeight) / 2

	positions := make(map[string]types.Position)

	colOffsets := make([]int, columns)
	currentX := startX
	for c := 0; c < columns; c++ {
		colOffsets[c] = currentX
		currentX += colWidths[c] + spacingX
	}

	rowOffsets := make([]int, rows)
	currentY := startY
	for r := 0; r < rows; r++ {
		rowOffsets[r] = currentY
		currentY += rowHeights[r] + spacingY
	}

	for i, elem := range elements {
		col := i % columns
		row := i / columns
		width, height := elem.GetSize()
		x := colOffsets[col]
		y := rowOffsets[row]
		positions[elementIDs[i]] = types.Position{
			X:      x,
			Y:      y,
			Width:  width,
			Height: height,
		}
	}

	return positions
}
