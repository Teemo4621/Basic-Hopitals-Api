# 🌟 Basic Hospitals API

## Overview
🏥 Basic Hospitals API is a simple RESTful API built with Go (Golang) to manage hospital-related data, including hospitals, patients, and staff. This project is designed to demonstrate a clean architecture with modular structure, making it easy to extend or integrate with other systems. 🚀

## Features
- 🏢 Manage hospital information (create, read, update, delete).
- 👤 Manage patient records with personal details and hospital affiliation.
- 👩‍⚕️ Manage staff accounts with authentication support.
- 🔗 RESTful endpoints with proper error handling.
- ✅ Mock implementations for testing.

## Project Structure
- `configs/`: 📋 Configuration files (e.g., database settings, environment variables).
- `modules/`: 🔧 Core modules like hospitals, JWT, patient, and staff logic.
- `entities/`: 📊 Data models or structs for hospitals, patients, and staff.
- `hospitals/`, `patients/`, `staff/`: 📂 Subdirectories for each entity with:
  - `controllers/`: 🎛️ Handle API logic and routes.
  - `repositories/`: 🗃️ Database interaction (includes mock files for testing).
  - `usecases/`: 🤖 Business logic implementation.
- `mocks/`: 🧪 Mock files for unit testing.
- `servers/`: 🖥️ Server setup and handlers.
- `pkgs/`: 📦 Custom packages or utilities.
- `docker-compose.yml` & `Dockerfile`: 🐳 Docker configuration for containerization.
- `go.mod` & `go.sum`: 📜 Go module dependencies.

## Prerequisites
- Go (version 1.18 or higher) 🐹
- Docker (for containerized environment) 🐳

## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/Teemo4621/Basic-Hopitals-Api.git
   cd Basic-Hopitals-Api
   ```
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Set up environment variables:
   - Create a `.env` file based on `.env.development` and configure database credentials. 🔧
4. Run with Docker (optional):
   ```bash
   docker-compose up -d --build
   ```
5. Run locally:
   ```bash
   go run main.go
   ```

## API Endpoints
- `GET /hospitals`: 📋 List all hospitals.
- `POST /hospitals`: ➕ Create a new hospital.
- `GET /patients`: 📋 List all patients.
- `POST /patients`: ➕ Add a new patient.
- `GET /staff`: 📋 List all staff.
- `POST /staff`: ➕ Add a new staff member.
(For detailed schema, refer to the `entities` directory or API documentation.)

## Usage
- Access the API at `http://localhost:8080` (default port). 🌐
- Use tools like Postman or curl to test endpoints. 🛠️
- Authentication may be required for certain routes (e.g., staff management) using JWT. 🔒

## Development
- **Testing**: Run `go test ./...` to execute unit tests. ✅
- **Contributing**: Fork the repository, create a feature branch, and submit a pull request. 🤝
- **Issues**: Report bugs or suggest features via the Issues tab. 🐛