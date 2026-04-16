package tui

import "github.com/charmbracelet/lipgloss"

func buildResyncModal(m dashboardModel, safeTextWidth int, paneHeight int) string {
	modalWidth := 60
	if modalWidth > safeTextWidth {
		modalWidth = safeTextWidth
	}

	content := lipgloss.NewStyle().
		Width(modalWidth).
		Align(lipgloss.Center).
		Render(
			lipgloss.JoinVertical(lipgloss.Center,
				lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true).Render("⚠️  WARNING: DATABASE RESET"),
				"\nThis will delete ALL locally cached tasks and settings.",
				"The application will need to re-sync everything.",
				"\nProceed? (y/n)",
			),
		)

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("#FF0000")).
		Padding(1, 4)

	modal := lipgloss.Place(m.width, paneHeight, lipgloss.Center, lipgloss.Center, modalStyle.Render(content))
	return modal
}
