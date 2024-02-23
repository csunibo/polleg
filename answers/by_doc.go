package answers

import (
	"log/slog"
	"net/http"

	"github.com/kataras/muxie"
)

func ByDoc(res http.ResponseWriter, req *http.Request) {
	doc := muxie.GetParam(res, "id")
	slog.Info("Fetching answers", "doc", doc)
}
