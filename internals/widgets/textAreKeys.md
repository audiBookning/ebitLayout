Below is an exhaustive list of typical key interactions with a `<textarea>` element or any other text field in a browser. These interactions are common across most operating systems, although behavior might slightly vary between platforms (Windows, macOS, Linux).

## Implemented Functions

### Navigation Keys

1. **Arrow Keys**:
   - **Left Arrow**: Move the cursor one character to the left.
   - **Right Arrow**: Move the cursor one character to the right.
   - **Up Arrow**: Move the cursor up one line.
   - **Down Arrow**: Move the cursor down one line.

2. **Home/End**:
   - **Home**: Move the cursor to the beginning of the current line.
   - **End**: Move the cursor to the end of the current line.

3. **Ctrl + Arrow Keys** (or Command + Arrow Keys on macOS):
   - **Ctrl + Left Arrow**: Move the cursor one word to the left.
   - **Ctrl + Right Arrow**: Move the cursor one word to the right.

4. **Page Up/Page Down**:
   - **Page Up**: Scroll up the content by one visible page without moving the cursor.
   - **Page Down**: Scroll down the content by one visible page without moving the cursor.

5. **Ctrl + Home/End**:
- **Ctrl + Home**: Move the cursor to the very beginning of the text (top of the textarea).
- **Ctrl + End**: Move the cursor to the very end of the text (bottom of the textarea).

6. **Ctrl + Arrow Keys**:
   - **Ctrl + Up Arrow**: Move the cursor to the beginning of the current paragraph.
   - **Ctrl + Down Arrow**: Move the cursor to the end of the current paragraph.

### Selection Keys

1. **Shift + Arrow Keys**:
   - **Shift + Left Arrow**: Select one character to the left.
   - **Shift + Right Arrow**: Select one character to the right.
   - **Shift + Up Arrow**: Select text from the current line upwards.
   - **Shift + Down Arrow**: Select text from the current line downwards.

2. **Shift + Home/End**:
   - **Shift + Home**: Select from the current cursor position to the beginning of the line.
   - **Shift + End**: Select from the current cursor position to the end of the line.

3. **Ctrl + Shift + Home/End**:
   - **Ctrl + Shift + Home**: Select all text from the cursor to the beginning of the text.
   - **Ctrl + Shift + End**: Select all text from the cursor to the end of the text.

### Clipboard Operations

1. **Ctrl + C** (or Command + C on macOS): Copy the selected text to the clipboard.
2. **Ctrl + X** (or Command + X on macOS): Cut the selected text to the clipboard.
3. **Ctrl + V** (or Command + V on macOS): Paste the content from the clipboard at the cursor's location.

### Text Editing Keys

1. **Backspace**: Delete the character to the left of the cursor.
2. **Delete**: Delete the character to the right of the cursor.
3. **Ctrl + Backspace**: Delete the word to the left of the cursor.
4. **Ctrl + Delete**: Delete the word to the right of the cursor.
5. **Enter**: Insert a new line at the cursor's position.
6. **Tab**: Insert a tab character (when allowed, though it can sometimes switch focus between form fields).

### Other Useful Keys

1. **Ctrl + Z**: Undo the last action.
2. **Ctrl + Y**: Redo the last undone action.
3. **Ctrl + A**: Select all text within the textarea.

## Not Implemented Functions

### Navigation Keys

done

### Selection Keys

1. **Ctrl + Shift + Arrow Keys**:
   - **Ctrl + Shift + Left Arrow**: Select one word to the left.
   - **Ctrl + Shift + Right Arrow**: Select one word to the right.
   - **Ctrl + Shift + Up Arrow**: Select from the current cursor position to the beginning of the paragraph.
   - **Ctrl + Shift + Down Arrow**: Select from the current cursor position to the end of the paragraph.



### Other Useful Keys

1. **Shift + Tab**: Generally, this navigates focus to the previous form element. Some text editors or configurations may allow it to un-indent text.
