# micro-app

### To Create ca-certificate

```ssh
openssl req -x509 -nodes -newkey rsa:4096 \
-keyout server.key -out server.crt -days 365 \
-subj "/C=US/ST=State/L=City/O=Organization/OU=Department/CN=localhost" \
-addext "subjectAltName=DNS:localhost,DNS:example.com,IP:127.0.0.1,IP:192.168.1.1" \
-config /dev/null
```

### gRPC Request
```curl
    grpcurl -d '{"movie_id":"1"}' localhost:8091 MetadataService/GetMetadata
```
