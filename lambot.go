package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/esoterick/lambot/log"
	"github.com/esoterick/lambot/transmission"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

var tor *transmission.Tor

func init() {
	viper.SetConfigName("config")                // name of config file (without extension)
	viper.AddConfigPath("/etc/lambot/")          // path to look for the config file in
	viper.AddConfigPath("$HOME/.config/lambot/") // call multiple times to add many search paths
	viper.AddConfigPath(".")                     // optionally look for config in the working directory

	viper.SetEnvPrefix("lambot")
	viper.BindEnv("token")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatalf("fatal error config file: %s\n", err)
	}
}

func main() {
	token := viper.GetString("token")
	log.Info("token: %s", token)

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("error creating discord session %s\n", err)
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Fatalf("error opening connection %s\n", err)
	}

	// Connecting to Transmission
	log.Info("Connecting to Transmission...")

	proto := viper.GetString("transmission.proto")
	host := viper.GetString("transmission.host")
	port := viper.GetInt("transmission.port")
	path := viper.GetString("transmission.path")

	url := fmt.Sprintf("%s://%s:%d%s", proto, host, port, path)

	username := viper.GetString("transmission.username")
	password := viper.GetString("transmission.password")

	tor, err = transmission.NewTor(url, username, password)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Bot has loaded...")
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// If the message is "ping" reply with "Pong!"
	if m.Content == "?ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "?pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}

	if m.Content == "?torrents get" {
		ts, err := tor.GetTorrents()
		if err != nil {
			log.Fatal(err)
		}

		msg := "```\n"
		for _, t := range ts {
			msg = msg + fmt.Sprintf("%4d. [%dx%d] {%d} %s\n", t.ID, t.RateDownload, t.RateUpload, t.Status, t.Name)
		}
		msg = msg + "```"

		s.ChannelMessageSend(m.ChannelID, msg)
	}
}
