package auth

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/csunibo/polleg/util"
)

func WhoAmIHandler(res http.ResponseWriter, req *http.Request) {
	user := GetUser(req)
	token, _ := req.Context().Value("token").(string)
	fmt.Println(token)
	if err := util.WriteJson(res, user); err != nil {
		_ = util.WriteError(res, http.StatusInternalServerError, "")
		slog.Error("could not encode json:", "error", err)
	}
}
