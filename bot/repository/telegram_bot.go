package repository

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/drhidians/testbot/bot"
	"github.com/drhidians/testbot/middleware/jwtauth"
	"github.com/drhidians/testbot/models"
	tg "github.com/drhidians/testbot/telegram"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type botRepository struct {
	botAPI    tg.Bot
	tokenAuth *jwtauth.JWTAuth
	domainurl string
}

// NewBotRepository will create an object that represent the bot.Repository interface
func NewBotRepository(botAPI tg.Bot, tokenAuth *jwtauth.JWTAuth, domainurl string) bot.Repository {

	br := botRepository{botAPI, tokenAuth, domainurl}
	//TO DO move to bot_commands.go, add extra layer (?)
	cmd := tg.NewCommands(botAPI.Username())
	cmd.Add("/start", br.sendToken())
	callCommand = func(u *tg.Update) error {
		if err, _ := cmd.Run(u); err != nil {
			return err
		}
		return nil
	}
	//====================================================

	return &br
}

func (b *botRepository) Get(ctx context.Context) (bot *models.Bot, err error) {

	bot = new(models.Bot)

	respUser, err := b.botAPI.GetMe(ctx)
	if err != nil {
		return nil, err
	}

	bot.ID = respUser.ID
	bot.Name = respUser.FirstName
	bot.Username = respUser.Username
	return bot, nil
}

func (b *botRepository) Update(ctx context.Context, upd tg.Update) (user *models.User, err error) {

	user = new(models.User)

	user.CreatedAt = time.Now()
	user.ExternalID = upd.Message.From.ID
	user.Language = upd.Message.From.LanguageCode
	user.Name = upd.Message.From.FirstName

	if upd.Message.From.LastName != nil {
		user.Name += " " + *upd.Message.From.LastName
	}
	user.Username = upd.Message.From.Username

	userAvatar, err := b.botAPI.GetUserProfilePhotos(ctx, tg.UserProfilePhotosConfig{
		UserID: user.ExternalID,
	})

	if err != nil {
		return nil, err
	}

	//TO DO to ugly
	if len(userAvatar.Photos) != 0 {
		avatarID := userAvatar.Photos[0][len(userAvatar.Photos[0])-1].FileID
		avatarURL := b.domainurl + "media/" + avatarID
		user.Avatar = &avatarURL
	}

	/*userFile, err := b.botAPI.GetFile(ctx, tg.FileConfig{
		FileID: avatarID,
	})

	if err != nil {
		return err
	}

	err := ioutil.WriteFile("/media/"+avatarID, d1, 0644)
	downloadFile("media/" + avatarID, b.botAPI.)
	if err != nil {
		return err
	}*/

	err = callCommand(&upd)

	return user, err
}

func (b *botRepository) GetFile(ctx context.Context, fileID string) (file []byte, err error) {

	fileURL, err := b.botAPI.GetFileDirectURL(ctx, tg.FileConfig{
		FileID: fileID,
	})

	if err != nil {
		return nil, err
	}

	file, err = getFile(fileURL)

	return file, err
}

func getFile(url *string) (bytes []byte, err error) {

	if url == nil {
		return nil, fmt.Errorf("url is  null")
	}
	// Get the data
	resp, err := http.Get(*url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	file, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	return file, nil
}
