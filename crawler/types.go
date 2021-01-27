package crawler

import "time"

type BasicFund struct {
	ID        string
	ShortCode string

	// TODO: could it be enum?
	NameCN string
	Type   string
	// pin yin name
	NamePY string
}

//go:generate mockgen -destination=mocks/types.go -package=mocks . FundCrawlerCallBack,FundCrawler
type FundCrawlerCallBack interface {
	OnBasicFund([]*BasicFund)
	OnHistoryValue(*ApiData)
}

type DefaultCrawlerCallBack struct {
}

func (cb *DefaultCrawlerCallBack) OnBasicFund([]*BasicFund) {}
func (cb *DefaultCrawlerCallBack) OnHistoryValue(*ApiData)  {}

type FundCrawler interface {
	SetCallBack(FundCrawlerCallBack)
	GetAllBasicFund(sync bool)
	GetHistoryValue(sync bool, code string, page int, sdate string, edate string)
	ConsumeHistoryValueQueue()
}

type FundValue struct {
	Date          time.Time
	NetAssetValue float64
	DailyDelta    float64
	BuyState      string
	SellState     string
}
