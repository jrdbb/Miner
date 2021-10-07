package tests

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/jrdbb/Miner/crawler"
	"github.com/golang/mock/gomock"
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

		f.fundCallback.EXPECT().OnBasicFund(gomock.Any()).Do(
			func(funds []*crawler.BasicFund) {
				Convey("Check args", t, func() {
					So(funds, ShouldNotBeNil)
					So(len(funds), ShouldEqual, 2)
					So(*funds[0], ShouldResemble, crawler.BasicFund{"000001", "HXCZHH", "华夏成长混合", "混合型", "HUAXIACHENGZHANGHUNHE"})
					So(*funds[1], ShouldResemble, crawler.BasicFund{"000002", "HXCZHH", "华夏成长混合(后端)", "混合型", "HUAXIACHENGZHANGHUNHE"})
				})
			},
		)

		f.fundCrawler.GetAllBasicFund(true)
	})
}

func TestHistoryValue(t *testing.T) {
	Convey("TestHistoryValue", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		f := newFixture(ctrl)
		crawler.URLCenter[crawler.HistoryValue] = f.GetCompleteURL("/historyValue")

		code := "020111"
		page := 1
		sdate := "2020-01-15"
		edate := "2020-02-02"

		f.server.HandleFunc("/historyValue", func(rw http.ResponseWriter, r *http.Request) {
			Convey("check request args", t, func() {
				So(r.URL.Query().Get("type"), ShouldEqual, "lsjz")
				So(r.URL.Query().Get("code"), ShouldEqual, code)
				So(r.URL.Query().Get("page"), ShouldEqual, strconv.Itoa(page))
				So(r.URL.Query().Get("sdate"), ShouldEqual, sdate)
				So(r.URL.Query().Get("edate"), ShouldEqual, edate)
			})
			rw.Write([]byte(`var apidata={ content:"<table class='w782 comm lsjz'><thead><tr><th class='first'>净值日期</th><th>单位净值</th><th>累计净值</th><th>日增长率</th><th>申购状态</th><th>赎回状态</th><th class='tor last'>分红送配</th></tr></thead><tbody><tr><td>2021-01-25</td><td class='tor bold'>1.4530</td><td class='tor bold'>3.9640</td><td class='tor bold red'>2.54%</td><td>开放申购</td><td>开放赎回</td><td class='red unbold'></td></tr><tr><td>2021-01-22</td><td class='tor bold'>1.4170</td><td class='tor bold'>3.9280</td><td class='tor bold red'>2.61%</td><td>开放申购</td><td>开放赎回</td><td class='red unbold'></td></tr><tr><td>2021-01-21</td><td class='tor bold'>1.3810</td><td class='tor bold'>3.8920</td><td class='tor bold red'>1.32%</td><td>开放申购</td><td>开放赎回</td><td class='red unbold'></td></tr><tr><td>2021-01-20</td><td class='tor bold'>1.3630</td><td class='tor bold'>3.8740</td><td class='tor bold red'>2.64%</td><td>开放申购</td><td>开放赎回</td><td class='red unbold'></td></tr><tr><td>2021-01-19</td><td class='tor bold'>1.3280</td><td class='tor bold'>3.8390</td><td class='tor bold grn'>-1.26%</td><td>开放申购</td><td>开放赎回</td><td class='red unbold'></td></tr><tr><td>2021-01-18</td><td class='tor bold'>1.3450</td><td class='tor bold'>3.8560</td><td class='tor bold grn'>-0.37%</td><td>开放申购</td><td>开放赎回</td><td class='red unbold'></td></tr><tr><td>2021-01-15</td><td class='tor bold'>1.3500</td><td class='tor bold'>3.8610</td><td class='tor bold grn'>-0.44%</td><td>开放申购</td><td>开放赎回</td><td class='red unbold'></td></tr><tr><td>2021-01-14</td><td class='tor bold'>1.3560</td><td class='tor bold'>3.8670</td><td class='tor bold grn'>-2.16%</td><td>开放申购</td><td>开放赎回</td><td class='red unbold'></td></tr><tr><td>2021-01-13</td><td class='tor bold'>1.3860</td><td class='tor bold'>3.8970</td><td class='tor bold grn'>-0.29%</td><td>开放申购</td><td>开放赎回</td><td class='red unbold'></td></tr><tr><td>2021-01-12</td><td class='tor bold'>1.3900</td><td class='tor bold'>3.9010</td><td class='tor bold red'>2.06%</td><td>开放申购</td><td>开放赎回</td><td class='red unbold'></td></tr></tbody></table>",records:4637,pages:464,curpage:1};`))
		})

		f.fundCallback.EXPECT().OnHistoryValue(gomock.Any()).Do(
			func(arg *crawler.ApiData) {
				Convey("check history value", t, func() {
					So(arg, ShouldNotBeNil)
					So(arg.Code, ShouldEqual, code)
					So(len(arg.Content), ShouldEqual, 10)
					So(arg.CurrPage, ShouldEqual, 1)
					So(arg.Records, ShouldEqual, 4637)
					So(arg.Pages, ShouldEqual, 464)
				})
			},
		)

		f.fundCrawler.GetHistoryValue(true, code, page, sdate, edate)
	})
}
