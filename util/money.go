package util

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func addMoney(userid string, amount int, self *discordgo.Session) {
	for i, user := range userdata.Users {
		if user.UserID == userid {
			userdata.Users[i].Money += amount
			self.ChannelMessageSend(Config.GChannel, "<@"+userid+"> gained "+strconv.Itoa(amount)+" money!")
			fmt.Println(userid, "gained ", amount)
			writeData()
			return
		}
	}
	userdata.Users = append(userdata.Users, UserData{
		UserID:  userid,
		Money:   amount,
		Supplys: []string{},
	})
	fmt.Println("Added entry for ", userid)
	writeData()
}

func rmMoney(userid string, amount int, self *discordgo.Session) bool {
	for i, user := range userdata.Users {
		if user.UserID == userid {
			userdata.Users[i].Money -= amount
			self.ChannelMessageSend(Config.GChannel, "<@"+userid+"> lost "+strconv.Itoa(amount)+" money!")
			fmt.Println(userid, "lost ", amount)
			writeData()
			return true
		}
	}
	return false
}
