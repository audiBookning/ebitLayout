package pagemodel

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more04/navigator"
	"example.com/menu/cmd02/more04/responsive"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type SidebarPageBase struct {
	ID            string
	Label         string
	MainUI        *responsive.UI
	SidebarUI     *responsive.UI
	SubNavigator  *navigator.Navigator
	PrevWidth     int
	PrevHeight    int
	SidebarWidth  int
	Navigator     *navigator.Navigator
	BackgroundClr color.Color
}

func NewSidebarPage(mainNav *navigator.Navigator, screenWidth, screenHeight int, id string, label string) *SidebarPageBase {

	subNav := navigator.NewNavigator(nil)

	level01 := NewSubPage(subNav, screenWidth, screenHeight)
	level02 := NewSubPage(subNav, screenWidth, screenHeight)

	subNav.AddPage("sub01", level01)
	subNav.AddPage("sub02", level02)

	subNav.SwitchTo("sub01")

	mainBreakpoints := []responsive.Breakpoint{
		{Width: 1200, LayoutMode: responsive.LayoutGrid},
		{Width: 800, LayoutMode: responsive.LayoutVertical},
		{Width: 0, LayoutMode: responsive.LayoutHorizontal},
	}
	mainButtons := []*responsive.Button{}

	const sub01PageID = "sub01"
	const sub02PageID = "sub02"
	const backPageID = "main"

	const sub01ButtonText = "Sub 1"
	const sub02ButtonText = "Sub 2"
	const backButtonText = "Back"

	sidebarButtons := []*responsive.Button{
		responsive.NewButton(sub01ButtonText, func() { subNav.SwitchTo(sub01PageID) }),
		responsive.NewButton(sub02ButtonText, func() { subNav.SwitchTo(sub02PageID) }),
		responsive.NewButton(backButtonText, func() { mainNav.SwitchTo(backPageID) }),
	}

	mainUI := responsive.NewUI("Start Game", mainBreakpoints, mainButtons)

	sidebarBreakpoints := []responsive.Breakpoint{
		{Width: 0, LayoutMode: responsive.LayoutVertical},
	}

	sidebarUI := responsive.NewUI(label, sidebarBreakpoints, sidebarButtons)

	const sidebarFixedWidth = 200

	mainUI.Update(screenWidth-sidebarFixedWidth, screenHeight)
	sidebarUI.Update(sidebarFixedWidth, screenHeight)

	page := &SidebarPageBase{
		ID:            "sidebar",
		Label:         "Sidebar",
		MainUI:        mainUI,
		SidebarUI:     sidebarUI,
		SubNavigator:  subNav,
		PrevWidth:     screenWidth,
		PrevHeight:    screenHeight,
		SidebarWidth:  sidebarFixedWidth,
		Navigator:     mainNav,
		BackgroundClr: color.RGBA{0x3E, 0x3E, 0x3E, 0xFF},
	}

	page.ResetAllButtonStates()

	return page
}

func (p *SidebarPageBase) AddSubPages(subPages ...*SubPage) {
	for _, subPage := range subPages {
		p.SubNavigator.AddPage(subPage.ID, subPage)
	}
}

func (p *SidebarPageBase) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth != p.PrevWidth || outsideHeight != p.PrevHeight {
		log.Printf("SidebarPage: Window resized to %dx%d\n", outsideWidth, outsideHeight)
		p.PrevWidth = outsideWidth
		p.PrevHeight = outsideHeight
		//p.ui.Update(p.prevWidth, p.prevHeight)
		p.MainUI.Update(p.PrevWidth-p.SidebarWidth, p.PrevHeight)
		p.SidebarUI.Update(p.SidebarWidth, p.PrevHeight)
	}
	return outsideWidth, outsideHeight
}

func (p *SidebarPageBase) Update() error {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		p.HandleInput(x, y)
	}

	p.SubNavigator.CurrentActivePage().Update()

	return nil
}

func (p *SidebarPageBase) HandleInput(x, y int) {
	if x < p.SidebarWidth {
		p.SidebarUI.HandleClick(x, y)
	} else {

		if p.SubNavigator.CurrentActivePage() != nil {
			p.SubNavigator.CurrentActivePage().HandleInput(x-p.SidebarWidth, y)
		}
	}
}

func (p *SidebarPageBase) DrawBackGround(screen *ebiten.Image) {
	screen.Fill(p.BackgroundClr)
}

func (p *SidebarPageBase) Draw(screen *ebiten.Image) {
	p.DrawBackGround(screen)

	p.SidebarUI.Draw(screen)
	p.MainUI.Draw(screen)

	if p.SubNavigator.CurrentActivePage() != nil {
		screenWidth, screenHeight := screen.Size()

		playRenderSpace := ebiten.NewImage(screenWidth-p.SidebarWidth, screenHeight)
		p.SubNavigator.CurrentActivePage().Draw(playRenderSpace)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(p.SidebarWidth), 0)
		screen.DrawImage(playRenderSpace, op)
	}

	separatorColor := color.RGBA{0x00, 0x00, 0x00, 0xFF}
	separatorImg := ebiten.NewImage(2, p.PrevHeight)
	separatorImg.Fill(separatorColor)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.SidebarWidth), 0)
	screen.DrawImage(separatorImg, op)
}

func (p *SidebarPageBase) ResetAllButtonStates() {
	p.MainUI.ResetButtonStates()
	p.SidebarUI.ResetButtonStates()
	if p.SubNavigator.CurrentActivePage() != nil {
		p.SubNavigator.CurrentActivePage().ResetButtonStates()
	}
}

func (p *SidebarPageBase) ResetButtonStates() {
	p.MainUI.ResetButtonStates()
	p.SidebarUI.ResetButtonStates()
	if p.SubNavigator.CurrentActivePage() != nil {
		p.SubNavigator.CurrentActivePage().ResetButtonStates()
	}
}
