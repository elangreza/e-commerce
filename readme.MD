# E-COMMERCE (Orchestrated Saga Pattern)

this is a simple e-commerce application that allows users to place orders for products. The application is built with scalability and reliability in mind, utilizing best practices for handling orders and payments.

Built with the Saga orchestration pattern and a central orchestrator. The orchestrator runs each step needed to complete an order, makes sure steps happen in the right order, and handles failures cleanly. This project does not use a choreography (event-driven / decentralized) approach.

# TASKS

## HIGH PRIORITY TASKS

// TODO Add status High priority. Full visibility into order lifecycle
// TODO Write integration tests. High priority. Confidence to ship faster

## MEDIUM PRIORITY TASKS

// TODO Build compensation retry system. Medium priority. Ensures system recovers from failures
// TODO Add structured logging & tracing. Medium priority. Faster debugging

## LOW PRIORITY TASKS

// TODO Add basic health/metrics. Low priority. Operational awareness
