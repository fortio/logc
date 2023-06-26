all: tests manual-check

tests:
	go test -v ./...

manual-check:
	@echo "======= With Timestamp: ==========="
	go run ./levelsDemo 2>&1 | TZ=UTC go run -race .
	@echo "======= Without: =================="
	go run ./levelsDemo -logger-timestamp=false 2>&1 | TZ=UTC go run -race .

.PHONY: tests manual-check
