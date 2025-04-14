lint:
	C:\ProgramData\chocolatey\bin\golangci-lint.exe run

test:
	go test ./... -parallel 10

start:
	docker-compose up -d

stop:
	docker-compose down

build-import:
	go build -ldflags "-s -w" -ldflags "-X main.version=${VERSION}" -o "dist/app" github.com/aviseu/jobs-backoffice/cmd/import

build-api:
	go build -ldflags "-s -w" -ldflags "-X main.version=${VERSION}" -o "dist/app" github.com/aviseu/jobs-backoffice/cmd/api

build-schedule:
	go build -ldflags "-s -w" -ldflags "-X main.version=${VERSION}" -o "dist/app" github.com/aviseu/jobs-backoffice/cmd/schedule

migrate-create:
	sh -c "migrate create -ext sql -dir config/migrations -seq $(name)"

migrate-up:
	sh -c "migrate -path config/migrations -database postgres://jobs:pwd@localhost:5433/jobs?sslmode=disable up"

migrate-down:
	sh -c "migrate -path config/migrations -database postgres://jobs:pwd@localhost:5433/jobs?sslmode=disable down"
