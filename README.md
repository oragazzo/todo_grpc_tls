# gRPC TLS Todo Application

This project demonstrates the implementation of a secure gRPC service using Mutual TLS (Transport Layer Security) authentication. The implementation is based on the article by [BB Engfort](https://bbengfort.github.io/2017/03/secure-grpc/).

## Overview

This project showcases how to implement secure communication between a gRPC client and server using TLS certificates. It implements a simple Todo service while focusing on the security aspects of gRPC communication.

## Features

- Secure gRPC communication using TLS
- Client-Server architecture
- Certificate-based authentication
- Todo service implementation
- Environment-based configuration
- Makefile for easy project management

## Prerequisites

- Go 1.23.5 or higher
- Protocol Buffers compiler
- Make
- PostgreSQL database

## Project Structure

```
.
├── cert/               # Directory containing TLS certificates and generation scripts
├── cmd/                # Application entry points
│   ├── client/        # gRPC client implementation
│   │   ├── .env      # Client environment configuration
│   │   └── main.go   # Client entry point
│   └── server/        # gRPC server implementation
│       ├── .env      # Server environment configuration
│       └── main.go   # Server entry point
├── internal/          # Internal packages
├── proto/             # Protocol Buffers definitions
├── Makefile          # Build and run commands
└── README.md         # This file
```

## Environment Configuration

Both the server and client use environment variables for configuration. You can provide these variables in two ways:

1. Using `.env` files
2. Setting environment variables directly

### Server Configuration

Create a `cmd/server/.env` file with the following variables:

```env
# Database Configuration
DSN="host=localhost user=postgres password=postgres dbname=todos port=5432 sslmode=disable"

# Server TLS Configuration
TLS_CERT_FILE=/path/to/server.crt
TLS_KEY_FILE=/path/to/server.key
TLS_CA_FILE=/path/to/ca.crt
TLS_SERVER_NAME=localhost

# gRPC Server Configuration
PORT=8080
```

### Client Configuration

Create a `cmd/client/.env` file with:

```env
# Client TLS Configuration
TLS_CERT_FILE=/path/to/client.crt
TLS_KEY_FILE=/path/to/client.key
TLS_CA_FILE=/path/to/ca.crt
TLS_SERVER_NAME=localhost

# gRPC Client Configuration
SERVER_ADDRESS=localhost:8080
```

### Environment File Location

You can specify a custom location for the `.env` file using the `--env-path` flag:

```bash
# Run server with custom env file
go run cmd/server/main.go --env-path=/path/to/server.env

# Run client with custom env file
go run cmd/client/main.go --env-path=/path/to/client.env
```

The Makefile includes commands that automatically use the correct env files:

```bash
# Run server with default env file
make run_server

# Run client with default env file
make run_client
```

## Getting Started

1. Generate Protocol Buffers code:
```bash
make gen_proto
```

2. Generate TLS certificates:
```bash
make gen_certs
```

3. Set up environment files:
   - Copy `cmd/server/.env.example` to `cmd/server/.env` and adjust values
   - Copy `cmd/client/.env.example` to `cmd/client/.env` and adjust values

4. Run the server:
```bash
make run_server
```

5. In a new terminal, run the client:
```bash
make run_client
```

## Security Implementation

The project implements TLS security using X.509 certificates. Here's how the security is implemented:

1. **Certificate Generation**: A script in the `cert/` directory generates:
   - A Certificate Authority (CA)
   - Server certificates signed by the CA

2. **Server Security**: 
   - The server is configured with its certificate and private key
   - It also has the CA certificate to verify client certificates
   - TLS is enforced for all connections

3. **Client Security**:
   - The client verifies the server's certificate
   - It uses the CA certificate to establish trust
   - Client certificates for mutual TLS authentication

## Understanding TLS in gRPC

gRPC uses TLS for securing communication channels between clients and servers. This implementation demonstrates:

1. **Server Authentication**: The server presents its certificate to prove its identity
2. **Encryption**: All communication is encrypted using TLS
3. **Trust Chain**: Certificates are validated against a trusted Certificate Authority
4. **Secure Connections**: No plaintext communication is allowed

## License

This project is licensed under the MIT License - see the LICENSE file for details.