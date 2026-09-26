# Configra
### Reliable Configuration, Without the Complexity.

**Configra** is a cloud-native configuration management and feature flag service designed for high-velocity engineering teams. It provides immutable versioning, schema validation, and atomic rollbacks to eliminate configuration-induced outages.

Built on strict "Twelve-Factor App" principles, Configra ensures 100% parity between local development and cloud production environments.

---

## The Problem
Configuration errors are a leading cause of production incidents. Existing solutions are often:
1.  **Too Complex**: Requiring heavy enterprise SaaS contracts or complex Kubernetes operators.
2.  **Too Risky**: Allowing "blind" updates without validation or easy rollback paths.
3.  **Too Fragmented**: Splitting secrets, flags, and JSON config across different tools.

## The Solution
Configra treats configuration as a first-class citizen with the same rigor as compiled code.

*   **Immutable History**: Every change creates a new version. You can always see *who* changed *what* and *when*.
*   **Schema Enforcement**: Configurations are validated against strict JSON schemas before they are ever accepted.
*   **Zero-Downtime Rollbacks**: Instantly revert to any previous known-good state via the CLI or API.
*   **Environment Parity**: The exact same Go binary runs locally (Docker Compose) and in production (Google Cloud Run).

---

## Features

| Feature | Description |
| :--- | :--- |
| **Atomic Versioning** | Every update is transactional. No partial states. |
| **Strict Validation** | Type checking (Int, String, Enum, Boolean) prevents bad data entry. |
| **Project Security** | **API Key Authentication** ensures only authorized clients can push configs. |
| **CLI First** | Validate configs locally (`configra validate`) before pushing them. |
| **Auto-Migration** | The service self-manages its database schema on startup. |
| **Cloud Native** | Stateless architecture ready for Serverless (Cloud Run, Render, Fly.io). |

---

## Technology Stack

This project leverages a modern, maintainable stack designed for scale:

*   **Core**: Go (Golang) 1.21+
*   **Database**: PostgreSQL 15 (Managed or Containerized)
*   **Architecture**: REST API with clean "Service-Repository" layering
*   **Infrastructure**: Docker & Terraform (IaC)
*   **CI/CD**: GitHub Actions

---

## Quick Start

### 1. Local Development
Run the complete stack with Docker Compose. This starts the API and a local PostgreSQL instance.

```bash
docker-compose up --build
```

The API will be available at `http://localhost:8080`.

### 2. Using the CLI
Configra comes with a dedicated CLI for local workflows.

```bash
# Install the CLI
go install ./cmd/cli

# 1. Validate a local config file against a schema
configra validate -schema schema.json -config config.json

# 2. Push the validated config to the server
configra push -file config.json -project 1

# 3. Rollback to a previous version (Emergency)
configra rollback -project 1 -key feature_flags -version 1

```

---

## Deployment

Configra is container-first. You can deploy it to any platform that supports Docker.

### Option A: Professional Cloud (GCP)
We provide full Terraform scripts in `deploy/gcp/` to provision:
*   **Google Cloud Run** (Serverless Compute)
*   **Cloud SQL** (Managed Database)

**Setup:**
1.  Navigate to `deploy/gcp`.
2.  Run `terraform apply`.
3.  The CI/CD pipeline (`.github/workflows/deploy.yml`) will automatically build and deploy updates.

### Option B: Zero-Config Cloud (Render / Railway)
For a "Vercel-like" experience for backend containers:

1.  Push this repo to **GitHub**.
2.  Connect it to **Render.com** or **Koyeb**.
3.  Add a PostgreSQL database (e.g., via **Neon.tech** or Render's free tier).
4.  Set the `DB_HOST`, `DB_USER`, etc., environment variables.
5.  Deploy.

*Note: Vercel is not recommended as it does not natively support long-running Docker containers.*

---

## API Reference

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/v1/validate` | Dry-run validation of a config payload. |
| `POST` | `/v1/configs` | Create a new configuration version. |
| `GET` | `/health` | Service health check. |

---

## License

MIT License. Free for commercial and non-commercial use.
