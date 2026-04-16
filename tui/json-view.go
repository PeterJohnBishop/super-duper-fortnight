package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func buildJSONView(m dashboardModel, safeTextWidth int, paneHeight int) string {
	var jsonContent string
	modalWidth := safeTextWidth - 4

	headerText := "[ SHIFT+J: Close | SHIFT+S: Copy ]"
	copiedStatus := ""
	if m.jsonCopied {
		copiedStatus = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true).
			Render(" [ Copied to Clipboard! ]")
	}

	modalHeader := lipgloss.JoinHorizontal(lipgloss.Top, headerText, copiedStatus)

	rawJSON := m.getCurrentJSON()
	jsonLines := strings.Split(rawJSON, "\n")

	usableHeight := paneHeight - 3
	if usableHeight < 1 {
		usableHeight = 1
	}

	maxScroll := len(jsonLines) - usableHeight
	if maxScroll < 0 {
		maxScroll = 0
	}
	if m.jsonScrollOffset > maxScroll {
		m.jsonScrollOffset = maxScroll
	}

	endIdx := m.jsonScrollOffset + usableHeight
	if endIdx > len(jsonLines) {
		endIdx = len(jsonLines)
	}

	var visibleJSON string
	if len(jsonLines) > 0 && m.jsonScrollOffset <= endIdx {
		visibleJSON = strings.Join(jsonLines[m.jsonScrollOffset:endIdx], "\n")
	}

	modalStyle := lipgloss.NewStyle().
		Width(modalWidth).
		Height(paneHeight-2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#9D4EDD")).
		Padding(0, 1)

	jsonContent = lipgloss.JoinVertical(lipgloss.Center,
		lipgloss.NewStyle().Width(modalWidth).Align(lipgloss.Center).Render(modalHeader),
		modalStyle.Render(visibleJSON),
	)
	return jsonContent
}
