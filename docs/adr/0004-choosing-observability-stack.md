# ADR-0004: Choosing observability stack (Metrics, Logs, Monitoring)

* **Status:** accepted
* **Deciders:** Yehor Reil (Lead Software Engineer)
* **Date:** 2026-04-06

---

## Context and Problem Statement

FinScale is a distributed microservices-basead system where failures, performance issues, and unexpected behavior may occur across multiple services.

To ensure system reliability and maintainability, it is necessary to implement a comprehensive observability solution that provides visibility into:
- system health,
- application performance,
- logs for debugging,
- real-time monitoring and alerting.

The challenge is to select an observability stack that is cost-effective, scalable, and well-suited for cloud-native environments.

---

## Decision Drivers

* **Visibility:** Ability to monitor system behavior in real time.
* **Debugging Efficiency:** Fast access to logs and metrics for issue investigation.
* **Scalability:** Support for increasing number of services and data volume.
* **Cost Efficiency:** Minimize infrastructure and licensing costs.
* **Integration:** Compatibility with Kubernetes / Docker environments.
* **Ecosystem Maturity:** Availability of community support and integrations.

---

## Considered Options

* **ELK Stack (Elasticsearch, Logstash, Kibana)**
* **Grafana Stack (Prometheus, Loki, Grafana, Promtail)**
* **Cloud-native solutions (AWS CloudWatch, Datadog)**

---

## Decision Outcome

Chosen option: **Grafana Stack (Prometheus + Loki + Grafana + Promtail)**

The Grafana ecosystem was selected as it provides a unified, cost-effective, and cloud-native observability solution.

This stack enables:
- metrics collection via Prometheus,
- log aggregation via Loki,
- log shipping via Promtail,
- visualization and alerting via Grafana.

It offers tight integration between logs and metrics while maintaining lower infrastructure costs compared to ELK and commercial SaaS solutions.

---

## Consequences

### Positive

* Unified observability stack with tight integration between components.
* Lower resource consumption compared to ELK stack.
* Cost-effective (open-source, no vendor lock-in).
* Powerful visualization and alerting via Grafana.
* Scales well in containerized environments.

---

### Negative

* Loki is less feature-rich for log search compared to Elasticsearch.
* Requires manual setup and maintenance.
* Alerting and configuration can be complex.
* No built-in advanced APM compared to Datadog or similar tools.
* Requires additional setup for distributed tracing (e.g., Tempo).

---

## Scope

This decision applies to all backend microservices and infrastructure components.

All services must:
- expose metrics in Prometheus format,
- write structured logs (JSON),
- integrate with centralized logging via Promtail,
- be included in Grafana dashboards and alerting rules.

---

## Pros and Cons of the Options

### ELK Stack

* **Pros:**
    - Powerful log search and analytics
    - Mature ecosystem

* **Cons:**
    - High resource consumption
    - Expensive to operate at scale
    - More complex setup

---

### Grafana Stack (Prometheus + Loki)

* **Pros:**
    - Lightweight and cost-efficient
    - Native integration between metrics and logs
    - Designed for cloud-native environments

* **Cons:**
    - Less powerful log querying than ELK
    - Requires assembling multiple components

---

### Cloud-native Solutions (AWS CloudWatch, Datadog)

* **Pros:**
    - Fully managed
    - Easy to set up
    - Advanced features (APM, tracing)

* **Cons:**
    - Vendor lock-in
    - High cost at scale
    - Less flexibility