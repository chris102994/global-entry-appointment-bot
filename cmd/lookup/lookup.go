package lookup

import (
	"github.com/chris102994/global-entry-appointment-bot/internal/appointmentLocation"
	"github.com/chris102994/global-entry-appointment-bot/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var inputConfig *config.Config

var lookupCmd = &cobra.Command{
	Use:   "lookup",
	Short: "Lookup location Information for global-entry-appointment-bot",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.WithFields(log.Fields{
			"lookupConfig": inputConfig.Lookup,
		}).Trace("Running global-entry-appointment-bot lookup")

		appointmentLocations, err := appointmentLocation.GetAppointmentLocations(inputConfig.Lookup.State)
		cobra.CheckErr(err)

		if inputConfig.Lookup.City != "" {
			// Normalize the city string (case-insensitive comparison)
			lookupCity := strings.ToLower(inputConfig.Lookup.City)

			// Filter the locations based on city
			var filteredLocations []*appointmentLocation.AppointmentLocation
			for _, appointmentLocation := range appointmentLocations {
				if strings.ToLower(appointmentLocation.City) == lookupCity {
					filteredLocations = append(filteredLocations, appointmentLocation)
				}
			}
			// Use the filtered locations
			appointmentLocations = filteredLocations
		}

		for _, appointmentLocation := range appointmentLocations {
			log.WithFields(log.Fields{
				"ID":      appointmentLocation.ID,
				"Name":    appointmentLocation.Name,
				"State":   appointmentLocation.State,
				"City":    appointmentLocation.City,
				"Address": appointmentLocation.AddressAdditional,
			}).Println("Appointment Location Found")
		}

		return nil
	},
}

func init() {
	lookupCmd.Flags().StringP("city", "c", "", "The city to sort lookups on")
	lookupCmd.Flags().StringP("state", "s", "", "The state to sort lookups on")

	lookupCmd.MarkFlagsOneRequired("city", "state")

	viper.BindPFlag("Lookup.City", lookupCmd.Flag("city"))
	viper.BindPFlag("Lookup.State", lookupCmd.Flag("state"))
}

func NewLookupCmd(c *config.Config) *cobra.Command {
	inputConfig = c
	return lookupCmd
}
