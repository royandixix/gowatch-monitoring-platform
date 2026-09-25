# GoWatch

<p align="center">
  <strong>Website & API Monitoring Platform</strong>
</p>

<p align="center">
  Monitor website availability, HTTP status, response time, uptime, and monitoring history from one dashboard.
</p>

---

## About GoWatch

**GoWatch** is a full-stack website and API monitoring platform designed to monitor the availability and performance of web services.

The application periodically checks registered websites, stores monitoring results in PostgreSQL, and presents the information through an interactive dashboard.

GoWatch provides authentication, website monitoring, uptime statistics, response-time analytics, monitoring history, user profile management, and backend security protection.

The project is built using **Go**, **Gin**, **PostgreSQL**, **SvelteKit**, and **TypeScript**.

---

## Main Features

### Authentication

- User registration
- User login
- JWT-based authentication
- Protected application routes
- Current user session
- Logout
- Password hashing using bcrypt

### User Profile

- View profile information
- Update name
- Update email
- Change password
- Current password verification
- Prevent reusing the same password

### Website Monitoring

Users can:

- Add a website monitor
- View all monitors
- View monitor details
- Update a monitor
- Delete a monitor
- Activate or deactivate monitoring
- Configure monitoring interval

Each monitor stores information such as:

- Monitor name
- Website URL
- Monitor type
- Monitoring interval
- Active status
- Creation date
- Update date

### Automatic Monitoring Worker

GoWatch includes a background monitoring worker that automatically checks active monitors according to their configured interval.

Monitoring results include:

- Website status
- HTTP status code
- Response time
- Error message
- Check timestamp

### Dashboard Analytics

The dashboard displays monitoring information such as:

- Total monitors
- Active monitors
- Monitors currently up
- Monitors currently down
- Total checks
- Average uptime
- Average response time
- Recent monitoring activity
- Response-time visualization

### Monitoring History

Monitoring history can be viewed and filtered by:

- Monitor
- Status
- Date
- Pagination

Each monitoring record contains:

- Monitor information
- Website status
- HTTP status
- Response time
- Error information
- Check time

### Monitor Details

Each monitor has a detailed monitoring page containing:

- Current status
- Uptime percentage
- Average response time
- Total checks
- Successful checks
- Failed checks
- Latest monitoring results
- Response-time chart

---

## Security

GoWatch includes several backend protections.

### SSRF Protection

Before a URL can be monitored, GoWatch validates the destination and blocks unsafe network targets.

Examples of blocked destinations include:

- localhost
- loopback addresses
- private IPv4 networks
- private IPv6 networks
- link-local addresses
- metadata addresses
- multicast addresses
- local network hostnames
- unsafe redirect destinations

Only HTTP and HTTPS URLs are accepted.

### DNS Validation

DNS resolution is checked before establishing a connection.

Resolved IP addresses are validated to prevent requests from reaching restricted internal networks.

### Redirect Protection

HTTP redirect destinations are also validated.

The number of redirects is limited to prevent redirect loops and unsafe redirection.

### Request Timeout

Monitoring requests use network timeouts to prevent hanging requests.

### Rate Limiting

Rate limiting is implemented on sensitive endpoints such as:

- Registration
- Login
- Single URL check
- Multiple URL check

### Authentication

Protected endpoints require a valid JWT token.

User-owned resources are queried using the authenticated user ID to prevent users from accessing another user's monitoring data.

---

## Tech Stack

### Backend

- Go
- Gin Web Framework
- pgx PostgreSQL Driver
- JWT
- bcrypt
- Go Context
- Background Worker

### Frontend

- SvelteKit
- Svelte
- TypeScript
- Vite
- Tailwind CSS
- SweetAlert2

### Database

- PostgreSQL

### Development & Production Tools

- Docker
- Docker Compose
- Nginx
- Git
- GitHub

---

## Architecture

```text
                  ┌─────────────────────┐
                  │       Browser       │
                  └──────────┬──────────┘
                             │
                             ▼
                  ┌─────────────────────┐
                  │      SvelteKit      │
                  │      Frontend       │
                  └──────────┬──────────┘
                             │
                             │ REST API
                             ▼
                  ┌─────────────────────┐
                  │       Go API        │
                  │        Gin          │
                  └──────────┬──────────┘
                             │
              ┌──────────────┴──────────────┐
              │                             │
              ▼                             ▼
    ┌──────────────────┐          ┌──────────────────┐
    │    PostgreSQL    │          │ Monitoring Worker │
    │     Database     │          │  Website Checker  │
    └──────────────────┘          └─────────┬────────┘
                                           │
                                           ▼
                                  ┌──────────────────┐
                                  │ External Website │
                                  │     / API        │
                                  └──────────────────┘
```

---

## Production Architecture

GoWatch also includes Docker production configuration.

```text
Browser
   │
   ▼
Nginx Gateway
   │
   ├── /api/* ───────► Go Backend
   │
   └── /* ───────────► SvelteKit Frontend
                           │
Go Backend ────────────────┤
   │
   ▼
PostgreSQL
```

Production deployment can be completed later on a VPS or cloud server.

---

## Project Structure

```text
gowatch/
│
├── backend/
│   │
│   ├── cmd/
│   │   └── api/
│   │       └── main.go
│   │
│   ├── internal/
│   │   ├── auth/
│   │   ├── config/
│   │   ├── handler/
│   │   ├── middleware/
│   │   ├── model/
│   │   ├── monitor/
│   │   ├── repository/
│   │   ├── service/
│   │   └── worker/
│   │
│   ├── migrations/
│   ├── Dockerfile
│   ├── .dockerignore
│   ├── go.mod
│   └── go.sum
│
├── frontend/
│   │
│   ├── src/
│   │   ├── lib/
│   │   │   ├── api/
│   │   │   ├── components/
│   │   │   ├── stores/
│   │   │   ├── types/
│   │   │   └── utils/
│   │   │
│   │   └── routes/
│   │       ├── (app)/
│   │       │   ├── dashboard/
│   │       │   ├── history/
│   │       │   ├── monitors/
│   │       │   └── profile/
│   │       │
│   │       └── (auth)/
│   │           ├── login/
│   │           └── register/
│   │
│   ├── static/
│   ├── Dockerfile
│   ├── .dockerignore
│   ├── package.json
│   └── svelte.config.js
│
├── deploy/
│   └── nginx/
│       └── default.conf
│
├── docker-compose.yml
├── docker-compose.prod.yml
├── .env.production.example
├── .gitignore
└── README.md
```

---

## Requirements

For local development, install:

- Go
- Node.js
- npm
- Docker Desktop
- Git

Recommended versions used during development:

```text
Go      : 1.27+
Node.js : 24+
npm     : 11+
Docker  : Docker Desktop
```

---

# Local Development

## 1. Clone Repository

```bash
git clone https://github.com/royandixix/gowatch-monitoring-platform.git

cd gowatch-monitoring-platform
```

---

## 2. Start PostgreSQL

GoWatch provides a Docker Compose configuration for PostgreSQL.

```bash
docker compose up -d
```

Check the container:

```bash
docker compose ps
```

Development PostgreSQL is available through:

```text
localhost:5433
```

---

## 3. Backend Configuration

Enter the backend directory:

```bash
cd backend
```

Create the environment file:

```bash
cp .env.example .env
```

Configure the backend environment.

Example:

```env
APP_NAME=GoWatch
APP_ENV=development
APP_PORT=8081

DB_HOST=localhost
DB_PORT=5433
DB_USER=gowatch_user
DB_PASSWORD=YOUR_DATABASE_PASSWORD
DB_NAME=gowatch
DB_SSLMODE=disable

JWT_SECRET=YOUR_RANDOM_JWT_SECRET

FRONTEND_URL=http://localhost:5173
```

Do not commit the real `.env` file.

---

## 4. Run Backend

From:

```text
gowatch/backend
```

run:

```bash
go run ./cmd/api
```

The backend API will run on:

```text
http://localhost:8081
```

Health check:

```text
http://localhost:8081/api/v1/health
```

---

## 5. Frontend Setup

Open another terminal:

```bash
cd frontend
```

Install dependencies:

```bash
npm install
```

Run development server:

```bash
npm run dev
```

Frontend:

```text
http://localhost:5173
```

---

# Application URLs

Development:

```text
Frontend
http://localhost:5173

Backend
http://localhost:8081

Health API
http://localhost:8081/api/v1/health

PostgreSQL Host Port
localhost:5433
```

---

# API Endpoints

## Public

```text
GET    /
GET    /api/v1/health

POST   /api/v1/auth/register
POST   /api/v1/auth/login
```

## Authentication

```text
GET    /api/v1/auth/me
```

## Dashboard

```text
GET    /api/v1/dashboard
```

## Website Check

```text
GET    /api/v1/check
GET    /api/v1/check-multiple
```

## Monitor Management

```text
POST   /api/v1/monitors
GET    /api/v1/monitors
GET    /api/v1/monitors/:id
PUT    /api/v1/monitors/:id
DELETE /api/v1/monitors/:id
```

## Monitoring Results

```text
GET    /api/v1/monitors/:id/results
GET    /api/v1/monitors/:id/stats
GET    /api/v1/history
```

## Profile

```text
PUT    /api/v1/profile
PUT    /api/v1/profile/password
```

Protected routes require:

```http
Authorization: Bearer <JWT_TOKEN>
```

---

# Backend Testing

Enter the backend:

```bash
cd backend
```

Run all tests:

```bash
go test ./...
```

Verbose tests:

```bash
go test -v ./...
```

Race detector:

```bash
go test -race ./...
```

Static analysis:

```bash
go vet ./...
```

Build verification:

```bash
go build ./...
```

The automated test suite covers areas including:

- JWT token handling
- Authentication middleware
- Rate limiting
- URL validation
- SSRF protection
- Monitor service validation
- Website checker
- API handlers

---

# Frontend Validation

Enter:

```bash
cd frontend
```

Type and Svelte validation:

```bash
npm run check
```

Lint:

```bash
npm run lint
```

Production build:

```bash
npm run build
```

The frontend production build uses:

```text
@sveltejs/adapter-node
```

---

# Docker Production

GoWatch contains a production Docker architecture consisting of:

```text
PostgreSQL
Go Backend
SvelteKit Frontend
Nginx Gateway
```

Create production environment:

```bash
cp .env.production.example .env.production
```

Generate a secure PostgreSQL password:

```bash
openssl rand -hex 24
```

Generate a secure JWT secret:

```bash
openssl rand -hex 32
```

Edit:

```text
.env.production
```

Never commit `.env.production`.

Validate Docker Compose:

```bash
docker compose \
  --env-file .env.production \
  -f docker-compose.prod.yml \
  config --quiet
```

Build:

```bash
docker compose \
  --env-file .env.production \
  -f docker-compose.prod.yml \
  build
```

Start:

```bash
docker compose \
  --env-file .env.production \
  -f docker-compose.prod.yml \
  up -d
```

Check services:

```bash
docker compose \
  --env-file .env.production \
  -f docker-compose.prod.yml \
  ps
```

Production gateway:

```text
http://localhost:8080
```

Stop services without deleting the database volume:

```bash
docker compose \
  --env-file .env.production \
  -f docker-compose.prod.yml \
  down
```

Do not use `down -v` unless you intentionally want to remove the database volume.

---

# Monitoring Flow

When an active monitor is created, the monitoring process works approximately as follows:

```text
User creates monitor
        │
        ▼
Monitor stored in PostgreSQL
        │
        ▼
Background Worker
        │
        ▼
Find monitors due for checking
        │
        ▼
Validate destination URL
        │
        ▼
HTTP / HTTPS request
        │
        ▼
Collect response information
        │
        ├── Status
        ├── HTTP status code
        ├── Response time
        └── Error
        │
        ▼
Save monitor result
        │
        ▼
Dashboard / History / Statistics
```

---

# Monitor Status

A successful website response is stored as an `UP` monitoring result.

Failed checks can produce a `DOWN` result with error information.

Response time is recorded in milliseconds.

---

# Uptime Calculation

GoWatch currently calculates uptime from the ratio between successful checks and total monitoring checks.

Conceptually:

```text
Uptime (%) =
Successful Checks
----------------- × 100
Total Checks
```

The current uptime metric is therefore check-based rather than a time-weighted SLA calculation.

---

# Database

The main data includes:

```text
users
monitors
monitor_results
```

`users` stores registered users.

`monitors` stores website monitoring configuration.

`monitor_results` stores the result of each monitoring execution.

---

# Development Commands

Backend:

```bash
cd backend

go run ./cmd/api
```

Frontend:

```bash
cd frontend

npm run dev
```

Database:

```bash
docker compose up -d
```

Check Docker:

```bash
docker compose ps
```

---

# Production Readiness

The project already includes:

```text
Authentication
User profile management
Website monitoring
Background monitoring worker
Monitoring history
Monitor statistics
Dashboard analytics
SSRF protection
DNS validation
Redirect protection
Rate limiting
Graceful shutdown
Automated backend tests
Frontend linting and validation
Docker configuration
Nginx gateway configuration
```

Deployment to a public VPS or cloud server can be completed as a separate deployment step.

---

# Future Development

Possible future improvements include:

- Email notifications
- Telegram notifications
- Discord notifications
- Slack notifications
- SSL certificate monitoring
- Domain expiration monitoring
- API keyword validation
- Custom HTTP headers
- HTTP method configuration
- Incident management
- Public status pages
- Organization/team accounts
- Multiple monitoring regions
- Advanced uptime SLA calculation
- Prometheus metrics
- Grafana integration
- CI/CD pipeline
- Cloud deployment

---

# Security Notes

Never commit:

```text
.env
.env.production
database passwords
JWT secrets
access tokens
API keys
```

Use strong random secrets for production.

---

# Disclaimer

GoWatch is currently developed as a software engineering and portfolio project.

Before running it in a critical production environment, perform additional security review, infrastructure hardening, backup configuration, monitoring, and deployment testing.

---

# Author

**Royandi**

Software Developer

Primary areas:

```text
Web Development
Backend Development
Frontend Development
Software Engineering
```

---

# License

This project is intended for educational, portfolio, and development purposes.

A formal open-source license can be added if the project is intended for public reuse.

---

<p align="center">
  Made with Go, SvelteKit, PostgreSQL, and Docker.
</p># gowatch-monitoring-platform
