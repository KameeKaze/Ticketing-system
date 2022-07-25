package tests

import (
	"math/rand"
	"net/http"
	"testing"
	"time"
	"fmt"

	"github.com/steinfletcher/apitest"

	"github.com/KameeKaze/Ticketing-system/routes"
)

func init() {
	//set seed for random generator
	rand.Seed(time.Now().UnixNano())


	ran_str := make([]byte, 12)

	// Generating Random string
	for i := 0; i < 12; i++ {
		ran_str[i] = byte(40 + rand.Intn(82))
	}

	//create random username for test
	USERNAME = string(ran_str)

	for i := 0; i < 12; i++ {
		ran_str[i] = byte(40 + rand.Intn(82))
	}
	USERNAME2 = string(ran_str)

}

var (
	USERNAME  string
	USERNAME2 string
)

func TestHome(t *testing.T) {
	apitest.New().
			HandlerFunc(routes.Home).
			Get("/").  // request
			Expect(t).
			Body(`{"Message": "Ticketing system"}`).
			Status(http.StatusOK).
			End()
}

func TestSignUp(t *testing.T) {
	// create user
	apitest.New().
		HandlerFunc(routes.SignUp).
		Post("/register").JSON(fmt.Sprintf(
			`{
				"username":"%s",
				"password":"secretpassword123",
				"role":"programmer"
			}`,USERNAME)).
		Expect(t).
		Body(fmt.Sprintf(
			`{"Message": "Creating user %s"}`,USERNAME)).
		Status(http.StatusCreated).
		End()

	//create same username
	apitest.New().
		HandlerFunc(routes.SignUp).
		Post("/register").JSON(fmt.Sprintf(
			`{
				"username":"%s",
				"password":"secretpassword123",
				"role":"programmer"
			}`,USERNAME)).
		Expect(t).
		Body(`{"Message": "Username already taken"}`).
		Status(http.StatusConflict).
		End()

	// invalid post data
	apitest.New().
		HandlerFunc(routes.SignUp).
			Post("/register").JSON(fmt.Sprintf(
				`{
					"username":"",
					"":"programmer"
				}`)).
		Expect(t).
		Body(`{"Message": "Invalid request"}`).
		Status(http.StatusBadRequest).
		End()
	//empty requset
	apitest.New().
		HandlerFunc(routes.SignUp).
		Post("/register").JSON(``).
		Expect(t).
		Body(`{"Message": "Invalid request"}`).
		Status(http.StatusBadRequest).
		End()
}

func TestLogin(t *testing.T) {
	// login existring user
	apitest.New().
		HandlerFunc(routes.Login).
		Post("/login").JSON(fmt.Sprintf(
			`{
				"username":"%s",
				"password":"secretpassword123"
			}`,USERNAME)).
		Expect(t).
		Body(fmt.Sprintf(
			`{"Message": "Logging in %s"}`,USERNAME)).
		Status(http.StatusOK).
		End()

	// login non existring user
	apitest.New().
		HandlerFunc(routes.Login).
		Post("/login").JSON(fmt.Sprintf(
			`{
				"username":"%s",
				"password":"secretpassword123"
			}`,USERNAME2)).
		Expect(t).
		Body(fmt.Sprintf(
			`{"Message": "Invalid credentials"}`)).
		Status(http.StatusUnauthorized).
		End()
	// empty body
	apitest.New().
		HandlerFunc(routes.Login).
		Post("/login").JSON(fmt.Sprintf(
			``)).
		Expect(t).
		Body(fmt.Sprintf(
			`{"Message": "Invalid request"}`)).
		Status(http.StatusBadRequest).
		End()

	// invalid post data
	apitest.New().
		HandlerFunc(routes.Login).
		Post("/login").JSON(fmt.Sprintf(
			`{
				"s":"dsa",
				"username":"",
			}`)).
		Expect(t).
		Body(fmt.Sprintf(
			`{"Message": "Invalid request"}`)).
		Status(http.StatusBadRequest).
		End()
}
