package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const dataURL = "http://reportapi.eastmoney.com/report/list"

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
			i.fill(qType)
			reports = append(reports, i)
		}
	}

	return reports, nil
}

func download(downloadPath string, reports []*Report) error {
	for _, report := range reports {
		fmt.Println(report)
		tmpDir, dir := "/tmp/"+report.DownloadPath, downloadPath+report.DownloadPath
		if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
			if err = os.MkdirAll(tmpDir, 0755); err != nil {
				return err
			}
		}
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err = os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}
		resp, err := http.Get(report.DownloadURL)
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
		if err = ioutil.WriteFile(tmpDir+report.DownloadName, data, 0644); err != nil {
			return err
		}
		if err = os.Rename(tmpDir+report.DownloadName, dir+report.DownloadName); err != nil {
			return err
		}
	}
	return nil
}
