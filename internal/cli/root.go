// cli package implements CLI client features.
package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

const (
	contentTypeJSON         = "application/json"
	prefixToken             = "Bearer "
	authenticationHeaderKey = "Authorization"
)

// application information
var (
	cliVersion   = "v 1.0"
	cliHost      = "http://localhost:3000"
	cliName      = "geospace"
	cliToken     string // access token from ldf
	JSONResponse bool   // make json response
	Token        string // access token
)

// UpdateAppInfo performs an update of the application information.
func UpdateAppInfo(version, host, name, token string) {
	if strings.TrimSpace(version) != "" {
		cliVersion = version
	}

	if strings.TrimSpace(name) != "" {
		cliName = name
	}

	if strings.TrimSpace(token) != "" {
		cliToken = token
	}

	if strings.TrimSpace(host) != "" {
		if !strings.HasPrefix(host, "localhost") {
			cliHost = "http://localhost" + host
		} else if !strings.HasPrefix(host, "http") {
			cliHost = "http://" + host
		}
	}
}

// rootCmd represents the root cli command
var rootCmd = &cobra.Command{
	Use:     cliName,
	Version: cliVersion,
	Short:   cliName + " simple CLI client",
	Long: cliName + `service with CLI client for calculating distances
	 between cities and much more...`,
	Run: func(cmd *cobra.Command, args []string) {},
}

// Execute start the CLI client.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(fmt.Errorf("there was an error while executing your CLI '%s'", err))
	}
}

// prettyJSON performs pretty-printing json output.
func prettyJSON(jsonData []byte) string {
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, jsonData, "", "  ")
	if error != nil {
		return ""
	}

	return prettyJSON.String()
}
