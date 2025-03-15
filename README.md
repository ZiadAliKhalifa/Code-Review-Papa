# Code Review Papa 🧙‍♂️

An AI-powered code review assistant that automatically analyzes pull requests and provides insightful feedback.


[![Go Report Card](https://goreportcard.com/badge/github.com/ziadalikhalifa/code-review-papa)](https://goreportcard.com/report/github.com/ziadalikhalifa/code-review-papa)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- **Automated Code Reviews**: Analyzes pull requests when they're opened or updated.  
- **Intelligent Feedback**: Provides insights on code quality, potential bugs, and improvement suggestions.  
- **GitHub Integration**: Works as a GitHub App or with personal access tokens.  
- **Serverless Deployment**: Easily deploy as an AWS Lambda function.  

---

## Getting Started

### Prerequisites

Ensure you have the following installed:

- Go 1.16 or later  
- AWS CLI (for Lambda deployment)  
- GitHub account (for GitHub App creation)  
- DeepSeek API key  

---

## Installation

Clone the repository and install dependencies:

```sh
git clone https://github.com/ziadalikhalifa/code-review-papa.git
cd code-review-papa
go mod download
```

Build the application:

```sh
make build
```

---

## Configuration

### GitHub Authentication

Choose one of the following authentication methods:

#### **Personal Access Token (simpler)**  
```sh
GITHUB_TOKEN=your_github_personal_access_token
```

#### **GitHub App (recommended for production)**  
```sh
GITHUB_APP_ID=your_github_app_id
GITHUB_APP_PRIVATE_KEY=your_github_app_private_key
GITHUB_APP_INSTALLATION_ID=your_github_app_installation_id
```
Alternatively, provide the private key via a file:  
```sh
GITHUB_APP_PRIVATE_KEY_PATH=/path/to/private-key.pem
```

### AI Service Configuration

```sh
DEEPSEEK_KEY=your_deepseek_api_key
```

---

## Usage

### Local Development

Run the application locally to test with a specific PR:

```sh
go run main.go
```

By default, it will use the test PR defined in `main.go`. Modify the owner, repo, and PR number in the code to test with different PRs.

---

## AWS Lambda Deployment

1. **Build the Lambda package**  
   ```sh
   make build
   ```
   
2. **Create a new Lambda function (first time only)**  
   ```sh
   make create-lambda LAMBDA_ROLE_ARN=arn:aws:iam::your-account-id:role/your-lambda-role
   ```

3. **Update an existing Lambda function**  
   ```sh
   make deploy
   ```

4. **Create a function URL for webhook integration**  
   ```sh
   aws lambda create-function-url-config --function-name codeReviewPapa --auth-type NONE --region us-east-1
   ```

For more detailed deployment instructions, see [`LAMBDA_DEPLOYMENT.md`](LAMBDA_DEPLOYMENT.md).

---

## GitHub App Setup

1. Create a new GitHub App at [GitHub App Settings](https://github.com/settings/apps/new) with the following permissions:

   - **Repository permissions**:
     - Pull requests: Read & Write  
     - Contents: Read  
   - **Subscribe to events**:
     - Pull request  

2. Generate a private key and note your App ID.  
3. Install the app on your repositories and note the Installation ID.  
4. Configure your Lambda function with the GitHub App credentials.  
5. Set the webhook URL to your Lambda function URL.  

---

## How It Works

1. When a pull request is opened or updated, GitHub sends a webhook event to the Lambda function.  
2. The function fetches the PR diff from GitHub.  
3. The diff is sent to DeepSeek for AI analysis.  
4. The analysis is posted as a comment on the PR.

The comment includes:

- A summary of the changes  
- Potential issues or bugs  
- Suggestions for improvements  
- Security concerns  
- Code quality feedback  

---

## Project Structure

```
cmd/
  lambda/          # Lambda function entry point
config/            # Configuration handling
internal/
  ai/              # AI service integration
  analyzer/        # PR analysis logic
  github/          # GitHub API client
Makefile           # Build and deployment commands
main.go            # CLI entry point
```

---

## Contributing

Contributions are welcome! Follow these steps:

1. Fork the repository  
2. Create your feature branch  
   ```sh
   git checkout -b feature/amazing-feature
   ```
3. Commit your changes  
   ```sh
   git commit -m "Add some amazing feature"
   ```
4. Push to the branch  
   ```sh
   git push origin feature/amazing-feature
   ```
5. Open a Pull Request  

---

## License

This project is licensed under the MIT License. See the [`LICENSE`](LICENSE) file for details.

---

## Acknowledgements

- **DeepSeek** for providing AI code analysis capabilities.  
- **GitHub API** for PR integration.  
- **AWS Lambda** for serverless execution.  
- **Claude.ai** for helping me bang this out in 3 hours.

