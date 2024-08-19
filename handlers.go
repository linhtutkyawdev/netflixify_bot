package main

import (
	"strconv"
	"time"

	tele "gopkg.in/telebot.v3"
)

func channelHandlers(b *tele.Bot) tele.HandlerFunc {
	return func(c tele.Context) error {
		if c.Text() == "/id" {
			b.Edit(c.Message(), strconv.FormatInt(c.Chat().ID, 10))
			time.AfterFunc(3*time.Second, func() {
				b.Delete(c.Message())
			})
			// if hasRegistered(c.Chat().ID) {
			// 	b.Edit(c.Message(), "You have already registered this channel!")
			// 	time.AfterFunc(3*time.Second, func() {
			// 		b.Delete(c.Message())
			// 	})
			// 	return nil
			// }
			// keyboard := &tele.ReplyMarkup{}

			// b.Send(c.Sender(), "Registering channel...")
			// btnYes := keyboard.URL("✅", "https://t.me/netflixify_bot?start="+strconv.FormatInt(c.Chat().ID, 10))
			// b.Handle(&btnYes, func(c tele.Context) error {
			// 	// return a button linked to the bot
			// 	// registerChannel(c.Chat().ID, c.Chat().Title)
			// 	b.Edit(c.Message(), "Your channel is registered successfully!")
			// 	time.AfterFunc(3*time.Second, func() {
			// 		b.Delete(c.Message())
			// 	})
			// 	return nil
			// })

			// btnNo := keyboard.Data("❌", "no")
			// b.Handle(&btnNo, func(c tele.Context) error {
			// 	return b.Delete(c.Message())
			// })

			// keyboard.Inline(keyboard.Row(btnYes, btnNo))

			// b.Edit(c.Message(), "Are you sure you want to register this channel for netflixifying?", &tele.SendOptions{ReplyMarkup: keyboard})
		}
		// // starts wirth "@netflixify_bot"
		// if c.Text()[0] == '@' {
		// 	keyboard := &tele.ReplyMarkup{}
		// 	btnCheckHere := keyboard.URL("Check out here", "https://t.me/netflixify_bot")
		// 	keyboard.Inline(keyboard.Row(btnCheckHere))
		// 	b.Edit(c.Message(), "I've personally sent you a message!", &tele.SendOptions{ReplyMarkup: keyboard})
		// 	time.AfterFunc(3*time.Second, func() {
		// 		b.Delete(c.Message())
		// 	})
		// }
		return nil
	}
}
