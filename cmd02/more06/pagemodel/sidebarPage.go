package pagemodel

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more06/navigator"
	"example.com/menu/cmd02/more06/responsive"
	"example.com/menu/cmd02/more06/textwrapper"
	"example.com/menu/cmd02/more06/types"
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

func NewSidebarPageBase(mainNav *navigator.Navigator, textWrapper *textwrapper.TextWrapper, id, label string, screenWidth, screenHeight int) *SidebarPageBase {

	subNav := navigator.NewNavigator(nil)

	sub1 := NewSubPageBase(textWrapper, "sub01", "Sub 1", screenWidth, screenHeight)
	sub2 := NewSubPageBase(textWrapper, "sub02", "Sub 2", screenWidth, screenHeight)

	subNav.AddPage(sub1.ID, sub1)
	subNav.AddPage(sub2.ID, sub2)

	subNav.SwitchTo("sub01")

	mainBreakpoints := []responsive.Breakpoint{
		{Width: 1200, LayoutMode: responsive.LayoutGrid},
		{Width: 800, LayoutMode: responsive.LayoutVertical},
		{Width: 0, LayoutMode: responsive.LayoutHorizontal},
	}
	mainFields := []types.Element{}
	mainUI := responsive.NewUI(label, mainBreakpoints, mainFields, textWrapper, responsive.AlignCenter)

	sidebarBreakpoints := []responsive.Breakpoint{
		{Width: 0, LayoutMode: responsive.LayoutVertical},
	}
	sidebarFields := []types.Element{
		responsive.NewButton("Sub 1", func() { subNav.SwitchTo("sub01") }, textWrapper),
		responsive.NewButton("Sub 2", func() { subNav.SwitchTo("sub02") }, textWrapper),
		responsive.NewButton("Back", func() { mainNav.SwitchTo("main") }, textWrapper),
	}

	sidebarUI := responsive.NewUI("Sidebar Menu", sidebarBreakpoints, sidebarFields, textWrapper, responsive.AlignCenter)

	const sidebarFixedWidth = 200
	mainUI.LayoutUpdate(screenWidth-sidebarFixedWidth, screenHeight)
	sidebarUI.LayoutUpdate(sidebarFixedWidth, screenHeight)

	page := &SidebarPageBase{
		ID:            id,
		Label:         label,
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

func (p *SidebarPageBase) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth != p.PrevWidth || outsideHeight != p.PrevHeight {
		log.Printf("SidebarPageBase: Window resized to %dx%d\n", outsideWidth, outsideHeight)
		p.PrevWidth = outsideWidth
		p.PrevHeight = outsideHeight
		p.MainUI.LayoutUpdate(p.PrevWidth-p.SidebarWidth, p.PrevHeight)
		p.SidebarUI.LayoutUpdate(p.SidebarWidth, p.PrevHeight)
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

func (p *SidebarPageBase) Draw(screen *ebiten.Image) {
	p.DrawBackGround(screen)

	p.SidebarUI.Draw(screen)
	p.MainUI.Draw(screen)

	if p.SubNavigator.CurrentActivePage() != nil {
		screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()

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

func (p *SidebarPageBase) DrawBackGround(screen *ebiten.Image) {
	screen.Fill(p.BackgroundClr)
}

func (p *SidebarPageBase) ResetAllButtonStates() {
	p.ResetFieldStates()
	if p.SubNavigator.CurrentActivePage() != nil {
		p.SubNavigator.CurrentActivePage().ResetFieldStates()
	}
}

func (p *SidebarPageBase) ResetFieldStates() {
	p.MainUI.ResetFieldStates()
	p.SidebarUI.ResetFieldStates()
}
