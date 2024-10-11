package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"example.com/menu/internals/navigator"
	"example.com/menu/internals/page"
	"example.com/menu/internals/textwrapper"
	"example.com/menu/internals/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth     = 800
	screenHeight    = 600
	leftColumnWidth = 160
)

type Game struct {
	navigator     *navigator.Navigator
	lastKeyState  map[ebiten.Key]bool
	leftColumnMsg string
}

func NewGame(navigator *navigator.Navigator) *Game {
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

type RedPage struct {
	*page.BasePage
	rotationAngle float64
}

func NewRedPage(tw *textwrapper.TextWrapper, navigatorAreaWidth, navigatorAreaHeight float32) *RedPage {
	basePage := page.NewBasePage(
		color.RGBA{255, 0, 0, 255},
		"Red Page - Rotating",
		tw,
		0,
		0,
		navigatorAreaWidth,
		navigatorAreaHeight,
	)

	redPage := &RedPage{
		BasePage:      basePage,
		rotationAngle: 0,
	}

	return redPage
}

func (rp *RedPage) customDraw(screen *ebiten.Image) {
	rp.rotationAngle += 0.02

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(rp.Width)/2, -float64(rp.Height)/2)
	op.GeoM.Rotate(rp.rotationAngle)
	op.GeoM.Translate(float64(rp.Width)/2, float64(rp.Height)/2)
	screen.DrawImage(rp.PageArea, op)
}

type BluePage struct {
	*page.BasePage
	waveOffset float64
}

func NewBluePage(tw *textwrapper.TextWrapper, navigatorAreaWidth, navigatorAreaHeight float32) *BluePage {
	basePage := page.NewBasePage(
		color.RGBA{0, 0, 255, 255},
		"Blue Page - Wave",
		tw,
		0,
		0,
		navigatorAreaWidth,
		navigatorAreaHeight/2,
	)

	bluePage := &BluePage{
		BasePage:   basePage,
		waveOffset: 0,
	}

	bluePage.SetCustomDraw(bluePage.customDraw)

	return bluePage
}

func (bp *BluePage) customDraw(screen *ebiten.Image) {
	bp.waveOffset += 0.1

	for x := float64(0); x < float64(bp.Width); x++ {
		y := math.Sin(x*0.05+bp.waveOffset)*20 + float64(bp.Height)/2
		vector.StrokeLine(bp.PageArea, float32(x), float32(y), float32(x), float32(y+1), 2, color.White, false)
	}
}

type GreenPage struct {
	*page.BasePage
	particles []Particle
}

type Particle struct {
	X, Y   float64
	SpeedX float64
	SpeedY float64
	Radius float64
}

func NewGreenPage(tw *textwrapper.TextWrapper, navigatorAreaWidth, navigatorAreaHeight float32) *GreenPage {
	basePage := page.NewBasePage(
		color.RGBA{0, 255, 0, 255},
		"Green Page - Particles",
		tw,
		0,
		navigatorAreaHeight/2,
		navigatorAreaWidth,
		navigatorAreaHeight/2,
	)

	greenPage := &GreenPage{
		BasePage:  basePage,
		particles: make([]Particle, 20),
	}

	for i := range greenPage.particles {
		greenPage.particles[i] = Particle{
			X:      rand.Float64() * float64(greenPage.Width),
			Y:      rand.Float64() * float64(greenPage.Height),
			SpeedX: (rand.Float64() - 0.5) * 2,
			SpeedY: (rand.Float64() - 0.5) * 2,
			Radius: rand.Float64()*5 + 2,
		}
	}

	greenPage.SetCustomDraw(greenPage.customDraw)

	return greenPage
}

func (gp *GreenPage) customDraw(screen *ebiten.Image) {

	for i := range gp.particles {
		gp.particles[i].X += gp.particles[i].SpeedX
		gp.particles[i].Y += gp.particles[i].SpeedY

		if gp.particles[i].X < 0 || gp.particles[i].X > float64(gp.Width) {
			gp.particles[i].SpeedX *= -1
		}
		if gp.particles[i].Y < 0 || gp.particles[i].Y > float64(gp.Height) {
			gp.particles[i].SpeedY *= -1
		}
	}

	for _, p := range gp.particles {
		vector.DrawFilledCircle(
			gp.PageArea,
			float32(p.X),
			float32(p.Y),
			float32(p.Radius),
			color.White,
			false,
		)
	}
}

type YellowPage struct {
	*page.BasePage
	textScale float64
}

func NewYellowPage(tw *textwrapper.TextWrapper, navigatorAreaWidth, navigatorAreaHeight float32) *YellowPage {
	basePage := page.NewBasePage(
		color.RGBA{255, 255, 0, 255},
		"Yellow Page - Pulsating Text",
		tw,
		navigatorAreaWidth/4,
		navigatorAreaHeight/4,
		navigatorAreaWidth/2,
		navigatorAreaHeight/2,
	)

	yellowPage := &YellowPage{
		BasePage:  basePage,
		textScale: 1.0,
	}

	yellowPage.SetCustomDraw(yellowPage.customDraw)

	return yellowPage
}

func (yp *YellowPage) Update(navigatorOffsetX, navigatorOffsetY float32, isAnimating bool) error {
	yp.textScale = 1 + 0.2*math.Sin(float64(time.Now().UnixNano())/1e9)
	return yp.BasePage.Update(navigatorOffsetX, navigatorOffsetY, isAnimating)
}

func (yp *YellowPage) customDraw(screen *ebiten.Image) {

	scaledMessage := fmt.Sprintf("Scale: %.2f", yp.textScale)
	textWidth, _ := yp.TextWrapper.MeasureText(scaledMessage)
	x := float64(yp.Width)/2 - textWidth/2
	y := float64(yp.Height) - 30

	tw := yp.TextWrapper
	tw.ResetGeom()
	tw.SetGeomScale(yp.textScale, yp.textScale)
	tw.Color = color.Black
	tw.DrawText(screen, scaledMessage, x, y)
}

type MagentaPage struct {
	*page.BasePage
	gradient    *ebiten.Image
	needsUpdate bool
}

func NewMagentaPage(tw *textwrapper.TextWrapper, navigatorAreaWidth, navigatorAreaHeight float32) *MagentaPage {
	basePage := page.NewBasePage(
		color.RGBA{255, 0, 255, 255},
		"Magenta Page - Gradient",
		tw,
		50,
		50,
		300,
		200,
	)

	magentaPage := &MagentaPage{
		BasePage:    basePage,
		gradient:    ebiten.NewImage(int(basePage.Width), int(basePage.Height)),
		needsUpdate: true,
	}

	magentaPage.updateGradient()

	basePage.SetCustomDrawBackground(magentaPage.DrawBackground)

	//magentaPage.SetCustomDraw(magentaPage.customDraw)

	return magentaPage
}

func (mp *MagentaPage) updateGradient() {
	for y := 0; y < int(mp.Height); y++ {
		for x := 0; x < int(mp.Width); x++ {
			r := uint8(float64(x) / float64(mp.Width) * 255)
			b := uint8(float64(y) / float64(mp.Height) * 255)
			mp.gradient.Set(x, y, color.RGBA{r, 0, b, 255})
		}
	}
}

func (mp *MagentaPage) customDraw(screen *ebiten.Image) {
	if mp.needsUpdate {
		mp.updateGradient()
		mp.needsUpdate = false
	}

	mp.PageArea.DrawImage(mp.gradient, nil)
}

func (mp *MagentaPage) Update(navigatorOffsetX, navigatorOffsetY float32, isAnimating bool) error {
	if mp.needsUpdate {
		mp.updateGradient()
		mp.needsUpdate = false
	}
	return mp.BasePage.Update(navigatorOffsetX, navigatorOffsetY, isAnimating)
}

func (mp *MagentaPage) DrawBackground(screen *ebiten.Image) {

	screen.DrawImage(mp.gradient, nil)
}

func (mp *MagentaPage) GetType() string {
	return "MagentaPage"
}

func registerPages(
	tw *textwrapper.TextWrapper,
	navigatorAreaWidth, navigatorAreaHeight float32,
	navigator *navigator.Navigator) page.Page {

	redPage := NewRedPage(tw, navigatorAreaWidth, navigatorAreaHeight)
	bluePage := NewBluePage(tw, navigatorAreaWidth, navigatorAreaHeight)
	greenPage := NewGreenPage(tw, navigatorAreaWidth, navigatorAreaHeight)
	yellowPage := NewYellowPage(tw, navigatorAreaWidth, navigatorAreaHeight)
	magentaPage := NewMagentaPage(tw, navigatorAreaWidth, navigatorAreaHeight)

	redPage.AddButton(page.PageButton{
		X:     (navigatorAreaWidth - 100) / 2,
		Y:     navigatorAreaHeight - 60,
		Label: "To Blue",
	}, func() { navigator.Push(bluePage) })

	redPage.AddButton(page.PageButton{
		X:     20,
		Y:     20,
		Label: "To Yellow",
	}, func() { navigator.Push(yellowPage) })

	bluePage.AddButton(page.PageButton{
		X:     (navigatorAreaWidth - 100) / 2,
		Y:     (navigatorAreaHeight / 2) - 60,
		Label: "To Green",
	}, func() { navigator.Push(greenPage) })

	bluePage.AddButton(page.PageButton{
		X:     navigatorAreaWidth - 120,
		Y:     20,
		Label: "To Magenta",
	}, func() { navigator.Push(magentaPage) })

	greenPage.AddButton(page.PageButton{
		X:     (navigatorAreaWidth - 100) / 2,
		Y:     (navigatorAreaHeight / 2) - 60,
		Label: "To Yellow",
	}, func() { navigator.Push(yellowPage) })

	yellowPage.AddButton(page.PageButton{
		X:     (navigatorAreaWidth/4 - 100) / 2,
		Y:     (navigatorAreaHeight / 2) - 60,
		Label: "To Magenta",
	}, func() { navigator.Push(magentaPage) })

	magentaPage.AddButton(page.PageButton{
		X:     100,
		Y:     140,
		Label: "To Red",
	}, func() { navigator.Push(redPage) })

	return redPage
}

func main() {
	utils.InitGetFilepath()
	fontPath := utils.GetFilePath("assets/fonts/roboto_regularTTF.ttf")

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Navigator Example with Animations")

	navigator := navigator.NewNavigator()

	textWrapper, err := textwrapper.NewTextWrapper(fontPath, 16, false)
	if err != nil {
		log.Fatalf("Failed to create TextWrapper: %v", err)
	}

	navigatorAreaWidth := float32(screenWidth - leftColumnWidth)
	navigatorAreaHeight := float32(screenHeight)

	initialPage := registerPages(textWrapper, navigatorAreaWidth, navigatorAreaHeight, navigator)

	navigator.Push(initialPage)

	game := NewGame(navigator)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
