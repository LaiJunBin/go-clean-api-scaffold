GO_EXECUTABLE := /usr/local/go/bin/go

.PHONY: api
api:
	oapi-codegen --config=gen.config.yaml ./api/openapi.yaml

.PHONY: dev
dev:
	air

.PHONY: build
build:
	go build -o bin/app ./cmd/main.go

.PHONY: up
up:
	docker-compose up -d

.PHONY: down
down:
	docker-compose down

.PHONY: mocks
mocks:
	find internal/app -type d -name "mocks" | xargs rm -rf 2>/dev/null; true; \
	mockery --config "" --all --outpkg sharedmocks --dir internal/app/shared --output internal/app/shared/mocks; \
	for domain_dir in internal/app/*_domain/*/; do \
		mockery --config "" --all --dir "$$domain_dir" --output "$${domain_dir}mocks"; \
	done

TEST_FLAGS := -v -count=1
TEST_FEATURE_PACKAGES := ./tests/features/...

.PHONY: test\:unit
test\:unit:
	go test $(TEST_FLAGS) ./internal/... -failfast

.PHONY: test\:feature
test\:feature:
	go run tests/test.go -count=1 $(TEST_FEATURE_PACKAGES)

.PHONY: migrate
migrate:
	go run ./cmd/migrate/main.go
