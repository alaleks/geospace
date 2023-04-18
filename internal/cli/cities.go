package cli

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
)

const (
	//commands
	commandCitiesNearbybyName            = "nearbyname"
	commandCitiesNearbybyNameShort       = "nn"
	commandCitiesNearbybyCoordinate      = "nearbycoordinate"
	commandCitiesNearbybyCoordinateShort = "nc"
	//urls
	urlFindByName     = "/v1/user/find-by-name"
	urlFindByCoord    = "/v1/user/find-by-coord"
	urlFindByNameAPI  = "/v1/api/find-by-name"
	urlFindByCoordAPI = "/v1/api/find-by-coord"
)

var (
	DistanceTo int64
	Lat        float64
	Lon        float64
)

// find nearby cities from the current city by name.
var findnn = &cobra.Command{
	Use:     commandCitiesNearbybyName,
	Aliases: []string{commandCitiesNearbybyNameShort},
	Short:   "find nearby cities from the current city by name.",
	Run: func(cmd *cobra.Command, args []string) {
		msg, err := findNearbybyName()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(msg)
	},
}

// find nearby cities from the current coordinates point.
var findnc = &cobra.Command{
	Use:     commandCitiesNearbybyCoordinate,
	Aliases: []string{commandCitiesNearbybyCoordinateShort},
	Short:   "find nearby cities from the current coordinates point.",
	Run: func(cmd *cobra.Command, args []string) {
		msg, err := findNearbybyCoordinate()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(msg)
	},
}

func init() {
	// find nearby cities from the current city by name.
	findnn.Flags().StringVarP(&From, "from", "f", "", "departure city")
	findnn.Flags().Int64VarP(&DistanceTo, "distanceto", "d", 0, "distance to")
	findnn.Flags().BoolVarP(&JSONResponse, "json", "j", false, "make a JSON response")
	rootCmd.AddCommand(findnn)

	// find nearby cities from the current coordinates point.
	findnc.Flags().Float64VarP(&Lat, "lat", "a", 0, "latitude")
	findnc.Flags().Float64VarP(&Lon, "lon", "o", 0, "longitude")
	findnc.Flags().Int64VarP(&DistanceTo, "distanceto", "d", 0, "distance to")
	findnc.Flags().BoolVarP(&JSONResponse, "json", "j", false, "make a JSON response")
	rootCmd.AddCommand(findnc)
}

// findNearbybyName provides capability of get nearby cities by name from current.
func findNearbybyName() (string, error) {
	agent := fiber.AcquireAgent()
	token := cliToken
	if Token != "" {
		token = Token
	}

	url := urlFindByName
	if JSONResponse {
		url = urlFindByNameAPI
	}

	req := agent.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.Header.Add(authenticationHeaderKey, prefixToken+token)
	req.SetRequestURI(cliHost + url)
	req.URI().QueryArgs().Add("departure", From)
	req.URI().QueryArgs().Add("distanceTo", strconv.Itoa(int(DistanceTo)))

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

	return string(body), nil
}

// findNearbybyCoordinate provides capability of get nearby cities by coordinate from current point.
func findNearbybyCoordinate() (string, error) {
	agent := fiber.AcquireAgent()
	token := cliToken
	if Token != "" {
		token = Token
	}

	url := urlFindByCoord
	if JSONResponse {
		url = urlFindByCoordAPI
	}

	req := agent.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.Header.SetMethod(fiber.MethodGet)
	req.Header.Add(authenticationHeaderKey, prefixToken+token)
	req.SetRequestURI(cliHost + url)
	req.URI().QueryArgs().Add("lat", fmt.Sprintf("%v", Lat))
	req.URI().QueryArgs().Add("lon", fmt.Sprintf("%v", Lon))
	req.URI().QueryArgs().Add("distanceTo", strconv.Itoa(int(DistanceTo)))

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

	return string(body), nil
}
