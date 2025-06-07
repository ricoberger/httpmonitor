package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ricoberger/httpmonitor/pkg/config"
	"github.com/ricoberger/httpmonitor/pkg/target"
	"github.com/ricoberger/httpmonitor/pkg/ui"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Failed to get home directory: %#v", err)
	}

	var configFile string
	var url string
	var interval time.Duration
	var timeout time.Duration
	var method string
	var body string
	var username string
	var password string
	var token string
	var notification bool
	var notificationThreshold time.Duration

	flag.StringVar(&configFile, "config", home+"/.httpmonitor.yaml", "The path to the configuration file.")
	flag.StringVar(&url, "url", "", "The url of the target to monitor.")
	flag.StringVar(&method, "method", http.MethodGet, "The HTTP method to use for the checks.")
	flag.StringVar(&body, "body", "", "The body to send with the HTTP checks.")
	flag.StringVar(&username, "username", "", "The username which should used, when the target requires basic authentication.")
	flag.StringVar(&password, "password", "", "The password which should used, when the target requires basic authentication.")
	flag.StringVar(&token, "token", "", "The token which should used, when the target requires token authentication.")
	flag.BoolVar(&notification, "notification", false, "Enable desktop notifications, for failed checks.")
	flag.DurationVar(&notificationThreshold, "notification-threshold", 0*time.Second, "Enable desktop notifications, for checks which are longer than the threshold.")
	flag.DurationVar(&interval, "interval", 5*time.Second, "The interval to run the HTTP checks.")
	flag.DurationVar(&timeout, "timeout", 2*time.Second, "The timeout for the HTTP checks.")
	flag.Parse()

	config, err := getConfig(configFile, url, method, body, username, password, token, notification, notificationThreshold, interval, timeout)
	if err != nil {
		log.Printf("Failed to load configuration: %#v", err)
		os.Exit(1)
	}

	var targets []target.Client

	for _, t := range config.Targets {
		client := target.NewClient(t)
		go client.Start()
		targets = append(targets, client)
	}

	if err := ui.Start(targets); err != nil {
		log.Printf("Failed to start the UI: %#v", err)
		os.Exit(1)
	}
}

// If the url flag is set, the function will return a config with just the
// target from the flag. Otherwise we will return the config from the file.
func getConfig(file, url, method, body, username, password, token string, notification bool, notificationThreshold, interval, timeout time.Duration) (*config.Config, error) {
	if url != "" {
		return &config.Config{
			Targets: []target.Config{{
				Name:                  url,
				URL:                   url,
				Method:                method,
				Body:                  body,
				Username:              username,
				Password:              password,
				Token:                 token,
				Notification:          notification,
				NotificationThreshold: notificationThreshold,
				Interval:              interval,
				Timeout:               timeout,
			}},
		}, nil
	}

	config, err := config.New(file)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(config.Targets); i++ {
		if config.Targets[i].Name == "" || config.Targets[i].URL == "" {
			return nil, fmt.Errorf("name and url are required for all targets")
		}
		if config.Targets[i].Method == "" {
			config.Targets[i].Method = method
		}
		if config.Targets[i].Interval == 0 {
			config.Targets[i].Interval = interval
		}
		if config.Targets[i].Timeout == 0 {
			config.Targets[i].Timeout = timeout
		}
	}

	return config, nil
}
