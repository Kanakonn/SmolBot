package anilewds

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"io/ioutil"
	"time"
	"encoding/json"
	"net/http"

	"github.com/jonas747/dcmd"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/common"
	"github.com/jonas747/yagpdb/common/config"
)

var (
	ErrConnection = errors.New("Couldn't contact the API, or no results found...")
	logger = common.GetPluginLogger(&Plugin{})

	confUserName  = config.RegisterOption("yagpdb.danbooruusername", "Danbooru Username", "")
	confApiKey  = config.RegisterOption("yagpdb.danbooruapikey", "Danbooru API key", "")

	danbooruPostsURL = "https://danbooru.donmai.us/posts.json"
	danbooruPostURL = "https://danbooru.donmai.us/posts/%d"
)

type Plugin struct {}

func (p *Plugin) PluginInfo() *common.PluginInfo {
	return &common.PluginInfo{
		Name:     "AniLewd",
		SysName:  "ani_lewd",
		Category: common.PluginCategoryMisc,
	}
}

func RegisterPlugin() {
	common.RegisterPlugin(&Plugin{})
}

var _ commands.CommandProvider = (*Plugin)(nil)

func (p *Plugin) AddCommands() {
	commands.AddRootCommands(p,
		// questionable
		FeetCommand,
		YuriCommand,
		FutaCommand,
		KemonoCommand,
		SoloCommand,
		BoobsCommand,
		PussyCommand,

		// explicit
		EroFeetCommand,
		EroYuriCommand,
		EroFutaCommand,
		EroKemonoCommand,
		EroSoloCommand,
		EroBoobsCommand,
		EroPussyCommand,
	)
}

type PostResult struct {
	Id int64 `json:"id"`
	TagStringGeneral string `json:"tag_string_general"`
	TagStringCharacter string `json:"tag_string_character"`
	TagStringCopyright string `json:"tag_string_copyright"`
	FileUrl string `json:"file_url"`
}

func getJson(request *http.Request) (error, []PostResult) {
	var myClient = &http.Client{Timeout: 3 * time.Second}

	r, err := myClient.Do(request)
	if err != nil {
		return err, nil
	}
	defer r.Body.Close()

	if r.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("%s", err)
		}
		bodyString := string(bodyBytes)
		fmt.Printf("Resp: %s\n", bodyString)
		var result []PostResult
		decodeErr := json.Unmarshal([]byte(bodyString), &result)
		return decodeErr, result
	} else {
		logger.Error(fmt.Sprintf("Error during JSON Get: Status code %d", r.StatusCode))
		return nil, nil
	}
}

func basicAuth(username string, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func getImage(tags string, data *dcmd.Data) (interface{}, error) {
	var result []PostResult

	request, _ := http.NewRequest("GET", danbooruPostsURL, nil)
	request.Header.Add("Authorization", "Basic " + basicAuth(confUserName.GetString(), confApiKey.GetString()))

	q := request.URL.Query()
    q.Add("limit", "1")
    q.Add("random", "true")
    q.Add("tags", tags)
    request.URL.RawQuery = q.Encode()

	err, result := getJson(request)

	if err != nil {
		logger.Warn(fmt.Sprintf("Error on JSON get: %s", err))
		return nil, err
	}

	if len(result) > 0 {
		var footer *discordgo.MessageEmbedFooter

		if result[0].TagStringCharacter != "" && result[0].TagStringCopyright != "" {
			footer = &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Characters: %s\nCopyright: %s", result[0].TagStringCharacter, result[0].TagStringCopyright),
			}
		} else if result[0].TagStringCharacter != "" {
			footer = &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Characters: %s", result[0].TagStringCharacter),
			}
		} else if result[0].TagStringCopyright != "" {
			footer = &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Copyright: %s", result[0].TagStringCopyright),
			}
		}
		var tags *discordgo.MessageEmbedField
		if len(strings.Split(result[0].TagStringGeneral, " ")) > 20 {
			more := fmt.Sprintf(" (and %d more)", len(strings.Split(result[0].TagStringGeneral, " ")) - 20)
			tags = &discordgo.MessageEmbedField{
				Name:   "Tags",
				Value:  strings.Join(strings.Split(result[0].TagStringGeneral, " ")[:20], " ") + more,
				Inline: false,
			}
		} else {
			tags = &discordgo.MessageEmbedField{
				Name:   "Tags",
				Value:  strings.Join(strings.Split(result[0].TagStringGeneral, " "), " "),
				Inline: false,
			}
		}

		embed := &discordgo.MessageEmbed{
			URL: fmt.Sprintf(danbooruPostURL, result[0].Id),
			Title: fmt.Sprintf("Danbooru #%d", result[0].Id),
			Image: &discordgo.MessageEmbedImage{
				URL: result[0].FileUrl,
			},
			Fields: []*discordgo.MessageEmbedField{
				tags,
			},
			Footer: footer,
		}

		return embed, nil
	}

	logger.Warn(fmt.Sprintf("Error on result: %s", ErrConnection))
	return nil, ErrConnection
}


// Questionable commands
var FeetCommand = &commands.YAGCommand{
	Cooldown:    3,
	CmdCategory: commands.CategoryLewdAnime,
	Name:        "Feet",
	Description: "Grabs a random, semi-lewd feet image from Danbooru.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		return getImage("feet rating:questionable -loli -shota -toddlercon -animated -flash favcount:>100 status:active parent:none", data)
	},
}

var YuriCommand = &commands.YAGCommand{
	Cooldown:    3,
	CmdCategory: commands.CategoryLewdAnime,
	Name:        "Yuri",
	Description: "Grabs a random, semi-lewd yuri image from Danbooru.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		return getImage("yuri rating:questionable -loli -shota -toddlercon -animated -flash favcount:>100 status:active parent:none", data)
	},
}

var FutaCommand = &commands.YAGCommand{
	Cooldown:    3,
	CmdCategory: commands.CategoryLewdAnime,
	Name:        "Futa",
	Description: "Grabs a random, semi-lewd futa image from Danbooru.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		return getImage("futa rating:questionable -loli -shota -toddlercon -animated -flash favcount:>100 status:active parent:none", data)
	},
}

var KemonoCommand = &commands.YAGCommand{
	Cooldown:    3,
	CmdCategory: commands.CategoryLewdAnime,
	Name:        "Kemono",
	Description: "Grabs a random, semi-lewd kemono image from Danbooru.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		return getImage("animal_ears rating:questionable -loli -shota -toddlercon -animated -flash favcount:>100 status:active parent:none", data)
	},
}

var SoloCommand = &commands.YAGCommand{
	Cooldown:    3,
	CmdCategory: commands.CategoryLewdAnime,
	Name:        "Solo",
	Description: "Grabs a random, semi-lewd solo image from Danbooru.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		return getImage("solo 1girl rating:questionable -loli -shota -toddlercon -animated -flash favcount:>100 status:active parent:none", data)
	},
}

var BoobsCommand = &commands.YAGCommand{
	Cooldown:    3,
	CmdCategory: commands.CategoryLewdAnime,
	Name:        "Boobs",
	Description: "Grabs a random, semi-lewd boobs image from Danbooru.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		return getImage("breasts rating:questionable -loli -shota -toddlercon -animated -flash favcount:>100 status:active parent:none", data)
	},
}

var PussyCommand = &commands.YAGCommand{
	Cooldown:    3,
	CmdCategory: commands.CategoryLewdAnime,
	Name:        "Pussy",
	Description: "Grabs a random, semi-lewd pussy image from Danbooru.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		return getImage("pussy rating:questionable -loli -shota -toddlercon -animated -flash favcount:>100 status:active parent:none", data)
	},
}


// Explicit commands
var EroFeetCommand = &commands.YAGCommand{
	Cooldown:    3,
	CmdCategory: commands.CategoryLewdAnime2,
	Name:        "EroFeet",
	Description: "Grabs a random, very lewd feet image from Danbooru.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		return getImage("feet rating:explicit -loli -shota -toddlercon -animated -flash favcount:>100 status:active parent:none", data)
	},
}

var EroYuriCommand = &commands.YAGCommand{
	Cooldown:    3,
	CmdCategory: commands.CategoryLewdAnime2,
	Name:        "EroYuri",
	Description: "Grabs a random, very lewd yuri image from Danbooru.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		return getImage("yuri rating:explicit -loli -shota -toddlercon -animated -flash favcount:>100 status:active parent:none", data)
	},
}

var EroFutaCommand = &commands.YAGCommand{
	Cooldown:    3,
	CmdCategory: commands.CategoryLewdAnime2,
	Name:        "EroFuta",
	Description: "Grabs a random, very lewd futa image from Danbooru.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		return getImage("futa rating:explicit -loli -shota -toddlercon -animated -flash favcount:>100 status:active parent:none", data)
	},
}

var EroKemonoCommand = &commands.YAGCommand{
	Cooldown:    3,
	CmdCategory: commands.CategoryLewdAnime2,
	Name:        "EroKemono",
	Description: "Grabs a random, very lewd kemono image from Danbooru.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		return getImage("animal_ears rating:explicit -loli -shota -toddlercon -animated -flash favcount:>100 status:active parent:none", data)
	},
}

var EroSoloCommand = &commands.YAGCommand{
	Cooldown:    3,
	CmdCategory: commands.CategoryLewdAnime2,
	Name:        "EroSolo",
	Description: "Grabs a random, very lewd solo image from Danbooru.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		return getImage("solo 1girl rating:explicit -loli -shota -toddlercon -animated -flash favcount:>100 status:active parent:none", data)
	},
}

var EroBoobsCommand = &commands.YAGCommand{
	Cooldown:    3,
	CmdCategory: commands.CategoryLewdAnime2,
	Name:        "EroBoobs",
	Description: "Grabs a random, very lewd boobs image from Danbooru.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		return getImage("breasts rating:explicit -loli -shota -toddlercon -animated -flash favcount:>100 status:active parent:none", data)
	},
}

var EroPussyCommand = &commands.YAGCommand{
	Cooldown:    3,
	CmdCategory: commands.CategoryLewdAnime2,
	Name:        "EroPussy",
	Description: "Grabs a random, very lewd pussy image from Danbooru.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		return getImage("pussy rating:explicit -loli -shota -toddlercon -animated -flash favcount:>100 status:active parent:none", data)
	},
}
