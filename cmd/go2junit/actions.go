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
	r, err := getReader(c.String("input"))
	if err != nil {
		return fmt.Errorf("error opening input file: %w", err)
	}
	defer r.Close()

	w, err := getWriter(c.String("output"))
	if err != nil {
		return fmt.Errorf("error opening input file: %w", err)
	}
	defer w.Close()

	parse(w, r, c.Bool("fail"))
	return nil
}

func actionTest(c *cli.Context) error {
	w, err := getWriter(c.String("output"))
	if err != nil {
		return fmt.Errorf("error opening input file: %w", err)
	}
	defer w.Close()

	parseChan := make(chan error)
	pipeReader, pipeWriter := io.Pipe()
	go func() {
		parse(w, pipeReader, true)
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
