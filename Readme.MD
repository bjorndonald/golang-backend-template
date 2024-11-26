# Golang Backend template

## Overview
This is a production ready made template with Golang GIN with a Next.js demo. This should make your development of application a lot faster with the features already implemented.

This project is heavily inspired by Johnson Awah Alfred's work. Be sure to check his project out. [Link here](https://github.com/CeoFred/osalobo)

## 🚀 Features

- [x] [Gin Framework](https://github.com/gin-gonic/gin) for routing and middleware
- [x] Two-Factor authentication
- [x] JWT Tokens and Session Management
- [x] Location tracking and Device tracking
- [x] Password Recovery
- [x] PostgreSQL integration with migration support
- [x] Swagger API documentation
- [x] API monitoring with APIToolkit
- [x] Custom error handling and logging
- [x] CORS configuration
- [x] Docker and Docker Compose support
- [x] OTP management system
- [x] Static file serving with Cloudinary file upload
- [x] Environment configuration
- [x] Hot reload during development
- [x] Code security scanning with gosec
- [x] Event streaming with Apache Kafka
- [x] Transactional message processing
- [x] Consumer group management
- [x] Event broadcasting system

## 🚀 Future additions
- [x] Encryption support
- [x] Webhook support
- [x] AI support
- [x] More testing
- [x] s3 support

## 📋 Prerequisites

- Go 1.x
- Docker and Docker Compose
- PostgreSQL (if running locally)
- Make

## 🛠️ Installation

1. Clone the repository
```bash
git clone https://github.com/bjorndonald/golang-backend-template.git
cd gin-boilerplate
```

2. Copy the example environment file
```bash
cp .env.example .env
```

3. Install dependencies
```bash
make requirements
```

## 🔧 Configuration

Update the `.env` file with your configuration:

```env
PORT=8000
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=yourdatabase

JWT_SECRET=
RESEND_API_KEY=

CLOUDINARY_API_KEY=
CLOUDINARY_API_SECRET=
CLOUDINARY_NAME=
SEND_FROM_NAME=
SEND_FROM_EMAIL=
```

## 🚀 Running the Application

### Local Development

```bash
# Run the application with hot reload
make run-local
```

### Using Docker

```bash
# Build the Docker image
make build

# Start all services using Docker Compose
make service-start

# Start frontend demo
make run-frontend
```

## 📚 API Documentation

Swagger documentation is available at:
```
http://localhost:8000/swagger/index.html
```

To regenerate Swagger documentation:
```bash
make docs-generate
```

## 🛠️ Available Make Commands

- `make run-local` - Run the application locally with hot reload
- `make docs-generate` - Generate Swagger documentation
- `make requirements` - Install/update dependencies
- `make clean-packages` - Clean Go module cache
- `make build` - Build Docker image
- `make start-postgres` - Start PostgreSQL container
- `make stop-postgres` - Stop PostgreSQL container
- `make start` - Start application with Docker
- `make build-no-cache` - Build Docker image without cache
- `make service-stop` - Stop all Docker Compose services
- `make service-start` - Start all Docker Compose services

## 📁 Project Structure

```

.
├── constants/           # Application constants and configuration
├── database/           # Database connection and migrations
├── docs/              # Swagger documentation
├── frontend/         # Next.js Demo
├── internal/
│   ├── bootstrap/     # Application bootstrapping
│   ├── helpers/       # Helper functions
│   ├── otp/          # OTP management
│   ├── repository/    # Repository management
│   ├── routes/        # API routes
│   ├── services/        # services
│     ├── email.go      # email service
│     └── streaming/     # Kafka streaming implementation
│       ├── consumer.go  # Kafka consumer implementation
│       ├── producer.go  # Kafka producer implementation
│       └── events.go    # Event type definitions
├── utils/         # Utilities
├── resend/         # Resend client implementation
├── templates/         # Email templates
├── main.go           # Application entry point
├── Dockerfile        # Docker configuration
├── docker-compose.yml # Docker Compose configuration
└── Makefile          # Build and development commands
```

## 🔒 Security

The project includes security measures:
- Custom recovery middleware
- CORS configuration
- Request logging
- Security scanning with gosec

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.

## 👤 Contact

Bjorn-Donald Bassey - [bjorndonaldb@gmail.com](mailto:bjorndonaldb@gmail.com)

## ⭐️ Show your support

Give a ⭐️ if this project helped you!