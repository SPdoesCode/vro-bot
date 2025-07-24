package util

// the d & m structs

type UserData struct {
	UserID  string   `json:"userid"`
	Money   int      `json:"money"`
	Supplys []string `json:"supplys"`
}

type DataFile struct {
	Users []UserData `json:"users"`
}

// config structs

type channels struct {
	Server  string `toml:"server"`
	Channel string `toml:"channel"`
}

type CFG struct {
	Token    string     `toml:"token"`
	Prefix   string     `toml:"prefix"`
	ImageDir string     `toml:"image_dir"`
	Channel  []channels `toml:"channel"`
	Deaths   []string   `toml:"deaths"`
	GChannel string     `toml:"gchannel"`
	Images   []string
}
