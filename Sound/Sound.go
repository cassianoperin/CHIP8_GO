package Sound

import (
	"log"
	"os"
	"time"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var (
	Beep_buffer *beep.Buffer
)

func Initialize(file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/30000))

	Beep_buffer = beep.NewBuffer(format)
	Beep_buffer.Append(streamer)

	streamer.Close()
}


func PlaySound(sound *beep.Buffer){
	shot := sound.Streamer(0, sound.Len())
	done := make(chan bool)
	speaker.Play(beep.Seq(shot, beep.Callback(func() {
		done <- true
	})))
	<-done
}
