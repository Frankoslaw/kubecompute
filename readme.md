### TODO:
- Features:
    - [ ] HealthCheck probes
    - [ ] Lifecycle providers
    - [ ] Kubectl like cli
    - [ ] NATS + websocket access for streaming resources
    - [ ] Explicite templates
    - [ ] Explicite image
    - [ ] Artifiact stores
- Error handling:
    - [ ] Diffrentiate does not exist and version conflict errors
- Observability:
    - [ ] Zap
    - [ ] OpenTelemetry
- Testing:
    - [ ] Unit tests
    - [ ] Integration tests
- Safety:
    - [ ] Nuke command
    - [ ] Better resource insight
    - [ ] Scan for orphans
    - [ ] Dry runs
    - [ ] Mock adapters
    - [ ] Split status updates from spec updates
    - [ ] Permissions with JWT over namespace tenants
    - [ ] Gitleaks
    - [ ] External key store support ex. vault
- QOL:
    - [ ] Add support for postgres
    - [ ] Auto generate pulumi/terraform providers
    - [ ] Allow for manual access to pulumi resources


### Tools
- https://github.com/golangci/golangci-lint
- https://github.com/air-verse/air
- https://github.com/kisielk/errcheck
- https://github.com/dominikh/go-tools
- https://pkg.go.dev/cmd/vet
- https://taskfile.dev/