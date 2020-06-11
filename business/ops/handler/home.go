package handler

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

// Home 首页
func Home(w http.ResponseWriter, r *http.Request) {
	switch r.FormValue("action") {
	case "user":
		data := make([]struct {
			Y     string          `json:"y"`
			Item1 decimal.Decimal `json:"item1"`
			Item2 decimal.Decimal `json:"item2"`
		}, 15)
		n := time.Now().AddDate(0, 0, -len(data))
		rand.Seed(n.Unix())
		for i := 0; i < len(data); i++ {
			data[i].Y = n.AddDate(0, 0, i).Format(dateFormate)
			data[i].Item1 = decimal.New(rand.Int63n(8000), -2)
			data[i].Item2 = decimal.New(rand.Int63n(4000), -2)
		}
		jSuccess(w, data)
	default:
		rLayout(w, r, "index.tpl", nil)
	}
}
