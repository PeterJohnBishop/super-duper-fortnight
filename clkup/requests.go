package clkup

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func calculatePerformance(totalTasks int, start time.Time) Performance {
	elapsed := time.Since(start)
	pages := float64(totalTasks) / 100.0
	if pages < 1 {
		pages = 1
	}

	rpm := (pages / elapsed.Minutes())
	tps := float64(totalTasks) / elapsed.Seconds()

	return Performance{
		Duration: elapsed.Round(time.Millisecond).String(),
		RPM:      fmt.Sprintf("%.2f", rpm),
		TPS:      fmt.Sprintf("%.2f", tps),
	}
}

// oauth

func GetAccessToken(code string) (string, error) {
	client_id := os.Getenv("CLIENT_ID")
	client_secret := os.Getenv("CLIENT_SECRET")

	url := "https://api.clickup.com/api/v2/oauth/token"

	reqData := map[string]string{
		"client_id":     client_id,
		"client_secret": client_secret,
		"code":          code}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// user

func GetAuthorizedUser(token string) (User, error) {
	url := "https://api.clickup.com/api/v2/user"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return User{}, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return User{}, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var userResponse UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
		return User{}, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return userResponse.User, nil
}

// workspace

func GetAuthorizedWorkspaces(token string) ([]Workspace, error) {

	url := "https://api.clickup.com/api/v2/team"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)

	}

	var teamsResponse TeamsResponse
	if err := json.NewDecoder(resp.Body).Decode(&teamsResponse); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return teamsResponse.Teams, nil
}

// space

func GetSpaces(token string, teamID string) ([]Space, error) {

	url := fmt.Sprintf("https://api.clickup.com/api/v2/team/%s/space", teamID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var spaceResponse SpacesResponse
	if err := json.NewDecoder(resp.Body).Decode(&spaceResponse); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}
	return spaceResponse.Spaces, nil
}

// folder

func GetFolders(token string, spaceID string) ([]Folder, error) {
	url := fmt.Sprintf("https://api.clickup.com/api/v2/space/%s/folder", spaceID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var foldersResponse FoldersResponse
	if err := json.NewDecoder(resp.Body).Decode(&foldersResponse); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return foldersResponse.Folders, nil
}

// folderless list

func GetFolderlessLists(token string, spaceID string) ([]List, error) {
	url := fmt.Sprintf("https://api.clickup.com/api/v2/space/%s/list", spaceID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var listsResponse ListsResponse
	if err := json.NewDecoder(resp.Body).Decode(&listsResponse); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return listsResponse.Lists, nil
}

// list

func GetLists(token string, folderID string) ([]List, error) {
	url := fmt.Sprintf("https://api.clickup.com/api/v2/folder/%s/list", folderID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var listsResponse ListsResponse
	if err := json.NewDecoder(resp.Body).Decode(&listsResponse); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return listsResponse.Lists, nil
}

// task

func GetAllTasks(teamID string, token string) ([]Task, error) {
	var allTasks []Task
	var wg sync.WaitGroup

	limiter := rate.NewLimiter(rate.Every(time.Minute/1000), 1)

	taskChan := make(chan []Task, 100)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for tasks := range taskChan {
			allTasks = append(allTasks, tasks...)
		}
	}()

	// fmt.Println("starting concurrent fetch...")
	start := time.Now()

	for page := 0; ; page++ {
		select {
		case <-ctx.Done():
			goto WaitAndFinish
		default:
		}

		if err := limiter.Wait(context.Background()); err != nil {
			break
		}

		wg.Add(1)
		go func(p int) {
			defer wg.Done()

			tasks, err := fetchPage(teamID, token, p)
			if err != nil || len(tasks) == 0 {
				cancel()
				return
			}

			taskChan <- tasks
		}(page)

	}

WaitAndFinish:
	wg.Wait()
	close(taskChan)

	time.Sleep(100 * time.Millisecond)

	performance := calculatePerformance(len(allTasks), start)
	fmt.Printf("fetched %d tasks in %s seconds. RPM: %s\n",
		len(allTasks), performance.Duration, performance.RPM)

	return allTasks, nil
}

func fetchPage(teamID string, token string, page int) ([]Task, error) {
	url := fmt.Sprintf("https://api.clickup.com/api/v2/team/%s/task?page=%d&include_closed=true&subtasks=true", teamID, page)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		return nil, fmt.Errorf("rate limit hit")
	}

	var tasksResponse TasksResponse
	if err := json.NewDecoder(resp.Body).Decode(&tasksResponse); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return tasksResponse.Task, nil
}
