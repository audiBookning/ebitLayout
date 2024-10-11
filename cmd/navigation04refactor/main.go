package main

import (
	"image"
	"image/color"
	"log"
	"path/filepath"
	"runtime"

	"example.com/menu/internals/textwrapper"
	"example.com/menu/internals/widgets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth     = 800
	screenHeight    = 600
	leftColumnWidth = 160
)

var (
	filePathTxt          string
	Assets_Relative_Path = "../../"

	pageRegistry    = make(map[string]*Page)
	globalNavigator *Navigator
)

const (
	PageIDRed     = "redPage"
	PageIDBlue    = "bluePage"
	PageIDGreen   = "greenPage"
	PageIDYellow  = "yellowPage"
	PageIDMagenta = "magentaPage"
)

type PageConfig struct {
	id      string
	color   color.Color
	message string
	x, y    float32
	width   float32
	height  float32
}

type PageButton struct {
	targetPageID string
	x, y         float32
	label        string
}

func GetFilePath(fileName string) string {
	dir := filepath.Dir(filePathTxt)
	return filepath.Join(dir, Assets_Relative_Path, fileName)
}

type Game struct {
	navigator     *Navigator
	lastKeyState  map[ebiten.Key]bool
	leftColumnMsg string
}

func NewGame(navigator *Navigator) *Game {
	return &Game{
		navigator:     navigator,
		lastKeyState:  make(map[ebiten.Key]bool),
		leftColumnMsg: "Static Left Column",
	}
}

func (g *Game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && !g.lastKeyState[ebiten.KeyArrowLeft] {
		g.navigator.Pop()
		g.lastKeyState[ebiten.KeyArrowLeft] = true
	}

	if !ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.lastKeyState[ebiten.KeyArrowLeft] = false
	}

	navigatorOffsetX := float32(leftColumnWidth)
	navigatorOffsetY := float32(0)

	_, err := g.navigator.Update(navigatorOffsetX, navigatorOffsetY)
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	leftColumn := ebiten.NewImage(leftColumnWidth, screenHeight)
	leftColumn.Fill(color.RGBA{50, 50, 50, 255})
	ebitenutil.DebugPrintAt(leftColumn, g.leftColumnMsg, 10, 10)
	screen.DrawImage(leftColumn, nil)

	navigatorAreaRect := image.Rect(leftColumnWidth, 0, screenWidth, screenHeight)

	g.navigator.Draw(screen, navigatorAreaRect)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

type Navigator struct {
	stack      []*Page
	animating  bool
	transition float64
	direction  int
	Push       func(*Page)
	Pop        func()
}

func NewNavigator() *Navigator {
	return &Navigator{
		stack:      []*Page{},
		animating:  false,
		transition: 1.0,
	}
}

func (n *Navigator) PushPage(page *Page) {
	if len(n.stack) > 0 {
		n.animating = true
		n.transition = 0.0
		n.direction = 1
	}
	n.stack = append(n.stack, page)
}

func (n *Navigator) PopPage() {
	if len(n.stack) > 1 && !n.animating {
		n.animating = true
		n.transition = 0.0
		n.direction = -1
	}
}

func (n *Navigator) Update(navigatorOffsetX, navigatorOffsetY float32) (bool, error) {
	if len(n.stack) == 0 {
		return false, nil
	}

	if n.animating {
		n.transition += 0.05
		if n.transition >= 1.0 {
			n.transition = 1.0
			n.animating = false
			if n.direction == -1 {

				n.stack = n.stack[:len(n.stack)-1]
			}
		}
	}

	currentPage := n.stack[len(n.stack)-1]
	err := currentPage.Update(navigatorOffsetX, navigatorOffsetY, n.animating)
	return n.animating, err
}

func (n *Navigator) Draw(screen *ebiten.Image, navigatorAreaRect image.Rectangle) {

	navigatorArea := ebiten.NewImage(navigatorAreaRect.Dx(), navigatorAreaRect.Dy())
	navigatorArea.Fill(color.RGBA{30, 30, 30, 255})

	if n.animating && len(n.stack) > 1 {
		var prevOffsetX, currentOffsetX float64

		if n.direction == 1 {

			prevOffsetX = -n.transition * float64(navigatorAreaRect.Dx())

			currentOffsetX = float64(navigatorAreaRect.Dx()) * (1.0 - n.transition)
		} else if n.direction == -1 {

			currentOffsetX = float64(navigatorAreaRect.Dx()) * n.transition

			prevOffsetX = -float64(navigatorAreaRect.Dx()) * (1.0 - n.transition)
		}

		//fmt.Printf("PrevOffsetX: %f, CurrentOffsetX: %f\n", prevOffsetX, currentOffsetX)

		previousPage := n.stack[len(n.stack)-2]
		previousPage.Draw(navigatorArea, prevOffsetX, 0)

		currentPage := n.stack[len(n.stack)-1]
		currentPage.Draw(navigatorArea, currentOffsetX, 0)
	} else {

		if len(n.stack) > 0 {
			currentPage := n.stack[len(n.stack)-1]
			currentPage.Draw(navigatorArea, 0, 0)
		}
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(navigatorAreaRect.Min.X), float64(navigatorAreaRect.Min.Y))
	screen.DrawImage(navigatorArea, op)
}

func (n *Navigator) CurrentPage() *Page {
	if len(n.stack) == 0 {
		return nil
	}
	return n.stack[len(n.stack)-1]
}

type UIElement interface {
	Update(offsetX, offsetY float32, isAnimating bool)
	Draw(screen *ebiten.Image)
}

type Page struct {
	X, Y            float32
	Width, Height   float32
	backgroundColor color.Color
	message         string
	elements        []UIElement
	textWrapper     *textwrapper.TextWrapper
	NextPageID      string
}

func NewPage(
	bgColor color.Color,
	message string,
	tw *textwrapper.TextWrapper,
	x, y, width, height float32,
) *Page {
	return &Page{
		X:               x,
		Y:               y,
		Width:           width,
		Height:          height,
		backgroundColor: bgColor,
		message:         message,
		elements:        make([]UIElement, 0),
		textWrapper:     tw,
		NextPageID:      "",
	}
}

func (p *Page) AddElement(element UIElement) {
	p.elements = append(p.elements, element)
}

func (p *Page) Update(navigatorOffsetX, navigatorOffsetY float32, isAnimating bool) error {
	for _, element := range p.elements {
		element.Update(navigatorOffsetX+p.X, navigatorOffsetY+p.Y, isAnimating)
	}
	return nil
}

func (p *Page) Draw(navigatorArea *ebiten.Image, offsetX, offsetY float64) {
	pageArea := ebiten.NewImage(int(p.Width), int(p.Height))
	pageArea.Fill(p.backgroundColor)
	ebitenutil.DebugPrintAt(pageArea, p.message, 10, 10)

	for _, element := range p.elements {
		element.Draw(pageArea)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.X)+offsetX, float64(p.Y)+offsetY)
	navigatorArea.DrawImage(pageArea, op)
}

func (page *Page) addButton(btn PageButton) {
	targetPage, exists := pageRegistry[btn.targetPageID]
	if !exists {
		log.Fatalf("Page with ID '%s' does not exist", btn.targetPageID)
	}

	onNext := func(target *Page) func() {
		return func() {
			if globalNavigator == nil {
				log.Println("Navigator not initialized")
				return
			}
			globalNavigator.PushPage(target)
		}
	}(targetPage)

	button := widgets.NewButtonStd(
		btn.x,
		btn.y,
		100,
		40,
		btn.label,
		page.textWrapper,
		color.RGBA{0, 128, 255, 255},
		color.White,
		16,
		onNext,
	)

	page.AddElement(button)
}

func registerPages(tw *textwrapper.TextWrapper, navigatorAreaWidth, navigatorAreaHeight float32) {

	pageConfigs := []PageConfig{
		{
			id:      PageIDRed,
			color:   color.RGBA{255, 0, 0, 255},
			message: "Red Page - Full Size",
			x:       0,
			y:       0,
			width:   navigatorAreaWidth,
			height:  navigatorAreaHeight,
		},
		{
			id:      PageIDBlue,
			color:   color.RGBA{0, 0, 255, 255},
			message: "Blue Page - Top Half",
			x:       0,
			y:       0,
			width:   navigatorAreaWidth,
			height:  navigatorAreaHeight / 2,
		},
		{
			id:      PageIDGreen,
			color:   color.RGBA{0, 255, 0, 255},
			message: "Green Page - Bottom Half",
			x:       0,
			y:       navigatorAreaHeight / 2,
			width:   navigatorAreaWidth,
			height:  navigatorAreaHeight / 2,
		},
		{
			id:      PageIDYellow,
			color:   color.RGBA{255, 255, 0, 255},
			message: "Yellow Page - Smaller Window",
			x:       navigatorAreaWidth / 4,
			y:       navigatorAreaHeight / 4,
			width:   navigatorAreaWidth / 2,
			height:  navigatorAreaHeight / 2,
		},
		{
			id:      PageIDMagenta,
			color:   color.RGBA{255, 0, 255, 255},
			message: "Magenta Page - Custom Size",
			x:       50,
			y:       50,
			width:   300,
			height:  200,
		},
	}

	for _, cfg := range pageConfigs {
		page := NewPage(
			cfg.color,
			cfg.message,
			tw,
			cfg.x,
			cfg.y,
			cfg.width,
			cfg.height,
		)
		pageRegistry[cfg.id] = page
	}

	pageRegistry[PageIDRed].addButton(PageButton{
		targetPageID: PageIDBlue,
		x:            (navigatorAreaWidth - 100) / 2,
		y:            navigatorAreaHeight - 60,
		label:        "To Blue",
	})
	pageRegistry[PageIDRed].addButton(PageButton{
		targetPageID: PageIDYellow,
		x:            20,
		y:            20,
		label:        "To Yellow",
	})

	pageRegistry[PageIDBlue].addButton(PageButton{
		targetPageID: PageIDGreen,
		x:            (navigatorAreaWidth - 100) / 2,
		y:            (navigatorAreaHeight / 2) - 60,
		label:        "To Green",
	})
	pageRegistry[PageIDBlue].addButton(PageButton{
		targetPageID: PageIDMagenta,
		x:            navigatorAreaWidth - 120,
		y:            20,
		label:        "To Magenta",
	})

	pageRegistry[PageIDGreen].addButton(PageButton{
		targetPageID: PageIDYellow,
		x:            (navigatorAreaWidth - 100) / 2,
		y:            (navigatorAreaHeight / 2) - 60,
		label:        "To Yellow",
	})
	pageRegistry[PageIDYellow].addButton(PageButton{
		targetPageID: PageIDMagenta,
		x:            (navigatorAreaWidth/4 - 100) / 2,
		y:            (navigatorAreaHeight / 2) - 60,
		label:        "To Magenta",
	})
	pageRegistry[PageIDMagenta].addButton(PageButton{
		targetPageID: PageIDRed,
		x:            100,
		y:            140,
		label:        "To Red",
	})
}

func main() {
	_, filePathTxt, _, _ = runtime.Caller(0)
	fontPath := GetFilePath("assets/fonts/roboto_regularTTF.ttf")

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Navigator Example with Animations")

	navigator := NewNavigator()
	globalNavigator = navigator

	textWrapper, err := textwrapper.NewTextWrapper(fontPath, 16, false)
	if err != nil {
		log.Fatalf("Failed to create TextWrapper: %v", err)
	}

	navigatorAreaWidth := float32(screenWidth - leftColumnWidth)
	navigatorAreaHeight := float32(screenHeight)

	registerPages(textWrapper, navigatorAreaWidth, navigatorAreaHeight)

	navigator.Push = func(page *Page) {
		navigator.PushPage(page)
	}
	navigator.Pop = func() {
		navigator.PopPage()
	}

	initialPage, exists := pageRegistry[PageIDRed]
	if !exists {
		log.Fatal("Initial page 'redPage' not found in registry")
	}
	navigator.PushPage(initialPage)

	game := NewGame(navigator)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
