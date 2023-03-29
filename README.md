# geospace

**geospace** geo service api for calculating distances between cities and much more...

## Installation

```
go get github.com/alaleks/geospace
```

## Configuration of server

### Required Options
```
-d Name of the database
-u User name of the database
-p Password of the database
-s Socket of connection to database
-t Port of the database
```
### Optional Options
```
-n Name of the database
-a Port for running the application
-r Max request quantity in seconds
-e Expiration period in seconds
```

### Example run server
```
 go run main.go -d=db_name -u=db_user -p=password -s=unix_socket -r=100

```
## Build client

```
go build -ldflags "-X main.Version=v1 -X main.Host=:3000 -X main.Name=geo" main.go

where:

Version is the version number app
Name is the name of the app
Host is the host to connect to app
```

## Run client

```
go run -ldflags "-X main.Version=v1 -X main.Host=:3000 -X main.Name=geo" main.go 
```