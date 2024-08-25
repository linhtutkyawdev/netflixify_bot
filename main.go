package main

import (
	"log"
	"os"
	"strconv"
	"strings"
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
			// payload starts with "del"
			if strings.HasPrefix(c.Message().Payload, "del") {
				// remove prefix
				err := deleteVideo(c.Message().Payload[3:])
				if err != nil {
					log.Fatal(err)
					c.Send("Cannot delete the video you are locking for!")
				}

				return c.Send("Successfully deleted the video!")
			}

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
		f, err := b.FileByID(c.Message().Text)
		if err != nil {
			log.Fatal(err)
			return c.Send("No file found")
		}
		c.Send(f.FilePath)
		return c.Send(&tele.Photo{File: f})

		// return c.Send("Please click on /start!")
	})

	b.Handle(tele.OnChannelPost, channelHandlers(b))

	// ticker := time.NewTicker(time.Hour)

	// // Run the scheduled function in a goroutine so it doesn't block the main program
	// go func() {
	// 	for {
	// 		msg := refreshPosts(b)
	// 		if msg != "Post paths updated successfully!" {
	// 			log.Fatal(msg)
	// 		}
	// 		<-ticker.C
	// 	}
	// }()

	log.Println("Bot is now running.  Press CTRL-C to exit.")
	b.Start()
}
