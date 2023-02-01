# Wikimedia Enterprise Realtime API

This is short installation guide to okapi realtime api.

## Getting started:
1. Create `.env` file in the project root. Here's an example of the file:

```bash
API_MODE=debug
API_PORT=4040
KAFKA_BROKER=localhost
BASIC_AUTH=admin:sql
AWS_AUTH_REGION=ap-northeast-1
AWS_AUTH_KEY=password
AWS_AUTH_ID=admin
IP_RANGE=
COGNITO_CLIENT_ID=client_id
REDIS_ADDR=cache:6379
```

3. Testing, in order to run test suite yo need to do:
```bash
go test ./... -v
```

4. Running the application:
```bash
go run main.go
```

## To update/test RBAC:

Use [this](https://github.com/prabhat393/rbac-example) for prototyping. Once you have a working `model.conf` and `policy.csv`, replace the ones in the base folder with your new model/policy. Currently, we are using [RBAC with transitive user roles](https://github.com/casbin/casbin/blob/master/examples/rbac_with_hierarchy_policy.csv).
