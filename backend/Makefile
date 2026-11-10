.PHONY: run build test tidy db-up db-down
run:   ; go run ./cmd/server
build: ; CGO_ENABLED=0 go build -o bin/server ./cmd/server
test:  ; go test ./...
tidy:  ; go mod tidy
db-up: ; docker compose -f ../docker-compose.yml up -d db
db-down: ; docker compose -f ../docker-compose.yml down
