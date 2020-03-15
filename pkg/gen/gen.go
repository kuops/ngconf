package gen

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// upstreamTmpl 生成 nginx 的 upstream 字段

var upstreamTmpl = `
{{- $beport := .BackendPort -}}
upstream  {{ .ProjectName }} {
    {{- range $index,$value := .Nodes }}
    server  {{ $value }}:{{ $beport }};
    {{- end }}
    {{- if eq .IsCloud true }}
    include site-enables/medusa-online-service;
    {{- end }}
}
`

// serverTmpl 生成 nginx 的 server 字段
var serverTmpl = `
{{- if eq .ForceHTTPS true }}
server {

	listen 80;
    {{- if eq .Preview true }}
    server_name {{ .ProjectName }}-preview.{{ .Domain }};
	{{ else }}
    server_name {{ .ProjectName }}.{{ .Domain }};
	{{- end }}

    location / {
         return 301 https://{{ .ProjectName }}.{{ .Domain }}$request_uri;
    }

}
{{- end }}

server {
    {{ if eq .ForceHTTPS false }}
    listen 80;
    {{- end }}
    listen 443 ssl;

    ssl_session_timeout  5m;
    ssl_protocols  TLSv1 TLSv1.1 TLSv1.2;
    ssl_ciphers  HIGH:!RC4:!MD5:!aNULL:!eNULL:!NULL:!DH:!EDH:!EXP:+MEDIUM;
    ssl_prefer_server_ciphers   on;

    ssl_certificate     {{ .Cert }};
    ssl_certificate_key {{ .Key }};
    {{- if eq .Preview true }}
    server_name {{ .ProjectName }}-preview.{{ .Domain }};
	{{ else }}
    server_name {{ .ProjectName }}.{{ .Domain }};
	{{- end }}

    location / {
        proxy_http_version  1.1;
        proxy_pass  http://{{ .ProjectName }};
        include  custom_proxy_header.conf;
    }

    access_log  /data/service_logs/nginx/{{ .ProjectName }}_access.log  misc;
    error_log   /data/service_logs/nginx/{{ .ProjectName }}_error.log;
}
`

// Options is ngconf gen options
type Options struct {
	ProjectName   string
	Domains       []string
	IsCloud       bool
	Nodes         []string
	BackendPort   uint
	ForceHTTPS    bool
	WriteFileName string
	Preview       bool
}

// DefaultOptions is generate a default options
func DefaultOptions() *Options {
	return &Options{
		ProjectName:   "example",
		Domains:       []string{"kuops.com"},
		IsCloud:       false,
		Nodes:         []string{},
		BackendPort:   80,
		ForceHTTPS:    false,
		WriteFileName: "",
		Preview:       false,
	}
}

type sslCert struct {
	KuopsSslKey      string
	KuopsSslCert     string
	DnskuopsSslKey   string
	DnskuopsSslCert  string
	KuopscorpSslKey  string
	KuopscorpSslCert string
}

func newSslCert() *sslCert {
	return &sslCert{
		KuopsSslKey:      "/usr/local/openresty/nginx/ssl/server.key",
		KuopsSslCert:     "/usr/local/openresty/nginx/ssl/server.crt",
		DnskuopsSslKey:   "/usr/local/openresty/nginx/ssl/dns_server.key",
		DnskuopsSslCert:  "/usr/local/openresty/nginx/ssl/dns_server.crt",
		KuopscorpSslKey:  "/usr/local/openresty/nginx/ssl/corp.key",
		KuopscorpSslCert: "/usr/local/openresty/nginx/ssl/corp.crt",
	}
}

type serverGenVars struct {
	Domain string
	Key    string
	Cert   string
	Options
}

func upstreamGen(opts *Options) string {
	var upConf bytes.Buffer
	if len(opts.Nodes[0:]) != 0 || opts.IsCloud == true {
		tmpl, err := template.New("upstream").Parse(upstreamTmpl)
		if err != nil {
			fmt.Printf("error is: %v\n", err)
		}
		tmpl.Execute(&upConf, *opts)
	}
	return upConf.String()
}

func serverGen(opts *Options) string {
	var srConf bytes.Buffer
	srtmpl, err := template.New("server").Parse(serverTmpl)
	if err != nil {
		fmt.Printf("error is: %v\n", err)
	}
	serverVars := &serverGenVars{
		Options: *opts,
	}

	for _, v := range opts.Domains {
		serverVars.Domain = v
		switch {
		case v == "kuops.com":
			serverVars.Key = newSslCert().KuopsSslKey
			serverVars.Cert = newSslCert().KuopsSslCert
		case v == "kuops-corp.com":
			serverVars.Key = newSslCert().KuopscorpSslKey
			serverVars.Cert = newSslCert().KuopscorpSslCert
		case v == "dns.kuops.com":
			serverVars.Key = newSslCert().DnskuopsSslKey
			serverVars.Cert = newSslCert().DnskuopsSslCert
		}
		srtmpl.Execute(&srConf, serverVars)
	}
	return srConf.String()
}

func getIP(o *Options) {
	for i, v := range o.Nodes {
		h := v + ".dns.kuops.com"
		ips, err := net.LookupIP(h)
		if err == nil {
			o.Nodes[i] = ips[0].String()
		} else {
			fmt.Printf("\033[31m# dose not have nodes: %v \033[0m\n\n\n", o.Nodes[i])
			o.Nodes = append(o.Nodes[:i], o.Nodes[i+1:]...)
		}
	}
}

func chkUpstreamExists(o *Options) error {
	basepath := filepath.Dir(o.WriteFileName)
	return filepath.Walk(basepath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if strings.Contains(path, ".conf") {
				body, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				r, _ := regexp.Compile(o.ProjectName)
				if r.Match(body) {
					return fmt.Errorf("projectname exists in file %v", path)
				}
			}
		}
		return nil
	})
}

// GenConfig is generate a nginx config to stdout or file
func GenConfig(o *Options) {
	getIP(o)
	ngconf := upstreamGen(o) + serverGen(o)
	if o.WriteFileName != "" {
		if _, err := os.Stat(o.WriteFileName); err == nil {
			fmt.Printf("file %v is exists.\n", o.WriteFileName)
		} else {
			err := chkUpstreamExists(o)
			if err != nil {
				fmt.Printf("err is: %v\n", err)
			} else {
				ioutil.WriteFile(o.WriteFileName, []byte(ngconf), 0644)
			}
		}
	} else {
		fmt.Printf("%v", ngconf)
	}

}
