package controllers

import (
	"html"
	"log"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/luxcgo/go-gallery/models"
	"gopkg.in/headzoo/surf.v1"
)

func (g *Galleries) UpdateSd() {
	lastGallery, err := g.gs.GetLast(2)
	if err != nil {
		return
	}

	u := "http://www.cwl.gov.cn/cwl_admin/front/cwlkj/search/kjxx/findDrawNotice?name=3d&issueCount=30"

	bow := surf.NewBrowser()
	if err := bow.Open(u); err != nil {
		log.Println(err.Error())
		return
	}

	resp := html.UnescapeString((bow.Body()))
	var auto TdAutoGenerated
	if err := jsoniter.UnmarshalFromString(resp, &auto); err != nil {
		log.Println(err.Error())
		return
	}
	for _, v := range auto.Result {
		_3d := new(models.Gallery)
		_3d.No, _ = strconv.Atoi(v.Code)
		if _3d.No <= lastGallery.No {
			return
		}
		noArr := strings.Split(v.Red, ",")
		if len(noArr) != 3 {
			return
		}
		_3d.Num1, _ = strconv.Atoi(noArr[0])
		_3d.Num2, _ = strconv.Atoi(noArr[1])
		_3d.Num3, _ = strconv.Atoi(noArr[2])
		_3d.Date = v.Date
		_3d.LotteryType = 2
		g.gs.Create(_3d)
	}
}

type TdAutoGenerated struct {
	State     int    `json:"state"`
	Message   string `json:"message"`
	PageCount int    `json:"pageCount"`
	CountNum  int    `json:"countNum"`
	Tflag     int    `json:"Tflag"`
	Result    []struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		DetailsLink string `json:"detailsLink"`
		VideoLink   string `json:"videoLink"`
		Date        string `json:"date"`
		Week        string `json:"week"`
		Red         string `json:"red"`
		Blue        string `json:"blue"`
		Blue2       string `json:"blue2"`
		Sales       string `json:"sales"`
		Poolmoney   string `json:"poolmoney"`
		Content     string `json:"content"`
		Addmoney    string `json:"addmoney"`
		Addmoney2   string `json:"addmoney2"`
		Msg         string `json:"msg"`
		Z2Add       string `json:"z2add"`
		M2Add       string `json:"m2add"`
		Prizegrades []struct {
			Type      int    `json:"type"`
			Typenum   string `json:"typenum"`
			Typemoney string `json:"typemoney"`
		} `json:"prizegrades"`
	} `json:"result"`
}
