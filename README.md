# Akeneo Serverless Connector

A multi-cloud serverless webhook connector for Akeneo PIM that receives webhook events and publishes them to message queues for downstream processing.

## Overview

This project provides cloud-native, serverless solutions for handling Akeneo PIM webhook events across multiple cloud providers. Each implementation follows cloud-specific best practices while maintaining a consistent event processing model.

## Supported Cloud Providers

| Provider | Status | Documentation |
|----------|--------|---------------|
| AWS | 🚧 WIP | [AWS Documentation](./aws/README.md) |
| Azure | 🚧 Coming Soon | [Azure Documentation](./azure/README.md) |
| Google Cloud | 🚧 Coming Soon | [GCP Documentation](./google_cloud/README.md) |

## Architecture

Each cloud implementation follows a similar pattern:

```
Akeneo PIM → API Gateway/HTTP Trigger → Serverless Function → Message Queue → Downstream Services
```

### AWS Implementation
- **Function**: AWS Lambda (Go)
- **API**: API Gateway
- **Queue**: SNS Topic
- **Status**: Production Ready

### Azure Implementation (Planned)
- **Function**: Azure Functions
- **API**: HTTP Trigger
- **Queue**: Azure Service Bus / Event Grid
- **Status**: Coming Soon

### Google Cloud Implementation (Planned)
- **Function**: Cloud Functions
- **API**: Cloud Functions HTTP Trigger
- **Queue**: Pub/Sub
- **Status**: Coming Soon

## Features

- **Multi-Cloud Support**: Deploy to AWS, Azure, or Google Cloud
- **Event Validation**: Ensures payload integrity before publishing
- **Retry Logic**: Automatic retry with exponential backoff
- **Structured Logging**: Cloud-native logging for each provider
- **Error Handling**: Comprehensive error handling with detailed codes
- **Scalability**: Auto-scaling based on load

## Project Structure

```
.
├── aws/                    # AWS Lambda implementation (Go)
│   ├── cmd/
│   ├── internal/
│   ├── go.mod
│   └── README.md
├── azure/                  # Azure Functions implementation (coming soon)
│   └── README.md
├── google_cloud/           # Google Cloud Functions implementation (coming soon)
│   └── README.md
└── README.md              # This file
```

## Quick Start

Choose your cloud provider and follow the specific documentation:

### AWS
```bash
cd aws
go mod download
GOOS=linux GOARCH=amd64 go build -o bootstrap cmd/lambda/main.go
# See aws/README.md for deployment instructions
```

### Azure
Coming soon - Azure Functions implementation

### Google Cloud
Coming soon - Cloud Functions implementation

## Event Format

All implementations accept and process Akeneo webhook events in the following format:

### Input (Akeneo Webhook)

```json
{
  "event_id": "unique-event-id",
  "event_type": "product.updated",
  "timestamp": "2024-10-01T10:00:00Z",
  "author": "username",
  "data": {
    "sku": "PRODUCT-SKU",
    "additional": "fields"
  }
}
```

### Required Fields

- `event_id`: Unique identifier for the event
- `event_type`: Type of event (e.g., product.updated, product.created)
- `timestamp`: ISO 8601 timestamp of when the event occurred

## API Responses

### Success (202 Accepted)

```json
{
  "status": "accepted",
  "event_id": "unique-event-id"
}
```

### Error (4xx/5xx)

```json
{
  "error_code": "INVALID_PAYLOAD",
  "message": "Invalid webhook payload",
  "details": "specific error details"
}
```

### Common Error Codes

- `INVALID_PAYLOAD`: Malformed JSON or invalid payload structure
- `MISSING_FIELD`: Required field is missing from the payload
- `PUBLISH_FAILED`: Failed to publish event to message queue

## Choosing a Cloud Provider

Consider these factors when selecting a cloud provider:

### AWS
- **Best for**: Organizations already using AWS services
- **Pros**: Mature ecosystem, extensive documentation, SNS/SQS integration
- **Language**: Go (high performance, low cold start)

### Azure
- **Best for**: Organizations using Microsoft ecosystem
- **Pros**: Seamless integration with Azure services, Event Grid
- **Language**: TBD (C#, Node.js, or Python)

### Google Cloud
- **Best for**: Organizations using Google Cloud Platform
- **Pros**: Simple deployment, Pub/Sub integration, competitive pricing
- **Language**: TBD (Go, Node.js, or Python)

## CI/CD

This project includes comprehensive GitHub Actions workflows for:

- **Continuous Integration**: Automated testing, linting, and security scanning
- **Pull Request Checks**: PR validation, size checks, and coverage reports
- **Automated Releases**: Multi-architecture builds with changelog generation
- **Continuous Deployment**: Automated deployment to AWS Lambda
- **Security**: CodeQL analysis and dependency updates via Dependabot

See [CI/CD Documentation](.github/workflows/README.md) for detailed setup instructions.

## Development Roadmap

- [x] AWS Lambda implementation (Go)
- [x] CI/CD pipeline with GitHub Actions
- [ ] Azure Functions implementation
- [ ] Google Cloud Functions implementation
- [ ] Dead letter queue support across all providers
- [ ] Event replay mechanism
- [ ] Monitoring and alerting templates
- [ ] Terraform/IaC modules for each provider

## Contributing

Contributions are welcome! Whether you want to:
- Implement Azure or Google Cloud versions
- Improve existing implementations
- Add tests or documentation
- Report bugs or suggest features

Please feel free to submit a Pull Request or open an issue.

### Development Guidelines

1. Each cloud implementation should be self-contained in its directory
2. Follow cloud-specific best practices and idioms
3. Maintain consistent event format across implementations
4. Include comprehensive tests
5. Document deployment and configuration steps

## Testing

Each implementation includes its own test suite. See provider-specific documentation:

- [AWS Testing Guide](./aws/README.md#testing)
- Azure Testing Guide (coming soon)
- Google Cloud Testing Guide (coming soon)

## Security

- Never commit credentials or secrets to the repository
- Use cloud-native secret management (AWS Secrets Manager, Azure Key Vault, GCP Secret Manager)
- Follow principle of least privilege for IAM/permissions
- Enable encryption at rest and in transit
- Regularly update dependencies

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

quimnieto

## Support

For issues and questions:
- Open an issue on GitHub
- Check provider-specific documentation
- Review cloud provider documentation for platform-specific issues

## Acknowledgments

Built to integrate Akeneo PIM with modern cloud-native architectures across multiple cloud providers.
