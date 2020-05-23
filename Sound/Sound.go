package Sound

import (
	"log"
	"os"
	"time"
	"github.com/faiface/beep"
	"github.com/faiface/beep/wav"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/effects"
)

var (
	Shot			beep.StreamSeeker
	AudioCtrl		*beep.Ctrl
)


func AudioDaemonStart(file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/600))

	Beep_buffer := beep.NewBuffer(format)
	Beep_buffer.Append(streamer)

	streamer.Close()

	Shot = Beep_buffer.Streamer(0, Beep_buffer.Len())

	loop := beep.Loop(-1, Shot)

	AudioCtrl = &beep.Ctrl{Streamer: loop, Paused: false}
	// speaker.Play(ctrl)
	volume := &effects.Volume{
		Streamer: AudioCtrl,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}
	// speaker.Play(volume)
	speedy := beep.ResampleRatio(4, 1, volume)

	// Start Paused
	speaker.Lock()
	AudioCtrl.Paused = true
	speaker.Unlock()

	// PLAY
	speaker.Play(speedy)

}
