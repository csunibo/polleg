package answers

import (
	"log/slog"
	"net/http"

	"github.com/csunibo/stackunibo/util"
	"github.com/kataras/muxie"
)

func AnswerHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		handleGet(res, req)
	case http.MethodPut:
		handlePut(res, req)
	default:
		_ = util.WriteError(res, http.StatusMethodNotAllowed, "invalid method")
	}
}

func handlePut(res http.ResponseWriter, req *http.Request) {
}

func handleGet(res http.ResponseWriter, req *http.Request) {
	answer := muxie.GetParam(res, "id")
	slog.Info("Fetching answers", "doc", answer)
}
