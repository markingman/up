package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

const configPath = "/app/data/conf.json"

type Config struct {
	To     string   `json:"to"`
	Tick   uint16   `json:"tick"`
	Checks []string `json:"checks"`
	SMTP   SMTP     `json:"smtp"`
	Sites  []Site   `json:"sites"`
}

type SMTP struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

type Site struct {
	URL string `json:"url"`
}

type Result struct {
	httpCode int
	url      string
	complete bool
}

func getConfig(configPath string) Config {
	config := Config{}

	file, err := os.Open(configPath)
	if err != nil {
		log.Panic("Error opening config:", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Panic("Error closing config:", err)
		}
	}(file)

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		log.Panic("Error decoding JSON:", err)
	}

	if config.To == "" {
		log.Println("email recipient (To:) not set")
	}

	if config.Tick == 0 {
		log.Println("tick (seconds) not set, defaulting to 600")
		config.Tick = 600
	}

	if len(config.Checks) == 0 {
		log.Println("check times not set")
	}

	if config.SMTP.Host == "" {
		log.Println("SMTP host not set")
	}

	if config.SMTP.Port == 0 {
		log.Println("SMTP port not set")
	}

	if config.SMTP.User == "" {
		log.Println("SMTP user not set")
	}

	if config.SMTP.Pass == "" {
		log.Println("SMTP password not set")
	}

	if len(config.Sites) == 0 {
		log.Println("no sites set")
	}

	for i, site := range config.Sites {
		if site.URL == "" {
			log.Println(fmt.Sprintf("Site %d is not set", i))
		}
	}

	return config
}

func sendEmail(to string, from string, subject string, body string, configSMTP SMTP) error {
	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n\r\n", from, to, subject, body)

	// connection -- note tls.Dial not smtp.Dial
	conn, err := tls.Dial("tcp",
		fmt.Sprintf("%s:%d", configSMTP.Host, configSMTP.Port),
		&tls.Config{ServerName: configSMTP.Host})

	if err != nil {
		log.Panic("Error dialing SMTP server:", err)
	}

	// client
	c, err := smtp.NewClient(conn, configSMTP.Host)
	if err != nil {
		log.Panic("Error creating SMTP client:", err)
	}

	// auth
	if err = c.Auth(smtp.PlainAuth("", configSMTP.User, configSMTP.Pass, configSMTP.Host)); err != nil {
		log.Panic("Error authenticating to SMTP:", err)
	}

	// To: and From:
	if err = c.Mail(from); err != nil {
		log.Panic("Error creating From: recipient:", err)
	}

	if err = c.Rcpt(to); err != nil {
		log.Panic("Error creating To: recipient:", err)
	}

	// write email
	w, err := c.Data()
	if err != nil {
		log.Panic("Error initiating DATA command with SMTP server:", err)
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		log.Panic("Error writing email body:", err)
	}

	err = w.Close()
	if err != nil {
		log.Panic("Error closing DATA command with SMTP server:", err)
	}

	// exit
	err = c.Quit()
	if err != nil {
		log.Panic("Error sending QUIT command:", err)
	}

	log.Println("Email sent")

	return nil
}

func sendEmailUpActive(config Config) error {
	to := config.To
	from := config.SMTP.User
	subject := "UP ACTIVE"
	body := "UP ACTIVE\n\n"

	for _, site := range config.Sites {
		body += site.URL + "\n"
	}

	return sendEmail(to, from, subject, body, config.SMTP)
}

func sendErrsEmail(config Config, results []Result) error {
	to := config.To
	from := config.SMTP.User
	subject := ""
	body := ""

	lenErrs := 0
	for _, result := range results {
		if result.httpCode == 0 || !result.complete {
			body += result.url + "\n"
			lenErrs++
		}
	}

	if lenErrs == 1 {
		subject += "ALERT: 1 SERVER DOWN"
	} else {
		subject += fmt.Sprintf("ALERT: %d SERVERS DOWN", lenErrs)
	}

	return sendEmail(to, from, subject, body, config.SMTP)
}

func callURL(url string, ch chan<- Result) {
	httpCode := 0

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error calling '%s': %s", url, err)
	} else {
		httpCode = resp.StatusCode
		err := resp.Body.Close()
		if err != nil {
			log.Println("Error closing response body:", err)
		}
	}

	ch <- Result{
		httpCode: httpCode,
		url:      url,
		complete: true,
	}
}

func nextConfirmTime(now time.Time, schedule []string) time.Duration {
	var nextTime time.Time

	for _, t := range schedule {
		parsedTime, err := time.Parse("15:04", t)
		if err != nil {
			fmt.Println("Error parsing time:", err)
			continue
		}

		scheduledTime := time.Date(now.Year(), now.Month(), now.Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, now.Location())

		if scheduledTime.Before(now) {
			scheduledTime = scheduledTime.Add(24 * time.Hour)
		}

		if nextTime.IsZero() || scheduledTime.Before(nextTime) {
			nextTime = scheduledTime
		}
	}

	duration := nextTime.Sub(now)
	if duration == 0 {
		duration = 3600
	}

	return duration
}

func checkSites(config Config) {
	var results []Result

	for _, site := range config.Sites {
		results = append(results, Result{0, site.URL, false})
	}

	resultsChan := make(chan Result)

	for i := range results {
		go callURL(results[i].url, resultsChan)
	}

	err := false
	for i := 0; i < len(results); i++ {
		result := <-resultsChan
		for j := range results {
			if results[j].url == result.url {
				results[j] = result
				if result.httpCode != 200 {
					err = true
				}
			}
		}
	}

	for _, result := range results {
		log.Printf("URL: %s, HTTP Code: %d, Complete: %v\n", result.url, result.httpCode, result.complete)
	}

	if err {
		if err := sendErrsEmail(config, results); err != nil {
			recover()
			log.Println("Could not send errors email", err)
		}
	}
}

func main() {
	log.Println("Loading config")
	config := getConfig(configPath)

	// run once on startup

	if err := sendEmailUpActive(config); err != nil {
		recover()
		log.Println("Could not send initial is active email", err)
	}

	checkSites(config)

	// send periodic "is active" emails
	if len(config.Checks) > 0 {
		go func() {
			for {
				nextCheck := nextConfirmTime(time.Now(), config.Checks)

				fmt.Printf("Waiting %v until the next scheduled is active time\n", nextCheck.Truncate(time.Second))
				time.Sleep(nextCheck)

				if err := sendEmailUpActive(config); err != nil {
					recover()
					log.Println("Could not send is active email", err)
				}
			}
		}()
	}

	// check sites each tick
	tick := int(config.Tick)
	for {
		now := time.Now()

		nextTick := tick - ((now.Minute()*60 + now.Second()) % tick)
		if nextTick == 0 {
			nextTick = tick
		}

		fmt.Printf("Waiting %vs until the next tick\n", nextTick)
		time.Sleep(time.Duration(nextTick) * time.Second)

		go checkSites(config)
	}
}
