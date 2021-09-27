package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/fasmat/go2junit/types"
)

func parse(w io.Writer, r io.Reader, p io.Writer, fail bool) {
	scanner := bufio.NewScanner(r)

	// suites -> suite
	suites := make(map[string]*types.Testsuite)

	// suites -> suite -> testcase
	cases := make(map[string]map[string]*types.Testcase)

	// keeps track if any test failed
	testfailed := false

	for scanner.Scan() {
		var event types.TestEvent
		if err := json.NewDecoder(bytes.NewBuffer(scanner.Bytes())).Decode(&event); err != nil {
			continue
		}

		if _, found := suites[event.Package]; !found {
			suites[event.Package] = &types.Testsuite{
				NameAttr:      event.Package,
				TimestampAttr: event.Time.String(),
				Systemout:     &types.Systemout{},
			}
		}
		suite := suites[event.Package]

		if _, found := cases[event.Package]; !found {
			cases[event.Package] = make(map[string]*types.Testcase)
		}
		testcases := cases[event.Package]

		// action related to the package and not a specific test
		if event.Test == "" {
			switch event.Action {
			case "output":
				suite.Systemout.Text += event.Output
				if _, err := p.Write([]byte(event.Output)); err != nil {
					log.Fatalf("error printing test package output: %v", err)
				}
				continue
			case "fail":
				testfailed = true
				fallthrough
			case "pass":
				fallthrough
			case "skip":
				suite.TimeAttr = strconv.FormatFloat(event.Elapsed, 'f', 2, 64)
				continue
			default:
				log.Printf("unknown package action found: %+v\n", event)
				log.Fatal("failed to parse input, please report this error to github.com/fasmat/go2junit")
			}
		}

		if _, found := testcases[event.Test]; !found {
			testcases[event.Test] = &types.Testcase{
				ClassnameAttr: event.Package,
				NameAttr:      event.Test,
				Systemout:     &types.Systemout{},
			}
		}
		testcase := testcases[event.Test]

		switch event.Action {
		case "run":
			// only indicating that test is run
			continue
		case "output":
			testcase.Systemout.Text += event.Output
			if _, err := p.Write([]byte(event.Output)); err != nil {
				log.Fatalf("error printing test event output: %v", err)
			}
		case "pass":
			testcase.TimeAttr = strconv.FormatFloat(event.Elapsed, 'f', 2, 64)
		case "fail":
			testfailed = true
			testcase.TimeAttr = strconv.FormatFloat(event.Elapsed, 'f', 2, 64)
			testcase.Failure = &types.Failure{
				TypeAttr:    "Error",
				MessageAttr: "test failed",
				Text:        testcase.Systemout.Text,
			}
		default:
			log.Printf("unknown test action found: %+v\n", event)
			log.Fatal("failed to parse input, please report this error to github.com/fasmat/go2junit")
		}
	}

	testsuites := &types.Testsuites{}
	testtotal := 0
	errortotal := 0
	failtotal := 0

	for _, s := range suites {
		testcount := 0
		errorcount := 0
		failcount := 0
		skipcount := 0
		for _, c := range cases[s.NameAttr] {
			s.Testcase = append(s.Testcase, c)
			testcount++

			if c.Error != nil {
				errorcount++
			}

			if c.Failure != nil {
				failcount++
			}

			if c.Skipped != nil {
				skipcount++
			}
		}

		s.TestsAttr = strconv.Itoa(testcount)
		s.ErrorsAttr = strconv.Itoa(errorcount)
		s.FailuresAttr = strconv.Itoa(failcount)
		s.SkippedAttr = strconv.Itoa(skipcount)

		testsuites.Testsuite = append(testsuites.Testsuite, s)
		testtotal += testcount + skipcount
		errortotal += errorcount
		failtotal += failcount
	}

	testsuites.TestsAttr = strconv.Itoa(testtotal)
	testsuites.ErrorsAttr = strconv.Itoa(errortotal)
	testsuites.FailuresAttr = strconv.Itoa(failtotal)

	if err := xml.NewEncoder(w).Encode(testsuites); err != nil {
		log.Println("failed to encode xml")
		log.Fatal(err)
	}

	if fail && testfailed {
		os.Exit(1)
	}
}
