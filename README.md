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

NAME:
   go2junit parse - parse input from `FILE` (defaults to stdin if not set)

USAGE:
   go2junit parse [command options] [arguments...]

OPTIONS:
   --input FILE, -i FILE   parse input from FILE (defaults to stdin if not set)
   --output FILE, -o FILE  write output to FILE (defaults to stdout if not set)
   --print, -p             print test log to stdout, requires --output to be set or junit report will be discarded (default: false)
   --fail                  return with a non-zero exit status in the case a parsed test failed (default: false)
   --help, -h              show help (default: false)
```

## How to use go2junit as your test runner

`go2unit` can be used as a test runner. It will parse the output of the go test command and convert it to a junit xml report.

```bash
go2junit test -- [arguments to ``go test``] > out.xml
```

Again by default `go2junit` will write to stdout. You can alter that behavior and specify a file to write to. For details see:

```bash
go2junit test -h

NAME:
   go2junit test - execute go test and output result as junit xml report

USAGE:
   go2junit test -- [arguments...]
   
      the test subcommand executes 'go test' with the given arguments and directly converts its output to junit xml

OPTIONS:
   --output FILE, -o FILE  write output to FILE (defaults to stdout if not set)
   --print, -p             print test log to stdout, requires --output to be set or junit report will be discarded (default: false)
   --help, -h              show help (default: false)
```

## Open points

* Add unit tests
* Add reference data from diverse test runs of different projects
