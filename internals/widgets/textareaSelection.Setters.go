package widgets

func (t *TextAreaSelection) setCursorPos(pos int) {
	t.cursorPos = clamp(pos, 0, len(t.text))
}

func (t *TextAreaSelection) setSelectionStart(pos int) {
	t.selectionStart = pos
	t.updateSelectionBounds()
}

func (t *TextAreaSelection) setSelectionEnd(pos int) {
	t.selectionEnd = pos
	t.updateSelectionBounds()
}

func (t *TextAreaSelection) SetScrollOffset(offset int) {
	t.scrollOffset = offset
	//t.updateSelectionBounds()
}
