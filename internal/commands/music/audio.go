package music

import (
	"io"
	"log/slog"
	"sync"

	"github.com/Dahl99/discordbot/internal/discord"

	"github.com/bwmarrin/discordgo"
	"github.com/jung-m/dca"
)

type VoiceInstance struct {
	voice      *discordgo.VoiceConnection
	session    *discordgo.Session
	encoder    *dca.EncodeSession
	stream     *dca.StreamingSession
	queueMutex sync.Mutex
	nowPlaying Song
	queue      []Song
	GuildID    string
	speaking   bool
	pause      bool
	stop       bool
	skip       bool
}

type Song struct {
	ChannelID string
	UserName  string
	ID        string
	VideoID   string
	Title     string
	VideoURL  string
}

type PkgSong struct {
	data Song
	v    *VoiceInstance
}

var (
	VoiceInstances = map[string]*VoiceInstance{}
	mutex          sync.Mutex
	SongSignal     chan PkgSong
)

func globalPlay(songSig chan PkgSong) {
	for {
		select {
		case song := <-songSig:
			go song.v.PlayQueue(song.data)
		}
	}
}

// Skip skips the current playing song in the queue.
func (v *VoiceInstance) Skip() bool {
	if v.speaking {
		if v.pause {
			return true
		} else {
			if v.encoder != nil {
				v.encoder.Cleanup()
			}
		}
	}
	return false
}

// Stop stops the audio.
func (v *VoiceInstance) Stop() {
	v.stop = true
	if v.encoder != nil {
		v.encoder.Cleanup()
	}
}

// QueueAdd adds a song to the queue.
func (v *VoiceInstance) QueueAdd(song Song) {
	v.queueMutex.Lock()
	defer v.queueMutex.Unlock()
	slog.Info("adding song to queue", "song", song)
	v.queue = append(v.queue, song)
}

// QueueGetSong gets the first song from the queue.
func (v *VoiceInstance) QueueGetSong() (song Song) {
	v.queueMutex.Lock()
	defer v.queueMutex.Unlock()
	if len(v.queue) != 0 {
		slog.Info("Getting first song from queue")
		return v.queue[0]
	}
	return
}

// QueueRemoveFirst removes the first song from the queue.
func (v *VoiceInstance) QueueRemoveFirst() {
	v.queueMutex.Lock()
	defer v.queueMutex.Unlock()
	if len(v.queue) != 0 {
		v.queue = v.queue[1:]
	}
}

// QueueClear clears the entire queue.
func (v *VoiceInstance) QueueClear() {
	v.queueMutex.Lock()
	defer v.queueMutex.Unlock()
	v.queue = []Song{}
}

// DCA streams the song to the Discord voice channel.
func (v *VoiceInstance) DCA(url string) {
	opts := dca.StdEncodeOptions
	opts.RawOutput = true
	opts.Bitrate = 64
	opts.Application = "lowdelay"

	encodeSession, err := dca.EncodeFile(url, opts)
	if err != nil {
		slog.Error("Failed to create an encoding session", "error", err)
	}
	defer encodeSession.Cleanup()

	v.encoder = encodeSession
	done := make(chan error)
	v.stream = dca.NewStream(encodeSession, v.voice, done)
	for {
		select {
		case err := <-done:
			if err != nil && err != io.EOF {
				slog.Error("An error occured", "error", err)
			}

			return
		}
	}
}

func (v *VoiceInstance) PlayQueue(song Song) {
	// add song to queue
	v.QueueAdd(song)

	if v.speaking {
		// the bot is playing
		return
	}

	go func() {
		for {
			if len(v.queue) == 0 {
				slog.Info("song queue is empty")
				discord.SendChannelMessage(v.nowPlaying.ChannelID, "**[Music]** End of queue")
				return
			}

			v.nowPlaying = v.QueueGetSong()
			go discord.SendChannelMessage(v.nowPlaying.ChannelID, "**[Music]** Now playing: **"+v.nowPlaying.Title+"**")

			v.stop = false
			v.skip = false
			v.speaking = true
			v.pause = false
			err := v.voice.Speaking(true)
			if err != nil {
				slog.Error("Failed to send speaking notification", "error", err)
				return
			}

			v.DCA(v.nowPlaying.VideoURL)

			v.QueueRemoveFirst()
			if v.stop {
				v.QueueClear()
			}

			v.stop = false
			v.skip = false
			v.speaking = false

			err = v.voice.Speaking(false)
			if err != nil {
				slog.Error("Failed to stop sending speaking notification", "error", err)
			}
		}
	}()
}
