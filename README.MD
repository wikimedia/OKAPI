# Okapi - Ocean Of Knowledge API

## Technologies:

1. [Golang](https://golang.org/)
2. [Postgresql](https://www.postgresql.org/)
3. [Redis](https://redis.io/)

## Getting started:

1. Create `.env` file in the project root (you need Redis, Postgress server to start the application, alternatively look at docker spin up section below). Here's an example of the file: 
```
REDIS_ADDR=localhost:6379
DB_ADDR=localhost:5432
DB_NAME=okapi
DB_USER=admin
DB_PASSWORD=sql
AWS_REGION=us-east-1
AWS_BUCKET=test
AWS_KEY=test4api12key123
AWS_ID=ABAC123AD13FF
WIKI_USER_AGENT=test@gmail.com
API_PORT=4000
API_MODE=dev
API_AUTH=admin:password,user:password
API_TEST_URL=http://app:3000
VOLUME_MOUNT_PATH=tmp
```
2. Running migrations:

```bash
$ go run migrations/*.go migrate
```

3. Running web server (you have an ability to pass `API_PORT` in `.env` file or by flag `-port` through the command line, otherwise it will use `4000` as a default port http://localhost:4000/):
```bash
$ go run main.go
```

3. Running your first task (this one will populate the database with projectss). You can specify a task name to run and jobs that need to be handled by providing the flags, also you can specify amount of workers to deal with them (this will override default number):
```bash
$ go run main.go -task=projects -workers=10
```

4. Generating your first dump. We need to run 3 tasks in order to get things done. First one would be `scan` - this one gets pages info and metadata into the database. Second one is `sync` this one gets HTML files and puts them into the storage. And final one `bundle` this one generates dump archive and puts into the storage as well. Lets do a small project just to get things rolling (list of projects https://dumps.wikimedia.org/other/pagetitles/20200410/):
```bash
$ go run main.go -task=scan -project=afwikibooks -workers=5
$ go run main.go -task=sync -project=afwikibooks -workers=5
$ go run main.go -task=bundle -project=afwikibooks -workers=5
```

5. Now let's spin up our queue server like this:
```bash
$ go run main.go -server=queue -workers=10
```

6. To start stream server which watchs for pages updates in Wikipedia and updates your versions you can do:
```bash
$ go run main.go -server=stream -workers=10
```

7. Production use. Yo can build and run all of the commands specified from the binary file like this:
```bash
# build the app
$ go build main.go

# and run stream server for example
$ ./main -server=stream -workers=10
```



## Using docker: 

1. Make sure `docker` and `docker-compose` is installed on yor machine.

2. Add `POSTGRES` variables and change other variables in your `.env` file to fit docker network connection criteria, take a closer look at our `REDIS_ADDR` and `DB_ADDR` variables, we change `localhost` to `db` and `cache` to fit docker network. Also make sure your API port is set to `3001`.  Your `.env` file should look something like this:

```
POSTGRES_USER=admin
POSTGRES_PASSWORD=test_password
POSTGRES_DB=okapi
CACHE_ADDR=cache:6379
CACHE_PASSWORD=testpass123
DB_ADDR=db:5432
DB_NAME=okapi
DB_USER=admin
DB_PASSWORD=test_password
AWS_REGION=us-east-1
AWS_BUCKET=test
AWS_KEY=test4api12key123
AWS_ID=ABAC123AD13FF
USER_AGENT=test@gmail.com
API_PORT=3001
API_MODE=dev
API_AUTH=admin:password,user:password
API_TEST_URL=http://app:3000
GRAYLOG_MONGODB_URI=mongodb://log_db/graylog
GRAYLOG_PASSWORD_SECRET=somepasswordpepper
GRAYLOG_ROOT_PASSWORD_SHA2=8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918
GRAYLOG_HTTP_EXTERNAL_URI=http://127.0.0.1:9000/
GRAYLOG_PASSWORD=supersecretpassword
VOLUME_MOUNT_PATH=tmp
```

3. Then we can just run `docker-compose up` and we should be able to run api requests.

4. To run jobs you would need to ssh into `app` container. Here's little cheat sheet how to do that:
```bash
# list containers
$ (sudo) docker container ls

# ssh into container
$ (sudo) docker exec -it <container_id> bash

```

## HTML files storage
To use local file storage for HTML files change local storage interface inside `lib/storage/local.go` inside `localClient()` function to your preferred method of storage:

```go
func localClient() Connection {
	return local.NewStorage() // Change this library
}
```

## API documentation:

We use postman to keep track of API endpoints, postman collection:
https://www.getpostman.com/collections/7bb09db7210e049a1375

## Docker:
We have `main.sh` script, that handles app build process depending on `API_MODE` variable in `.env` file. For more info please take a look at the file itself.

## Packages:

https://github.com/go-pg/pg - ORM package

https://github.com/go-pg/pg/wiki/Writing-Queries - Writing queries with ORM

https://github.com/robinjoseph08/go-pg-migrations - Running migrations

https://github.com/aws/aws-sdk-go - AWS SDK package

https://github.com/go-resty/resty - HTTP request client (used like wrapper for wikipedia API)

https://github.com/gin-gonic/gin - HTTP router for API

https://github.com/go-redis/redis - Redis cache client

https://github.com/joho/godotenv - `.env` support for golang

https://github.com/gookit/event - event manager