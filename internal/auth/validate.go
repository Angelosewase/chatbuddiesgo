package auth

import (
	"errors"
	"fmt"
	"net/http"
)

func ValidateUserCredentials(res http.ResponseWriter, req *http.Request) error {

	 cookie,err:=req.Cookie("jwt")
	 if err != nil{
		return errors.New("user not authorised")
	 }
     _,err=ValidateToken(cookie.Value)

	 if err !=nil{
		return  fmt.Errorf("invalid credentials : %v",err)
	 }

	 return nil
     
}

