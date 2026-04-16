package tui

import (
	"github.com/charmbracelet/lipgloss"
)

func buildCenterContent(m dashboardModel, occupiedHeight int, safeTextWidth int) string {
	paneHeight := m.height - occupiedHeight - 1
	if paneHeight < 5 {
		paneHeight = 5
	}

	paneWidthLeft := safeTextWidth / 2
	paneWidthRight := safeTextWidth - paneWidthLeft

	if paneWidthLeft < 10 {
		paneWidthLeft = 10
	}
	if paneWidthRight < 10 {
		paneWidthRight = 10
	}

	leftItems, leftTitle, leftCursor := m.getLeftPane()
	rightItems, rightTitle, rightRawText := m.getRightPane()

	leftActive := !m.focusRight && (m.depth != DepthTaskDetails)
	rightActive := m.focusRight || (m.depth == DepthTaskDetails)
	leftPane := renderPane(leftItems, leftTitle, "", leftCursor, 0, paneWidthLeft, paneHeight, leftActive)
	rightPane := renderPane(rightItems, rightTitle, rightRawText, -1, m.taskScrollOffset, paneWidthRight, paneHeight, rightActive)

	var centerContent string

	if m.state == stateResetConfirm {
		centerContent = buildResyncModal(m, safeTextWidth, paneHeight)
	} else if m.showJSONPopup {
		centerContent = buildJSONView(m, safeTextWidth, paneHeight)
	} else {
		splitPanes := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)
		centerContent = splitPanes
	}
	return centerContent
}
