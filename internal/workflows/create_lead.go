package workflows

import (
	"context"

	"temporal-crm-ingestor/internal/crm"
	"temporal-crm-ingestor/internal/utils"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// CreateLeadWorkflow is the Temporal workflow entry point
func CreateLeadWorkflow(ctx workflow.Context, payload map[string]interface{}) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("🚀 Starting CreateLeadWorkflow", "payload", payload)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: utils.DefaultActivityTimeout,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    utils.InitialRetryInterval,
			BackoffCoefficient: 2.0,
			MaximumInterval:    utils.MaxRetryInterval,
			MaximumAttempts:    5,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var contactID, dealID string

	// Step 1: Create Contact
	err := workflow.ExecuteActivity(ctx, CreateContactActivity, payload).Get(ctx, &contactID)
	if err != nil {
		logger.Error("❌ Failed to create contact", "error", err)
		return err
	}

	// Step 2: Create Deal
	dealInput := map[string]interface{}{
		"contact_id": contactID,
		"payload":    payload,
	}
	err = workflow.ExecuteActivity(ctx, CreateDealActivity, dealInput).Get(ctx, &dealID)
	if err != nil {
		logger.Error("❌ Failed to create deal, compensating by deleting contact", "error", err)
		_ = workflow.ExecuteActivity(ctx, DeleteContactActivity, contactID).Get(ctx, nil)
		return err
	}

	logger.Info("✅ Workflow completed successfully", "contactID", contactID, "dealID", dealID)
	return nil
}

// CreateContactActivity creates a contact in Zoho CRM
func CreateContactActivity(ctx context.Context, payload map[string]interface{}) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("📬 Executing CreateContactActivity", "payload", payload)

	contact := map[string]interface{}{
		"First_Name": payload["firstName"],
		"Last_Name":  payload["lastName"],
		"Email":      payload["email"],
		"Phone":      payload["phone"],
	}

	return crm.CreateContactWithRefresh(ctx, contact)
}

// CreateDealActivity creates a deal linked to a contact
func CreateDealActivity(ctx context.Context, input map[string]interface{}) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("💼 Executing CreateDealActivity", "input", input)

	deal := map[string]interface{}{
		"Contact_Name": input["contact_id"],
		"Deal_Name":    "Sample Deal", // or extract from input["payload"]
		"Stage":        "Qualification",
	}

	return crm.CreateDealWithRefresh(ctx, deal)
}

// DeleteContactActivity deletes a contact in case of rollback
func DeleteContactActivity(ctx context.Context, contactID string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("🧹 Executing DeleteContactActivity", "contactID", contactID)

	return crm.DeleteContact(ctx, contactID)
}
