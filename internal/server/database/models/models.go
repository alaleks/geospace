// Package models contains data structures for current app.
package models

type (
	City struct {
		Name             string  `db:"name"`              // Name of the city
		NameASCII        string  `db:"name_ascii"`        // Name of the city ASCII
		AlternativeNames string  `db:"alternative_names"` // Alternative names of the city
		CountryCode      string  `db:"country_code"`      // Code of country with this city located
		Country          string  `db:"country"`           // Name of the country with this city located
		Timezone         string  `db:"timezone"`          // Name of the timezone with this city located
		CreatedAt        int64   `db:"created_at"`        // Date when the city was created formated by Unix timestamp
		ID               int     `db:"cid"`               // ID of the city (inside application)
		Latitude         float64 `db:"latitude"`          // Latitude of the city
		Longitude        float64 `db:"longitude"`         // Longitude of the city
	}

	User struct {
		Name      string `db:"name"`       // Name of the user
		Email     string `db:"email"`      // Email of the user
		Password  string `db:"password"`   // Password of the user
		UID       int    `db:"uid"`        // ID of the user
		CreatedAt int64  `db:"created_at"` // Date when the user was created
	}
)
