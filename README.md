# OpsWS

OpsWS is a task runner platform similar to Jenkins, built in Golang. It supports parsing Jenkins pipelines in YAML format, running tasks on agents, logging, and more.

## Features
- Master-server with API for pipeline management
- Agents for task execution
- SQLite database backend
- Redis for queuing and logging
- Support for sequential/parallel stages and steps
- Step types: sh, py, go, scm, artifact, notify
- Inline code or script files with arguments
- Agent selection based on labels
- Real-time log streaming via WebSocket
- Swagger API documentation
- Credentials and Secrets Management
- SCM integration (Git)
- Artifacts Management
- Notifications (email, Slack)
- Authentication and RBAC
- Pipeline Scheduling/Cron
- Workspace Management and Cleanup

## Setup
1. Install dependencies: `go mod tidy`
2. Run master: `go run master/cmd/main.go`
3. Run agent: `go run agent/cmd/main.go`
4. Use API to create pipelines and schedule tasks

Generate Swagger docs: `swag init`
Access Swagger UI at `/swagger/index.html`