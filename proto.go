package main

import (
	"fmt"
	"strings"
)

const downloadURL = "https://pdf.dfcfw.com/pdf/"

type reports struct {
	Hits      int `json:"hits"`
	Size      int `json:"size"`
	TotalPage int `json:"TotalPage"`
	PageNo    int `json:"pageNo"`

	Data []*report `json:"data"`
}

type report struct {
	// Custom
	DownloadType  string `json:"download_type"`
	DownloadDate  string `json:"download_date"`
	DownloadTitle string `json:"download_title"`
	DownloadPath  string `json:"download_path"`
	DownloadName  string `json:"download_name"`
	DownloadURL   string `json:"download_url"`

	// Origin
	Title                 string      `json:"title"`
	StockName             string      `json:"stockName"`
	StockCode             string      `json:"stockCode"`
	OrgCode               string      `json:"orgCode"`
	OrgName               string      `json:"orgName"`
	OrgSName              string      `json:"orgSName"`
	PublishDate           string      `json:"publishDate"`
	InfoCode              string      `json:"infoCode"`
	Column                string      `json:"column"`
	PredictNextTwoYearEps string      `json:"predictNextTwoYearEps"`
	PredictNextTwoYearPe  string      `json:"predictNextTwoYearPe"`
	PredictNextYearEps    string      `json:"predictNextYearEps"`
	PredictNextYearPe     string      `json:"predictNextYearPe"`
	PredictThisYearEps    string      `json:"predictThisYearEps"`
	PredictThisYearPe     string      `json:"predictThisYearPe"`
	PredictLastYearEps    string      `json:"predictLastYearEps"`
	PredictLastYearPe     string      `json:"predictLastYearPe"`
	ActualLastTwoYearEps  string      `json:"actualLastTwoYearEps"`
	ActualLastYearEps     string      `json:"actualLastYearEps"`
	IndustryCode          string      `json:"industryCode"`
	IndustryName          string      `json:"industryName"`
	EmIndustryCode        string      `json:"emIndustryCode"`
	IndvInduCode          string      `json:"indvInduCode"`
	IndvInduName          string      `json:"indvInduName"`
	EmRatingCode          string      `json:"emRatingCode"`
	EmRatingValue         string      `json:"emRatingValue"`
	EmRatingName          string      `json:"emRatingName"`
	LastEmRatingCode      string      `json:"lastEmRatingCode"`
	LastEmRatingValue     string      `json:"lastEmRatingValue"`
	LastEmRatingName      string      `json:"lastEmRatingName"`
	RatingChange          interface{} `json:"ratingChange"`
	ReportType            int         `json:"reportType"`
	Author                interface{} `json:"author"`
	IndvIsNew             string      `json:"indvIsNew"`
	Researcher            string      `json:"researcher"`
	NewListingDate        string      `json:"newListingDate"`
	NewPurchaseDate       string      `json:"newPurchaseDate"`
	NewIssuePrice         interface{} `json:"newIssuePrice"`
	NewPeIssueA           interface{} `json:"newPeIssueA"`
	IndvAimPriceT         string      `json:"indvAimPriceT"`
	IndvAimPriceL         string      `json:"indvAimPriceL"`
	AttachType            string      `json:"attachType"`
	AttachSize            int         `json:"attachSize"`
	AttachPages           int         `json:"attachPages"`
	EncodeURL             string      `json:"encodeUrl"`
	SRatingName           string      `json:"sRatingName"`
	SRatingCode           string      `json:"sRatingCode"`
	Market                string      `json:"market"`
	AuthorID              interface{} `json:"authorID"`
	Count                 int         `json:"count"`
	OrgType               string      `json:"orgType"`
}

func (r *report) String() string {
	return fmt.Sprintf("%s %s", r.DownloadName, r.DownloadURL)
}

func (r *report) fill(qType string) {
	r.DownloadType = qType
	r.DownloadDate = strings.ReplaceAll(r.PublishDate[:10], "-", "")
	r.DownloadURL = fmt.Sprintf("%sH3_%s_1.pdf", downloadURL, r.InfoCode)
	r.DownloadTitle = fixTitle(r.Title)

	switch qType {
	case "0":
		r.DownloadName = fmt.Sprintf("%s-%s-%s-%s.pdf", r.StockName, r.DownloadDate, r.OrgSName, r.DownloadTitle)
		r.DownloadPath = fmt.Sprintf("研报/%s/个股研报/%s/", r.DownloadDate[:4], r.IndvInduName)
	case "1":
		r.DownloadName = fmt.Sprintf("%s-%s-%s.pdf", r.DownloadDate, r.OrgSName, r.DownloadTitle)
		r.DownloadPath = fmt.Sprintf("研报/%s/行业研报/%s/", r.DownloadDate[:4], r.IndustryName)
	case "2":
		r.DownloadName = fmt.Sprintf("%s-%s-%s.pdf", r.DownloadDate, r.OrgSName, r.DownloadTitle)
		r.DownloadPath = fmt.Sprintf("研报/%s/宏观研究/", r.DownloadDate[:4])
	default:
		panic("unknown type")
	}
}

func fixTitle(title string) string {
	title = strings.ReplaceAll(title, "/", "")
	title = strings.ReplaceAll(title, "|", "；")
	return title
}
