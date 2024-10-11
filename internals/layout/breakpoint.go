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
	Columns    int
}

func (layout BreakpointLayout) DrawLayout(screen *ebiten.Image) {

	screen.Fill(color.RGBA{240, 240, 240, 255})

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

func (ls *BreakpointLayoutSystem) DetermineBreakpoint(width, height int) Breakpoint {

	type bpWidth struct {
		bp    Breakpoint
		width int
	}

	var bpWidths []bpWidth
	for bp, layout := range ls.Layouts {
		bpWidths = append(bpWidths, bpWidth{bp: bp, width: layout.Width})
	}

	sort.Slice(bpWidths, func(i, j int) bool {
		return bpWidths[i].width < bpWidths[j].width
	})

	for _, bpw := range bpWidths {
		if width < bpw.width {
			return bpw.bp
		}
	}

	return Bp_ExtraExtraLarge
}
