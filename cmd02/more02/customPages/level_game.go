package customPages

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"example.com/menu/cmd02/more02/responsive"
	"example.com/menu/cmd02/more02/types"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type LevelGamePage struct {
	mainUI       *responsive.UI
	sidebarUI    *responsive.UI
	subPages     map[string]types.Page
	currentSub   types.Page
	prevWidth    int
	prevHeight   int
	sidebarWidth int
	switchPage   func(string)
}

const sidebarFixedWidth = 200

func NewLevelGamePage(switchPage func(string)) *LevelGamePage {

	mainBreakpoints := []responsive.Breakpoint{
		{Width: 1200, LayoutMode: responsive.LayoutGrid},
		{Width: 800, LayoutMode: responsive.LayoutVertical},
		{Width: 0, LayoutMode: responsive.LayoutHorizontal},
	}
	mainButtons := []*responsive.Button{}

	sidebarButtons := []*responsive.Button{
		responsive.NewButton("Level 1", nil),
		responsive.NewButton("Level 2", nil),
		responsive.NewButton("Back", nil),
	}

	mainUI := responsive.NewUI("Start Game", mainBreakpoints, mainButtons)

	sidebarBreakpoints := []responsive.Breakpoint{
		{Width: 0, LayoutMode: responsive.LayoutVertical},
	}

	sidebarUI := responsive.NewUI("Menu", sidebarBreakpoints, sidebarButtons)

	screenWidth, screenHeight := 800, 600
	mainUI.Update(screenWidth-sidebarFixedWidth, screenHeight)
	sidebarUI.Update(sidebarFixedWidth, screenHeight)

	subPages := make(map[string]types.Page)
	subPages["level01"] = NewLevel01Page()
	subPages["level02"] = NewLevel02Page()

	page := &LevelGamePage{
		mainUI:       mainUI,
		sidebarUI:    sidebarUI,
		subPages:     subPages,
		currentSub:   subPages["level01"],
		prevWidth:    screenWidth,
		prevHeight:   screenHeight,
		sidebarWidth: sidebarFixedWidth,
		switchPage:   switchPage,
	}

	page.setupSidebarButtons()

	page.ResetAllButtonStates()

	return page
}

func (p *LevelGamePage) setupSidebarButtons() {
	p.sidebarUI.Buttons[0].OnClickFunc = func() { p.SwitchSubPage("level01") }
	p.sidebarUI.Buttons[1].OnClickFunc = func() { p.SwitchSubPage("level02") }
	p.sidebarUI.Buttons[2].OnClickFunc = func() {
		log.Println("Back clicked")
		p.switchPage("main")
	}
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

	if p.currentSub != nil {
		p.currentSub.Update()
	}

	return nil
}

func (p *LevelGamePage) HandleInput(x, y int) {
	if x < p.sidebarWidth {
		p.sidebarUI.HandleClick(x, y)
	} else {

		if p.currentSub != nil {
			p.currentSub.HandleInput(x, y)
		}
	}
}

func (p *LevelGamePage) SwitchSubPage(pageName string) {
	if page, exists := p.subPages[pageName]; exists {
		log.Printf("Switching to subpage: %s\n", pageName)
		p.currentSub = page
		p.ResetAllButtonStates()
	} else {
		log.Printf("Subpage %s does not exist!\n", pageName)
	}
}

func (p *LevelGamePage) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{0x3E, 0x3E, 0x3E, 0xFF})

	p.sidebarUI.Draw(screen)
	p.mainUI.Draw(screen)

	if p.currentSub != nil {
		screenWidth, screenHeight := screen.Size()

		playRenderSpace := ebiten.NewImage(screenWidth-p.sidebarWidth, screenHeight)
		p.currentSub.Draw(playRenderSpace)

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
	for _, subPage := range p.subPages {
		if subPageWithUI, ok := subPage.(interface{ ResetButtonStates() }); ok {
			subPageWithUI.ResetButtonStates()
		}
	}
}
