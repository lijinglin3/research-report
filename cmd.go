package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	types              []string
	beginTime, endTime string
	minPages           int

	now     = time.Now().Format("2006-01-02")
	rootCmd = &cobra.Command{Use: "research-report"}
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
	rootCmd.PersistentFlags().IntVarP(&minPages, "min-pages", "m", 15, "min pages limit")
}

func runList(_ *cobra.Command, _ []string) {
	for _, qt := range types {
		items, err := list(qt, beginTime, endTime, minPages)
		if err != nil {
			log.Fatalln(err)
		}
		for _, i := range items {
			fmt.Println(i)
		}
	}
}

func runDownload(_ *cobra.Command, args []string) {
	downloadPath := "./"
	if len(args) != 0 {
		downloadPath = args[0]
		if !strings.HasSuffix(downloadPath, "/") {
			downloadPath = downloadPath + "/"
		}
	}
	for _, qt := range types {
		items, err := list(qt, beginTime, endTime, minPages)
		if err != nil {
			log.Fatalln(err)
		}

		if err = download(downloadPath, items); err != nil {
			log.Fatalln(err)
		}
	}
}

func execute() error {
	rootCmd.AddCommand(listCmd, downloadCmd)
	return rootCmd.Execute()
}
