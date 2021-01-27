package crawler

import (

	"github.com/gocolly/colly/v2"
)

func newDefaultCollector() *colly.Collector {
	c := colly.NewCollector(
		// Turn on asynchronous requests
		colly.Async(),
	)
	// Limit the number of threads started by colly to two
	c.Limit(&colly.LimitRule{
		DomainRegexp: ".*",
		Parallelism:  1,
	})

	return c
}

type crawlerImpl struct {
	mCollectors map[string]*colly.Collector
	mCallBack   FundCrawlerCallBack
}

func NewFundCrawler() FundCrawler {
	return &crawlerImpl{
		mCollectors: make(map[string]*colly.Collector),
		mCallBack:   &DefaultCrawlerCallBack{},
	}
}

func (cl *crawlerImpl) SetCallBack(cb FundCrawlerCallBack) {
	cl.mCallBack = cb
}
