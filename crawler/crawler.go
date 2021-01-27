package crawler

import (
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"

	log "github.com/sirupsen/logrus"
)

func newDefaultCollector() *colly.Collector {
	c := colly.NewCollector(
		// Turn on asynchronous requests
		colly.Async(),
	)
	// Limit the number of threads started by colly to two
	c.Limit(&colly.LimitRule{
		DomainRegexp: ".*",
		Parallelism:  4,
	})

	return c
}

type crawlerImpl struct {
	mCollectors map[string]*colly.Collector
	mCallBack   FundCrawlerCallBack
	// create a request queue with 2 consumer threads
	mQueue *queue.Queue
}

func NewFundCrawler() FundCrawler {
	c := &crawlerImpl{
		mCollectors: make(map[string]*colly.Collector),
		mCallBack:   &DefaultCrawlerCallBack{},
	}
	var err error
	c.mQueue, err = queue.New(
		4, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 999999}, // Use default queue storage
	)
	if err != nil {
		log.Fatal(err)
	}

	return c
}

func (cl *crawlerImpl) SetCallBack(cb FundCrawlerCallBack) {
	cl.mCallBack = cb
}
