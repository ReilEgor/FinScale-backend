# [ADR-0001] Choosing the main programming language for FinScale microservices

* **Status:** accepted
* **Deciders:** Yehor Reil(Lead Software Engineer)
* **Date:** 2026-03-29

## Context and Problem Statement

The FinScale project requires the development of a scalable microservice architecture (services: Auth, Transactions, Analytics, Budgets). The system must process financial transactions in real time, provide high type safety, and operate efficiently on a cloud infrastructure (AWS).

You need to choose a primary development language that balances performance, development speed, and infrastructure ownership costs.

## Decision Drivers
* **Type Safety:** Minimizing errors when working with monetary data.
* **Concurrency:** Efficient parallel processing of queries and data streams.
* **Resource Efficiency:** Low RAM/CPU consumption to optimize AWS costs.
* **Maintainability:** Easy to maintain and scale the code base.

## Considered Options
* **Go (Golang)**
* **Python (FastAPI)**
* **Rust**
* **C# (ASP.NET Core)**

## Decision Outcome
Chosen option: **Go (Golang)** because it offers the best balance between out-of-the-box performance and speed of writing clean, maintainable code for microservices.

### Consequences
* **Good:** Fast launch of services and small size of Docker images.
* **Good:** Native gRPC support and strong typing simplify integration between services.
* **Good:** The goroutine model allows you to easily handle thousands of concurrent connections.
* **Bad:** Error handling via `if err != nil` increases code size.
* **Bad:** Fewer libraries for deep Data Science analysis compared to Python.

---

## Pros and Cons of the Options

### Go (Golang)
* **Pros:** Compilation to static binaries; minimal resource consumption; high execution speed.
* **Cons:** Lack of classic OOP (class hierarchy), which requires a different approach to architecture.

### Python (FastAPI)
* **Pros:** Huge number of libraries for analytics and finance; very fast prototyping.
* **Cons:** Performance issues under high load (GIL); high memory consumption.

### Rust
* **Pros:** Maximum memory safety and C++-level performance.
* **Cons:** High cognitive load; long compilation times, which slows down CI/CD cycles.

### C# (ASP.NET Core)
* **Pros:** Mature ecosystem, excellent development tools.
* **Cons:** Heavier runtime compared to Go; higher entry barrier for cloud optimization.

## More Information
This decision secures Go as the primary backend language. FinScale will continue to use React for the frontend, and REST/gRPC will be used for integration.