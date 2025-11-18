# E-COMMERCE (Orchestrated Saga Pattern)

this is a simple e-commerce application that allows users to place orders for products. The application is built with scalability and reliability in mind, utilizing best practices for handling orders and payments.

Built with the Saga orchestration pattern and a central orchestrator. The orchestrator runs each step needed to complete an order, makes sure steps happen in the right order, and handles failures cleanly. This project does not use a choreography (event-driven / decentralized) approach.

# TASKS

## HIGH PRIORITY TASKS

// TODO Write integration tests. High priority. Confidence to ship faster

## MEDIUM PRIORITY TASKS

// TODO Build compensation retry system. Medium priority. Ensures system recovers from failures
// TODO Add structured logging & tracing. Medium priority. Faster debugging

## LOW PRIORITY TASKS

// TODO Add basic health/metrics. Low priority. Operational awareness

TODO change the structure of the project to the following:
since the API service will call multiple other services, it makes sense to have a clear separation of concerns and organize the code accordingly.
The proposed structure is as follows:
see this link https://chat.qwen.ai/c/f8326237-936a-4c0b-a991-064610f0b4a6
order-service/
├── cmd/
│ └── server/
│ └── main.go
├── internal/
│ ├── config/
│ ├── server/
│ ├── service/ # Business logic
│ │ └── order_service.go # Works directly with pb types
│ ├── repository/ # DB logic
│ │ └── order_repo.go
│ └── pb/ # Generated protos
│ └── order.pb.go
└── proto/
└── order.proto
