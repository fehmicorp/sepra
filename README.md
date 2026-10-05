# Sepra

### Secured Endpoint Protection Responsive Agent

Sepra is an endpoint security and management platform designed to provide centralized control, endpoint visibility, policy enforcement, telemetry, software/file distribution, and security response across managed devices.

> **Core principle:** Cloud/CMC decides. DSM distributes. Agent enforces.

Sepra is being designed as a modular platform supporting:

- Windows
- Linux
- macOS
- iOS
- Android

The project is currently being developed as a **solo-developer project**, so the architecture and roadmap intentionally prioritize small, testable increments over premature complexity.

---

## Project Status

> 🚧 **Active Development — Early Architecture / MVP**

Sepra is not yet production-ready.

The current development priority is to establish a stable core platform before implementing advanced EDR/XDR capabilities.

### Current priorities

- Secure endpoint identity
- Endpoint enrollment
- Endpoint inventory
- Agent heartbeat
- Policy management
- Remote commands
- Audit logging
- File/software distribution
- DSM
- Telemetry
- Alerts
- Basic incident management
- CMC

Advanced detection and response capabilities will come later.

---

# Architecture

Sepra is divided into four major components:

```text
                         SEPRA
                           │
          ┌────────────────┼────────────────┐
          │                │                │
          ▼                ▼                ▼
    Cloud Console         CMC             DSM
     / Main Server       Management      Distribution
          │                │                │
          └────────────────┼────────────────┘
                           │
                           ▼
                    Endpoint Agent
                           │
             ┌─────────────┼─────────────┐
             │             │             │
          Windows        Linux         macOS
             │
          iOS / Android
```

### Component responsibilities

| Component | Responsibility |
|---|---|
| **Cloud Console** | Central control plane |
| **CMC** | Customer management console |
| **DSM** | LAN relay, cache, and distribution |
| **Endpoint Agent** | Local enforcement, telemetry, commands |

---

# Cloud Console

The Cloud Console is the central Sepra control plane.

It is responsible for:

- Organizations
- Users
- Authentication
- RBAC
- Endpoint enrollment
- Endpoint inventory
- Device groups
- Policies
- Commands
- Telemetry
- Alerts
- Incidents
- Distribution
- DSM management
- CMC management
- Audit logs
- Reporting

## Technology

```text
Backend       → Go
Frontend      → Next.js / React
Database      → PostgreSQL
Cache         → Redis
Object Store  → S3-compatible storage
Messaging     → Message bus / worker system
Deployment    → Docker
```

The first implementation will use a **modular monolith**.

Microservices should only be introduced when there is a clear operational or scalability reason.

---

# Central Management Console

The **Central Management Console (CMC)** is the customer-side management layer.

CMC may be:

- Installed on Windows
- Installed on Linux
- Hosted by Sepra

CMC is responsible for customer-side operational management and should remain synchronized with the central control plane.

CMC is **not intended to become a second independent source of truth**.

---

# Distribution Service Manager

The **Distribution Service Manager (DSM)** provides local LAN distribution and caching.

Its primary purpose is to reduce WAN traffic and improve endpoint management reliability.

```text
                    SEPRA CLOUD
                         │
                         │
                    Package/File
                         │
                         ▼
                        DSM
                  ┌──────┼──────┐
                  │      │      │
                  ▼      ▼      ▼
                Agent  Agent  Agent
```

For example, if 100 endpoints require the same 500 MB package:

Without DSM:

```text
Cloud → 100 × 500 MB
       = 50 GB WAN traffic
```

With DSM:

```text
Cloud → DSM
       = 500 MB

DSM → LAN endpoints
       = Local network traffic
```

DSM therefore acts as:

- Local cache
- Distribution point
- LAN relay
- Optional endpoint communication checkpoint

DSM must have its own authenticated identity.

Endpoints must never blindly trust a device simply because it is on the same LAN.

---

# Endpoint Agent

The Endpoint Agent is the enforcement component installed on managed devices.

```text
Endpoint Agent
│
├── Identity
├── Secure Communication
├── Policy Engine
├── Telemetry
├── Command Engine
├── Update Engine
├── Local State
└── Platform Modules
    ├── Windows
    ├── Linux
    ├── macOS
    ├── iOS
    └── Android
```

The agent should share:

- Protocol
- Identity model
- Policy model
- Event schema
- Command model

while keeping security enforcement platform-specific.

---

# Connectivity Model

Sepra supports both LAN-connected and roaming endpoints.

### Endpoint with DSM available

```text
Agent
  │
  ▼
 DSM
  │
  ▼
Cloud
```

### Roaming endpoint

```text
Agent
  │
  ▼
Cloud
```

The agent should automatically use DSM when appropriate and fall back to cloud connectivity when DSM is unavailable.

DSM must therefore **never become a mandatory single point of failure**.

---

# Data Architecture

Sepra intentionally separates persistent state, cache, artifacts, and asynchronous processing.

```text
                 SEPRA DATA LAYER

        ┌──────────────┐
        │ PostgreSQL   │
        │              │
        │ Source of    │
        │ Truth        │
        └──────┬───────┘
               │
     ┌─────────┼──────────┐
     │         │          │
     ▼         ▼          ▼
  Redis    Object Store  Message Bus
  Cache      Files       Events/Jobs
```

## PostgreSQL

PostgreSQL is the primary system of record.

Expected domains include:

```text
organizations
users
roles
permissions
devices
device_groups
device_credentials
policies
policy_versions
policy_assignments
alerts
incidents
commands
command_results
software
software_versions
distribution_jobs
dsms
cmcs
audit_logs
licenses
integrations
```

## Redis

Redis is used for:

- Cache
- Rate limiting
- Distributed locks
- Short-lived state
- Job coordination
- Real-time state

Redis is **not** the permanent source of truth.

## Object Storage

Object storage is used for:

- Agent installers
- Software packages
- Configuration files
- Security artifacts
- Distributed files
- Update packages

S3-compatible storage should be supported.

For self-hosted deployments, MinIO or another S3-compatible implementation may be used.

## Message Bus

Asynchronous work should use a proper message/event mechanism.

Potential implementations include:

- NATS
- RabbitMQ
- Kafka
- Redis Streams

The exact implementation may change as the project evolves.

MongoDB should not be used simply as a generic replacement for both object storage and messaging.

---

# Security Model

Security is a foundational requirement of Sepra.

The platform follows these principles:

### Zero Trust

LAN presence does not imply trust.

### Least Privilege

Components receive only the permissions they require.

### Strong Device Identity

Endpoints should have unique cryptographic identities.

### Secure Transport

Network communication should use TLS 1.3 or stronger where supported.

### mTLS

mTLS should be used where strong infrastructure/device authentication is required.

### Signed Updates

Agent updates must be cryptographically verified before installation.

### Artifact Integrity

Packages should be verified using checksums and signatures where applicable.

### RBAC

Administrative operations must be authorized using role-based permissions.

### Auditability

Security-sensitive operations must generate audit records.

### Offline Safety

Temporary loss of connectivity must not disable the endpoint's last valid security configuration.

---

# Repository Structure

The repository is expected to evolve toward:

```text
sepra/
│
├── server/
│   ├── cmd/
│   │   ├── api/
│   │   ├── worker/
│   │   └── migrate/
│   │
│   ├── internal/
│   │   ├── auth/
│   │   ├── tenant/
│   │   ├── organization/
│   │   ├── device/
│   │   ├── policy/
│   │   ├── command/
│   │   ├── telemetry/
│   │   ├── detection/
│   │   ├── alert/
│   │   ├── incident/
│   │   ├── distribution/
│   │   ├── dsm/
│   │   ├── cmc/
│   │   ├── audit/
│   │   ├── reporting/
│   │   └── update/
│   │
│   ├── migrations/
│   └── deployments/
│
├── console/
│
├── agent/
│
├── dsm/
│
├── cmc/
│
├── protocol/
│
├── docs/
│
└── docker-compose.yml
```

The structure is intentionally modular without requiring separate microservices.

---

# API

Sepra APIs should be versioned.

Example:

```text
/api/v1/
```

Major API areas:

```text
/api/v1/auth
/api/v1/admin
/api/v1/agent
/api/v1/dsm
/api/v1/organizations
/api/v1/devices
/api/v1/groups
/api/v1/policies
/api/v1/alerts
/api/v1/incidents
/api/v1/commands
/api/v1/distribution
/api/v1/software
/api/v1/files
/api/v1/audit
```

Agent-facing and administrator-facing APIs should remain logically separated.

---

# Event Model

Sepra events should use a common envelope.

```text
event_id
tenant_id
device_id
event_type
timestamp
schema_version
severity
source
correlation_id
payload
```

This allows telemetry, alerts, commands, and future event-processing infrastructure to evolve without constantly changing the entire platform.

---

# Endpoint Enrollment

The basic enrollment flow is:

```text
Administrator
      │
      ▼
Enrollment Token / Policy
      │
      ▼
Install Agent
      │
      ▼
Agent → Cloud / DSM
      │
      ▼
Identity Verification
      │
      ▼
Device Identity
      │
      ▼
Device Registered
      │
      ▼
Policy Assigned
      │
      ▼
Initial Inventory
      │
      ▼
Heartbeat
```

---

# Policy Model

Policies follow a hierarchical model:

```text
Global Policy
      │
      ▼
Organization Policy
      │
      ▼
Group Policy
      │
      ▼
Device Policy
```

Policies should be:

- Versioned
- Scoped
- Signed where appropriate
- Assignable
- Auditable
- Deterministic when conflicts occur

---

# Remote Commands

The command lifecycle should follow:

```text
Requested
    │
    ▼
Authorized
    │
    ▼
Queued
    │
    ▼
Delivered
    │
    ▼
Executed
    │
    ▼
Result
    │
    ▼
Audited
```

Initial commands may include:

- Refresh policy
- Collect inventory
- Restart agent
- Restart device
- Shutdown device
- Start scan
- Stop process
- Quarantine file
- Restore file
- Network isolation
- Network release
- Install package
- Uninstall package

Command availability depends on platform capabilities.

---

# Development Philosophy

Sepra is currently being developed by a single developer with limited daily development time.

The project therefore follows a **small vertical slice** development model.

The goal is not:

> Build everything as quickly as possible.

The goal is:

> Build one small capability, connect it end-to-end, test it, and keep it working.

---

# Weekly Development Cycle

Each development week follows:

```text
Monday
  ↓
Plan / Design

Tuesday
  ↓
Implementation

Wednesday
  ↓
Implementation

Thursday
  ↓
Integration

Friday
  ↓
Finish / Fix

Saturday
  ↓
Optional Polish / Documentation

Sunday
  ↓
TEST ONLY
```

## Sunday Rule

**Sunday is reserved for testing.**

No major feature development should be performed on Sunday.

Sunday should be used for:

- Integration testing
- Regression testing
- Failure testing
- Restart testing
- Network-loss testing
- Authentication testing
- Permission testing
- Data consistency testing
- Documentation
- Bug recording

If a feature fails Sunday testing, it is **not considered complete**.

---

# Development Roadmap

The initial roadmap is intentionally conservative.

## Phase 1 — Foundation

**Weeks 1–4**

### Week 1
Project foundation.

- Go server
- Next.js console
- Docker Compose
- Configuration
- Logging
- Health checks

### Week 2
Database foundation.

- PostgreSQL
- Migrations
- Tenant model
- Organization model

### Week 3
Authentication.

- Login
- Sessions
- Initial RBAC

### Week 4
Endpoint identity.

- Enrollment token
- Device registration
- Device identity

### Phase 1 result

A user can authenticate and register a test endpoint.

---

# Phase 2 — Core Endpoint Management

**Weeks 5–8**

### Week 5

Agent heartbeat and endpoint status.

### Week 6

Endpoint inventory.

### Week 7

Device groups.

### Week 8

Basic policy delivery.

### Phase 2 result

Sepra can manage a small collection of endpoints.

---

# Phase 3 — Policy & Control

**Weeks 9–12**

### Week 9

Real policy enforcement.

### Week 10

Remote commands.

### Week 11

Audit logging.

### Week 12

MVP checkpoint and cleanup.

### Phase 3 result

Sepra has a basic usable control plane.

---

# Phase 4 — Distribution & DSM

**Weeks 13–18**

### Week 13

Object storage and artifacts.

### Week 14

Distribution jobs.

### Week 15

DSM registration and cache.

### Week 16

DSM-based distribution.

### Week 17

Roaming endpoint routing.

### Week 18

Agent update mechanism.

### Phase 4 result

Sepra can securely distribute files/packages through LAN infrastructure.

---

# Phase 5 — Security Telemetry

**Weeks 19–21**

### Week 19

Telemetry pipeline.

### Week 20

Alerts.

### Week 21

Basic incidents.

### Phase 5 result

Sepra begins behaving like a security monitoring platform rather than only an endpoint management system.

---

# Phase 6 — CMC & Hardening

**Weeks 22–24**

### Week 22

CMC v1.

### Week 23

Operational hardening.

- Metrics
- Health checks
- Backup/recovery
- Failure handling
- Security review

### Week 24

Release candidate.

- Bug fixing
- Documentation
- Demo environment
- Full regression testing

### Phase 6 result

**Sepra Core MVP Release Candidate**

---

# What MVP Means

The first Sepra MVP is **not** intended to compete feature-for-feature with mature EDR platforms.

The MVP should demonstrate that the architecture works end-to-end.

The MVP should be able to:

- Authenticate administrators
- Enroll an endpoint
- Establish endpoint identity
- Display endpoint inventory
- Track endpoint health
- Create groups
- Assign policies
- Enforce basic policies
- Send remote commands
- Record audit events
- Store artifacts
- Distribute files/packages
- Use DSM for LAN distribution
- Fall back to cloud connectivity
- Update the agent
- Receive telemetry
- Generate alerts
- Create basic incidents
- Connect a basic CMC

---

# Platform Strategy

Do **not** implement all platforms simultaneously.

Recommended order:

```text
1. One primary desktop/server platform
             ↓
2. Stabilize common agent protocol
             ↓
3. Second desktop/server platform
             ↓
4. Remaining desktop/server platforms
             ↓
5. iOS
             ↓
6. Android
```

The first platform should become the reference implementation for the common agent architecture.

iOS and Android should be treated as platform-specific implementations sharing the Sepra control plane rather than forcing desktop security assumptions onto mobile operating systems.

---

# Testing Strategy

Testing should happen continuously, with a dedicated Sunday stability cycle.

## Every Feature

At minimum test:

### Happy Path

Does the feature work normally?

### Restart

What happens after restarting the relevant service?

### Network Failure

What happens when connectivity disappears?

### Invalid Input

What happens with malformed or unauthorized data?

### Duplicate Operation

What happens when the same command/request is submitted twice?

### Permission

Can an unauthorized user perform the operation?

### Recovery

Does the system recover after the failure?

### Regression

Do previous features still work?

---

# Example Local Environment

The initial development environment should be simple:

```text
Docker Compose
│
├── sepra-api
├── sepra-console
├── postgres
├── redis
└── object-storage
```

Later:

```text
Docker / Kubernetes
│
├── API
├── Workers
├── Console
├── PostgreSQL
├── Redis
├── Object Storage
└── Message Bus
```

Do not introduce Kubernetes just because Sepra may eventually need it.

---

# Local Development

## Requirements

Expected development tools:

- Git
- Docker
- Docker Compose
- Go
- Node.js
- npm/pnpm
- PostgreSQL tooling
- A supported endpoint test environment

---

## Start the Development Environment

Once the initial Compose configuration is available:

```bash
docker compose up -d
```

Check running services:

```bash
docker compose ps
```

View logs:

```bash
docker compose logs -f
```

Stop the environment:

```bash
docker compose down
```

---

# Development Commands

Expected commands as the project matures:

```bash
# Start infrastructure
docker compose up -d

# Run backend
go run ./server/cmd/api

# Run frontend
npm run dev

# Run Go tests
go test ./...

# Run frontend tests
npm test

# Run database migrations
go run ./server/cmd/migrate

# Format Go
gofmt -w .

# Build
go build ./...
```

Exact commands may change as the repository structure evolves.

---

# Branching Strategy

Keep Git workflow simple.

Recommended:

```text
main
 │
 ├── feature/endpoint-enrollment
 ├── feature/policy-engine
 ├── feature/dsm-cache
 └── fix/heartbeat-timeout
```

### `main`

Should remain reasonably stable.

### Feature branches

Used for individual weekly capabilities.

Avoid maintaining a large number of long-lived branches.

---

# Commit Strategy

Prefer small meaningful commits.

Examples:

```text
feat(agent): add device enrollment
feat(api): add endpoint heartbeat
feat(policy): add policy versioning
feat(dsm): add artifact cache
fix(agent): handle offline policy sync
test(agent): add heartbeat reconnect tests
docs: document local development
```

---

# Security Reporting

Do not publish sensitive security vulnerabilities as normal GitHub issues.

If this repository becomes public, provide a dedicated security reporting process, such as:

```text
SECURITY.md
```

Security reports should contain enough information to reproduce the issue without exposing unnecessary sensitive information.

---

# Documentation

Documentation should eventually be organized as:

```text
docs/
├── architecture/
├── api/
├── agent/
├── dsm/
├── cmc/
├── deployment/
├── security/
├── development/
├── testing/
└── decisions/
```

Architecture Decision Records can be maintained under:

```text
docs/decisions/
```

Example:

```text
ADR-001-postgresql-source-of-truth.md
ADR-002-modular-monolith.md
ADR-003-dsm-routing.md
ADR-004-agent-identity.md
```

---

# Design Principles

Sepra follows these principles:

1. **Security first**
2. **Simple before distributed**
3. **Modular before microservices**
4. **Strong identity**
5. **Least privilege**
6. **Offline resilience**
7. **Explicit trust**
8. **Auditable actions**
9. **Versioned protocols**
10. **Test before expanding scope**

---

# Long-Term Vision

After the Core MVP is stable, Sepra can evolve into a broader endpoint security platform.

Potential future capabilities:

```text
Core Endpoint Management
        │
        ▼
Security Telemetry
        │
        ▼
Detection
        │
        ▼
EDR
        │
        ▼
Automated Response
        │
        ▼
Threat Intelligence
        │
        ▼
XDR / Security Operations
```

Potential future modules:

- Advanced EDR
- Behavioral detection
- Threat hunting
- Vulnerability management
- Application control
- Device control
- DLP
- Threat intelligence
- SIEM integration
- SOAR integration
- Identity integrations
- Compliance reporting
- AI-assisted investigation

These capabilities should be added only after the core architecture is stable.

---

# Current Development Rule

> **Do not build the entire security platform at once.**

For every weekly milestone:

```text
Build
  ↓
Connect
  ↓
Run
  ↓
Break
  ↓
Fix
  ↓
Test Sunday
  ↓
Document
  ↓
Next milestone
```

A feature is complete only when it survives its Sunday testing cycle.

---

# License

License: **TBD**

The project licensing model should be decided before public distribution.

---

# Disclaimer

Sepra is an actively developed security platform.

Until production readiness is explicitly declared, it should be considered experimental software and should not be deployed into production environments where failure could create unacceptable security or operational risk.

---

## Sepra

**Secured Endpoint Protection Responsive Agent**

> Cloud/CMC decides.  
> DSM distributes.  
> Agent enforces.