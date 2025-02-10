# gRPC TLS Todo Application

This project demonstrates the implementation of a secure gRPC service using Mutual TLS (Transport Layer Security) authentication. The implementation is based on the article by [BB Engfort](https://bbengfort.github.io/2017/03/secure-grpc/).

## Overview

This project showcases how to implement secure communication between a gRPC client and server using TLS certificates. It implements a simple Todo service while focusing on the security aspects of gRPC communication.

## Features

- Secure gRPC communication using TLS
- Client-Server architecture
- Certificate-based authentication
- Todo service implementation
- Makefile for easy project management

## Prerequisites

- Go 1.23.5 or higher
- Protocol Buffers compiler
- Make

## Project Structure

```
.
├── cert/               # Directory containing TLS certificates and generation scripts
├── cmd/                # Application entry points
│   ├── client/        # gRPC client implementation
│   └── server/        # gRPC server implementation
├── proto/              # Protocol Buffers definitions
├── Makefile           # Build and run commands
└── README.md          # This file
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

3. Run the server:
```bash
make run_server
```

4. In a new terminal, run the client:
```bash
make run_client
```

## Security Implementation

The project implements TLS security using X.509 certificates. Here's how the security is implemented:

1. **Certificate Generation**: A script in the `cert/` directory generates:
   - A Certificate Authority (CA)
   - Server certificates signed by the CA
   - Client certificates (when needed)

2. **Server Security**: 
   - The server is configured with its certificate and private key
   - It also has the CA certificate to verify client certificates
   - TLS is enforced for all connections

3. **Client Security**:
   - The client verifies the server's certificate
   - It uses the CA certificate to establish trust
   - (Optional) Client certificates for mutual TLS authentication

## Understanding TLS in gRPC

gRPC uses TLS for securing communication channels between clients and servers. This implementation demonstrates:

1. **Server Authentication**: The server presents its certificate to prove its identity
2. **Encryption**: All communication is encrypted using TLS
3. **Trust Chain**: Certificates are validated against a trusted Certificate Authority
4. **Secure Connections**: No plaintext communication is allowed

## License

This project is licensed under the MIT License - see the LICENSE file for details.