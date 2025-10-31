# CRUSH.md - Helix CLI Development Guide

## Build Commands
- `go build` - Build the project
- `go test ./...` - Run all tests
- `go test -v ./path/to/package` - Run specific package tests
- `go run main.go` - Run the application

## Code Style Guidelines
- **Language**: Go (Golang)
- **Imports**: Group standard library, third-party, and local imports
- **Formatting**: Use `gofmt` or `goimports` for consistent formatting
- **Naming**: camelCase for variables/functions, PascalCase for exported types
- **Error Handling**: Always handle errors explicitly, use proper error wrapping
- **Types**: Use strong typing with structs and interfaces
- **Documentation**: Document exported functions with Go doc comments

## Project Structure
- Components are decoupled Go programs bound by bash scripts
- Follow reusability principle with generics and interfaces
- Each component can be used separately or as part of larger application

## Testing Requirements
- 100% code coverage required for all test types
- Unit, integration, automation, end-to-end testing
- SonarQube and Snyk scanning via Docker containers
- All tests must execute successfully

## Development Workflow
- Use the provided `commit` script for commits with project context
- Follow the CLI specifications in `Specification/CLI_Specs.md`
- Components include: LLM APIs, CLI UI, database, debugger, planner, builder, refactorer, tester

## Configuration
- User config stored in `~/.config/Helix/helix.json`
- Project-specific config in `Helix.md` files
- Support for multiple LLM providers and models