package widgets

type SelectionBounds struct {
	selectionStart  int
	selectionEnd    int
	isSelecting     bool
	minSelectionPos int
	maxSelectionPos int
}

// NewSelection initializes a new Selection instance
func NewSelectionBounds() *SelectionBounds {
	return &SelectionBounds{
		selectionStart:  0,
		selectionEnd:    0,
		isSelecting:     false,
		minSelectionPos: 0,
		maxSelectionPos: 0,
	}
}

func (t *SelectionBounds) getSelectionBoundsStart() int {
	minPos, _ := t.getSelectionBounds()
	return minPos
}

// getSelectionBoundsEnd returns the end position of the current selection
func (t *SelectionBounds) getSelectionBoundsEnd() int {
	_, maxPos := t.getSelectionBounds()
	return maxPos
}

func (t *SelectionBounds) updateSelectionBounds() {
	if t.selectionStart <= t.selectionEnd {
		t.minSelectionPos = t.selectionStart
		t.maxSelectionPos = t.selectionEnd
	} else {
		t.minSelectionPos = t.selectionEnd
		t.maxSelectionPos = t.selectionStart
	}
}

func (t *SelectionBounds) getSelectionBounds() (int, int) {
	t.updateSelectionBounds()
	return t.minSelectionPos, t.maxSelectionPos
}

func (t *SelectionBounds) setSelectionStart(pos int) {
	t.selectionStart = pos
	t.updateSelectionBounds()
}

func (t *SelectionBounds) setSelectionEnd(pos int) {
	t.selectionEnd = pos
	t.updateSelectionBounds()
}

func (t *SelectionBounds) ClearSelection(pos int) {
	t.setSelectionStart(pos)
	t.setSelectionEnd(pos)
}

func (t *SelectionBounds) SetIsSelecting(isSelecting bool) {
	t.isSelecting = isSelecting
}
