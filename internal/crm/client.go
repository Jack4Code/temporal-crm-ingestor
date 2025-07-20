package crm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"temporal-crm-ingestor/config"
)

// -------------------- AUTH --------------------

func GetAccessToken(ctx context.Context) (string, error) {
	data := fmt.Sprintf(
		"refresh_token=%s&client_id=%s&client_secret=%s&grant_type=refresh_token",
		config.Cfg.Zoho.RefreshToken, config.Cfg.Zoho.ClientID, config.Cfg.Zoho.ClientSecret,
	)

	resp, err := http.Post(
		"https://accounts.zoho.com/oauth/v2/token",
		"application/x-www-form-urlencoded",
		strings.NewReader(data),
	)
	if err != nil {
		return "", fmt.Errorf("failed to request token: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	token, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access_token not found in response: %s", body)
	}

	return token, nil
}

// -------------------- PUBLIC "WITH REFRESH" API --------------------

func CreateContactWithRefresh(ctx context.Context, contactData map[string]interface{}) (string, error) {
	token, err := GetAccessToken(ctx)
	if err != nil {
		return "", fmt.Errorf("unable to get access token: %w", err)
	}
	return createContact(token, contactData)
}

func CreateDealWithRefresh(ctx context.Context, dealData map[string]interface{}) (string, error) {
	token, err := GetAccessToken(ctx)
	if err != nil {
		return "", fmt.Errorf("unable to get access token: %w", err)
	}
	return createDeal(token, dealData)
}

func DeleteContactWithRefresh(ctx context.Context, contactID string) error {
	token, err := GetAccessToken(ctx)
	if err != nil {
		return fmt.Errorf("unable to get access token: %w", err)
	}
	return deleteContactRaw(token, contactID)
}

// -------------------- PUBLIC DIRECT TOKEN API --------------------

func CreateLead(token string, leadData map[string]interface{}) (string, error) {
	return createLeadRaw(token, leadData)
}

func DeleteLead(token, leadID string) error {
	return deleteLeadRaw(token, leadID)
}

func CreateDeal(token string, dealData map[string]interface{}) (string, error) {
	return createDealRaw(token, dealData)
}

func DeleteDeal(token, dealID string) error {
	return deleteDealRaw(token, dealID)
}

func DeleteContact(token, contactID string) error {
	return deleteContactRaw(token, contactID)
}
