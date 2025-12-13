# E-Commerce Microservices Code Review

## Executive Summary

This is a **microservices-based e-commerce application** built with Go, implementing an **orchestrated Saga pattern** for distributed transaction management. The system consists of 6 services communicating via gRPC, with a Backend-for-Frontend (BFF) API gateway exposing REST endpoints.

**Overall Assessment**: The codebase demonstrates solid architectural decisions and good separation of concerns. However, it lacks critical production-readiness features like comprehensive testing, structured logging, and robust error handling.

---

## Architecture Overview

### Services

1. **API Service** - BFF/Gateway (REST → gRPC)
2. **Order Service** - Saga orchestrator for order processing
3. **Product Service** - Product catalog management
4. **Payment Service** - Payment processing (mocked)
5. **Warehouse Service** - Inventory and stock management
6. **Shop Service** - Shop/merchant management

### Technology Stack

- **Language**: Go
- **Communication**: gRPC (inter-service), REST (client-facing)
- **Database**: SQLite3 (dev) / PostgreSQL (production-ready)
- **Router**: go-chi
- **Authentication**: JWT tokens
- **Deployment**: Docker Compose

---

## Strengths

### 1. **Excellent Architectural Patterns**

#### ✅ Orchestrated Saga Pattern

The order service implements a well-structured Saga orchestration:

```go
// Order creation with coordinated steps
1. Validate cart and products
2. Create order (PENDING status)
3. Reserve stock → Update to STOCK_RESERVED
4. Process payment
5. Handle failures with compensating transactions
```

**Benefits**:

- Centralized coordination logic
- Clear transaction boundaries
- Explicit rollback mechanisms

#### ✅ Clean Separation of Concerns

Each service follows a layered architecture:

- `cmd/` - Entry points
- `internal/service/` - Business logic
- `internal/sqlitedb/` - Data access
- `internal/entity/` - Domain models
- `internal/rest/` or `internal/server/` - Transport layer

### 2. **Good Code Organization**

#### ✅ Shared Package (`pkg/`)

Reusable utilities across services:

- `dbsql` - Database abstraction with migration support
- `money` - Currency handling (prevents floating-point errors)
- `extractor` - Metadata extraction utilities
- `gracefulshutdown` - Proper shutdown handling
- `interceptor` - gRPC middleware

#### ✅ Interface-Based Design

Services use interfaces for dependencies, enabling testability:

```go
type orderRepo interface {
    CreateOrder(ctx context.Context, order entity.Order) (uuid.UUID, error)
    UpdateOrder(ctx context.Context, payloads map[string]any, orderID uuid.UUID) error
    // ...
}
```

### 3. **Idempotency Support**

Order creation uses idempotency keys to prevent duplicate orders:

```go
ord, err := s.orderRepo.GetOrderByIdempotencyKey(ctx, idempotencyKey)
if ord != nil {
    return ord.GetGenOrder(), nil // Return existing order
}
```

### 4. **Database Migration Support**

Automated migrations using `golang-migrate`:

- Version-controlled schema changes
- Support for both SQLite and PostgreSQL
- Separate seeder support

### 5. **Money Handling**

Proper currency handling using protobuf `Money` type:

- Avoids floating-point arithmetic errors
- Enforces single-currency transactions
- Safe multiplication and addition operations

---

## Critical Issues

### 🔴 1. **No Testing**

**Severity**: CRITICAL

**Issue**: Zero test coverage across the entire codebase.

**Impact**:

- No confidence in code correctness
- High risk of regressions
- Difficult to refactor safely
- Cannot validate Saga compensation logic

**Recommendation**:

```go
// Example: Order service unit test
func TestCreateOrder_Success(t *testing.T) {
    mockOrderRepo := mock.NewMockOrderRepo(ctrl)
    mockCartRepo := mock.NewMockCartRepo(ctrl)
    // ... setup mocks

    svc := NewOrderService(mockOrderRepo, mockCartRepo, ...)
    order, err := svc.CreateOrder(ctx, req)

    assert.NoError(t, err)
    assert.Equal(t, constanta.OrderStatusStockReserved, order.Status)
}
```

**Priority**: HIGH - Add unit tests for business logic first, then integration tests.

---

### 🔴 2. **Inadequate Error Handling**

**Severity**: HIGH

**Issues**:

1. **Printf instead of proper logging**:

```go
// ❌ Bad - from order_service.go:317
fmt.Printf("Error during rollback: %v", rollbackErr)
```

2. **Silent error swallowing**:

```go
// ❌ Bad - from order_service.go:378
if err != nil {
    fmt.Println("err when Update status", err)
}
// Continues execution despite error
```

3. **Debug print statements in production code**:

```go
// ❌ Bad - from payment_service.go:115-154
fmt.Println("cek", 1)
fmt.Println("cek", 2)
// ... scattered throughout
```

4. **Inconsistent error context**:

```go
// ❌ Bad
return nil, errors.New("failed to fetch products")

// ✅ Good
return nil, fmt.Errorf("failed to fetch products: %w", err)
```

**Recommendations**:

1. **Implement structured logging**:

```go
import "log/slog"

logger.Error("rollback failed",
    slog.String("order_id", orderID.String()),
    slog.String("error", rollbackErr.Error()),
)
```

2. **Return errors instead of logging**:

```go
// Let caller decide how to handle
if err != nil {
    return fmt.Errorf("failed to update order status: %w", err)
}
```

3. **Remove debug prints** - Use proper logging levels instead.

---

### 🔴 3. **No Observability**

**Severity**: HIGH

**Missing**:

- ❌ Structured logging
- ❌ Distributed tracing
- ❌ Metrics/monitoring
- ❌ Health checks
- ❌ Request IDs across services

**Impact**: Impossible to debug production issues in distributed system.

**Recommendations**:

1. **Add structured logging** (slog, zap, or zerolog)
2. **Implement OpenTelemetry** for tracing
3. **Add Prometheus metrics**
4. **Health check endpoints** (`/health`, `/ready`)
5. **Propagate request IDs** via gRPC metadata

---

### 🟡 4. **Incomplete Saga Implementation**

**Severity**: MEDIUM

**Issues**:

1. **Missing payment success callback**:

```go
// TODO from README.md:
// - TODO send callback from payment service
// - TODO add background worker to limit payment waiting time
```

Currently, payment is created with `WAITING` status but there's no mechanism to:

- Confirm payment success
- Timeout expired payments
- Notify order service of payment completion

2. **No compensation retry logic**:

```go
// From order_service.go:314-318
_, err = s.warehouseServiceClient.ReserveStock(ctx, req)
if err != nil {
    rollbackErr := rollback()
    // ❌ No retry mechanism if rollback fails
    return nil, fmt.Errorf("failed to reserve stock: %w", err)
}
```

**Recommendations**:

1. **Add payment callback endpoint**:

```go
// In order service
func (s *orderService) HandlePaymentCallback(ctx context.Context, req *gen.PaymentCallbackRequest) (*gen.Empty, error) {
    if req.Status == "PAID" {
        // Confirm stock, update order to CONFIRMED
    } else {
        // Release stock, mark order as FAILED
    }
}
```

2. **Implement background worker** for payment timeout:

```go
// Periodically check for expired payments
func (s *orderService) CleanupExpiredPayments() {
    payments := s.getPaymentsOlderThan(15 * time.Minute)
    for _, payment := range payments {
        s.rollbackPayment(payment)
    }
}
```

3. **Add retry mechanism** with exponential backoff for compensations.

---

### 🟡 5. **Security Concerns**

**Severity**: MEDIUM

**Issues**:

1. **No authorization/RBAC**:

```go
// ❌ Anyone with valid JWT can access warehouse management
r.Use(authMiddleware.MustAuthMiddleware())
r.Post("/warehouse/status", oh.SetWarehouseStatus())
```

2. **CORS allows all origins**:

```go
// ❌ api/cmd/server/main.go:48
AllowedOrigins: []string{"*"},
```

3. **No rate limiting** on API endpoints.

4. **No input sanitization** beyond basic validation.

**Recommendations**:

1. **Implement RBAC**:

```go
type UserRole string
const (
    RoleAdmin    UserRole = "admin"
    RoleCustomer UserRole = "customer"
)

func (m *AuthMiddleware) RequireRole(role UserRole) func(http.Handler) http.Handler {
    // Check user role from JWT claims
}
```

2. **Restrict CORS** to specific domains in production.
3. **Add rate limiting** (e.g., using `golang.org/x/time/rate`).

---

### 🟡 6. **Database Concerns**

**Severity**: MEDIUM

**Issues**:

1. **Timestamps not in UTC**:

```go
// TODO from README.md:
// - TODO save time in UTC format in the database
```

2. **No connection pooling configuration** in most services.

3. **Potential N+1 queries**:

```go
// order_service.go:228-234
// Fetches all products in one call ✅
products, err := s.productServiceClient.GetProducts(ctx, &gen.GetProductsRequest{
    Ids: cart.GetProductIDs(),
})
```

This is actually well-done, but ensure similar patterns throughout.

**Recommendations**:

1. **Use UTC everywhere**:

```go
time.Now().UTC()
```

2. **Configure connection pools**:

```go
dbsql.WithDBConnectionPool(25, 5, 5*time.Minute)
```

---

### 🟡 7. **Code Quality Issues**

**Severity**: LOW-MEDIUM

**Issues**:

1. **Inconsistent error messages**:

```go
// Some use status.Error, some use fmt.Errorf
return status.Errorf(codes.NotFound, "cart not found")
return fmt.Errorf("failed to get cart: %w", err)
```

2. **Magic numbers**:

```go
// order_service.go:91
if req.Quantity > product.Stock {
    // No max quantity constant
}
```

3. **Commented-out code and TODOs in production**:

```go
// payment_service.go - multiple "cek" debug prints
```

4. **Inconsistent naming**:

```go
type AutService interface {} // Should be AuthService
```

**Recommendations**:

1. **Standardize error handling** - Choose gRPC status codes for gRPC services.
2. **Extract constants**:

```go
const MaxOrderQuantity = 1000
```

3. **Remove debug code** before production.
4. **Fix typos** and naming inconsistencies.

---

## Medium Priority Issues

### 🟢 8. **Missing Features**

From the README TODO list:

1. **Order details endpoint** - Users can't view order history
2. **Product details endpoint** - Limited product information
3. **Shop active status filtering** - Can't filter by active shops
4. **Payment UI** - No way to actually "pay" in the mock service

**Recommendations**: Prioritize based on user needs.

---

### 🟢 9. **Docker & Deployment**

**Current State**: Basic Docker Compose setup works well.

**Issues**:

1. **No health checks** in docker-compose.yaml
2. **No resource limits**
3. **No CI/CD pipeline**

**Recommendations**:

```yaml
services:
  api:
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
    deploy:
      resources:
        limits:
          cpus: "0.5"
          memory: 512M
```

---

## Positive Patterns to Maintain

### ✅ 1. **Graceful Shutdown**

```go
gs := gracefulshutdown.New(context.Background(), 5*time.Second,
    gracefulshutdown.Operation{
        Name: "server",
        ShutdownFunc: func(ctx context.Context) error {
            return srv.Shutdown(ctx)
        }},
)
```

### ✅ 2. **Context Propagation**

User ID properly propagated via gRPC metadata:

```go
func AppendUserIDintoContextGrpcClient(ctx context.Context, userID uuid.UUID) context.Context {
    md := metadata.New(map[string]string{string(globalcontanta.UserIDKey): userID.String()})
    return metadata.NewOutgoingContext(ctx, md)
}
```

### ✅ 3. **Validation at Boundaries**

```go
if err := body.Validate(); err != nil {
    sendErrorResponse(w, http.StatusBadRequest, err)
    return
}
```

---

## Actionable Recommendations (Prioritized)

### 🔥 High Priority (Do First)

1. **Add Unit Tests** (1-2 weeks)

   - Start with order service Saga logic
   - Mock external dependencies
   - Target 70%+ coverage for business logic

2. **Implement Structured Logging** (2-3 days)

   - Replace all `fmt.Printf` with proper logger
   - Add log levels (debug, info, warn, error)
   - Include context (request ID, user ID, order ID)

3. **Remove Debug Code** (1 day)

   - Remove all `fmt.Println("cek", ...)` statements
   - Clean up commented code
   - Fix naming inconsistencies

4. **Complete Payment Flow** (3-5 days)

   - Add payment callback endpoint
   - Implement background worker for payment timeout
   - Add order confirmation flow

5. **Add Health Checks** (1 day)
   - `/health` endpoint for each service
   - Check database connectivity
   - Update docker-compose with health checks

### 🟡 Medium Priority (Next Sprint)

6. **Add Integration Tests** (1 week)

   - Test full order flow end-to-end
   - Test failure scenarios and rollbacks
   - Use testcontainers for database

7. **Implement Observability** (1-2 weeks)

   - Add OpenTelemetry tracing
   - Add Prometheus metrics
   - Set up Grafana dashboards

8. **Add Authorization/RBAC** (3-5 days)

   - Define user roles
   - Protect admin endpoints
   - Add role-based middleware

9. **Improve Error Handling** (1 week)
   - Standardize error types
   - Add error codes for client handling
   - Implement retry logic for compensations

### 🟢 Low Priority (Future)

10. **Add CI/CD Pipeline**

    - GitHub Actions or GitLab CI
    - Automated testing
    - Docker image building

11. **Add Monitoring & Alerting**

    - Set up Prometheus + Grafana
    - Alert on high error rates
    - Monitor Saga completion rates

12. **Performance Optimization**
    - Add caching layer (Redis)
    - Optimize database queries
    - Add connection pooling

---

## Code Examples: Before & After

### Error Handling

#### ❌ Before

```go
if rollbackErr != nil {
    fmt.Printf("Error during rollback: %v", rollbackErr)
}
```

#### ✅ After

```go
if rollbackErr != nil {
    logger.Error("order rollback failed",
        slog.String("order_id", orderID.String()),
        slog.String("user_id", userID.String()),
        slog.String("error", rollbackErr.Error()),
    )
    // Consider alerting for manual intervention
    return fmt.Errorf("failed to rollback order %s: %w", orderID, rollbackErr)
}
```

### Testing

#### ✅ Add This

```go
func TestOrderService_CreateOrder_StockReservationFails(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockWarehouse := mock.NewMockWarehouseServiceClient(ctrl)
    mockWarehouse.EXPECT().
        ReserveStock(gomock.Any(), gomock.Any()).
        Return(nil, errors.New("insufficient stock"))

    svc := NewOrderService(orderRepo, cartRepo, mockWarehouse, ...)

    _, err := svc.CreateOrder(ctx, validRequest)

    assert.Error(t, err)
    assert.Contains(t, err.Error(), "failed to reserve stock")

    // Verify order status is FAILED
    order, _ := orderRepo.GetOrderByIdempotencyKey(ctx, idempotencyKey)
    assert.Equal(t, constanta.OrderStatusFailed, order.Status)
}
```

---

## Conclusion

### Summary

This is a **well-architected microservices application** with solid foundations:

- ✅ Clean architecture and separation of concerns
- ✅ Proper Saga orchestration pattern
- ✅ Good use of Go idioms and interfaces
- ✅ Idempotency and transaction safety

However, it's **not production-ready** due to:

- ❌ Zero test coverage
- ❌ Poor error handling and logging
- ❌ No observability
- ❌ Incomplete Saga implementation
- ❌ Security gaps

### Recommendation

**Focus on the High Priority items** before considering this production-ready. The architecture is sound, but operational concerns (testing, logging, monitoring) need immediate attention.

**Estimated effort to production-ready**: 4-6 weeks with 2 developers.

### Final Grade: **B-**

- Architecture: A
- Code Quality: B
- Testing: F
- Observability: F
- Security: C
- Documentation: B+

**Overall**: Good foundation, needs production hardening.
