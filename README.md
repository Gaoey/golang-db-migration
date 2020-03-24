# GOLANG Dabase Migration 
KTB Golang Version Control for Database Schemas
using golang-migrate
see more. https://github.com/golang-migrate/migrate

tutorial by mysql
# Docker Compose
``` 
$ docker-compose build --no-cache
$ docker-compuse up -d
```

# Manual Step
## Pre-setup
1. install docker
2. pull mysql image
3. set username & pass word for mysql container as 
```
username: root
password: password
```
4. create table "todolist"

## Project base
1. migrate up
```go run main.go -step up```
2. migrate down
```go run main.go -step down```
3. force version
```go run main.go -force <version>```


## CLI
1. install migrate CLI see here: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
2. try to create database schemas for up and down (up for new schema to upgrade your database, down for rollback when found problem while migration)
```
migrate create -ext sql create_todo_table
```
3. try to migrate up, down by command
```
migrate -path ./migration -database mysql://root:password@/todolist up
```
```
migrate -path ./migration -database mysql://root:password@/todolist down
```
4. WARNING: if show Dirty status -> try to force version 
example
```
migrate -path migration/ -database mysql://root:password@/todolist force 1
```