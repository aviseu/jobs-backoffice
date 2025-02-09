lint:
	C:\ProgramData\chocolatey\bin\golangci-lint.exe run

test:
	go test -v ./...

start:
	docker-compose up -d

stop:
	docker-compose down
