package customPages

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more03/navigator"
	"example.com/menu/cmd02/more03/responsive"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type LevelGamePage struct {
	mainUI       *responsive.UI
	sidebarUI    *responsive.UI
	subNavigator *navigator.Navigator
	prevWidth    int
	prevHeight   int
	sidebarWidth int
	navigator    *navigator.Navigator
}

const sidebarFixedWidth = 200

func NewLevelGamePage(mainNav *navigator.Navigator, screenWidth, screenHeight int) *LevelGamePage {

	subNav := navigator.NewNavigator(nil)

	level01 := NewLevel01Page(subNav, screenWidth, screenHeight)
	level02 := NewLevel02Page(subNav, screenWidth, screenHeight)

	subNav.AddPage("level01", level01)
	subNav.AddPage("level02", level02)

	subNav.SwitchTo("level01")

	mainBreakpoints := []responsive.Breakpoint{
		{Width: 1200, LayoutMode: responsive.LayoutGrid},
		{Width: 800, LayoutMode: responsive.LayoutVertical},
		{Width: 0, LayoutMode: responsive.LayoutHorizontal},
	}
	mainButtons := []*responsive.Button{}

	sidebarButtons := []*responsive.Button{
		responsive.NewButton("Level 1", func() { subNav.SwitchTo("level01") }),
		responsive.NewButton("Level 2", func() { subNav.SwitchTo("level02") }),
		responsive.NewButton("Back", func() { mainNav.SwitchTo("main") }),
	}

	mainUI := responsive.NewUI("Start Game", mainBreakpoints, mainButtons)

	sidebarBreakpoints := []responsive.Breakpoint{
		{Width: 0, LayoutMode: responsive.LayoutVertical},
	}

	sidebarUI := responsive.NewUI("Menu", sidebarBreakpoints, sidebarButtons)

	mainUI.Update(screenWidth-sidebarFixedWidth, screenHeight)
	sidebarUI.Update(sidebarFixedWidth, screenHeight)

	page := &LevelGamePage{
		mainUI:       mainUI,
		sidebarUI:    sidebarUI,
		subNavigator: subNav,
		prevWidth:    screenWidth,
		prevHeight:   screenHeight,
		sidebarWidth: sidebarFixedWidth,
		navigator:    mainNav,
	}

	page.ResetAllButtonStates()

	return page
}

func (p *LevelGamePage) Layout(outsideWidth, outsideHeight int) (int, int) {
	p.prevWidth = outsideWidth
	p.prevHeight = outsideHeight
	return outsideWidth, outsideHeight
}

func (p *LevelGamePage) Update() error {
	screenWidth, screenHeight := ebiten.WindowSize()

	if screenWidth != p.prevWidth || screenHeight != p.prevHeight {
		log.Printf("LevelGamePage: Window resized to %dx%d\n", screenWidth, screenHeight)
		p.prevWidth = screenWidth
		p.prevHeight = screenHeight
	}

	p.mainUI.Update(screenWidth-p.sidebarWidth, screenHeight)
	p.sidebarUI.Update(p.sidebarWidth, screenHeight)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		p.HandleInput(x, y)
	}

	p.subNavigator.CurrentActivePage().Update()

	return nil
}

func (p *LevelGamePage) HandleInput(x, y int) {
	if x < p.sidebarWidth {
		p.sidebarUI.HandleClick(x, y)
	} else {

		if p.subNavigator.CurrentActivePage() != nil {
			p.subNavigator.CurrentActivePage().HandleInput(x-p.sidebarWidth, y)
		}
	}
}

func (p *LevelGamePage) DrawBackGround(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x3E, 0x3E, 0x3E, 0xFF})
}

func (p *LevelGamePage) Draw(screen *ebiten.Image) {
	p.DrawBackGround(screen)

	p.sidebarUI.Draw(screen)
	p.mainUI.Draw(screen)

	if p.subNavigator.CurrentActivePage() != nil {
		screenWidth, screenHeight := screen.Size()

		playRenderSpace := ebiten.NewImage(screenWidth-p.sidebarWidth, screenHeight)
		p.subNavigator.CurrentActivePage().Draw(playRenderSpace)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(p.sidebarWidth), 0)
		screen.DrawImage(playRenderSpace, op)
	}

	separatorColor := color.RGBA{0x00, 0x00, 0x00, 0xFF}
	separatorImg := ebiten.NewImage(2, p.prevHeight)
	separatorImg.Fill(separatorColor)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.sidebarWidth), 0)
	screen.DrawImage(separatorImg, op)
}

func (p *LevelGamePage) ResetAllButtonStates() {
	p.mainUI.ResetButtonStates()
	p.sidebarUI.ResetButtonStates()
	if p.subNavigator.CurrentActivePage() != nil {
		p.subNavigator.CurrentActivePage().ResetButtonStates()
	}
}

func (p *LevelGamePage) ResetButtonStates() {
	p.mainUI.ResetButtonStates()
	p.sidebarUI.ResetButtonStates()
	if p.subNavigator.CurrentActivePage() != nil {
		p.subNavigator.CurrentActivePage().ResetButtonStates()
	}
}
