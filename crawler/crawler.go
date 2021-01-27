package crawler

import (
	"time"

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
		Parallelism: 5,
		Delay:      1 * time.Millisecond,
		RandomDelay: 100 * time.Millisecond,
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
