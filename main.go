package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/overjt/webmonitor/coreapp"
)

type MonitorService struct {
	URL        string    `json:"url"`
	Timeout    int       `json:"timeout"`
	Interval   int       `json:"interval"`
	Emails     []string  `json:"emails"`
	SmsNumbers []string  `json:"smsNumbers"`
	Enabled    bool      `json:"enabled"`
	Name       string    `json:"name"`
	LastCheck  time.Time `json:"lastCheck"`
}

type Config struct {
	Services []MonitorService `json:"services"`
	CoreApp  coreapp.CoreApp  `json:"coreApp"`
}

//Load Json file to Coreapp struct
func LoadConfig(filename string) *Config {
	//Load json file
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	//Unmarshal json file
	var config Config
	json.Unmarshal(file, &config)
	return &config
}

//Returns True if the url responds within the timeout
func checkUrl(url string, timeout int) bool {
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func main() {
	//Load config file
	config := LoadConfig("config.json")
	wg := sync.WaitGroup{}
	for _, service := range config.Services {
		wg.Add(1)
		go func(service *MonitorService) {
			for {
				if !checkUrl(service.URL, service.Timeout) {
					//Send email
					message := fmt.Sprintf("%s is down", service.URL)
					subject := fmt.Sprintf("[%s] Website down, Last Check %s", service.Name, service.LastCheck)
					if len(service.Emails) > 0 {
						config.CoreApp.SendEmail(service.Emails, message, subject)
					}

					//Send SMS
					if len(service.SmsNumbers) > 0 {
						config.CoreApp.SendSMS(service.SmsNumbers, message)
					}

				} else {
					//save last check date
					service.LastCheck = time.Now()
				}
				time.Sleep(time.Duration(service.Interval) * time.Second)
			}
		}(&service)
	}
	wg.Wait()

}
