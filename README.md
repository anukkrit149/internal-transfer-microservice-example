# Internal Transfer Microservice

A Go REST server for managing account transfers built with domain-driven design principles using Gin, GORM, PostgreSQL, and Redis. This microservice provides APIs for creating accounts, retrieving account information, and transferring money between accounts with concurrency control and deadlock prevention.

## Architecture

The project follows a domain-driven design architecture with the following layers:

- **Domain Layer**: Contains the core business logic and entities
- **Repository Layer**: Handles data access and persistence
- **Service Layer**: Implements business logic and orchestrates repositories
- **Controller Layer**: Handles HTTP requests and responses
- **Infrastructure Layer**: Provides technical capabilities (database, cache)
- **Routes Layer**: Defines API endpoints
- **Factory Layer**: Creates and wires up components
- **Config Layer**: Manages application configuration

## Project Structure

```
internal-transfer-microservice/
├── main.go                           # Application entry point
├── config/                           # Configuration files
│   └── env.yaml                      # YAML configuration example
├── internal/
│   ├── config/                       # Configuration package
│   │   ├── config.go                 # Configuration structs and loading
│   │   └── getter.go                 # Configuration getters
│   ├── domain/                       # Domain models and interfaces
│   │   ├── base_model.go             # Base model for all domain models
│   │   └── account/                  # Account domain
│   │       ├── model.go              # Account model
│   │       ├── interface.go          # Account interfaces
│   │       └── structs.go            # Account-related request/response structs
│   ├── repository/                   # Repository implementations
│   │   └── account.go                # Account repository implementation
│   ├── service/                      # Service implementations
│   │   ├── account.go                # Account service implementation
│   │   └── account_test.go           # Tests for account service
│   ├── controller/                   # Controller implementations
│   │   └── account.go                # Account controller implementation
│   ├── routes/                       # Route definitions
│   │   └── account.go                # Account routes
│   ├── infrastructure/               # Infrastructure components
│   │   ├── db/                       # Database connections
│   │   │   ├── interface.go          # Database interface
│   │   │   └── postgres.go           # PostgreSQL implementation
│   │   └── cache/                    # Cache implementations
│   │       ├── interface.go          # Cache interface
│   │       └── redis.go              # Redis implementation
│   └── factory/                      # Factory pattern implementations
│       └── factory.go                # Application factory
├── pkg/
│   └── logger/                       # Logging package
│       ├── interface.go              # Logger interface
│       ├── logger.go                 # Logger implementation
│       └── logrus.go                 # Logrus logger implementation
├── Dockerfile                        # Docker configuration
└── Makefile                          # Build and deployment automation
```

## Design Patterns

- **Repository Pattern**: Abstracts data access logic
- **Factory Pattern**: Creates and wires up components
- **Dependency Injection**: Components receive their dependencies
- **Interface Segregation**: Interfaces define specific behaviors

## Key Features

### Account Management
- Create accounts with initial balance
- Retrieve account information by ID
- Validate account existence and balance

### Money Transfer with Concurrency Control
- Transfer money between accounts with transaction support
- Prevent insufficient balance transfers
- Ensure data consistency with database transactions

### Deadlock Prevention
- Implement resource ordering to prevent deadlocks
- Use distributed locks with Redis for concurrent access control
- Handle concurrent transfers between the same accounts in opposite directions
- Ensure consistent final balances regardless of execution order

## API Endpoints

- `GET /api/v1/accounts/:id`: Get an account by ID
- `POST /api/v1/accounts`: Create a new account with initial balance
- `POST /api/v1/accounts/transfer`: Transfer money between accounts
- `GET /health`: Health check endpoint

## Prerequisites

- Go 1.24 or higher
- PostgreSQL
- Redis

## Configuration

The application supports configuration through YAML, JSON, or environment variables. 

### YAML Configuration

Create a YAML file (e.g., `config/env.yaml`):

```yaml
# Server configuration
server:
  port: "3000"  # Default port for Docker deployment
  gin_mode: "debug"
  shutdown_timeout: 5  # Graceful shutdown timeout in seconds

# Database configuration
database:
  host: "localhost"
  port: "5432"
  user: "myuser"
  password: "mypassword"
  name: "account_db"
  sslmode: "disable"

# Redis configuration
redis:
  host: "localhost"
  port: "6379"
  password: ""
  db: 0
```

### Environment Variables

You can also use environment variables:

```
SERVER_PORT=3000
SERVER_GIN_MODE=debug
SERVER_SHUTDOWN_TIMEOUT=5
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=myuser
DATABASE_PASSWORD=mypassword
DATABASE_NAME=account_db
DATABASE_SSLMODE=disable
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

## Running the Application

1. Clone the repository
2. Set up PostgreSQL and Redis
3. Create a configuration file (YAML, JSON, or ENV)
4. Run database migrations:

```bash
# Using the default configuration
go run main.go migrate

# Using a specific configuration file
go run main.go migrate --config config/env.yaml
```

5. Start the API server:

```bash
# Using the default configuration
go run main.go api

# Using a specific configuration file
go run main.go api --config config/env.yaml
```

## Building the Application

You can build the application using the provided Makefile:

```bash
# Build the application
make build

# Clean build artifacts
make clean
```

## Running Tests

```bash
# Run all tests
make test
```

The project includes comprehensive tests for the account service, covering:

1. Creating an account
2. Getting an existing account
3. Getting a non-existent account
4. Simple money transfer between accounts
5. Concurrent transfers on different accounts
6. Deadlock prevention for concurrent transfers between the same accounts in opposite directions

The tests use mock implementations of the repository and cache interfaces to isolate the service layer for unit testing.

## Docker Support

The application can be containerized using Docker.

### Building the Docker Image

```bash
# Build the Docker image
make docker-build
```

### Running with Docker

```bash
# Run the application in Docker
make docker-run  # This will expose the application on port 3000

# Stop the Docker container
make docker-stop
```

### Pushing to DockerHub

```bash
# Push the Docker image to DockerHub
make docker-push
```

Note: Before pushing to DockerHub, update the `DOCKER_REPO` variable in the Makefile with your DockerHub username.

## Makefile Commands

The project includes a Makefile with the following commands:

- `make build` - Build the application
- `make clean` - Clean build artifacts
- `make run` - Run the application
- `make run-with-config` - Run with a custom config file
- `make migrate` - Run database migrations
- `make migrate-with-config` - Run migrations with a custom config file
- `make test` - Run tests
- `make docker-build` - Build Docker image
- `make docker-push` - Push Docker image to registry
- `make docker-run` - Run the application in Docker
- `make docker-stop` - Stop and remove Docker container
- `make help` - Show help information
