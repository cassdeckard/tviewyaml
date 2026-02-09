# tviewyaml Feature Gaps

**Purpose**: This document catalogs features provided by the tviewyaml package that aren't sufficiently demonstrated in the example application.

**Last Updated**: February 9, 2026

---

## How to Use This Document

This serves as a roadmap for improving the example app. Each gap represents an opportunity to:
- Add a new example configuration file
- Enhance existing examples
- Improve user onboarding and feature discovery

Items are prioritized based on:
- Frequency of use in typical applications
- Importance for understanding the package
- Complexity of the feature

---

## Medium Priority Gaps

### 1. InputField onChange Callback

**Description**: Callback fired on every keystroke in an input field

**Why It Matters**: Enables real-time validation, search-as-you-type, character counting, and other interactive patterns.

**Current Status**: `onChanged` callback exists but isn't demonstrated for InputField.

**Suggested Example**: Add `inputfield-live.yaml` or update existing:
- InputField with `onChanged` callback
- TextView showing character count or validation status
- Live feedback via state binding

**Related Files**:
- `config/types.go` - `OnChanged` in FormItem
- `builder/builder.go` - `SetChangedFunc` in addFormItems

---

### 2. Nested Pages Container

**Description**: Pages primitive used within flex/grid layouts (not just as root)

**Why It Matters**: Enables sub-navigation patterns and tab-like interfaces within sections of the UI.

**Current Status**: Pages only demonstrated as root container in `app.yaml`.

**Suggested Example**: Add `nested-pages.yaml`:
- Flex layout with sidebar and content area
- Content area uses nested pages container
- Buttons in sidebar switch content pages

**Related Files**:
- `builder/factory.go` - Pages factory
- Package supports nesting inherently

---

### 3. Extended Special Keys

**Description**: Special keys beyond the basics: Insert, Delete, PgUp, PgDn, Home, End, F2-F11

**Why It Matters**: Many terminal applications use these keys (vim-style navigation, IDE-style shortcuts).

**Current Status**: Only Escape, Ctrl+C, Ctrl+Q demonstrated.

**Suggested Example**: Update `app.yaml` or add comments:
- `F1` for help
- `PgUp`/`PgDn` for page navigation
- `Home`/`End` for first/last page
- Comment listing all supported special keys

**Related Files**:
- `template/keybinding.go` - Full special key support
- `keys/keys_test.go` - All keys tested

---

### 4. TreeView Selection Mode Differences

**Description**: Three distinct selection modes: `auto`, `true`, `false` with different behaviors

**Why It Matters**: Understanding mode differences helps users choose appropriate behavior for their use case.

**Current Status**: Modes mentioned but behavioral differences unclear.

**Suggested Example**: Create `treeview-modes.yaml`:
- Three separate trees demonstrating each mode
- Clear labels explaining behavior
- Comments describing when to use each mode

**Related Files**:
- `builder/builder.go` - `populateTreeView` implements mode logic
- `config/types.go` - `Selectable` field

---

## Lower Priority / Advanced Gaps

### 6. OnStateChange Subscriptions

**Description**: Subscribe to state changes with callback functions (reactive pattern)

**Why It Matters**: Enables reactive programming patterns and decoupled components responding to state changes.

**Current Status**: Method exists but isn't demonstrated. Examples only show `SetStateDirect` and `bindState`.

**Suggested Example**: Add Go code example or advanced YAML:
- Custom template function that uses `OnStateChange`
- Multiple components reacting to same state key
- Document in comments or example code

**Related Files**:
- `template/context.go` - `OnStateChange` method
- Used internally for bound views

---

### 7. YAML-Configured Modals

**Description**: Modal primitive defined in YAML configuration (not just via `showSimpleModal`)

**Why It Matters**: Static modals for confirmations, about boxes, or error messages can be pre-configured.

**Current Status**: Package supports `modal` primitive type but examples only show programmatic creation.

**Suggested Example**: Add `modal-yaml.yaml`:
- Modal defined in YAML with buttons
- Compare to `showSimpleModal` approach
- Show when each approach is appropriate

**Related Files**:
- `builder/factory.go` - Modal factory exists
- Modal creation supported

---

### 8. Nested Form Submission Pattern

**Description**: Form as a nested primitive with its own `onSubmit`, triggered via `runFormSubmit`

**Why It Matters**: Enables complex multi-section forms where sections can be submitted independently.

**Current Status**: `runFormSubmit` shown but not nested form pattern.

**Suggested Example**: Add `form-nested.yaml`:
- Multiple forms in grid/flex layout
- Buttons triggering `{{ runFormSubmit "section1" }}`
- Demonstrate multi-step or multi-section forms

**Related Files**:
- `builder/builder.go` - `RegisterFormSubmit` and `RunFormSubmit`
- Template function registered

---

### 9. Grid Size Constraints

**Description**: `minHeight` and `minWidth` constraints on grid items

**Why It Matters**: Controls responsive behavior and prevents grid items from becoming too small.

**Current Status**: Fields exist but behavior/impact unclear in example.

**Suggested Example**: Update `grid.yaml`:
- Grid items with explicit `minHeight`/`minWidth`
- Comments explaining constraint behavior
- Show impact on layout

**Related Files**:
- `config/types.go` - GridItem fields
- `builder/builder.go` - `AddItem` with constraints

---

### 10. Dynamic Page Lifecycle Management

**Description**: Programmatic page creation, removal, and navigation patterns

**Why It Matters**: Enables wizard flows, drill-down interfaces, and dynamic content presentation.

**Current Status**: Static pages only. Dynamic capabilities exist but not shown.

**Suggested Example**: Create advanced example in Go code:
- Create pages dynamically from data
- Use `removePage` for cleanup
- Show pattern for detail/master workflow

**Related Files**:
- `builder/callbacks.go` - Page management functions
- Context Pages field

---

## Additional Missing Features

### 11. noop Function

**Description**: No-operation template function

**Why It Matters**: Useful as placeholder or for conditionally disabling actions.

**Current Status**: Built-in function not demonstrated.

**Suggested Example**: Add comment in existing examples showing usage for disabled buttons.

**Related Files**: Built-in template function

---

### 12. Complex State Binding Patterns

**Description**: Multiple bound views updating from single state key

**Why It Matters**: Shows reactive UI pattern where multiple components react to state changes.

**Current Status**: Single bound view (clock) shown but not multi-view pattern.

**Suggested Example**: Add example with:
- One state key
- Multiple TextViews bound to same key
- Shows broadcast update pattern

**Related Files**: 
- `template/context.go` - BoundView registration
- Background refresh system

---

## Already Well Demonstrated

The following features are adequately covered in existing examples:

### Primitives
- **Box** - box.yaml (borders, titles, alignment)
- **Button** - button.yaml (callbacks, modals, notifications)
- **Checkbox** - checkbox.yaml (checked state, labels)
- **Dropdown** - dropdown.yaml (options, multiple dropdowns)
- **Flex** - flex.yaml, box.yaml (direction, sizing, nesting)
- **Form** - form.yaml (all item types including textarea, submission, cancellation, onSubmit/onCancel callbacks, runFormSubmit/runFormCancel)
- **Grid** - grid.yaml (rows, columns, spans, borders)
- **InputField** - inputfield.yaml (validation, password mode, placeholder)
- **List** - list.yaml, main.yaml (items, shortcuts, selection)
- **Table** - table.yaml (headers, rows, selection, fixed rows, fixed columns, column colors)
- **TextView** - textview.yaml (colors, regions, formatting)
- **TextArea** - form.yaml (multi-line input, placeholder)
- **TreeView** - treeview.yaml, treeview-standalone.yaml (hierarchy, colors)

### Template Features
- **State binding** - clock.yaml (`{{ bindState clock }}`)
- **Function calls** - button.yaml, modal.yaml (multiple function types)
- **Custom functions** - clock.go (registration and implementation)
- **Multi-parameter calls** - modal.yaml (`showSimpleModal` with multiple args)

### Navigation & Layout
- **Page switching** - Demonstrated throughout via `{{ switchToPage }}`
- **Page removal** - dynamic-pages.yaml, temp-detail.yaml (removePage function, dynamic page lifecycle)
- **Global key bindings** - app.yaml (Escape, Ctrl+C/Q, Alt+1-9, Shift+F1, Meta+H/M - all modifiers)
- **Page container** - app.yaml (root pages structure)
- **Responsive layouts** - grid.yaml, flex.yaml (flexible sizing)

### Callbacks & Events
- **Button onSelected** - button.yaml
- **List item selection** - list.yaml
- **Form submission** - form.yaml
- **Table cell selection** - table.yaml
- **Tree node selection** - treeview.yaml
- **Dropdown onChange** - Shown in dropdown.yaml

### Application Features
- **Mouse support** - app.yaml (`enableMouse: true`)
- **App lifecycle** - main.go (`app.Run()`, defer `app.Stop()`)
- **Error handling** - main.go (page errors, fatal errors)
- **Escape passthrough** - app.yaml configuration

### Advanced Patterns
- **Background goroutines** - clock.go (state updates from ticker)
- **Thread-safe state** - clock.go (QueueUpdateDraw pattern)
- **Fluent API** - main.go (builder chaining with `With`)

---

## Implementation Priority Recommendations

When adding examples for missing features, consider this order:

**Phase 1: Essential Features**
1. TextArea primitive (common use case)
2. Table columnColors (newly added feature)
3. Form onCancel (important UX pattern)

**Phase 2: Power User Features**
4. Advanced key modifiers (Alt, Shift, Meta)
5. removePage function (dynamic UI)
6. Table fixedColumns (wide tables)

**Phase 3: Interactive Patterns**
7. InputField onChanged (real-time feedback)
8. Extended special keys (F1-F12, PgUp, etc.)
9. TreeView selection modes (behavioral clarity)

**Phase 4: Advanced Topics**
10. OnStateChange subscriptions (reactive patterns)
11. Nested pages container (sub-navigation)
12. Complex state binding (multi-view updates)

**Phase 5: Edge Cases & Polish**
13. YAML-configured modals (alternative approach)
14. Nested form patterns (wizard UIs)
15. Grid constraints (responsive control)
16. Dynamic page lifecycle (advanced workflows)

---

## Notes

- Some features may be intentionally undocumented if they're internal or advanced
- Feature coverage should balance completeness with example complexity
- Consider combining related features in single examples where appropriate
- Each new example should serve a clear use case, not just demonstrate syntax

---

## Contributing

When adding examples for these gaps:

1. Create YAML file in `example/config/`
2. Add page reference to `app.yaml`
3. Update main menu in `main.yaml` with new shortcut
4. Add custom functions to `main.go` or separate file if needed
5. Test thoroughly to ensure example works
6. Update this document to move feature from gaps to "Already Well Demonstrated"

---

## Summary Statistics

- **Total Features Identified**: 11 gaps
- **High Priority**: 0 items
- **Medium Priority**: 4 items  
- **Low Priority**: 7 items
- **Well Demonstrated**: 35+ features

**Coverage Estimate**: The example app demonstrates approximately 76% of package functionality, with most common use cases well covered.
