package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/text/gregex"
	"github.com/gorilla/mux"
	"github.com/luxcgo/go-gallery/models"
	"github.com/luxcgo/go-gallery/views"
)

const (
	maxMultipartMem = 1 << 20 // 1 megabyte

	IndexGalleries  = "index_galleries"
	Graph2Galleries = "graph2_galleries"
	Graph3Galleries = "graph3_galleries"
	Graph5Galleries = "graph5_galleries"
	ShowGallery     = "show_gallery"
	EditGallery     = "edit_gallery"
)

var graphTypeMap = map[int]string{
	1: "排三原始数据",
	2: "排三组三和组六",
	3: "排三组三和组六压缩1",
	4: "排三组三和组六压缩2",
	5: "排三组三和组六压缩2延展",
	6: "排三组六",
	7: "排三组六压缩1",
	8: "排三组六压缩2",
}

type Galleries struct {
	New        *views.View
	ShowView   *views.View
	EditView   *views.View
	IndexView  *views.View
	Graph2View *views.View
	Graph3View *views.View
	Graph5View *views.View
	Graph8View *views.View
	Graph9View *views.View
	gs         models.GalleryService
	r          *mux.Router
}

type GalleryForm struct {
	Title string `schema:"title"`
}

type QueryForm struct {
	PermutationKey int `schema:"permutationKey"`
	LotteryType    int `schema:"lotteryType"`
	GraphType      int `schema:"graphType"`
	FourNumber     int `schema:"fourNumber"`
	Offset         int `schema:"offset"`
}

func NewGalleries(gs models.GalleryService,
	r *mux.Router) *Galleries {
	return &Galleries{
		IndexView:  views.NewView("bootstrap", "galleries/index"),
		Graph2View: views.NewView("bootstrap", "galleries/graph2"),
		Graph3View: views.NewView("bootstrap", "galleries/graph3"),
		Graph5View: views.NewView("bootstrap", "galleries/graph5"),
		Graph8View: views.NewView("bootstrap", "galleries/graph8"),
		Graph9View: views.NewView("bootstrap", "galleries/graph9"),
		gs:         gs,
		r:          r,
	}
}

func (g *Galleries) Run() {
	g.UpdatePs()
	g.UpdateSd()
}

func (g *Galleries) UpdatePsV1() {
	lastGallery, err := g.gs.GetLast(1)
	if err != nil {
		return
	}

	client := ghttp.NewClient()
	client.SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.120 Safari/537.36")
	// client.SetTimeOut(time.Second * 10)

	response, err := client.Get("http://datachart.500.com/pls/history/inc/history.php?limit=20000")
	if err != nil {
		fmt.Println("error")
	}
	ct := response.ReadAllString()
	match, err := gregex.MatchAllString(`<tr class="t_tr1">\s*?<!--<td>2</td>-->\s*?<td class="t_tr1">(\d{5})</td>\s*?<td class="cfont2">(\d{1})\s*?(\d{1})\s*?(\d{1})</td>.*?<td class="t_tr1">([0-9,\-]*?)<\/td>\s*?</tr>`, ct)
	response.Close()
	// fmt.Println(match)
	for _, arr := range match {
		_3d := new(models.Gallery)
		_3d.No, _ = strconv.Atoi(arr[1])
		if _3d.No <= lastGallery.No {
			return
		}
		_3d.Num1, _ = strconv.Atoi(arr[2])
		_3d.Num2, _ = strconv.Atoi(arr[3])
		_3d.Num3, _ = strconv.Atoi(arr[4])
		_3d.Date = arr[5]
		_3d.LotteryType = 1
		// fmt.Println("%+v", _3d)
		g.gs.Create(_3d)
	}
}

func (g *Galleries) UpdateSdV1() {
	lastGallery, err := g.gs.GetLast(2)
	if err != nil {
		return
	}

	client := ghttp.NewClient()
	client.SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.120 Safari/537.36")
	// client.SetTimeOut(time.Second * 10)

	response, err := client.Get("http://datachart.500.com/sd/history/inc/history.php?limit=20000")
	if err != nil {
		fmt.Println("error")
	}
	ct := response.ReadAllString()
	match, err := gregex.MatchAllString(`<tr class="t_tr1">\s*?<!--<td>2</td>-->\s*?<td>(\d{7})</td>\s*?<td class="cfont2">(\d{1})\s*?(\d{1})\s*?(\d{1})</td>.*?<td class="t_tr1">([0-9,\-]*?)<\/td>\s*?</tr>`, ct)
	response.Close()
	// fmt.Println(match)
	for _, arr := range match {
		_3d := new(models.Gallery)
		_3d.No, _ = strconv.Atoi(arr[1])
		if _3d.No <= lastGallery.No {
			return
		}
		_3d.Num1, _ = strconv.Atoi(arr[2])
		_3d.Num2, _ = strconv.Atoi(arr[3])
		_3d.Num3, _ = strconv.Atoi(arr[4])
		_3d.Date = arr[5]
		_3d.LotteryType = 2
		g.gs.Create(_3d)
	}

}

// POST /galleries
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data

	var form QueryForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.IndexView.Render(w, r, vd)
		return
	}
	vd.Form = form

	// fmt.Printf("%+v", form)
	// fmt.Println(form.PermutationKey)
	switch form.GraphType {
	case 1:
		galleries, err := g.gs.GetAllPs(form.LotteryType)
		if err != nil {
			http.Error(w, "Something went wrong.", http.StatusInternalServerError)
			return
		}
		var vd views.Data
		vd.Yield = galleries
		g.IndexView.Render(w, r, vd)
	case 2:
		res, _ := g.gs.GetGraph2(form.LotteryType, form.PermutationKey)
		var vd views.Data
		vd.Yield = res
		g.Graph2View.Render(w, r, vd)
	case 3:
		res, _ := g.gs.GetGraph3(form.LotteryType, form.PermutationKey)
		var vd views.Data
		vd.Yield = res.Arr
		g.Graph3View.Render(w, r, vd)
	case 4:
		res, _ := g.gs.GetGraph4(form.LotteryType, form.PermutationKey)
		var vd views.Data
		vd.Yield = res.Arr
		g.Graph3View.Render(w, r, vd)
	case 5:
		res, _ := g.gs.GetGraph5(form.LotteryType, form.PermutationKey)
		var vd views.Data
		vd.Yield = res
		g.Graph5View.Render(w, r, vd)
	case 6:
		res, _ := g.gs.GetGraph6(form.LotteryType, form.PermutationKey)
		var vd views.Data
		vd.Yield = res.Arr
		g.Graph3View.Render(w, r, vd)
	case 7:
		res, _ := g.gs.GetGraph7(form.LotteryType, form.PermutationKey)
		var vd views.Data
		vd.Yield = res.Arr
		g.Graph3View.Render(w, r, vd)
	case 8:
		res, _ := g.gs.GetGraph8(form.LotteryType, form.FourNumber)
		var vd views.Data
		vd.Yield = res.Arr
		g.Graph8View.Render(w, r, vd)
	case 9:
		res, _ := g.gs.GetGraph9(form.LotteryType, form.FourNumber, form.Offset)
		var vd views.Data
		vd.Yield = res.Arr
		g.Graph9View.Render(w, r, vd)
	}

	g.IndexView.Render(w, r, vd)
}

// GET /galleries
func (g *Galleries) Index(w http.ResponseWriter, r *http.Request) {
	galleries, err := g.gs.GetAllPs(1)
	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	var vd views.Data
	vd.Yield = galleries

	var form QueryForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.IndexView.Render(w, r, vd)
		return
	}
	fmt.Printf("%+v", form)
	fmt.Println(form.PermutationKey)
	g.IndexView.Render(w, r, vd)
}
