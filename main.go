package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/proxy" // Required package for proxy support
)

// Telegram Bot API URL
const telegramAPI = "https://api.telegram.org/bot<BotAPI>/sendMessage"

// Telegram Chat ID (where the bot will send messages)
const chatID = "numbers"

// Target site and timeout threshold
const targetSite = "https://example.com"
const slowResponseThreshold = 100 * time.Millisecond // 100 ms

// Tor proxy address
const torProxyAddress = "127.0.0.1:9150"

// Function to check if Tor is working
func checkTorConnection() error {
	dialer, err := proxy.SOCKS5("tcp", torProxyAddress, nil, proxy.Direct)
	if err != nil {
		return fmt.Errorf("failed to connect to Tor: %w", err)
	}

	// Send an HTTP request through Tor to test the connection
	client := &http.Client{
		Transport: &http.Transport{
			Dial: dialer.Dial,
		},
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get("http://check.torproject.org/")
	if err != nil {
		return fmt.Errorf("failed to reach Tor check site: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Tor check site returned non-OK status: %s", resp.Status)
	}

	log.Println("Successfully connected to Tor")
	return nil
}

// Function to send a message to Telegram via Tor
func sendTelegramMessage(message string) {
	payload := map[string]string{
		"chat_id": chatID,
		"text":    message,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("JSON marshalling failed: %s\n", err)
		return
	}

	// Create HTTP client via Tor proxy
	client, err := getTorHTTPClient()
	if err != nil {
		log.Printf("Failed to create Tor client: %s\n", err)
		return
	}

	// Send a POST request to Telegram
	resp, err := client.Post(telegramAPI, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Failed to send message to Telegram: %s\n", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Message sent to Telegram: %s\n", message)
}

// Function to check if the site is accessible and measure response time via Tor
func checkSiteStatus() {
	start := time.Now()

	// Create HTTP client via Tor proxy
	client, err := getTorHTTPClient()
	if err != nil {
		log.Printf("Failed to create Tor client: %s\n", err)
		return
	}

	// Send GET request to the target site
	resp, err := client.Get(targetSite)
	if err != nil {
		message := fmt.Sprintf("Site down: %s\n", targetSite)
		sendTelegramMessage(message)
		log.Println(message)
		return
	}
	defer resp.Body.Close()

	// Calculate response time
	duration := time.Since(start)

	// Send a warning or status message based on response time
	if duration > slowResponseThreshold {
		message := fmt.Sprintf("Warning: Site %s is slow. Response time: %.2f ms\n", targetSite, duration.Seconds()*1000)
		sendTelegramMessage(message)
	} else {
		message := fmt.Sprintf("Site %s is up. Response time: %.2f ms\n", targetSite, duration.Seconds()*1000)
		sendTelegramMessage(message)
	}

	log.Printf("Site status: %s - Response time: %.2f ms\n", targetSite, duration.Seconds()*1000)
}

// Function to create an HTTP client that routes requests through Tor
func getTorHTTPClient() (*http.Client, error) {
	// Connect to the Tor SOCKS5 proxy
	dialer, err := proxy.SOCKS5("tcp", torProxyAddress, nil, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("failed to create SOCKS5 proxy: %w", err)
	}

	// Configure HTTP transport with custom dialer
	transport := &http.Transport{
		Dial: dialer.Dial,
	}

	// Create an HTTP client that sends requests via the proxy
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second, // Set timeout duration
	}

	return client, nil
}

func main() {
	// Check Tor connection
	if err := checkTorConnection(); err != nil {
		log.Fatalf("Tor connection failed: %s\n", err)
	}

	// Create a ticker to check the site every 5 minutes
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	// First site check
	checkSiteStatus()

	// Continue checking the site every 5 minutes
	for range ticker.C {
		checkSiteStatus()
	}
}
