package tests

import (
	"net"
	"net/http"
	"net/url"

	"github.com/jrdbb/Miner/crawler"
	"github.com/jrdbb/Miner/crawler/mocks"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

// this is not a beautiful fixture
// the fucking colly doesn't offer an interface for its controller
type fixture struct {
	fundCrawler  crawler.FundCrawler
	fundCallback *mocks.MockFundCrawlerCallBack
	server       *http.ServeMux
	listener     net.Listener
}

func newFixture(ctrl *gomock.Controller) *fixture {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	f := &fixture{
		fundCrawler:  crawler.NewFundCrawler(),
		fundCallback: mocks.NewMockFundCrawlerCallBack(ctrl),
		server:       http.NewServeMux(),
		listener:     listener,
	}
	f.fundCrawler.SetCallBack(f.fundCallback)
	go http.Serve(listener, f.server)
	return f
}

func (f *fixture) GetCompleteURL(resource string) *url.URL {
	return &url.URL{
		Scheme: "http",
		Host:   f.listener.Addr().String(),
		Path:   resource,
	}
}
