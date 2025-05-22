.PHONY: build clean deploy test create-lambda

# Configuration
LAMBDA_FUNCTION_NAME=codeReviewPapa
AWS_REGION=us-east-1
BINARY_NAME=bootstrap
BUILD_DIR=./build
# LAMBDA_ROLE_ARN should be defined, e.g., as an environment variable or directly below
# LAMBDA_ROLE_ARN=arn:aws:iam::123456789012:role/your-lambda-role

build:
	@echo "Building Lambda function..."
	# Ensure build directory exists
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/lambda
	cd $(BUILD_DIR) && zip function.zip $(BINARY_NAME)
	@echo "Build complete"

clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)
	@echo "Clean complete"

test:
	@echo "Running tests..."
	go test -v ./...

deploy: build
	@echo "Deploying Lambda function..."
	aws lambda update-function-code \
		--function-name $(LAMBDA_FUNCTION_NAME) \
		--zip-file fileb://$(BUILD_DIR)/function.zip \
		--region $(AWS_REGION)
	@echo "Deployment complete"

create-lambda: build
	@echo "Creating Lambda function..."
	@if [ -z "$(LAMBDA_ROLE_ARN)" ]; then \
		echo "Error: LAMBDA_ROLE_ARN is not set. Please set it as an environment variable or define it in the Makefile."; \
		exit 1; \
	fi
	aws lambda create-function \
		--function-name $(LAMBDA_FUNCTION_NAME) \
		--runtime provided.al2 \
		--role $(LAMBDA_ROLE_ARN) \
		--handler bootstrap \
		--zip-file fileb://$(BUILD_DIR)/function.zip \
		--region $(AWS_REGION) \
		--timeout 30 \
		--memory-size 256
	@echo "Lambda function created"