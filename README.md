# 🚀 WasaText:

## 1.How to  Run locally

1.0 in gormSQL.go in 1 line set db as:
```shell
dbFile := "database.db" // SQLite This will be your sql file
```

1.1 Run the backend server:
```shell
go run cmd/main.go
```
1.2.1 Frontend server
```shell
./open-node.sh           
```
1.2.2
```shell 
yarn run dev
```

Enjoy the project

## OR RUN EVERYTHING WITH DOCKER
1.0 in gormSQL.go in 1 line set db as:
```shell
dbFile := "/app/data/database.db" // Matches your Docker volume setup
```
2.0 Build docker in terminal with
```
docker-compose up --build
```





## 📂 **Project Structure from Fantastic Coffee**


WASAText aims to accomplish the following objectives:

-Define APIs using the OpenAPI standard.

-Design and develop the server side ("backend") in Go.

-Design and develop the client side ("frontend") in JavaScript.

-Create a Docker container image for deployment.






```plaintext
WasaText
├── cmd/                       
│   ├── database/              # Database configuration
│   │   ├── gorm/              # GORM utilities for database setup
│   │   └── models/            # Database models (e.g., User, Group, Message)
│   └── webapi/                # API setup and server initialization
│       ├── server.go          
│   main.go   
         
├── doc/                       # API documentation
│   └── api.yaml               # OpenAPI specification

├── internal/                  # Core application logic
│   ├── api/                   # HTTP handlers and middleware
│   ├── consts/                # Application constants (e.g., colors, message types)
│   ├── helpers/               # Utility functions
│   ├── repositories/          # Database repository patterns
│   └── service/               # Business logic

# Frontend
├── webui/                     
│   ├── public/                # Static assets
│   └── src/                   # Frontend source code


├── .env                       # Environment variables for the backend
├── docker-compose.yml         # Docker Compose configuration
├── Dockerfile.backend         # Backend Dockerfile
├── Dockerfile.frontend        # Frontend Dockerfile
├── go.mod                     # Go module dependencies
├── go.sum                     # Go module checksums
├── reset.sh                   # Script to reset the database
└── README.md                  # Project documentation
```





