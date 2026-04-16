package tui

import "github.com/charmbracelet/lipgloss"

func buildKeybindings(safeTextWidth int) string {
	helpStr := "Nav: h j k l | Back: esc | Select: enter | JSON: Shift+J | Sync: r | Cycle Auto-Sync: SHIFT+F | Open: o | Quit: q | Reset Database: CTRL+X"

	if lipgloss.Width(helpStr) > safeTextWidth {
		// Medium width fallback
		helpStr = "Nav: hjkl | esc: back | enter: sel | J: json | r: sync | F: auto-sync | o: open | q: quit | ctrl+x: reset data"

		if lipgloss.Width(helpStr) > safeTextWidth {
			// Small width fallback
			helpStr = "hjkl:nav | esc:back | enter:sel | J:json | r:sync | F:auto-sync | o:web | q:quit | ctrl+x:reset"

			if lipgloss.Width(helpStr) > safeTextWidth {
				// Extreme squish fallback
				runes := []rune(helpStr)
				helpStr = string(runes[:safeTextWidth-1]) + "…"
			}
		}
	}
	helpText := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).MarginBottom(1).Render(helpStr)
	return helpText
}
