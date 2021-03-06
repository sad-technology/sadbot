package bot

import (
	"github.com/DiscoreMe/sadbot/cache"
	"github.com/DiscoreMe/sadbot/calculator"
	"github.com/DiscoreMe/sadbot/dict"
	"github.com/DiscoreMe/sadbot/weather"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
	"time"
	"unicode/utf8"
)

type Bot struct {
	bot  *tb.Bot
	w    *weather.Weather
	c    *cache.Cache
	d    *dict.Dict
	calc *calculator.Cal
}

type BotSettings struct {
	Token   string
	Weather *weather.Weather
	Cache   *cache.Cache
	Dict    *dict.Dict
	Calc    *calculator.Cal
}

func NewBot(settings BotSettings) (*Bot, error) {
	b, err := tb.NewBot(tb.Settings{
		Token:  settings.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		bot:  b,
		c:    settings.Cache,
		d:    settings.Dict,
		w:    settings.Weather,
		calc: settings.Calc,
	}
	bot.setup()

	return bot, nil
}

func (b *Bot) Listen() {
	b.bot.Start()
}

func (b *Bot) setup() {
	b.bot.Handle(".hello", func(m *tb.Message) {
		b.bot.Send(m.Sender, "Hello World!")
	})
	b.bot.Handle(tb.OnText, b.CmdHandler)
}

func (b *Bot) CmdHandler(m *tb.Message) {
	if m.Text == "" {
		return
	}
	if m.Text[0] != '.' {
		b.SpeakHandler(m)
		return
	}
	if utf8.RuneCountInString(m.Text) <= 1 {
		return
	}

	cmd := strings.Split(m.Text, " ")[0][1:]
	switch cmd {
	case "погода":
		b.WeatherHandler(m)
	case "адик":
		b.SpeakAddHandler(m)
	case "эбауте":
		b.about(m)
	case "кл":
		b.CalcHandler(m)
	default:
		b.SpeakHandler(m)
	}
}

func (b *Bot) about(m *tb.Message) {
	b.bot.Send(m.Chat, "Исходный код:\nhttps://github.com/DiscoreMe/sadbot")
}
