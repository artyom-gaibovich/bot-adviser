package telegram

import (
	"bot-adviser/clients/telegram"
	"bot-adviser/lib/e"
	"bot-adviser/storage"
	"errors"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)
	log.Printf("got new command '%s' from '%s", text, username)
	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}
	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessages(chatID, msgUnknownCommand)

	}

}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessages(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessages(chatID, msgHello)
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()
	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMessages(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := p.tg.SendMessages(chatID, msgSaved); err != nil {
		return err
	}

	return nil

}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't random", err) }()
	sendMessage := createSendMessageFunc(chatID, p.tg)

	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessages(chatID, msgNoSavedPages)
	}
	if err := sendMessage(page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)

}

func createSendMessageFunc(chatID int, tg *telegram.Client) func(msg string) error {

	return func(msg string) error {
		return tg.SendMessages(chatID, msg)
	}
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
