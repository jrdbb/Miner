package crawler

import (
	"errors"
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/robertkrimen/otto"
	log "github.com/sirupsen/logrus"
)

func toBasicFund(elements []*otto.Value) (*BasicFund, error) {
	if len(elements) != 5 {
		return nil, errors.New("toBasicFund Error: len(elements) != 5")
	}
	for i, e := range elements {
		if !e.IsString() {
			return nil, errors.New(fmt.Sprintf("toBasicFund Error: element(%d) is not string", i))
		}
	}
	bf := &BasicFund{
		ID:        elements[0].String(),
		ShortCode: elements[1].String(),
		Type:      elements[2].String(),
		NameCN:    elements[3].String(),
		NamePY:    elements[4].String(),
	}
	return bf, nil
}

func (cl *crawlerImpl) GetAllBasicFund(sync bool) {
	log.Debugf("Start collecting all basic fund")
	c, ok := cl.mCollectors[basicFundCollector]
	if !ok {
		vm := otto.New()
		c = newDefaultCollector()
		c.OnResponse(func(r *colly.Response) {
			log.Infof("GetAllBasicFund OnResponse")
			_, err := vm.Run(r.Body)
			if err != nil {
				log.Errorf("GetAllBasicFund Err(%v)", err)
				return
			}
			ret, err := vm.Get("r")
			if err != nil {
				log.Errorf("GetAllBasicFund Err(%v)", err)
				return
			}
			jsFundArray, err := ParseArray(&ret)
			if err != nil {
				log.Errorf("GetAllBasicFund Err(%v)", err)
				return
			}

			funds := make([]*BasicFund, 0, len(jsFundArray))
			for _, v := range jsFundArray {
				fundElements, err := ParseArray(v)
				if err != nil {
					log.Errorf("GetAllBasicFund Err(%v)", err)
					return
				}
				basicFund, err := toBasicFund(fundElements)
				if err != nil {
					log.Errorf("GetAllBasicFund Err(%v)", err)
					return
				}
				funds = append(funds, basicFund)
			}
			cl.mCallBack.OnBasicFund(funds)
		})

		cl.mCollectors[basicFundCollector] = c
	}

	c.Visit(URLCenter[FundCodeSearch])
	if sync {
		c.Wait()
	}
}
