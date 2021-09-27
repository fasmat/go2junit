package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

func getReader(name string) (io.ReadCloser, error) {
	if name == "" {
		return os.Stdin, nil
	}

	return os.Open(name)
}

func getWriter(name string) (io.WriteCloser, error) {
	if name == "" {
		return os.Stdout, nil
	}

	return os.Create(name)
}

func actionParse(c *cli.Context) error {
	var writer io.Writer
	var printer io.Writer = io.Discard

	reader, err := getReader(c.String("input"))
	if err != nil {
		return fmt.Errorf("error opening input file: %w", err)
	}
	defer reader.Close()

	w, err := getWriter(c.String("output"))
	if err != nil {
		return fmt.Errorf("error opening input file: %w", err)
	}
	writer = w
	defer w.Close()

	if c.Bool("print") {
		printer = os.Stdout
	}

	if c.Bool("print") && writer == os.Stdout {
		writer = io.Discard
	}

	parse(writer, reader, printer, c.Bool("fail"))
	return nil
}

func actionTest(c *cli.Context) error {
	var writer io.Writer
	var printer io.Writer = io.Discard

	w, err := getWriter(c.String("output"))
	if err != nil {
		return fmt.Errorf("error opening input file: %w", err)
	}
	writer = w
	defer w.Close()

	if c.Bool("print") {
		printer = os.Stdout
	}

	if c.Bool("print") && writer == os.Stdout {
		writer = io.Discard
	}

	parseChan := make(chan error)
	pipeReader, pipeWriter := io.Pipe()
	go func() {
		parse(writer, pipeReader, printer, true)
		close(parseChan)
		pipeReader.Close()
	}()

	testErrChan := make(chan error)
	go func() {
		cmd := exec.Command("go")
		cmd.Args = append(cmd.Args, "test", "-json")
		cmd.Args = append(cmd.Args, c.Args().Slice()...)
		cmd.Stdout = pipeWriter
		cmd.Stderr = c.App.ErrWriter
		err := cmd.Run()

		pipeWriter.Close()
		if err != nil {
			testErrChan <- err
		}
		close(testErrChan)
	}()

	<-parseChan

	for err := range testErrChan {
		return err
	}

	return nil
}
