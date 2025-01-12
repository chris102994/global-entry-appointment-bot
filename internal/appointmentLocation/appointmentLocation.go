package appointmentLocation

import (
	"crypto/tls"
	"encoding/json"
	"github.com/chris102994/global-entry-appointment-bot/internal/appointment"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"time"
)

type AppointmentLocation struct {
	ID                  int                    `json:"id"`
	Name                string                 `json:"name"`
	ShortName           string                 `json:"shortName"`
	LocationType        string                 `json:"locationType"`
	LocationCode        string                 `json:"locationCode"`
	Address             string                 `json:"address"`
	AddressAdditional   string                 `json:"addressAdditional"`
	City                string                 `json:"city"`
	State               string                 `json:"state"`
	PostalCode          string                 `json:"postalCode"`
	CountryCode         string                 `json:"countryCode"`
	TZData              string                 `json:"tzData"`
	PhoneNumber         string                 `json:"phoneNumber"`
	PhoneAreaCode       string                 `json:"phoneAreaCode"`
	PhoneCountryCode    string                 `json:"phoneCountryCode"`
	PhoneExtension      string                 `json:"phoneExtension"`
	PhoneAltNumber      string                 `json:"phoneAltNumber"`
	PhoneAltAreaCode    string                 `json:"phoneAltAreaCode"`
	PhoneAltCountryCode string                 `json:"phoneAltCountryCode"`
	PhoneAltExtension   string                 `json:"phoneAltExtension"`
	FaxNumber           string                 `json:"faxNumber"`
	FaxAreaCode         string                 `json:"faxAreaCode"`
	FaxCountryCode      string                 `json:"faxCountryCode"`
	FaxExtension        string                 `json:"faxExtension"`
	EffectiveDate       appointment.CustomTime `json:"effectiveDate"`
	Temporary           bool                   `json:"temporary"`
	InviteOnly          bool                   `json:"inviteOnly"`
	Operational         bool                   `json:"operational"`
	Directions          string                 `json:"directions"`
	Notes               string                 `json:"notes"`
	MapFileName         string                 `json:"mapFileName"`
	RemoteInd           bool                   `json:"remoteInd"`
	Services            []Service              `json:"services"`
}

// Service represents the services offered at a location
type Service struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetAppointmentLocations(state string) ([]*AppointmentLocation, error) {
	var appointmentLocations []*AppointmentLocation

	baseUrl := "https://ttp.cbp.dhs.gov/schedulerapi/locations/"

	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	client := &http.Client{
		Transport: customTransport,
		Timeout:   10 * time.Second,
	}

	requestParams := url.Values{
		"invitesOnly": {"false"},
		"operational": {"true"},
		"serviceName": {"Global Entry"},
	}

	if state != "" {
		requestParams.Add("state", state)
	}

	requestUrl := baseUrl + "?" + requestParams.Encode()

	log.WithFields(log.Fields{
		"requestUrl": requestUrl,
	}).Trace("Scraping endpoint")

	resp, err := client.Get(requestUrl)
	if err != nil {
		return appointmentLocations, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return appointmentLocations, err
	}

	err = json.Unmarshal(bodyBytes, &appointmentLocations)
	if err != nil {
		return appointmentLocations, err
	}

	return appointmentLocations, nil
}
