# Wikimedia Enterprise

Thi file contains short installation guide for engineers to start working with the system. As this is gRPC server by large you might need client to work with the application, we recommend [BloomRPC](https://github.com/uw-labs/bloomrpc) as fast and easy to use app for development.

## Getting started:
1. Create `.env` file in the project root (you need Redis, Postgress server to start the application, alternatively look at docker spin up section below). Here's an example of the file: 

```bash
# App settings
AWS_REGION=us-east-2
AWS_BUCKET=test-s3-bucket
AWS_KEY=asdasijoj12312ojasdap1p/asd
AWS_ID=FFDWQIRJ123
REDIS_ADDR=cache:6379
REDIS_PASSWORD=J7Bcxp6kjQjVbGG7
DB_ADDR=db:5432
DB_NAME=wikimedia_enterprise
DB_USER=admin
DB_PASSWORD=sql
ELASTIC_URL=http://search:9200
HTML_VOL=tmp/html
GEN_VOL=tmp/gen
JSON_VOL=tmp/json
WT_VOL=tmp/wt
```
2. Running migrations:

```bash
go run migrations/*.go migrate
```

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
# Queue processors
go run queues/main.go
```

```bash
# Steams listener
go run streams/main.go
```

5. Updating gRPC server. If you need to make changes to `.proto` files don't forget to push you changes into `/protos` git submodule and re-generate static files by running (make sure [protoc compiler](https://grpc.io/docs/protoc-installation/) is installed on you machine):
```bash
make protos
```

## Using docker: 

1. Make sure `docker` and `docker-compose` is installed on yor machine.

2. Add docker configuration variables to `.env` file, here's an example of complete `.env` file with `app` and `docker` configuration:

```bash
# App settings
AWS_REGION=us-east-2
AWS_BUCKET=test-s3-bucket
AWS_KEY=asdasijoj12312ojasdap1p/asd
AWS_ID=FFDWQIRJ123
REDIS_ADDR=cache:6379
REDIS_PASSWORD=J7Bcxp6kjQjVbGG7
DB_ADDR=db:5432
DB_NAME=wikimedia_enterprise
DB_USER=admin
DB_PASSWORD=sql
ELASTIC_URL=http://search:9200

# Docker settings
POSTGRES_USER=admin
POSTGRES_PASSWORD=sql
POSTGRES_DB=wikimedia_enterprise
```

3. After that's done just run:
```bash
sudo docker-compose up
```

4. In order to `ssh` into the container you can do:
```bash
sudo docker-compose exec server bash
```
