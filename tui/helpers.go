package tui

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"super-duper-fortnight/clkup"

	"github.com/charmbracelet/lipgloss"
)

func (m dashboardModel) getLeftPane() ([]ListItem, string, int) {
	wd := m.workspaceCache[m.activeTeamID]

	switch m.depth {
	case DepthWorkspaces:
		var items []ListItem
		for _, w := range m.workspaces {
			items = append(items, ListItem{ID: string(w.ID), Name: w.Name, Type: "workspace"})
		}
		return items, "Workspaces", m.cursorWorkspace

	case DepthSpaces:
		if wd == nil {
			return nil, "Spaces", 0
		}
		var items []ListItem
		for _, s := range wd.Spaces {
			items = append(items, ListItem{ID: string(s.ID), Name: s.Name, Type: "space"})
		}
		return items, "Spaces", m.cursorSpace

	case DepthFolders:
		if wd == nil || len(wd.Spaces) == 0 {
			return nil, "Folders & Standalone Lists", 0
		}
		var items []ListItem
		sID := string(wd.Spaces[m.cursorSpace].ID)
		for _, f := range wd.FoldersBySpace[sID] {
			items = append(items, ListItem{ID: string(f.ID), Name: fmt.Sprintf("📁 %s", f.Name), Type: "folder"})
		}
		for _, l := range wd.ListsBySpace[sID] {
			items = append(items, ListItem{ID: string(l.ID), Name: fmt.Sprintf("📄 %s", l.Name), Type: "list"})
		}
		return items, "Folders & Standalone Lists", m.cursorFolder

	case DepthLists:
		if wd == nil || len(wd.Spaces) == 0 {
			return nil, "Lists", 0
		}
		var items []ListItem
		sID := string(wd.Spaces[m.cursorSpace].ID)
		folders := wd.FoldersBySpace[sID]
		if m.cursorFolder < len(folders) {
			fID := string(folders[m.cursorFolder].ID)
			for _, l := range wd.ListsByFolder[fID] {
				items = append(items, ListItem{ID: string(l.ID), Name: fmt.Sprintf("📄 %s", l.Name), Type: "list"})
			}
		}
		return items, "Lists", m.cursorList

	case DepthTasks, DepthTaskDetails:
		if wd == nil {
			return nil, "Tasks", 0
		}
		var items []ListItem
		listID := m.getActiveListID()
		for _, t := range wd.TasksByList[listID] {
			items = append(items, ListItem{ID: string(t.Id), Name: t.Name, Type: "task", Subtitle: t.Status.Status})
		}
		return items, "Tasks", m.cursorTask
	}
	return nil, "", 0
}

func (m dashboardModel) getRightPane() ([]ListItem, string, string) {
	wd := m.workspaceCache[m.activeTeamID]

	switch m.depth {
	case DepthWorkspaces:
		if len(m.workspaces) == 0 {
			return nil, "Instructions", "\n  <-- Press Enter to fetch Workspace data."
		}
		hoveredWS := string(m.workspaces[m.cursorWorkspace].ID)
		if hwd, ok := m.workspaceCache[hoveredWS]; ok {
			var items []ListItem
			for _, s := range hwd.Spaces {
				items = append(items, ListItem{ID: string(s.ID), Name: s.Name, Type: "space"})
			}
			return items, "Spaces Preview", ""
		}
		return nil, "Instructions", "\n  <-- Press Enter to fetch Workspace data."

	case DepthSpaces:
		if wd == nil || len(wd.Spaces) == 0 || m.cursorSpace >= len(wd.Spaces) {
			return nil, "Folders & Standalone Lists", ""
		}
		var items []ListItem
		sID := string(wd.Spaces[m.cursorSpace].ID)
		for _, f := range wd.FoldersBySpace[sID] {
			items = append(items, ListItem{ID: string(f.ID), Name: fmt.Sprintf("📁 %s", f.Name), Type: "folder"})
		}
		for _, l := range wd.ListsBySpace[sID] {
			items = append(items, ListItem{ID: string(l.ID), Name: fmt.Sprintf("📄 %s", l.Name), Type: "list"})
		}
		return items, "Folders & Standalone Lists", ""

	case DepthFolders:
		if wd == nil || len(wd.Spaces) == 0 {
			return nil, "", ""
		}
		var items []ListItem
		sID := string(wd.Spaces[m.cursorSpace].ID)
		folders := wd.FoldersBySpace[sID]

		if m.cursorFolder < len(folders) {
			fID := string(folders[m.cursorFolder].ID)
			for _, l := range wd.ListsByFolder[fID] {
				items = append(items, ListItem{ID: string(l.ID), Name: fmt.Sprintf("📄 %s", l.Name), Type: "list"})
			}
			return items, "Lists", ""
		} else {
			idx := m.cursorFolder - len(folders)
			lists := wd.ListsBySpace[sID]
			if idx >= 0 && idx < len(lists) {
				lID := string(lists[idx].ID)
				for _, t := range wd.TasksByList[lID] {
					items = append(items, ListItem{ID: string(t.Id), Name: t.Name, Type: "task", Subtitle: t.Status.Status})
				}
			}
			return items, "Tasks", ""
		}

	case DepthLists:
		if wd == nil {
			return nil, "Tasks", ""
		}
		var items []ListItem
		lID := m.getHoveredListID()
		for _, t := range wd.TasksByList[lID] {
			items = append(items, ListItem{ID: string(t.Id), Name: t.Name, Type: "task", Subtitle: t.Status.Status})
		}
		return items, "Tasks", ""

	case DepthTasks, DepthTaskDetails:
		t := m.getHoveredTask()
		if t != nil {
			b, _ := json.MarshalIndent(t, "", "  ")
			return nil, "Task Details JSON", string(b)
		}
		return nil, "Task Details", "No task selected"
	}
	return nil, "", ""
}

func (m dashboardModel) getCurrentSelectionURL() string {
	if m.state != stateLoaded && m.state != stateIdle {
		return ""
	}
	teamID := m.activeTeamID

	switch m.depth {
	case DepthWorkspaces:
		if len(m.workspaces) > 0 && m.cursorWorkspace < len(m.workspaces) {
			return fmt.Sprintf("https://app.clickup.com/%s", m.workspaces[m.cursorWorkspace].ID)
		}
	case DepthSpaces:
		wd := m.workspaceCache[teamID]
		if wd != nil && m.cursorSpace < len(wd.Spaces) {
			return fmt.Sprintf("https://app.clickup.com/%s/v/s/%s", teamID, wd.Spaces[m.cursorSpace].ID)
		}
	case DepthFolders:
		wd := m.workspaceCache[teamID]
		if wd != nil && len(wd.Spaces) > 0 {
			sID := string(wd.Spaces[m.cursorSpace].ID)
			folders := wd.FoldersBySpace[sID]
			if m.cursorFolder < len(folders) {
				return fmt.Sprintf("https://app.clickup.com/%s/v/f/%s", teamID, folders[m.cursorFolder].ID)
			} else {
				idx := m.cursorFolder - len(folders)
				lists := wd.ListsBySpace[sID]
				if idx >= 0 && idx < len(lists) {
					return fmt.Sprintf("https://app.clickup.com/%s/v/l/li/%s", teamID, lists[idx].ID)
				}
			}
		}
	case DepthLists:
		wd := m.workspaceCache[teamID]
		if wd != nil {
			lID := m.getHoveredListID()
			if lID != "" {
				return fmt.Sprintf("https://app.clickup.com/%s/v/l/li/%s", teamID, lID)
			}
		}
	case DepthTasks, DepthTaskDetails:
		t := m.getHoveredTask()
		if t != nil {
			return fmt.Sprintf("https://app.clickup.com/t/%s", t.Id)
		}
	}
	return ""
}

func OpenBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	case "linux":
		cmd = "xdg-open"
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func (m dashboardModel) getStats() (spaces, folders, lists, tasks string) {
	teamID := m.activeTeamID
	if m.depth == DepthWorkspaces {
		if len(m.workspaces) == 0 {
			return "-", "-", "-", "-"
		}
		teamID = string(m.workspaces[m.cursorWorkspace].ID)
	}

	wd, exists := m.workspaceCache[teamID]
	if !exists {
		return "-", "-", "-", "-"
	}

	switch m.depth {
	case DepthWorkspaces:
		fCount, lCount := 0, 0
		for _, f := range wd.FoldersBySpace {
			fCount += len(f)
		}
		for _, l := range wd.ListsByFolder {
			lCount += len(l)
		}
		for _, l := range wd.ListsBySpace {
			lCount += len(l)
		}
		return fmt.Sprint(len(wd.Spaces)), fmt.Sprint(fCount), fmt.Sprint(lCount), fmt.Sprint(len(wd.Tasks))

	case DepthSpaces:
		if len(wd.Spaces) == 0 {
			return "0", "0", "0", "0"
		}
		sID := string(wd.Spaces[m.cursorSpace].ID)
		fCount := len(wd.FoldersBySpace[sID])
		lCount := len(wd.ListsBySpace[sID])
		for _, f := range wd.FoldersBySpace[sID] {
			lCount += len(wd.ListsByFolder[string(f.ID)])
		}
		tCount := 0
		var listIDs []string
		for _, l := range wd.ListsBySpace[sID] {
			listIDs = append(listIDs, string(l.ID))
		}
		for _, f := range wd.FoldersBySpace[sID] {
			for _, l := range wd.ListsByFolder[string(f.ID)] {
				listIDs = append(listIDs, string(l.ID))
			}
		}
		for _, lid := range listIDs {
			tCount += len(wd.TasksByList[lid])
		}
		return "1", fmt.Sprint(fCount), fmt.Sprint(lCount), fmt.Sprint(tCount)

	case DepthFolders:
		if len(wd.Spaces) == 0 {
			return "-", "0", "0", "0"
		}
		sID := string(wd.Spaces[m.cursorSpace].ID)
		folders := wd.FoldersBySpace[sID]

		if m.cursorFolder < len(folders) {
			fID := string(folders[m.cursorFolder].ID)
			lCount := len(wd.ListsByFolder[fID])
			tCount := 0
			for _, l := range wd.ListsByFolder[fID] {
				tCount += len(wd.TasksByList[string(l.ID)])
			}
			return "-", "1", fmt.Sprint(lCount), fmt.Sprint(tCount)
		} else {
			idx := m.cursorFolder - len(folders)
			lists := wd.ListsBySpace[sID]
			if idx >= 0 && idx < len(lists) {
				lID := string(lists[idx].ID)
				tCount := len(wd.TasksByList[lID])
				return "-", "-", "1", fmt.Sprint(tCount)
			}
			return "-", "-", "-", "-"
		}

	case DepthLists:
		lID := m.getHoveredListID()
		if lID != "" {
			tCount := len(wd.TasksByList[lID])
			return "-", "-", "1", fmt.Sprint(tCount)
		}
		return "-", "-", "-", "-"

	case DepthTasks, DepthTaskDetails:
		lID := m.getActiveListID()
		if lID != "" {
			tCount := len(wd.TasksByList[lID])
			return "-", "-", "-", fmt.Sprint(tCount)
		}
		return "-", "-", "-", "-"
	}

	return "-", "-", "-", "-"
}

func (m dashboardModel) getActiveListID() string {
	wd := m.workspaceCache[m.activeTeamID]
	if wd == nil || len(wd.Spaces) == 0 {
		return ""
	}
	sID := string(wd.Spaces[m.cursorSpace].ID)
	folders := wd.FoldersBySpace[sID]
	if m.cursorFolder < len(folders) {
		fID := string(folders[m.cursorFolder].ID)
		lists := wd.ListsByFolder[fID]
		if m.cursorList < len(lists) {
			return string(lists[m.cursorList].ID)
		}
	} else {
		idx := m.cursorFolder - len(folders)
		lists := wd.ListsBySpace[sID]
		if idx >= 0 && idx < len(lists) {
			return string(lists[idx].ID)
		}
	}
	return ""
}

func (m dashboardModel) getHoveredListID() string {
	wd := m.workspaceCache[m.activeTeamID]
	if wd == nil || len(wd.Spaces) == 0 {
		return ""
	}
	sID := string(wd.Spaces[m.cursorSpace].ID)
	folders := wd.FoldersBySpace[sID]
	if m.cursorFolder < len(folders) {
		fID := string(folders[m.cursorFolder].ID)
		lists := wd.ListsByFolder[fID]
		if m.cursorList < len(lists) {
			return string(lists[m.cursorList].ID)
		}
	}
	return ""
}

func (m dashboardModel) getHoveredTask() *clkup.Task {
	wd := m.workspaceCache[m.activeTeamID]
	if wd == nil {
		return nil
	}
	listID := m.getActiveListID()
	tasksInList := wd.TasksByList[listID]

	if m.cursorTask >= 0 && m.cursorTask < len(tasksInList) {
		return &tasksInList[m.cursorTask]
	}
	return nil
}

func getListIDFromTask(t clkup.Task) string {
	b, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	var temp map[string]interface{}
	json.Unmarshal(b, &temp)

	if listObj, ok := temp["list"].(map[string]interface{}); ok {
		if id, ok := listObj["id"].(string); ok {
			return id
		}
	}
	return ""
}

func (m dashboardModel) getBreadcrumbs() string {
	if m.state != stateLoaded && m.state != stateIdle {
		return ""
	}
	crumbs := []string{}
	if len(m.workspaces) > 0 && m.cursorWorkspace < len(m.workspaces) {
		crumbs = append(crumbs, m.workspaces[m.cursorWorkspace].Name)
	}

	wd := m.workspaceCache[m.activeTeamID]
	if wd == nil {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#9D4EDD")).Render(strings.Join(crumbs, " > "))
	}

	if m.depth >= DepthSpaces && m.cursorSpace < len(wd.Spaces) {
		crumbs = append(crumbs, wd.Spaces[m.cursorSpace].Name)
	}
	if m.depth >= DepthFolders && len(wd.Spaces) > 0 {
		sID := string(wd.Spaces[m.cursorSpace].ID)
		folders := wd.FoldersBySpace[sID]
		if m.cursorFolder < len(folders) {
			crumbs = append(crumbs, folders[m.cursorFolder].Name)
		} else {
			idx := m.cursorFolder - len(folders)
			lists := wd.ListsBySpace[sID]
			if idx >= 0 && idx < len(lists) {
				crumbs = append(crumbs, lists[idx].Name)
			}
		}
	}
	if m.depth >= DepthLists && len(wd.Spaces) > 0 {
		sID := string(wd.Spaces[m.cursorSpace].ID)
		folders := wd.FoldersBySpace[sID]
		if m.cursorFolder < len(folders) {
			fID := string(folders[m.cursorFolder].ID)
			lists := wd.ListsByFolder[fID]
			if m.cursorList < len(lists) {
				crumbs = append(crumbs, lists[m.cursorList].Name)
			}
		}
	}
	if m.depth >= DepthTasks {
		t := m.getHoveredTask()
		if t != nil {
			crumbs = append(crumbs, t.Name)
		}
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#9D4EDD")).Render(strings.Join(crumbs, " > "))
}

func renderPane(items []ListItem, title string, rawText string, cursor int, scrollOffset int, width int, height int, isActive bool) string {
	innerW := width - 2
	if innerW < 5 {
		innerW = 5
	}

	innerH := height - 2
	if innerH < 3 {
		innerH = 3
	}

	paneStyle := lipgloss.NewStyle().
		Width(innerW).
		MaxWidth(innerW).
		Height(innerH).
		MaxHeight(innerH).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#5A189A"))

	if isActive {
		paneStyle = paneStyle.BorderForeground(lipgloss.Color("#E0AAFF"))
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7B2CBF"))

	var uiLines []string

	titleRunes := []rune(title)
	if len(titleRunes) > innerW {
		title = string(titleRunes[:innerW-3]) + "..."
	}
	uiLines = append(uiLines, titleStyle.Render(title))
	uiLines = append(uiLines, "")
	if rawText != "" {
		lines := strings.Split(rawText, "\n")
		maxLines := innerH - 2
		if maxLines <= 0 {
			maxLines = 1
		}

		startIdx := scrollOffset
		if startIdx > len(lines)-maxLines {
			startIdx = len(lines) - maxLines
		}
		if startIdx < 0 {
			startIdx = 0
		}

		endIdx := startIdx + maxLines
		if endIdx > len(lines) {
			endIdx = len(lines)
		}

		if startIdx < len(lines) {
			for _, line := range lines[startIdx:endIdx] {
				line = strings.ReplaceAll(line, "\t", "  ")

				runes := []rune(line)
				if len(runes) > innerW {
					uiLines = append(uiLines, string(runes[:innerW-3])+"...")
				} else {
					uiLines = append(uiLines, line)
				}
			}
		}

	} else {
		maxLines := innerH - 2
		if maxLines <= 0 {
			maxLines = 1
		}
		startIdx := 0
		if cursor >= maxLines {
			startIdx = cursor - maxLines + 1
		}
		for i := startIdx; i < len(items) && i < startIdx+maxLines; i++ {
			item := items[i]
			prefix := "  "
			style := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

			if i == cursor && isActive {
				prefix = "> "
				style = style.Foreground(lipgloss.Color("#E0AAFF")).Bold(true)
			} else if i == cursor && !isActive {
				style = style.Foreground(lipgloss.Color("#9D4EDD"))
			}

			nameStr := item.Name
			if item.Subtitle != "" {
				nameStr += " [" + item.Subtitle + "]"
			}

			maxNameW := innerW - 2
			if maxNameW < 3 {
				maxNameW = 3
			}

			runes := []rune(nameStr)
			if len(runes) > maxNameW {
				nameStr = string(runes[:maxNameW-3]) + "..."
			}

			uiLines = append(uiLines, prefix+style.Render(nameStr))
		}
	}

	content := strings.Join(uiLines, "\n")
	return paneStyle.Render(content)
}
