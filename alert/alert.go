package alert

import (
	"fmt"
	"sync"

	"github.com/asticode/go-texttospeech/texttospeech"
)

var once sync.Once
var tts texttospeech.TextToSpeech

func Say(msg string) {
	once.Do(func() {
		tts = texttospeech.NewTextToSpeech()
	})

	fmt.Printf("ALERT! %s\n", msg)

	go tts.Say(msg)
}

func Sayf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)

	Say(msg)
}
