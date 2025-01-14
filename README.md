# 🚀 WasaText:



WASAText aims to accomplish the following objectives:

-Define APIs using the OpenAPI standard.

-Design and develop the server side ("backend") in Go.

-Design and develop the client side ("frontend") in JavaScript.

-Create a Docker container image for deployment.

## 📂 **Project Structure from Fantastic Coffee**

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


## How to run 

1.Clone the repository:
```shell
git clone https://github.com/your-repository.git
cd wasaText
```

2.Set up environment variables by creating a .env file with following structure:
```shell
SERVER_PORT=3000
API_SECRET=your-secret
TTL_HOUR=12
```

3.Run the backend server:
```shell
go run ./cmd/webapi/main.go
```
or Everything with Docker
```
docker-compose up --build
```


