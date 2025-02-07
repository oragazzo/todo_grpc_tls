#!/bin/bash

# Create config file for the CA
cat > ./cert/ca.conf << EOF
[req]
distinguished_name = req_distinguished_name
x509_extensions = v3_ca
prompt = no

[req_distinguished_name]
C = BR
ST = SP
L = SP
O = Example-CA
CN = example-ca

[v3_ca]
basicConstraints = critical,CA:TRUE
keyUsage = critical,digitalSignature,keyEncipherment,keyCertSign
EOF

# Create config file for the server
cat > ./cert/server.conf << EOF
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_req
prompt = no

[req_distinguished_name]
C = BR
ST = SP
L = SP
O = Example-Server
CN = localhost

[v3_req]
basicConstraints = CA:FALSE
keyUsage = digitalSignature,keyEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
IP.1 = 127.0.0.1
EOF

# Create config file for the client
cat > ./cert/client.conf << EOF
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_req
prompt = no

[req_distinguished_name]
C = BR
ST = SP
L = SP
O = Example-Client
CN = todo-client

[v3_req]
basicConstraints = CA:FALSE
keyUsage = digitalSignature,keyEncipherment
extendedKeyUsage = clientAuth
EOF

# Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ./cert/ca.key -out ./cert/ca.crt -config ./cert/ca.conf -extensions v3_ca

echo "Generated CA certificate and private key"

# Generate server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout ./cert/server.key -out ./cert/server.csr -config ./cert/server.conf

# Generate client's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout ./cert/client.key -out ./cert/client.csr -config ./cert/client.conf

# Use CA's private key to sign server's CSR and get back the signed certificate
openssl x509 -req -in ./cert/server.csr -days 365 -CA ./cert/ca.crt -CAkey ./cert/ca.key -CAcreateserial -out ./cert/server.crt -extfile ./cert/server.conf -extensions v3_req

echo "Generated server certificate and private key"

# Use CA's private key to sign client's CSR and get back the signed certificate
openssl x509 -req -in ./cert/client.csr -days 365 -CA ./cert/ca.crt -CAkey ./cert/ca.key -CAcreateserial -out ./cert/client.crt -extfile ./cert/client.conf -extensions v3_req

echo "Generated client certificate and private key"

# Clean up
rm ./cert/server.csr ./cert/client.csr ./cert/ca.key ./cert/ca.srl ./cert/ca.conf ./cert/server.conf ./cert/client.conf

echo "Cleaned up temporary files"
echo "Generated files in ./cert/: ca.crt, server.key, server.crt, client.key, client.crt"
