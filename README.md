Log Ingestor and Threat Analyzer Microservices
This project consists of two Go-based microservices: Log Ingestor Service and Threat Analyzer Service. These services work together to ingest logs, store them in a MySQL database, and analyze them for potential threats, secured with JWT authentication and containerized using Docker.

Log Ingestor Service: Ingests, stores, and queries logs.
Threat Analyzer Service: Analyzes logs to detect threats.

Both services use the Gin framework, GORM for database interactions, and Swagger for API documentation. They are orchestrated with Docker Compose for easy deployment.

Features Included

JWT Authentication: Secure endpoints with JWT-based authentication.
Rate Limiting: Prevents abuse with configurable rate limits.
Swagger UI: Interactive API documentation.
MySQL Storage: Persistent storage for users, logs, and threat data.
Dockerized: Containerized services for consistent deployment.

Project Structure
.
├── docker-compose.yml                  # Orchestrates services and MySQL
|── README.md
├── log-ingestor-service                # Log Ingestor Service
│   ├── Dockerfile                      # Docker config
│   ├── go.mod, go.sum                 # Go dependencies
│   └── src
│       ├── config                      # MySQL and Swagger setup
│       ├── controllers                 # Auth and log endpoints
│       ├── docs                        # Swagger docs
│       ├── dtos                        # Data Transfer Objects
│       ├── genric_error                # Error responses
│       ├── main.go                     # Service entry point
│       ├── middleware                  # JWT auth and rate limiting
│       ├── models                      # User and log models
│       ├── routes                      # Public and protected routes
│       ├── services                    # Log ingestion logic
│       └── utils                       # Utility functions
├                  
└── threat-analyzer-service             # Threat Analyzer Service
    ├── Dockerfile                      # Docker config
    ├── go.mod, go.sum                 # Go dependencies
    └── src
        ├── config                      # MySQL and Swagger setup
        ├── controllers                 # Auth and threat endpoints
        ├── docs                        # Swagger docs
        ├── dto                         # Data Transfer Objects
        ├── genric_error                # Error responses
        ├── main.go                     # Service entry point
        ├── middleware                  # JWT auth and rate limiting
        ├── models                      # User, log, and threat models
        ├── routes                      # Public and protected routes
        ├── services                    # Threat analysis logic
        └── utility                     # Threat detection utilities

Prerequisites

Go 1.18+
Docker and Docker Compose
MySQL 8.0+
Swaggo CLI:go install github.com/swaggo/swag/cmd/swag@latest



Setup

Clone the Repository:
git clone https://github.com/yatenderpareek7887/security-service1.git
cd security-service


Configure Environment Variables:

Create .env files in log-ingestor-service and threat-analyzer-service:PORT=8080/8081        

BASE_PATH=/api/
JWT_SECRET_KEY=secure-secret-key-32-chars-long-1234567890
DB_USER=root
DB_PASSWORD=root
DB_HOST=mysql
DB_PORT=3306
DB_NAME=log_ingestor_db



Install Dependencies:
cd log-ingestor-service
go mod tidy
cd ../threat-analyzer-service
go mod tidy


Generate Swagger Docs:
cd log-ingestor-service/src
swag init
cd ../threat-analyzer-service/src
swag init


Running the Services, Using Docker Compose:

Start services and MySQL:
- docker compose up -d --build

Running the Services, Manually:

Start MySQL locally or via Docker.
Run each service:cd log-ingestor-service/src
go run main.go
cd ../threat-analyzer-service/src
go run main.go

Access:
Log Ingestor: http://localhost:8080
Threat Analyzer: http://localhost:8081



API Usage

Swagger UI:
Log Ingestor: http://localhost:8080/swagger/api/docs/index.html
Threat Analyzer: http://localhost:8081/swagger/api/docs/index.html

Database
log_ingestor_db

Tables:
users: Stores user data (id, username, password, email, created_at, update_at, deleted_at).
log_data (Log Ingestor): Stores logs (id, username, message, source, created_at, update_at).
threats (Threat Analyzer): Stores threat analyses (id, username, log_id, threat_level, description, created_at, update_at).

Docker Commands::
- docker compose ps

- docker compose up -d --build

- docker compose logs log-ingestor-service

- docker compose logs threat-analyzer-service

- docker compose logs mysql

- docker compose down
