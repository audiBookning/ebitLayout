package widgets

import (
	"example.com/menu/internals/textwrapper"
	//"example.com/menu/internals/textwrapper02"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.design/x/clipboard"
)

// KeyState tracks the repeat state of a specific key
type KeyState struct {
	InitialPress    bool // Indicates if the initial press has been handled
	FramesHeld      int  // Number of frames the key has been held down
	FramesUntilNext int  // Frames remaining until the next action
}

type TextState struct {
	Text      string
	CursorPos int
}

type TextArea struct {
	textWrapper *textwrapper.TextWrapper
	text        string
	selection   *SelectionBounds
	hasFocus    bool
	cursorPos   int
	counter     int
	//selectionStart       int
	//selectionEnd         int
	//isSelecting          bool
	x, y, w, h           int
	maxLines             int
	cursorBlinkRate      int
	tabWidth             int
	lineHeight           float64
	heldKeys             map[ebiten.Key]*KeyState
	undoStack            []TextState
	redoStack            []TextState
	desiredCursorCol     int
	lastClickTime        int  // Frame count of the last click
	clickCount           int  // Number of consecutive clicks
	doubleClickThreshold int  // Threshold frames to consider as a double-click
	doubleClickHandled   bool // Indicates if a double-click has just been handled
	scrollOffset         int

	isMouseLeftPressed bool

	// Scrollbar Fields
	scrollbarWidth  int     // Width of the scrollbar
	scrollbarX      int     // X position of the scrollbar
	scrollbarY      int     // Y position of the scrollbar
	scrollbarHeight int     // Height of the scrollbar
	scrollbarThumbY float64 // Y position of the scrollbar thumb
	scrollbarThumbH float64 // Height of the scrollbar thumb
	isDraggingThumb bool    // Indicates if the scrollbar thumb is being dragged
	dragOffsetY     float64 // Offset between mouse position and thumb position during drag

	// Key repeat constants
	keyRepeatInitialDelay int
	keyRepeatInterval     int
	// performance
	cachedLines   []string
	isTextChanged bool
	//minSelectionPos int
	//maxSelectionPos int
	// Minimum movement to consider as drag
	paddingLeft   int
	paddingTop    int
	paddingBottom int
	clicked       bool

	stepX float64
	stepY float64
}

// func NewTextAreaSelection(textWrapper *textwrapper02.TextWrapper, x, y, w, h int, startTxt string) *TextAreaSelection {
func NewTextAreaSelection(textWrapper *textwrapper.TextWrapper, x, y, w, h int, startTxt string) *TextArea {
	err := clipboard.Init()
	if err != nil {
		fmt.Println("Clipboard initialization failed:", err)
		return nil
	}

	// Calculate line height based on font metrics
	metrics := textWrapper.GetFontMetrics()
	/*
		lineHeight := float64(metrics.Height)
	*/
	lineHeight := metrics.HAscent + metrics.HDescent + metrics.HLineGap
	//lineHeight := textWrapper.MeasureTextHeightWrap(startTxt)
	//_, lineHeight := textWrapper.MeasureText(startTxt)

	//monospaceWidth := textWrapper.GetMonospaceWidth()
	monospaceWidth, _ := textWrapper.MeasureString("s")

	// Calculate maxLines based on the height of the TextAreaSelection and the line height
	padding := 10
	maxLines := int((h - 2*padding) / int(lineHeight))

	return &TextArea{
		textWrapper:          textWrapper,
		selection:            NewSelectionBounds(),
		x:                    x,
		y:                    y,
		w:                    w,
		h:                    h,
		maxLines:             maxLines, // Use the calculated maxLines
		cursorBlinkRate:      30,
		tabWidth:             4,
		lineHeight:           float64(lineHeight),
		heldKeys:             make(map[ebiten.Key]*KeyState),
		desiredCursorCol:     -1,
		lastClickTime:        0,     // Frame count of the last click
		clickCount:           0,     // Number of consecutive clicks
		doubleClickThreshold: 30,    // Threshold frames to consider as a double-click
		doubleClickHandled:   false, // Indicates if a double-click has just been handled
		scrollOffset:         0,
		text:                 startTxt, // Default text added here

		keyRepeatInitialDelay: 30,
		keyRepeatInterval:     5,
		isTextChanged:         true,
		paddingLeft:           padding,
		paddingTop:            padding,
		paddingBottom:         padding,
		clicked:               false,
		isMouseLeftPressed:    false,

		stepX: monospaceWidth,
		stepY: float64(lineHeight),
	}
}
