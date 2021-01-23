package tests

import (
	"net"
	"net/http"

	"github.com/CommonProsperity/Miner/crawler"
	"github.com/CommonProsperity/Miner/crawler/mocks"
	"github.com/golang/mock/gomock"
)

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

func (f *fixture) GetCompleteURL(resource string) string {
	return "http://" + f.listener.Addr().String() + resource
}
