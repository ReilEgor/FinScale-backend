# [ADR-0002] Choosing the database management systems (DBMS) for microservices

* **Status:** accepted
* **Deciders:** Yehor Reil (Lead Software Engineer)
* **Date:** 2026-04-06

## Context and Problem Statement

The FinScale project is based on a microservice architecture (Auth, Transactions, Analytics, Budgets), where each service owns its data. The system processes financial data, requires high reliability, consistency for critical operations, and scalability for analytical workloads.

You need to choose suitable DBMS solutions that balance consistency, scalability, performance, and operational complexity for different service needs.

## Decision Drivers
* **Data Consistency:** Strong guarantees for financial transactions.
* **Scalability:** Ability to scale independently per service.
* **Performance:** Low latency for transactional operations and high throughput for analytics.
* **Flexibility:** Support for different data models across services.
* **Operational Complexity:** Ease of deployment and maintenance in AWS.

## Considered Options
* **PostgreSQL**
* **MongoDB**
* **MySQL**
* **DynamoDB**
* **Redis (as primary DB)**

## Decision Outcome
Chosen option: **Polyglot persistence approach**, primarily using **PostgreSQL** for transactional services and **MongoDB** for analytics and flexible data storage, with **Redis** as a caching layer.

### Consequences
* **Good:** Strong consistency and ACID guarantees for financial data using PostgreSQL.
* **Good:** Flexible schema in MongoDB allows rapid iteration for analytics and reporting services.
* **Good:** Redis improves performance via caching and reduces load on primary databases.
* **Bad:** Increased operational complexity due to multiple database technologies.
* **Bad:** Requires careful data synchronization and consistency management across services.

---

## Pros and Cons of the Options

### PostgreSQL
* **Pros:** Strong ACID compliance; rich querying capabilities; mature ecosystem.
* **Cons:** Vertical scaling limitations; more complex sharding.

### MongoDB
* **Pros:** Schema flexibility; easy horizontal scaling; good fit for semi-structured data.
* **Cons:** Weaker consistency guarantees compared to relational databases.

### MySQL
* **Pros:** Widely adopted; simple to operate; good performance for read-heavy workloads.
* **Cons:** Less advanced features compared to PostgreSQL; limited extensibility.

### DynamoDB
* **Pros:** Fully managed; automatic scaling; high availability out of the box.
* **Cons:** Vendor lock-in (AWS); limited querying flexibility; cost at scale.

### Redis (as primary DB)
* **Pros:** Extremely fast; ideal for caching and ephemeral data.
* **Cons:** Not suitable as a primary persistent store for critical financial data.

## More Information
Each microservice is responsible for its own database (Database per Service pattern). PostgreSQL is used for core financial domains (Transactions, Auth), while MongoDB supports Analytics. Redis is used for caching, session storage, and rate limiting.
