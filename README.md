# geospace

**geospace** service with client for calculating distances between cities and much more...

## Installation

```
go get github.com/alaleks/geospace
```

## Stack

### Server 

- MariaDB (Database)
- Fiber (Web Framework)
- SQLX (Library which provides using database sql)
- Zap (Logger)

### Client

- Fiber (Web Framework)
- Pterm (Printer terminal)

## Configuration of server

### Required Options

-d Name of the database

-u User name of the database

-p Password of the database

-s Socket of connection to database

-t Port of the database

### Optional Options

-n Name of the database

-a Port for running the application

-r Max request quantity in seconds

-e Expiration period in seconds


### First run

When the server is first started, the connection to the database is checked. If the connection to the database is successful, a configuration file "config.yaml" is created in the "сfg" folder in the root directory of the project. If it was not possible to create a folder and write a file, then the settings are valid only in current session. It also creates table schemas in the database and imports the necessary data.

### Change configuration parameters

If a configuration file was created, then if you need to change settings, you need to make changes in it, and not through flags. Or you can delete the configuration file and start the server with the configuration flags.

### Example run server

```
 go run main.go -d=db_name -u=db_user -p=password -s=unix_socket -r=100
```

## Build client

```
go build -ldflags "-X main.Version=v1 -X main.Host=:3000 -X main.Name=geo" main.go
```

Where:

- Version is the version number app
- Name is the name of the app
- Host is the host to connect to app

## Run client

```
go run -ldflags "-X main.Version=v1 -X main.Host=:3000 -X main.Name=geo" main.go 
```

## Methods

 - /ping - check server health. If server is healthy return 200.
 - /v1/register - provides sign up. 
 
If is registered successfully returned 200:

```
{
    "token": "string"
}
```

Token need transfer to:
- Header as parameter Authorization in format [Bearer token]
- Cookie access_token

 ```
 POST application/json

 {
    "name":"string",
    "email": "string",
    "password": "string",
}
 ```
- /v1/login - provides log in. 

If is log in successfully returned 200:

```
{
    "token": "string"
}
```

Token need transfer to:
- Header as parameter Authorization in format [Bearer token]
- Cookie access_token

 ```
 POST application/json

{
    "email": "string",
    "password": "string"
}
```
- /v1/logout - provides log out.

### Client

- /v1/user/distance - provides calculate distance between two points by coordinates.

 ```
GET application/json

http --follow --timeout 3600 GET 'http://localhost:3000/v1/user/distance?departure=Moscow,Ru&destination=Краснодар&Bearer=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODAzNDM4ODksInVpZCI6MH0.L_kAIH_8FmrIAwGEkJ4CZU13QOsollvH9Xebufjxfxw' \
```
Where:

- departure - city of departure
- destination - city of destination

- /v1/user/find-by-name - provides find nearest points (cities) by departure city by name

```
http --follow --timeout 3600 GET 'http://localhost:3000/v1/user/find-by-name?departure=Rome, It&distanceTo=100' \
 Authorization:'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODA5NTQ5NDgsInVpZCI6OH0.rwwfDiJKTf1yzxlG3-JeKvmYu6PVb3HHuVOCguYCNhY'
```

Where:

- departure - city of departure
- distanceTo - in what radius (at what distance) to look for cities in km

- /v1/user/find-by-coord - provides find nearest points (cities) from passed coordinates

```
http --follow --timeout 3600 GET 'http://localhost:3000/v1/user/find-by-coord?lat=43.77925&lon=11.24626&distanceTo=50' \
 Authorization:'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODA5NTQ5NDgsInVpZCI6OH0.rwwfDiJKTf1yzxlG3-JeKvmYu6PVb3HHuVOCguYCNhY'
```

Where:

- lat - latitude of point
- lon - longitude of point
- distanceTo - in what radius (at what distance) to look for cities in km

### Api
- /v1/api/distance - provides calculate distance between two points by coordinates
```
http --follow --timeout 3600 GET 'http://localhost:3000/v1/api/distance?destination=Venice,Italy&departure=Rome,It' \
Authorization:'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODA5NTQ5NDgsInVpZCI6OH0.rwwfDiJKTf1yzxlG3-JeKvmYu6PVb3HHuVOCguYCNhY'
```

Where:

- departure - city of departure
- destination - city of destination

Response Error:
```
{
    "message": "string",
    "code": int
}
```

Response Ok:
```
{
    "departure": {
        "name": "string", // city name
        "name_ascii": "string", // city name in ASCII
        "country_code": "string", // country code
        "country": "Italy", // country name
        "timezone": "string", // time zone
        "city_id": int, // city id in database
        "latitude": float, 
        "longitude": float,
    },
    "destination": {
        "name": "string", // city name
        "name_ascii": "string", // city name in ASCII
        "country_code": "string", // country code
        "country": "Italy", // country name
        "timezone": "string", // time zone
        "city_id": int, // city id in database
        "latitude": float, 
        "longitude": float,
    },
    "distance_straight": int, // distance between two city in straight line
    "distance_road":int // distance between two city by road
}
```

- /v1/api/find-by-name - provides find nearest points (cities) by departure city by name

```
http --follow --timeout 3600 GET 'http://localhost:3000/v1/api/find-by-name?departure=Rome, It&distanceTo=100' \
Authorization:'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODA5NTQ5NDgsInVpZCI6OH0.rwwfDiJKTf1yzxlG3-JeKvmYu6PVb3HHuVOCguYCNhY'
 ```

 Where:

- departure - city of departure
- distanceTo - in what radius (at what distance) to look for cities in km

Response Error:
```
{
    "message": "string",
    "code": int
}
```

Response Ok:
```
{
   "cities_nearby": [
        {
            "name": "string", // city name
            "name_ascii": "string", // city name in ASCII
            "country_code": "string", // country code
            "country": "Italy", // country name
            "timezone": "string", // time zone
            "city_id": int, // city id in database
            "latitude": float, 
            "longitude": float,
        },
        {
            "name": "string", // city name
            "name_ascii": "string", // city name in ASCII
            "country_code": "string", // country code
            "country": "Italy", // country name
            "timezone": "string", // time zone
            "city_id": int, // city id in database
            "latitude": float, 
            "longitude": float,
        }
   ]
    "departure": {
        "name": "string", // city name
        "name_ascii": "string", // city name in ASCII
        "country_code": "string", // country code
        "country": "Italy", // country name
        "timezone": "string", // time zone
        "city_id": int, // city id in database
        "latitude": float, 
        "longitude": float,
    },
    "distance_to": int, // in what radius (at what distance) to look for cities in km
    "qty_nearby": int // number of cities in response
}
```

- /v1/api/find-by-coord - provides find nearest points (cities) from passed coordinates 

```
http --follow --timeout 3600 GET 'http://localhost:3000/v1/api/find-by-coord?lat=43.77925&lon=11.24626&distanceTo=50' \
Authorization:'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODA5NTQ5NDgsInVpZCI6OH0.rwwfDiJKTf1yzxlG3-JeKvmYu6PVb3HHuVOCguYCNhY'
```

Where:

- lat - latitude of point
- lon - longitude of point
- distanceTo - in what radius (at what distance) to look for cities in km

Response Error:
```
{
    "message": "string",
    "code": int
}
```

Response Ok:
```
{
   "cities_nearby": [
        {
            "name": "string", // city name
            "name_ascii": "string", // city name in ASCII
            "country_code": "string", // country code
            "country": "Italy", // country name
            "timezone": "string", // time zone
            "city_id": int, // city id in database
            "latitude": float, 
            "longitude": float,
        },
        {
            "name": "string", // city name
            "name_ascii": "string", // city name in ASCII
            "country_code": "string", // country code
            "country": "Italy", // country name
            "timezone": "string", // time zone
            "city_id": int, // city id in database
            "latitude": float, 
            "longitude": float,
        }
   ]
    "distance_to": int, // in what radius (at what distance) to look for cities in km
    "qty_nearby": int // number of cities in response
}
```