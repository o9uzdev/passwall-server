package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/pass-wall/passwall-server/internal/app"
	"github.com/pass-wall/passwall-server/internal/storage"
	"github.com/pass-wall/passwall-server/model"

	"github.com/gorilla/mux"
)

// FindAll ...
func FindAllLogins(s storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		logins := []model.Login{}

		fields := []string{"id", "created_at", "updated_at", "url", "username"}
		argsStr, argsInt := SetArgs(r, fields)

		logins, err = s.Logins().FindAll(argsStr, argsInt)

		if err != nil {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		logins = app.DecryptLoginPasswords(logins)
		RespondWithJSON(w, http.StatusOK, logins)
	}
}

// FindByID ...
func FindLoginsByID(s storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		login, err := s.Logins().FindByID(uint(id))
		if err != nil {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		uLogin, err := app.DecryptLoginPassword(s, &login)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		RespondWithJSON(w, http.StatusOK, model.ToLoginDTO(uLogin))
	}
}

// Create ...
func CreateLogin(s storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginDTO model.LoginDTO

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&loginDTO); err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
			return
		}
		defer r.Body.Close()

		if loginDTO.Password == "" {
			loginDTO.Password = app.Password()
		}

		createdLogin, err := app.CreateLogin(s, &loginDTO)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		RespondWithJSON(w, http.StatusOK, model.ToLoginDTO(createdLogin))
	}
}

// Update ...
func UpdateLogin(s storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		var loginDTO model.LoginDTO
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&loginDTO); err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
			return
		}
		defer r.Body.Close()

		login, err := s.Logins().FindByID(uint(id))
		if err != nil {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		updatedLogin, err := app.UpdateLogin(s, &login, &loginDTO)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		RespondWithJSON(w, http.StatusOK, model.ToLoginDTO(updatedLogin))
	}
}

// Delete ...
func DeleteLogin(s storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		login, err := s.Logins().FindByID(uint(id))
		if err != nil {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		err = s.Logins().Delete(login.ID)
		if err != nil {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		response := model.Response{
			Code:    http.StatusOK,
			Status:  "Success",
			Message: "Login deleted successfully!",
		}
		RespondWithJSON(w, http.StatusOK, response)
	}
}
