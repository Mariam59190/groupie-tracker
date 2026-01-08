package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const artistsURL = "https://groupietrackers.herokuapp.com/api/artists"

type Artist struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

func fetchArtists() ([]Artist, error) {
	resp, err := http.Get(artistsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var artists []Artist
	err = json.Unmarshal(body, &artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Groupie Tracker")
	myWindow.Resize(fyne.NewSize(800, 600))

	title := widget.NewLabelWithStyle(
		"🎵 Groupie Tracker",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	content := container.NewVBox(title)

	artists, err := fetchArtists()
	if err != nil {
		content.Add(widget.NewLabel("Erreur lors du chargement de l'API"))
	} else {
		for _, artist := range artists {
			card := widget.NewCard(
				artist.Name,
				fmt.Sprintf("Création : %d | Premier album : %s",
					artist.CreationDate, artist.FirstAlbum),
				widget.NewLabel("Membres : "+fmt.Sprint(artist.Members)),
			)
			content.Add(card)
		}
	}

	scroll := container.NewVScroll(content)
	myWindow.SetContent(scroll)
	myWindow.ShowAndRun()
}
