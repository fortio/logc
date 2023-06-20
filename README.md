# logc
Converts [fortio.org/log](https://github.com/fortio/log#log) JSON structured log back to console/text output with colors

## Example
```
go run ./levelsDemo 2>&1 | logc
```

![Example console color output](example.png)

If you don't want colors; pass `-no-color`

## Installation

If you have a recent go installation already:
```shell
CGO_ENABLED=0 go install fortio.org/logc@latest
```

Or get one of the [binary releases](https://github.com/fortio/logc/releases)

Or using the docker image
```shell
docker run fortio/logc
```

Or using brew (mac)
```shell
brew install fortio/tap/logc
```

# Development

Run make for both go tests and human check colorization example
```
make
```
