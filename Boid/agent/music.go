package agent

import (
	"fmt"
	"math"
	"os"
	"time"
	"unicode"

	pkg "github.com/DylanMeeus/GoAudio/wave"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/gdamore/tcell"
)

type MusicAgent struct {
	mp3FileTitle string // pour lire la musique et la controler via l'utilisation de Beep
	wavFileTite  string // pour obtenir l'amplitude du signal sonore de la musique via l'utilisation de GoAudio/wave (implique le besoin d'ajouter le fichier de musique en version .wav en plus d'une version .mp3) 
	sync         chan string // pour comuniquer les variations d'amplitude au jeu (game) afin que les comportements des poisons puissent être modifier en fonction de la musique
}

func NewMusicAgent(mp3FileTitle string, wavFileTite string, sync chan string) *MusicAgent {
	return &MusicAgent{mp3FileTitle, wavFileTite, sync}
}

func (musicAgent *MusicAgent) Start() {
	go func() {
		// lecture du fichier mp3 pour lire et lancer la musique
		f, err := os.Open(musicAgent.mp3FileTitle)
		if err != nil {
			report(err)
		}
		streamer, format, err := mp3.Decode(f)
		if err != nil {
			report(err)
		}
		defer streamer.Close()

		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/30))

		screen, err := tcell.NewScreen()
		if err != nil {
			report(err)
		}
		err = screen.Init()
		if err != nil {
			report(err)
		}
		defer screen.Fini()

		// récupération de l'amplitude du signal audio de la musique au travers du fichie.wav
		wave, err := pkg.ReadWaveFile(musicAgent.wavFileTite)
		// wave de type Wave possède comme attribut: un tableau de Frames d’indice i représentant l'amplitude du signal provenant du fichier .wav pour chaque échantillon i du signal.
		if err != nil {
			panic("Could not parse wave file")
		}
		fmt.Printf("Read %v samples\n", len(wave.Frames))

		ap := newAudioPanel(format.SampleRate, streamer, wave.Frames, musicAgent.sync)

		screen.Clear()
		ap.draw(screen)
		screen.Show()

		ap.play()

		Seconds := time.NewTicker(time.Second)
		MicroSeconds := time.NewTicker(time.Microsecond)
		events := make(chan tcell.Event)

		go func() {
			for {
				events <- screen.PollEvent()
			}
		}()

		// loop:
		esc := false
		for !esc {
			select {
			case event := <-events:
				changed, quit := ap.handle(event)
				if quit {
					esc = true
				}

				if changed {
					screen.Clear()
					ap.draw(screen)
					screen.Show()
				}
			case <-Seconds.C:
				screen.Clear()
				ap.draw(screen)
				screen.Show()

			//Observation de la musique toute les microsecondes:
			case <-MicroSeconds.C:
				speaker.Lock()
				// Grace à l’attribut streamer de notre audioPanel (voir suite du code) donnant l’échantillon i du signal de la musique en cours de lecture, pouvons récupérer l’amplitude du signal audio toutes les microsecondes et le stocker dans amplitude
				sr := ap.streamer.Position()
				amplitude := ap.samples[sr]
				speaker.Unlock()

				// communication avec game en fonction des variations d'amplitude du signal sonore:
				if amplitude > 0.5 {
					ap.sync <- "very hard drop"
				} else if amplitude > 0.3 {
					ap.sync <- "hard drop"
				} else if sr > 0 && math.Sqrt(math.Pow(float64(amplitude-ap.samples[sr-1]), 2)) > 0.5 {
					ap.sync <- "medium drop"
				} else if sr > 0 && math.Sqrt(math.Pow(float64(amplitude-ap.samples[sr-1]), 2)) > 0.3 {
					ap.sync <- "small drop"
				} else if sr > 0 && math.Sqrt(math.Pow(float64(amplitude-ap.samples[sr-1]), 2)) > 0.2 {
					ap.sync <- "small drop"
				} else {
					ap.sync <- "R"
				}
			}
		}
	}()
}

func report(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func drawTextLine(screen tcell.Screen, x, y int, s string, style tcell.Style) {
	for _, r := range s {
		screen.SetContent(x, y, r, nil, style)
		x++
	}
}

// Comme vu plus haut Beep nous a permis de  lire un fichier mp3 via Go au travers d’un agent musicale: MusicAgent. En plus de cela Beep nous a permis de générer dans la fonction Start de l'agent musical: un audioPanel à l’aide du tutoriel présent à l’adresse suivante: https://github.com/faiface/beep/wiki.
// Ainsi nous pouvions modifier la position, le volume et la vitesse de la musique en lecture.

type audioPanel struct {
	sampleRate beep.SampleRate
	streamer   beep.StreamSeeker
	ctrl       *beep.Ctrl
	resampler  *beep.Resampler
	volume     *effects.Volume
	samples    []pkg.Frame
	sync       chan string
}

func newAudioPanel(sampleRate beep.SampleRate, streamer beep.StreamSeeker, samples []pkg.Frame, c1 chan string) *audioPanel {
	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer)}
	resampler := beep.ResampleRatio(4, 1, ctrl)
	volume := &effects.Volume{Streamer: resampler, Base: 2}
	return &audioPanel{sampleRate, streamer, ctrl, resampler, volume, samples, c1}
}

func (ap *audioPanel) play() {
	speaker.Play(ap.volume)
}

func (ap *audioPanel) draw(screen tcell.Screen) {
	mainStyle := tcell.StyleDefault.
		Background(tcell.NewHexColor(0x473437)).
		Foreground(tcell.NewHexColor(0xD7D8A2))
	statusStyle := mainStyle.
		Foreground(tcell.NewHexColor(0xDDC074)).
		Bold(true)

	screen.Fill(' ', mainStyle)

	drawTextLine(screen, 0, 0, "Welcome to the Music Video Player!", mainStyle)
	drawTextLine(screen, 0, 1, "Press [ESC] to quit.", mainStyle)
	drawTextLine(screen, 0, 2, "Press [SPACE] to pause/resume.", mainStyle)
	drawTextLine(screen, 0, 3, "Use keys in (?/?) to turn the buttons.", mainStyle)

	speaker.Lock()
	sr := ap.streamer.Position()
	amplitude := ap.samples[sr]
	position := ap.sampleRate.D(ap.streamer.Position())
	length := ap.sampleRate.D(ap.streamer.Len())
	volume := ap.volume.Volume
	speed := ap.resampler.Ratio()
	speaker.Unlock()

	amplitudeStatus := fmt.Sprintf("%f", amplitude)
	positionStatus := fmt.Sprintf("%v / %v", position.Round(time.Second), length.Round(time.Second))
	volumeStatus := fmt.Sprintf("%.1f", volume)
	speedStatus := fmt.Sprintf("%.3fx", speed)

	drawTextLine(screen, 0, 5, "Position (Q/W):", mainStyle)
	drawTextLine(screen, 16, 5, positionStatus, statusStyle)

	drawTextLine(screen, 0, 6, "Volume   (A/S):", mainStyle)
	drawTextLine(screen, 16, 6, volumeStatus, statusStyle)

	drawTextLine(screen, 0, 7, "Speed    (Z/X):", mainStyle)
	drawTextLine(screen, 16, 7, speedStatus, statusStyle)

	drawTextLine(screen, 0, 8, "amplitude:", mainStyle)
	drawTextLine(screen, 16, 8, amplitudeStatus, statusStyle)

	// Si l'écran s'affiche toutes les microsecondes le code qui suit permet de voir ce qui est envoyer aux slimes:
	// if amplitude > 0.5 {
	// 	drawTextLine(screen, 0, 9, "amplitude BooM:", mainStyle)
	// } else if amplitude > 0.3 {
	// 	drawTextLine(screen, 0, 9, "amplitude BooM:", mainStyle)
	// } else if sr > 0 && math.Sqrt(math.Pow(float64(amplitude-ap.samples[sr-1]), 2)) > 0.5 {
	// 	drawTextLine(screen, 0, 9, "amplitude BooM:", mainStyle)
	// } else if sr > 0 && math.Sqrt(math.Pow(float64(amplitude-ap.samples[sr-1]), 2)) > 0.3 {
	// 	drawTextLine(screen, 0, 9, "amplitude BOOM:", mainStyle)
	// } else if sr > 0 && math.Sqrt(math.Pow(float64(amplitude-ap.samples[sr-1]), 2)) > 0.2 {
	// 	drawTextLine(screen, 0, 9, "amplitude boom:", mainStyle)
	// } else {
	// 	drawTextLine(screen, 0, 9, "R boom:", mainStyle)
	// }
}

func (ap *audioPanel) handle(event tcell.Event) (changed, quit bool) {
	switch event := event.(type) {
	case *tcell.EventKey:
		if event.Key() == tcell.KeyESC {
			return false, true
		}

		if event.Key() != tcell.KeyRune {
			return false, false
		}

		switch unicode.ToLower(event.Rune()) {
		case ' ':
			speaker.Lock()
			ap.ctrl.Paused = !ap.ctrl.Paused
			speaker.Unlock()
			return false, false

		case 'q', 'w':
			speaker.Lock()
			newPos := ap.streamer.Position()
			if event.Rune() == 'q' {
				newPos -= ap.sampleRate.N(time.Second)
			}
			if event.Rune() == 'w' {
				newPos += ap.sampleRate.N(time.Second)
			}
			if newPos < 0 {
				newPos = 0
			}
			if newPos >= ap.streamer.Len() {
				newPos = ap.streamer.Len() - 1
			}
			if err := ap.streamer.Seek(newPos); err != nil {
				report(err)
			}
			speaker.Unlock()
			return true, false

		case 'a':
			speaker.Lock()
			ap.volume.Volume -= 0.1
			speaker.Unlock()
			return true, false

		case 's':
			speaker.Lock()
			ap.volume.Volume += 0.1
			speaker.Unlock()
			return true, false

		case 'z':
			speaker.Lock()
			ap.resampler.SetRatio(ap.resampler.Ratio() * 15 / 16)
			speaker.Unlock()
			return true, false

		case 'x':
			speaker.Lock()
			ap.resampler.SetRatio(ap.resampler.Ratio() * 16 / 15)
			speaker.Unlock()
			return true, false

		// raccourcis secret par la channel de la musique:
		case 'l':
			// commenter lock/unlock si on ne veux pas que la musique s'arrete
			speaker.Lock()
			ap.sync <- "1"
			time.Sleep((1 * time.Second))
			speaker.Unlock()
			return true, false
		case 'o':
			speaker.Lock()
			ap.sync <- "2"
			time.Sleep((1 * time.Second))
			speaker.Unlock()
			return true, false
		case 'u':
			speaker.Lock()
			ap.sync <- "3"
			time.Sleep((1 * time.Second))
			speaker.Unlock()
			return true, false
		case 'i':
			speaker.Lock()
			ap.sync <- "4"
			time.Sleep((1 * time.Second))
			speaker.Unlock()
			return true, false
		case 'j':
			speaker.Lock()
			ap.sync <- "5"
			time.Sleep((1 * time.Second))
			speaker.Unlock()
			return true, false
		}
	}
	return false, false
}
