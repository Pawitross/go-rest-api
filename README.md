# Book Management REST API in Go

A RESTful API written in Go using Gin framework and MariaDB database, designed for managing books, authors, genres and languages.
Supports JWT authentication, includes OpenAPI (Swagger) documentation and provides load testing with Grafana k6.

## Getting started

### Prerequisites
 - Go (1.24 or higher)
 - MariaDB (11.4.5 or higher)
 - Docker (optional, if you want to run the database in a container)
 - Grafana k6 (optional, for load testing)

### Clone the repository

```sh
git clone --depth=1 https://github.com/foo/bar
cd bar
```

### Usage

Copy the initial server configuration file:
```sh
cp env.yaml.initial env.yaml
```

Start MariaDB using Docker (via [docker-compose.yaml](/docker-compose.yaml)):
```sh
docker compose up
```
Or install MariaDB locally and run scripts from the [sql directory](/sql) to initialize the database.

Start the server (by default available at `localhost:8080`):
```sh
go run cmd/api/main.go
```

Open a web browser and navigate to `http://localhost:8080/swagger/index.html` to access Swagger docs.

## Configuration

The default [`env.yaml.initial` file](/env.yaml.initial) provides a minimal setup to run the API.\
By default, the database is initialized with two users:
 - `user` - has full access to the `paw` database
 - `user_test` - has full access to `paw_test` database used for testing

To start the server, you must configure these three keys: `DBUSER`, `DBNAME`, `SECRET`.\
There are three ways to configure the server:
 - via the `env.yaml` file
 - environment variables
 - CLI flags

Additional configuration options are listed below:
| Config key / environment variable | Description             | Default value |
| --------------------------------- | ----------------------- | ------------- |
| **`DBUSER`**                      | Database user           | -             |
| `DBPASS`                          | Database user password  | empty         |
| **`DBNAME`**                      | Database name           | -             |
| `DBHOST`                          | Database host address   | `127.0.0.1`   |
| `DBPORT`                          | Database port           | `3306`        |
| **`SECRET`**                      | JWT token secret        | -             |

The server can be configured using CLI flags, the `env.yaml` config file or environment variables:
| CLI flag  | Config key / environment variable | Description                               | Default value                  |
| --------- | --------------------------------- | ----------------------------------------- | ------------------------------ |
| `--https` | `HTTPS`                           | Use HTTPS to run the server               | `false`                        |
| `--port`  | `PORT`                            | Server serving port                       | `8080` (HTTP) / `8443` (HTTPS) |
| `--cert`  | `TLS_CERT`                        | TLS certificate file location (for HTTPS) | `keys/server.pem`              |
| `--key`   | `TLS_KEY`                         | TLS private key file location (for HTTPS) | `keys/server.key`              |

## Documentation

The API documentation is available in the [docs directory](/docs) in [JSON](/docs/swagger.json) and [YAML](/docs/swagger.yaml) formats.\
The documentation is also available at [`http://localhost:8080/swagger/index.html`](`http://localhost:8080/swagger/index.html`) after starting the server.

## API Endpoints

All endpoints are available under the `/api/v1` prefix.
```
/api/v1
 ├── /books
 │    ├── GET, POST, OPTIONS
 │    └── /:id  GET, PUT, PATCH, DELETE, OPTIONS
 ├── /authors
 │    ├── GET, POST, OPTIONS
 │    └── /:id  GET, PUT, PATCH, DELETE, OPTIONS
 ├── /genres
 │    ├── GET, POST, OPTIONS
 │    └── /:id  GET, PUT, DELETE, OPTIONS
 ├── /languages
 │    ├── GET, POST, OPTIONS
 │    └── /:id  GET, PUT, DELETE, OPTIONS
 └── /login
      └── POST
```

### Authorization

To access resource endpoints, you need to provide a JWT bearer token in the request `Authorization` header.\
To authenticate, send a POST request to the `/login` endpoint with the following bodies.\
Get a user token:
```json
{
    "return_admin_token": false
}
```
Get an admin access token (used for modifying resources):
```json
{
    "return_admin_token": true
}
```

The response will include a JWT token:
```json
{"admin":true,"token":"jwt_token"}
```

Example for retrieving the token (using cURL):
```sh
curl -X POST 'http://localhost:8080/api/v1/login' \
  -H 'Accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{ "return_admin_token": true }'
```

Use the token like this - put the JWT in the `Authorization` header:
```sh
curl -X GET 'http://localhost:8080/api/v1/authors' \
  -H 'Accept: application/json' \
  -H 'Authorization: Bearer jwt_token'
```

## Testing

### Code tests

Run full code tests (requires a running database server):
```sh
go test ./...
```

Run unit tests using a mocked database:
```sh
go test -short ./...
```

### Performance/load tests

First, reconfigure the server to use `user_test` database user and `paw_test` database (password: `testpass`).\
Example in the `env.yaml` file:
```yaml
DBUSER: "user_test"
DBPASS: "testpass"
DBNAME: "paw_test"
```

Then, [install Grafana k6](https://grafana.com/docs/k6/latest/set-up/install-k6/),
and run JS files with the `Test` suffix, located inside [loadtests directory](/loadtests):
```sh
k6 run loadtests/loadTest.js    # Load test
k6 run loadtests/stressTest.js  # Stress test
k6 run loadtests/spikeTest.js   # Spike test (not recommended on low-end systems)
```

## Stack

 - [Go](https://go.dev/) - main programming language
 - [MariaDB](https://mariadb.org/) - relational database
 - [Gin](https://github.com/gin-gonic/gin) - web framework
 - [Testify](https://github.com/stretchr/testify) - testing toolkit
 - [Grafana k6](https://k6.io/) - performance/load testing
 - [Docker](https://www.docker.com/) - containerization
