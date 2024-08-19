package main

import (
	"crypto/rand"
	"log"
	"math/big"
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

	// Create main menu
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}

	b.Handle("/start", func(c tele.Context) error {
		url := os.Getenv("APP_URL") + "/thumbnail?tgId=" + strconv.FormatInt(c.Sender().ID, 10) + "&tgUsername=" + c.Sender().Username
		// Default buttons.
		btnWatch := menu.WebApp("🍿 Watch", &tele.WebApp{URL: os.Getenv("APP_URL")})
		btnCreateThumbnail := menu.WebApp("🪄 Create Thumbnail", &tele.WebApp{URL: url})
		// Pages
		home_page := []tele.Row{
			menu.Row(btnWatch),
			menu.Row(btnCreateThumbnail),
		}

		menu.Reply(home_page...)

		b.Handle(tele.OnText, func(c tele.Context) error {
			return c.Send("You can click buttons to do actions!")
		})

		return c.Send("Please click the \"Watch 🍿\" button!", menu)
	})
	b.Handle("/admin", func(c tele.Context) error {
		url := os.Getenv("APP_URL") + "/thumbnail?tgId=" + strconv.FormatInt(c.Sender().ID, 10) + "&tgUsername=" + c.Sender().Username
		// Menu buttons
		btnCreateThumbnail := menu.WebApp("🪄 Create Thumbnail", &tele.WebApp{URL: url})
		btnRegister := menu.Text("📝 Register")
		btnLogin := menu.Text("👤 Login")
		btnCreatePost := menu.Text("🪄 Create Post")
		btnSettings := menu.Text("⚙ Settings")

		// Setting buttons
		btnHelp := menu.Text("ℹ Help")
		btnBack := menu.Text("⬅ Back")

		// Pages
		home_page := []tele.Row{
			menu.Row(btnCreateThumbnail),
			menu.Row(btnRegister, btnLogin),
		}

		logined_page := []tele.Row{
			menu.Row(btnCreateThumbnail),
			menu.Row(btnCreatePost, btnSettings),
		}

		settings_page := []tele.Row{
			menu.Row(btnHelp),
			menu.Row(btnBack),
		}

		//  = c.Sender().ID
		menu.Reply(home_page...)

		b.Handle(tele.OnText, func(c tele.Context) error {
			return c.Send("You can click buttons to do actions!")
		})

		b.Handle(&btnRegister, func(c tele.Context) error {
			// when an id is sent
			b.Handle(tele.OnText, func(c tele.Context) error {
				// set default on error
				b.Handle(tele.OnText, func(c tele.Context) error {
					return c.Send("You can click buttons to do actions!")
				})
				id, err := strconv.ParseInt(c.Text(), 10, 64)
				if err != nil {
					return c.Send("Invalid channel id")
				}
				chat, err := b.ChatByID(id)
				if err != nil {
					return c.Send("Please make sure that the channel id provided has already been added the bot as an admin.")
				}
				//send 6 random otp
				otp, err := rand.Int(rand.Reader, big.NewInt(999999))
				if err != nil {
					return c.Send("Failed to generate OTP")
				}
				b.Send(chat, "Your OTP is \""+otp.String()+"\"")

				// when an otp is sent
				b.Handle(tele.OnText, func(c tele.Context) error {
					// set default on error
					b.Handle(tele.OnText, func(c tele.Context) error {
						return c.Send("You can click buttons to do actions!")
					})

					if c.Text() != otp.String() {
						return c.Send("Your OTP is incorrect")
					}

					// when a password is sent
					b.Handle(tele.OnText, func(c tele.Context) error {
						if len(c.Text()) < 6 {
							return c.Send("Password Must Be At Least 6 Characters Long! Try sending a valid one again.")
						}

						// set default on error
						b.Handle(tele.OnText, func(c tele.Context) error {
							return c.Send("You can click buttons to do actions!")
						})

						if hasRegistered(id) {
							return c.Send("This channel has already been registered! Please try logging in.")
						}

						return c.Send(registerChannel(id, chat.Title, c.Text()))
					})
					return c.Send("Your OTP is correct! Please send a password to register.")
				})
				return c.Send("Please enter the OTP sent to " + chat.Title + ".")
			})
			return c.Send("Please send the channel id to register.")
		})

		b.Handle(&btnLogin, func(c tele.Context) error {
			// when an id is sent
			b.Handle(tele.OnText, func(c tele.Context) error {
				// set default on error
				b.Handle(tele.OnText, func(c tele.Context) error {
					return c.Send("You can click buttons to do actions!")
				})
				id, err := strconv.ParseInt(c.Text(), 10, 64)
				if err != nil {
					return c.Send("Invalid channel id")
				}
				_, err = b.ChatByID(id)
				if err != nil {
					return c.Send("Please make sure that the channel id provided has already been added the bot as an admin.")
				}

				// when a password is sent
				b.Handle(tele.OnText, func(c tele.Context) error {
					if len(c.Text()) < 6 {
						return c.Send("Password Must Be At Least 6 Characters Long! Try sending a valid one again.")
					}

					// set default on error
					b.Handle(tele.OnText, func(c tele.Context) error {
						return c.Send("You can click buttons to do actions!")
					})

					msg := loginChannel(id, c.Text())
					if msg != "" {
						return c.Send(msg)
					}
					// change menu
					menu.Reply(logined_page...)
					return c.Send("Logged in successfully!", &tele.SendOptions{
						ReplyMarkup: menu,
					})
				})

				return c.Send("Please send the password to login.")
			})
			return c.Send("Please send the channel id to login.")
		})
		// b.Handle(&btnCreatePost, func(c tele.Context) error {
		// 	b.Handle(tele.OnVideo, func(c tele.Context) error {
		// 		b.Handle(tele.OnVideo, nil)
		// 		// selector.Inline(
		// 		// 	selector.Row(btnYes, btnNo),
		// 		// )
		// 		_, err := c.Message().Video.Send(b, &tele.Chat{ID: -1002085983072}, &tele.SendOptions{
		// 			ReplyMarkup: selector,
		// 		})

		// 		return err
		// 	})
		// 	return c.Send("Please Send a Video")
		// })

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

		return c.Send("Please click the \"Watch 🍿\" button!", menu)
	})

	b.Handle(tele.OnText, func(c tele.Context) error {
		return c.Send("Please click on /start!")
	})

	b.Handle(tele.OnChannelPost, channelHandlers(b))

	b.Start()
}
