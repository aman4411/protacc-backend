# 📦 protacc-backend

Backend service for a Protacc website, built using **Go Fiber**, **PostgreSQL**, and hosted on **Render**. 
It includes user authentication, admin control over services, and is designed to be cost-effective with free-tier services.

---

## 🚀 Features

- ✅ RESTful API with Go Fiber
- ✅ PostgreSQL integration (via Neon)
- ✅ Structured logging (Info/Error)
- ✅ Environment-based profiling (Dev/Prod)
- ✅ Database migrations using `golang-migrate`
- ✅ Render-ready deployment
- ✅ CORS-enabled for frontend-backend communication
- 🔜 JWT Authentication, Admin Panel, Payment Integration

---

## ⚙️ Tech Stack

| Tool             | Description                            |
|------------------|----------------------------------------|
| Go (Golang)      | Core backend language                  |
| Fiber            | Web framework for routing & middleware |
| PostgreSQL       | Relational DB (Neon - free tier)       |
| golang-migrate   | SQL migration tool                     |
| godotenv         | Loads environment variables            |
| Render           | Cloud hosting for backend              |

---

## 📁 Project Structure

```plaintext
protacc-backend/
├── main.go                # Entry point of the application
├── .env                   # Environment variables (local only, gitignored)
├── go.mod / go.sum        # Go module dependencies
│
├── config/                # Configuration loaders (e.g. logger, env, constants)
│   └── config.go
│
├── db/                    # Database layer
│   ├── db.go              # DB connection setup (PostgreSQL)
│   └── migrations/        # SQL migration files (.up.sql and .down.sql)
│
├── handlers/              # HTTP route handlers (e.g. auth, service)
│   └── user_handler.go
│
├── models/                # DB model structs (User, Service, etc.)
│   └── user.go
│
├── utils/                 # Utility functions (e.g. logger, helpers)
│   └── logger.go
│
└── README.md              # Project documentation
```

## 🧑‍💻 Local Setup

Follow these steps to run the project locally:

### 1. Clone the repository

```bash
    git clone https://github.com/yourusername/protacc-backend.git
    cd protacc-backend
```

### 2. Create a `.env` file
Create a `.env` file in the root directory and add your environment variables. Example:

```plaintext
    PORT=8080
    DB_URL=your_postgresql_connection_string
    JWT_SECRET=your_jwt_secret
```

### 3. Install dependencies
```bash
    go mod tidy
```

### 4. Install `golang-migrate` CLI
```bash
    brew install golang-migrate
```

### 5. Create migrations
```bash
    migrate create -ext sql -dir db/migrations <name>
```

### 6. Run migrations
```bash
    migrate -path db/migrations -database "your_postgresql_connection_string" up
```

### 7. Run the application
```bash
    go run main.go
```

### 8. Test the API
```bash
    curl http://localhost:8080/ping
    # Expected response: "pong"
```





