# 🌀 HelixCode - Distributed AI Development Platform

**Version**: 1.0.0  
**Build**: 2025-11-01_02:53:21  
**Commit**: 42a36df

## 🚀 Overview

HelixCode is an enterprise-grade distributed AI development platform that enables intelligent task division, work preservation, and cross-platform development workflows. Built with Go and designed for scalability, HelixCode provides a robust foundation for distributed computing with automatic checkpointing, rollback functionality, and real-time monitoring.

## ✨ Key Features

### 🎯 Phase 1: Foundation (Completed)
- **✅ Database Schema**: Complete PostgreSQL schema with 11 tables
- **✅ Authentication System**: JWT-based auth with session management
- **✅ Worker Management**: Distributed worker registration and health monitoring
- **✅ Task Management**: Intelligent task division with work preservation
- **✅ Logo Integration**: Automatic asset generation with color extraction
- **✅ REST API**: Comprehensive HTTP API with Gin framework
- **✅ Configuration System**: Flexible config with environment variables

### 🔮 Upcoming Phases
- **Phase 2**: Core Services (LLM integration, MCP protocol)
- **Phase 3**: Workflows (Project management, development modes)
- **Phase 4**: Clients & Integration (Terminal UI, CLI, Cross-platform)

## 🏗️ Architecture

```
HelixCode Architecture
├── API Layer (REST + WebSocket)
├── Core Services
│   ├── Authentication
│   ├── Worker Management
│   ├── Task Management
│   ├── Project Management
│   └── Session Management
├── Database Layer (PostgreSQL + Redis)
├── Distributed Workers
└── Cross-Platform Clients
```

## 🛠️ Installation

### Prerequisites
- Go 1.21+
- PostgreSQL 15+
- Redis 7+

### Quick Start

1. **Clone and build**:
   ```bash
   cd HelixCode
   make build
   ```

2. **Setup database**:
   ```bash
   # Create database and user
   createdb helixcode
   createuser helixcode
   ```

3. **Configure environment**:
   ```bash
   export HELIX_DATABASE_PASSWORD=your_password
   export HELIX_AUTH_JWT_SECRET=your_jwt_secret
   ```

4. **Run server**:
   ```bash
   ./bin/helixcode
   ```

## 📁 Project Structure

```
HelixCode/
├── cmd/
│   ├── server/          # Main server application
│   └── cli/             # CLI client (upcoming)
├── internal/
│   ├── auth/            # Authentication system
│   ├── config/          # Configuration management
│   ├── database/        # Database layer
│   ├── logo/            # Logo processing & assets
│   ├── server/          # HTTP server & API
│   ├── task/            # Task management
│   ├── theme/           # Color themes from logo
│   └── worker/          # Worker management
├── assets/
│   ├── colors/          # Color schemes
│   ├── icons/           # Platform icons
│   └── images/          # Logo & ASCII art
├── config/
│   └── config.yaml      # Configuration file
└── scripts/
    └── logo/            # Asset generation scripts
```

## 🔧 Configuration

### Environment Variables

```bash
# Database
HELIX_DATABASE_HOST=localhost
HELIX_DATABASE_PORT=5432
HELIX_DATABASE_USER=helixcode
HELIX_DATABASE_PASSWORD=your_password
HELIX_DATABASE_DBNAME=helixcode

# Authentication
HELIX_AUTH_JWT_SECRET=your_jwt_secret
HELIX_AUTH_TOKEN_EXPIRY=86400

# Server
HELIX_SERVER_ADDRESS=0.0.0.0
HELIX_SERVER_PORT=8080
```

### Configuration File

See `config/config.yaml` for complete configuration options.

## 🎨 Design System

HelixCode features a comprehensive design system extracted from the project logo:

- **Primary Color**: #C2E95B (Extracted from logo)
- **Secondary Color**: #C0E853
- **Accent Color**: #B8ECD7
- **Text Color**: #2D3047
- **Background**: #F5F5F5

All platform icons and themes are automatically generated from the source logo.

## 📊 Database Schema

### Core Tables
- **users**: User accounts and authentication
- **user_sessions**: Active user sessions
- **workers**: Distributed worker nodes
- **worker_metrics**: Worker performance metrics
- **distributed_tasks**: Task management with work preservation
- **task_checkpoints**: Automatic checkpointing system
- **projects**: Project management
- **sessions**: Development sessions

### Work Preservation Features
- Automatic checkpointing for long-running tasks
- Dependency-based task execution
- Criticality-based pausing
- Graceful degradation during worker failures
- Comprehensive rollback functionality

## 🔌 API Endpoints

### Health Check
- `GET /health` - System health status

### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/logout` - User logout
- `POST /api/v1/auth/refresh` - Token refresh

### Users
- `GET /api/v1/users/me` - Get current user
- `PUT /api/v1/users/me` - Update current user
- `DELETE /api/v1/users/me` - Delete current user

### Workers
- `GET /api/v1/workers` - List workers
- `POST /api/v1/workers` - Register worker
- `GET /api/v1/workers/:id` - Get worker details
- `PUT /api/v1/workers/:id` - Update worker
- `DELETE /api/v1/workers/:id` - Delete worker
- `POST /api/v1/workers/:id/heartbeat` - Worker heartbeat
- `GET /api/v1/workers/:id/metrics` - Worker metrics

### Tasks
- `GET /api/v1/tasks` - List tasks
- `POST /api/v1/tasks` - Create task
- `GET /api/v1/tasks/:id` - Get task details
- `PUT /api/v1/tasks/:id` - Update task
- `DELETE /api/v1/tasks/:id` - Delete task
- `POST /api/v1/tasks/:id/assign` - Assign task to worker
- `POST /api/v1/tasks/:id/start` - Start task execution
- `POST /api/v1/tasks/:id/complete` - Complete task
- `POST /api/v1/tasks/:id/fail` - Mark task as failed
- `POST /api/v1/tasks/:id/checkpoint` - Create checkpoint
- `GET /api/v1/tasks/:id/checkpoints` - List checkpoints
- `POST /api/v1/tasks/:id/retry` - Retry failed task

## 🧪 Development

### Build Commands

```bash
make build          # Build the application
make test           # Run all tests
make clean          # Clean build artifacts
make logo-assets    # Generate logo assets
make setup-deps     # Setup dependencies
make fmt            # Format code
make lint           # Lint code
make dev            # Run development server
make prod           # Build for production
```

### Testing

```bash
# Run all tests
go test ./...

# Run specific package tests
go test -v ./internal/auth

# Run with coverage
go test -cover ./...
```

## 🔒 Security

- JWT-based authentication
- Password hashing with bcrypt
- CORS and security headers
- Input validation
- SQL injection protection
- Environment-based secret management

## 📈 Monitoring

- Database health checks
- Worker connectivity monitoring
- Task progress tracking
- System metrics collection
- Real-time status updates

## 🚦 Roadmap

### Phase 2: Core Services (Weeks 5-8)
- [ ] LLM provider integration
- [ ] MCP protocol implementation
- [ ] Advanced reasoning engine
- [ ] Multi-channel notifications

### Phase 3: Workflows (Weeks 9-12)
- [ ] Project management system
- [ ] Development workflows
- [ ] Testing mode implementation
- [ ] Refactoring capabilities

### Phase 4: Clients & Integration (Weeks 13-16)
- [ ] Terminal UI with BubbleTea
- [ ] CLI implementation
- [ ] Cross-platform clients
- [ ] Performance optimization

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

- **Documentation**: See `/docs` directory
- **Issues**: Create GitHub issues for bugs and feature requests
- **Discussions**: Join project discussions for questions

---

**Built with ❤️ using Go, PostgreSQL, and distributed computing principles**