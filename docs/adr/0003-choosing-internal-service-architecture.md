# ADR-0003: Choosing internal service architecture (Layered vs Hexagonal vs Clean Architecture)

* **Status:** accepted
* **Deciders:** Yehor Reil (Lead Software Engineer)
* **Date:** 2026-04-06

---

## Context and Problem Statement

FinScale is built as a microservices-based system. Each service (Auth, Transactions, Analytics, Budgets) must handle complex business logic related to financial operations.

Given the critical nature of financial data, the architecture must:
- ensure strict separation between business logic and infrastructure,
- provide high testability,
- support long-term maintainability,
- prevent tight coupling to frameworks and external systems.

The key decision is selecting an internal architectural style for microservices.

---

## Decision Drivers

* **Separation of Concerns:** Clear isolation of domain logic from infrastructure.
* **Testability:** Ability to test core logic independently from external systems.
* **Maintainability:** Ease of evolving complex business rules.
* **Dependency Control:** Enforcing clear dependency direction.
* **Scalability of Codebase:** Ability to grow without becoming chaotic.
* **Consistency:** Unified approach across all services.

---

## Considered Options

* **Layered Architecture (N-tier)**
* **Hexagonal Architecture (Ports and Adapters)**
* **Clean Architecture**

---

## Decision Outcome

Chosen option: **Clean Architecture**

Clean Architecture was selected because it provides the most explicit control over dependency direction and enforces a strict separation between domain, application logic, and infrastructure.

This approach ensures that:
- business logic remains independent from frameworks and external systems,
- dependencies always point inward (toward the domain),
- the system remains maintainable as complexity grows,
- core logic can be tested in isolation without infrastructure concerns.

While Hexagonal Architecture offers similar principles, Clean Architecture provides a more structured and standardized model, which improves consistency across services.

---

## Consequences

### Positive

* Strong isolation of business logic from infrastructure.
* Clear dependency rules reduce architectural erosion over time.
* High testability of domain and application layers.
* Improved maintainability for complex financial logic.
* Consistent structure across all microservices.

---

### Negative

* Increased complexity and verbosity compared to simpler architectures.
* More boilerplate code (interfaces, layers, DTOs).
* Slower initial development speed.
* Requires strong discipline to correctly implement boundaries.
* Potential overengineering for small or simple services.

---

## Scope

This decision applies to all backend microservices in FinScale.

Each service should follow Clean Architecture principles:
- **Domain layer** contains core business logic and entities.
- **Application layer** handles use cases and orchestration.
- **Infrastructure layer** implements external integrations (DB, APIs, messaging).
- Dependencies must always point inward.

---

## Pros and Cons of the Options

### Layered Architecture (N-tier)

* **Pros:**
    - Simple and easy to understand
    - Fast to implement

* **Cons:**
    - Weak separation of concerns over time
    - Business logic often leaks into infrastructure
    - Harder to test independently

---

### Hexagonal Architecture (Ports and Adapters)

* **Pros:**
    - Good separation via ports and adapters
    - Flexible integration with external systems
    - High testability

* **Cons:**
    - Less explicit dependency model
    - Can vary in implementation between teams
    - Slightly less structured than Clean Architecture

---

### Clean Architecture

* **Pros:**
    - Strong and explicit dependency rules
    - Clear separation of layers
    - High maintainability and testability
    - Well-suited for complex domains

* **Cons:**
    - More complex to implement
    - Boilerplate overhead
    - Requires architectural discipline

---