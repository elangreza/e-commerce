module github.com/elangreza/e-commerce/api

go 1.24.3

replace github.com/elangreza/e-commerce/gen => ../gen

require (
	github.com/elangreza/e-commerce/gen v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.75.0
)

require (
	github.com/golang/protobuf v1.5.4 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
