# Configra

## Project name

Configra

### What this is

Configra is a lightweight configuration and feature flag service designed for small teams and solo developers. It focuses on safety, validation, versioning, and rollback to prevent configuration related outages.

This project demonstrates backend systems design, cloud deployment, and operational thinking.

---

### Why this exists

Configuration errors are one of the most common causes of production incidents. Many existing solutions are either too heavy or too abstract for small teams.

This project explores a simpler approach:
- clear APIs
- explicit versioning
- local validation before deployment
- safe rollback paths

---

### Core features

- User authentication
- Project and environment management
- Configuration or feature flag CRUD
- Schema validation before activation
- Version history for every change
- Rollback to previous versions
- Rate limiting and audit logging
- CLI integration for local validation and deployment

---

### Tech stack

**Backend**
- Go
- Chi HTTP router
- Standard library where possible

**Data**
- PostgreSQL

**Auth**
- JWT based authentication

**Infrastructure**
- Docker
- GCP Cloud Run (Serverless)
- GCP Cloud SQL (PostgreSQL)
- Terraform (IaC)
- GitHub Actions (CI/CD)


**Observability**
- Structured logging
- Health check endpoint

**CLI**
- Go based CLI
- Local config validation
- Remote push and fetch

---

### Architecture overview

The system is composed of three parts.

1. An HTTP API that handles authentication, configuration management, and versioning.
2. A PostgreSQL database that stores users, projects, environments, and config versions.
3. A CLI tool that validates configuration locally and interacts with the API.

The design favors clarity over abstraction and avoids unnecessary complexity.

---

### API overview

- Authentication endpoints for login and registration.
- Project and environment management endpoints.
- Configuration or feature flag endpoints with versioning support.
- Rollback endpoints for safe recovery.
- Health endpoint for monitoring.

---

### Validation and safety model

- All configurations are validated before becoming active.
- Invalid configurations are rejected early.
- Each change creates a new immutable version.
- Rollbacks restore the last known good state.

This reduces the blast radius of mistakes.

---

### Failure scenarios considered

- Invalid configuration submission
- Concurrent updates
- Database unavailability
- Authentication failure
- Rollback failure

Each scenario is documented with expected behavior.

---

### Deployment

The API is containerized using Docker.
Secrets are injected using environment variables.
The service is deployed to a managed cloud platform.

The deployment prioritizes simplicity and reproducibility.

---

### CLI usage

- Local validation of configuration files.
- Push validated configurations to the server.
- Fetch active configurations.
- Trigger rollbacks when necessary.

---

### Trade-offs

- No Kubernetes to reduce operational overhead.
- Explicit versioning instead of automatic merging.
- PostgreSQL chosen for clarity over distributed stores.

These decisions favor maintainability and learning value.

---

### What this project demonstrates

- Backend API design
- Cloud deployment and Infrastructure Awareness
- **Twelve-Factor App principles** (Dev/Prod parity)
- Failure handling
- Documentation and communication

---

### Future improvements

- Role based access control
- Better metrics and alerting
- Configuration diffing
- Multi region support

---

### Deployment Infrastructure
This project includes production-ready Infrastructure as Code (IaC) and CI/CD pipelines.

**Database Management**
The API is configured to **automatically apply database migrations on startup**. This means you do not need to manually connect to Cloud SQL to create tables; the application handles its own schema consistency.

**Infrastructure as Code (Terraform)**
Located in `deploy/gcp/`, this Terraform configuration provisions:
1.  **GCP Cloud Run**: Serverless container hosting for the API.
2.  **Cloud SQL (PostgreSQL)**: Managed database with high availability.
3.  **IAM Roles**: Secure access policies.

**CI/CD (GitHub Actions)**
The `.github/workflows/deploy.yml` pipeline handles:
1.  Building the Docker container.
2.  Pushing to Google Artifact Registry.
3.  Deploying the new revision to Cloud Run.

### How to use
**Local Development**
1.  `go run ./cmd/api` - Starts the server
2.  `go run ./cmd/cli validate -schema schema.json -config config.json` - Validates config

**Deploying to Cloud (First Time Setup)**
1.  **Prerequisites**:
    - A Google Cloud Platform (GCP) Account.
    - [GCloud CLI](https://cloud.google.com/sdk/docs/install) installed and authenticated (`gcloud auth login`).
    - [Terraform](https://developer.hashicorp.com/terraform/downloads) installed.

2.  **Initial Provisioning**:
    - Create a new Project in GCP Console.
    - Enable the "Cloud Run Admin" and "Cloud SQL Admin" APIs.
    - Navigate to `deploy/gcp` in your terminal.
    - Create a `sequel.auto.tfvars` file with your secrets:
      ```hcl
      project_id  = "your-project-id"
      db_password = "secure-password-here"
      container_image = "us-docker.pkg.dev/cloudrun/container/hello" # Placeholder for first run
      ```
    - Run `terraform init` to download plugins.
    - Run `terraform apply` to create the database and service.

3.  **Connect GitHub**:
    - Create a Service Account key in GCP with permissions to push to Artifact Registry and deploy to Cloud Run.
    - Add this JSON key as `GCP_CREDENTIALS` in your GitHub Repository Secrets.
    - Push your code to `main` to trigger the actual deployment!



```
configra/
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── cli/
│       └── main.go
│
├── internal/
│   ├── auth/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── middleware.go
│   │
│   ├── projects/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repository.go
│   │
│   ├── configs/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── validator.go
│   │   └── repository.go
│   │
│   ├── audit/
│   │   └── logger.go
│   │
│   ├── middleware/
│   │   ├── auth.go
│   │   └── ratelimit.go
│   │
│   ├── db/
│   │   ├── db.go
│   │   └── migrations/
│   │
│   └── config/
│       └── env.go
│
├── pkg/
│   └── utils/
│       └── response.go
│
├── scripts/
│   ├── migrate.sh
│   └── seed.sh
│
├── deploy/
│   └── gcp/
│       ├── main.tf
│       └── variables.tf
│
├── .github/
│   └── workflows/
│       └── deploy.yml
│
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```
