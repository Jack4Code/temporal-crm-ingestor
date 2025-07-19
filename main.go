package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
)

// const (
// 	expectedToken = "supersecrettokengo123"

// 	clientID     = "1000.JZ10KQZK3BE18W80TJ5JYJKL1RH6OE"
// 	clientSecret = "8f6e3c4e438cf7b1666be3ab3ed60faee4b77f3c84"
// 	refreshToken = "1000.20cb118d5491933a35add68aa861e7cc.ac094b63f45b4cfd8e93b4872ae58bca"
// )

type Config struct {
	ExpectedToken string `toml:"expected_token"`
	Zoho          ZohoConfig
}

type ZohoConfig struct {
	ClientID     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
	RefreshToken string `toml:"refresh_token"`
}

var cfg Config

func loadConfig(path string) error {
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}
	return nil
}

type ZohoPayload map[string]interface{}

func getAccessToken() (string, error) {
	url := "https://accounts.zoho.com/oauth/v2/token"
	data := fmt.Sprintf(
		"refresh_token=%s&client_id=%s&client_secret=%s&grant_type=refresh_token",
		cfg.Zoho.RefreshToken, cfg.Zoho.ClientID, cfg.Zoho.ClientSecret,
	)

	resp, err := http.Post(url, "application/x-www-form-urlencoded", bytes.NewBufferString(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	token, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access_token not found in response: %s", body)
	}

	return token, nil
}

func createLead(accessToken string, leadData map[string]interface{}) error {
	url := "https://www.zohoapis.com/crm/v2/Leads"

	payload := map[string]interface{}{
		"data": []map[string]interface{}{leadData},
	}

	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Zoho-oauthtoken "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("📬 Zoho CRM Response:", string(body))

	if resp.StatusCode >= 300 {
		return fmt.Errorf("CRM API error: %s", body)
	}

	return nil
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("X-Auth-Token")
	if authToken != cfg.ExpectedToken {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println("❌ Unauthorized attempt: missing or invalid token")
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var payload ZohoPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Printf("✅ Authenticated webhook received: %+v\n", payload)

	accessToken, err := getAccessToken()
	if err != nil {
		log.Println("❌ Failed to get access token:", err)
		http.Error(w, "Auth error", http.StatusInternalServerError)
		return
	}

	lead := map[string]interface{}{
		"First_Name":  payload["firstName"],
		"Last_Name":   payload["lastName"],
		"Email":       payload["email"],
		"Phone":       payload["phone"],
		"Company":     payload["companyName"],
		"Lead_Source": "Webhook",
	}

	err = createLead(accessToken, lead)
	if err != nil {
		log.Println("❌ Failed to create lead:", err)
		http.Error(w, "CRM error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Lead received and created in CRM")
}

func oauthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code param", http.StatusBadRequest)
		return
	}

	log.Printf("✅ Received OAuth code: %s\n", code)
	fmt.Fprintf(w, "OAuth code received. You can close this window.")
}

// func main() {
// 	http.HandleFunc("/webhook", webhookHandler)
// 	http.HandleFunc("/oauth/callback", oauthCallbackHandler)

// 	fmt.Println("🚀 Listening on http://localhost:8080")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

func main() {
	if err := loadConfig("config.toml"); err != nil {
		log.Fatalf("❌ Config error: %v", err)
	}

	http.HandleFunc("/webhook", webhookHandler)
	http.HandleFunc("/oauth/callback", oauthCallbackHandler)

	fmt.Println("🚀 Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
