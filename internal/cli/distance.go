package cli

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
)

const (
	// commands
	commandDistance      = "distance"
	commandDistanceShort = "dist"
	// urls
	urlDistance    = "/v1/user/distance"
	urlDistanceApi = "/v1/api/distance"
)

var (
	From  string // departure city
	Where string // destination city
)

var distance = &cobra.Command{
	Use:     commandDistance,
	Aliases: []string{commandDistanceShort},
	Short:   "calculate distance between two cities",
	Run: func(cmd *cobra.Command, args []string) {
		msg, err := calcDistance()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(msg)
	},
}

func init() {
	distance.Flags().StringVarP(&From, "from", "f", "", "departure city")
	distance.Flags().StringVarP(&Where, "where", "w", "", "destination city")
	distance.Flags().StringVarP(&Token, "token", "t", "", "access token")
	distance.Flags().BoolVarP(&JSONResponse, "json", "j", false, "make a JSON response")
	rootCmd.AddCommand(distance)
}

// calcDistance provides capability of calculate distance between two cities.
func calcDistance() (string, error) {
	agent := fiber.AcquireAgent()
	token := cliToken
	if Token != "" {
		token = Token
	}

	url := urlDistance
	if JSONResponse {
		url = urlDistanceApi
	}

	req := agent.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.Header.Add(authenticationHeaderKey, prefixToken+token)
	req.SetRequestURI(cliHost + url)
	req.URI().QueryArgs().Add("departure", From)
	req.URI().QueryArgs().Add("destination", Where)

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
