package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type SoundRequest struct {
	MP3Link string `json:"mp3Link"`
}

func playSound(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var t SoundRequest
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()

		// Download mp3 file
		resp, err := http.Get(t.MP3Link)
		if err != nil {
			log.Fatalln(err)
		}

		// Read mp3 file into byte array
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		resp.Body.Close()

		// Decode mp3 data
		s, format, err := mp3.Decode(ioutil.NopCloser(bytes.NewReader(body)))
		if err != nil {
			http.Error(w, "Could not fetch file", http.StatusBadRequest)
			return
		}
		defer s.Close()

		// Initialize speaker with sample rate
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

		// Play the sound
		speaker.Play(beep.Seq(s, beep.Callback(func() {
			fmt.Println("Sound finished playing")
		})))

		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/play-sound", playSound)
	log.Fatal(http.ListenAndServe(":44342", nil))
}
