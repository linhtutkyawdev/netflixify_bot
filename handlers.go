package main

import (
	"crypto/rand"
	"math/big"
	"os"
	"strconv"
	"time"

	tele "gopkg.in/telebot.v3"
)

func channelHandlers(b *tele.Bot) tele.HandlerFunc {
	return func(c tele.Context) error {
		if c.Text() == "/id" {
			b.Edit(c.Message(), strconv.FormatInt(c.Chat().ID, 10))
			time.AfterFunc(10*time.Second, func() {
				b.Delete(c.Message())
			})
		}
		return nil
	}
}

func registerHandler(b *tele.Bot) tele.HandlerFunc {
	return func(c tele.Context) error {
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
			if hasRegistered(id) {
				return c.Send("This channel has already been registered!")
			}
			//send 6 random otp
			otp, err := rand.Int(rand.Reader, big.NewInt(999999))
			if err != nil {
				return c.Send("Failed to generate OTP")
			}
			msg, err := b.Send(chat, "Your OTP is \""+otp.String()+"\"")
			if err != nil {
				return c.Send("Failed to send OTP")
			}
			time.AfterFunc(time.Minute, func() {
				b.Delete(msg)
			})

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
	}
}

func loginHandler(b *tele.Bot) tele.HandlerFunc {
	return func(c tele.Context) error {
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

			if !hasRegistered(id) {
				return c.Send("This channel has not been registered! Please register first.")
			}

			channel, err := b.ChatByID(id)
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
				// when logg in success
				menu := &tele.ReplyMarkup{ResizeKeyboard: true}
				url := os.Getenv("APP_URL") + "/thumbnail?tgId=" + strconv.FormatInt(c.Sender().ID, 10) + "&tgUsername=" + c.Sender().Username

				btnCreatePost := menu.Text("🪄 Create Post")
				btnCreateThumbnail := menu.WebApp("🪄 Create Thumbnail", &tele.WebApp{URL: url})
				btnSettings := menu.Text("⚙ Settings")

				// Setting buttons
				btnRefresh := menu.Text("⟳ Refresh")
				btnHelp := menu.Text("ℹ Help")
				btnBack := menu.Text("⬅ Back")

				b.Handle(&btnSettings, func(c tele.Context) error {
					menu.Reply(menu.Row(btnRefresh), menu.Row(btnHelp, btnBack))
					return c.Send("Settings", menu)
				})

				b.Handle(&btnRefresh, func(c tele.Context) error {
					return c.Send(refreshPosts(b))
				})

				b.Handle(&btnHelp, func(c tele.Context) error {
					return c.Send("Help")
				})

				b.Handle(&btnBack, func(c tele.Context) error {
					menu.Reply(menu.Row(btnCreateThumbnail), menu.Row(btnCreatePost, btnSettings))
					return c.Send("Home", menu)
				})

				b.Handle(&btnCreatePost, func(c tele.Context) error {
					selector := &tele.ReplyMarkup{}
					btnCancel := selector.Data("❌ Cancel", "Cancel")
					selector.Inline(selector.Row(btnCancel))

					b.Handle(&btnCancel, func(c tele.Context) error {
						b.Handle(tele.OnText, func(c tele.Context) error {
							return c.Send("You can click buttons to do actions!")
						})
						b.Handle(tele.OnPhoto, func(c tele.Context) error {
							return nil
						})
						b.Handle(tele.OnVideo, func(c tele.Context) error {
							return nil
						})
						return c.Send("Operation cancelled!")
					})

					b.Handle(tele.OnText, func(c tele.Context) error {
						title := c.Text()
						b.Handle(tele.OnText, func(c tele.Context) error {
							rating, err := strconv.Atoi(c.Text())
							if err != nil || rating < 1 || rating > 100 {
								return c.Send("Invalid video rating! Please try again.", selector)
							}

							b.Handle(tele.OnText, func(c tele.Context) error {
								description := c.Text()

								b.Handle(tele.OnText, func(c tele.Context) error {
									tags := c.Text()

									b.Handle(tele.OnText, func(c tele.Context) error {
										return c.Send("You need to send a video", selector)
									})

									b.Handle(tele.OnVideo, func(c tele.Context) error {
										if c.Message().Video == nil {
											return c.Send("Video is invalid! Please try again.", selector)
										}

										video, err := b.FileByID(c.Message().Video.FileID)
										if err != nil {
											return c.Send("File not found! Please try again.", selector)
										}

										b.Handle(tele.OnText, func(c tele.Context) error {
											return c.Send("You need to send a video", selector)
										})

										b.Handle(tele.OnPhoto, func(c tele.Context) error {
											if c.Message().Photo == nil {
												return c.Send("Thumbnail is invalid! Please try again.", selector)
											}
											if c.Message().Photo.FileSize > 10*1024*1024 {
												return c.Send("Thumbnail must be less than 10MB! Please try again.", selector)
											}
											if c.Message().Photo.Width < c.Message().Photo.Height {
												return c.Send("Thumbnail must be landscape! Please try again.", selector)
											}

											thumbnail, err := b.FileByID(c.Message().Photo.FileID)
											if err != nil {
												return c.Send("File not found! Please try again.", selector)
											}

											inlineButtons := &tele.ReplyMarkup{}
											btnYes := inlineButtons.Data("✅ Yes", "Yes")
											btnNo := inlineButtons.Data("❌ No", "No")

											inlineButtons.Inline(inlineButtons.Row(btnYes, btnNo))

											b.Handle(&btnNo, func(c tele.Context) error {
												generated_thumbnail := &tele.Photo{File: thumbnail, Caption: "Description: " + description}
												photoInlineButtons := &tele.ReplyMarkup{}
												btnPublish := inlineButtons.Data("🌐 Publish", "Publish")
												btnWatch := inlineButtons.URL("🍿 Watch", os.Getenv("BOT_URL")+"?start="+video.FileID[15:])
												btnTakeDown := inlineButtons.URL("🔥 Take Down The Post", os.Getenv("BOT_URL")+"?start=del"+video.FileID[20:])
												btnCheckChannel := inlineButtons.URL("🔎 Check", channel.InviteLink)
												btnCheck := inlineButtons.URL("🔎 Check", os.Getenv("BOT_URL"))
												photoInlineButtons.Inline(photoInlineButtons.Row(btnPublish, btnCheckChannel), photoInlineButtons.Row(btnTakeDown))

												msg, err := b.Send(c.Sender(), generated_thumbnail, photoInlineButtons)
												if err != nil {
													return c.Send("Cannnot send thumbnail! Please try again.", selector)
												}

												b.Handle(&btnPublish, func(c tele.Context) error {
													status, err := b.Send(c.Sender(), "Publishing...")
													if err != nil {
														return c.Send("Cannnot set status! Please try again.", selector)
													}
													c.Send(createPost(id, title, rating, description, tags, video.FileID, video.FilePath, thumbnail.FileID, thumbnail.FilePath, "", ""))
													photoInlineButtons.Inline(photoInlineButtons.Row(btnWatch, btnCheck))
													b.Send(channel, generated_thumbnail, photoInlineButtons)
													b.Edit(status, "And the post is sent to "+channel.Title+"!")
													return nil
												})

												c.Send(msg.Photo.FilePath)
												return os.Remove("tmp/thumbnail.png")
											})

											b.Handle(&btnYes, func(c tele.Context) error {
												buf := screenshot(os.Getenv("APP_URL")+"/thumbnail?title="+title+"&subtitle="+"Rating+%3A+"+strconv.Itoa(rating)+"%25&categories="+tags+"&imgSrc="+os.Getenv("TG_API_URL")+"/file/bot"+os.Getenv("BOT_TOKEN")+"/"+thumbnail.FilePath, `#thumbnail-container`)
												err := os.WriteFile("tmp/thumbnail.png", buf, 0644)
												if err != nil {
													return c.Send("Cannnot create temp file! Please try again.", selector)
												}
												generated_thumbnail := &tele.Photo{File: tele.FromDisk("tmp/thumbnail.png"), Caption: "Description: " + description}
												photoInlineButtons := &tele.ReplyMarkup{}
												btnPublish := inlineButtons.Data("🌐 Publish", "Publish")
												btnWatch := inlineButtons.URL("🍿 Watch", os.Getenv("BOT_URL")+"?start="+video.FileID[15:])
												btnTakeDown := inlineButtons.URL("🔥 Take Down The Post", os.Getenv("BOT_URL")+"?start=del%20"+video.FileID[20:])
												btnCheckChannel := inlineButtons.URL("🔎 Check", channel.InviteLink)
												btnCheck := inlineButtons.URL("🔎 Check", os.Getenv("BOT_URL"))
												photoInlineButtons.Inline(photoInlineButtons.Row(btnPublish, btnCheckChannel), photoInlineButtons.Row(btnTakeDown))

												msg, err := b.Send(c.Sender(), generated_thumbnail, photoInlineButtons)
												if err != nil {
													return c.Send("Cannnot send thumbnail! Please try again.", selector)
												}

												b.Handle(&btnPublish, func(c tele.Context) error {
													status, err := b.Send(c.Sender(), "Publishing...")
													if err != nil {
														return c.Send("Cannnot set status! Please try again.", selector)
													}

													g_thumbnail, err := b.FileByID(msg.Photo.FileID)
													if err != nil {
														return c.Send("Cannnot fetch generated thumbnail! Please try again.", selector)
													}

													c.Send(createPost(id, title, rating, description, tags, video.FileID, video.FilePath, thumbnail.FileID, thumbnail.FilePath, g_thumbnail.FileID, g_thumbnail.FilePath))
													photoInlineButtons.Inline(photoInlineButtons.Row(btnWatch, btnCheck))
													b.Send(channel, generated_thumbnail, photoInlineButtons)
													b.Edit(status, "And the post is sent to "+channel.Title+"!")
													return nil
												})

												c.Send(msg.Photo.FilePath)
												return os.Remove("tmp/thumbnail.png")
											})

											return c.Send("Do you want to modify thumbnail with video details?", inlineButtons)
										})

										b.Handle(tele.OnVideo, func(c tele.Context) error {
											if c.Message().Video == nil {
												return c.Send("Thumbnail is invalid", selector)
											}
											if c.Message().Video.FileSize > 10*1024*1024 {
												return c.Send("Thumbnail must be less than 10MB! Please try again.", selector)
											}
											if c.Message().Video.Width < c.Message().Video.Height {
												return c.Send("Thumbnail must be landscape", selector)
											}
											if c.Message().Video.Duration > 10 {
												return c.Send("Thumbnail must not be longer than 10 seconds", selector)
											}

											thumbnail := c.Message().Video
											// create post with thumbnail
											b.Handle(tele.OnVideo, func(c tele.Context) error {
												if c.Message().Video == nil {
													return c.Send("Video is invalid")
												}
												return c.Send("Creating post with following data.\nTitle: " + title + "\nDescription: " + description + "\nTags: " + tags + "\nThumbnail: " + thumbnail.FileID + "\nVideo: " + c.Message().Video.FileID)
											})

											return c.Send("Please send the video.")
										})

										return c.Send("Please send a thumbnail for the video.")
									})
									return c.Send("Please send the video.")
								})
								return c.Send("Please send the video tags. 	(e.g. \"tag1, tag2, tag3\")")
							})

							return c.Send("Please send the video description.")
						})

						return c.Send("Please send the rating. (between 1 to 100)", selector)
					})
					return c.Send("Please send the post title.", selector)
				})

				menu.Reply(menu.Row(btnCreateThumbnail), menu.Row(btnCreatePost, btnSettings))

				return c.Send("Logged in successfully!", menu)
			})

			return c.Send("Please send the password to login.")
		})
		return c.Send("Please send the channel id to login.")
	}
}
