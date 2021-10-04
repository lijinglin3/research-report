package main

type Reports struct {
	Hits      int `json:"hits"`
	Size      int `json:"size"`
	TotalPage int `json:"TotalPage"`
	PageNo    int `json:"pageNo"`

	Data []*Report `json:"data"`
}

type Report struct {
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
	EncodeUrl             string      `json:"encodeUrl"`
	SRatingName           string      `json:"sRatingName"`
	SRatingCode           string      `json:"sRatingCode"`
	Market                string      `json:"market"`
	AuthorID              interface{} `json:"authorID"`
	Count                 int         `json:"count"`
	OrgType               string      `json:"orgType"`
}
