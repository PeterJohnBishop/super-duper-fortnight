package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func buildHierarchyStats(m dashboardModel) string {
	sCount, fCount, lCount, tCount := m.getStats()
	statSpaces := statBoxStyle.Render(fmt.Sprintf("%s\n%s", statLabelStyle.Render("Spaces"), statValueStyle.Render(sCount)))
	statFolders := statBoxStyle.Render(fmt.Sprintf("%s\n%s", statLabelStyle.Render("Folders"), statValueStyle.Render(fCount)))
	statLists := statBoxStyle.Render(fmt.Sprintf("%s\n%s", statLabelStyle.Render("Lists"), statValueStyle.Render(lCount)))
	statTasks := statBoxStyle.Render(fmt.Sprintf("%s\n%s", statLabelStyle.Render("Tasks"), statValueStyle.Render(tCount)))

	statsRow := lipgloss.NewStyle().MarginBottom(1).Render(
		lipgloss.JoinHorizontal(lipgloss.Top, statSpaces, statFolders, statLists, statTasks),
	)

	return statsRow
}
