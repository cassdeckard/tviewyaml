# tview Component Coverage in tviewyaml

This table maps [tview](https://pkg.go.dev/github.com/rivo/tview) primitives and features to their implementation status in tviewyaml. Based on the [tview package documentation](https://pkg.go.dev/github.com/rivo/tview).

| tview Component | Implemented | Page-Level | Nested/Child | Primary Example Config | Notes |
|-----------------|-------------|------------|--------------|------------------------|-------|
| **Box** | Yes | No | Yes | [box.yaml](../example/config/box.yaml) | Basic container; base class for most primitives |
| **Button** | Yes | No | Yes (in Form) | [button.yaml](../example/config/button.yaml) | Via form buttons or standalone in flex/grid |
| **Checkbox** | Yes | No | Yes (Form item) | [checkbox.yaml](../example/config/checkbox.yaml) | Form item type `checkbox` |
| **DropDown** | Yes | No | Yes (Form item) | [dropdown.yaml](../example/config/dropdown.yaml) | Form item type `dropdown` |
| **Flex** | Yes | Yes | Yes | [flex.yaml](../example/config/flex.yaml) | Row/column layout; `direction: row` or default column |
| **Form** | Yes | Yes | Yes | [form.yaml](../example/config/form.yaml) | With InputField, Checkbox, Dropdown, Button, TextArea |
| **Frame** | No | No | No | — | Wrapper with header/footer text; not in factory |
| **Grid** | Yes | Yes | Yes | [grid.yaml](../example/config/grid.yaml) | Row/column sizing via `gridRows`, `gridColumns` |
| **Image** | No | No | No | — | Form.AddImage exists in tview; no YAML support |
| **InputField** | Yes | No | Yes (Form item) | [inputfield.yaml](../example/config/inputfield.yaml) | Form item type `inputfield`; supports placeholder, acceptance; `onDone` for Enter/Escape when standalone |
| **List** | Yes | Yes | Yes | [list.yaml](../example/config/list.yaml) | MainText, secondaryText, shortcut, onSelected |
| **Modal** | Yes | Yes | No | [modal.yaml](../example/config/modal.yaml), [modal-about.yaml](../example/config/modal-about.yaml), [modal-confirm.yaml](../example/config/modal-confirm.yaml), [modal-help.yaml](../example/config/modal-help.yaml) | Page-level `type: modal`; text + buttons |
| **Pages** | Yes | Yes (root) | Yes | [app.yaml](../example/config/app.yaml), [nested-pages.yaml](../example/config/nested-pages.yaml) | Tab-like page switching; `ref` to YAML files |
| **Table** | Yes | Yes | Yes | [table.yaml](../example/config/table.yaml) | Headers, rows, borders, fixed rows/columns; `onDone` for Enter/Escape |
| **TextArea** | Yes | No | Yes (Form item) | [form.yaml](../example/config/form.yaml) | Form item type `textarea`; multi-line input |
| **TextView** | Yes | No | Yes | [textview.yaml](../example/config/textview.yaml) | Dynamic colors, regions, scrollable; `onDone` for Enter/Escape; `onHighlighted` for region clicks |
| **TreeView** | Yes | Yes | Yes | [treeview.yaml](../example/config/treeview.yaml), [treeview-standalone.yaml](../example/config/treeview-standalone.yaml), [treeview-modes.yaml](../example/config/treeview-modes.yaml) | Nodes with children; selectable modes |

## Form Item Types

Form items supported in `formItems`:

| Form Item Type | tview Primitive | Example |
|----------------|-----------------|---------|
| `inputfield` | InputField | [form.yaml](../example/config/form.yaml), [inputfield.yaml](../example/config/inputfield.yaml) |
| `checkbox` | Checkbox | [form.yaml](../example/config/form.yaml), [checkbox.yaml](../example/config/checkbox.yaml) |
| `dropdown` | DropDown | [form.yaml](../example/config/form.yaml), [dropdown.yaml](../example/config/dropdown.yaml) |
| `textarea` | TextArea | [form.yaml](../example/config/form.yaml) |
| `button` | Button | All form configs |

## tview Features Not Yet Implemented

| Feature | tview API | Notes |
|---------|-----------|-------|
| **Frame** | `tview.NewFrame(primitive)` | Header/footer text around a primitive |
| **Image** | `tview.NewImage()`, `Form.AddImage()` | Terminal image display; form image field |
| **PasswordField** | `Form.AddPasswordField()` | Implemented via `passwordMode: true` on inputfield + SetMaskCharacter('*'). **TODO:** Support custom `maskCharacter` (e.g. `•`) with `'*'` as default. |

## Callback Support: onDone

TextView, InputField (standalone), and Table support `onDone`—a template expression that runs when the user presses Enter or Escape. State `__doneKey` is set to `"Enter"` or `"Escape"` before the callback runs, enabling template branching. See [ondone-demo.yaml](../example/config/ondone-demo.yaml).

Note: tview's Table does not call SetDoneFunc for Enter when rows are selectable—Enter is used for selection. Escape still triggers onDone on selectable tables.

## Callback Support: onHighlighted (TextView)

TextView with `regions: true` supports `onHighlighted`—a template expression that runs when the highlighted region changes (from clicks or Tab/Enter navigation). State `__highlightedRegion` is set to the first newly highlighted region ID before the callback runs. Use `switchToPage "{{ bindState __highlightedRegion }}"` for presentation-style info bar navigation. See [textview.yaml](../example/config/textview.yaml).

## Additional Example Configs (Feature Demos)

| Config | Purpose |
|--------|---------|
| [main.yaml](../example/config/main.yaml) | Demo menu (list of all features) |
| [state-binding.yaml](../example/config/state-binding.yaml) | Template state binding, reactive views |
| [dynamic-pages.yaml](../example/config/dynamic-pages.yaml) | Dynamic page add/remove |
| [modal-yaml.yaml](../example/config/modal-yaml.yaml) | Modals opened from YAML actions |
| [nested-pages-example.yaml](../example/config/nested-pages-example.yaml) | Nested Pages + Form |
| [clock.yaml](../example/config/clock.yaml) | Live-updating clock (template refresh) |
| [ondone-demo.yaml](../example/config/ondone-demo.yaml) | TextView, InputField, Table with `onDone` (Enter/Escape callbacks) |
| [about.yaml](../example/config/about.yaml) | Simple TextView page |
| [help.yaml](../example/config/help.yaml) | Modal help overlay |
