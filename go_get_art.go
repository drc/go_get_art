package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

var SPOTIFY_ID, SPOTIFY_SECRET string

func main() {

	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		PaddingTop(2).
		PaddingLeft(4).
		PaddingRight(4).
		PaddingBottom(2)

	var args struct {
		SongID string `arg:"positional,required" help:"The Spotify URI"`
	}

	arg.MustParse(&args)

	if SPOTIFY_ID == "" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		SPOTIFY_ID = os.Getenv("SPOTIFY_ID")
		SPOTIFY_SECRET = os.Getenv("SPOTIFY_SECRET")
	}

	song_id := strings.Split(args.SongID, ":")[2]

	ctx := context.Background()
	config := &clientcredentials.Config{
		ClientID:     SPOTIFY_ID,
		ClientSecret: SPOTIFY_SECRET,
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)
	msg, err := client.GetTrack(ctx, spotify.ID(song_id))
	if err != nil {
		log.Fatalf("couldn't get track info: %v", err)
	}

	fmt.Println(style.Render(msg.Name, " ", msg.Album.ReleaseDate))

	fmt.Println(msg.Album.Images[0].URL)

	resp, err := client.GetArtist(ctx, msg.Artists[0].ID)

	if err != nil {
		log.Fatalf("couldnt get artist info: %v", err)
	}

	for _, genre := range resp.Genres {
		fmt.Println(" ", genre)
	}

}
