package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alaleks/geospace/internal/server/database/models"
	"github.com/alaleks/geospace/pkg/distance"
	"github.com/gofiber/fiber/v2"
)

// CalculateDistance performs a distance between two cities.
func (h *Hdls) CalculateDistance(c *fiber.Ctx) error {
	var (
		departure   = c.Query("departure")
		destination = c.Query("destination")
	)

	if strings.TrimSpace(departure) == "" {
		err := fmt.Errorf("departure %v", ErrEmptyParam)
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	if strings.TrimSpace(destination) == "" {
		err := fmt.Errorf("destination %v", ErrEmptyParam)
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	var (
		cityDeparture     models.City
		cityDestination   models.City
		distStraight      int
		chErr             = make(chan error, 1)
		cityDepartureCh   = make(chan models.City, 1)
		cityDestinationCh = make(chan models.City, 1)
	)

	// run concurrently findind
	go h.db.FindCityConc(departure, chErr, cityDepartureCh)
	go h.db.FindCityConc(destination, chErr, cityDestinationCh)

	for {
		select {
		case err := <-chErr:
			err = fmt.Errorf("error find city in db: %v", err)
			return h.errorApiRequest(c, fiber.StatusBadRequest, err)
		case cityDeparture = <-cityDepartureCh:
			continue
		case cityDestination = <-cityDestinationCh:
			continue
		default:
			if cityDeparture == (models.City{}) || cityDestination == (models.City{}) {
				continue
			}

			distStraight = int(distance.CalcGreatCircle(
				cityDeparture.Latitude, cityDeparture.Longitude,
				cityDestination.Latitude, cityDestination.Longitude))

			response := fmt.Sprintf("distance between %s, %s and %s, %s by straight line %d km",
				cityDeparture.Name, cityDeparture.Country,
				cityDestination.Name, cityDestination.Country, distStraight)

			distanceRoad, err := h.getDistancebyRoad(cityDeparture.Longitude, cityDeparture.Latitude,
				cityDestination.Longitude, cityDestination.Latitude)
			if err == nil || distanceRoad != 0 {
				response += fmt.Sprintf(" / by road %d km", distanceRoad)
			}

			return c.SendString(response)
		}
	}
}

// FindObjectsNearByName performs search for all objects at a distance
// until n km from the object passed in the query.
func (h *Hdls) FindObjectsNearByName(c *fiber.Ctx) error {
	var (
		departure  = c.Query("departure")
		distanceTo = c.Query("distanceTo")
	)

	if strings.TrimSpace(departure) == "" {
		err := fmt.Errorf("departure %v", ErrEmptyParam)
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	if strings.TrimSpace(distanceTo) == "" {
		err := fmt.Errorf("distanceTo %v", ErrEmptyParam)
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}
	// convert distance to int
	dist, err := strconv.Atoi(distanceTo)
	if err != nil {
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	ciyDeparture, cities, err := h.db.FindObjectsNearByName(departure, dist)
	if err != nil {
		err = fmt.Errorf("error finding cities nearby in database: %v", err)
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	respCities := make([]string, 0, len(cities))
	for _, city := range cities {
		dist := distance.CalcGreatCircle(ciyDeparture.Latitude, ciyDeparture.Longitude,
			city.Latitude, city.Longitude)
		respCities = append(respCities, fmt.Sprintf("%s, %s (%d km)", city.Name, city.Country, int(dist)))
	}

	response := fmt.Sprintf("There are %d cities at a distance %d km from %s, %s\n",
		len(respCities), dist, ciyDeparture.Name, ciyDeparture.Country)

	response += fmt.Sprintf("List:\n %s", strings.Join(respCities, ", "))

	return c.SendString(response)
}

// FindObjectsNearByCoord performs search for all objects at a distance
// until n km from coordinates passed in the query.
func (h *Hdls) FindObjectsNearByCoord(c *fiber.Ctx) error {
	var (
		latStr     = c.Query("lat")
		lonStr     = c.Query("lon")
		distanceTo = c.Query("distanceTo")
	)

	if strings.TrimSpace(distanceTo) == "" {
		err := fmt.Errorf("distanceTo %v", ErrEmptyParam)
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	if strings.TrimSpace(latStr) == "" {
		err := fmt.Errorf("lat %v", ErrEmptyParam)
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	if strings.TrimSpace(lonStr) == "" {
		err := fmt.Errorf("lon %v", ErrEmptyParam)
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	dist, err := strconv.Atoi(distanceTo)
	if err != nil {
		err = fmt.Errorf("error convert distance to int: %v", err)
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		err = fmt.Errorf("error convert lat to decimal number: %v", err)
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		err = fmt.Errorf("error convert lon to decimal number: %v", err)
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	cities, err := h.db.FindObjectsNearByCoord(lat, lon, dist)
	if err != nil {
		err = fmt.Errorf("error finding cities nearby in database: %v", err)
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	respCities := make([]string, 0, len(cities))
	for _, city := range cities {
		dist := distance.CalcGreatCircle(lat, lon, city.Latitude, city.Longitude)
		respCities = append(respCities, fmt.Sprintf("%s, %s (%d km)", city.Name, city.Country, int(dist)))
	}

	response := fmt.Sprintf("There are %d cities at a distance %d km\n", len(respCities), dist)

	if len(respCities) > 0 {
		response += fmt.Sprintf("List:\n%s", strings.Join(respCities, ","))
	}

	return c.SendString(response)
}
