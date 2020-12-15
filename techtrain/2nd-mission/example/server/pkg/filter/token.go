package filter

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"todo-example/server/pkg/common"
	"todo-example/server/pkg/config"
)

func CheckToken() Decorator {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userToken := r.Header.Get("Authorization")
			if userToken == "" {
				common.RespondErr(w, 401, "Token is required!")
				return
			} else {
				userToken = strings.TrimPrefix(userToken, "Bearer ")
			}
			logger := r.Context().Value("logger").(*common.Log)
			logger.Info(log.Fields{"token": userToken}, "Start token filter")
			if userToken == "" {
				logger.Info(log.Fields{}, "Failed token filter: X-Auth-Token missing")
				common.RespondErr(w, 401, "Access-Token is required")
			} else {
				logger.Debug(log.Fields{}, "Try to validate user token")
				userID, err := validateToken(userToken)
				if err != nil {
					logger.Info(log.Fields{"error": err.Error(), "userToken": userToken}, "Failed token filter")
					common.RespondErr(w, 401, err.Error())
				} else {
					vars := mux.Vars(r)
					if userID != vars["userID"] {
						common.RespondErr(w, 403, "You don't have permission")
					} else {
						logger.Info(log.Fields{"userId": userID}, "Passed token filter")
						h.ServeHTTP(w, r)
					}
				}
			}
		})
	}
}

func validateToken(token string) (string, error) {
	c := config.GetConf()
	req, _ := http.NewRequest("GET", c.LineAPI+"/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	client := new(http.Client)
	rep, err := client.Do(req)
	if err != nil {
		return "", errors.New("line request failed: " + err.Error())
	}
	res, _ := ioutil.ReadAll(rep.Body)
	if rep.StatusCode != 200 {
		return "", errors.New("Invalid token: " + string(res))
	} else {
		return gjson.Get(string(res), "userId").String(), nil
	}
}
