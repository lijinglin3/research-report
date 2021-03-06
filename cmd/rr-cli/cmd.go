package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/lijinglin3/research-report/report"
)

var (
	types              []string
	beginTime, endTime string
	minPages           int

	now     = time.Now().Format("2006-01-02")
	rootCmd = &cobra.Command{Use: "rr-cli"}
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List research reports",
		Args:  cobra.NoArgs,
		Run:   runList,
	}
	downloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download research reports",
		Args:  cobra.MaximumNArgs(1),
		Run:   runDownload,
	}
)

func init() {
	rootCmd.PersistentFlags().StringSliceVarP(&types, "type", "t", []string{"0", "1", "2"},
		"report type, 0: individual stocks, 1: industry, 2: macro")
	rootCmd.PersistentFlags().StringVarP(&beginTime, "begin", "b", now, "begin time")
	rootCmd.PersistentFlags().StringVarP(&endTime, "end", "e", now, "end time")
	rootCmd.PersistentFlags().IntVarP(&minPages, "min-pages", "m", 10, "min pages limit")
}

func runList(_ *cobra.Command, _ []string) {
	sort.Sort(sort.Reverse(sort.StringSlice(types)))
	for _, qt := range types {
		items, err := report.List(qt, beginTime, endTime, minPages)
		if err != nil {
			log.Fatalln(err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.SetStyle(table.StyleRounded)
		t.SetTitle("%s", table.Row{report.Type(qt)})
		switch qt {
		case "0":
			t.AppendHeader(table.Row{"日期", "券商", "行业", "股票", "标题", "链接"})
			for _, i := range items {
				t.AppendSeparator()
				t.AppendRow([]interface{}{i.Date, i.Org, i.Industry, i.Stock, i.Title, i.URL})
			}
		case "1":
			t.AppendHeader(table.Row{"日期", "券商", "行业", "标题", "链接"})
			for _, i := range items {
				t.AppendSeparator()
				t.AppendRow([]interface{}{i.Date, i.Org, i.Industry, i.Title, i.URL})
			}
		case "2":
			t.AppendHeader(table.Row{"日期", "券商", "标题", "链接"})
			for _, i := range items {
				t.AppendSeparator()
				t.AppendRow([]interface{}{i.Date, i.Org, i.Title, i.URL})
			}
		default:
			panic(report.ErrUnknownType)
		}

		t.Render()
	}
}

func runDownload(_ *cobra.Command, args []string) {
	sort.Sort(sort.Reverse(sort.StringSlice(types)))
	downloadPath := "./"
	if len(args) != 0 {
		downloadPath = args[0]
		if !strings.HasSuffix(downloadPath, "/") {
			downloadPath = downloadPath + "/"
		}
	}
	for _, qt := range types {
		items, err := report.List(qt, beginTime, endTime, minPages)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("开始下载 %s\n", report.Type(qt))
		if err = report.Download(downloadPath, items); err != nil {
			log.Fatalln(err)
		}
		fmt.Println()
	}
}

func execute() error {
	rootCmd.AddCommand(listCmd, downloadCmd)
	return rootCmd.Execute()
}
