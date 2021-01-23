package crawler

import (
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	log "github.com/sirupsen/logrus"
)

func newDefaultCollector() *colly.Collector {
	return colly.NewCollector(
		// Turn on asynchronous requests
		colly.Async(),
		// Attach a debugger to the collector
		colly.Debugger(&debug.LogDebugger{
			Prefix: "CollyDebug",
			Output: log.StandardLogger().Out,
		}),
	)
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
