# go test output to junit converter

This tool is a small helper to convert the output of `go test` to a junit xml report.

## How to install

```bash
go install github.com/fasmat/go2junit/cmd/go2junit@latest
```

## How to use the parser

Execute your go tests with `-json` output and pipe the output to the converter:

```bash
go test -json ./path/to/pkg | go2junit parse > out.xml
```

By default `go2junit` will read from stdin and write to stdout. You can also specify a file to read from and a file to write to. For details see:

```bash
go2junit parse -h
```

## How to use go2junit as your test runner

`go2unit` can be used as a test runner. It will parse the output of the go test command and convert it to a junit xml report.

```bash
go2junit test -- [arguments to ``go test``] > out.xml
```

Again by default `go2junit` will write to stdout. You can alter that behavior and specify a file to write to. For details see:

```bash
go2junit test -h
```

## Open points

* Add CI to github actions
* Add reference data from diverse test runs of different projects
* Handle non-trivial cases that might not be handled yet
* Improve documentation
