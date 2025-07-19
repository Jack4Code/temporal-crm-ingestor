# 🧠 temporal-crm-ingestor

A reliable, fault-tolerant ingestion pipeline for syncing Zoho Form or Survey submissions to Zoho CRM using [Temporal](https://temporal.io/) for orchestration.

> Designed to solve challenges in multi-record creation, API delay handling, and record linking across CRM entities.

---

## 🔧 What It Does

- Accepts webhooks from Zoho Forms or Surveys  
- Parses and transforms submissions into CRM-friendly formats  
- Uses Temporal to orchestrate the creation of:
  - Uses Temporal to orchestrate the creation of one or more CRM records per submission  
  - Designed to support complex workflows involving multiple related entities (e.g., Leads, Contacts, Accounts)  
  - Prevents partial writes by ensuring all steps succeed or the entire workflow retries  
- Handles failures, retries, and queue latency gracefully  
- Uses TOML for config management (credentials, tokens, etc.)

---

## 📦 Tech Stack

- **Go** for webhook server & workers  
- **Temporal** (self-hosted or Cloud) for workflow orchestration  
- **Zoho CRM** (REST API) for record creation  
- **Ngrok** for local testing (optional)  
- **TOML** config files (with `.gitignore` protections)

---

## 🚀 Quick Start

```bash
# Clone the repo
git clone https://github.com/your-username/temporal-crm-ingestor.git
cd temporal-crm-ingestor

# Setup configuration
cp config.example.toml config.toml
# Fill in your Zoho tokens and expected webhook token

# Start Temporal worker and server (in separate terminals)
go run cmd/server/main.go
go run cmd/worker/main.go
```


## 🛠 Features

- ✅ Secure webhook with static token header  
- ✅ OAuth2-based Zoho token refresh logic  
- ✅ Upsert & deduplication logic for CRM records  
- ✅ Graceful handling of queue latency / Zoho API lag  
- ✅ Extensible: easily add new form field mappings or CRM modules

---

## 📊 Roadmap

- [ ] Add support for survey branching logic  
- [ ] CLI for token refresh + test submission  
- [ ] Form validation + required field enforcement  
- [ ] Temporal Web UI integration  
- [ ] Terraform module for deploy infra (optional)

---

## 📁 Project Structure

```bash
.
├── cmd/
│   ├── server/         # Webhook HTTP server (receives Zoho Forms submissions)
│   └── worker/         # Temporal worker for processing workflows
├── config/
│   └── config.toml     # Contains sensitive credentials (in .gitignore)
├── internal/
│   ├── crm/            # Handles all Zoho CRM API interactions
│   ├── workflows/      # Temporal workflows and activity definitions
│   └── utils/          # Utility functions (e.g. JSON parsing, logging)
├── go.mod              # Go module definition
└── README.md
```


---

## 🔒 Security

- 🔐 All sensitive configuration (OAuth credentials, tokens, secrets) is stored in `config.toml`, which is included in `.gitignore`
- 🔑 Webhook endpoint validates an `X-Auth-Token` header for basic security
- 🔄 Zoho access tokens are refreshed using a stored refresh token — no manual intervention needed
- ✅ No secrets are stored in source code or committed to version control

---

## 🧪 Contributions

Contributions are welcome!

If you'd like to:
- Report a bug
- Suggest a feature
- Submit a pull request

Please open an issue or create a PR with clear documentation and use case. The project is designed to be extensible to other CRMs and forms platforms in the future.

---

## 📊 Observability & Success Metrics

A production-ready deployment should include:

- 🔁 Retry and error handling with visibility via Temporal’s built-in UI
- 📜 Structured logs for webhook receipt, token exchange, and CRM API calls
- 📈 Optional metrics integration with tools like Prometheus + Grafana or DataDog for:
  - Workflow execution success/failure rates
  - CRM API latency and errors
  - Duplicate or missing field detection
- ✅ Alerting thresholds configurable for retry exhaustion, API rate limits, etc.


---

## 📜 License

MIT License.  
Use it, fork it, ship it — just don’t hardcode your Zoho secrets in public repos 😄


