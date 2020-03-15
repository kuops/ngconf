package hck

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// Options is hck options
type Options struct {
	Upstream string
	Domain   string
	HckFile  string
	URI      string
}

var healthTmpl = `    ok, err = hc.spawn_checker {
	shm = "healthcheck",
	upstream = "{{ .Upstream }}",
	type = "http",
	http_req = "GET {{ .URI }} HTTP/1.1\r\nHost: {{ .Domain }}\r\n\r\n",
	interval = 2000,
	timeout = 2000,
	fall = 3,
	rise = 2,
	valid_statuses = {200},
	concurrency = 30,
    }`

var appendBuff bytes.Buffer
var healthCheckConfig []string

func hckExists(o *Options) error {
	if o.HckFile == "" {
		return fmt.Errorf("need option --hckfile")
	}
	_, err := os.Stat(o.HckFile)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadFile(o.HckFile)
	if err != nil {
		return err
	}
	r, err := regexp.Compile("upstream = \"" + o.Upstream + "\",")
	if err != nil {
		return err
	}
	if r.Match(body) {
		return fmt.Errorf("upstream %v exist in file %v", o.Upstream, o.HckFile)
	}
	return nil
}

// AppendHck is add healthcheckconf to healthcheck.conf
func AppendHck(o *Options) error {
	err := hckExists(o)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	tmpl, err := template.New("hck").Parse(healthTmpl)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	tmpl.Execute(&appendBuff, *o)

	r, _ := regexp.Compile("local hc = require \"resty.upstream.healthcheck\"")

	body, _ := ioutil.ReadFile(o.HckFile)

	lines := strings.Split(string(body), "\n")

	insertLineIndex := 0

	for i := range lines {
		if r.MatchString(lines[i]) {
			insertLineIndex = i + 1
		}
	}

	frontField := lines[:insertLineIndex]
	lastField := lines[insertLineIndex:]
	appendField := strings.Split(appendBuff.String(), "\n")

	result := []string{}

	result = append(result, frontField...)
	result = append(result, appendField...)
	result = append(result, lastField...)

	resultstr := strings.Join(result, "\n")
	resultbyte := []byte(resultstr)

	ioutil.WriteFile(o.HckFile, resultbyte, 0644)
	return nil
}
