package oauth

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

var (
	ClientID     string
	ClientSecret string
	RedirectURI  string
)

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

func Authenticate() {
	if ClientID == "" {
		ClientID = os.Getenv("CLIENT_ID")
	}
	if RedirectURI == "" {
		RedirectURI = os.Getenv("REDIRECT_URI")
	}
	url := fmt.Sprintf("https://app.clickup.com/api?client_id=%s&redirect_uri=%s", ClientID, RedirectURI)
	OpenBrowser(url)
}
