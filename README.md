# Homelab — Kubernetes, GitOps & Observability

A bare-metal Kubernetes homelab built on two old laptops running Ubuntu 26.04. The project focuses on learning and operating cloud-native infrastructure, automating application delivery, managing secrets securely, and building an observable Kubernetes environment.

The infrastructure runs on physical hardware, with Tailscale providing connectivity between nodes. Kubernetes is bootstrapped using `kubeadm`, while Argo CD manages application deployments through GitOps.

## Architecture

```text
                  Developer
                      |
                      v
                GitHub Repository
                      |
                      v
              GitHub Actions CI/CD
              +------------------+
              | Lint             |
              | Gitleaks         |
              | Trivy Scan       |
              | Build Image      |
              | Push to GHCR     |
              +------------------+
                      |
                      v
              GitOps Configuration
                      |
                      v
                  Argo CD
                      |
                      v
             Kubernetes Cluster
              +---------------------------+
              |                           |
              v                           v
         Laptop Node 1               Laptop Node 2
         Ubuntu 26.04                Ubuntu 26.04
         12 GB RAM                   8 GB RAM
              |                           |
              +------------+--------------+
                           |
                      Tailscale
                           |
              +------------+--------------+
              |            |              |
              v            v              v
           Vault        Prometheus      Grafana
              |
              v
      External Secrets Operator

      Node-level monitoring:
      - Disk Usage Exporter
      - Scaphandre
```

*Note: The diagram illustrates the logical components and relationships, not the exact network topology or pod placement.*

## Infrastructure

| Component | Description |
|---|---|
| Physical hardware | Two repurposed laptops |
| Operating system | Ubuntu 26.04 |
| CPU | 8 CPU cores per node |
| Memory | 12 GB RAM on one node; 8 GB RAM on the other |
| Kubernetes | kubeadm-based cluster |
| Node connectivity | Tailscale |
| Container runtime | Containerd |
| Source control | GitHub |
| Container registry | GitHub Container Registry (GHCR) |
| Deployment model | GitOps with Argo CD |

The cluster runs on physical machines rather than managed Kubernetes services, providing hands-on experience with node provisioning, networking, resource constraints, and infrastructure troubleshooting.

## Technology Stack

### Kubernetes and Infrastructure

- **Kubernetes / kubeadm** — Cluster bootstrap and workload orchestration.
- **Ubuntu Server** — Bare-metal operating system.
- **Tailscale** — Connectivity between machines over a private network.
- **Argo CD** — Declarative application deployment and GitOps reconciliation.

### CI/CD and Container Security

- **GitHub Actions** — Automated CI/CD workflows.
- **Go tooling** — Linting and building Go applications where applicable.
- **Gitleaks** — Detects hardcoded secrets and credentials in source code and Git history.
- **Trivy** — Vulnerability scanning.
- **GitHub Container Registry (GHCR)** — Container image storage.

### Secrets Management

- **HashiCorp Vault** — Centralized secrets management.
- **External Secrets Operator (ESO)** — Synchronizes secrets from external secret-management systems into Kubernetes.
- **Automated Vault unseal job** — Automates the Vault unseal workflow.

### Observability and Monitoring

- **Prometheus** — Metrics collection and time-series monitoring.
- **Grafana** — Metrics visualization and dashboards.
- **Disk Usage Exporter** — Exposes filesystem and disk-usage metrics.
- **Scaphandre** — Monitors power consumption and energy-related metrics at the node level.

## CI/CD Pipeline

The CI/CD pipeline automates the process of validating, scanning, building, and publishing container images.

```text
       Code Changes
            |
            v
       GitHub Actions
            |
            v
       Lint Source Code
            |
            v
       Gitleaks Secret Scan
            |
            v
       Trivy Security Scan
            |
            v
       Build Container Image
            |
            v
       Push Image to GHCR
            |
            v
       Update GitOps Configuration
            |
            v
       Argo CD Reconciliation
            |
            v
       Kubernetes Deployment
```

Pipeline responsibilities:

1. **Lint:** Detect source-code quality issues before publishing an image.
2. **Security scanning:** Use Trivy to identify known vulnerabilities.
3. **Build:** Package the application into a container image.
4. **Publish:** Push the image to GitHub Container Registry.
5. **Deploy:** Use Argo CD to reconcile Kubernetes resources with the desired state stored in Git.

The GitOps configuration is the source of truth for application deployment. This separates application build and image publication from the process of deploying workloads into Kubernetes.

## Secrets Management

Vault and External Secrets Operator provide a foundation for managing application secrets without making plaintext credentials part of the application deployment workflow.

The current setup includes:

- Vault as the central secrets backend.
- External Secrets Operator for synchronizing secrets into Kubernetes.
- An automated job for Vault unseal operations.

This setup provides practical experience with secret synchronization, Kubernetes service accounts and access control, and the operational requirements of running a secrets-management system.

**Security note:** Automated unsealing must be designed carefully. The security of the unseal mechanism depends on where unseal material is stored, how it is accessed, and which identities are authorized to use it. Automation should not eliminate Vault's security boundary.

## Observability

Prometheus collects metrics from the homelab, while Grafana provides dashboards for monitoring system behavior.

The current monitoring components include:

| Component | Purpose |
|---|---|
| Prometheus | Metrics collection and storage |
| Grafana | Dashboards and visualization |
| Disk Usage Exporter | Filesystem and disk-usage visibility |
| Scaphandre | Node-level energy and power-consumption monitoring |

These components help investigate resource consumption and understand how infrastructure workloads behave on resource-constrained hardware.

The monitoring environment also serves as a foundation for future work on alerting, centralized logging, service discovery, and Kubernetes observability.

## Design Principles

- **Infrastructure as code:** Keep infrastructure and deployment configuration version-controlled where practical.
- **GitOps first:** Use Argo CD to manage application deployment and continuously reconcile the desired state.
- **Automated delivery:** Run validation and security scanning before publishing container images.
- **Secrets outside source control:** Use Vault and External Secrets Operator for secret management.
- **Observable systems:** Collect metrics to understand resource usage and diagnose operational problems.
- **Bare-metal learning:** Operate Kubernetes directly on physical machines and learn from real infrastructure constraints.
- **Incremental improvement:** Introduce new components through small, testable changes rather than deploying an entire platform at once.

## Current Status

The following components are already running in the homelab:

- [x] Two Ubuntu 26.04 bare-metal nodes
- [x] Kubernetes cluster bootstrapped with kubeadm
- [x] Tailscale connectivity
- [x] GitHub Actions CI/CD pipeline
- [x] Source-code linting
- [x] Gitleaks secret scanning
- [x] Trivy security scanning
- [x] Container image build and publication to GHCR
- [x] Argo CD GitOps deployment
- [x] HashiCorp Vault
- [x] External Secrets Operator
- [x] Automated Vault unseal job
- [x] Prometheus
- [x] Disk Usage Exporter
- [x] Scaphandre
- [x] Grafana

## Future Improvements

Planned improvements may include:

- [ ] Add centralized logging with Grafana Loki and Grafana Alloy.
- [ ] Improve Prometheus alerting with Prometheus Operator, `PrometheusRule`, and actionable runbooks.
- [ ] Introduce durable, cluster-aware storage with Longhorn and evaluate its failure modes.
- [ ] Practice node failure, workload recovery, storage recovery, and Kubernetes upgrade scenarios.
- [ ] Add TLS certificate automation with cert-manager.
- [ ] Explore Gateway API for Kubernetes ingress and traffic management.
- [ ] Improve observability dashboards and resource-capacity monitoring.
- [ ] Document disaster recovery procedures, including Vault recovery and etcd snapshot restoration.
- [ ] Evaluate high-availability trade-offs within the constraints of two physical nodes.

## Goals

This project is an ongoing practical environment for exploring:

- Kubernetes internals and cluster operations.
- Cloud-native application delivery and GitOps.
- Infrastructure automation and platform engineering.
- Secrets management and security practices.
- Metrics collection, observability, and troubleshooting.
- Linux systems, networking, and resource management.
- Stateful workloads, persistent storage, and failure recovery.

The long-term goal is to build a reliable, reproducible, and observable platform while developing the systems knowledge required to contribute to cloud-native infrastructure and open-source projects.
