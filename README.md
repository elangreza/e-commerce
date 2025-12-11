# E-COMMERCE (Orchestrated Saga Pattern)

this is a simple e-commerce application that allows users to place orders for products. The application is built with scalability and reliability in mind, utilizing best practices for handling orders and payments.

Built with the **Saga orchestration pattern and a central orchestrator**. The orchestrator runs each step needed to complete an order, makes sure steps happen in the right order, and handles failures cleanly. This project does not use a **choreography (event-driven / decentralized)** approach.

# TASKS

## HIGH PRIORITY TASKS

- TODO ADD flow order success

  - TODO send callback from payment service
  - TODO add background worker to limit payment waiting time, must be less than order timeout process in order service

- TODO ADD order details in order service
  - TODO ADD product details in product service
  - TODO ADD product details in BFF / API service
- TODO API get product with shop is active or not
- TODO Add unit tests. High priority. Confidence to ship faster
- TODO Make sure the mocked payment service can be accessed via HTTP / own UI
- TODO Write integration tests. High priority. Confidence to ship faster
- TODO Add CI/CD pipeline. High priority. Automate testing & deployment
- TODO add better error handling. High priority. Improve reliability
- TODO add logging. High priority. Easier debugging, observability
- TODO add user role, can get from content-management system repo (CMS) later and restrict the warehouse management endpoints to admin users only
- TODO save time in UTC format in the database

## MEDIUM PRIORITY TASKS

- TODO Build compensation retry system. Medium priority. Ensures system recovers from failures
- TODO Add structured logging & tracing. Medium priority. Faster debugging

## LOW PRIORITY TASKS

- TODO Add basic health/metrics. Low priority. Operational awareness
- TODO Add authentication/authorization. Low priority. Secure access in API service

# ARCHITECTURE OVERVIEW

The application is composed of the following microservices:

- API service. It runs as a gateway to route requests to other services. This is the public access only Includes:
  - User authentication,
  - Product catalog,
  - Order processing,
  - Warehouse management.
- Order service,
- Product service,
- Shop service,
- Warehouse service

here's the list of technologies used in this project:

- API gateway: go-chi router
- Authentication: JWT tokens
- Communication between services: gRPC with go
- Communication between client: Backend for Frontend (BFF) pattern with REST API
- Database: can be run with either Sqlite3 or PostgreSQL

here's the list of API endpoints exposed by the API service:

### Register a new user

| Field            | Value                                                        |
| ---------------- | ------------------------------------------------------------ |
| **Endpoint**     | `POST /auth/register`                                        |
| **URL**          | `http://localhost:8080/auth/register`                        |
| **Content-Type** | `application/json`                                           |
| **Success Code** | `201 Created`                                                |
| **Description**  | Registers a new user account with email, password, and name. |

<details>
<summary><b><i>Click here for the curl!</i></b></summary>

```bash
curl --location 'http://localhost:8080/auth/register' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email":"test@test.com",
    "password":"test",
    "name":"test"
}'
```

</details>

---

### Login and obtain a JWT token

| Field            | Value                                                           |
| ---------------- | --------------------------------------------------------------- |
| **Endpoint**     | `POST /auth/login`                                              |
| **URL**          | `http://localhost:8080/auth/login`                              |
| **Content-Type** | `application/json`                                              |
| **Success Code** | `200 OK`                                                        |
| **Description**  | Authenticates a user and returns a JWT for protected endpoints. |

<details>
<summary><b><i>Click here for the curl!</i></b></summary>

```bash
curl --location 'http://localhost:8080/auth/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email":"test@test.com",
    "password":"test"
}'
```

</details>

---

### Get a list of products

| Field            | Value                                                                                                                 |
| ---------------- | --------------------------------------------------------------------------------------------------------------------- |
| **Endpoint**     | `GET /products`                                                                                                       |
| **URL**          | `http://localhost:8080/products`                                                                                      |
| **Content-Type** | —                                                                                                                     |
| **Success Code** | `200 OK`                                                                                                              |
| **Description**  | Retrieves a paginated, optionally filtered list of products. Supports `page`, `limit`, and `search` query parameters. |

<details>
<summary><b><i>Click here for the curl!</i></b></summary>

```bash
curl --location 'http://localhost:8080/products?page=2&limit=10&search=men'
```

</details>

---

### Add a product to the cart

| Field             | Value                                                                    |
| ----------------- | ------------------------------------------------------------------------ |
| **Endpoint**      | `POST /cart`                                                             |
| **URL**           | `http://localhost:8080/cart`                                             |
| **Content-Type**  | `application/json`                                                       |
| **Authorization** | `Bearer <JWT>`                                                           |
| **Success Code**  | `201 Created`                                                            |
| **Description**   | Adds a specified quantity of a product to the authenticated user’s cart. |

<details>
<summary><b><i>Click here for the curl!</i></b></summary>

```bash
curl --location 'http://localhost:8080/cart' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {{token from login API}}' \
--data '{
    "product_id":"019394d0-4d5e-7d6a-9c4b-8a3f2e1d5ca2",
    "quantity":79
}'
```

</details>

---

### Get the current cart contents

| Field             | Value                                                                |
| ----------------- | -------------------------------------------------------------------- |
| **Endpoint**      | `GET /cart`                                                          |
| **URL**           | `http://localhost:8080/cart`                                         |
| **Content-Type**  | —                                                                    |
| **Authorization** | `Bearer <JWT>`                                                       |
| **Success Code**  | `200 OK`                                                             |
| **Description**   | Returns the full contents of the authenticated user’s shopping cart. |

<details>
<summary><b><i>Click here for the curl!</i></b></summary>

```bash
curl --location 'http://localhost:8080/cart' \
--header 'Authorization: Bearer {{token from login API}}'
```

</details>

---

### Create a new order based on the cart

| Field             | Value                                                                                                                |
| ----------------- | -------------------------------------------------------------------------------------------------------------------- |
| **Endpoint**      | `POST /order`                                                                                                        |
| **URL**           | `http://localhost:8080/order`                                                                                        |
| **Content-Type**  | `application/json`                                                                                                   |
| **Authorization** | `Bearer <JWT>`                                                                                                       |
| **Success Code**  | `201 Created`                                                                                                        |
| **Description**   | Converts the user’s current cart into a confirmed order. Uses an `idempotency_key` to prevent duplicate submissions. |

<details>
<summary><b><i>Click here for the curl!</i></b></summary>

```bash
curl --location 'http://localhost:8080/order' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {{token from login API}}' \
--data '{
    "idempotency_key":"75b12b36-8547-4c02-9783-d42007f6a92a"
}'
```

</details>

---

### Set warehouse status (active/inactive)

| Field             | Value                                                                                             |
| ----------------- | ------------------------------------------------------------------------------------------------- |
| **Endpoint**      | `POST /warehouse/status`                                                                          |
| **URL**           | `http://localhost:8080/warehouse/status`                                                          |
| **Content-Type**  | `application/json`                                                                                |
| **Authorization** | `Bearer <JWT>`                                                                                    |
| **Success Code**  | `200 OK`                                                                                          |
| **Description**   | Updates the operational status (`is_active`) of a warehouse. Typically restricted to admin users. |

<details>
<summary><b><i>Click here for the curl!</i></b></summary>

```bash
curl --location 'http://localhost:8080/warehouse/status' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {{token from login API}}' \
--data '{
    "warehouse_id":1,
    "is_active": true
}'
```

</details>

---

### Transfer stock between warehouses

| Field             | Value                                                                                             |
| ----------------- | ------------------------------------------------------------------------------------------------- |
| **Endpoint**      | `POST /warehouse/transfer`                                                                        |
| **URL**           | `http://localhost:8080/warehouse/transfer`                                                        |
| **Content-Type**  | `application/json`                                                                                |
| **Authorization** | `Bearer <JWT>`                                                                                    |
| **Success Code**  | `200 OK`                                                                                          |
| **Description**   | Moves a specified quantity of a product from one warehouse to another. Requires admin privileges. |

<details>
<summary><b><i>Click here for the curl!</i></b></summary>

```bash
curl --location 'http://localhost:8080/warehouse/transfer' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {{token from login API}}' \
--data '{
    "from_warehouse_id": 1,
    "to_warehouse_id": 3,
    "product_id": "019394d0-4d5e-7d6a-9c4b-8a3f2e1d5ca2",
    "quantity": 10
}'
```

</details>
