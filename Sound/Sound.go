package Sound

import (
	"log"
	"os"
	"time"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/wav"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/effects"
	"Chip8/CPU"
	"Chip8/Global"
)

var (
	Beep_buffer *beep.Buffer
)

func Initialize(file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10000))

	Beep_buffer = beep.NewBuffer(format)
	Beep_buffer.Append(streamer)

	streamer.Close()

	shot := Beep_buffer.Streamer(0, Beep_buffer.Len())

	loop := beep.Loop(20, shot)

	ctrl := &beep.Ctrl{Streamer: loop, Paused: false}
	// speaker.Play(ctrl)
	volume := &effects.Volume{
		Streamer: ctrl,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}
	// speaker.Play(volume)
	speedy := beep.ResampleRatio(4, 1, volume)

	// PLAY
	speaker.Play(speedy)


	// Sound Infinite Loop
	for {

		// Necessary to dont hang
		fmt.Sprint("")

		if Global.PlaySound {
			// Avoid multiple sound starts and stops while playing
			if !Global.SpeakerPlaying {
				speaker.Lock()
					ctrl.Paused = false
					// volume.Volume += 0.5
					// speedy.SetRatio(speedy.Ratio() + 0.1) // <-- right here
					// fmt.Println(format.SampleRate.D(shot.Position()).Round(time.Second))

					// newPos := shot.Position()
					// newPos = 0
					// shot.Seek(newPos)
				speaker.Unlock()

				if CPU.Debug {
					fmt.Print("Start playing sound\n")
				}
				Global.SpeakerPlaying = true
				Global.SpeakerStopped = false
			}

		} else {
			// Flag used to avoid multiple sound starts and stops while playing
			if !Global.SpeakerStopped {
				speaker.Lock()
					ctrl.Paused = true
					// volume.Volume += 0.5
					// speedy.SetRatio(speedy.Ratio() + 0.1) // <-- right here
					// fmt.Println(format.SampleRate.D(shot.Position()).Round(time.Second))

					newPos := shot.Position()
					newPos = 0
					shot.Seek(newPos)
				speaker.Unlock()

				if CPU.Debug {
					fmt.Print("Stop playing sound\n")
				}
				Global.SpeakerPlaying = false
				Global.SpeakerStopped = true
			}
		}
	}
}
