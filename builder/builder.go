package builder

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
	"github.com/cassdeckard/tviewyaml/config"
	"github.com/cassdeckard/tviewyaml/template"
)

// BuildContext tracks the component path for better error messages
type BuildContext struct {
	path []string
}

// NewBuildContext creates a new build context
func NewBuildContext() *BuildContext {
	return &BuildContext{
		path: make([]string, 0),
	}
}

// Push adds a component to the path
func (bc *BuildContext) Push(component string) {
	bc.path = append(bc.path, component)
}

// Pop removes the last component from the path
func (bc *BuildContext) Pop() {
	if len(bc.path) > 0 {
		bc.path = bc.path[:len(bc.path)-1]
	}
}

// Path returns the current component path as a string
func (bc *BuildContext) Path() string {
	if len(bc.path) == 0 {
		return ""
	}
	return strings.Join(bc.path, " -> ")
}

// Errorf formats an error with the current component path
func (bc *BuildContext) Errorf(format string, args ...interface{}) error {
	if path := bc.Path(); path != "" {
		return fmt.Errorf("%s: "+format, append([]interface{}{path}, args...)...)
	}
	return fmt.Errorf(format, args...)
}

// Builder orchestrates the building of tview UI from configuration
type Builder struct {
	factory  *Factory
	mapper   *PropertyMapper
	attacher *CallbackAttacher
	executor *template.Executor
	context  *template.Context
}

// assertPrimitiveType safely asserts a primitive to a specific type, returning an error if the type doesn't match.
// This prevents panics when factory and builder get out of sync.
func assertPrimitiveType[T tview.Primitive](p tview.Primitive) (T, error) {
	var zero T
	if p == nil {
		return zero, fmt.Errorf("primitive is nil")
	}
	if result, ok := p.(T); ok {
		return result, nil
	}
	return zero, fmt.Errorf("primitive type mismatch: expected %T, got %T", zero, p)
}

// NewBuilder creates a new UI builder
func NewBuilder(ctx *template.Context, registry *template.FunctionRegistry) *Builder {
	executor := template.NewExecutor(ctx, registry)
	return &Builder{
		factory:  NewFactory(),
		mapper:   NewPropertyMapper(ctx, executor),
		attacher: NewCallbackAttacher(),
		executor: executor,
		context:  ctx,
	}
}

// BuildFromConfig builds a tview primitive from a page configuration
func (b *Builder) BuildFromConfig(pageConfig *config.PageConfig) (tview.Primitive, error) {
	bc := NewBuildContext()
	bc.Push(fmt.Sprintf("page:%s", pageConfig.Type))
	defer bc.Pop()

	// Create the top-level primitive
	primitive, err := b.factory.CreatePrimitiveFromPageConfig(pageConfig)
	if err != nil {
		return nil, bc.Errorf("%w", err)
	}

	// Apply page-level properties
	if err := b.mapper.ApplyPageProperties(primitive, pageConfig); err != nil {
		return nil, bc.Errorf("%w", err)
	}

	// Build based on type
	switch pageConfig.Type {
	case "list":
		list, err := assertPrimitiveType[*tview.List](primitive)
		if err != nil {
			return nil, bc.Errorf("failed to build list: %w", err)
		}
		return b.buildList(list, pageConfig, bc)
	case "flex":
		flex, err := assertPrimitiveType[*tview.Flex](primitive)
		if err != nil {
			return nil, bc.Errorf("failed to build flex: %w", err)
		}
		return b.buildFlex(flex, pageConfig, bc)
	case "form":
		form, err := assertPrimitiveType[*tview.Form](primitive)
		if err != nil {
			return nil, bc.Errorf("failed to build form: %w", err)
		}
		return b.buildForm(form, pageConfig, bc)
	case "table":
		table, err := assertPrimitiveType[*tview.Table](primitive)
		if err != nil {
			return nil, bc.Errorf("failed to build table: %w", err)
		}
		return b.buildTable(table, pageConfig, bc)
	case "treeView":
		tree, err := assertPrimitiveType[*tview.TreeView](primitive)
		if err != nil {
			return nil, bc.Errorf("failed to build treeView: %w", err)
		}
		return b.buildTreeView(tree, pageConfig, bc)
	default:
		return primitive, nil
	}
}

// buildList populates a list with items
func (b *Builder) buildList(list *tview.List, cfg *config.PageConfig, bc *BuildContext) (tview.Primitive, error) {
	for i, item := range cfg.ListItems {
		bc.Push(fmt.Sprintf("listItem[%d]", i))
		shortcut := rune(0)
		if len(item.Shortcut) > 0 {
			shortcut = rune(item.Shortcut[0])
		}

		// Create callback from template
		var callback func()
		if item.OnSelected != "" {
			cb, err := b.executor.ExecuteCallback(item.OnSelected)
			if err != nil {
				bc.Pop()
				return nil, bc.Errorf("failed to execute callback: %w", err)
			}
			callback = cb
		}

		list.AddItem(item.MainText, item.SecondaryText, shortcut, callback)
		bc.Pop()
	}

	return list, nil
}

// buildFlex populates a flex container with items
func (b *Builder) buildFlex(flex *tview.Flex, cfg *config.PageConfig, bc *BuildContext) (tview.Primitive, error) {
	for i, item := range cfg.Items {
		if item.Primitive == nil {
			continue
		}

		bc.Push(fmt.Sprintf("flex[%d]", i))
		child, err := b.buildPrimitive(item.Primitive, bc)
		if err != nil {
			bc.Pop()
			return nil, err
		}
		bc.Pop()

		flex.AddItem(child, item.FixedSize, item.Proportion, item.Focus)
	}

	return flex, nil
}

// buildForm populates a form with items
func (b *Builder) buildForm(form *tview.Form, cfg *config.PageConfig, bc *BuildContext) (tview.Primitive, error) {
	_, err := b.addFormItems(form, cfg.FormItems, bc)
	if err != nil {
		return nil, err
	}
	// Setup form callbacks (cancel and submit)
	if err := b.setupFormCallbacks(form, cfg.OnCancel, cfg.OnSubmit, cfg.Name, bc); err != nil {
		return nil, err
	}
	return form, nil
}

// addFormItems adds form items to a form (shared logic for both page-level and nested forms)
func (b *Builder) addFormItems(form *tview.Form, formItems []config.FormItem, bc *BuildContext) (*tview.Form, error) {
	for i, item := range formItems {
		bc.Push(fmt.Sprintf("formItem[%d]:%s", i, item.Type))
		switch item.Type {
		case "inputfield":
			var acceptFunc func(textToCheck string, lastChar rune) bool
			switch item.AcceptanceFunc {
			case "integer":
				acceptFunc = tview.InputFieldInteger
			case "float":
				acceptFunc = tview.InputFieldFloat
			case "maxlength":
				if item.MaxLength > 0 {
					acceptFunc = tview.InputFieldMaxLength(item.MaxLength)
				}
			}

			needCustomInput := item.Placeholder != "" || item.PasswordMode || item.OnChanged != ""
			if needCustomInput {
				input := tview.NewInputField().
					SetLabel(item.Label).
					SetText(item.Value).
					SetFieldWidth(item.FieldWidth)
				if acceptFunc != nil {
					input.SetAcceptanceFunc(acceptFunc)
				}
				if item.Placeholder != "" {
					input.SetPlaceholder(item.Placeholder)
				}
				if item.PasswordMode {
					input.SetMaskCharacter('*')
				}
				if item.OnChanged != "" {
					cb, err := b.executor.ExecuteCallback(item.OnChanged)
					if err != nil {
						bc.Pop()
						return nil, bc.Errorf("failed to execute callback for inputfield %q: %w", item.Label, err)
					}
					input.SetChangedFunc(func(text string) { cb() })
				}
				form.AddFormItem(input)
			} else {
				form.AddInputField(item.Label, item.Value, item.FieldWidth, acceptFunc, nil)
			}

		case "button":
			callback := func() {}
			if item.OnSelected != "" {
				cb, err := b.executor.ExecuteCallback(item.OnSelected)
				if err != nil {
					bc.Pop()
					return nil, bc.Errorf("failed to execute callback for button: %w", err)
				}
				callback = cb
			}
			form.AddButton(item.Label, callback)
		case "checkbox":
			var changedFunc func(checked bool)
			if item.OnChanged != "" {
				cb, err := b.executor.ExecuteCallback(item.OnChanged)
				if err != nil {
					bc.Pop()
					return nil, bc.Errorf("failed to execute callback for checkbox %q: %w", item.Label, err)
				}
				changedFunc = func(checked bool) { cb() }
			}
			form.AddCheckbox(item.Label, item.Checked, changedFunc)
		case "dropdown":
			var selectedFunc func(text string, index int)
			if item.OnChanged != "" {
				cb, err := b.executor.ExecuteCallback(item.OnChanged)
				if err != nil {
					bc.Pop()
					return nil, bc.Errorf("failed to execute callback for dropdown %q: %w", item.Label, err)
				}
				selectedFunc = func(text string, index int) { cb() }
			}
			form.AddDropDown(item.Label, item.Options, 0, selectedFunc)
		}
		bc.Pop()
	}

	return form, nil
}

// setupFormCallbacks configures the cancel and submit callbacks for a form
// This is shared logic used by both buildForm and populateFormItems
func (b *Builder) setupFormCallbacks(form *tview.Form, onCancel, onSubmit, name string, bc *BuildContext) error {
	// Register cancel callback if provided
	if onCancel != "" && name != "" {
		cb, err := b.executor.ExecuteCallback(onCancel)
		if err != nil {
			return bc.Errorf("failed to execute onCancel callback: %w", err)
		}
		b.context.RegisterFormCancel(name, cb)
	}
	// SetCancelFunc runs when user presses Escape on the form; use OnCancel if set, else OnSubmit
	if onCancel != "" || onSubmit != "" {
		expr := onCancel
		if expr == "" {
			expr = onSubmit
		}
		cb, err := b.executor.ExecuteCallback(expr)
		if err != nil {
			return bc.Errorf("failed to execute form cancel callback: %w", err)
		}
		form.SetCancelFunc(cb)
	}
	if onSubmit != "" && name != "" {
		cb, err := b.executor.ExecuteCallback(onSubmit)
		if err != nil {
			return bc.Errorf("failed to execute onSubmit callback: %w", err)
		}
		b.context.RegisterFormSubmit(name, cb)
	}
	return nil
}

// buildTable populates a table with data
func (b *Builder) buildTable(table *tview.Table, cfg *config.PageConfig, bc *BuildContext) (tview.Primitive, error) {
	if cfg.TableData == nil {
		return table, nil
	}

	// Add headers
	for col, header := range cfg.TableData.Headers {
		cell := tview.NewTableCell(header).
			SetTextColor(b.context.Colors.Parse("yellow")).
			SetAlign(tview.AlignCenter).
			SetSelectable(false)
		table.SetCell(0, col, cell)
	}

	// Add rows
	for row, rowData := range cfg.TableData.Rows {
		for col, cellData := range rowData {
			cell := tview.NewTableCell(cellData).
				SetAlign(tview.AlignLeft)
			table.SetCell(row+1, col, cell)
		}
	}

	table.SetBorder(true)
	table.SetSelectable(true, false)
	return table, nil
}

// buildPrimitive builds a primitive from a Primitive config (recursive)
func (b *Builder) buildPrimitive(prim *config.Primitive, bc *BuildContext) (tview.Primitive, error) {
	primName := prim.Type
	if prim.Name != "" {
		primName = fmt.Sprintf("%s:%s", prim.Type, prim.Name)
	}
	bc.Push(primName)
	defer bc.Pop()

	// Create primitive
	primitive, err := b.factory.CreatePrimitive(prim)
	if err != nil {
		return nil, bc.Errorf("%w", err)
	}

	// Apply properties
	if err := b.mapper.ApplyProperties(primitive, prim); err != nil {
		return nil, bc.Errorf("%w", err)
	}

	// Handle callbacks
	if prim.OnSelected != "" {
		callback, err := b.executor.ExecuteCallback(prim.OnSelected)
		if err != nil {
			return nil, bc.Errorf("failed to execute callback: %w", err)
		}
		b.attacher.AttachCallback(primitive, callback)
	}

	// Handle nested items for specific types
	switch v := primitive.(type) {
	case *tview.Flex:
		if err := b.populateFlexItems(v, prim, bc); err != nil {
			return nil, err
		}
	case *tview.List:
		if err := b.populateListItems(v, prim, bc); err != nil {
			return nil, err
		}
	case *tview.Form:
		if err := b.populateFormItems(v, prim, bc); err != nil {
			return nil, err
		}
	case *tview.Table:
		if err := b.populateTableData(v, prim, bc); err != nil {
			return nil, err
		}
	case *tview.TreeView:
		if err := b.populateTreeView(v, prim, bc); err != nil {
			return nil, err
		}
	case *tview.Grid:
		if err := b.populateGridItems(v, prim, bc); err != nil {
			return nil, err
		}
	}

	return primitive, nil
}

// populateFlexItems adds items to a flex container
func (b *Builder) populateFlexItems(flex *tview.Flex, prim *config.Primitive, bc *BuildContext) error {
	for i, item := range prim.Items {
		if item.Primitive == nil {
			continue
		}

		bc.Push(fmt.Sprintf("flex[%d]", i))
		child, err := b.buildPrimitive(item.Primitive, bc)
		bc.Pop()
		if err != nil {
			return err
		}

		flex.AddItem(child, item.FixedSize, item.Proportion, item.Focus)
	}
	return nil
}

// populateListItems adds items to a list
func (b *Builder) populateListItems(list *tview.List, prim *config.Primitive, bc *BuildContext) error {
	for i, item := range prim.ListItems {
		bc.Push(fmt.Sprintf("listItem[%d]", i))
		shortcut := rune(0)
		if len(item.Shortcut) > 0 {
			shortcut = rune(item.Shortcut[0])
		}

		var callback func()
		if item.OnSelected != "" {
			cb, err := b.executor.ExecuteCallback(item.OnSelected)
			if err != nil {
				bc.Pop()
				return bc.Errorf("failed to execute callback: %w", err)
			}
			callback = cb
		}

		list.AddItem(item.MainText, item.SecondaryText, shortcut, callback)
		bc.Pop()
	}
	return nil
}

// populateFormItems adds items to a form (delegates to shared logic)
func (b *Builder) populateFormItems(form *tview.Form, prim *config.Primitive, bc *BuildContext) error {
	_, err := b.addFormItems(form, prim.FormItems, bc)
	if err != nil {
		return err
	}
	// Setup form callbacks (cancel and submit)
	return b.setupFormCallbacks(form, prim.OnCancel, prim.OnSubmit, prim.Name, bc)
}

// populateTableData populates table with data from primitive config
func (b *Builder) populateTableData(table *tview.Table, prim *config.Primitive, bc *BuildContext) error {
	// Use configured column colors if provided, otherwise use defaults
	colors := prim.ColumnColors
	if len(colors) == 0 {
		colors = []string{"white", "green", "blue", "red"}
	}
	
	// Set borders before adding cells (if specified)
	if prim.Borders {
		table.SetBorders(true)
	}
	
	if len(prim.Columns) > 0 {
		// Add headers
		for col, header := range prim.Columns {
			cell := tview.NewTableCell(header).
				SetTextColor(b.context.Colors.Parse("yellow")).
				SetAlign(tview.AlignCenter).
				SetSelectable(false)
			table.SetCell(0, col, cell)
		}
	}

	if len(prim.Rows) > 0 {
		// Add rows
		startRow := 0
		if len(prim.Columns) > 0 {
			startRow = 1
		}

		for row, rowData := range prim.Rows {
			for col, cellData := range rowData {
				// Cycle through colors for each column
				color := colors[col%len(colors)]
				cell := tview.NewTableCell(cellData).
					SetTextColor(b.context.Colors.Parse(color)).
					SetAlign(tview.AlignCenter)
				table.SetCell(startRow+row, col, cell)
			}
		}
	}

	// Set fixed rows/columns after populating
	if prim.FixedRows > 0 || prim.FixedColumns > 0 {
		table.SetFixed(prim.FixedRows, prim.FixedColumns)
	}

	if prim.OnCellSelected != "" {
		table.SetSelectedFunc(func(row int, column int) {
			cellText := ""
			if cell := table.GetCell(row, column); cell != nil {
				cellText = cell.Text
			}
			b.context.SetStateDirect("__selectedCellText", cellText)
			b.context.SetStateDirect("__selectedRow", row)
			b.context.SetStateDirect("__selectedCol", column)
			if cb, err := b.executor.ExecuteCallback(prim.OnCellSelected); err == nil {
				cb()
			}
		})
	}

	table.SetBorder(true)
	table.SetSelectable(true, false)
	return nil
}

// populateTreeView populates a tree view from primitive config
func (b *Builder) populateTreeView(tree *tview.TreeView, prim *config.Primitive, bc *BuildContext) error {
	if len(prim.Nodes) == 0 {
		// No nodes defined, return empty tree
		return nil
	}

	// Build a map of node name to tview.TreeNode and selectable mode
	tviewNodeMap := make(map[string]*tview.TreeNode)
	selectableModeMap := make(map[string]string) // node name -> selectable mode ("true", "auto", "false")

	// Create all nodes first
	for _, node := range prim.Nodes {
		tviewNode := tview.NewTreeNode(node.Text)
		if node.Color != "" {
			tviewNode.SetColor(b.context.Colors.Parse(node.Color))
		}
		// Parse selectable mode: "true", "auto", "false", or default to "auto"
		selectableMode := node.Selectable
		if selectableMode == "" {
			selectableMode = "auto"
		}
		selectableModeMap[node.Name] = selectableMode
		// Set tview selectable: "true" and "auto" are selectable, "false" is not
		tviewNode.SetSelectable(selectableMode != "false")
		// Store node name in Reference so we can look it up later
		tviewNode.SetReference(node.Name)
		tviewNodeMap[node.Name] = tviewNode
	}

	// Now connect children with validation
	for _, node := range prim.Nodes {
		parent := tviewNodeMap[node.Name]
		for _, childName := range node.Children {
			child, ok := tviewNodeMap[childName]
			if !ok {
				// Build list of available node names for error message
				availableNodes := make([]string, 0, len(tviewNodeMap))
				for name := range tviewNodeMap {
					availableNodes = append(availableNodes, fmt.Sprintf("%q", name))
				}
				return bc.Errorf("node %q references unknown child %q (available nodes: %s)",
					node.Name, childName, strings.Join(availableNodes, ", "))
			}
			parent.AddChild(child)
		}
	}

	// Set root and current node
	var root *tview.TreeNode
	if prim.RootNode != "" {
		root = tviewNodeMap[prim.RootNode]
	}
	if root == nil && len(prim.Nodes) > 0 {
		// Default to first node if root not specified
		root = tviewNodeMap[prim.Nodes[0].Name]
	}

	if root != nil {
		tree.SetRoot(root)

		// Set current node
		if prim.CurrentNode != "" {
			if current, ok := tviewNodeMap[prim.CurrentNode]; ok {
				tree.SetCurrentNode(current)
			} else {
				tree.SetCurrentNode(root)
			}
		} else {
			tree.SetCurrentNode(root)
		}
	}

	// Handle node selection
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		children := node.GetChildren()
		isParent := len(children) > 0

		// Get selectable mode from stored node name
		nodeName, ok := node.GetReference().(string)
		selectableMode := "auto" // default
		if ok {
			if mode, exists := selectableModeMap[nodeName]; exists {
				selectableMode = mode
			}
		}

		if selectableMode == "auto" {
			// Default behavior: modal for leaf, toggle expansion for parent (ignore onNodeSelected)
			if isParent {
				// Toggle expansion
				node.SetExpanded(!node.IsExpanded())
			} else {
				// Leaf node - show info
				modal := tview.NewModal().
					SetText(fmt.Sprintf("Selected: %s", node.GetText())).
					AddButtons([]string{"OK"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						b.context.Pages.RemovePage("tree-modal")
					})
				b.context.Pages.AddPage("tree-modal", modal, false, true)
			}
		} else if selectableMode == "true" {
			// Always run onNodeSelected if set, and toggle expansion for parent nodes
			if prim.OnNodeSelected != "" {
				b.context.SetStateDirect("__selectedNodeText", node.GetText())
				if cb, err := b.executor.ExecuteCallback(prim.OnNodeSelected); err == nil {
					cb()
				}
			}
			// Still toggle expansion for parent nodes (preserve default UX)
			if isParent {
				node.SetExpanded(!node.IsExpanded())
			}
		}
		// selectableMode == "false" shouldn't happen (node wouldn't be selectable), but handle gracefully
	})

	return nil
}

// populateGridItems configures a grid and adds items to it
func (b *Builder) populateGridItems(grid *tview.Grid, prim *config.Primitive, bc *BuildContext) error {
	// Set rows (0 = flexible)
	if len(prim.GridRows) > 0 {
		grid.SetRows(prim.GridRows...)
	}

	// Set columns (0 = flexible)
	if len(prim.GridColumns) > 0 {
		grid.SetColumns(prim.GridColumns...)
	}

	// Set borders
	if prim.GridBorders {
		grid.SetBorders(true)
	}

	// Add items
	for _, item := range prim.GridItems {
		if item.Primitive == nil {
			continue
		}

		bc.Push(fmt.Sprintf("grid[%d,%d]", item.Row, item.Column))
		child, err := b.buildPrimitive(item.Primitive, bc)
		bc.Pop()
		if err != nil {
			return err
		}

		// Default spans to 1 if not specified
		rowSpan := item.RowSpan
		if rowSpan == 0 {
			rowSpan = 1
		}
		colSpan := item.ColSpan
		if colSpan == 0 {
			colSpan = 1
		}

		grid.AddItem(child, item.Row, item.Column, rowSpan, colSpan, item.MinHeight, item.MinWidth, item.Focus)
	}

	return nil
}

// buildTreeView populates a tree view from page config (for page-level type: treeView)
func (b *Builder) buildTreeView(tree *tview.TreeView, cfg *config.PageConfig, bc *BuildContext) (tview.Primitive, error) {
	prim := &config.Primitive{
		OnNodeSelected: cfg.OnNodeSelected,
		RootNode:       cfg.RootNode,
		CurrentNode:    cfg.CurrentNode,
		Nodes:          cfg.Nodes,
	}
	return tree, b.populateTreeView(tree, prim, bc)
}
