module github.com/elangreza/e-commerce/product

go 1.24.3

replace github.com/elangreza/e-commerce/gen => ../gen

replace github.com/elangreza/e-commerce/pkg => ../pkg

require (
	github.com/elangreza/e-commerce/gen v0.0.0-00010101000000-000000000000
	github.com/elangreza/e-commerce/pkg v0.0.0-00010101000000-000000000000
	github.com/golang-migrate/migrate/v4 v4.19.0
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.77.0
)

require (
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-sqlite3 v1.14.32 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	golang.org/x/net v0.46.1-0.20251013234738-63d1a5100f82 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
