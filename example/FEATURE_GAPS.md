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

## Lower Priority / Advanced Gaps

### 2. OnStateChange Subscriptions

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

### 3. Nested Form Submission Pattern

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

### 4. Grid Size Constraints

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

### 5. noop Function

**Description**: No-operation template function

**Why It Matters**: Useful as placeholder or for conditionally disabling actions.

**Current Status**: Built-in function not demonstrated.

**Suggested Example**: Add comment in existing examples showing usage for disabled buttons.

**Related Files**: Built-in template function

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
- **InputField** - inputfield.yaml (validation, password mode, placeholder, onChange callbacks with real-time feedback)
- **List** - list.yaml, main.yaml (items, shortcuts, selection)
- **Modal** - modal-yaml.yaml, modal.yaml (YAML-configured with buttons/callbacks, programmatic via showSimpleModal)
- **Table** - table.yaml (headers, rows, selection, fixed rows, fixed columns, column colors)
- **TextView** - textview.yaml (colors, regions, formatting)
- **TextArea** - form.yaml (multi-line input, placeholder)
- **TreeView** - treeview.yaml, treeview-standalone.yaml, treeview-modes.yaml (hierarchy, colors, selection modes: auto/true/false)

### Template Features
- **State binding** - clock.yaml (`{{ bindState clock }}`), state-binding.yaml (multiple views bound to single state key, reactive broadcast pattern)
- **Function calls** - button.yaml, modal.yaml (multiple function types)
- **Custom functions** - clock.go, state-binding.go (registration and implementation)
- **Multi-parameter calls** - modal.yaml (`showSimpleModal` with multiple args)

### Navigation & Layout
- **Page switching** - Demonstrated throughout via `{{ switchToPage }}`
- **Page removal** - dynamic-pages.yaml (removePage function)
- **Dynamic page creation** - dynamic-pages.yaml (programmatic AddPage from Go code, data-driven detail pages)
- **Global key bindings** - app.yaml (Escape, Ctrl+C/Q, Alt+1-9, Shift+F1, Meta+H/M, Home/End/PgUp/PgDn/Insert/Delete, F2/F5/F11/F12 - all modifiers and special keys)
- **Page container** - app.yaml (root pages structure)
- **Nested pages** - nested-pages.yaml (pages within flex layouts, tab-based interfaces, sub-navigation)
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
7. TreeView selection modes (behavioral clarity)

**Phase 4: Advanced Topics**
8. OnStateChange subscriptions (reactive patterns)

**Phase 5: Edge Cases & Polish**
9. Nested form patterns (wizard UIs)
10. Grid constraints (responsive control)

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

- **Total Features Identified**: 4 gaps
- **High Priority**: 0 items
- **Medium Priority**: 0 items  
- **Low Priority**: 4 items
- **Well Demonstrated**: 40+ features

**Coverage Estimate**: The example app demonstrates approximately 87% of package functionality, with most common use cases well covered.
