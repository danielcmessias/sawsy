package utils

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
    Up            key.Binding
    Down          key.Binding
    FirstLine     key.Binding
    LastLine      key.Binding
    TogglePreview key.Binding
    Refresh       key.Binding
    NextTab       key.Binding
    PrevTab       key.Binding
    NextCol       key.Binding
    PrevCol       key.Binding
    StartSearch   key.Binding
    EndSearch     key.Binding
    InspectField  key.Binding
    Home          key.Binding
    Help          key.Binding
    Quit          key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
    return []key.Binding{k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
    return [][]key.Binding{
        {k.Up, k.Down},
        {k.PrevCol, k.NextCol},
        {k.PrevTab, k.NextTab},
        {k.InspectField, k.Home},
        {k.StartSearch, k.Refresh},
        {k.Help, k.Quit},
    }
}

var Keys = KeyMap{
    Up: key.NewBinding(
        key.WithKeys("up", "k"),
        key.WithHelp("↑/k", "move up"),
    ),
    Down: key.NewBinding(
        key.WithKeys("down", "j"),
        key.WithHelp("↓/j", "move down"),
    ),
    NextTab: key.NewBinding(
        key.WithKeys("tab"),
        key.WithHelp("tab", "next tab"),
    ),
    PrevTab: key.NewBinding(
        key.WithKeys("shift+tab"),
        key.WithHelp("shift+tab", "previous tab"),
    ),
    NextCol: key.NewBinding(
        key.WithKeys("right", "l"),
        key.WithHelp("/l", "next col"),
    ),
    PrevCol: key.NewBinding(
        key.WithKeys("left", "h"),
        key.WithHelp("/h", "previous col"),
    ),
    StartSearch: key.NewBinding(
        key.WithKeys("/"),
        key.WithHelp("/", "search"),
    ),
    EndSearch: key.NewBinding(
        key.WithKeys("esc", "enter"),
    ),
    InspectField: key.NewBinding(
        key.WithKeys("enter"),
        key.WithHelp("enter", "inspect"),
    ),
    Home: key.NewBinding(
        key.WithKeys("backspace"),
        key.WithHelp("backspace", "homepage"),
    ),
    FirstLine: key.NewBinding(
        key.WithKeys("g", "home"),
        key.WithHelp("g/home", "first item"),
    ),
    LastLine: key.NewBinding(
        key.WithKeys("G", "end"),
        key.WithHelp("G/end", "last item"),
    ),
    TogglePreview: key.NewBinding(
        key.WithKeys("p"),
        key.WithHelp("p", "open in Preview"),
    ),
    Refresh: key.NewBinding(
        key.WithKeys("r"),
        key.WithHelp("r", "refresh"),
    ),
    Help: key.NewBinding(
        key.WithKeys("?"),
        key.WithHelp("?", "toggle help"),
    ),
    Quit: key.NewBinding(
        key.WithKeys("q", "esc", "ctrl+c"),
        key.WithHelp("q", "quit"),
    ),
}
