package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func buildLogsDisplay(m dashboardModel, safeTextWidth int) string {
	logInnerWidth := safeTextWidth - 4
	if logInnerWidth < 10 {
		logInnerWidth = 10
	}

	logBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Width(logInnerWidth).
		Height(4).
		Padding(0, 1)

	logTitle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("--- Live API Event Log ---")

	displayLogs := m.logs
	if len(displayLogs) > 3 {
		displayLogs = displayLogs[len(displayLogs)-3:]
	}

	var safeLogs []string
	for _, l := range displayLogs {
		runes := []rune(l)
		if len(runes) > logInnerWidth {
			safeLogs = append(safeLogs, string(runes[:logInnerWidth-3])+"...")
		} else {
			safeLogs = append(safeLogs, l)
		}
	}
	logContent := strings.Join(safeLogs, "\n")
	logsDisplay := logBoxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, logTitle, logContent))
	return logsDisplay
}
