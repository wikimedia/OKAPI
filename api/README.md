# Wikimedia Enterprise API

## Getting stated

1. Create `.env` file in the project root (for more info look inside `env/env.go`). Here's an example:
```bash
# App settings
API_PORT=3000
API_MODE=debug
AWS_URL=http://minio:9000
AWS_REGION=ap-northeast-1
AWS_BUCKET=wme-data
AWS_KEY=password
AWS_ID=admin
REDIS_ADDR=cache:6379
REDIS_PASSWORD=J7Bcxp6kjQjVbGG7
PROJECTS_EXPIRE=86400
PAGES_EXPIRE=3600
IP_RANGE=
IP_RANGE_REQUESTS_LIMIT=10
AWS_AUTH_REGION=ap-northeast-1
AWS_AUTH_KEY=password
AWS_AUTH_ID=admin
GROUP=group_1

# Docker settings
MINIO_ROOT_USER=admin
MINIO_ROOT_PASSWORD=password
```

2. Run the app:
```bash
$ go run main.go
```

## Deploying the app

1. Build the application with a following command:
```bash
$ go build *.go
```

2. Don't forget to switch the application in release mode inside the env file:
```bash
API_MODE=release
# ... your other variables 
```

## Working with documentation

1. If you are running linux/unix and made changes to the docs you can re-generate the app docs by running:
```bash
$ make swag
```

2. To access the docs you need to go to [http://localhost:8080/v1/docs/index.html](http://localhost:8080/v1/docs/index.html) or another port if you've changed the `.env` file.


## Testing

1. To run all tests please use:
```bash
$ go test ./... -v
```

2. To run a single test (where `TestInit` is the name of the test and `./env/` is path to the module):
```bash
$ go test ./env/ -run TestInit -v
```

## Using docker

1. Running the app:
```bash
$ (sudo) docker-compose up #optionally you can add -d flag
```

## To update/test RBAC

Use [this](https://github.com/prabhat393/rbac-example) for prototyping. Once you have a working `model.conf` and `policy.csv`, replace the ones in the base folder with your new model/policy. Currently, we are using [RBAC with transitive user roles](https://github.com/casbin/casbin/blob/master/examples/rbac_with_hierarchy_policy.csv).
