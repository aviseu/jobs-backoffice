lint:
	C:\ProgramData\chocolatey\bin\golangci-lint.exe run

test:
	go test -v ./...

start:
	docker-compose up -d

stop:
	docker-compose down

build-import:
	go build -ldflags "-s -w" -ldflags "-X main.version=${VERSION}" -o "dist/app" github.com/aviseu/jobs/cmd/import
