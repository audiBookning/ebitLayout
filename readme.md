# Randomly trying out basic UI and layouts in Ebitengine.

None of this is much organized. It is just random code.

## Examples:

### WIDGETS:

- accordion widget

    - `go run .\cmd\accordion\`

- Button widget - with input manager

    - `go run .\cmd\button\`

- slider widgets

    - `go run .\cmd\slider01\` // no keyboard. drag slide and click slide with mouse interaction

    - `go run .\cmd\slider02\` // no keyboard, no drag slide, only clcik slide?

- slider content

    - `go run .\cmd\imageSlider01\` // just keyboard, no snapp

    - `go run .\cmd\imageSlider02\` // with snapp - mouse drag buggy

    - `go run .\cmd\imageSlider03\` // with "physical snapping" - mouse drag buggy

    - `go run .\cmd\imageSlider04\` // simpler and looks better - no animation

    - `go run .\cmd\imageSlider05\` // simpler and looks better - animation

    - `go run .\cmd\imageSlider06\` // like imageSlider05 with navigator and inputmanager - unfinished and Buggy

- toggle widget

    // very basic
    - `go run .\cmd\buttonToggle01` 
    - `go run .\cmd\buttonToggle02` 

    - `go run .\cmd\buttonToggle03` // with a label

- textArea input widget

    - `go run .\cmd\textarea\` // basic draft

    - `go run .\cmd\textareaSelection\` // textArea input widget with many more features like keyboard selection , tabs indent, etc. Work in progress. Still very buggy


### LAYOUT:

- breakpoints layouts: bootstrap like

    - `go run .\cmd\breakpoints01\` // first draft
    
    - `go run .\cmd\breakpoints02\` // first draft + separate ui elements

    - `go run .\cmd\breakpoints03\` // ui elements have their own breakpoints (more like bootstrap)

- container: bootstrap like 

    - `go run .\cmd\container01\`

- flex layout, css inspired

    - `go run .\cmd\flex\` 

- flex 12 column layout, css inspired

    - `go run .\cmd\flexColumn\` // example showing how to use flex for a 12 column layout

- grid, static or css inspired

    - `go run .\cmd\gridCss\`

    - `go run .\cmd\gridStatic\`
    - 
    - `go run .\cmd\gridCss12Column\` // example showing how to use grid for a 12 column. still buggy when trying to maintain the aspect ratio at certain snaller tresholds. and the margins and such are not yet responsive.

- stacked Layout

    - `go run .\cmd\stackedLayout\`

- tabbed Layout

    - `go run .\cmd\tabbedLayout\`


### MIXED:

- Base Layout - uses stack navigator and input manager

    - `go run .\cmd\base01\`

### VARIOUS:

- basic navigations and others...

- navigation: android-like stack navigation - navigator

    - `go run .\cmd\navigation01\` // uses keys for nav
    
    - `go run .\cmd\navigation02\` // uses buttons and mouse and shortcut keys for nav
    
    - `go run .\cmd\navigation03Stack\` // uses a stack for the naviagtion as in mobile nav
  
    - `go run .\cmd\navigation04refactor\` // much refactoring and additon of some more features
    
    - `go run .\cmd\navigation05` // refactor into separate Page and Navigator packages
                                  // refactor into more custom pages

    - `go run .\cmd\navigation06\` // incomplete - refactor into a custom draw function

- Page Type has several ways to be customized: 
  - Override the whole Draw method.
  - Only Override the DrawBackground method with SetCustomDrawBackground
  - Only Override the DrawElements method with SetCustomDrawElements
  - Use the SetCustomDraw method to add custom drawing code at he end of the normal draw method.
  - - in the same way we can add as many custom logic as we want or also custom update code etc

- scroll widget with some textArea

    - `go run .\cmd\scrolling\`

- sidebar menu with top bar

    - `go run .\cmd\sidebar01\` // basic form

    - `go run .\cmd\sidebar02\` // uses input manager

    - `go run .\cmd\sidebar03\` // uses input manager and navigator, no internal pages

    - `go run .\cmd\sidebar04\` // uses input manager and navigator, several example pages for navigation

    - `go run .\cmd\sidebar05\` // 

    - `go run .\cmd\sidebar04\` // Buggy

- graph examples

    - `go run .\cmd\graph01\` // basic sine wave graph

    - `go run .\cmd\graph02\` // basic sine wave graph with numbered ticks

    - `go run .\cmd\graph03\` // basic sine wave graph with numbered ticks and labels

    - `go run .\cmd\graph04\` // basic sine wave graph with numbered ticks and labels


## Utils

- textwrapper util type to ease text rendering


## Notes

- the order of the button registration is as important as the order of drawning

- the input manager gives a better control over the button events

- the sidebarControler allows a better control over the sidebar and its dependant outside click area


## References

- https://devsquad.com/blog/user-interface-layouts
