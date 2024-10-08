package layout

import (
	"image/color"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Breakpoint int

const (
	Bp_ExtraSmall Breakpoint = iota
	Bp_Small
	Bp_Medium
	Bp_Large
	Bp_ExtraLarge
	Bp_ExtraExtraLarge
)

type BreakpointLayout struct {
	Breakpoint Breakpoint
	Width      int
	Height     int
	Columns    int // Number of columns in this layout
}

func (layout BreakpointLayout) DrawLayout(screen *ebiten.Image) {
	// Fill background with a color
	screen.Fill(color.RGBA{240, 240, 240, 255}) // Light gray background

	// Optionally, draw grid lines or other layout indicators
	// Example: Draw column guides
	// Note: This is optional and for visual debugging purposes
	for i := 1; i < layout.Columns; i++ {
		x := i * layout.Width / layout.Columns
		ebitenutil.DrawLine(screen, float64(x), 0, float64(x), float64(layout.Height), color.RGBA{200, 200, 200, 255})
	}
}

type BreakpointLayoutSystem struct {
	Layouts map[Breakpoint]BreakpointLayout
}

func NewLayoutSystem(customLayouts map[Breakpoint]BreakpointLayout) *BreakpointLayoutSystem {
	defaultLayouts := map[Breakpoint]BreakpointLayout{
		Bp_ExtraSmall:      {Breakpoint: Bp_ExtraSmall, Width: 320, Height: 480, Columns: 1},
		Bp_Small:           {Breakpoint: Bp_Small, Width: 576, Height: 768, Columns: 1},
		Bp_Medium:          {Breakpoint: Bp_Medium, Width: 768, Height: 1024, Columns: 2},
		Bp_Large:           {Breakpoint: Bp_Large, Width: 992, Height: 1280, Columns: 2},
		Bp_ExtraLarge:      {Breakpoint: Bp_ExtraLarge, Width: 1200, Height: 1600, Columns: 3},
		Bp_ExtraExtraLarge: {Breakpoint: Bp_ExtraExtraLarge, Width: 1400, Height: 1920, Columns: 4},
	}

	if customLayouts != nil {
		for bp, layout := range customLayouts {
			defaultLayouts[bp] = layout
		}
	}

	return &BreakpointLayoutSystem{
		Layouts: defaultLayouts,
	}
}

func (ls *BreakpointLayoutSystem) GetLayout(breakpoint Breakpoint) BreakpointLayout {
	return ls.Layouts[breakpoint]
}

// DetermineBreakpoint determines the current breakpoint based on the window width.
func (ls *BreakpointLayoutSystem) DetermineBreakpoint(width, height int) Breakpoint {
	// Create a slice to hold breakpoints with their corresponding widths
	type bpWidth struct {
		bp    Breakpoint
		width int
	}

	var bpWidths []bpWidth
	for bp, layout := range ls.Layouts {
		bpWidths = append(bpWidths, bpWidth{bp: bp, width: layout.Width})
	}

	// Sort the breakpoints by width in ascending order
	sort.Slice(bpWidths, func(i, j int) bool {
		return bpWidths[i].width < bpWidths[j].width
	})

	// Iterate through the sorted breakpoints to find the appropriate one
	for _, bpw := range bpWidths {
		if width < bpw.width {
			return bpw.bp
		}
	}

	// If none matched, return the largest breakpoint
	return Bp_ExtraExtraLarge
}
