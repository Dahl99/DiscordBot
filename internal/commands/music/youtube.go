package music

import (
	"context"
	"errors"
	"log/slog"

	"google.golang.org/api/option"

	"github.com/Dahl99/discordbot/internal/config"
	"github.com/Dahl99/discordbot/internal/discord"

	"github.com/bwmarrin/discordgo"
	ytdl "github.com/kkdai/youtube/v2"

	"google.golang.org/api/youtube/v3"
)

type Video struct {
	VideoID    string
	VideoTitle string
}

func SearchVideoByName(name string) (*Video, error) {
	slog.Info("Searching video by name", "name", name)

	var err error
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(config.GetYoutubeApiKey()))
	if err != nil {
		slog.Warn("failed to create YoutubeService (youtube-api-v3)", "error", err)
		return nil, err
	}

	call := youtubeService.Search.
		List([]string{"id", "snippet"}).
		Q(name).
		MaxResults(1)

	res, err := call.Do()
	if err != nil {
		slog.Warn("failed to search YouTube for video", "error", err)
		return nil, err
	}

	var (
		videoID, videoTitle string
	)

	for _, item := range res.Items {
		videoID = item.Id.VideoId
		videoTitle = item.Snippet.Title
	}

	if videoID == "" {
		slog.Info("video not found on YouTube")
		return nil, errors.New("video not found on YouTube")
	}

	slog.Info("found video on YouTube", "videoID", videoID, "videoTitle", videoTitle)

	return &Video{
		VideoID:    videoID,
		VideoTitle: videoTitle,
	}, nil
}

func findVideoByVideoID(videoID string) (*Video, error) {
	slog.Info("finding video by ID", "videoID", videoID)

	ctx := context.Background()
	var err error
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(config.GetYoutubeApiKey()))
	if err != nil {
		slog.Warn("failed to create YoutubeService", "error", err)
		return nil, err
	}

	call := youtubeService.Videos.List([]string{"id", "snippet"}).Id(videoID).MaxResults(1)
	res, err := call.Do()
	if err != nil {
		slog.Warn("failed to find video on YouTube (youtube-api-v3)", "error", err)
		return nil, err
	}

	var videoTitle string

	for _, item := range res.Items {
		videoTitle = item.Snippet.Title
	}

	if videoID == "" {
		slog.Info("video not found on YouTube")
		return nil, errors.New("video not found on YouTube")
	}

	return &Video{
		VideoID:    videoID,
		VideoTitle: videoTitle,
	}, nil
}

func getSongDataByVideo(video *Video, v *VoiceInstance, m *discordgo.MessageCreate) (songStruct PkgSong, err error) {
	client := ytdl.Client{
		HTTPClient:  nil,
		MaxRoutines: 0,
		ChunkSize:   0,
	}

	youTubeVideo, err := client.GetVideo("https://www.youtube.com/watch?v=" + video.VideoID)
	if err != nil {
		slog.Warn("failed to get video from YouTube (ytdl-Client)", "error", err)
		return
	}

	youTubeVideo.Formats.Sort()
	formatURL := youTubeVideo.Formats[len(youTubeVideo.Formats)-1].URL

	guildID := discord.SearchGuildByChannelID(m.ChannelID)
	member, _ := v.session.GuildMember(guildID, m.Author.ID)
	userName := ""

	if member.Nick == "" {
		userName = m.Author.Username
	} else {
		userName = member.Nick
	}

	song := Song{
		m.ChannelID,
		userName,
		m.Author.ID,
		video.VideoID,
		video.VideoTitle,
		formatURL,
	}

	// var song_struct PkgSong
	songStruct.data = song
	songStruct.v = v

	return songStruct, nil
}
