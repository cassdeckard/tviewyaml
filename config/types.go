package config

// AppConfig represents the top-level application configuration
type AppConfig struct {
	Application ApplicationElement `yaml:"application"`
}

// ApplicationElement contains application-level settings
type ApplicationElement struct {
	Name                   string       `yaml:"name,omitempty"`
	EnableMouse            *bool        `yaml:"enableMouse,omitempty"` // nil = default true
	GlobalKeyBindings      []KeyBinding `yaml:"globalKeyBindings,omitempty"`
	EscapePassthroughPages []string     `yaml:"escapePassthroughPages,omitempty"` // pages where Escape is not captured globally (e.g. so form SetCancelFunc runs)
	Root                   RootElement `yaml:"root"`
}

// KeyBinding represents a global keyboard shortcut
type KeyBinding struct {
	Key    string `yaml:"key"`    // "Escape", "Ctrl+Q", "F1", etc.
	Action string `yaml:"action"` // Template expression
}

// RootElement contains the list of pages (or can be any view type in the future)
type RootElement struct {
	Type  string    `yaml:"type"` // "pages"
	Pages []PageRef `yaml:"pages"`
}

// PageRef references a page configuration file
type PageRef struct {
	Name string `yaml:"name"`
	Ref  string `yaml:"ref"` // Path to YAML file
}

// PageConfig represents a single page/screen configuration
type PageConfig struct {
	Type       string                 `yaml:"type"` // "list", "flex", "form", etc.
	Name       string                 `yaml:"name,omitempty"` // optional name (e.g. for form runFormSubmit)
	Direction  string                 `yaml:"direction,omitempty"`
	Border     bool                   `yaml:"border,omitempty"`
	Title      string                 `yaml:"title,omitempty"`
	TitleAlign string                 `yaml:"titleAlign,omitempty"`
	Items      []FlexItem             `yaml:"items,omitempty"`
	ListItems  []ListItem             `yaml:"listItems,omitempty"`
	FormItems  []FormItem             `yaml:"formItems,omitempty"`
	OnSubmit   string                 `yaml:"onSubmit,omitempty"` // Template expression for runFormSubmit (e.g. Submit button)
	OnCancel   string                 `yaml:"onCancel,omitempty"` // Template expression when form is cancelled (Escape); if unset and OnSubmit set, Escape runs OnSubmit
	TableData  *TableData             `yaml:"tableData,omitempty"`
	// TreeView-specific (for page-level type: treeView)
	OnNodeSelected string     `yaml:"onNodeSelected,omitempty"` // Template expression when a node is selected (state: __selectedNodeText)
	RootNode       string     `yaml:"rootNode,omitempty"`
	CurrentNode    string     `yaml:"currentNode,omitempty"`
	Nodes          []TreeNode `yaml:"nodes,omitempty"`
	// Modal-specific (for page-level type: modal)
	Text    string        `yaml:"text,omitempty"`    // Modal text content
	Buttons []ModalButton `yaml:"buttons,omitempty"` // Modal buttons
	Properties     map[string]interface{} `yaml:",inline"` // Catch-all for other properties
}

// FlexItem represents an item in a flex container
type FlexItem struct {
	Primitive  *Primitive `yaml:"primitive"`
	Spacer     bool       `yaml:"spacer,omitempty"` // if true, treat as spacer (nil primitive)
	FixedSize  int        `yaml:"fixedSize,omitempty"`
	Proportion int        `yaml:"proportion,omitempty"`
	Focus      bool       `yaml:"focus,omitempty"`
}

// Primitive represents a tview primitive configuration
type Primitive struct {
	Name       string `yaml:"name,omitempty"`
	Type       string `yaml:"type"`
	Border     bool   `yaml:"border,omitempty"`
	Title      string `yaml:"title,omitempty"`
	TitleAlign string `yaml:"titleAlign,omitempty"`
	Text       string `yaml:"text,omitempty"`
	TextAlign  string `yaml:"textAlign,omitempty"`
	TextColor  string `yaml:"textColor,omitempty"`
	// TextView-specific properties
	DynamicColors bool       `yaml:"dynamicColors,omitempty"` // Enable color tags in text
	Regions       bool       `yaml:"regions,omitempty"`       // Enable region tags in text
	Label         string     `yaml:"label,omitempty"`
	Checked       bool       `yaml:"checked,omitempty"`
	OnSelected    string     `yaml:"onSelected,omitempty"` // Template expression
	OnChanged     string     `yaml:"onChanged,omitempty"`  // Template expression
	Items         []FlexItem `yaml:"items,omitempty"`
	ListItems     []ListItem `yaml:"listItems,omitempty"`
	Direction     string     `yaml:"direction,omitempty"`
	Columns       []string   `yaml:"columns,omitempty"`
	Rows          [][]string `yaml:"rows,omitempty"`
	Options       []string   `yaml:"options,omitempty"`
	FormItems     []FormItem `yaml:"formItems,omitempty"`
	OnSubmit      string     `yaml:"onSubmit,omitempty"` // Template expression for runFormSubmit (nested form)
	OnCancel      string     `yaml:"onCancel,omitempty"` // Template expression when form is cancelled (Escape); if unset and OnSubmit set, Escape runs OnSubmit
	// Table-specific properties
	OnCellSelected string   `yaml:"onCellSelected,omitempty"` // Template expression when a cell is selected (state: __selectedCellText, __selectedRow, __selectedCol)
	Borders        bool     `yaml:"borders,omitempty"`        // Show borders between cells
	FixedRows      int      `yaml:"fixedRows,omitempty"`      // Number of fixed rows
	FixedColumns   int      `yaml:"fixedColumns,omitempty"`   // Number of fixed columns
	ColumnColors   []string `yaml:"columnColors,omitempty"`   // Colors for each column (cycles if fewer colors than columns)
	// TreeView-specific properties
	OnNodeSelected string     `yaml:"onNodeSelected,omitempty"` // Template expression when a node is selected (state: __selectedNodeText)
	RootNode       string     `yaml:"rootNode,omitempty"`       // Name of the root node
	CurrentNode    string     `yaml:"currentNode,omitempty"`    // Name of the initial current node
	Nodes          []TreeNode `yaml:"nodes,omitempty"`          // List of tree nodes
	// Grid-specific properties
	GridRows    []int        `yaml:"gridRows,omitempty"`    // Row heights (0 = flexible)
	GridColumns []int        `yaml:"gridColumns,omitempty"` // Column widths (0 = flexible)
	GridBorders bool         `yaml:"gridBorders,omitempty"` // Show borders between grid cells
	GridItems   []GridItem   `yaml:"gridItems,omitempty"`   // Items to place in grid
	// Pages-specific properties (for nested pages containers)
	Pages []PageRef `yaml:"pages,omitempty"` // List of pages for nested pages container
	// Modal-specific properties
	Buttons    []ModalButton          `yaml:"buttons,omitempty"` // Buttons with callbacks for modal dialogs
	Properties map[string]interface{} `yaml:",inline"`           // Catch-all for other properties
}

// TreeNode represents a node in a tree view
type TreeNode struct {
	Name       string   `yaml:"name"`                 // Unique identifier for the node
	Text       string   `yaml:"text"`                 // Display text
	Color      string   `yaml:"color,omitempty"`      // Text color
	Selectable string   `yaml:"selectable,omitempty"` // "true" (always run onNodeSelected), "auto" (default behavior), "false" (not selectable). Defaults to "auto" if unset.
	Children   []string `yaml:"children,omitempty"`   // Names of child nodes
}

// ModalButton represents a button in a modal dialog
type ModalButton struct {
	Label      string `yaml:"label"`                // Button text
	OnSelected string `yaml:"onSelected,omitempty"` // Template expression when clicked
}

// GridItem represents an item in a grid layout
type GridItem struct {
	Primitive *Primitive `yaml:"primitive"`           // The primitive to place in the grid
	Row       int        `yaml:"row"`                 // Starting row (0-based)
	Column    int        `yaml:"column"`              // Starting column (0-based)
	RowSpan   int        `yaml:"rowSpan,omitempty"`   // Number of rows to span (default 1)
	ColSpan   int        `yaml:"colSpan,omitempty"`   // Number of columns to span (default 1)
	MinHeight int        `yaml:"minHeight,omitempty"` // Minimum height
	MinWidth  int        `yaml:"minWidth,omitempty"`  // Minimum width
	Focus     bool       `yaml:"focus,omitempty"`     // Whether this item should receive focus
}

// ListItem represents an item in a list
type ListItem struct {
	MainText      string `yaml:"mainText"`
	SecondaryText string `yaml:"secondaryText,omitempty"`
	Shortcut      string `yaml:"shortcut,omitempty"`
	OnSelected    string `yaml:"onSelected,omitempty"` // Template expression
}

// FormItem represents an item in a form
type FormItem struct {
	Type           string   `yaml:"type"` // "inputfield", "button", "checkbox", "dropdown", "textarea"
	Label          string   `yaml:"label"`
	Value          string   `yaml:"value,omitempty"`
	Options        []string `yaml:"options,omitempty"`
	Checked        bool     `yaml:"checked,omitempty"`
	OnSelected     string   `yaml:"onSelected,omitempty"` // Template expression
	OnChanged      string   `yaml:"onChanged,omitempty"`  // Template expression
	FieldWidth     int      `yaml:"fieldWidth,omitempty"`
	PasswordMode   bool     `yaml:"passwordMode,omitempty"`
	AcceptanceFunc string   `yaml:"acceptanceFunc,omitempty"` // "integer", "float", etc.
	MaxLength      int      `yaml:"maxLength,omitempty"`
	Placeholder    string   `yaml:"placeholder,omitempty"`
}

// TableData represents data for a table
type TableData struct {
	Headers []string   `yaml:"headers"`
	Rows    [][]string `yaml:"rows"`
}
