package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func buildFooter(m dashboardModel, safeTextWidth int) string {
	var footer string
	if (m.state == stateLoaded || m.state == stateIdle) && m.activeTeamID != "" {
		if perf, ok := m.teamPerf[m.activeTeamID]; ok {
			perfText := fmt.Sprintf("Last fetch completed in %s | Tasks Per Second: %s | Est. RPM: %s",
				perf.Duration, perf.TPS, perf.RPM)
			if lipgloss.Width(perfText) > safeTextWidth {
				runes := []rune(perfText)
				perfText = string(runes[:safeTextWidth-1]) + "…"
			}
			footer = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#E0AAFF")).
				MarginTop(1).
				Italic(true).
				Render(perfText)
		}
	}
	return footer
}
