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

	return workflow.ExecuteActivity(ctx, CreateLeadActivity, payload).Get(ctx, nil)
}

// CreateLeadActivity actually calls the Zoho CRM API
func CreateLeadActivity(ctx context.Context, payload map[string]interface{}) error {
	logger := activity.GetLogger(ctx)
	logger.Info("📬 Executing CreateLeadActivity", "payload", payload)

	lead := map[string]interface{}{
		"First_Name":  payload["firstName"],
		"Last_Name":   payload["lastName"],
		"Email":       payload["email"],
		"Phone":       payload["phone"],
		"Company":     payload["companyName"],
		"Lead_Source": "Webhook",
	}

	return crm.CreateLeadWithRefresh(ctx, lead)
}
