package users

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		token = strings.Replace(token, "Bearer ", "", -1)

		if !validateToken(token) {
			fmt.Println("InvalidToken")
		} else {
			fmt.Println("Valid User!!")
		}

		next.ServeHTTP(rw, r)
	})
}

var hmacSampleSecret []byte

func validateToken(stringToken string) bool {

	pubKey, err := ioutil.ReadFile("/home/adrian.castan/cert/id_rsa.pub")
	if err != nil {
		log.Fatalln(err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(pubKey)

	if err != nil {
		return false
	}

	tok, err := jwt.Parse(stringToken, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return key, nil
	})
	if err != nil {
		return false
	}

	_, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return false
	}

	return true

}
