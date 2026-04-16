package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func buildHeader(m dashboardModel, safeTextWidth int) string {
	var header string
	if m.state != stateInit {
		headerInnerWidth := safeTextWidth - 2

		leftSide := fmt.Sprintf("%s | %s - %s", m.user.ID, m.user.Initials, m.user.Email)
		rightTop := fmt.Sprintf("[%s]", m.user.Timezone)

		syncColor := "#5A189A"
		if m.syncInterval != SyncOff {
			syncColor = "#00FF00"
		}
		rightBot := lipgloss.NewStyle().Foreground(lipgloss.Color(syncColor)).Render(m.syncInterval.String())

		leftW := lipgloss.Width(leftSide)
		rightTopW := lipgloss.Width(rightTop)

		if leftW+rightTopW >= headerInnerWidth {
			availLeft := headerInnerWidth - rightTopW - 1
			if availLeft > 3 {
				runes := []rune(leftSide)
				leftSide = string(runes[:availLeft-1]) + "…"
			} else {
				leftSide = ""
			}
		}

		leftW = lipgloss.Width(leftSide)
		spaceCountTop := headerInnerWidth - leftW - rightTopW
		if spaceCountTop < 0 {
			spaceCountTop = 0
		}
		line1 := leftSide + strings.Repeat(" ", spaceCountTop) + rightTop

		rightBotW := lipgloss.Width(rightBot)
		spaceCountBot := headerInnerWidth - rightBotW
		if spaceCountBot < 0 {
			spaceCountBot = 0
		}
		line2 := strings.Repeat(" ", spaceCountBot) + rightBot

		headerStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E0AAFF")).
			Padding(0, 1).
			MarginBottom(1).
			Italic(true)

		headerContent := line1 + "\n" + line2
		header = headerStyle.Render(headerContent)
	}

	return header
}
