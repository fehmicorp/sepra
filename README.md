# Sepra Infrastructure Platform (`fehmicorp/sepra`)

[![API Version](https://img.shields.io/badge/api-v1-blue.svg)](#)
[![License](https://img.shields.io/badge/license-Proprietary-red.svg)](#)

**Sepra** is a multi-tenant, distributed edge management platform built for secure bi-directional telemetry, dynamic remote execution, and hybrid deployment models for Enterprise and SME clients.

---

## 🏗 Deployment Topologies & Network Routing

Sepra supports two distinct deployment pipelines tailored to customer infrastructure needs:

### Model A: Enterprise Client-Side Architecture (CMD + Local DSM)
For enterprise clients requiring local management autonomy. Analytics and license validation stream upstream to Central Cloud, while primary control sits locally via CMD.


┌──────────────┐     Commands / Config Push     ┌──────────────┐     Task Dispatch / Deploy     ┌──────────────┐
│              ├───────────────────────────────►│              ├───────────────────────────────►│              │
│   cloud/v1   │                                │    dsm/v1    │                                │   agent/v1   │
│ (Controller) │                                │ (Data/Bridge)│                                │ (Edge Node)  │
│              │◄───────────────────────────────┤              │◄───────────────────────────────┤              │
└──────────────┘      Aggregated Telemetry      └──────────────┘     Heartbeats & Task Results  └──────────────┘


┌─────────────────────────┐          Licensing & Analytics Sync          ┌─────────────────────────┐          Commands & Dispatch           ┌─────────────────────────┐
│     Central Cloud       │◄────────────────────────────────────────────►│           CMD           │◄──────────────────────────────────────►│         DSM/v1          │
│ (sepra.fehmicorp.in/v1) │                                              │ (cmd.clientdomain.com)  │                                        │(dsm.clientdomain.com/v1)│
└─────────────────────────┘                                              └─────────────────────────┘                                        └────────────┬────────────┘
│
┌───────────────────────────────────────────────────────────────────────────────┘
│ Heartbeats, Results & Config Sync
▼
┌───────────────────────┐
│       agent/v1        │
│  (Local/Roaming Host) │
└───────────────────────┘