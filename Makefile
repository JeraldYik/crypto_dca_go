run:
	go run main.go

dev:
	source conf/dev.env && go run main.go

test:
	go test -count=1 ./...