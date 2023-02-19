http-server
=====

This project implemented a HTTP/SQLite server. The server maintains a request multiplexer to process incoming requests.
For each PUT/GET request, a corresponding store/query operation is implemented using the SQLite database.


## Usage 

### Server
To build the project:

```sh
go run main.go
```

The server will listen on localhost:8080/

### Client

Insert a row:
```sh
curl -X PUT http://localhost:8080 -H 'Content-Type: application/json' -d '{"key": "mykey", "value": "myvalue", "timestamp" : 1673524092123456}'
```

Fetch a row:
```sh
curl -X GET http://localhost:8080 -H 'Content-Type: application/json' -d '{"key":"mykey", "timestamp": 1673524092123456}'
```
