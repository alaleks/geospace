package client

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/pterm/pterm"
)

var (
	ErrInvalidCommand        = errors.New("invalid command specified")
	ErrInvalidAuthentication = errors.New("invalid authentication")
)

const (
	commandLogIn  = "log in"
	commandSignUp = "sign up"
)

// authentication performs sign up or login to app.
func (c *Client) authentication() error {
	options := [...]string{
		commandSignUp,
		commandLogIn,
	}

	printer := pterm.DefaultInteractiveSelect.WithOptions(options[:])
	selectedOptions, err := printer.Show()
	if err != nil {
		return err
	}

	switch selectedOptions {
	case commandSignUp:
		return c.signUp()
	case commandLogIn:
		return c.login()
	default:
		return ErrInvalidCommand
	}
}

// signUp provides capability of registering a user.
func (c *Client) signUp() error {
	var user struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	email, err := inputWithReslult("email*")
	if err != nil {
		return err
	}

	user.Email = email

	name, err := inputWithReslult("name")
	if err != nil {
		return err
	}

	user.Name = name

	password, err := inputWithReslult("password*")
	if err != nil {
		return err
	}

	user.Password = password

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	req := c.Agent.Request()
	req.SetBody(data)
	req.Header.SetMethod(fiber.MethodPost)
	req.Header.SetContentType(contentTypeJSON)
	req.SetRequestURI(c.Host + "/v1/register")

	if err := c.Agent.Parse(); err != nil {
		return err
	}

	code, body, _ := c.Agent.Bytes()
	if code != 200 {
		return fmt.Errorf(string(body))
	}

	var response struct {
		Token string `json:"token"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return ErrInvalidAuthentication
	}

	c.Token = response.Token

	return nil
}

// login provides capability of log in a user.
func (c *Client) login() error {
	var user struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	email, err := inputWithReslult("email*")
	if err != nil {
		return err
	}

	user.Email = email

	password, err := inputWithReslult("password*")
	if err != nil {
		return err
	}

	user.Password = password

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	req := c.Agent.Request()
	req.SetBody(data)
	req.Header.SetMethod(fiber.MethodPost)
	req.Header.SetContentType(contentTypeJSON)
	req.SetRequestURI(c.Host + "/v1/login")

	if err := c.Agent.Parse(); err != nil {
		return err
	}

	code, body, _ := c.Agent.Bytes()
	if code != 200 {
		return fmt.Errorf(string(body))
	}

	var response struct {
		Token string `json:"token"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return ErrInvalidAuthentication
	}

	c.Token = response.Token

	return nil
}
