# F5 Todo API

A simple Todo REST API built with Go and PostgreSQL. The project supports creating, reading, updating, and deleting tasks, with validation, pagination, filtering, database migrations, Docker support, and unit tests.

## How to Run

### Requirements

- Go 1.26+
- Docker and Docker Compose

### Run with Docker

Clone the repository and move into the project directory:

- git clone <repository-url>
- cd F5

Create a .env file in the project root based on the .env.example file:

POSTGRES_USER=root
POSTGRES_PASSWORD=your_password
POSTGRES_DB=f5_todo
API_URL=":8082"

Start the application:

- docker compose up --build

The API will be available at:

- http://localhost:8085

PostgreSQL is exposed on port 5433 on the host machine.

### Run locally

Install the dependencies:

- go mod download

Run the API:

- go run ./cmd/api

The application expects the database connection settings to be available through environment variables.

### Run tests

Run all tests:

- go test ./...

Run tests with coverage:

- go test ./... -cover

To generate a detailed coverage report:

- go test ./... -coverprofile=coverage.out
- go tool cover -html=coverage.out


## Technology Choices

### Go

I chose Go for the backend because it is relatively simple to structure, has good performance, and has useful built-in support for HTTP servers, testing, and concurrency. It also helped me keep the API small without relying on a large framework.

### PostgreSQL

PostgreSQL was used as the database because the project mainly deals with structured relational data. It also provides good support for UUIDs, timestamps, constraints, and transactions.

### pgx

The project uses pgx to communicate with PostgreSQL. It is a lightweight PostgreSQL driver for Go and gives direct access to PostgreSQL features without adding an ORM.

### Docker

Docker and Docker Compose are used to run the API, PostgreSQL database, and migrations together. This makes the project easier to run without having to manually configure PostgreSQL on the host machine.

## What I Would Improve

If I had more time, I would improve a few things:

- Add more integration tests that test the API together with a real PostgreSQL database.
- Improve the API documentation, for example by adding Swagger/OpenAPI.
- Add more structured logging instead of relying mostly on standard logging.
- Improve the Docker setup for production instead of mainly focusing on local development.
- Add more database constraints and indexes where they make sense.
- Add CI checks for tests, formatting, and linting.

I would also spend more time improving the error handling and making the API responses more consistent.

## Assumptions

A few decisions had to be made where the requirements were not completely specific:

- New tasks are created with the default pending status.
- Due dates use the YYYY-MM-DD format.
- Due dates are optional.
- Pagination uses limit and offset.
- An invalid UUID is treated as a bad request.