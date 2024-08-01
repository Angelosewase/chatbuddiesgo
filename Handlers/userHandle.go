package Handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Angelosewase/chatbuddiesgo/internal/database"
	// "github.com/google/uuid"
)

type ApiCfg struct {
	DB *database.Queries
}

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

	_, err := db.CreateUser(req.Context(), database.CreateUserParams{
		ID:        "",
		Firstname: sql.NullString{Valid: true, String: parameters.FirstName},
		Lastname:  sql.NullString{Valid: true, String: parameters.LastName},
		Email:     parameters.Email,
		Password:  parameters.Password,
	})

	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}
	return nil
}

func SignUpHandler(db *database.Queries) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		if err := SignUp(res, req, db); err != nil {
			log.Fatal(err)
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

	parameters := Parameters{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&parameters)

	if err != nil {
		return fmt.Errorf("failed parsing the request body %v", err)
	}

	//database logic to fectch the user with the given parameters 
	//generate the jwt token using the user credentials
	//sending the cookies to the client containing th jwt token 

	return nil

}

func LoginHandler(db *database.Queries)(func (res http.ResponseWriter,req *http.Request)){

	return func(res http.ResponseWriter,req *http.Request){
		err:=LogIn(res,req,db)
		if err !=nil{
           res.Write([]byte("Error logging in "))
		}
	}
}
