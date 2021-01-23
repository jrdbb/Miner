package crawler

type BasicFund struct {
	ID        string
	ShortCode string

	// TODO: could it be enum?
	Type   string
	NameCN string
	// pin yin name
	NamePY string
}

//go:generate mockgen -destination=mocks/types.go -package=mocks . FundCrawlerCallBack,FundCrawler
type FundCrawlerCallBack interface {
	OnBasicFund([]*BasicFund)
}

type DefaultCrawlerCallBack struct {
}

func (cb *DefaultCrawlerCallBack) OnBasicFund([]*BasicFund) {}

type FundCrawler interface {
	SetCallBack(FundCrawlerCallBack)
	GetAllBasicFund(sync bool)
}
