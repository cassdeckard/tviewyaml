package builder

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/cassdeckard/tviewyaml/config"
	"github.com/cassdeckard/tviewyaml/template"
)

// Builder orchestrates the building of tview UI from configuration
type Builder struct {
	factory  *Factory
	mapper   *PropertyMapper
	attacher *CallbackAttacher
	executor *template.Executor
	context  *template.Context
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
	// Create the top-level primitive
	primitive, err := b.factory.CreatePrimitiveFromPageConfig(pageConfig)
	if err != nil {
		return nil, err
	}

	// Apply page-level properties
	if err := b.mapper.ApplyPageProperties(primitive, pageConfig); err != nil {
		return nil, err
	}

	// Build based on type
	switch pageConfig.Type {
	case "list":
		return b.buildList(primitive.(*tview.List), pageConfig)
	case "flex":
		return b.buildFlex(primitive.(*tview.Flex), pageConfig)
	case "form":
		return b.buildForm(primitive.(*tview.Form), pageConfig)
	case "table":
		return b.buildTable(primitive.(*tview.Table), pageConfig)
	case "treeView":
		return b.buildTreeView(primitive.(*tview.TreeView), pageConfig)
	default:
		return primitive, nil
	}
}

// buildList populates a list with items
func (b *Builder) buildList(list *tview.List, cfg *config.PageConfig) (tview.Primitive, error) {
	for _, item := range cfg.ListItems {
		shortcut := rune(0)
		if len(item.Shortcut) > 0 {
			shortcut = rune(item.Shortcut[0])
		}

		// Create callback from template
		var callback func()
		if item.OnSelected != "" {
			cb, err := b.executor.ExecuteCallback(item.OnSelected)
			if err != nil {
				return nil, fmt.Errorf("failed to execute callback for list item: %w", err)
			}
			callback = cb
		}

		list.AddItem(item.MainText, item.SecondaryText, shortcut, callback)
	}

	return list, nil
}

// buildFlex populates a flex container with items
func (b *Builder) buildFlex(flex *tview.Flex, cfg *config.PageConfig) (tview.Primitive, error) {
	for _, item := range cfg.Items {
		if item.Primitive == nil {
			continue
		}

		child, err := b.buildPrimitive(item.Primitive)
		if err != nil {
			return nil, err
		}

		flex.AddItem(child, item.FixedSize, item.Proportion, item.Focus)
	}

	return flex, nil
}

// buildForm populates a form with items
func (b *Builder) buildForm(form *tview.Form, cfg *config.PageConfig) (tview.Primitive, error) {
	_, err := b.addFormItems(form, cfg.FormItems)
	if err != nil {
		return nil, err
	}
	// SetCancelFunc runs when user presses Escape on the form; use OnCancel if set, else OnSubmit
	if cfg.OnCancel != "" || cfg.OnSubmit != "" {
		expr := cfg.OnCancel
		if expr == "" {
			expr = cfg.OnSubmit
		}
		cb, err := b.executor.ExecuteCallback(expr)
		if err != nil {
			return nil, fmt.Errorf("failed to execute form cancel callback: %w", err)
		}
		form.SetCancelFunc(cb)
	}
	if cfg.OnSubmit != "" && cfg.Name != "" {
		cb, err := b.executor.ExecuteCallback(cfg.OnSubmit)
		if err != nil {
			return nil, fmt.Errorf("failed to execute onSubmit callback: %w", err)
		}
		b.context.RegisterFormSubmit(cfg.Name, cb)
	}
	return form, nil
}

// addFormItems adds form items to a form (shared logic for both page-level and nested forms)
func (b *Builder) addFormItems(form *tview.Form, formItems []config.FormItem) (*tview.Form, error) {
	for _, item := range formItems {
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
						return nil, fmt.Errorf("failed to execute callback for inputfield %q: %w", item.Label, err)
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
					return nil, fmt.Errorf("failed to execute callback for button: %w", err)
				}
				callback = cb
			}
			form.AddButton(item.Label, callback)
		case "checkbox":
			var changedFunc func(checked bool)
			if item.OnChanged != "" {
				cb, err := b.executor.ExecuteCallback(item.OnChanged)
				if err != nil {
					return nil, fmt.Errorf("failed to execute callback for checkbox %q: %w", item.Label, err)
				}
				changedFunc = func(checked bool) { cb() }
			}
			form.AddCheckbox(item.Label, item.Checked, changedFunc)
		case "dropdown":
			var selectedFunc func(text string, index int)
			if item.OnChanged != "" {
				cb, err := b.executor.ExecuteCallback(item.OnChanged)
				if err != nil {
					return nil, fmt.Errorf("failed to execute callback for dropdown %q: %w", item.Label, err)
				}
				selectedFunc = func(text string, index int) { cb() }
			}
			form.AddDropDown(item.Label, item.Options, 0, selectedFunc)
		}
	}

	return form, nil
}

// buildTable populates a table with data
func (b *Builder) buildTable(table *tview.Table, cfg *config.PageConfig) (tview.Primitive, error) {
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
func (b *Builder) buildPrimitive(prim *config.Primitive) (tview.Primitive, error) {
	// Create primitive
	primitive, err := b.factory.CreatePrimitive(prim)
	if err != nil {
		return nil, err
	}

	// Apply properties
	if err := b.mapper.ApplyProperties(primitive, prim); err != nil {
		return nil, err
	}

	// Handle callbacks
	if prim.OnSelected != "" {
		callback, err := b.executor.ExecuteCallback(prim.OnSelected)
		if err != nil {
			return nil, fmt.Errorf("failed to execute callback: %w", err)
		}
		b.attacher.AttachCallback(primitive, callback)
	}

	// Handle nested items for specific types
	switch v := primitive.(type) {
	case *tview.Flex:
		if err := b.populateFlexItems(v, prim); err != nil {
			return nil, err
		}
	case *tview.List:
		if err := b.populateListItems(v, prim); err != nil {
			return nil, err
		}
	case *tview.Form:
		if err := b.populateFormItems(v, prim); err != nil {
			return nil, err
		}
	case *tview.Table:
		if err := b.populateTableData(v, prim); err != nil {
			return nil, err
		}
	case *tview.TreeView:
		if err := b.populateTreeView(v, prim); err != nil {
			return nil, err
		}
	case *tview.Grid:
		if err := b.populateGridItems(v, prim); err != nil {
			return nil, err
		}
	}

	return primitive, nil
}

// populateFlexItems adds items to a flex container
func (b *Builder) populateFlexItems(flex *tview.Flex, prim *config.Primitive) error {
	for _, item := range prim.Items {
		if item.Primitive == nil {
			continue
		}

		child, err := b.buildPrimitive(item.Primitive)
		if err != nil {
			return err
		}

		flex.AddItem(child, item.FixedSize, item.Proportion, item.Focus)
	}
	return nil
}

// populateListItems adds items to a list
func (b *Builder) populateListItems(list *tview.List, prim *config.Primitive) error {
	for _, item := range prim.ListItems {
		shortcut := rune(0)
		if len(item.Shortcut) > 0 {
			shortcut = rune(item.Shortcut[0])
		}

		var callback func()
		if item.OnSelected != "" {
			cb, err := b.executor.ExecuteCallback(item.OnSelected)
			if err != nil {
				return fmt.Errorf("failed to execute callback for list item: %w", err)
			}
			callback = cb
		}

		list.AddItem(item.MainText, item.SecondaryText, shortcut, callback)
	}
	return nil
}

// populateFormItems adds items to a form (delegates to shared logic)
func (b *Builder) populateFormItems(form *tview.Form, prim *config.Primitive) error {
	_, err := b.addFormItems(form, prim.FormItems)
	if err != nil {
		return err
	}
	if prim.OnCancel != "" || prim.OnSubmit != "" {
		expr := prim.OnCancel
		if expr == "" {
			expr = prim.OnSubmit
		}
		cb, err := b.executor.ExecuteCallback(expr)
		if err != nil {
			return fmt.Errorf("failed to execute form cancel callback: %w", err)
		}
		form.SetCancelFunc(cb)
	}
	if prim.OnSubmit != "" && prim.Name != "" {
		cb, err := b.executor.ExecuteCallback(prim.OnSubmit)
		if err != nil {
			return fmt.Errorf("failed to execute onSubmit callback: %w", err)
		}
		b.context.RegisterFormSubmit(prim.Name, cb)
	}
	return nil
}

// populateTableData populates table with data from primitive config
func (b *Builder) populateTableData(table *tview.Table, prim *config.Primitive) error {
	colors := []string{"white", "green", "blue", "red"}
	
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
func (b *Builder) populateTreeView(tree *tview.TreeView, prim *config.Primitive) error {
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

	// Now connect children
	for _, node := range prim.Nodes {
		parent := tviewNodeMap[node.Name]
		for _, childName := range node.Children {
			if child, ok := tviewNodeMap[childName]; ok {
				parent.AddChild(child)
			}
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
func (b *Builder) populateGridItems(grid *tview.Grid, prim *config.Primitive) error {
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

		child, err := b.buildPrimitive(item.Primitive)
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
func (b *Builder) buildTreeView(tree *tview.TreeView, cfg *config.PageConfig) (tview.Primitive, error) {
	prim := &config.Primitive{
		OnNodeSelected: cfg.OnNodeSelected,
		RootNode:       cfg.RootNode,
		CurrentNode:    cfg.CurrentNode,
		Nodes:          cfg.Nodes,
	}
	return tree, b.populateTreeView(tree, prim)
}
