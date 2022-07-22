package tests

import (
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/steinfletcher/apitest"

	"github.com/KameeKaze/Ticketing-system/routes"
)

func init() {
	rand.Seed(time.Now().UnixNano())

	ran_str := make([]byte, 12)

	// Generating Random string
	for i := 0; i < 12; i++ {
		ran_str[i] = byte(33 + rand.Intn(93))
	}

	// Displaying the random string
	USERNAME = string(ran_str)
}

var (
	USERNAME string
)

func TestHome(t *testing.T) {
	apitest.New(). // configuration
			HandlerFunc(routes.Home).
			Get("/").  // request
			Expect(t). // expectations
			Body(`{"Message": "Ticketing system"}`).
			Status(http.StatusOK).
			End()
}

func TestSignUp(t *testing.T) {
	apitest.New(). // configuration
			HandlerFunc(routes.Register).
			Post("/register").JSON(`{
					"username":"` + USERNAME + `",
					"password":"secretpassword123",
					"role":"programmer"
				}`).
		Expect(t). // expectations
		Body(`{"Message": "Creating user ` + USERNAME + `"}`).
		Status(http.StatusCreated).
		End()
}
