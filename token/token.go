package token

import (
	"errors"
	"os"

	"github.com/gin-gonic/gin"
)

var DbLogins = getDbLogins()

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getDbLogins() string {
	f, err := os.Open("token/dblogins.yml")
	check(err)

	_, err = f.Seek(11, 0)
	check(err)
	b := make([]byte, 46)
	_, err = f.Read(b)
	check(err)

	f.Close()

	return string(b)

}

func GetAPIToken() string {

	f, err := os.Open("token/token.yml")
	check(err)

	_, err = f.Seek(11, 0)
	check(err)
	b := make([]byte, 60)
	_, err = f.Read(b)
	check(err)

	f.Close()

	return string(b)
}

func CheckAPIToken(tokenAPI, tokenHeader string, c *gin.Context) error {
	if tokenHeader != tokenAPI {
		return errors.New("Not Autorized")
	}
	return nil
}
