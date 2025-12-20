# QingLong Panel (Go Implementation)

This is a Go implementation of the popular QingLong Panel, a task management system for automated scripts.

## Features

- Task scheduling and management
- RESTful API
- User authentication
- Environment variable management
- Script file management
- Task execution and logging
- Web interface
- Docker support

## Project Structure

```
├── main.go              # Application entry point
├── go.mod               # Go modules definition
├── go.sum               # Go modules checksums
├── Dockerfile           # Docker configuration
├── docker-compose.yml   # Docker Compose configuration
├── configs/             # Configuration files
├── internal/
│   ├── controllers/     # HTTP handlers
│   ├── models/          # Data structures
│   ├── services/        # Business logic
│   └── utils/           # Utility functions
├── web/
│   ├── static/          # Static assets (CSS, JS, images)
│   └── templates/       # HTML templates
├── data/                # Data storage
├── logs/                # Log files
└── scripts/             # User scripts
```

## Getting Started

### Using Go directly

1. Install Go (version 1.21 or higher)
2. Clone the repository
3. Run `go mod tidy` to install dependencies
4. Run `go run main.go` to start the server

### Using Docker

1. Install Docker and Docker Compose
2. Run `docker-compose up -d` to start the server

The server will start on port 8080.

## API Endpoints

### Authentication
- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration

### Tasks
- `GET /api/tasks` - Get all tasks
- `POST /api/tasks` - Create a new task
- `GET /api/tasks/:id` - Get a specific task
- `PUT /api/tasks/:id` - Update a specific task
- `DELETE /api/tasks/:id` - Delete a specific task

### Task Execution
- `POST /api/execute/task/:id` - Execute a task by ID
- `POST /api/execute/command` - Execute a command directly
- `GET /api/execute/results` - Get last execution results

### Environment Variables
- `GET /api/env` - Get all environment variables
- `POST /api/env` - Create a new environment variable
- `GET /api/env/:id` - Get a specific environment variable
- `PUT /api/env/:id` - Update a specific environment variable
- `DELETE /api/env/:id` - Delete a specific environment variable

### Scripts
- `GET /api/scripts` - Get all scripts
- `POST /api/scripts` - Create a new script
- `GET /api/scripts/:id` - Get a specific script
- `PUT /api/scripts/:id` - Update a specific script
- `DELETE /api/scripts/:id` - Delete a specific script

## Web Interface

- `/` - Dashboard
- `/login` - Login page
- `/tasks` - Task management
- `/scripts` - Script management
- `/environments` - Environment variable management

## Default Admin User

Username: `admin`
Password: `admin123`

## Configuration

The application can be configured using the `configs/config.json` file:

```json
{
  "server": {
    "port": 8080,
    "host": "0.0.0.0"
  },
  "database": {
    "type": "sqlite",
    "path": "./data/ql.db"
  },
  "security": {
    "jwt_secret": "ql_panel_secret_key",
    "password_salt": "ql_panel_salt"
  },
  "task": {
    "default_timeout": 3600,
    "log_retention_days": 30
  }
}
```

## Docker Support

The application includes Docker support for easy deployment:

1. Build the Docker image: `docker build -t baihu .`
2. Run the container: `docker run -p 8080:8080 baihu`

Or use Docker Compose:
```bash
docker-compose up -d
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.