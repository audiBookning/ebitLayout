package widgets

func (t *TextAreaSelection) setCursorPos(pos int) {
	t.cursorPos = clamp(pos, 0, len(t.text))
}

func (t *TextAreaSelection) setSelectionStart(pos int) {
	t.selectionStart = clamp(pos, 0, len(t.text))
	t.updateSelectionBounds()
}

func (t *TextAreaSelection) setSelectionEnd(pos int) {
	t.selectionEnd = clamp(pos, 0, len(t.text))
	t.updateSelectionBounds()
}

func (t *TextAreaSelection) SetScrollOffset(offset int) {
	t.scrollOffset = offset
	//t.updateSelectionBounds()
}

func (t *TextAreaSelection) SetIsSelecting(isSelecting bool) {
	t.isSelecting = isSelecting
}

func (t *TextAreaSelection) SetIsDraggingThumb(isDragging bool) {
	t.isDraggingThumb = isDragging
}
