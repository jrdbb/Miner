package tests

import (
	"net/http"
	"testing"

	"github.com/CommonProsperity/Miner/crawler"
	"github.com/golang/mock/gomock"
	"github.com/smartystreets/assertions"
	. "github.com/smartystreets/goconvey/convey"
	// log "github.com/sirupsen/logrus"
)

func TestBaicFund(t *testing.T) {
	Convey("TestBaicFund", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		f := newFixture(ctrl)
		crawler.URLCenter[crawler.FundCodeSearch] = f.GetCompleteURL("/fundcode")
		f.server.HandleFunc("/fundcode", func(rw http.ResponseWriter, r *http.Request) {
			rw.Write([]byte(`var r = [["000001","HXCZHH","华夏成长混合","混合型","HUAXIACHENGZHANGHUNHE"],["000002","HXCZHH","华夏成长混合(后端)","混合型","HUAXIACHENGZHANGHUNHE"]]`))
		})

		var funds []*crawler.BasicFund = nil
		f.fundCallback.EXPECT().OnBasicFund(gomock.Any()).Do(
			func(arg []*crawler.BasicFund) {
				funds = arg
			},
		)

		f.fundCrawler.GetAllBasicFund(true)

		So(funds, assertions.ShouldNotBeNil)
		So(len(funds), assertions.ShouldEqual, 2)
		So(*funds[0], ShouldResemble, crawler.BasicFund{"000001", "HXCZHH", "华夏成长混合", "混合型", "HUAXIACHENGZHANGHUNHE"})
		So(*funds[1], ShouldResemble, crawler.BasicFund{"000002","HXCZHH","华夏成长混合(后端)","混合型","HUAXIACHENGZHANGHUNHE"})
	})
}
