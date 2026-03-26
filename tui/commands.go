package tui

import (
	"fmt"
	"super-duper-fortnight/clkup"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/sync/errgroup"
)

func fetchPlanCmd(client *clkup.APIClient, teamID string) tea.Cmd {
	return func() tea.Msg {
		plan, err := client.GetPlan(teamID)
		if err != nil {
			return ErrMsg{err}
		}
		return PlanLoadedMsg{
			TeamID: teamID,
			PlanID: plan.PlanID,
		}
	}
}

func fetchInitDataCmd(client *clkup.APIClient) tea.Cmd {
	return func() tea.Msg {
		var user clkup.User
		var workspaces []clkup.Workspace
		var err error

		for attempts := 0; attempts < 3; attempts++ {
			user, err = client.GetAuthorizedUser()
			if err == nil {
				break
			}
			time.Sleep(1 * time.Second)
		}
		if err != nil {
			return ErrMsg{fmt.Errorf("failed to fetch user after 3 attempts: %w", err)}
		}

		for attempts := 0; attempts < 3; attempts++ {
			workspaces, err = client.GetAuthorizedWorkspaces()
			if err == nil {
				break
			}
			time.Sleep(1 * time.Second)
		}
		if err != nil {
			return ErrMsg{fmt.Errorf("workspace API error after 3 attempts: %w", err)}
		}

		if len(workspaces) == 0 {
			return ErrMsg{fmt.Errorf("success, but workspace array was empty")}
		}

		return InitDataMsg{
			User:       user,
			Workspaces: workspaces,
		}
	}
}

func fetchHierarchyCmd(client *clkup.APIClient, teamID string) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()
		var g errgroup.Group
		var mu sync.Mutex

		var finalSpaces []clkup.Space
		var finalTasks []clkup.Task
		foldersBySpace := make(map[string][]clkup.Folder)
		listsByFolder := make(map[string][]clkup.List)
		listsBySpace := make(map[string][]clkup.List)

		g.Go(func() error {
			tasks, err := client.GetAllTasks(teamID)
			if err == nil {
				mu.Lock()
				finalTasks = tasks
				mu.Unlock()
			}
			return err
		})

		g.Go(func() error {
			spaces, err := client.GetSpaces(teamID)
			if err != nil {
				return err
			}

			mu.Lock()
			finalSpaces = spaces
			mu.Unlock()

			for _, space := range spaces {
				sID := string(space.ID)

				g.Go(func() error {
					folders, err := client.GetFolders(sID)
					if err != nil {
						return err
					}

					mu.Lock()
					foldersBySpace[sID] = folders
					mu.Unlock()

					for _, folder := range folders {
						fID := string(folder.ID)
						g.Go(func() error {
							lists, err := client.GetLists(fID)
							if err != nil {
								return err
							}

							mu.Lock()
							listsByFolder[fID] = lists
							mu.Unlock()
							return nil
						})
					}
					return nil
				})

				g.Go(func() error {
					folderlessLists, err := client.GetFolderlessLists(sID)
					if err != nil {
						return err
					}

					mu.Lock()
					listsBySpace[sID] = folderlessLists
					mu.Unlock()
					return nil
				})
			}

			return nil
		})

		if err := g.Wait(); err != nil {
			return ErrMsg{err}
		}

		tasksByList := make(map[string][]clkup.Task)
		for _, t := range finalTasks {
			lID := getListIDFromTask(t)
			if lID != "" {
				tasksByList[lID] = append(tasksByList[lID], t)
			}
		}

		perf := clkup.CalculatePerformance(len(finalTasks), start)

		wd := &WorkspaceData{
			Spaces:         finalSpaces,
			FoldersBySpace: foldersBySpace,
			ListsByFolder:  listsByFolder,
			ListsBySpace:   listsBySpace,
			Tasks:          finalTasks,
			TasksByList:    tasksByList,
			Performance:    perf,
		}

		return FanOutCompleteMsg{
			TeamID: teamID,
			Data:   wd,
		}
	}
}
