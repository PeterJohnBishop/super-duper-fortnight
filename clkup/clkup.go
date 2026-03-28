package clkup

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type APIClient struct {
	Client  *http.Client
	Token   string
	Limiter *rate.Limiter
	LogChan chan string
}

// 1=Free (100), 2=Unlimited (100), 3=Business (1000), 4=Enterprise (10000)
func GetRateLimit(planID int) rate.Limit {
	var rpm int
	switch planID {
	case 3:
		rpm = 1000
	case 4:
		rpm = 10000
	default:
		rpm = 100
	}
	safeRpm := float64(rpm) * 0.95
	return rate.Every(time.Minute / time.Duration(safeRpm))
}

func (c *APIClient) Do(req *http.Request) (*http.Response, error) {
	err := c.Limiter.Wait(context.Background())
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", c.Token)
	req.Header.Add("Content-Type", "application/json")

	if c.LogChan != nil {
		// Use a non-blocking select so the API engine never stalls
		// if the UI happens to fall behind in reading the logs.
		select {
		case c.LogChan <- fmt.Sprintf("[%s] %s %s", time.Now().Format("15:04:05.000"), req.Method, req.URL.Path):
		default:
		}
	}

	return c.Client.Do(req)
}
