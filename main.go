package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	dataURL     = "http://reportapi.eastmoney.com/report/list"
	downloadURL = "https://pdf.dfcfw.com/pdf/"
)

func main() {
	now := time.Now().Format("2006-01-02")
	qType := flag.String("query-type", "0,1,2", "研报类型，0：个股研报，1：行业研报，2：宏观研究")
	beginTime := flag.String("begin-time", now, "开始时间")
	endTime := flag.String("end-time", now, "结束时间")
	minPages := flag.Int("min-pages", 15, "最小页数")
	downloadPath := flag.String("download-path", "/tmp", "下载路径")
	flag.Parse()

	for _, qt := range strings.Split(*qType, ",") {
		items, err := list(qt, *beginTime, *endTime, *minPages)
		if err != nil {
			log.Fatalln(err)
		}

		if err = download(qt, *downloadPath, items); err != nil {
			log.Fatalln(err)
		}
	}
}

func list(qType, beginTime, endTime string, minPages int) ([]*Report, error) {
	u, err := url.Parse(dataURL)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Add("pageSize", "100")
	q.Add("qType", qType)
	q.Add("beginTime", beginTime)
	q.Add("endTime", endTime)

	curPage, maxPage := 0, 1
	reports := make([]*Report, 0)

	for curPage <= maxPage {
		q.Del("pageNo")
		q.Add("pageNo", strconv.Itoa(curPage+1))
		u.RawQuery = q.Encode()
		resp, err := http.Get(u.String())
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			log.Fatalln(resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		tmp := new(Reports)
		if err := json.Unmarshal(body, &tmp); err != nil {
			return nil, err
		}
		curPage = tmp.PageNo
		maxPage = tmp.TotalPage

		for _, i := range tmp.Data {
			if i.AttachPages < minPages {
				continue
			}
			reports = append(reports, i)
		}
	}

	return reports, nil
}

func download(qType, downloadPath string, reports []*Report) error {
	downloadPath = strings.TrimSuffix(downloadPath, "/")
	for _, report := range reports {
		report.PublishDate = fixDate(report.PublishDate)
		report.Title = fixTitle(report.Title)
		u := genDownloadURL(report.InfoCode)
		dir, name := "", ""
		switch qType {
		case "0":
			name = fmt.Sprintf("%s-%s-%s-%s.pdf", report.StockName, report.PublishDate, report.OrgSName, report.Title)
			dir = fmt.Sprintf("%s/研报/%s/个股研报/%s/", downloadPath, report.PublishDate[:4], report.IndvInduName)
		case "1":
			name = fmt.Sprintf("%s-%s-%s.pdf", report.PublishDate, report.OrgSName, report.Title)
			dir = fmt.Sprintf("%s/研报/%s/行业研报/%s/", downloadPath, report.PublishDate[:4], report.IndustryName)
		case "2":
			name = fmt.Sprintf("%s-%s-%s.pdf", report.PublishDate, report.OrgSName, report.Title)
			dir = fmt.Sprintf("%s/研报/%s/宏观研究/", downloadPath, report.PublishDate[:4])
		}

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err = os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}
		fmt.Println(dir+name, u)
		resp, err := http.Get(u)
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("unexpect code %d", resp.StatusCode)
		}
		data, err := ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			return err
		}
		if err = ioutil.WriteFile(dir+name, data, 0644); err != nil {
			return err
		}
	}
	return nil
}

func fixTitle(title string) string {
	title = strings.ReplaceAll(title, "/", "")
	title = strings.ReplaceAll(title, "|", "；")
	return title
}

func fixDate(date string) string {
	return strings.ReplaceAll(strings.Fields(date)[0], "-", "")
}

func genDownloadURL(code string) string {
	return fmt.Sprintf("%sH3_%s_1.pdf", downloadURL, code)
}
