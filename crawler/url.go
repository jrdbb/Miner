package crawler

type TargetURL string

const (
	FundCodeSearch TargetURL = "FundCodeSearch"
)

var URLCenter map[TargetURL]string

func init() {
	URLCenter = make(map[TargetURL]string)
	URLCenter[FundCodeSearch] = "http://fund.eastmoney.com/js/fundcode_search.js"
}
