package checker

import "time"

type GeoIPAPIResponse struct {
	IP      string `json:"ip"`
	Type    string `json:"type"`
	Country struct {
		IsEuMember   bool        `json:"is_eu_member"`
		CurrencyCode string      `json:"currency_code"`
		Continent    string      `json:"continent"`
		Name         string      `json:"name"`
		CountryCode  string      `json:"country_code"`
		State        string      `json:"state"`
		City         string      `json:"city"`
		Zip          interface{} `json:"zip"`
		Timezone     string      `json:"timezone"`
	} `json:"country"`
	Location struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	Asn struct {
		Number  int    `json:"number"`
		Name    string `json:"name"`
		Network string `json:"network"`
		Type    string `json:"type"`
	} `json:"asn"`
	Duration time.Duration `json:"-"`
}
