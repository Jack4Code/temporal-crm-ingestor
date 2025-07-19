# 🧠 temporal-crm-ingestor

A reliable, fault-tolerant ingestion pipeline for syncing Zoho Form or Survey submissions to Zoho CRM using [Temporal](https://temporal.io/) for orchestration.

> Designed to solve challenges in multi-record creation, API delay handling, and record linking across CRM entities.

---

## 🔧 What It Does

- Accepts webhooks from Zoho Forms or Surveys  
- Parses and transforms submissions into CRM-friendly formats  
- Uses Temporal to orchestrate the creation of:
  - Members  
  - Caregivers  
  - Emergency contacts  
  - Staff or partners (customizable)  
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

```text
.
├── cmd/
│   ├── server/         # Webhook server
│   └── worker/         # Temporal workflows & activities
├── config/
│   └── config.toml     # Auth, tokens, secrets
├── internal/
│   ├── crm/            # Zoho API logic
│   ├── workflows/      # Temporal workflows
│   └── utils/          # Helpers
├── README.md
└── go.mod
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

## 📈 Suggested KPIs

- 🧮 **Error rate**: < 1% on CRM record creation  
- ⏱ **Average queue delay**: < 2 minutes per submission  
- 🎯 **Match/link accuracy**: > 98% between submission and CRM contact  

Log events, retry metrics, and Temporal UI dashboards are recommended for real-time observability.

---

## 📜 License

MIT License.  
Use it, fork it, ship it — just don’t hardcode your Zoho secrets in public repos 😄


