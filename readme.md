# Randomly trying out basic UI and layouts in Ebitengine


### WIDGETS:

- accordion widget

    - `go run .\cmd\accordion\`

- Button widget

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

- grid, static or css inspired

    - `go run .\cmd\gridCss\`

    - `go run .\cmd\gridStatic\`

- stacked Layout

    - `go run .\cmd\stackedLayout\`

- tabbed Layout

    - `go run .\cmd\tabbedLayout\`


### VARIOUS:

- navigation: android-like stack navigation

    - `go run .\cmd\navigation01\` // uses keys for nav
    
    - `go run .\cmd\navigation02\` // uses buttons and mouse for nav

- scroll widget with some textArea

    - `go run .\cmd\scrolling\`

- sidebar menu with top bar

    - `go run .\cmd\sidebar\`






## Notes

- the order of the button registration is as important as the order of drawning

- the input manager gives a better control over the button events

- the sidebarControler allows a better control over the sidebar and its dependant outside click area


## References

- https://devsquad.com/blog/user-interface-layouts
