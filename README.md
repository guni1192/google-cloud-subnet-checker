# google-cloud-subnet-checker

A command-line tool to check for CIDR overlap before creating a new subnet in Google Cloud Platform.

## Features

- Retrieves existing subnet CIDR ranges in the specified region
- Checks both primary and secondary IP ranges for overlap
- Validates CIDR format

## Installation

```
go install github.com/guni1192/google-cloud-subnet-checker@latest
```

### Options

- `--project`: (Required) Google Cloud Project ID
- `--region`: (Required) Region name (e.g., us-central1)
- `--cidr`: (Required) Desired CIDR range for the new subnet (e.g., 192.168.1.0/24)
- `--debug`: Enable debug logging

### Example

```console
google-cloud-subnet-checker \
--project=my-project \
--region=asia-northeast1 \
--cidr=192.168.1.0/24
```
