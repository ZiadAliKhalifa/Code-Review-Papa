.PHONY: build clean deploy

# Configuration
LAMBDA_FUNCTION_NAME=codeReviewPapa
AWS_REGION=us-east-1
BINARY_NAME=bootstrap
BUILD_DIR=./build

build:
	@echo "Building Lambda function..."
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/lambda
	cd $(BUILD_DIR) && zip function.zip $(BINARY_NAME)
	@echo "Build complete"

clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)
	@echo "Clean complete"

deploy: build
	@echo "Deploying Lambda function..."
	aws lambda update-function-code \
		--function-name $(LAMBDA_FUNCTION_NAME) \
		--zip-file fileb://$(BUILD_DIR)/function.zip \
		--region $(AWS_REGION)
	@echo "Deployment complete"

create-lambda:
	@echo "Creating Lambda function..."
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