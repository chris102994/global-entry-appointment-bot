package run

import (
	"github.com/chris102994/global-entry-appointment-bot/internal/appointment"
	"github.com/chris102994/global-entry-appointment-bot/internal/config"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var inputConfig *config.Config

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the global-entry-appointment-bot",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.WithFields(log.Fields{
			"runConfig": inputConfig.Run,
		}).Trace("Running global-entry-appointment-bot run")

		// If cron expression is provided, run on a schedule
		if inputConfig.Cron.Expression != "" {
			log.WithFields(log.Fields{
				"cron": inputConfig.Cron,
			}).Info("Cron mode")

			logger := log.StandardLogger()
			c := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(logger)))
			_, err := c.AddFunc(inputConfig.Cron.Expression, func() {
				// Use the new helper function for repeated logic
				err := handleAvailableAppointments()
				if err != nil {
					cobra.CheckErr(err)
				}
				log.Trace("Cron Schedule ran")
			})
			cobra.CheckErr(err)

			c.Run()

		} else {
			// Run once if no cron expression
			log.WithFields(log.Fields{}).Trace("No Cron Schedule detected.")
			err := handleAvailableAppointments()
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	viper.SetDefault("Run.Limit", 1)
	viper.SetDefault("Run.Minimum", 1)
	viper.SetDefault("Run.OrderBy", "soonest")

	runCmd.Flags().StringP("cron-expression", "c", "", "The cron schedule to run the command")
	runCmd.Flags().StringP("limit", "n", "1", "The number of results returned by the query (useful when trying to narrow down multiple results)")
	runCmd.Flags().StringArrayP("location-ids", "i", nil, "The Location IDs where you'd like your appointment (hint: use the lookup command to find your ID)")
	runCmd.Flags().StringP("minimum", "m", "1", "The number of minimum available appointments")
	runCmd.Flags().StringP("order-by", "b", "soonest", "How to order the results")
	runCmd.Flags().StringP("discord-token", "", "", "Discord Notification Token")
	runCmd.Flags().StringP("discord-user-id", "", "", "Discord User ID")

	runCmd.MarkFlagsRequiredTogether("discord-token", "discord-user-id")

	viper.BindPFlag("Cron.Expression", runCmd.Flag("cron-expression"))
	viper.BindPFlag("Run.Limit", runCmd.Flag("limit"))
	viper.BindPFlag("Run.LocationIds", runCmd.Flag("location-ids"))
	viper.BindPFlag("Run.Minimum", runCmd.Flag("minimum"))
	viper.BindPFlag("Run.OrderBy", runCmd.Flag("order-by"))
	viper.BindPFlag("Notify.Discord.Token", runCmd.Flag("discord-token"))
	viper.BindPFlag("Notify.Discord.UserID", runCmd.Flag("discord-user-id"))
}

func NewRunCmd(c *config.Config) *cobra.Command {
	inputConfig = c
	return runCmd
}

func handleAvailableAppointments() error {
	availableAppointmentSlots, err := appointment.GetAppointments(*inputConfig.Run)
	if err != nil {
		return err
	}

	if len(availableAppointmentSlots) == 0 {
		log.Println("No Appointments Available")
	} else {
		for _, appointmentSlot := range availableAppointmentSlots {
			log.WithFields(log.Fields{
				"Appointment": appointmentSlot,
			}).Println("Appointment Found")
			if inputConfig.Notify != nil {
				err = inputConfig.Notify.Run(appointmentSlot)
				cobra.CheckErr(err)
			}
		}
	}

	return nil
}
