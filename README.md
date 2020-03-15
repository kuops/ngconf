# ngconf

# 源代码获取

```
go get -v github.com/kuops/ngconf
```

# 编译

```
# mac
make build-darwin

# linux
make build-linux

# all
make all

# clean
make clean
```

# 使用说明

生成 nginx 配置

```
genterate proxy nginx config files

Usage:
  ngconf gen [--cloud=bool] [--domains="domain1,domain2,..."] [--https=bool] [--nodes="node1","node2"] [--port=int] --projectname="value" [--file=demo.conf] [--preview=bool]

Examples:

        ngconf gen --cloud=false --domains=kuops.com,kuops.io --https=false --nodes=g1-sre-jenkins-v01,g1-sre-jenkins-v02 --port=8080 --projectname=demo --file=demo.conf --preview=true


Flags:
      --cloud                is cloud project include medusa-online-service
      --domains strings      domain name (default [a.com])
      --file string          write to files
  -h, --help                 help for gen
      --https                use https redirect
      --nodes strings        nginx upstream backend server
      --port uint            port is backend port (default 80)
      --preview              is preview config files
      --projectname string   projectname is git repo name, use to domain prefix. (default "example")
```


生成 healthcheck 配置

```
add healthcheck to healthcheck.conf

Usage:
  ngconf hck --upstream=upstream_name --domain=healthcheck_domain [--uri=uri] --hckfile=configfile

Examples:
ngconf hck  --upstream=demo --domain=demo.kuops.com --uri=/healthCheck --hckfile=healthcheck.conf

Flags:
      --domain string     domain is healthcheck header Host field
      --hckfile string    hckfile is healthcheck.conf file path
  -h, --help              help for hck
      --upstream string   upstram is upstream_name same common as projectname.
      --uri string        hckfile is healthcheck.conf file path
```
