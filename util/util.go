package util

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var Config CFG
var userdata DataFile
var Images []string

func SendRand(self *discordgo.Session, Channel string) {

	if len(Config.Images) == 0 {
		self.ChannelMessageSend(Channel, "No images configured.")
		fmt.Println("WARN: No images in Config.Images")
		return
	}
	num := rand.Intn(len(Images))
	filename := Config.ImageDir + Images[num] // get the random file
	fmt.Println("Selected image: ", filename)
	file, err := os.Open(filename)
	if err != nil {
		self.ChannelMessageSend(Channel, "Couldnt find image, sorry")
		fmt.Println("WARN: Failed to find image ", filename, ": ", err)
		return
	} // warn and alert to the error
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
		fmt.Println("WARN: Failed to find image ", filename, ": ", err)
	}

}

// for commands only
func CtrlMessages(self *discordgo.Session, message *discordgo.MessageCreate) {
	//fmt.Println("Message received: ", message.Content)

	if message.Author.ID == self.State.User.ID {
		return // return nothing cuz its our own message
	}

	if strings.HasPrefix(message.Content, Config.Prefix) { // check for the prefix
		GetConfig()
		fmt.Println("Reloaded Config...")
		args := strings.Fields(message.Content)
		cmd := strings.TrimPrefix(args[0], Config.Prefix) // remove it
		fmt.Println("Got command (", message.Content, ") and stripped prefix ", Config.Prefix)
		switch cmd { // switch it to see what command it is
		case "getvro":
			fmt.Println("Command is getvro with args as ", args, " sent by ", message.Author)
			fmt.Println("Sending sendRand() to channel ", message.ChannelID)
			SendRand(self, message.ChannelID)

		case "sex": // this is a joke command
			fmt.Println("Command is sex with args as ", args, " sent by ", message.Author)
			if rand.Intn(2) == 0 {
				self.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+"> had sex with "+strings.Join(args[1:], " ")+" and made them pregnant!")
				fmt.Println("Sent: \"<@"+message.Author.ID+"> had sex with "+strings.Join(args[1:], " ")+" and made them pregnant!\" to ", message.ChannelID)
			} else {
				self.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+"> had sex with "+strings.Join(args[1:], " ")+" and failed to make them pregnant!")
				fmt.Println("Sent: \"<@"+message.Author.ID+"> had sex with "+strings.Join(args[1:], " ")+" and failed to make them pregnant!\" to ", message.ChannelID)
			}

		case "kill":
			fmt.Println("Command is kill with args as ", args, " sent by ", message.Author)
			if len(Config.Deaths) == 0 {
				self.ChannelMessageSend(message.ChannelID, "No deaths configured.")
				fmt.Println("WARN: No images in Config.Deaths")
				break
			}
			dnum := rand.Intn(len(Config.Deaths))

			self.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+"> killed "+strings.Join(args[1:], " ")+" with a "+Config.Deaths[dnum])
			fmt.Println("Sent: \"<@"+message.Author.ID+"> killed "+strings.Join(args[1:], " ")+" with a "+Config.Deaths[dnum]+"\" to ", message.ChannelID)

		case "goonto": // yet another obv joke command
			fmt.Println("Command is goonto with args as ", args, " sent by ", message.Author)
			self.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+"> gooned to "+strings.Join(args[1:], " ")+"!!!!11!")

		case "eat":
			fmt.Println("Command is hug with args as ", args, " sent by ", message.Author)
			self.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+"> ate "+strings.Join(args[1:], " ")+" kindly with love <3")
			fmt.Println("Sent \"<@"+message.Author.ID+"> ate "+strings.Join(args[1:], " ")+" kindly with love <3\" to ", message.ChannelID)

		case "hug":
			fmt.Println("Command is eat with args as ", args, " sent by ", message.Author)
			if rand.Intn(56+56*2) == 8 {
				self.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+"> huged "+strings.Join(args[1:], " ")+" kindly... then started to eat them...")
				fmt.Println("Sent \"<@"+message.Author.ID+"> huged "+strings.Join(args[1:], " ")+" kindly... then started to eat them...\" to ", message.ChannelID)
			} else {
				self.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+"> huged "+strings.Join(args[1:], " ")+" kindly... nothing else...")
				fmt.Println("Sent \"<@"+message.Author.ID+"> huged "+strings.Join(args[1:], " ")+" kindly... nothing else...\" to ", message.ChannelID)
			}

		case "work":
			fmt.Println("Command is work with args as ", args, " sent by ", message.Author)
			num := rand.Intn(201)
			self.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+"> wroked and earned "+strconv.Itoa(num)+" money!")
			addMoney(message.Author.ID, num, self)
			fmt.Println("Sent \"<@"+message.Author.ID+"> wroked and earned "+strconv.Itoa(num)+" money!\" to ", message.ChannelID)

		case "rob":
			fmt.Println("Command is rob with args as ", args, " sent by ", message.Author)
			victim := strings.Join(args[1:], " ")
			num := rand.Intn(101)
			victim = strings.TrimPrefix(victim, "<@")
			victim = strings.TrimPrefix(victim, "!")
			victim = strings.TrimSuffix(victim, ">")
			e := rmMoney(victim, num, self)
			if !e {
				self.ChannelMessageSend(message.ChannelID, "User not found!")
				fmt.Println("No user found!")
				break
			}
			self.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+"> robbed <@"+victim+"> and stole "+strconv.Itoa(num)+" money!")
			addMoney(message.Author.ID, num, self)
			fmt.Println("Sent \"<@"+message.Author.ID+"> robbed <@"+victim+"> and stole "+strconv.Itoa(num)+" money!\" to ", message.ChannelID)

		case "help":
			fmt.Println("Command is help with args as ", args, " sent by ", message.Author)
			self.ChannelMessageSend(message.ChannelID, "Commands: help, kill, sex, getvro, eat, hug, goonto")
			fmt.Println("Sent help message to ", message.ChannelID)

		default:
			self.ChannelMessageSend(message.ChannelID, "Command "+cmd+" not found! Try the help command for more commands!")
			fmt.Println("No ", cmd, "command found!")
		}
	}
}

// every hour message hanndler
func HourlyMessage(self *discordgo.Session) {
	go func() {
		ticker := time.NewTicker(1 * time.Hour) // set up the hour control
		defer ticker.Stop()

		for {
			<-ticker.C
			fmt.Println("Hourly Message Sending to channels ", Config.Channel)

			for _, ch := range Config.Channel {
				server, err := self.Guild(ch.Server)
				if err != nil {
					fmt.Println("WARN: Couldnt get guild info:", err)
				}
				fmt.Println("Sending hourly message to ", ch.Channel, "in the guild", ch.Server, "with name", server.Name)

				SendRand(self, ch.Channel)
				fmt.Println("Sent!")

			}
			fmt.Println("Sent to all configured servers!")
		}
	}()
}
