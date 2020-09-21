package nekoslife

import (
	"errors"
	"fmt"
	"time"
	"encoding/json"
	"net/http"

	"github.com/jonas747/dcmd"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/dstate/v2"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/common"
)

var (
	ErrConnection = errors.New("Couldn't contact the API...")
	logger = common.GetPluginLogger(&Plugin{})

	TickleText = "<@!%d> tickles <@!%d>!"
	FeedText = "<@!%d> feeds <@!%d>!"
	PokeText = "<@!%d> pokes <@!%d>!"
	SlapText = "<@!%d> slaps <@!%d>!"
	PatText = "<@!%d> pats <@!%d>!"
	KissText = "<@!%d> kisses <@!%d>!"
	CuddleText = "<@!%d> cuddles <@!%d>!"
	HugText = "<@!%d> hugs <@!%d>!"
	BakaText = "<@!%d> calls <@!%d> a BAKA!"
)

type Plugin struct {}

func (p *Plugin) PluginInfo() *common.PluginInfo {
	return &common.PluginInfo{
		Name:     "Nekos.Life",
		SysName:  "nekos_life",
		Category: common.PluginCategoryMisc,
	}
}

func RegisterPlugin() {
	common.RegisterPlugin(&Plugin{})
}

var _ commands.CommandProvider = (*Plugin)(nil)

func (p *Plugin) AddCommands() {
	commands.AddRootCommands(p,
		// non-lewd
		WallpaperCommand,
		NekoGifCommand,
		GECGCommand,
		KemonomimiCommand,
		AvatarCommand,
		HoloCommand,
		WaifuCommand,
		NekoCommand,
		FoxgirlCommand,
		SmugCommand,
		WoofCommand,

		TickleCommand,
		FeedCommand,
		PokeCommand,
		SlapCommand,
		PatCommand,
		KissCommand,
		CuddleCommand,
		HugCommand,
		BakaCommand,

		// misc
		CatCommand,
		WhyCommand,
		FactCommand,
	)
}

type ImageResult struct {
	Url string
}

type CatResult struct {
	Cat string
}

type WhyResult struct {
	Why string
}

type FactResult struct {
	Fact string
}

func getJson(url string, target interface{}) error {
	var myClient = &http.Client{Timeout: 3 * time.Second}
    r, err := myClient.Get(url)
    if err != nil {
        return err
    }
    defer r.Body.Close()

    return json.NewDecoder(r.Body).Decode(target)
}

func getImage(endpoint string, data *dcmd.Data) (interface{}, error) {
	result := &ImageResult{}
	url := fmt.Sprintf("https://nekos.life/api/v2/img/%s", endpoint)
	err := getJson(url, result)

	if err != nil {
		return nil, err
	}

	embed := &discordgo.MessageEmbed{
		Image: &discordgo.MessageEmbedImage{
			URL: result.Url,
		},
	}

	return embed, nil
}

func getImageText(endpoint string, format string, data *dcmd.Data) (interface{}, error) {
	result := &ImageResult{}
	url := fmt.Sprintf("https://nekos.life/api/v2/img/%s", endpoint)
	err := getJson(url, result)
	
	var member *dstate.MemberState
	if data.Args[0].Value != nil {
		member = data.Args[0].Value.(*dstate.MemberState)
	} else {
		member = nil
	}

	if err != nil {
		return nil, err
	}

	var message string
	var embed interface{}

	if member == nil {
		message = fmt.Sprintf("Are you trying to %s the void...?", endpoint)
		embed = &discordgo.MessageEmbed{
			Description: message,
		}
	} else if member.ID == data.MS.ID {
		message = fmt.Sprintf("Sorry to see you alone, <@!%d> ;-;", member.ID)
		embed = &discordgo.MessageEmbed{
			Description: message,
		}
	} else {
		embed = &discordgo.MessageEmbed{
			Description: fmt.Sprintf(format, data.MS.ID, member.ID),
			Image: &discordgo.MessageEmbedImage{
				URL: result.Url,
			},
		}
	}

	return embed, nil
}

// Simple image commands (SFW)
var WallpaperCommand = &commands.YAGCommand{
	Cooldown:    5,
	CmdCategory: commands.CategoryAnime,
	Name:        "Wallpaper",
	Description: "Grabs a wallpaper from nekos.life.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) { return getImage("wallpaper", data) },
}

var NekoGifCommand = &commands.YAGCommand{
	Cooldown:    5,
	CmdCategory: commands.CategoryAnime,
	Name:        "NekoGif",
	Description: "Grabs a NekoGif from nekos.life.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) { return getImage("ngif", data) },
}

var GECGCommand = &commands.YAGCommand{
	Cooldown:    5,
	CmdCategory: commands.CategoryAnime,
	Name:        "GECG",
	Description: "Grabs a GeneticallyEngineeredCatGirl meme from nekos.life.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) { return getImage("gecg", data) },
}

var KemonomimiCommand = &commands.YAGCommand{
	Cooldown:    5,
	CmdCategory: commands.CategoryAnime,
	Name:        "Kemonomimi",
	Description: "Grabs a Kemonomimi from nekos.life.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) { return getImage("kemonomimi", data) },
}

var AvatarCommand = &commands.YAGCommand{
	Cooldown:    5,
	CmdCategory: commands.CategoryAnime,
	Name:        "AnimeAvatar",
	Description: "Grabs a Avatar from nekos.life.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) { return getImage("avatar", data) },
}

var HoloCommand = &commands.YAGCommand{
	Cooldown:    5,
	CmdCategory: commands.CategoryAnime,
	Name:        "Holo",
	Description: "Grabs a Holo from nekos.life.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) { return getImage("holo", data) },
}

var WaifuCommand = &commands.YAGCommand{
	Cooldown:    5,
	CmdCategory: commands.CategoryAnime,
	Name:        "Waifu",
	Description: "Grabs a Waifu from nekos.life.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) { return getImage("waifu", data) },
}

var NekoCommand = &commands.YAGCommand{
	Cooldown:    5,
	CmdCategory: commands.CategoryAnime,
	Name:        "Neko",
	Description: "Grabs a Neko from nekos.life.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) { return getImage("neko", data) },
}

var FoxgirlCommand = &commands.YAGCommand{
	Cooldown:    5,
	CmdCategory: commands.CategoryAnime,
	Name:        "Foxgirl",
	Description: "Grabs a Foxgirl from nekos.life.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) { return getImage("fox_girl", data) },
}

var SmugCommand = &commands.YAGCommand{
	Cooldown:    5,
	CmdCategory: commands.CategoryAnime,
	Name:        "Smug",
	Description: "Grabs a Smug from nekos.life.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) { return getImage("smug", data) },
}

var WoofCommand = &commands.YAGCommand{
	Cooldown:    5,
	CmdCategory: commands.CategoryAnime,
	Name:        "Woof",
	Description: "Grabs a Woof from nekos.life.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) { return getImage("woof", data) },
}


// Commands to kiss/poke/pat/etc someone.
var TickleCommand = &commands.YAGCommand{
	Cooldown:    5,
	CmdCategory: commands.CategoryAnime,
	Name:        "Tickle",
	Description: "Tickle someone with a GIF from nekos.life.",
	Arguments: []*dcmd.ArgDef{
		{Name: "User", Type: &commands.MemberArg{}},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) { return getImageText("tickle", TickleText, data) },
}

var FeedCommand = &commands.YAGCommand{
	Cooldown:    5,
	CmdCategory: commands.CategoryAnime,
	Name:        "Feed",
	Description: "Feed someone with a GIF from nekos.life.",
	Arguments: []*dcmd.ArgDef{
		{Name: "User", Type: &commands.MemberArg{}},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) { return getImageText("feed", FeedText, data) },
}

var PokeCommand = &commands.YAGCommand{
	Cooldown:    5,
	CmdCategory: commands.CategoryAnime,
	Name:        "Poke",
	Description: "Poke someone with a GIF from nekos.life.",
	Arguments: []*dcmd.ArgDef{
		{Name: "User", Type: &commands.MemberArg{}},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) { return getImageText("poke", PokeText, data) },
}

var SlapCommand = &commands.YAGCommand{
	Cooldown:    5,
	CmdCategory: commands.CategoryAnime,
	Name:        "Slap",
	Description: "Slap someone with a GIF from nekos.life.",
	Arguments: []*dcmd.ArgDef{
		{Name: "User", Type: &commands.MemberArg{}},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) { return getImageText("slap", SlapText, data) },
}

var PatCommand = &commands.YAGCommand{
	Cooldown:    5,
	CmdCategory: commands.CategoryAnime,
	Name:        "Pat",
	Description: "Pat someone with a GIF from nekos.life.",
	Arguments: []*dcmd.ArgDef{
		{Name: "User", Type: &commands.MemberArg{}},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) { return getImageText("pat", PatText, data) },
}

var KissCommand = &commands.YAGCommand{
	Cooldown:    5,
	CmdCategory: commands.CategoryAnime,
	Name:        "Kiss",
	Description: "Kiss someone with a GIF from nekos.life.",
	Arguments: []*dcmd.ArgDef{
		{Name: "User", Type: &commands.MemberArg{}},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) { return getImageText("kiss", KissText, data) },
}

var CuddleCommand = &commands.YAGCommand{
	Cooldown:    5,
	CmdCategory: commands.CategoryAnime,
	Name:        "Cuddle",
	Description: "Cuddle someone with a GIF from nekos.life.",
	Arguments: []*dcmd.ArgDef{
		{Name: "User", Type: &commands.MemberArg{}},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) { return getImageText("cuddle", CuddleText, data) },
}

var HugCommand = &commands.YAGCommand{
	Cooldown:    5,
	CmdCategory: commands.CategoryAnime,
	Name:        "Hug",
	Description: "Hug someone with a GIF from nekos.life.",
	Arguments: []*dcmd.ArgDef{
		{Name: "User", Type: &commands.MemberArg{}},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) { return getImageText("hug", HugText, data) },
}

var BakaCommand = &commands.YAGCommand{
	Cooldown:    5,
	CmdCategory: commands.CategoryAnime,
	Name:        "Baka",
	Description: "Call someone a BAKA with a GIF from nekos.life.",
	Arguments: []*dcmd.ArgDef{
		{Name: "User", Type: &commands.MemberArg{}},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) { return getImageText("baka", BakaText, data) },
}



// Misc text-only commands
var CatCommand = &commands.YAGCommand{
	Cooldown:    5,
	CmdCategory: commands.CategoryAnime,
	Name:        "CatEmoji",
	Description: "Grabs a random cat text emoji from nekos.life.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		result := &CatResult{}
		url := fmt.Sprintf("https://nekos.life/api/v2/cat")
		err := getJson(url, result)

		if err != nil {
			return nil, err
		}

		return result.Cat, nil
	},
}

var WhyCommand = &commands.YAGCommand{
	Cooldown:    5,
	CmdCategory: commands.CategoryAnime,
	Name:        "Why",
	Description: "Grabs a random why question from nekos.life.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		result := &WhyResult{}
		url := fmt.Sprintf("https://nekos.life/api/v2/why")
		err := getJson(url, result)

		if err != nil {
			return nil, err
		}

		return result.Why, nil
	},
}

var FactCommand = &commands.YAGCommand{
	Cooldown:    5,
	CmdCategory: commands.CategoryAnime,
	Name:        "Fact",
	Description: "Grabs a random fact from nekos.life.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		result := &FactResult{}
		url := fmt.Sprintf("https://nekos.life/api/v2/fact")
		err := getJson(url, result)

		if err != nil {
			return nil, err
		}

		return result.Fact, nil
	},
}
