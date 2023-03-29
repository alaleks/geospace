package client

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/pterm/pterm"
)

// calcDistance provides capability of calculate distance between two cities.
func (c *Client) calcDistance() error {
	departure, err := inputWithReslult("Departure  city*")
	if err != nil {
		return err
	}

	destination, err := inputWithReslult("Destination city*")
	if err != nil {
		return err
	}

	req := c.Agent.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.Header.Add("Authorization", "Bearer "+c.Token)
	req.SetRequestURI(c.Host + "/v1/user/distance")
	req.URI().QueryArgs().Add("departure", departure)
	req.URI().QueryArgs().Add("destination", destination)

	if err := c.Agent.Parse(); err != nil {
		return err
	}

	code, body, _ := c.Agent.Bytes()
	if code != 200 {
		return fmt.Errorf(string(body))
	}

	pterm.DefaultHeader.WithFullWidth().Println(string(body))

	return nil
}
