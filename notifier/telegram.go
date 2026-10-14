package notifier

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"depop-alerter/depop"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramNotifier struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

func NewTelegramNotifier(token string, chatID int64) (*TelegramNotifier, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}

	return &TelegramNotifier{
		bot:    bot,
		chatID: chatID,
	}, nil
}

func (n *TelegramNotifier) NotifyNewProduct(query string, p depop.Product) error {
	var caption bytes.Buffer
	caption.WriteString(fmt.Sprintf("🚨 <b>New Match for:</b> %s\n\n", query))
	
	// Handle Title/Description (sometimes Title is named differently or we just use slug)
	title := p.Slug
	if p.Title != nil {
		title = *p.Title
	}
	// Truncate title if it's too long
	if len(title) > 100 {
		title = title[:100] + "..."
	}
	caption.WriteString(fmt.Sprintf("<b>Item:</b> %s\n", tgbotapi.EscapeText(tgbotapi.ModeHTML, title)))

	price := "N/A"
	if p.Price.PriceString != "" {
		price = p.Price.PriceString
	} else if p.Price.Amount != "" {
		price = fmt.Sprintf("%s %s", p.Price.Amount, p.Price.CurrencyCode)
	}
	caption.WriteString(fmt.Sprintf("<b>Price:</b> %s\n", tgbotapi.EscapeText(tgbotapi.ModeHTML, price)))

	if len(p.Sizes) > 0 {
		caption.WriteString(fmt.Sprintf("<b>Size:</b> %s\n", tgbotapi.EscapeText(tgbotapi.ModeHTML, strings.Join(p.Sizes, ", "))))
	}

	// Figure out seller username
	username := p.SlugUsername
	if username == "" && p.Seller.Username != "" {
		username = p.Seller.Username
	}
	if username != "" {
		caption.WriteString(fmt.Sprintf("<b>Seller:</b> <a href=\"https://www.depop.com/%s/\">%s</a>\n", username, username))
	}

	itemURL := fmt.Sprintf("https://www.depop.com/products/%s/", p.Slug)
	caption.WriteString(fmt.Sprintf("<b>Link:</b> <a href=\"%s\">View Item</a>\n", itemURL))

	// Get first image
	var imageURL string
	if len(p.Pictures) > 0 && len(p.Pictures[0]) > 0 {
		imageURL = p.Pictures[0][0].URL
	} else if len(p.Images) > 0 {
		imageURL = p.Images[0].URL
	}

	var msg tgbotapi.Chattable
	if imageURL != "" {
		photoMsg := tgbotapi.NewPhoto(n.chatID, tgbotapi.FileURL(imageURL))
		photoMsg.Caption = caption.String()
		photoMsg.ParseMode = tgbotapi.ModeHTML
		msg = photoMsg
	} else {
		textMsg := tgbotapi.NewMessage(n.chatID, caption.String())
		textMsg.ParseMode = tgbotapi.ModeHTML
		textMsg.DisableWebPagePreview = true
		msg = textMsg
	}

	_, err := n.bot.Send(msg)
	return err
}

func (n *TelegramNotifier) NotifyError(errMsg string) {
	text := fmt.Sprintf("⚠️ <b>Depop Alerter Error:</b>\n\n<code>%s</code>", tgbotapi.EscapeText(tgbotapi.ModeHTML, errMsg))
	msg := tgbotapi.NewMessage(n.chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	
	if _, err := n.bot.Send(msg); err != nil {
		log.Printf("Failed to send error notification: %v", err)
	}
}
