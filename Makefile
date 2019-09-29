.PHONY: build

BucketName = sam-postory-deploy-bucket
TestShippoToken = ''

# build lambda handlers
build:
	GOOS=linux GOARCH=amd64 go build -o ./build/track handler/track/main.go
	GOOS=linux GOARCH=amd64 go build -o ./build/history handler/history/main.go

# start local server using test shippo API
local:
	make build
	sam local start-api --host localhost

# clean build artifacts
clean:
	rm -rf ./build

# run unit tests
test:
	go test ./... -v

# run integration tests
test-integration:
	SHIPPO_TOKEN=$(TestShippoToken) go test ./... -tags=integration -v

# run all tests
test-all:
	make test
	make test-integration

# deploy to aws
deploy:
	make build
	sam package \
		--template-file template.yaml \
		--output-template-file build/packaged-template.yml \
		--s3-bucket $(BucketName)
	sam deploy \
		--template-file build/packaged-template.yml \
		--stack-name postory \
		--capabilities CAPABILITY_IAM \
		--parameter-overrides 'Stage=prod'
