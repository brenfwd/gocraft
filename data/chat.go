package data

import (
	"encoding/json"
	"errors"
	"regexp"
)

type ChatColor string

const (
	ChatColorBlack       ChatColor = "black"
	ChatColorDarkBlue    ChatColor = "dark_blue"
	ChatColorDarkGreen   ChatColor = "dark_green"
	ChatColorDarkAqua    ChatColor = "dark_aqua"
	ChatColorDarkRed     ChatColor = "dark_red"
	ChatColorDarkPurple  ChatColor = "dark_purple"
	ChatColorGold        ChatColor = "gold"
	ChatColorGray        ChatColor = "gray"
	ChatColorDarkGray    ChatColor = "dark_gray"
	ChatColorBlue        ChatColor = "blue"
	ChatColorGreen       ChatColor = "green"
	ChatColorAqua        ChatColor = "aqua"
	ChatColorRed         ChatColor = "red"
	ChatColorLightPurple ChatColor = "light_purple"
	ChatColorYellow      ChatColor = "yellow"
	ChatColorWhite       ChatColor = "white"
)

func (c ChatColor) String() string {
	return string(c)
}

const colorHexRegexp string = `^#[0-9a-fA-F]{6}$`

func ChatColorHex(hex string) (ChatColor, error) {
	re := regexp.MustCompile(colorHexRegexp)
	if re.MatchString(hex) {
		return ChatColor(hex), nil
	}
	return "", errors.New("invalid hex color")
}

type ChatFont string

const (
	// Default font
	ChatFontDefault ChatFont = "minecraft:default"
	// Unicode font
	ChatFontUnicode ChatFont = "minecraft:uniform"
	// Enchanting Table font
	ChatFontAlt ChatFont = "minecraft:alt"
	// Unused
	ChatFontIllagerAlt ChatFont = "minecraft:illageralt"
)

type Chat struct {
	Text  *string    `json:"text,omitempty"`
	Color *ChatColor `json:"color,omitempty"`
	Font  *ChatFont  `json:"font,omitempty"`

	Extra []*Chat `json:"extra,omitempty"`

	// Styles
	Bold          bool `json:"bold,omitempty"`
	Italic        bool `json:"italic,omitempty"`
	Underlined    bool `json:"underlined,omitempty"`
	Strikethrough bool `json:"strikethrough,omitempty"`
	Obfuscated    bool `json:"obfuscated,omitempty"`
}

func MakeChat() *Chat {
	return new(Chat)
}

func (c *Chat) String() (string, error) {
	j, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

func (c *Chat) SetText(value string) *Chat {
	c.Text = &value
	return c
}

func (c *Chat) RemoveText() *Chat {
	c.Text = nil
	return c
}

func (c *Chat) SetColor(value ChatColor) *Chat {
	c.Color = &value
	return c
}

func (c *Chat) RemoveColor() *Chat {
	c.Color = nil
	return c
}

func (c *Chat) SetFont(value ChatFont) *Chat {
	c.Font = &value
	return c
}

func (c *Chat) RemoveFont() *Chat {
	c.Font = nil
	return c
}

func (c *Chat) AddExtra(newChat ...*Chat) *Chat {
	c.Extra = append(c.Extra, newChat...)
	return c
}

func (c *Chat) BuildExtra(callback func(*Chat) *Chat) *Chat {
	var c2 Chat
	c.AddExtra(callback(&c2))
	return c
}

func (c *Chat) SetBold(value bool) *Chat {
	c.Bold = value
	return c
}

func (c *Chat) SetItalic(value bool) *Chat {
	c.Italic = value
	return c
}

func (c *Chat) SetUnderlined(value bool) *Chat {
	c.Underlined = value
	return c
}

func (c *Chat) SetStrikethrough(value bool) *Chat {
	c.Strikethrough = value
	return c
}

func (c *Chat) SetObfuscated(value bool) *Chat {
	c.Obfuscated = value
	return c
}
