package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"go2junit/types"
)

func main() {
	f, err := os.Open("test.log")
	if err != nil {
		log.Fatal(err)
	}
	parse(os.Stdout, f)
}

func parse(w io.Writer, r io.Reader) {
	scanner := bufio.NewScanner(r)

	// suites -> suite
	suites := make(map[string]*types.Testsuite)

	// suites -> suite -> testcase
	cases := make(map[string]map[string]*types.Testcase)

	for scanner.Scan() {
		var event types.TestEvent
		if err := json.NewDecoder(bytes.NewBuffer(scanner.Bytes())).Decode(&event); err != nil {
			continue
		}

		if _, found := suites[event.Package]; !found {
			suites[event.Package] = &types.Testsuite{
				NameAttr:      event.Package,
				TimestampAttr: event.Time.String(),
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
				continue
			case "pass":
				fallthrough
			case "fail":
				fallthrough
			case "skip":
				suite.TimeAttr = strconv.FormatFloat(event.Elapsed, 'f', 2, 64)
				continue
			default:
				fmt.Println("unknown package action found:")
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "   ")
				_ = enc.Encode(event)
				return
			}
		}

		if _, found := testcases[event.Test]; !found {
			testcases[event.Test] = &types.Testcase{
				ClassnameAttr: event.Package,
				NameAttr:      event.Test,
			}
		}
		testcase := testcases[event.Test]

		switch event.Action {
		case "run":
			// only indicating that test is run
			continue
		case "output":
			testcase.Systemout.Text += event.Output
		case "pass":
			testcase.TimeAttr = strconv.FormatFloat(event.Elapsed, 'f', 2, 64)
		case "fail":
			testcase.TimeAttr = strconv.FormatFloat(event.Elapsed, 'f', 2, 64)
			testcase.Error = &types.Error{
				TypeAttr:    "Error",
				MessageAttr: "test failed",
			}
		default:
			fmt.Println("unknown test action found:")
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "   ")
			_ = enc.Encode(event)
			return
		}
	}

	testsuites := &types.Testsuites{}
	testtotal := 0
	errortotal := 0

	for _, s := range suites {
		testcount := 0
		errorcount := 0
		for _, c := range cases[s.NameAttr] {
			// c.Systemout = strings.Replace(c.Systemout, "\\n", "\n", -1)
			// c.Systemout = strings.Replace(c.Systemout, "\\t", "\t", -1)
			s.Testcase = append(s.Testcase, c)
			testcount++

			if c.Error != nil {
				errorcount++
			}
		}

		s.TestsAttr = strconv.Itoa(testcount)
		s.ErrorsAttr = strconv.Itoa(errorcount)
		// s.Systemout = strings.Replace(s.Systemout, "\\n", "\n", -1)
		// s.Systemout = strings.Replace(s.Systemout, "\\t", "\t", -1)

		testsuites.Testsuite = append(testsuites.Testsuite, s)
		testtotal += testcount
		errortotal += errorcount
	}

	testsuites.TestsAttr = strconv.Itoa(testtotal)
	testsuites.ErrorsAttr = strconv.Itoa(errortotal)

	err := xml.NewEncoder(os.Stdout).Encode(testsuites)
	if err != nil {
		log.Fatal(err)
	}
}
