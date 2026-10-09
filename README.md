# Configra
### Versioned, validated config with one-click rollback.

**Configra** stores versioned application configuration, validates updates against a schema, and lets teams roll back through the API or dashboard.

---

## The Problem
Configuration errors are a leading cause of production incidents. Existing solutions are often:
1.  **Too Complex**: Requiring heavy enterprise SaaS contracts or complex Kubernetes operators.
2.  **Too Risky**: Allowing "blind" updates without validation or easy rollback paths.

## The Solution
Configra treats configuration as a first-class citizen with the same rigor as compiled code.

*   **Immutable History**: Every change creates a new version with its timestamp and the project owner or project credential label.
*   **Schema Enforcement**: Configurations are validated against strict JSON schemas before they are ever accepted.
*   **Zero-Downtime Rollbacks**: Instantly revert to any previous known-good state via the CLI or API.

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

The API will be available at `http://localhost:8080`. Browse the OpenAPI document at `http://localhost:8080/docs`.

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

Set `CONFIGRA_API_KEY` before running `push`, `fetch`, or `rollback`. Push sends the project credential in the `X-API-Key` header. The `-project` flag is retained for CLI compatibility; authorization scope is derived from the API key.

## Verification and release evidence

CI runs `go test ./...` with a PostgreSQL 15 service. The tests exercise configuration validation, append-only history, historical reads, diffs, rollback, concurrent version allocation, migration idempotency, and API-key acceptance/rejection. Tests needing PostgreSQL skip locally unless `CONFIGRA_TEST_DATABASE_URL` is set; CI sets it and provisions PostgreSQL.

CI also builds the production Docker image, starts it against PostgreSQL, checks liveness/readiness and the landing/dashboard routes, confirms protected access is denied without credentials, and exercises authenticated config write/read/rollback plus CLI fetch. A green run is evidence for that commit and those checks; it is not a substitute for a production deployment or load/security review.

The `/health` endpoint reports process liveness. `/ready` reports readiness only while PostgreSQL is reachable. The API exits on startup when it cannot connect to PostgreSQL or apply migrations.

For each release, retain the GitHub Actions run URL and commit SHA. After deployment, retain the Cloud Run deployment run, the image digest, and the post-deployment `/health` and `/ready` results. A successful workflow file or local test run alone does not prove a cloud deployment succeeded. Schema validation uses Configra's documented `version`/`rules` format; it is not a general JSON Schema implementation.

### Evidence matrix

| Claim | Automated evidence | Remaining release evidence |
| --- | --- | --- |
| Atomic versioning, history, rollback | PostgreSQL lifecycle and concurrent-writer integration tests | CI run on the release commit |
| Strict config validation | Validator unit tests and invalid-push integration test | Confirm schemas match the documented Configra rules format |
| Migrations | Fresh migration and rerun/idempotency integration test | Successful deployment migration log |
| CLI push/fetch/rollback | CLI fetch and API rollback smoke operations run against the image in CI; CLI credentials use `CONFIGRA_API_KEY` | Release artifact/installation check on supported platforms |
| Dashboard and landing page | HTTP smoke checks in container CI | Visual review if UI changes |
| Container operation | Production image build and PostgreSQL-backed route smoke tests | Published image digest for each release |
| Cloud deployment | Workflow deploys then checks `/health` and `/ready` | Green deployment job, URL, and image digest |
| Twelve-Factor practices | External config, stdout logs, stateless API, health/readiness; DB required at startup | Operational review of secrets, scaling, backups, and managed DB settings |

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

The OpenAPI 3.1 document is served at `/docs` and available in [openapi.yaml](openapi.yaml). Protected routes require `X-API-Key`.

---

## License

MIT License. Free for commercial and non-commercial use.
