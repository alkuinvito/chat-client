package user

import (
	"encoding/json"
	"net/http"
)

type UserController struct {
	userService *UserService
}

type IUserController interface {
	PairUser() func(w http.ResponseWriter, r *http.Request)
	Routes() http.Handler
}

func NewUserController(userService *UserService) *UserController {
	return &UserController{userService}
}

func (uc *UserController) PairUser() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input RequestPairSchema
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&input)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		pubkey, err := uc.userService.PairUser(input)
		if err != nil {
			http.Error(w, "Unknown error", http.StatusInternalServerError)
			return
		}

		res, err := json.Marshal(&pubkey)
		if err != nil {
			http.Error(w, "Unknown error", http.StatusInternalServerError)
			return
		}

		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}

func (uc *UserController) Routes() http.Handler {
	userRoute := http.NewServeMux()
	userRoute.HandleFunc("POST /pair", uc.PairUser())

	return userRoute
}
