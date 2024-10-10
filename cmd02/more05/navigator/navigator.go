package navigator

import (
	"log"
	"sync"

	"example.com/menu/cmd02/more04/types"
)

// Navigator manages page navigation.
type Navigator struct {
	pages   map[string]types.Page
	current types.Page
	mu      sync.RWMutex
	onExit  func()
}

// NewNavigator initializes a new Navigator.
// onExit is called when the "exit" page is triggered.
func NewNavigator(onExit func()) *Navigator {
	return &Navigator{
		pages:  make(map[string]types.Page),
		onExit: onExit,
	}
}

// AddPage adds a new page to the navigator.
func (n *Navigator) AddPage(name string, page types.Page) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.pages[name] = page
}

// Layout all pages
func (n *Navigator) Layout(outsideWidth, outsideHeight int) (int, int) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	for _, page := range n.pages {
		page.Layout(outsideWidth, outsideHeight)
	}
	return outsideWidth, outsideHeight
}

// SwitchTo switches the current page to the specified page.
func (n *Navigator) SwitchTo(pageName string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if pageName == "exit" {
		log.Println("Exit requested")
		if n.onExit != nil {
			n.onExit()
		}
		return
	}

	if page, exists := n.pages[pageName]; exists {
		log.Printf("Switching to page: %s\n", pageName)
		n.current = page
		// Reset button states if applicable
		if pageWithReset, ok := page.(interface{ ResetButtonStates() }); ok {
			pageWithReset.ResetButtonStates()
		}
	} else {
		log.Printf("Page %s does not exist!\n", pageName)
	}
}

func (n *Navigator) CurrentActivePage() types.Page {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.current
}

func (n *Navigator) ResetAllPagesButtonStates() {
	n.mu.RLock()
	defer n.mu.RUnlock()
	for _, page := range n.pages {
		if pageWithReset, ok := page.(interface{ ResetButtonStates() }); ok {
			pageWithReset.ResetButtonStates()
		}
	}
}
