# Wikimedia Enterprise Diffs Service

Thi file contains short installation guide for engineers to start working with the system. As this is gRPC server by large you might need client to work with the application, we recommend [BloomRPC](https://github.com/uw-labs/bloomrpc) as fast and easy to use app for development.

## Getting started:
1. Create `.env` file in the project root (you need Kafka broker to start the application, alternatively look at docker spin up section below). Here's an example of the file:

```bash
# App settings
AWS_URL=http://minio:9000
AWS_REGION=ap-northeast-1
AWS_BUCKET=wme-data
AWS_KEY=password
AWS_ID=admin
VOL=/root/.vol
KAFKA_BROKER=localhost

# Docker settings
MINIO_ROOT_USER=admin
MINIO_ROOT_PASSWORD=password
```
2. Create a bucket called `wme-data` using minio console. It should be available on [http://localhost:9202/](http://localhost:9202/). Username and password are from your `.env` file. Not that if you don't see bucket in the UI after it was created take a look inside `.storage` directory inside the project. If you see folder called `wme-data` inside than everything is fine.

3. Testing, in order to run test suite yo need to do:
```bash
go test ./... -v
```

4. Running the application:
```bash
# Main app (gRPC server)
go run server/main.go
```

```bash
# Event processors
go run listener/main.go
```

5. Updating gRPC server. If you need to make changes to `.proto` files don't forget to push you changes into `/protos` git submodule and re-generate static files by running (make sure [protoc compiler](https://grpc.io/docs/protoc-installation/) is installed on you machine):
```bash
make protos
```

## Using docker: 

1. Make sure `docker` and `docker-compose` is installed on yor machine.

2. Then just run (note you still need to set `.env` variables for this to work):
```bash
sudo docker-compose up
```
