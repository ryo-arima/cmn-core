# AWS Secrets Manager Integration

This document describes how to use AWS Secrets Manager for configuration management in the cmn-core application.

## Overview

cmn-core supports two configuration sources:
1. **File-based configuration** (default): Loads configuration from `etc/app.yaml`
2. **AWS Secrets Manager**: Loads configuration from AWS Secrets Manager (supports both LocalStack for development and production AWS)

## Configuration

### Environment Variables

The configuration source is controlled by environment variables:

```bash
# Configuration source selection
USE_SECRETSMANAGER=true    # Set to "true" to use Secrets Manager, "false" or unset to use file-based config

# AWS Secrets Manager configuration
SECRET_ID=cmn-core/config/app # The ID of the secret containing application configuration

# LocalStack configuration (for local development)
USE_LOCALSTACK=true        # Set to "true" to use LocalStack, "false" to use production AWS
AWS_ENDPOINT_URL=http://localhost:4566  # LocalStack endpoint URL
AWS_REGION=us-east-1       # AWS region
```

See `.env.example` for a complete example of environment variable configuration.

## Local Development with LocalStack

### Prerequisites

- Docker and Docker Compose installed
- AWS CLI installed (optional, for manual testing)

### Setup

1. **Start LocalStack**:
   ```bash
   docker-compose up -d localstack
   ```

2. **Initialize Secrets Manager** (automatically runs on container start):
   The secrets initialization is handled by `./scripts/main.sh secrets init`, which creates the `cmn-core/config/app` secret with the configuration from `etc/app.dev.yaml`.

3. **Verify secret creation** (optional):
   ```bash
   aws --endpoint-url=http://localhost:4566 secretsmanager list-secrets --region us-east-1
   aws --endpoint-url=http://localhost:4566 secretsmanager get-secret-value --secret-id cmn-core/config/app --region us-east-1
   ```

4. **Configure environment variables**:
   ```bash
   export USE_SECRETSMANAGER=true
   export USE_LOCALSTACK=true
   export SECRET_ID=cmn-core/config/app
   export AWS_ENDPOINT_URL=http://localhost:4566
   export AWS_REGION=us-east-1
   ```

5. **Run the application**:
   ```bash
   go run cmd/server/main.go
   ```

   The application will load configuration from LocalStack Secrets Manager instead of `etc/app.yaml`.

### Manual Secret Update

To update the secret manually:

```bash
# Create a JSON configuration file
cat > /tmp/config.json <<EOF
{
  "application": {
    "common": {
      "name": "cmn-core-updated",
      "version": "1.0.0",
      "env": "dev"
    },
    ...
  }
}
EOF

# Update the secret
aws --endpoint-url=http://localhost:4566 secretsmanager put-secret-value \
  --secret-id cmn-core/config/app \
  --secret-string file:///tmp/config.json \
  --region us-east-1
```

## Production Deployment

### Prerequisites

- AWS account with Secrets Manager access
- IAM credentials or IAM role with `secretsmanager:GetSecretValue` permission

### Setup

1. **Create secret in AWS Secrets Manager**:
   ```bash
   # Create JSON configuration file
   cat > /tmp/prod-config.json <<EOF
   {
     "application": {
       "common": {
         "name": "cmn-core",
         "version": "1.0.0",
         "env": "production"
       },
       "server": {
         "host": "0.0.0.0",
         "port": "8080",
         "jwtSecretKey": "your-production-jwt-secret",
         "jwtExpiresIn": 3600
       },
       ...
     }
   }
   EOF

   # Create the secret
   aws secretsmanager create-secret \
     --name cmn-core/config/app \
     --secret-string file:///tmp/prod-config.json \
     --region us-east-1
   ```

2. **Configure environment variables**:
   ```bash
   export USE_SECRETSMANAGER=true
   export USE_LOCALSTACK=false
   export SECRET_ID=cmn-core/config/app
   export AWS_REGION=us-east-1
   # AWS credentials will be loaded from default credential chain (environment variables, IAM role, etc.)
   ```

3. **Run the application**:
   ```bash
   ./cmn-server
   ```

## Fallback Behavior

If Secrets Manager configuration fails (e.g., network issues, invalid credentials), the application will automatically fall back to file-based configuration (`etc/app.yaml`).

Log messages will indicate which configuration source was used:
- `Successfully loaded configuration from Secrets Manager` - Secrets Manager was used
- `Successfully loaded configuration from file (etc/app.yaml)` - File-based config was used
- `Failed to load config from Secrets Manager: <error>, falling back to file-based config` - Secrets Manager failed, fell back to file

## Configuration Format

The secret value must be a JSON object matching the YAML structure in `etc/app.yaml`:

```json
{
  "application": {
    "common": {
      "name": "cmn-core",
      "version": "1.0.0",
      "env": "dev"
    },
    "server": {
      "host": "localhost",
      "port": "8080",
      "jwtSecretKey": "dev-jwt-secret-key",
      "jwtExpiresIn": 3600
    },
    "client": {
      "serverEndpoint": "http://localhost:8080",
      "userEmail": "user@example.com",
      "userPassword": "password"
    }
  },
  "admin": {
    "emails": ["admin@cmn.local"]
  },
  "mysql": {
    "host": "localhost",
    "user": "root",
    "pass": "password",
    "port": "3306",
    "db": "cmn_core"
  },
  "redis": {
    "host": "localhost",
    "port": "6379",
    "pass": "",
    "db": "0"
  }
}
```

## Troubleshooting

### "Failed to create Secrets Manager client"

- Check that AWS credentials are properly configured
- Verify AWS_REGION is set correctly
- For LocalStack: ensure LocalStack is running and AWS_ENDPOINT_URL is correct

### "Failed to get secret"

- Verify the SECRET_ID matches the secret name in Secrets Manager
- Check IAM permissions include `secretsmanager:GetSecretValue`
- For LocalStack: ensure the init script ran successfully

### "Failed to unmarshal secret as JSON"

- Verify the secret value is valid JSON
- Check that the JSON structure matches the expected configuration format

### Application falls back to file-based config unexpectedly

- Check log messages for specific error details
- Verify USE_SECRETSMANAGER=true is set
- Ensure SECRET_ID is not empty
