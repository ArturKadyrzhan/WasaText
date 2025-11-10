## Preference
**WasaText — Simple text messanger**

**This project was developed for educational purposes. Project based exam of "Web and Software Architecture" course.**

Distributed full-stack web application designed to simulate a modern messaging system.

It demonstrates the integration of a **Golang REST API backend**, **Vue 3 frontend**, and **SQLite database**,emphasizing modular architecture, scalability, and maintainability.

## 1 Overview
The project implements:
- RESTful communication between client and server  
- Authentication and persistent user state  
- Group and private chat interactions  
- Local data storage with SQLite  
- Containerized deployment via Docker


## 2 System Architecture

```text
Frontend (Vue 3 + Vite)
 ├─ Axios Client & Auth Store
 └─ UI Components (SPA)
         │
         ▼
Backend (Go REST API)
 ├─ API Handlers & Database Layer
 └─ SQLite (WAL Mode)
```



## 3 Execution
### Local
```bash
cd cmd/webapi && go run .
cd webui && npm install && npm run dev
```

### Docker
```bash
docker build -t my-go-app -f Dockerfile.backend .
docker run -p 3000:3000 --name my-running-app my-go-app
docker build -t my-frontend -f Dockerfile.frontend .
docker run -p 4000:80 --name frontend-dev my-frontend
```


