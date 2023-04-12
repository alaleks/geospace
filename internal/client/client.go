// Package client provides a command line client for interactions server.
package client

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

const (
	contentTypeJSON                 = "application/json"
	commandDistance                 = "distance"
	commandCalculate                = "calculate"
	commandCities                   = "cities"
	commandCitiesNearbybyName       = "nearby by Name"
	commandCitiesNearbybyCoordinate = "nearby by Coordinate"
	commandBack                     = "back"
	commandExit                     = "exit"
)

// typical errors
var (
	ErrServerInternal = errors.New("server internal error")
)

// Client contains agent instance and technical information.
type Client struct {
	Agent   *fiber.Agent
	Version string
	Host    string
	Name    string
	Token   string
}

// New returns a new pointer instance a app of client.
func New(version, host, name string) *Client {
	client := Client{
		Agent:   fiber.AcquireAgent(),
		Version: version,
		Name:    name,
	}

	if !strings.HasPrefix(host, "localhost") {
		client.Host = "http://localhost" + host
	} else if !strings.HasPrefix(host, "http") {
		client.Host = "http://" + host
	}

	return &client
}

// Run performs running client in command line.
func (c *Client) Run() {
	// check server health
	if err := c.checkServer(); err != nil {
		printErrWithExit(ErrServerInternal)
	}

	// print greeting
	c.greet()

	// new user login or registration
	for {
		err := c.authentication()
		if err != nil {
			printErr(err)
			printExitMessage()

			continue
		}

		break
	}

	printExitMessage()

	// create commands
	commandsMainMenu := [...]string{
		commandDistance,
		commandCities,
		commandExit,
	}

	commandsMenuDistance := [...]string{
		commandCalculate,
		commandBack,
		commandExit,
	}

	commandsMenuCities := [...]string{
		commandCitiesNearbybyName,
		commandCitiesNearbybyCoordinate,
		commandBack,
		commandExit,
	}

mainMenu:
	for {
		printer := pterm.DefaultInteractiveSelect.WithOptions(commandsMainMenu[:])
		selectedOptions, err := printer.Show()
		if err != nil {
			printErr(err)

			break mainMenu
		}

		switch selectedOptions {
		case commandDistance:
		distance:
			for {
				printer = pterm.DefaultInteractiveSelect.WithOptions(commandsMenuDistance[:])
				selectedOptions, err = printer.Show()
				if err != nil {
					printErr(err)

					break distance
				}

				switch selectedOptions {
				case commandCalculate:
					err = c.calcDistance()
					if err != nil {
						printErr(err)
					}
				case commandBack:
					pterm.Info.Println("return to the main menu")

					break distance
				case commandExit:
					pterm.Info.Println("client closed")

					break mainMenu
				default:
					printErrWithExit(ErrInvalidCommand)

					break mainMenu
				}
			}
		case commandCities:
		cities:
			for {
				printer = pterm.DefaultInteractiveSelect.WithOptions(commandsMenuCities[:])
				selectedOptions, err = printer.Show()
				if err != nil {
					printErr(err)

					break cities
				}

				switch selectedOptions {
				case commandCitiesNearbybyName:
					err = c.getNearbyCitiesbyName()
					if err != nil {
						printErr(err)
					}

					continue
				case commandCitiesNearbybyCoordinate:
					err = c.getNearbyCitiesbyCoord()
					if err != nil {
						printErr(err)
					}

					continue
				case commandBack:
					pterm.Info.Println("return to the main menu")

					break cities
				case commandExit:
					pterm.Info.Println("client closed")

					break mainMenu
				default:
					printErrWithExit(ErrInvalidCommand)

					break mainMenu
				}
			}
		case commandExit:
			pterm.Info.Println("client closed")

			break mainMenu
		default:
			printErrWithExit(ErrInvalidCommand)
		}
	}
}

// checkServer performs checking work of server using url /ping.
func (c *Client) checkServer() error {
	req := c.Agent.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.SetRequestURI(c.Host + "/ping")

	if err := c.Agent.Parse(); err != nil {
		return err
	}

	code, _, _ := c.Agent.Bytes()
	if code != 200 {
		return ErrServerInternal
	}

	printExitMessage()

	return nil
}

// greet prints the greeting message to the terminal.
func (c *Client) greet() {
	fmt.Printf("\n\n")
	s, err := pterm.DefaultBigText.
		WithLetters(putils.LettersFromString(c.Name)).Srender()
	if err != nil {
		return
	}

	pterm.DefaultCenter.Println(s)
	pterm.DefaultCenter.WithCenterEachLineSeparately().Println(c.Version)
}

// printExitMessage displays message how to exit the client.
func printExitMessage() {
	pterm.DefaultCenter.WithCenterEachLineSeparately().
		Println("press ctrl + c to exit")
}

// printErr displays the error.
func printErr(err error) {
	pterm.Error.Println(err.Error())
}

// printErr displays the error and os exit.
func printErrWithExit(err error) {
	pterm.Error.Println(err.Error())
	os.Exit(1)
}
