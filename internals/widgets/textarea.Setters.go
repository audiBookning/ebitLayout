package widgets

func (t *TextArea) setCursorPos(pos int) {
	t.cursorPos = clamp(pos, 0, len(t.text))
}

func (t *TextArea) SetScrollOffset(offset int) {
	t.scrollOffset = offset
	//t.updateSelectionBounds()
}

func (t *TextArea) SetIsDraggingThumb(isDragging bool) {
	t.isDraggingThumb = isDragging
}
