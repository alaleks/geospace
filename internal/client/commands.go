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

// getNearbyCities provides capability of get nearby cities.
func (c *Client) getNearbyCities() error {
	departure, err := inputWithReslult("Departure  city*")
	if err != nil {
		return err
	}

	distanceTo, err := inputWithReslult("Distance to (in km)*")
	if err != nil {
		return err
	}

	req := c.Agent.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.Header.Add("Authorization", "Bearer "+c.Token)
	req.SetRequestURI(c.Host + "/v1/user/find-by-name")
	req.URI().QueryArgs().Add("departure", departure)
	req.URI().QueryArgs().Add("distanceTo", distanceTo)

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

// getNearbyCitiesbyCoord provides capability of get nearby cities by coordinates.
func (c *Client) getNearbyCitiesbyCoord() error {
	lat, err := inputWithReslult("Latitude, decimal number*")
	if err != nil {
		return err
	}

	lon, err := inputWithReslult("Longitude, decimal number*")
	if err != nil {
		return err
	}

	distanceTo, err := inputWithReslult("Distance to (in km)*")
	if err != nil {
		return err
	}

	req := c.Agent.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.Header.Add("Authorization", "Bearer "+c.Token)
	req.SetRequestURI(c.Host + "/v1/user/find-by-coord")
	req.URI().QueryArgs().Add("lat", lat)
	req.URI().QueryArgs().Add("lon", lon)
	req.URI().QueryArgs().Add("distanceTo", distanceTo)

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
