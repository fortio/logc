all: tests manual-check

tests:
	go test -v ./...

manual-check:
	@echo "=========== With Everything:                     ==========="
	go run ./levelsDemo 2>&1 | TZ=UTC go run -race .
	@echo "=========== Without Timestamp nor go Routine ID: ==========="
	go run ./levelsDemo -logger-timestamp=false -logger-goroutine=false 2>&1 | go run -race .
	@echo "=========== Without file/line                    ==========="
	go run ./levelsDemo -logger-file-line=false -logger-timestamp=false -logger-goroutine=false 2>&1 | go run -race .
	@echo "=========== Without Color:                       ==========="
	go run ./levelsDemo 2>&1 | go run -race . -no-color

.PHONY: tests manual-check
