all: tests manual-check

tests:
	go test -v ./...

manual-check:
	go run ./levelsDemo 2>&1 | TZ=UTC go run -race .

.PHONY: tests manual-check
