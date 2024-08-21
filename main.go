package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
)

func main() {
	// Load .env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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

		// on payload
		if c.Message().Payload != "" {
			vid, err := findVideoId(c.Message().Payload)
			if err != nil {
				c.Send("Cannot find the video you are locking for!")
			}
			videoFile, err := b.FileByID(vid)
			if err != nil {
				c.Send("Cannot find the video you are locking for!")
			}
			return c.Send(&tele.Video{
				File: videoFile,
			})
		}

		url := os.Getenv("APP_URL") + "/thumbnail?tgId=" + strconv.FormatInt(c.Sender().ID, 10) + "&tgUsername=" + c.Sender().Username
		menu := &tele.ReplyMarkup{ResizeKeyboard: true}

		btnWatch := menu.WebApp("🍿 Watch", &tele.WebApp{URL: os.Getenv("APP_URL")})
		btnCreateThumbnail := menu.WebApp("🪄 Create Thumbnail", &tele.WebApp{URL: url})

		menu.Reply(menu.Row(btnWatch), menu.Row(btnCreateThumbnail))

		b.Handle(tele.OnText, func(c tele.Context) error {
			return c.Send("You can click buttons to do actions!")
		})

		return c.Send("Please click the \"Watch 🍿\" button!", menu)
	})

	b.Handle("/admin", func(c tele.Context) error {
		url := os.Getenv("APP_URL") + "/thumbnail?tgId=" + strconv.FormatInt(c.Sender().ID, 10) + "&tgUsername=" + c.Sender().Username
		menu := &tele.ReplyMarkup{ResizeKeyboard: true}
		// Menu buttons
		btnCreateThumbnail := menu.WebApp("🪄 Create Thumbnail", &tele.WebApp{URL: url})
		btnRegister := menu.Text("📝 Register")
		btnLogin := menu.Text("👤 Login")

		menu.Reply(menu.Row(btnCreateThumbnail), menu.Row(btnRegister, btnLogin))

		b.Handle(tele.OnText, func(c tele.Context) error {
			return c.Send("You can click buttons to do actions!")
		})

		b.Handle(&btnRegister, registerHandler(b))

		b.Handle(&btnLogin, loginHandler(b))

		return c.Send("You can click buttons to do actions!", menu)
	})

	b.Handle(tele.OnText, func(c tele.Context) error {
		return c.Send("Please click on /start!")
	})

	b.Handle(tele.OnPhoto, func(c tele.Context) error {
		photo, err := b.FileByID(c.Message().Photo.FileID)
		if err != nil {
			return nil
		}
		return c.Send(photo.FilePath)
	})

	b.Handle(tele.OnVideo, func(c tele.Context) error {
		video, err := b.FileByID(c.Message().Video.FileID)
		if err != nil {
			return nil
		}
		return c.Send(video.FilePath)
	})

	b.Handle(tele.OnChannelPost, channelHandlers(b))

	b.Start()
}
