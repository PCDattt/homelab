# Homelab: Kubernetes, GitOps & Observability

A self-hosted Kubernetes homelab running on two repurposed Ubuntu laptops. The project is a practical environment for learning Linux infrastructure, Kubernetes operations, GitOps delivery, secrets management, and observability.

The cluster is bootstrapped with kubeadm, uses containerd as its runtime, and connects its nodes over Tailscale. Workloads and infrastructure configuration are maintained in Git and reconciled by Argo CD.

The observability stack covers both metrics and logs:

- **VictoriaLogs Collector** collects Kubernetes logs and forwards them to VictoriaLogs.
- **VictoriaLogs** stores and queries log data.
- **Node Exporter** and **Scaphandre** expose node-level resource metrics.
- **Prometheus** collects infrastructure and application metrics.
- **Grafana** provides dashboards and exploration.

## Goals

This homelab is an ongoing environment for developing practical skills in:

- Kubernetes internals and cluster operations.
- Linux systems, networking, and resource management.
- GitOps, CI/CD, and platform engineering.
- Secrets management and workload security.
- Metrics, centralized logging, and observability.
- Failure analysis, reliability, and operational recovery.

The long-term goal is to build a reproducible, reliable, and observable platform while developing the skills and experience needed to become an Infrastructure Platform Engineer.

## Architecture

```text
                         Developer
                             |
                             v
                      GitHub Repository
                             |
                             v
                       GitHub Actions
                  +-----------------------+
                  | Lint / validation     |
                  | Gitleaks              |
                  | Trivy                 |
                  | Build container image|
                  | Push image to GHCR    |
                  +-----------------------+
                             |
                             v
                    GitOps configuration
                             |
                             v
                          Argo CD
                             |
                             v
                    Kubernetes cluster
                  +-----------------------+
                  |                       |
             Laptop node 1           Laptop node 2
             Ubuntu 26.04            Ubuntu 26.04
             12 GB RAM               8 GB RAM
                  |                       |
                  +-----------+-----------+
                              |
                          Tailscale
                              |
          +-------------------+--------------------+
          |                   |                    |
          v                   v                    v
     Secrets              Metrics                Logs
   +-----------+       +-------------+      +-------------+
   |  Vault    |       | Prometheus  |      |     VLC     |
   +-----------+       +-------------+      | Log collector|
          |                   |             +------+------+
          v                   v                    |
   +-----------+          Grafana                 v
   |    ESO    |                         +------------------+
   +-----------+                         |  VictoriaLogs    |
                                         | storage & query  |
                                         +------------------+

       Node-level metrics: Node Exporter, Scaphandre
```

## Hardware and infrastructure

| Component | Configuration |
|---|---|
| Nodes | 2 old laptops |
| Operating system | Ubuntu 26.04 |
| CPU | 8 CPU cores per node |
| Memory | 12 GB on one node; 8 GB on the other |
| Kubernetes | kubeadm |
| Container runtime | containerd |
| Node connectivity | Tailscale |
| Source control | GitHub |
| Container registry | GitHub Container Registry (GHCR) |
| Deployment model | Argo CD GitOps |

The control-plane node is also allowed to run workloads, making better use of the limited resources available in this two-node lab.

## Technology stack

### Kubernetes and platform

| Technology | Role |
|---|---|
| Kubernetes | Container orchestration and workload management |
| kubeadm | Cluster bootstrap and lifecycle management |
| Ubuntu | Host operating system |
| containerd | Kubernetes container runtime |
| Tailscale | Private connectivity between the nodes |
| Argo CD | GitOps reconciliation and application deployment |

### CI/CD and supply-chain checks

| Technology | Role |
|---|---|
| GitHub Actions | Automated validation and image build workflow |
| golangci-lint / Go tooling | Go source-code quality checks |
| Gitleaks | Secret detection |
| Trivy | Container and dependency vulnerability scanning |
| GHCR | Container image registry |

### CI/CD and GitOps workflow

The delivery workflow separates building an application image from deploying it to the cluster.

```text
Code change
    |
    v
GitHub Actions
    |
    +--> Lint and validate
    |
    +--> Scan for secrets (Gitleaks)
    |
    +--> Scan for vulnerabilities (Trivy)
    |
    +--> Build container image
    |
    +--> Push image to GHCR
              |
              v
      GitOps configuration
              |
              v
           Argo CD
              |
              v
       Kubernetes cluster
```

The intended responsibilities are:

1. **Validate** source code and configuration before publishing.
2. **Scan** for accidentally committed secrets and known vulnerabilities.
3. **Build** the application container image.
4. **Publish** the image to GHCR.
5. **Reconcile** the desired Kubernetes state through Argo CD.

Git remains the source of truth for declarative deployment configuration. Argo CD continuously compares the live cluster state with the desired state stored in the repository.

### Secrets management

| Technology | Role |
|---|---|
| HashiCorp Vault | Central secrets backend |
| External Secrets Operator (ESO) | Synchronizes external secrets into Kubernetes Secrets |
| Vault unseal Job | Automates the configured Vault unseal workflow |

Secrets are kept out of application manifests and source control. The security of automated unsealing depends on how unseal material and the identity authorized to use it are protected.

### Secrets management workflow

Vault and ESO provide a workflow for consuming external secrets without committing plaintext credentials to Git.

```text
HashiCorp Vault
       |
       | authenticated secret retrieval
       v
External Secrets Operator
       |
       | reconcile
       v
Kubernetes Secret
       |
       v
Application Pod
```

The Vault unseal Job automates the configured unseal procedure. This is an operational convenience, not a replacement for a secure key-management design. Unseal material should not be exposed in Git, container images, ordinary ConfigMaps, or broadly accessible logs.

### Observability

| Technology | Role |
|---|---|
| Prometheus | Metrics collection and time-series storage |
| Grafana | Dashboards and metrics/log exploration |
| VictoriaLogs | Centralized log storage and querying |
| VictoriaLogs Collector | Kubernetes log collection and forwarding |
| Node Exporter | Node metrics |
| Scaphandre | Node-level power and energy metrics |

### Log collection with VictoriaLogs and VLC

The logging pipeline is designed around a collector/storage separation:

```text
Kubernetes workloads
        |
        | container stdout / stderr
        v
VictoriaLogs Collector
        |
        | collected log records
        v
VictoriaLogs
        |
        +----> LogsQL querying via VictoriaLogs APIs
```

**VictoriaLogs** is the log database: it receives, stores, indexes, and serves queries over ingested log records.

**VictoriaLogs Collector** is the collection layer: it runs in the Kubernetes environment, reads workload logs, enriches them with Kubernetes metadata, and ships them to the VictoriaLogs ingestion endpoint.

Keeping collection separate from storage makes the responsibilities easier to reason about:

- VLC is responsible for getting logs out of workloads and delivering them.
- VictoriaLogs is responsible for durable log storage and query execution.
- Grafana can provide a common interface for metrics and logs when the corresponding data source is configured.

The collector is deployed using the VictoriaMetrics Helm chart repository. The repository's deployment configuration has used the victoria-logs-collector chart.

### Metrics collection with Node Exporter, Scaphandre, and Prometheus

The metrics pipeline collects host-level resource metrics and makes them available for monitoring and analysis.

```text
              Kubernetes Nodes
                     |
          +----------+----------+
          |                     |
          v                     v
     Node Exporter          Scaphandre
          |                     |
          |                     |
          +----------+----------+
                     |
                     v
                 Prometheus
                     |
                     | scrape and store metrics
                     v
                  Grafana
                     |
                     v
             Dashboards & Alerts
```

**Node Exporter** exposes hardware and operating-system metrics from Linux hosts, including CPU, memory, filesystem, disk I/O, and network statistics.

**Scaphandre** collects power and energy consumption metrics, allowing the monitoring stack to observe host-level and process-level energy usage.

**Prometheus** periodically scrapes metrics from Node Exporter, Scaphandre, and other configured targets. It stores the collected time-series data and provides PromQL for querying and alerting.

**Grafana** visualizes the collected metrics through dashboards and supports monitoring infrastructure health and resource utilization.

## Possible next steps

- Add actionable Prometheus alert rules and link each alert to a runbook.
- Define and test VictoriaLogs retention, disk-capacity limits, and backup/recovery procedures.
- Verify that VLC checkpoints and buffering behave correctly across Pod restarts.
- Test log delivery during VictoriaLogs unavailability and collector restarts.
- Evaluate persistent-volume behavior and the impact of disk I/O contention on stateful workloads.
- Practice node failure, workload rescheduling, and recovery from a damaged or unavailable data directory.
- Automate TLS certificate lifecycle management.
- Explore Kubernetes Gateway API and network-policy enforcement.
- Document Vault recovery and etcd snapshot restoration procedures.
