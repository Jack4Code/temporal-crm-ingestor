package crm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// createLeadRaw sends a lead to Zoho using a pre-obtained access token
func createLeadRaw(accessToken string, leadData map[string]interface{}) error {
	url := "https://www.zohoapis.com/crm/v2/Leads"

	payload := map[string]interface{}{
		"data": []map[string]interface{}{leadData},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Zoho-oauthtoken "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return fmt.Errorf("CRM API error: %s", body)
	}

	return nil
}
