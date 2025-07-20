package crm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// createContact creates a contact and returns the ID
func createContact(accessToken string, contactData map[string]interface{}) (string, error) {
	url := "https://www.zohoapis.com/crm/v2/Contacts"

	payload := map[string]interface{}{
		"data": []map[string]interface{}{contactData},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Zoho-oauthtoken "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("CRM API error: %s", body)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	data := result["data"].([]interface{})
	first := data[0].(map[string]interface{})
	id, ok := first["details"].(map[string]interface{})["id"].(string)
	if !ok {
		return "", fmt.Errorf("ID not found in response: %s", body)
	}

	return id, nil
}

// createDeal creates a deal and returns the ID
func createDeal(accessToken string, dealData map[string]interface{}) (string, error) {
	url := "https://www.zohoapis.com/crm/v2/Deals"

	payload := map[string]interface{}{
		"data": []map[string]interface{}{dealData},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Zoho-oauthtoken "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("CRM API error: %s", body)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	data := result["data"].([]interface{})
	first := data[0].(map[string]interface{})
	id, ok := first["details"].(map[string]interface{})["id"].(string)
	if !ok {
		return "", fmt.Errorf("ID not found in response: %s", body)
	}

	return id, nil
}

// deleteContactByID deletes a contact by its ID
func deleteContactByID(accessToken string, contactID string) error {
	url := fmt.Sprintf("https://www.zohoapis.com/crm/v2/Contacts/%s", contactID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Zoho-oauthtoken "+accessToken)

	resp, err := http.DefaultClient.Do(req)
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
