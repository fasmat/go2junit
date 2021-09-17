# go test output to junit converter

This tool is a small helper to convert the output of `go test` to a junit xml report.

## How to install

```bash
go install github.com/fasmat/go2junit/cmd/go2junit@latest
```

## How to use

```bash
go test -json ./path/to/pkg | go2junit > out.xml
```

## Open points

* [ ] Add CI to github actions
* [ ] Add multiple reference test runs from various projects
* [ ] Handle non-trivial cases that might not be handled yet
* [ ] Allow the use of go2junit as drop-in replacement for go test
* [ ] Add license
* [ ] Improve documentation
