# GitHub Actions CI/CD Pipelines

This directory contains the CI/CD workflows for the Akeneo Serverless Connector project.

## Workflows

### 1. CI Pipeline (`ci.yml`)

Runs on every push and pull request to `main` and `develop` branches.

**Jobs:**
- **test-aws**: Runs Go tests with race detection and coverage
- **lint-aws**: Runs golangci-lint for code quality checks
- **build-aws**: Builds the Lambda function and creates deployment artifact
- **security-scan**: Runs Gosec security scanner

**Artifacts:**
- Test coverage report (HTML)
- Lambda function zip file (retained for 7 days)

### 2. PR Checks (`pr-checks.yml`)

Runs additional checks on pull requests.

**Jobs:**
- **validate**: Validates PR title follows conventional commits
- **size-check**: Ensures Lambda package size is within limits
- **comment-coverage**: Posts test coverage report as PR comment

### 3. Release (`release.yml`)

Triggered when a version tag (e.g., `v1.0.0`) is pushed.

**Jobs:**
- **release**: Builds Lambda for both amd64 and arm64 architectures
  - Creates GitHub release with changelog
  - Attaches both architecture builds as release assets

**Usage:**
```bash
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

### 4. Deploy to AWS (`deploy-aws.yml`)

Deploys the Lambda function to AWS.

**Triggers:**
- Push to `main` branch (auto-deploy to staging)
- Manual workflow dispatch (choose environment)

**Jobs:**
- **deploy**: Builds and deploys Lambda function to AWS
  - Updates function code
  - Updates environment variables
  - Creates deployment tag
  - Sends Slack notification (optional)

### 5. CodeQL Analysis (`codeql.yml`)

Security and quality analysis using GitHub CodeQL.

**Triggers:**
- Push to `main` and `develop`
- Pull requests
- Weekly schedule (Mondays at midnight)

### 6. Dependabot (`dependabot.yml`)

Automated dependency updates for:
- Go modules (weekly)
- GitHub Actions (weekly)

## Required Secrets

Configure these secrets in your GitHub repository settings:

### AWS Deployment
- `AWS_ACCESS_KEY_ID`: AWS access key for deployment
- `AWS_SECRET_ACCESS_KEY`: AWS secret key for deployment
- `AWS_REGION`: AWS region (e.g., us-east-1)
- `LAMBDA_FUNCTION_NAME`: Name of your Lambda function
- `SNS_TOPIC_ARN`: ARN of the SNS topic

### Optional
- `SLACK_WEBHOOK`: Slack webhook URL for deployment notifications
- `LOG_LEVEL`: Lambda log level (default: INFO)

## GitHub Environments

Set up environments in your repository for deployment:

1. **staging**: For staging deployments
2. **production**: For production deployments (with protection rules)

### Environment Protection Rules (Recommended)

For production environment:
- Require reviewers before deployment
- Wait timer (e.g., 5 minutes)
- Restrict to specific branches (main only)

## Setting Up CI/CD

### 1. Configure AWS Credentials

Create an IAM user with the following permissions:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "lambda:UpdateFunctionCode",
        "lambda:UpdateFunctionConfiguration",
        "lambda:GetFunction",
        "lambda:PublishVersion"
      ],
      "Resource": "arn:aws:lambda:REGION:ACCOUNT:function:FUNCTION_NAME"
    }
  ]
}
```

### 2. Add Secrets to GitHub

1. Go to repository Settings → Secrets and variables → Actions
2. Add the required secrets listed above

### 3. Create Environments

1. Go to repository Settings → Environments
2. Create `staging` and `production` environments
3. Add environment-specific secrets if needed
4. Configure protection rules for production

### 4. Enable Workflows

Workflows are automatically enabled when you push them to the repository.

## Manual Deployment

To manually deploy to a specific environment:

1. Go to Actions tab in GitHub
2. Select "Deploy to AWS" workflow
3. Click "Run workflow"
4. Choose environment (staging/production)
5. Click "Run workflow"

## Monitoring Workflow Runs

- View all workflow runs in the Actions tab
- Each workflow run shows detailed logs for each job
- Failed runs will show error messages and stack traces
- Artifacts can be downloaded from successful runs

## Troubleshooting

### Build Failures

Check the build logs for:
- Go compilation errors
- Missing dependencies
- Test failures

### Deployment Failures

Common issues:
- Invalid AWS credentials
- Insufficient IAM permissions
- Lambda function not found
- Network connectivity issues

### Lint Failures

Fix linting issues locally:
```bash
cd aws
golangci-lint run
```

### Test Failures

Run tests locally:
```bash
cd aws
go test -v ./...
```

## Best Practices

1. **Always run tests locally** before pushing
2. **Keep secrets secure** - never commit them
3. **Use semantic versioning** for releases
4. **Review PR checks** before merging
5. **Monitor deployment status** in Slack or GitHub
6. **Use staging environment** for testing before production

## Customization

### Adding New Workflows

1. Create a new `.yml` file in `.github/workflows/`
2. Define triggers, jobs, and steps
3. Test with a pull request
4. Document in this README

### Modifying Existing Workflows

1. Edit the workflow file
2. Test changes in a feature branch
3. Review workflow run results
4. Merge when validated

## Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [AWS Lambda Deployment](https://docs.aws.amazon.com/lambda/latest/dg/welcome.html)
- [golangci-lint](https://golangci-lint.run/)
- [CodeQL](https://codeql.github.com/)
