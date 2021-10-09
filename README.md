# research-report

[![build](https://github.com/lijinglin3/research-report/actions/workflows/build.yml/badge.svg)](https://github.com/lijinglin3/research-report/actions/workflows/build.yml)

研报下载助手

## Usage

```text
$ research-report -h
Usage:
  research-report [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  download    Download research reports
  help        Help about any command
  list        List research reports

Flags:
  -b, --begin string    begin time (default "2021-10-10")
  -e, --end string      end time (default "2021-10-10")
  -h, --help            help for research-report
  -m, --min-pages int   min pages limit (default 15)
  -t, --type strings    report type, 0: individual stocks, 1: industry, 2: macro (default [0,1,2])

Use "research-report [command] --help" for more information about a command.
```
