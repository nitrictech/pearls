# Pearls

A library of utilities and components for building [bubbletea](https://github.com/charmbracelet/bubbletea) programs.

## Composing Models

When including one model in another to create more complex user interface elements it's often necessary to forward `tea.Msg` messages from the parent's `Update` function to the child component's `Update` function. Model updates return a new copy of the model, instead of mutating the existing value, so you typically need to store this updated value back into the parent model.

Here is how you might think to do this:

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

	default:
		var cmd tea.Cmd
		// Child model updated and stored back in the parent.
		m.childModel, cmd = m.childModel.Update(msg)
		return m, cmd
	}
}
```

However, `tea.Model.Update()` returns a `tea.Model` interface, not the original model's type (i.e. the type is erased). So this results in an error like this:

> Cannot assign tea.Model to m.childModel (type child.Model) in multiple assignment

To assist with composition pearls all implement a second `Update` function which preserves the type information. This update function's name will match the model's name.

```go
m.childModel, cmd = m.childModel.UpdateChildModel(msg)
```