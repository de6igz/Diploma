package discord_repository

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
)

// DiscordRepository описывает интерфейс для отправки сообщений в Discord.
type DiscordRepository interface {
	SendMessage(channelID, message string) error
}

// discordRepository реализует DiscordRepository.
type discordRepository struct {
	session *discordgo.Session
	logger  *zerolog.Logger
}

// NewDiscordRepository создаёт экземпляр discordRepository.
func NewDiscordRepository(logger *zerolog.Logger, botToken string) (DiscordRepository, error) {
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		return nil, err
	}

	// Открываем websocket-соединение с Discord
	err = dg.Open()
	if err != nil {
		return nil, err
	}
	logger.Info().Msg("Discord session opened successfully")
	return &discordRepository{
		session: dg,
		logger:  logger,
	}, nil
}

// SendMessage отправляет сообщение в указанный канал Discord.
func (r *discordRepository) SendMessage(channelID, message string) error {
	_, err := r.session.ChannelMessageSend(channelID, message)
	if err != nil {
		r.logger.Error().Err(err).Msgf("Failed to send Discord message to channel %s", channelID)
		return err
	}
	r.logger.Info().Msgf("Successfully sent Discord message to channel %s", channelID)
	return nil
}
