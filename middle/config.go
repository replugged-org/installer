package middle

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/ProtonMail/go-appdir"
)

type Config struct {
	DiscordPath string `json:"discordPath"`
}

func getConfigPath() string {
	cfg, err := os.UserConfigDir()
	if err != nil {
		cfg = ""
	}
	return filepath.Join(cfg, "replugged-installer.json")
}

type UserData struct {
	Name string
	Ruid int
	Rgid int
}
func GetUserData() *UserData {
	cmd := exec.Command("logname")
	outName, _:= cmd.Output()
	name := strings.TrimSuffix(string(outName), "\n")

	cmd = exec.Command("id", "-u", name)
	outUid, _ := cmd.Output()
	uid := strings.TrimSuffix(string(outUid), "\n")

	cmd = exec.Command("id", "-g", name)
	outGid, _ := cmd.Output()
	gid := strings.TrimSuffix(string(outGid), "\n")

	ruid, _ := strconv.Atoi(uid)
	rgid, _ := strconv.Atoi(gid)
	return &UserData{name, ruid, rgid}
}

func GetDataPath() string {
	if runtime.GOOS != "darwin" && runtime.GOOS != "windows" {
		userData := GetUserData()
		if userData.Name != "" {
			usr := strings.TrimSuffix(string(userData.Name), "\n")
			return fmt.Sprintf("/home/%s/.local/share/replugged-installer", usr)
		} else {
			return ""
		}
	}

	dirs := appdir.New("replugged-installer")
	data := dirs.UserData()

	if err := os.MkdirAll(data, 0755); err != nil {
		panic(err)
	}

	return data
}

func ReadConfig() Config {
	cfg := &Config{}
	cfgFile, err := os.Open(getConfigPath())
	if err != nil {
		fmt.Printf("Failed to open config: %s\n", err.Error())
	}
	defer cfgFile.Close()
	if json.NewDecoder(cfgFile).Decode(cfg) != nil {
		fmt.Printf("Failed to decode config\n")
		return *cfg
	}
	WriteConfig(*cfg)
	return *cfg
}

func WriteConfig(cfg Config) {
	cfgPath := getConfigPath()
	os.MkdirAll(filepath.Dir(cfgPath), os.ModePerm)
	cfgFile, err := os.OpenFile(cfgPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to write config: %s\n", err.Error())
		return
	}
	defer cfgFile.Close()
	if json.NewEncoder(cfgFile).Encode(cfg) != nil {
		fmt.Printf("Failed to encode config\n")
	}
}
