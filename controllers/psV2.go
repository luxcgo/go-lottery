package controllers

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/gogf/gf/net/ghttp"
	"github.com/luxcgo/go-gallery/models"
)

func (g *Galleries) UpdatePsV2() {
	lastGallery, err := g.gs.GetLast(1)
	if err != nil {
		log.Fatal(err)
		return
	}

	client := ghttp.NewClient()
	client.SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.120 Safari/537.36")
	// client.SetTimeOut(time.Second * 10)

	for i := 1; ; i++ {
		url := fmt.Sprintf("https://www.lottery.gov.cn/historykj/history_%d.jspx?_ltype=pls", i)
		response, err := client.Get(url)
		if err != nil {
			fmt.Println("error")
			return
		}
		defer response.Body.Close()
		ct := response.ReadAllString()
		ct = strings.ReplaceAll(ct, " ", "")
		ct = strings.ReplaceAll(ct, "\n", "")
		ct = strings.ReplaceAll(ct, "\b", "")
		ct = strings.ReplaceAll(ct, "ececec", "f9f9f9")

		end := g.procV2(ct, lastGallery.No)
		if end {
			break
		}
	}

}

var (
	psFirstNo = 04001
	tr        = `<tr><tdheight="23"align="center"bgcolor="#f9f9f9">([\d]+)</td><tdalign="center"bgcolor="#f9f9f9"class="cpl">([\d]+)</td>`
	trback    = `<tdalign="center"bgcolor="#f9f9f9">([^<]+)</td></tr>`
	retr      = regexp.MustCompile(tr)
	retrback  = regexp.MustCompile(trback)
	content   = `<tr><tdheight="23"align="center"bgcolor="#f9f9f9">20145</td><tdalign="center"bgcolor="#f9f9f9"class="cpl">846</td><tdalign="center"bgcolor="#f9f9f9">4,886</td><tdalign="center"bgcolor="#f9f9f9">1,040</td><tdalign="center"bgcolor="#f9f9f9">0</td><tdalign="center"bgcolor="#f9f9f9">346</td><tdalign="center"bgcolor="#f9f9f9">17,960</td><tdalign="center"bgcolor="#f9f9f9">173</td><tdalign="center"bgcolor="#f9f9f9"><!--补位--><aid="pl312"href="javascript:"flag="pl312,pl3img12,28200,20145,/kjpls/,0"><imgid="pl3img12"src="/res/images/history/fangda.png"width="20"height="20"/></a></td><tdalign="center"bgcolor="#f9f9f9">22,776,768</td><tdalign="center"bgcolor="#f9f9f9">2020-07-13</td></tr>`
)

func (g *Galleries) procV2(c string, lastNo int) bool {
	submatch := retr.FindAllStringSubmatch(c, -1)
	var uid = []int{}
	var winNo = []int{}
	var date = []string{}
	for _, v := range submatch {
		if v1, err := strconv.Atoi(v[1]); err != nil {
			log.Println(err.Error())
			return true
		} else {
			uid = append(uid, v1)
		}

		if v2, err := strconv.Atoi(v[2]); err != nil {
			log.Println(err.Error())
			return true
		} else {
			winNo = append(winNo, v2)
		}

	}

	submatch = retrback.FindAllStringSubmatch(c, -1)
	for _, v := range submatch {
		date = append(date, v[1])
	}

	// log.Println(uid)
	// log.Println(winNo)
	// log.Println(date)

	if len(uid) != len(winNo) || len(winNo) != len(date) {
		return true
	}

	for i := 0; i < len(uid); i++ {
		_3d := new(models.Gallery)
		_3d.No = uid[i]
		if _3d.No <= lastNo {
			return true
		}
		_3d.Num1 = winNo[i] / 100 % 10
		_3d.Num2 = winNo[i] / 10 % 10
		_3d.Num3 = winNo[i] % 10
		_3d.Date = date[i]
		_3d.LotteryType = 1
		// fmt.Printf("%+v", _3d)
		g.gs.Create(_3d)

	}

	for _, v := range uid {
		if v == psFirstNo {
			return true
		}
	}
	return false
}
