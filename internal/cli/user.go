package cli

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/alaleks/geospace/internal/server/database/models"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
)

const (
	// commands
	commandSignUp      = "signup"
	commandSignUpShort = "sign"
	commandLogIn       = "login"
	commandLogInAlias  = "signin"
	// urls
	urlSignUp = "/v1/signup"
	urlLogin  = "/v1/login"
)

var (
	Username string // name of the user
	Password string // password of the user
	Email    string // email of the user
)

var sign = &cobra.Command{
	Use:     commandSignUp,
	Aliases: []string{commandSignUpShort},
	Short:   "application registration.",
	Run: func(cmd *cobra.Command, args []string) {
		token, err := signUp()
		if err != nil {
			log.Fatal(err)
		}

		if !JSONResponse {
			fmt.Println("Access token:")
		}
		fmt.Println(token)
	},
}

var signin = &cobra.Command{
	Use:     commandLogIn,
	Aliases: []string{commandLogInAlias},
	Short:   "login in application.",
	Run: func(cmd *cobra.Command, args []string) {
		token, err := login()
		if err != nil {
			log.Fatal(err)
		}

		if !JSONResponse {
			fmt.Println("Access token:")
		}
		fmt.Println(token)
	},
}

func init() {
	// sign up
	sign.Flags().StringVarP(&Username, "username", "u", "", "name of the user")
	sign.Flags().StringVarP(&Password, "password", "p", "", "name of the user")
	sign.Flags().StringVarP(&Email, "email", "e", "", "email of the user")
	sign.Flags().BoolVarP(&JSONResponse, "json", "j", false, "make a JSON response")
	rootCmd.AddCommand(sign)

	// sign in
	signin.Flags().StringVarP(&Password, "password", "p", "", "name of the user")
	signin.Flags().StringVarP(&Email, "email", "e", "", "email of the user")
	signin.Flags().BoolVarP(&JSONResponse, "json", "j", false, "make a JSON response")
	rootCmd.AddCommand(signin)
}

// signUp performs the registration in application.
func signUp() (any, error) {
	user := models.User{
		Name:     Username,
		Password: Password,
		Email:    Email,
	}

	data, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	agent := fiber.AcquireAgent()

	req := agent.Request()
	req.SetBody(data)
	req.Header.SetMethod(fiber.MethodPost)
	req.Header.SetContentType(contentTypeJSON)
	req.SetRequestURI(cliHost + urlSignUp)

	if err := agent.Parse(); err != nil {
		return "", err
	}

	code, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return "", fmt.Errorf("%v", errs)
	}

	if code != 200 {
		return "", fmt.Errorf(string(body))
	}

	if JSONResponse {
		return prettyJSON(body), nil
	}

	var response struct {
		Token string `json:"token"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("invalid authentication: %v", err)
	}

	return response.Token, nil
}

// login provides capability of log in a user.
func login() (any, error) {
	user := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    Email,
		Password: Password,
	}

	data, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	agent := fiber.AcquireAgent()

	req := agent.Request()
	req.SetBody(data)
	req.Header.SetMethod(fiber.MethodPost)
	req.Header.SetContentType(contentTypeJSON)
	req.SetRequestURI(cliHost + urlLogin)

	if err := agent.Parse(); err != nil {
		return "", err
	}

	code, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return "", fmt.Errorf("%v", errs)
	}

	if code != 200 {
		return "", fmt.Errorf(string(body))
	}

	if JSONResponse {
		return prettyJSON(body), nil
	}

	var response struct {
		Token string `json:"token"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("invalid authentication: %v", err)
	}

	return response.Token, nil
}
