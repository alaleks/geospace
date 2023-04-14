package client

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/pterm/pterm"
)

const (
	prefixToken             = "Bearer "
	authenticationHeaderKey = "Authorization"
	urlDistance             = "/v1/user/distance"
	urlFindByName           = "/v1/user/find-by-name"
	urlFindByCoord          = "/v1/user/find-by-coord"
	urlGetInfo              = "/v1/api/info"
)

// calcDistance provides capability of calculate distance between two cities.
func (c *Client) calcDistance() error {
	departure, err := inputWithResult("Departure  city*")
	if err != nil {
		return err
	}

	destination, err := inputWithResult("Destination city*")
	if err != nil {
		return err
	}

	req := c.Agent.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.Header.Add(authenticationHeaderKey, prefixToken+c.Token)
	req.SetRequestURI(c.Host + urlDistance)
	req.URI().QueryArgs().Add("departure", departure)
	req.URI().QueryArgs().Add("destination", destination)

	if err := c.Agent.Parse(); err != nil {
		return err
	}

	code, body, errs := c.Agent.Bytes()
	if len(errs) > 0 {
		return fmt.Errorf("%v", errs)
	}

	if code != 200 {
		return fmt.Errorf(string(body))
	}

	pterm.DefaultHeader.WithFullWidth().Println(string(body))

	return nil
}

// getNearbyCitiesbyName provides capability of get nearby cities.
func (c *Client) getNearbyCitiesbyName() error {
	departure, err := inputWithResult("Departure  city*")
	if err != nil {
		return err
	}

	distanceTo, err := inputWithResult("Distance to (in km)*")
	if err != nil {
		return err
	}

	req := c.Agent.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.Header.Add(authenticationHeaderKey, prefixToken+c.Token)
	req.SetRequestURI(c.Host + urlFindByName)
	req.URI().QueryArgs().Add("departure", departure)
	req.URI().QueryArgs().Add("distanceTo", distanceTo)

	if err := c.Agent.Parse(); err != nil {
		return err
	}

	code, body, errs := c.Agent.Bytes()
	if len(errs) > 0 {
		return fmt.Errorf("%v", errs)
	}

	if code != 200 {
		return fmt.Errorf(string(body))
	}

	pterm.DefaultHeader.WithFullWidth().Println(string(body))

	return nil
}

// getNearbyCitiesbyCoord provides capability of get nearby cities by coordinates.
func (c *Client) getNearbyCitiesbyCoord() error {
	lat, err := inputWithResult("Latitude, decimal number*")
	if err != nil {
		return err
	}

	lon, err := inputWithResult("Longitude, decimal number*")
	if err != nil {
		return err
	}

	distanceTo, err := inputWithResult("Distance to (in km)*")
	if err != nil {
		return err
	}

	req := c.Agent.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.Header.Add(authenticationHeaderKey, prefixToken+c.Token)
	req.SetRequestURI(c.Host + urlFindByCoord)
	req.URI().QueryArgs().Add("lat", lat)
	req.URI().QueryArgs().Add("lon", lon)
	req.URI().QueryArgs().Add("distanceTo", distanceTo)

	if err := c.Agent.Parse(); err != nil {
		return err
	}

	code, body, errs := c.Agent.Bytes()
	if len(errs) > 0 {
		return fmt.Errorf("%v", errs)
	}

	if code != 200 {
		return fmt.Errorf(string(body))
	}

	pterm.DefaultHeader.WithFullWidth().Println(string(body))

	return nil
}

// getInfoCityByName provides capability of get information about city as a short text.
func (c *Client) getInfoCityByName() error {
	name, err := inputWithResult("City name*")
	if err != nil {
		return err
	}

	req := c.Agent.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.Header.Add(authenticationHeaderKey, prefixToken+c.Token)
	req.SetRequestURI(c.Host + urlGetInfo)
	req.URI().QueryArgs().Add("name", name)

	pterm.Info.Println("please wait")

	if err := c.Agent.Parse(); err != nil {
		return err
	}

	code, body, errs := c.Agent.Bytes()
	if len(errs) > 0 {
		return fmt.Errorf("%v", errs)
	}

	if code != 200 {
		return fmt.Errorf(string(body))
	}

	var response struct {
		Text string `json:"text"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	pterm.DefaultHeader.WithFullWidth().Println(response.Text)

	return nil
}
