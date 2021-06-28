package main

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/uploader"
	"github.com/gotd/td/tg"
	"github.com/jessevdk/go-flags"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Options struct {
	AppID            int    `short:"a" long:"app-id" env:"APP_ID" required:"true" description:"telegram app id"`
	AppHash          string `short:"h" long:"app-hash" env:"APP_HASH" required:"true" description:"telegram app hash"`
	BotToken         string `short:"t" long:"bot-token" env:"BOT_TOKEN" required:"true" description:"telegram bot token"`
	UserID           int    `short:"u" long:"user-id" env:"USER_ID" required:"true" description:"the id of the telegram user to whom the bot will send the files"`
	InputFolderPath  string `short:"i" long:"input-folder" env:"INPUT_FOLDER" required:"true" description:"the folder where service get files for sending to telegram" default:"/data/input"`
	OutputFolderPath string `short:"o" long:"output-folder" env:"OUTPUT_FOLDER" required:"true" description:"output folder" default:"/data/output"`
}

func main() {
	opts := parseOpts()
	ctx := context.Background()
	filename := ""
	client := telegram.NewClient(opts.AppID, opts.AppHash, telegram.Options{})
	err := client.Run(ctx, func(ctx context.Context) error {

		if _, err := client.AuthBot(ctx, opts.BotToken); err != nil {
			return err
		}

		c := tg.NewClient(client)

		u := uploader.NewUploader(c)

		if !exists(opts.InputFolderPath) {
			log.Fatal("Input folder is not exist")
		}
		if !exists(opts.OutputFolderPath) {
			log.Fatal("Output folder is not exist")
		}

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Program has been started")
		defer watcher.Close()
		done := make(chan bool)
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					if event.Op&fsnotify.Create == fsnotify.Create {
						time.Sleep(10 * time.Second)
						filename = event.Name

						err = sendFile(ctx, opts.UserID, u, filename, *c)
						if err != nil {
							log.Println(err)
						}

						log.Printf("[info] move %s to output folder", event.Name)
						oldLocation := event.Name
						newLocation := filepath.Join(opts.OutputFolderPath, filepath.Base(event.Name))
						err = os.Rename(oldLocation, newLocation)
						if err != nil {
							log.Fatal(err)
						}

					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					log.Print("[error]", err)
				}
			}
		}()

		err = watcher.Add(opts.InputFolderPath)
		if err != nil {
			log.Fatal(err)
		}
		<-done

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func parseOpts() Options {
	var opts Options
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}
	return opts
}

func sendFile(ctx context.Context, userID int, u *uploader.Uploader, filename string, c tg.Client) error {
	rand.Seed(time.Now().Unix())
	ranId := rand.Int63()
	peer := tg.InputPeerUser{
		UserID:     userID,
		AccessHash: 0,
	}
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
	}
	mimeType := http.DetectContentType(fileBytes)
	log.Printf("[info] start uploading %s", filename)
	file, err := u.FromPath(ctx, filename)
	if err != nil {
		log.Println(err)
	}
	media := tg.InputMediaUploadedDocument{
		Flags:      0,
		ForceFile:  true,
		File:       file,
		MimeType:   mimeType,
		Attributes: []tg.DocumentAttributeClass{&tg.DocumentAttributeFilename{FileName: filepath.Base(filename)}},
		TTLSeconds: 5000,
	}
	rq := tg.MessagesSendMediaRequest{
		Peer:     &peer,
		Media:    &media,
		RandomID: ranId,
	}

	_, err = c.MessagesSendMedia(ctx, &rq)
	if err != nil {
		log.Println(err)
	}
	log.Println("[info] sent!")

	return err
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
