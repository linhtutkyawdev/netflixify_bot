package main

import (
	"log"
	"os"
	"strconv"
	"time"

	tele "gopkg.in/telebot.v3"
)

func main() {
	// Create a new bot instance
	pref := tele.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/start", func(c tele.Context) error {
		// Create main menu
		menu := &tele.ReplyMarkup{ResizeKeyboard: true}
		selector := &tele.ReplyMarkup{}

		url := "https://ad41-202-165-88-192.ngrok-free.app/thumbnail?tgId=" + strconv.FormatInt(c.Sender().ID, 10) + "&username=" + c.Sender().Username

		// Default buttons.
		btnCreateThumbnail := menu.WebApp("🪄 Create Thumbnail", &tele.WebApp{URL: url})

		btnCreatePost := menu.Text("🪄 Create Post")
		btnSettings := menu.Text("⚙ Settings")

		// Setting buttons
		btnHelp := menu.Text("ℹ Help")
		btnBack := menu.Text("⬅ Back")

		btnButton := selector.Data("Button", "Button")

		// Pages
		home_page := []tele.Row{
			menu.Row(btnCreatePost),
			menu.Row(btnCreateThumbnail),
			menu.Row(btnSettings),
		}

		settings_page := []tele.Row{
			menu.Row(btnHelp),
			menu.Row(btnBack),
		}

		//  = c.Sender().ID
		menu.Reply(home_page...)

		b.Handle(tele.OnText, func(c tele.Context) error {
			return c.Send(strconv.Itoa(c.Message().ID))
		})

		b.Handle(&btnButton, func(c tele.Context) error {
			return c.Send("Button")
		})

		// Default Button Handlers
		b.Handle(&btnCreateThumbnail, func(c tele.Context) error {
			return c.Send("Create T")
		})

		b.Handle(&btnCreatePost, func(c tele.Context) error {
			b.Handle(tele.OnVideo, func(c tele.Context) error {
				b.Handle(tele.OnVideo, nil)
				selector.Inline(
					selector.Row(btnButton),
				)
				_, err := c.Message().Video.Send(b, &tele.Chat{ID: -1002085983072}, &tele.SendOptions{
					ReplyMarkup: selector,
				})

				return err
			})
			return c.Send("Please Send a Video")
		})

		b.Handle(&btnSettings, func(c tele.Context) error {
			menu.Reply(settings_page...)
			return c.Send("Settings", menu)
		})

		b.Handle(&btnHelp, func(c tele.Context) error {
			return c.Send("Help")
		})

		b.Handle(&btnBack, func(c tele.Context) error {
			menu.Reply(home_page...)
			return c.Send("Home", menu)
		})

		return c.Send("Admin Menu Unlocked!", menu)
	})

	b.Handle("/admin", func(c tele.Context) error {
		return c.Send("To get started, use /start")
	})

	b.Handle(tele.OnText, func(c tele.Context) error {
		return c.Send(`Please use the blue "Netflixify" button to explore the world of Entertainment.`)
	})

	b.Handle(tele.OnChannelPost, func(c tele.Context) error {
		if c.Text() == "channel_id" {
			msg, err := b.Send(c.Chat(), strconv.FormatInt(c.Chat().ID, 10))

			if err != nil {
				return nil
			}
			time.AfterFunc(10*time.Second, func() {
				b.Delete(c.Message())
				b.Delete(msg)
			})
			return nil
		}
		return nil
	})
	b.Start()
}
