package auth

import (
	"encoding/json"
	api_helper "mpt_data/api/apihelper"
	"mpt_data/database/auth"
	apiModel "mpt_data/models/apimodel"

	"net/http"

	"github.com/gorilla/mux"
)

const packageName = "api.auth"

// RegisterRoutes adds all routes to a mux.Router
func RegisterRoutes(mux *mux.Router) {
	mux.HandleFunc(apiModel.LoginHref, login).Methods(http.MethodPost)
	mux.HandleFunc(apiModel.UserChangePWHref, CheckAuthentication(changePW)).Methods(http.MethodPost)
}

// @Summary		Login
// @Description	Login to Service
// @Tags			Users
// @Accept			json
// @Produce		json
// @Param			Auth-Information	body	apiModel.UserLogin	true	"Auth Information"
// @Type
// @Success	200
//
// @Header		200	{string}	Authorization	"Bearer-Token"
// @Router		/login [POST]
func login(w http.ResponseWriter, r *http.Request) {
	var user apiModel.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(401)
		return
	}

	jwt, err := auth.Login(user)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	w.Header().Add("Authorization", "Bearer "+jwt)
	w.WriteHeader(http.StatusOK)

}

// @Summary		Change Password
// @Description	Change the Password to login
// @Tags			Users
// @Accept			json
// @Produce		json
// @Param			Auth-Information	body	auth.changePW.pw	true	"Auth Information"
// @Security		ApiKeyAuth
// @Success		200	{object}	apiModel.Result
// @Router			/user/password [POST]
func changePW(w http.ResponseWriter, r *http.Request) {
	type pw struct {
		Password string
	}
	const funcName = packageName + ".changePW"
	id, err := auth.GetUserIDFromToken(r.Header.Get("Authorization"))
	if err != nil {
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "password not changed"}, err)
		return
	}
	var user pw
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		api_helper.ResponseBadRequest(w, funcName, apiModel.Result{Result: "missing data"}, err)
		return
	}

	if err := auth.ChangePassword(id, user.Password); err != nil {
		api_helper.InternalError(w, funcName, err.Error())
	}
	api_helper.ResponseJSON(w, funcName, apiModel.Result{Result: "password changed succesfull"})
}
