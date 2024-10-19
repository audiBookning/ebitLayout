package responsive

import (
	"example.com/menu/cmd02/more06/types"
)

type LayoutMode string

const (
	LayoutHorizontal LayoutMode = "horizontal"
	LayoutVertical   LayoutMode = "vertical"
	LayoutGrid       LayoutMode = "grid"
)

type Breakpoint struct {
	Width      int
	LayoutMode LayoutMode
}

type Alignment int

const (
	AlignLeft Alignment = iota
	AlignCenter
	AlignRight
)

func getElementSizes(elements []types.Element) []types.Position {
	sizes := make([]types.Position, len(elements))
	for i, elem := range elements {
		width, height := elem.GetSize()
		sizes[i] = types.Position{
			Width:  width,
			Height: height,
		}
	}
	return sizes
}

func getElementsByIds(ids []string) []types.Element {

	return nil
}
