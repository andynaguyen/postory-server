run:
	go run cmd/run.go

test:
	make test-unit
	make test-integration

test-integration:
	go test ./... -tags=integration -v

test-unit:
	go test ./... -v
