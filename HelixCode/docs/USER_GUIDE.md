# HelixCode User Guide

## 🎯 Introduction

Welcome to HelixCode! This guide will help you get started with the distributed AI development platform.

### What is HelixCode?

HelixCode is an enterprise-grade platform that enables:
- **Distributed Task Execution**: Split large tasks across multiple workers
- **Work Preservation**: Automatic checkpointing and recovery
- **Intelligent Task Division**: Smart task splitting based on capabilities
- **Cross-Platform Development**: Support for multiple development workflows

## 🚀 Quick Start

### 1. Installation

```bash
# Download and build
make build

# Or use pre-built binaries
# Download from releases page
```

### 2. Configuration

Create a configuration file at `~/.config/helixcode/config.yaml`:

```yaml
server:
  address: "0.0.0.0"
  port: 8080

database:
  host: "localhost"
  port: 5432
  user: "helixcode"
  dbname: "helixcode"

auth:
  jwt_secret: "your-secret-key"
```

### 3. Start the Server

```bash
./bin/helixcode
```

## 👤 User Management

### Registration

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "your_username",
    "email": "your@email.com",
    "password": "your_password",
    "display_name": "Your Name"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "your_username",
    "password": "your_password"
  }'
```

Response includes authentication token:

```json
{
  "token": "jwt_token_here",
  "user": {
    "id": "uuid",
    "username": "your_username",
    "email": "your@email.com"
  }
}
```

### Using Authentication

Include the token in subsequent requests:

```bash
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/users/me
```

## 🔧 Worker Management

### Registering a Worker

Workers are distributed nodes that execute tasks. To register a worker:

```bash
curl -X POST http://localhost:8080/api/v1/workers \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "hostname": "worker-1.local",
    "display_name": "Development Worker",
    "capabilities": ["code_generation", "testing"],
    "resources": {
      "cpu_count": 8,
      "total_memory": 17179869184,
      "total_disk": 107374182400
    },
    "ssh_config": {
      "host": "worker-1.local",
      "port": 22,
      "username": "deploy"
    }
  }'
```

### Worker Heartbeat

Workers should send regular heartbeats:

```bash
curl -X POST http://localhost:8080/api/v1/workers/<worker_id>/heartbeat \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "cpu_usage_percent": 45.2,
    "memory_usage_percent": 67.8,
    "disk_usage_percent": 23.1,
    "current_tasks_count": 2
  }'
```

### Monitoring Workers

View all workers:

```bash
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/workers
```

View worker metrics:

```bash
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/workers/<worker_id>/metrics
```

## 📋 Task Management

### Creating Tasks

Tasks represent units of work that can be distributed across workers.

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "task_type": "code_generation",
    "task_data": {
      "language": "go",
      "description": "Create REST API endpoint",
      "specifications": "..."
    },
    "priority": 5,
    "criticality": "normal",
    "estimated_duration": "1h"
  }'
```

### Task Types

- **code_generation**: Generate code based on specifications
- **code_analysis**: Analyze and review code
- **testing**: Run tests and generate reports
- **building**: Compile and build projects
- **refactoring**: Refactor existing code
- **debugging**: Debug and fix issues
- **planning**: Create development plans

### Task Assignment

Tasks are automatically assigned to available workers. You can also manually assign:

```bash
curl -X POST http://localhost:8080/api/v1/tasks/<task_id>/assign \
  -H "Authorization: Bearer <token>"
```

### Monitoring Task Progress

View task status:

```bash
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/tasks/<task_id>
```

List all tasks:

```bash
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/tasks
```

### Work Preservation

HelixCode automatically creates checkpoints for long-running tasks. To manually create a checkpoint:

```bash
curl -X POST http://localhost:8080/api/v1/tasks/<task_id>/checkpoint \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "checkpoint_name": "pre-processing",
    "checkpoint_data": {
      "processed_files": 15,
      "current_state": "processing"
    }
  }'
```

## 🗂️ Project Management

### Creating Projects

Projects organize related tasks and sessions.

```bash
curl -X POST http://localhost:8080/api/v1/projects \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "E-commerce API",
    "description": "REST API for e-commerce platform",
    "workspace_path": "/projects/ecommerce-api",
    "git_repository_url": "https://github.com/user/ecommerce-api"
  }'
```

### Development Sessions

Sessions represent focused development work within a project.

```bash
curl -X POST http://localhost:8080/api/v1/sessions \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "project_uuid",
    "name": "User Authentication",
    "description": "Implement JWT authentication",
    "session_type": "building"
  }'
```

## 📊 System Monitoring

### Health Check

```bash
curl http://localhost:8080/health
```

### System Statistics

```bash
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/system/stats
```

### System Status

```bash
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/system/status
```

## 🔄 Workflows

### Development Workflow

1. **Create Project**: Organize your work
2. **Start Session**: Begin focused development
3. **Create Tasks**: Break down work into manageable units
4. **Monitor Progress**: Track task completion
5. **Review Results**: Analyze completed work

### Distributed Workflow

1. **Register Workers**: Add computing resources
2. **Create Distributed Tasks**: Tasks that can be split
3. **Automatic Assignment**: HelixCode assigns tasks to workers
4. **Progress Tracking**: Monitor distributed execution
5. **Result Aggregation**: Combine results from multiple workers

## 🛠️ Advanced Features

### Task Dependencies

Create tasks that depend on others:

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "task_type": "testing",
    "task_data": { ... },
    "dependencies": ["task_uuid_1", "task_uuid_2"]
  }'
```

### Task Retry

Retry failed tasks:

```bash
curl -X POST http://localhost:8080/api/v1/tasks/<task_id>/retry \
  -H "Authorization: Bearer <token>"
```

### Critical Tasks

Mark tasks as critical to ensure they complete:

```bash
curl -X PUT http://localhost:8080/api/v1/tasks/<task_id> \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "criticality": "critical"
  }'
```

## 🔍 Troubleshooting

### Common Issues

1. **Authentication Failed**
   - Check token expiration
   - Verify username/password
   - Ensure proper Authorization header

2. **Task Not Assigned**
   - Check worker availability
   - Verify worker capabilities
   - Check task dependencies

3. **Worker Offline**
   - Verify worker connectivity
   - Check heartbeat frequency
   - Review worker configuration

### Logs and Debugging

Enable debug logging in configuration:

```yaml
logging:
  level: "debug"
  format: "text"
  output: "stdout"
```

### Performance Optimization

- Use appropriate task priorities
- Distribute work across multiple workers
- Monitor system resources
- Use checkpoints for long-running tasks

## 📚 Best Practices

### Task Design

- **Break Down Work**: Create small, focused tasks
- **Define Dependencies**: Specify task relationships
- **Set Realistic Estimates**: Provide accurate duration estimates
- **Use Appropriate Criticality**: Mark critical tasks appropriately

### Worker Management

- **Regular Heartbeats**: Maintain worker connectivity
- **Resource Monitoring**: Track worker resource usage
- **Capability Matching**: Assign tasks to capable workers
- **Load Balancing**: Distribute work evenly

### Project Organization

- **Clear Naming**: Use descriptive project and session names
- **Proper Documentation**: Document project requirements
- **Version Control**: Integrate with Git repositories
- **Regular Updates**: Keep project information current

## 🔗 Integration

### API Clients

HelixCode provides a comprehensive REST API that can be integrated with:

- **CI/CD Pipelines**: Automate development workflows
- **Monitoring Tools**: Track system performance
- **Development Environments**: Integrate with IDEs
- **Custom Applications**: Build custom interfaces

### Web Interface

Access the web interface at `http://localhost:8080` for a visual management interface.

---

**Need Help?**

- Check the documentation in `/docs`
- Create issues for bugs and feature requests
- Join community discussions
- Contact support for assistance

**Happy developing with HelixCode! 🚀**