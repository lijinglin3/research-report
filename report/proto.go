package report

import (
	"errors"
	"fmt"
	"strings"
)

const downloadURL = "https://pdf.dfcfw.com/pdf/"

// ErrUnknownType unknown type error
var ErrUnknownType = errors.New("unknown report type")

// Type report type
type Type string

const (
	reportTypeIndividualStock = "0"
	reportTypeIndustry        = "1"
	reportTypeMacro           = "2"
)

func (r Type) String() string {
	switch r {
	case reportTypeIndividualStock:
		return "个股研报"
	case reportTypeIndustry:
		return "行业研报"
	case reportTypeMacro:
		return "宏观研究"
	default:
		panic(ErrUnknownType)
	}
}

type reports struct {
	Hits      int `json:"hits"`
	Size      int `json:"size"`
	TotalPage int `json:"TotalPage"`
	PageNo    int `json:"pageNo"`

	Data []*rawReport `json:"data"`
}

type rawReport struct {
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
	ReportType            int         `json:"ReportType"`
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

// Report structure
type Report struct {
	Type      Type   `json:"type"`
	Date      string `json:"date"`
	DateShort string `json:"short_date"`
	DateMonth string `json:"date_month"`
	Title     string `json:"title"`
	Path      string `json:"path"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	Org       string `json:"org"`
	Industry  string `json:"industry"`
	Stock     string `json:"stock"`
}

func (r *rawReport) convert(qType string) *Report {
	ret := new(Report)
	ret.Type = Type(qType)
	ret.Date = r.PublishDate[:10]
	ret.DateShort = strings.ReplaceAll(r.PublishDate[:10], "-", "")
	ret.URL = fmt.Sprintf("%sH3_%s_1.pdf", downloadURL, r.InfoCode)
	ret.Title = fixTitle(r.Title)
	ret.Org = r.OrgSName
	ret.Industry = r.IndvInduName + r.IndustryName
	ret.Stock = r.StockName

	if ret.Stock != "" {
		ret.Path = fmt.Sprintf("研报/%s/%s/", ret.Type, ret.Stock)
		ret.Name = fmt.Sprintf("%s-%s-%s-%s.pdf", ret.DateShort, ret.Org, ret.Stock, ret.Title)
	} else {
		ret.Path = fmt.Sprintf("研报/%s/%s/", ret.Type, ret.Date[:7])
		ret.Name = fmt.Sprintf("%s-%s-%s.pdf", ret.DateShort, ret.Org, ret.Title)
	}
	return ret
}

func fixTitle(title string) string {
	title = strings.ReplaceAll(title, "/", "")
	title = strings.ReplaceAll(title, "|", "；")
	return title
}
