package crawler

import "net/url"

type TargetURL string

const (
	FundCodeSearch TargetURL = "FundCodeSearch"
	HistoryValue   TargetURL = "HistoryValue"
)

var URLCenter map[TargetURL]*url.URL

func init() {
	URLCenter = make(map[TargetURL]*url.URL)
	URLCenter[FundCodeSearch] = &url.URL{
		Scheme: "http",
		Host:   "fund.eastmoney.com",
		Path:   "js/fundcode_search.js",
	}

	URLCenter[HistoryValue] = &url.URL{
		Scheme: "http",
		Host:   "fund.eastmoney.com",
		Path:   "f10/F10DataApi.aspx",
	}
}
