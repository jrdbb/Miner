package crawler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/robertkrimen/otto"
	log "github.com/sirupsen/logrus"
)

type ApiData struct {
	Code     string
	Content  []FundValue
	Records  int64
	Pages    int64
	CurrPage int64
}

func toApiData(jsValue *otto.Value) (*ApiData, error) {
	if !jsValue.IsObject() {
		return nil, errors.New("toApiData: input is not an object")
	}
	jsContent, err := jsValue.Object().Get("content")
	if err != nil || !jsContent.IsString() {
		return nil, errors.New("toApiData: decode content fail")
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(jsContent.String()))
	if err != nil {
		return nil, errors.New("toApiData: NewDocumentFromReader fail")
	}
	content := make([]FundValue, 0)
	err = nil
	doc.Find("tbody>tr").Each(func(i int, s *goquery.Selection) {
		// 净值日期 单位净值 累计净值 日增长率 申购状态 赎回状态 分红送配
		tds := s.Find("td").Contents().Nodes
		if len(tds) != 6 && len(tds) != 7 {
			err = fmt.Errorf("toApiData: tbody len(nodes)(%d) != 6 or 7", len(tds))
			return
		}
		dt, err := time.Parse("2006-01-02", tds[0].Data)
		if err != nil {
			err = fmt.Errorf("toApiData: convert date(%s) to date fail", tds[0].Data)
			return
		}

		netAssetValue, err := strconv.ParseFloat(tds[1].Data, 64)
		if err != nil {
			err = fmt.Errorf("toApiData: convert netAssetValue(%s) to float fail", tds[1].Data)
			return
		}
		log.Debugf("toApiData: process %v", tds)
		content = append(content, FundValue{
			Date:          dt,
			NetAssetValue: netAssetValue,
			BuyState:      tds[4].Data,
			SellState:     tds[5].Data,
		})
	})
	if err != nil {
		return nil, err
	}

	// not check the error, as the default value is 0
	jsRecords, _ := jsValue.Object().Get("records")
	jsPages, _ := jsValue.Object().Get("pages")
	jsCurrPage, _ := jsValue.Object().Get("curpage")
	records, _ := jsRecords.ToInteger()
	pages, _ := jsPages.ToInteger()
	currPage, _ := jsCurrPage.ToInteger()
	return &ApiData{
		Content:  content,
		Records:  records,
		Pages:    pages,
		CurrPage: currPage,
	}, nil
}

func (cl *crawlerImpl) GetHistoryValue(sync bool, code string, page int, sdate string, edate string) {
	log.Debugf("Start collecting history value")
	c, ok := cl.mCollectors[historyValueCollector]
	if !ok {
		c = newDefaultCollector()
		c.OnResponse(func(r *colly.Response) {
			vm := otto.New()
			_, err := vm.Run(r.Body)
			if err != nil {
				log.Errorf("GetHistoryValue Err(%v)", err)
				return
			}
			ret, err := vm.Get("apidata")
			if err != nil {
				log.Errorf("GetHistoryValue Err(%v)", err)
				return
			}
			data, err := toApiData(&ret)
			if err != nil {
				log.Errorf("GetHistoryValue Err(%v)", err)
				return
			}
			data.Code = r.Request.URL.Query().Get("code")
			cl.mCallBack.OnHistoryValue(data)
		})

		cl.mCollectors[basicFundCollector] = c
	}
	target := *URLCenter[HistoryValue]
	qry := target.Query()
	qry.Add("type", "lsjz")
	qry.Add("code", code)
	qry.Add("page", strconv.Itoa(page))
	qry.Add("sdate", sdate)
	qry.Add("edate", edate)
	qry.Add("per", "49")
	target.RawQuery = qry.Encode()

	c.Visit(target.String())
	if sync {
		c.Wait()
	}
}
