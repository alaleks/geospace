package handlers

import (
	"strconv"
	"strings"

	"github.com/alaleks/geospace/internal/server/database/models"
	"github.com/alaleks/geospace/pkg/distance"
	"github.com/gofiber/fiber/v2"
)

// RespCity represents a data for response list near a citу,
// supplementing the data with the distance field.
type RespCity struct {
	models.City
	Distance int `json:"distance"`
}

// CalculateDistanceAPI performs a distance between two cities.
func (h *Hdls) CalculateDistanceAPI(c *fiber.Ctx) error {
	var (
		departure   = c.Query("departure")
		destination = c.Query("destination")
	)

	if strings.TrimSpace(departure) == "" ||
		strings.TrimSpace(destination) == "" {
		return h.errorApiRequest(c, fiber.StatusBadRequest,
			ErrEmptyDataForCalculateDisnance)
	}

	var (
		response struct {
			Departure        models.City `json:"departure"`
			Destination      models.City `json:"destination"`
			DistanceStraight int         `json:"distance_straight"`
			DistanceRoad     int         `json:"distance_road,omitempty"`
		}
		chErr             = make(chan error, 1)
		cityDepartureCh   = make(chan models.City, 1)
		cityDestinationCh = make(chan models.City, 1)
	)

	// run concurrently findind
	go h.db.FindCityConc(departure, chErr, cityDepartureCh)
	go h.db.FindCityConc(destination, chErr, cityDestinationCh)

	// in cycle for will wait when search is done
	// at the first error, immediately return the request error
	for {
		select {
		case err := <-chErr:
			return h.errorApiRequest(c, fiber.StatusBadRequest, err)
		case response.Departure = <-cityDepartureCh:
			continue
		case response.Destination = <-cityDestinationCh:
			continue
		default:
			if response.Departure == (models.City{}) || response.Destination == (models.City{}) {
				continue
			}

			response.DistanceStraight = int(distance.CalcGreatCircle(
				response.Departure.Latitude, response.Departure.Longitude,
				response.Destination.Latitude, response.Destination.Longitude))

			response.DistanceRoad, _ = h.getDistancebyRoad(
				response.Departure.Longitude, response.Departure.Latitude,
				response.Destination.Longitude, response.Destination.Latitude)

			return c.JSON(response)
		}
	}
}

// FindObjectsNearAPI performs search for all objects at a distance
// until n km from the object passed in the query.
func (h *Hdls) FindObjectsNearByNameAPI(c *fiber.Ctx) error {
	var (
		departure  = c.Query("departure")
		distanceTo = c.Query("distance-to")
	)

	if strings.TrimSpace(departure) == "" {
		return h.errorApiRequest(c, fiber.StatusBadRequest,
			ErrEmptyDataForCalculateDisnance)
	}

	dist, err := strconv.Atoi(distanceTo)
	if err != nil {
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	ciyDeparture, cities, err := h.db.FindObjectsNearByName(departure, dist)
	if err != nil {
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	var respCities = make([]RespCity, 0, len(cities))
	for _, city := range cities {
		respCities = append(respCities, RespCity{
			city,
			int(distance.CalcGreatCircle(ciyDeparture.Latitude, ciyDeparture.Longitude,
				city.Latitude, city.Longitude)),
		})
	}

	response := struct {
		Departure    models.City `json:"departure"`
		DistanceTo   int         `json:"distance_to"`
		QtyNearby    int         `json:"qty_nearby"`
		CitiesNearby []RespCity  `json:"cities_nearby"`
	}{
		Departure:    ciyDeparture,
		DistanceTo:   dist,
		QtyNearby:    len(respCities),
		CitiesNearby: respCities,
	}

	return c.JSON(response)
}

func (h *Hdls) FindObjectsNearByCoordAPI(c *fiber.Ctx) error {
	var (
		latStr     = c.Query("lat")
		lonStr     = c.Query("lon")
		distanceTo = c.Query("distance-to")
	)

	if strings.TrimSpace(distanceTo) == "" ||
		strings.TrimSpace(latStr) == "" ||
		strings.TrimSpace(lonStr) == "" {
		return h.errorApiRequest(c, fiber.StatusBadRequest,
			ErrEmptyDataForCalculateDisnance)
	}

	dist, err := strconv.Atoi(distanceTo)
	if err != nil {
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	cities, err := h.db.FindObjectsNearByCoord(lat, lon, dist)
	if err != nil {
		return h.errorApiRequest(c, fiber.StatusBadRequest, err)
	}

	var respCities = make([]RespCity, 0, len(cities))
	for _, city := range cities {
		respCities = append(respCities, RespCity{
			city,
			int(distance.CalcGreatCircle(lat, lon, city.Latitude, city.Longitude)),
		})
	}

	response := struct {
		DistanceTo   int        `json:"distance_to"`
		QtyNearby    int        `json:"qty_nearby"`
		CitiesNearby []RespCity `json:"cities_nearby"`
	}{
		DistanceTo:   dist,
		QtyNearby:    len(respCities),
		CitiesNearby: respCities,
	}

	return c.JSON(response)
}

// errorApiRequest performs send status code and message error.
func (h *Hdls) errorApiRequest(c *fiber.Ctx, code int, err error) error {
	var errReq = struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{
		Code:    code,
		Message: err.Error(),
	}

	return c.Status(code).JSON(errReq)
}
