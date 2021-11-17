# research-report

[![build](https://github.com/lijinglin3/research-report/actions/workflows/build.yml/badge.svg)](https://github.com/lijinglin3/research-report/actions/workflows/build.yml)

研报下载助手

## Usage

```text
$ rr-cli -h
Usage:
  rr-cli [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  download    Download research reports
  help        Help about any command
  list        List research reports

Flags:
  -b, --begin string    begin time (default "2021-11-17")
  -e, --end string      end time (default "2021-11-17")
  -h, --help            help for rr-cli
  -m, --min-pages int   min pages limit (default 20)
  -t, --type strings    report type, 0: individual stocks, 1: industry, 2: macro (default [0,1,2])

Use "rr-cli [command] --help" for more information about a command.
```
