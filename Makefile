.PHONY: audit
audit:
	go fmt ./...
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest -show verbose ./...
	go test -race -vet=off ./...

.PHONY: code-coverage
code-coverage:
	go test -v -coverprofile /tmp/cover.out ./...
	go tool cover -html /tmp/cover.out -o /tmp/cover.html
	xdg-open /tmp/cover.html

.PHONY: benchmark
benchmark:
	go test -benchmem -bench=. ./...

.PHONY: loki-up
loki-up:
	docker-compose -f docker-compose.loki.yml up -d
	@echo "🚀 Grafana Loki stack started!"
	@echo "📊 Grafana: http://localhost:3000 (admin/admin)"
	@echo "🔍 Loki: http://localhost:3100"
	@echo "📝 Run 'make loki-logs' to see logs"

.PHONY: loki-down
loki-down:
	docker-compose -f docker-compose.loki.yml down -v
	@echo "🛑 Grafana Loki stack stopped and cleaned up!"

.PHONY: loki-logs
loki-logs:
	docker-compose -f docker-compose.loki.yml logs -f

.PHONY: loki-status
loki-status:
	docker-compose -f docker-compose.loki.yml ps
