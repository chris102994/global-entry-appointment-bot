package notify

import (
	"github.com/bwmarrin/discordgo"
	"github.com/chris102994/global-entry-appointment-bot/internal/appointment"
	"github.com/chris102994/global-entry-appointment-bot/internal/appointmentLocation"
	log "github.com/sirupsen/logrus"
	"strings"
	"text/template"
)

type Notify struct {
	Discord *Discord `mapstructure:"discord,omitempty"`
}

type Discord struct {
	Token  string `mapstructure:"token,omitempty"`
	UserID string `mapstructure:"userid,omitempty"`
}

func (d *Discord) Run(message string) error {
	log.Info("Sending Discord Message")

	discord, err := discordgo.New("Bot " + d.Token)
	if err != nil {
		return err
	}

	err = discord.Open()
	if err != nil {
		return err
	}
	defer discord.Close()

	channel, err := discord.UserChannelCreate(d.UserID)
	if err != nil {
		return err
	}

	_, err = discord.ChannelMessageSend(channel.ID, message)
	if err != nil {
		return err
	}

	log.Info("Message sent to user successfully!")

	return nil
}

type NotificationData struct {
	Appointment         *appointment.Appointment
	AppointmentLocation *appointmentLocation.AppointmentLocation
}

func (n *Notify) Run(appointment *appointment.Appointment) error {
	log.WithFields(log.Fields{
		"notify":      n,
		"appointment": appointment,
	}).Trace("Sending Notifications")

	appointmentLocations, err := appointmentLocation.GetAppointmentLocations("")
	var theAppointmentLocation appointmentLocation.AppointmentLocation

	for _, _appointmentLocation := range appointmentLocations {
		if _appointmentLocation.ID == appointment.LocationID {
			theAppointmentLocation = *_appointmentLocation
			break
		}
	}

	log.WithFields(log.Fields{
		"theAppointmentLocation": theAppointmentLocation,
	}).Trace("Found the appointment location")

	messageTemplate, err := template.New("notificationMessage").Parse("# Appointment Detected!" +
		"\n* **Date:** {{.Appointment.StartTimestamp.DateFmt}}" +
		"\n* **Time:** {{.Appointment.StartTimestamp.TimeFmt}}" +
		"\n* **Name:** {{.AppointmentLocation.Name}}" +
		"\n* **State:** {{.AppointmentLocation.State}}" +
		"\n* **City:** {{.AppointmentLocation.City}}" +
		"\n* **Address:** {{.AppointmentLocation.AddressAdditional}}")
	if err != nil {
		return err
	}

	var messageBuilder strings.Builder

	notificationData := NotificationData{
		Appointment:         appointment,
		AppointmentLocation: &theAppointmentLocation,
	}

	err = messageTemplate.Execute(&messageBuilder, notificationData)
	if err != nil {
		return err
	}

	message := messageBuilder.String()

	log.WithFields(log.Fields{
		"message": message,
	}).Println("Appointment Information")

	if n.Discord != nil {
		if n.Discord.Token != "" && n.Discord.UserID != "" {
			err := n.Discord.Run(message)
			if err != nil {
				return err
			}
		} else {
			log.WithFields(log.Fields{
				"token":  n.Discord.Token,
				"userID": n.Discord.UserID,
			}).Trace("Unable to send Discord notification. One or both fields must not be blank.")
		}
	}

	return nil
}
