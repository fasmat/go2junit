package types

import (
	"encoding/xml"
)

// The schema we are implementing is the one used by jenkins:
//    https://llg.cubic.org/docs/junit/

type Failure struct {
	XMLName     xml.Name `xml:"failure"`
	TypeAttr    string   `xml:"type,attr,omitempty"`
	MessageAttr string   `xml:"message,attr,omitempty"`
	Text        string   `xml:",cdata"`
}

type Error struct {
	XMLName     xml.Name `xml:"error"`
	TypeAttr    string   `xml:"type,attr,omitempty"`
	MessageAttr string   `xml:"message,attr,omitempty"`
	Text        string   `xml:",cdata"`
}

type Skipped struct {
	XMLName     xml.Name `xml:"skipped"`
	TypeAttr    string   `xml:"type,attr,omitempty"`
	MessageAttr string   `xml:"message,attr,omitempty"`
}

type Properties struct {
	XMLName  xml.Name    `xml:"properties"`
	Property []*Property `xml:"property"`
}

type Property struct {
	XMLName   xml.Name `xml:"property"`
	NameAttr  string   `xml:"name,attr"`
	ValueAttr string   `xml:"value,attr"`
}

type Systemerr struct {
	XMLName xml.Name `xml:"system-err"`
	Text    string   `xml:",cdata"`
}

type Systemout struct {
	XMLName xml.Name `xml:"system-out"`
	Text    string   `xml:",cdata"`
}

type Testcase struct {
	XMLName       xml.Name   `xml:"testcase"`
	NameAttr      string     `xml:"name,attr"`
	TimeAttr      string     `xml:"time,attr,omitempty"`
	ClassnameAttr string     `xml:"classname,attr,omitempty"`
	GroupAttr     string     `xml:"group,attr,omitempty"`
	Skipped       *Skipped   `xml:"skipped"`
	Error         *Error     `xml:"error"`
	Failure       *Failure   `xml:"failure"`
	Systemout     *Systemout `xml:"system-out,omitempty"`
	Systemerr     *Systemerr `xml:"system-err,omitempty"`
}

type Testsuite struct {
	XMLName       xml.Name    `xml:"testsuite"`
	NameAttr      string      `xml:"name,attr"`
	TestsAttr     string      `xml:"tests,attr"`
	FailuresAttr  string      `xml:"failures,attr"`
	ErrorsAttr    string      `xml:"errors,attr"`
	GroupAttr     string      `xml:"group,attr,omitempty"`
	TimeAttr      string      `xml:"time,attr,omitempty"`
	SkippedAttr   string      `xml:"skipped,attr,omitempty"`
	TimestampAttr string      `xml:"timestamp,attr,omitempty"`
	HostnameAttr  string      `xml:"hostname,attr,omitempty"`
	IdAttr        string      `xml:"id,attr,omitempty"`
	PackageAttr   string      `xml:"package,attr,omitempty"`
	FileAttr      string      `xml:"file,attr,omitempty"`
	LogAttr       string      `xml:"log,attr,omitempty"`
	UrlAttr       string      `xml:"url,attr,omitempty"`
	VersionAttr   string      `xml:"version,attr,omitempty"`
	Testsuite     *Testsuite  `xml:"testsuite"`
	Properties    *Properties `xml:"properties"`
	Testcase      []*Testcase `xml:"testcase"`
	Systemout     *Systemout  `xml:"system-out,omitempty"`
	Systemerr     *Systemerr  `xml:"system-err,omitempty"`
}

type Testsuites struct {
	XMLName      xml.Name     `xml:"testsuites"`
	NameAttr     string       `xml:"name,attr,omitempty"`
	TimeAttr     string       `xml:"time,attr,omitempty"`
	TestsAttr    string       `xml:"tests,attr,omitempty"`
	FailuresAttr string       `xml:"failures,attr,omitempty"`
	ErrorsAttr   string       `xml:"errors,attr,omitempty"`
	Testsuite    []*Testsuite `xml:"testsuite"`
}
