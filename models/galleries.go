package models

import (
	"log"
	"sort"
	"strconv"

	"gorm.io/gorm"
)

const (

	// ErrNotFound is returned when a resource cannot be found
	// in the database.
	ErrNotFound modelError = "models: resource not found"
	RowCount               = 10
)

var (
	permutation = map[int][]int{
		1:   {0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		2:   {0, 1, 2, 3, 5, 4, 6, 7, 8, 9},
		3:   {0, 1, 2, 3, 6, 4, 5, 7, 8, 9},
		4:   {0, 1, 2, 3, 7, 4, 5, 6, 8, 9},
		5:   {0, 1, 2, 3, 8, 4, 5, 6, 7, 9},
		6:   {0, 1, 2, 3, 9, 4, 5, 6, 7, 8},
		7:   {0, 1, 2, 4, 5, 3, 6, 7, 8, 9},
		8:   {0, 1, 2, 4, 6, 3, 5, 7, 8, 9},
		9:   {0, 1, 2, 4, 7, 3, 5, 6, 8, 9},
		10:  {0, 1, 2, 4, 8, 3, 5, 6, 7, 9},
		11:  {0, 1, 2, 4, 9, 3, 5, 6, 7, 8},
		12:  {0, 1, 2, 5, 6, 3, 4, 7, 8, 9},
		13:  {0, 1, 2, 5, 7, 3, 4, 6, 8, 9},
		14:  {0, 1, 2, 5, 8, 3, 4, 6, 7, 9},
		15:  {0, 1, 2, 5, 9, 3, 4, 6, 7, 8},
		16:  {0, 1, 2, 6, 7, 3, 4, 5, 8, 9},
		17:  {0, 1, 2, 6, 8, 3, 4, 5, 7, 9},
		18:  {0, 1, 2, 6, 9, 3, 4, 5, 7, 8},
		19:  {0, 1, 2, 7, 8, 3, 4, 5, 6, 9},
		20:  {0, 1, 2, 7, 9, 3, 4, 5, 6, 8},
		21:  {0, 1, 2, 8, 9, 3, 4, 5, 6, 7},
		22:  {0, 1, 3, 4, 5, 2, 6, 7, 8, 9},
		23:  {0, 1, 3, 4, 6, 2, 5, 7, 8, 9},
		24:  {0, 1, 3, 4, 7, 2, 5, 6, 8, 9},
		25:  {0, 1, 3, 4, 8, 2, 5, 6, 7, 9},
		26:  {0, 1, 3, 4, 9, 2, 5, 6, 7, 8},
		27:  {0, 1, 3, 5, 6, 2, 4, 7, 8, 9},
		28:  {0, 1, 3, 5, 7, 2, 4, 6, 8, 9},
		29:  {0, 1, 3, 5, 8, 2, 4, 6, 7, 9},
		30:  {0, 1, 3, 5, 9, 2, 4, 6, 7, 8},
		31:  {0, 1, 3, 6, 7, 2, 4, 5, 8, 9},
		32:  {0, 1, 3, 6, 8, 2, 4, 5, 7, 9},
		33:  {0, 1, 3, 6, 9, 2, 4, 5, 7, 8},
		34:  {0, 1, 3, 7, 8, 2, 4, 5, 6, 9},
		35:  {0, 1, 3, 7, 9, 2, 4, 5, 6, 8},
		36:  {0, 1, 3, 8, 9, 2, 4, 5, 6, 7},
		37:  {0, 1, 4, 5, 6, 2, 3, 7, 8, 9},
		38:  {0, 1, 4, 5, 7, 2, 3, 6, 8, 9},
		39:  {0, 1, 4, 5, 8, 2, 3, 6, 7, 9},
		40:  {0, 1, 4, 5, 9, 2, 3, 6, 7, 8},
		41:  {0, 1, 4, 6, 7, 2, 3, 5, 8, 9},
		42:  {0, 1, 4, 6, 8, 2, 3, 5, 7, 9},
		43:  {0, 1, 4, 6, 9, 2, 3, 5, 7, 8},
		44:  {0, 1, 4, 7, 8, 2, 3, 5, 6, 9},
		45:  {0, 1, 4, 7, 9, 2, 3, 5, 6, 8},
		46:  {0, 1, 4, 8, 9, 2, 3, 5, 6, 7},
		47:  {0, 1, 5, 6, 7, 2, 3, 4, 8, 9},
		48:  {0, 1, 5, 6, 8, 2, 3, 4, 7, 9},
		49:  {0, 1, 5, 6, 9, 2, 3, 4, 7, 8},
		50:  {0, 1, 5, 7, 8, 2, 3, 4, 6, 9},
		51:  {0, 1, 5, 7, 9, 2, 3, 4, 6, 8},
		52:  {0, 1, 5, 8, 9, 2, 3, 4, 6, 7},
		53:  {0, 1, 6, 7, 8, 2, 3, 4, 5, 9},
		54:  {0, 1, 6, 7, 9, 2, 3, 4, 5, 8},
		55:  {0, 1, 6, 8, 9, 2, 3, 4, 5, 7},
		56:  {0, 1, 7, 8, 9, 2, 3, 4, 5, 6},
		57:  {0, 2, 3, 4, 5, 1, 6, 7, 8, 9},
		58:  {0, 2, 3, 4, 6, 1, 5, 7, 8, 9},
		59:  {0, 2, 3, 4, 7, 1, 5, 6, 8, 9},
		60:  {0, 2, 3, 4, 8, 1, 5, 6, 7, 9},
		61:  {0, 2, 3, 4, 9, 1, 5, 6, 7, 8},
		62:  {0, 2, 3, 5, 6, 1, 4, 7, 8, 9},
		63:  {0, 2, 3, 5, 7, 1, 4, 6, 8, 9},
		64:  {0, 2, 3, 5, 8, 1, 4, 6, 7, 9},
		65:  {0, 2, 3, 5, 9, 1, 4, 6, 7, 8},
		66:  {0, 2, 3, 6, 7, 1, 4, 5, 8, 9},
		67:  {0, 2, 3, 6, 8, 1, 4, 5, 7, 9},
		68:  {0, 2, 3, 6, 9, 1, 4, 5, 7, 8},
		69:  {0, 2, 3, 7, 8, 1, 4, 5, 6, 9},
		70:  {0, 2, 3, 7, 9, 1, 4, 5, 6, 8},
		71:  {0, 2, 3, 8, 9, 1, 4, 5, 6, 7},
		72:  {0, 2, 4, 5, 6, 1, 3, 7, 8, 9},
		73:  {0, 2, 4, 5, 7, 1, 3, 6, 8, 9},
		74:  {0, 2, 4, 5, 8, 1, 3, 6, 7, 9},
		75:  {0, 2, 4, 5, 9, 1, 3, 6, 7, 8},
		76:  {0, 2, 4, 6, 7, 1, 3, 5, 8, 9},
		77:  {0, 2, 4, 6, 8, 1, 3, 5, 7, 9},
		78:  {0, 2, 4, 6, 9, 1, 3, 5, 7, 8},
		79:  {0, 2, 4, 7, 8, 1, 3, 5, 6, 9},
		80:  {0, 2, 4, 7, 9, 1, 3, 5, 6, 8},
		81:  {0, 2, 4, 8, 9, 1, 3, 5, 6, 7},
		82:  {0, 2, 5, 6, 7, 1, 3, 4, 8, 9},
		83:  {0, 2, 5, 6, 8, 1, 3, 4, 7, 9},
		84:  {0, 2, 5, 6, 9, 1, 3, 4, 7, 8},
		85:  {0, 2, 5, 7, 8, 1, 3, 4, 6, 9},
		86:  {0, 2, 5, 7, 9, 1, 3, 4, 6, 8},
		87:  {0, 2, 5, 8, 9, 1, 3, 4, 6, 7},
		88:  {0, 2, 6, 7, 8, 1, 3, 4, 5, 9},
		89:  {0, 2, 6, 7, 9, 1, 3, 4, 5, 8},
		90:  {0, 2, 6, 8, 9, 1, 3, 4, 5, 7},
		91:  {0, 2, 7, 8, 9, 1, 3, 4, 5, 6},
		92:  {0, 3, 4, 5, 6, 1, 2, 7, 8, 9},
		93:  {0, 3, 4, 5, 7, 1, 2, 6, 8, 9},
		94:  {0, 3, 4, 5, 8, 1, 2, 6, 7, 9},
		95:  {0, 3, 4, 5, 9, 1, 2, 6, 7, 8},
		96:  {0, 3, 4, 6, 7, 1, 2, 5, 8, 9},
		97:  {0, 3, 4, 6, 8, 1, 2, 5, 7, 9},
		98:  {0, 3, 4, 6, 9, 1, 2, 5, 7, 8},
		99:  {0, 3, 4, 7, 8, 1, 2, 5, 6, 9},
		100: {0, 3, 4, 7, 9, 1, 2, 5, 6, 8},
		101: {0, 3, 4, 8, 9, 1, 2, 5, 6, 7},
		102: {0, 3, 5, 6, 7, 1, 2, 4, 8, 9},
		103: {0, 3, 5, 6, 8, 1, 2, 4, 7, 9},
		104: {0, 3, 5, 6, 9, 1, 2, 4, 7, 8},
		105: {0, 3, 5, 7, 8, 1, 2, 4, 6, 9},
		106: {0, 3, 5, 7, 9, 1, 2, 4, 6, 8},
		107: {0, 3, 5, 8, 9, 1, 2, 4, 6, 7},
		108: {0, 3, 6, 7, 8, 1, 2, 4, 5, 9},
		109: {0, 3, 6, 7, 9, 1, 2, 4, 5, 8},
		110: {0, 3, 6, 8, 9, 1, 2, 4, 5, 7},
		111: {0, 3, 7, 8, 9, 1, 2, 4, 5, 6},
		112: {0, 4, 5, 6, 7, 1, 2, 3, 8, 9},
		113: {0, 4, 5, 6, 8, 1, 2, 3, 7, 9},
		114: {0, 4, 5, 6, 9, 1, 2, 3, 7, 8},
		115: {0, 4, 5, 7, 8, 1, 2, 3, 6, 9},
		116: {0, 4, 5, 7, 9, 1, 2, 3, 6, 8},
		117: {0, 4, 5, 8, 9, 1, 2, 3, 6, 7},
		118: {0, 4, 6, 7, 8, 1, 2, 3, 5, 9},
		119: {0, 4, 6, 7, 9, 1, 2, 3, 5, 8},
		120: {0, 4, 6, 8, 9, 1, 2, 3, 5, 7},
		121: {0, 4, 7, 8, 9, 1, 2, 3, 5, 6},
		122: {0, 5, 6, 7, 8, 1, 2, 3, 4, 9},
		123: {0, 5, 6, 7, 9, 1, 2, 3, 4, 8},
		124: {0, 5, 6, 8, 9, 1, 2, 3, 4, 7},
		125: {0, 5, 7, 8, 9, 1, 2, 3, 4, 6},
		126: {0, 6, 7, 8, 9, 1, 2, 3, 4, 5},
	}
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

var _ GalleryDB = &galleryGorm{}

type Gallery struct {
	gorm.Model
	Date        string
	No          int `gorm:"uniqueIndex"`
	Num1        int
	Num2        int
	Num3        int
	LotteryType int
}

type Ps struct {
	Date  string
	No    int
	Index int
	Num1  int
	Num2  int
	Num3  int
}

type Graph2 struct {
	Date     string
	No       int
	Index    int
	Num1     int
	Num2     int
	Num3     int
	Omission int
	Diff1    int
	Diff2    int
	Diff3    int
	Diff4    int
}

type Graph3 struct {
	Arr [][]int
}

type Graph8 struct {
	Arr [][]string
}

type Graph5 struct {
	Date  string
	No    int
	Index int
	Num1  int
	Num2  int
	Num3  int
	Arr0  int
	Arr1  int
	Arr2  int
	Arr3  int
	Arr4  int
	Arr5  int
	Arr6  int
	Arr7  int
	Arr8  int
	Arr9  int
}

type GalleryService interface {
	GalleryDB
}

// GalleryDB is used to interact with the galleries database.
//
// For pretty much all single gallery queries:
// If the gallery is found, we will return a nil error
// If the gallery is not found, we will return ErrNotFound
// If there is another error, we will return an error with
// more information about what went wrong. This may not be
// an error generated by the models package.
type GalleryDB interface {
	ByID(id uint) (*Gallery, error)
	GetLast(int) (*Gallery, error)
	Create(gallery *Gallery) error
	GetAllPs(int) ([]Ps, error)
	GetGraph2(lotteryType int, permutationKey int) ([]Graph2, error)
	GetGraph3(lotteryType int, permutationKey int) (Graph3, error)
	GetGraph4(lotteryType int, permutationKey int) (Graph3, error)
	GetGraph5(lotteryType int, permutationKey int) ([]Graph5, error)
	GetGraph6(lotteryType int, permutationKey int) (Graph3, error)
	GetGraph7(lotteryType int, permutationKey int) (Graph3, error)
	GetGraph8(lotteryType int, fourNumber int) (Graph8, error)
}

type galleryGorm struct {
	db *gorm.DB
}

type galleryService struct {
	GalleryDB
}

type galleryValidator struct {
	GalleryDB
}

type galleryValFn func(*Gallery) error

func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}

func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryService{
		GalleryDB: &galleryValidator{
			GalleryDB: &galleryGorm{
				db: db},
		}}
}

func runGalleryValFns(gallery *Gallery, fns ...galleryValFn) error {
	for _, fn := range fns {
		if err := fn(gallery); err != nil {
			return err
		}
	}
	return nil
}

func (gv *galleryValidator) userIDRequired(g *Gallery) error {
	return nil
}

func (gv *galleryValidator) titleRequired(g *Gallery) error {
	return nil
}

func (gv *galleryValidator) Create(gallery *Gallery) error {
	err := runGalleryValFns(gallery)
	if err != nil {
		return err
	}
	return gv.GalleryDB.Create(gallery)
}

func (gg *galleryGorm) ByID(id uint) (*Gallery, error) {
	var gallery Gallery
	db := gg.db.Where("id = ?", id)
	err := first(db, &gallery)
	if err != nil {
		return nil, err
	}
	return &gallery, nil
}

func (gg *galleryGorm) GetLast(lotteryType int) (*Gallery, error) {
	var gallery Gallery
	gg.db.Order("no desc").Where("lottery_type = ?", lotteryType).Find(&gallery)
	return &gallery, nil
}

func (gg *galleryGorm) GetAllPs(lotteryType int) ([]Ps, error) {
	var ps []Ps
	var results []Gallery
	var idx int
	gg.db.Where("lottery_type = ?", lotteryType).Order("no asc").FindInBatches(&results, 10000, func(tx *gorm.DB, batch int) error {
		// batch processing found records
		for _, result := range results {
			idx++
			ps = append(ps, Ps{
				Date:  result.Date,
				No:    result.No,
				Index: idx,
				Num1:  result.Num1,
				Num2:  result.Num2,
				Num3:  result.Num3,
			})
		}

		return nil
	})

	return ps, nil
}

func (gg *galleryGorm) GetGraph2(lotteryType int, permutationKey int) ([]Graph2, error) {
	srcData, _ := gg.GetAllPs(lotteryType)
	var res []Graph2
	var state [5]Graph2
	for _, v := range srcData {
		var target int
		if v.Num1 == v.Num2 && v.Num1 == v.Num3 {
			continue
		}
		if v.Num1 == v.Num2 || v.Num1 == v.Num3 || v.Num2 == v.Num3 {
			target = 1
		} else {
			target = 2
		}
		var cnt int
		for _, value := range permutation[permutationKey][:5] {
			if v.Num1 == value || v.Num2 == value || v.Num3 == value {
				cnt++
			}
		}
		if target == 1 && cnt == 2 {
			cnt++
		}
		if cnt != 0 && cnt != 3 {
			continue
		}
		if cnt == 0 {
			target += 2
		}

		g := Graph2{
			Date:     v.Date,
			No:       v.No,
			Index:    v.Index,
			Num1:     v.Num1,
			Num2:     v.Num2,
			Num3:     v.Num3,
			Omission: v.Index - state[0].Index,
		}

		switch target {
		case 1:
			g.Diff1 = g.Index - state[target].Index
		case 2:
			g.Diff2 = g.Index - state[target].Index
		case 3:
			g.Diff3 = g.Index - state[target].Index
		case 4:
			g.Diff4 = g.Index - state[target].Index
		}
		state[0] = g
		state[target] = g
		res = append(res, g)
	}
	return res, nil
}

func (gg *galleryGorm) GetGraph3(lotteryType int, permutationKey int) (Graph3, error) {
	g2, _ := gg.GetGraph2(lotteryType, permutationKey)
	var temp [4][]int
	for _, v := range g2 {
		if v.Diff1 > 0 {
			temp[0] = append(temp[0], v.Diff1)
		}
		if v.Diff2 > 0 {
			temp[1] = append(temp[1], v.Diff2)
		}
		if v.Diff3 > 0 {
			temp[2] = append(temp[2], v.Diff3)
		}
		if v.Diff4 > 0 {
			temp[3] = append(temp[3], v.Diff4)
		}
	}
	var res [][]int
	for i, arr := range temp {
		var t = make([]int, RowCount)
		for i := range t {
			t[i] = -1
		}
		t[0] = i + 1
		res = append(res, t)
		for i := 0; i < len(arr); i += RowCount {
			ends := i + RowCount
			if ends > len(arr) {
				ends = len(arr)
			}
			res = append(res, arr[i:ends])
		}
	}
	return Graph3{res}, nil
}

func (gg *galleryGorm) GetGraph4(lotteryType int, permutationKey int) (Graph3, error) {
	g2, _ := gg.GetGraph2(lotteryType, permutationKey)
	var temp [4][]int
	for _, v := range g2 {
		if v.Diff1 > 19 {
			temp[0] = append(temp[0], v.Diff1)
		}
		if v.Diff2 > 19 {
			temp[1] = append(temp[1], v.Diff2)
		}
		if v.Diff3 > 19 {
			temp[2] = append(temp[2], v.Diff3)
		}
		if v.Diff4 > 19 {
			temp[3] = append(temp[3], v.Diff4)
		}
	}
	var res [][]int
	for i, arr := range temp {
		var t = make([]int, RowCount)
		for i := range t {
			t[i] = -1
		}
		t[0] = i + 1
		res = append(res, t)
		for i := 0; i < len(arr); i += RowCount {
			ends := i + RowCount
			if ends > len(arr) {
				ends = len(arr)
			}
			res = append(res, arr[i:ends])
		}
	}
	return Graph3{res}, nil
}

func (gg *galleryGorm) GetGraph5(lotteryType int, permutationKey int) ([]Graph5, error) {
	srcData, _ := gg.GetAllPs(lotteryType)
	var res []Graph5
	var state [10]Graph5
	for _, v := range srcData {
		if v.Num1 == v.Num2 || v.Num1 == v.Num3 || v.Num2 == v.Num3 {
			continue
		}
		var cnt int
		var target int
		endfor := false
		for i, value := range permutation[permutationKey] {
			if v.Num1 == value || v.Num2 == value || v.Num3 == value {
				cnt++
				target = i
			}
			if i == 4 {
				if cnt == 0 || cnt == 3 {
					endfor = true
				}
			}
			if i == 4 && cnt == 1 {
				break
			}
		}

		if endfor {
			continue
		}

		g := Graph5{
			Date:  v.Date,
			No:    v.No,
			Index: v.Index,
			Num1:  v.Num1,
			Num2:  v.Num2,
			Num3:  v.Num3,
		}

		switch target {
		case 0:
			g.Arr0 = g.Index - state[target].Index
		case 1:
			g.Arr1 = g.Index - state[target].Index
		case 2:
			g.Arr2 = g.Index - state[target].Index
		case 3:
			g.Arr3 = g.Index - state[target].Index
		case 4:
			g.Arr4 = g.Index - state[target].Index
		case 5:
			g.Arr5 = g.Index - state[target].Index
		case 6:
			g.Arr6 = g.Index - state[target].Index
		case 7:
			g.Arr7 = g.Index - state[target].Index
		case 8:
			g.Arr8 = g.Index - state[target].Index
		case 9:
			g.Arr9 = g.Index - state[target].Index
		}
		state[target] = g
		res = append(res, g)
	}
	return res, nil
}

func (gg *galleryGorm) GetGraph6(lotteryType int, permutationKey int) (Graph3, error) {
	g5, _ := gg.GetGraph5(lotteryType, permutationKey)
	var temp [10][]int
	for _, v := range g5 {
		if v.Arr0 > 0 {
			temp[0] = append(temp[0], v.Arr0)
		}
		if v.Arr1 > 0 {
			temp[1] = append(temp[1], v.Arr1)
		}
		if v.Arr2 > 0 {
			temp[2] = append(temp[2], v.Arr2)
		}
		if v.Arr3 > 0 {
			temp[3] = append(temp[3], v.Arr3)
		}
		if v.Arr4 > 0 {
			temp[4] = append(temp[4], v.Arr4)
		}
		if v.Arr5 > 0 {
			temp[5] = append(temp[5], v.Arr5)
		}
		if v.Arr6 > 0 {
			temp[6] = append(temp[6], v.Arr6)
		}
		if v.Arr7 > 0 {
			temp[7] = append(temp[7], v.Arr7)
		}
		if v.Arr8 > 0 {
			temp[8] = append(temp[8], v.Arr8)
		}
		if v.Arr9 > 0 {
			temp[9] = append(temp[9], v.Arr9)
		}
	}
	var res [][]int
	for i, arr := range temp {
		var t = make([]int, RowCount)
		for i := range t {
			t[i] = -1
		}
		t[0] = permutation[permutationKey][i]
		res = append(res, t)
		for i := 0; i < len(arr); i += RowCount {
			ends := i + RowCount
			if ends > len(arr) {
				ends = len(arr)
			}
			res = append(res, arr[i:ends])
		}
	}
	return Graph3{res}, nil
}

func (gg *galleryGorm) GetGraph7(lotteryType int, permutationKey int) (Graph3, error) {
	g5, _ := gg.GetGraph5(lotteryType, permutationKey)
	var temp [10][]int
	for _, v := range g5 {
		if v.Arr0 > 19 {
			temp[0] = append(temp[0], v.Arr0)
		}
		if v.Arr1 > 19 {
			temp[1] = append(temp[1], v.Arr1)
		}
		if v.Arr2 > 19 {
			temp[2] = append(temp[2], v.Arr2)
		}
		if v.Arr3 > 19 {
			temp[3] = append(temp[3], v.Arr3)
		}
		if v.Arr4 > 19 {
			temp[4] = append(temp[4], v.Arr4)
		}
		if v.Arr5 > 19 {
			temp[5] = append(temp[5], v.Arr5)
		}
		if v.Arr6 > 19 {
			temp[6] = append(temp[6], v.Arr6)
		}
		if v.Arr7 > 19 {
			temp[7] = append(temp[7], v.Arr7)
		}
		if v.Arr8 > 19 {
			temp[8] = append(temp[8], v.Arr8)
		}
		if v.Arr9 > 19 {
			temp[9] = append(temp[9], v.Arr9)
		}
	}
	var res [][]int
	for i, arr := range temp {
		var t = make([]int, RowCount)
		for i := range t {
			t[i] = -1
		}
		t[0] = permutation[permutationKey][i]
		res = append(res, t)
		for i := 0; i < len(arr); i += RowCount {
			ends := i + RowCount
			if ends > len(arr) {
				ends = len(arr)
			}
			res = append(res, arr[i:ends])
		}
	}
	return Graph3{res}, nil
}

const (
	plsOffset    = 22000
	threedOffset = 2022000
)

func (gg *galleryGorm) GetGraph8(lotteryType int, fourNumber int) (Graph8, error) {
	var offset int
	switch lotteryType {
	case 1:
		offset = plsOffset
	case 2:
		offset = threedOffset
	}
	log.Println("fournm", fourNumber)
	var res [][]string
	srcData, _ := gg.GetAllPs(lotteryType)
	var dateArr, noArr, numArr []string
	var config [][]string
	// Date  string
	// No    int
	// Index int
	// Num1  int
	// Num2  int
	// Num3  int
	for _, v := range srcData {
		if v.No < offset {
			continue
		}
		dateArr = append(dateArr, v.Date)
		noArr = append(noArr, strconv.Itoa(v.No))
		numArr = append(numArr, strconv.Itoa(v.Num1)+strconv.Itoa(v.Num2)+strconv.Itoa(v.Num3))
		tmp := []int{v.Num1, v.Num2, v.Num3}
		sort.Ints(tmp)
		config = append(config, []string{
			strconv.Itoa(tmp[0]) + strconv.Itoa(tmp[1]),
			strconv.Itoa(tmp[0]) + strconv.Itoa(tmp[2]),
			strconv.Itoa(tmp[1]) + strconv.Itoa(tmp[2]),
		})
	}

	emptyFields := []string{"", "", "", "", ""}

	res = append(res, append(emptyFields, dateArr...))
	res = append(res, append(emptyFields, noArr...))
	res = append(res, append(emptyFields, numArr...))

	headers := gg.getGraph8Headers(fourNumber)
	for _, v := range headers {
		res = append(res, gg.getGraph8Content(config, v))
	}

	res = transpose(res)

	return Graph8{res}, nil
}

func (gg *galleryGorm) getGraph8Headers(a int) [][]string {
	// log.Println(a)
	// println(a / 1000 % 10)
	n1 := a / 1000 % 10
	n2 := a / 100 % 10
	n3 := a / 10 % 10
	n4 := a / 1 % 10
	// println(n1, n2, n3, n4)

	pair1 := strconv.Itoa(n1) + strconv.Itoa(n2)
	pair2 := strconv.Itoa(n3) + strconv.Itoa(n4)

	// log.Println(pair1, pair2)

	var tmpArr []int
	for i := n3 + 1; i <= 9; i++ {
		if i != n4 {
			tmpArr = append(tmpArr, i)
		}
	}
	// log.Println(tmpArr)

	var pair3 []string
	for i := 0; i < len(tmpArr)-1; i++ {
		for j := i + 1; j < len(tmpArr); j++ {
			pair3 = append(pair3, strconv.Itoa(tmpArr[i])+strconv.Itoa(tmpArr[j]))
		}
	}

	// log.Println(pair3)

	var pairs [][]string

	for _, v := range pair3 {
		pair, _ := strconv.Atoi(v)
		n5 := pair / 10 % 10
		n6 := pair % 10

		tmpArr = nil
		for i := 0; i <= 9; i++ {
			switch i {
			case n1,
				n2,
				n3,
				n4,
				n5,
				n6:
			default:
				tmpArr = append(tmpArr, i)
			}
		}

		// log.Println(tmpArr)

		var pair4, pair5 string
		pair4 = strconv.Itoa(tmpArr[0]) + strconv.Itoa(tmpArr[1])
		pair5 = strconv.Itoa(tmpArr[2]) + strconv.Itoa(tmpArr[3])
		pairs = append(pairs, []string{
			pair1, pair2, v, pair4, pair5,
		})

		pair4 = strconv.Itoa(tmpArr[0]) + strconv.Itoa(tmpArr[2])
		pair5 = strconv.Itoa(tmpArr[1]) + strconv.Itoa(tmpArr[3])
		pairs = append(pairs, []string{
			pair1, pair2, v, pair4, pair5,
		})

		pair4 = strconv.Itoa(tmpArr[0]) + strconv.Itoa(tmpArr[3])
		pair5 = strconv.Itoa(tmpArr[1]) + strconv.Itoa(tmpArr[2])
		pairs = append(pairs, []string{
			pair1, pair2, v, pair4, pair5,
		})
	}

	// log.Println(pairs)
	return pairs
}

func (gg *galleryGorm) getGraph8Content(config [][]string, header []string) []string {
	headerSet := make(map[string]bool)
	for _, v := range header {
		headerSet[v] = true
	}
	// log.Println(config, len(config))
	res := make([]string, 0, len(config))
	var lastHit int
	for i, v := range config {
		var empty bool
		for _, val := range v {
			empty = true
			if headerSet[val] {
				empty = false
				res = append(res, strconv.Itoa(i-lastHit))
				lastHit = i
				break
			}
		}
		if empty {
			res = append(res, "")
		}
	}
	// log.Println("res", len(res))
	// log.Println(len(append(header, res...)))
	return append(header, res...)
}

// first will query using the provided gorm.DB and it will
// get the first item returned and place it into dst. If
// nothing is found in the query, it will return ErrNotFound
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

func transpose(matrix [][]string) [][]string {
	n, m := len(matrix), len(matrix[0])
	t := make([][]string, m)
	for i := range t {
		t[i] = make([]string, n)
		for j := range t[i] {
			t[i][j] = "hh"
		}
	}
	for i, row := range matrix {
		for j, v := range row {
			t[j][i] = v
		}
	}
	return t
}
