package Handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Angelosewase/chatbuddiesgo/helpers"
	"github.com/Angelosewase/chatbuddiesgo/internal/auth"
	"github.com/Angelosewase/chatbuddiesgo/internal/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(res http.ResponseWriter, req *http.Request, db *database.Queries) error {
	type Parameters struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	parameters := Parameters{}

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&parameters); err != nil {
		return fmt.Errorf("failed to parse the request body %V", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(parameters.Password), 10)
	if err != nil {
		return fmt.Errorf("error hashing password %v", err)
	}

	_, err = db.CreateUser(req.Context(), database.CreateUserParams{
		ID:        uuid.NewString(),
		Firstname: sql.NullString{Valid: true, String: parameters.FirstName},
		Lastname:  sql.NullString{Valid: true, String: parameters.LastName},
		Email:     parameters.Email,
		Password:  string(hashedPassword),
	})

	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}
	return nil
}

func SignUpHandler(db *database.Queries) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		if err := SignUp(res, req, db); err != nil {
			res.WriteHeader(http.StatusInternalServerError)

			errres, err := json.Marshal(err)
			if err != nil {
				res.Write([]byte("internal server error"))
				return
			}
			// res.Write([]byte("Error creating user"))
			res.Write(errres)
			return
		}

		res.WriteHeader(200)
		res.Write([]byte("user created successfully"))
	}
}

func LogIn(res http.ResponseWriter, req *http.Request, db *database.Queries) error {
	type Parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type User struct {
		Id         string
		First_name string
		Last_name  string
		Email      string
	}

	parameters := Parameters{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&parameters)

	if err != nil {
		return fmt.Errorf("failed parsing the request body %v", err)
	}

	user, err := db.GetUser(req.Context(), parameters.Email)

	if err != nil {
		return fmt.Errorf("error fetching user:%v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(parameters.Password))

	if err != nil {
		return fmt.Errorf("incorrect Password: %v", err)
	}
	token, err := auth.GenerateJwtToken(user.ID)

	if err != nil {
		return fmt.Errorf("error generating token :%v", err)

	}

	http.SetCookie(res, &http.Cookie{
		Expires:  time.Now().Add(time.Hour * 24),
		Name:     "jwt",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
	})

	resUser := User{
		Id:         user.ID,
		First_name: user.Firstname.String,
		Last_name:  user.Lastname.String,
		Email:      user.Email,
	}

	jsonuser, err := json.Marshal(resUser)
	if err != nil {
		return fmt.Errorf("error converting user into json: %v", err)
	}

	res.Write(jsonuser)
	return nil

}

func LoginHandler(db *database.Queries) func(res http.ResponseWriter, req *http.Request) {

	return func(res http.ResponseWriter, req *http.Request) {
		err := LogIn(res, req, db)
		if err != nil {
			res.Write([]byte(fmt.Sprintf("error logging in : %v", err)))
		}
	}
}

func LogoutHandler(res http.ResponseWriter, req *http.Request) {
	http.SetCookie(res, &http.Cookie{
		Expires:  time.Now().Add(-time.Hour),
		Name:     "jwt",
		Value:    "",
		HttpOnly: true,
	})
}

func GetUserByUserId(db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Parameters struct {
			UserId string `json:"userId"`
		}
		parameters := Parameters{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&parameters)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := db.GetUserByUserId(r.Context(), parameters.UserId)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		type User struct {
			Id         string
			First_name string
			Last_name  string
			Email      string
		}

		resUser := User{
			Id:         user.ID,
			First_name: user.Firstname.String,
			Last_name:  user.Lastname.String,
			Email:      user.Email,
		}

		err = helpers.RespondWithJson(w, r, resUser, 200)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func SearchHandler(db *database.Queries) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		querryParams := r.URL.Query()
		userName := querryParams.Get("q")

		users, err := db.GetUserByName(r.Context(), database.GetUserByNameParams{
			CONCAT:   userName,
			CONCAT_2: userName,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}
		helpers.RespondWithJson(w, r, users, http.StatusAccepted)
	}
}

func IsLoggedIn(db *database.Queries) http.HandlerFunc {
	type User struct {
		Id         string 
		First_name string
		Last_name  string
		Email      string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		jwtCookie, err := r.Cookie("jwt")
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		userId, err := auth.ValidateToken(jwtCookie.Value)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := db.GetUserByUserId(r.Context(), userId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resUser := User{
			Id:         user.ID,
			First_name: user.Firstname.String,
			Last_name:  user.Lastname.String,
			Email:      user.Email,
		}

		err = helpers.RespondWithJson(w, r, resUser, 200)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}
