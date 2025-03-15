# AWS Lambda Deployment Guide

This guide explains how to deploy Code Review Papa as an AWS Lambda function.

## Prerequisites

1. AWS CLI installed and configured with appropriate credentials
2. Go 1.16 or later
3. An AWS IAM role with Lambda execution permissions

## Deployment Steps

### 1. Create an IAM Role for Lambda

Create an IAM role with the following policies:
- `AWSLambdaBasicExecutionRole` (for CloudWatch Logs)
- Custom policy for any other AWS services your Lambda needs to access

Note the ARN of this role for the next step.

### 2. Build and Deploy the Lambda Function

```bash
# Create the build directory
mkdir -p build

# Set your Lambda role ARN
export LAMBDA_ROLE_ARN=arn:aws:iam::123456789012:role/your-lambda-role

# Build and create the Lambda function (first time)
make build
make create-lambda

# For subsequent deployments
make deploy