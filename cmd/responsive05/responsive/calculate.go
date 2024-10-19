package responsive

func calculateHorizontal(screenWidth, screenHeight int, elements []string) map[string]Position {
	numElements := len(elements)
	if numElements == 0 {
		return nil
	}

	buttonWidth := 200
	buttonHeight := 50
	spacing := 50

	totalWidth := numElements*buttonWidth + (numElements-1)*spacing
	startX := (screenWidth - totalWidth) / 2
	yPos := screenHeight - buttonHeight - 50

	positions := make(map[string]Position)
	for i, elem := range elements {
		x := startX + i*(buttonWidth+spacing)
		positions[elem] = Position{
			X:      x,
			Y:      yPos,
			Width:  buttonWidth,
			Height: buttonHeight,
		}
	}

	return positions
}

func calculateVertical(screenWidth, screenHeight int, elements []string) map[string]Position {
	numElements := len(elements)
	if numElements == 0 {
		return nil
	}

	buttonWidth := 150
	buttonHeight := 40
	spacing := 20

	totalHeight := numElements*buttonHeight + (numElements-1)*spacing
	startY := (screenHeight - totalHeight) / 2

	positions := make(map[string]Position)
	for i, elem := range elements {
		x := (screenWidth - buttonWidth) / 2
		y := startY + i*(buttonHeight+spacing)
		positions[elem] = Position{
			X:      x,
			Y:      y,
			Width:  buttonWidth,
			Height: buttonHeight,
		}
	}

	return positions
}

func calculateGrid(screenWidth, screenHeight int, elements []string) map[string]Position {
	numElements := len(elements)
	if numElements == 0 {
		return nil
	}

	columns := 2
	rows := (numElements + 1) / 2

	buttonWidth := 180
	buttonHeight := 45
	spacingX := 30
	spacingY := 30

	totalWidth := columns*buttonWidth + (columns-1)*spacingX
	totalHeight := rows*buttonHeight + (rows-1)*spacingY

	startX := (screenWidth - totalWidth) / 2
	startY := (screenHeight - totalHeight) / 2

	positions := make(map[string]Position)
	for i, elem := range elements {
		row := i / columns
		col := i % columns
		x := startX + col*(buttonWidth+spacingX)
		y := startY + row*(buttonHeight+spacingY)
		positions[elem] = Position{
			X:      x,
			Y:      y,
			Width:  buttonWidth,
			Height: buttonHeight,
		}
	}

	return positions
}
