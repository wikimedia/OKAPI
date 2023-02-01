# Wikimedia Enterprise

This file contains short installation guide for engineers to start working with the system. As this is gRPC server by large you might need client to work with the application, we recommend [BloomRPC](https://github.com/uw-labs/bloomrpc) as fast and easy to use app for development.

## Getting started:
1. Create `.env` file in the project root (you need Redis, PostgreSQL, Kafka etc. to start the application, alternatively look at docker spin up section below). Here's an example of the file: 

```bash
# App settings
AWS_URL=default
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
ELASTIC_USERNAME=admin
ELASTIC_PASSWORD=sql
GEN_VOL=/var/www/app/.vol
JSON_VOL=/var/www/app/.vol
KAFKA_BROKER=kafka:19092
```

3. Running migrations:

```bash
go run migrations/*.go migrate
```

4. Testing, in order to run test suite yo need to do:
```bash
go test ./... -v
```

5. Running the application:
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
AWS_URL=http://minio:9000
AWS_REGION=ap-northeast-1
AWS_BUCKET=wme-data
AWS_KEY=password
AWS_ID=admin
REDIS_ADDR=cache:6379
REDIS_PASSWORD=J7Bcxp6kjQjVbGG7
DB_ADDR=db:5432
DB_NAME=wikimedia_enterprise
DB_USER=admin
DB_PASSWORD=sql
ELASTIC_URL=http://search:9200
ELASTIC_USERNAME=admin
ELASTIC_PASSWORD=sql
GEN_VOL=/root/.vol
JSON_VOL=/root/.vol
KAFKA_BROKER=kafka:19092

# Docker settings
POSTGRES_USER=admin
POSTGRES_PASSWORD=sql
POSTGRES_DB=wikimedia_enterprise
MINIO_ROOT_USER=admin
MINIO_ROOT_PASSWORD=password
```

3. After that's done just run:
```bash
sudo docker-compose up
```

4. Create a bucket called `wme-data` using minio console. It should be available on [http://localhost:9202/](http://localhost:9202/). Username and password are from your `.env` file. Not that if you don't see bucket in the UI after it was created take a look inside `.storage` directory inside the project. If you see folder called `wme-data` inside than everything is fine.

5. In order to `ssh` into the container you can do:
```bash
sudo docker-compose exec server bash
```

## Getting baseline data:

1. Import proto files from `/protos` directory into the application. Specifically `pages.proto`, `namespaces.proto`, `projects.proto`, `search.proto`.

2. Then using the [BloomRPC](https://github.com/uw-labs/bloomrpc) or any other gRPC client. Execute following steps:
    
    * Run `projects.Fetch` method to collect the initial projects list.

    * Run `projects.Index` method to index the projects in elasticsearch.

    * Run `namespaces.Fetch` method to collect namespaces across the projects.

3. After all of that is done. You now can generate your first `export`. Like this:
    
    * Run `pages.Fetch` with particular database name (we recommend `afwikibooks` in namespace `0` as it's small and fast to execute). This is needed to collect initial dataset.

    * Run `pages.Export` with the same database name and namespace again to create the `export` file in one particular namespace.

    * Run `projects.Aggregate` to generate and export global metadata of projects for each namespace.

4. After exporting dumps and metadata, you are ready to create a `copy` for custom users.

    * Run `pages.Copy` with a list of projects and a namespace (e.g., db_names `["afwikibooks"]` and ns `0`). This will create copies of `afwikibooks` tar and metadata for namespace 0, as well as copy of global exports metadata for namespace 0.

