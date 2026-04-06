# ADR-0002: Choosing the architectural style (Monolith vs Microservices)

* **Status:** accepted
* **Deciders:** Yehor Reil (Lead Software Engineer)
* **Date:** 2026-04-06

---

## Context and Problem Statement

The FinScale system is designed as a financial platform that includes multiple domains such as authentication, transaction processing, budgeting, and analytics.

The system must:
- process financial transactions reliably and in near real-time,
- scale under increasing load,
- be maintainable as the number of features and services grows,
- support independent development and deployment of different system parts.

The key architectural decision is whether to implement the system as a modular monolith or adopt a microservices architecture from the start.

---

## Decision Drivers

* **Scalability:** Ability to scale independently under uneven load (e.g., transactions vs analytics).
* **Domain Isolation:** Separation of bounded contexts (Auth, Transactions, Analytics, Budgets).
* **Team Scalability:** Ability to develop features in parallel without conflicts.
* **Operational Complexity:** Infrastructure, deployment, and monitoring overhead.
* **Time to Market:** Speed of initial development and iteration.
* **Fault Isolation:** Preventing failures in one domain from affecting others.

---

## Considered Options

* **Modular Monolith**
* **Microservices Architecture**

---

## Decision Outcome

Chosen option: **Microservices Architecture**

Microservices were selected due to the need for clear domain separation, independent scalability, and long-term maintainability of a complex financial system.

This approach enables:
- independent deployment of services,
- isolation of critical components (e.g., transaction processing),
- flexible scaling strategies based on load patterns,
- better alignment with domain-driven design principles.

---

## Consequences

### Positive

* Services can be scaled independently (e.g., Transactions under heavy load).
* Fault isolation reduces risk of system-wide failures.
* Enables parallel development and clearer ownership boundaries.
* Easier to evolve individual services without affecting the entire system.
* Aligns well with cloud-native infrastructure (AWS, containers, orchestration).

---

### Negative

* Increased infrastructure complexity (Docker, orchestration, networking).
* Requires handling distributed system challenges (network failures, retries, timeouts).
* Data consistency becomes more complex (eventual consistency, distributed transactions).
* Higher operational costs compared to a monolith in early stages.
* Requires mature observability (logging, tracing, monitoring).

---

## Scope

This decision applies to all core backend components of FinScale.

A hybrid approach may be used internally:
- each microservice should follow a modular architecture (e.g., Clean Architecture),
- shared libraries should be minimized to avoid tight coupling.

---

## Pros and Cons of the Options

### Modular Monolith

* **Pros:**
    - Simpler deployment and infrastructure
    - Easier debugging and testing
    - Faster initial development

* **Cons:**
    - Difficult to scale specific components independently
    - Tight coupling between domains over time
    - Risk of becoming a “big ball of mud”

---

### Microservices Architecture

* **Pros:**
    - Independent scalability and deployment
    - Clear domain boundaries
    - Better fault isolation

* **Cons:**
    - High operational complexity
    - Requires distributed system expertise
    - Increased development overhead (communication, contracts)
