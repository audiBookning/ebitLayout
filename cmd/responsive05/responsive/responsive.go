package responsive

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

type Position struct {
	X, Y          int
	Width, Height int
}
