package main

/*
 *
 * 	Hourly vro bot by SP649/SPdoesCode
 *
 */

import (
	"fmt"
	"time"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
)

type CFG struct {
	Token     string   `toml:"token"`
	Prefix    string   `toml:"prefix"`
	ImageDir  string   `toml:"image_dir"`
	Images    []string `toml:"images"`
	Channel   string   `toml:"channel"`
}

var config CFG


func getConfig() {

	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		fmt.Println("ERROR: Couldnt load config.toml: ", err)
		os.Exit(1)
	}
	fmt.Println("Parsed toml")
	fmt.Println("Bot prefix is: ", config.Prefix)
}

func sendRand(self *discordgo.Session, Channel string) {
	rand.Seed(time.Now().UnixNano())
	if len(config.Images) == 0 {
		self.ChannelMessageSend(Channel, "No images configured.")
		fmt.Println("WARN: No images in config.Images")
		return
	}
	num := rand.Intn(len(config.Images))
	filename := config.ImageDir+config.Images[num] // get the random file
	fmt.Println("Selected image: ", filename)
	file, err := os.Open(filename)
	if err != nil { self.ChannelMessageSend(Channel, "Couldnt find image, sorry"); fmt.Println("WARN: Failed to find image ", filename,": ", err); return } // warn and alert to the error
	defer file.Close() // we need to close
	_, err = self.ChannelMessageSendComplex(Channel, &discordgo.MessageSend{
		Files: []*discordgo.File{
			{
				Name:   filename,
				Reader: file,
			},
		},
	})
	fmt.Println("Send image: ", file, filename)
	if err != nil {
		self.ChannelMessageSend(Channel, "Couldnt open image...")
		fmt.Println("WARN: Failed to find image ", filename,": ", err)
	}

}

// for commands only
func ctrlMessages(self *discordgo.Session, message *discordgo.MessageCreate) {
	fmt.Println("Message received: ", message.Content)

	if message.Author.ID == self.State.User.ID {
		return // return nothing cuz its our own message
	}

	if strings.HasPrefix(message.Content, config.Prefix) { // check for the prefix
		endcmd := strings.TrimPrefix(message.Content, config.Prefix) // remove it
		fmt.Println("Got command and stripped prefix: ", config.Prefix)
		switch endcmd { // switch it to see what command it is
			case "getvro":
				fmt.Println("Command is getvro, sending sendRand() to channel: ", message.ChannelID)
				sendRand(self, message.ChannelID)
		}
	}
}

// every hour message hanndler
func hourlyMessage(self *discordgo.Session) {
	go func(){
		ticker := time.NewTicker(1 * time.Hour) // set up the hour control
		defer ticker.Stop()

		for {
			<- ticker.C
			sendRand(self, config.Channel)

		}
	}()
}

func main() {

	getConfig()

	bot, err := discordgo.New("Bot "+config.Token) // grab the token and make a new session
	if err != nil {
		fmt.Println("ERROR: Couldnt create discord session: ", err)
		os.Exit(1)
	}

	bot.AddHandler(ctrlMessages) // grab messages into the message control

	err2 := bot.Open() // try to open the connection
	if err2 != nil {
		fmt.Println("ERROR: Coulnt open the connection: ", err2)
		os.Exit(1)
	}

	bot.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	defer bot.Close() // close the connection when done

	fmt.Println("Bot should be running properly now, If you need to stop do ctrl + c or ctrl + d!")

	quit := make(chan os.Signal, 1)

	hourlyMessage(bot) // start the hourly message

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<- quit // stop when indicated

	fmt.Println("Quiting...")

	os.Exit(0)
}
