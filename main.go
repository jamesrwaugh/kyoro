package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/jamesrwaugh/kyoro/acquisition"
	"github.com/jamesrwaugh/kyoro/anki"
	"github.com/jamesrwaugh/kyoro/verification"
	"github.com/urfave/cli"
)

func makeKyoroOptions(cli *cli.Context) Options {
	return Options{
		DeckName:       cli.String("deck-name"),
		ModelName:      cli.String("model-name"),
		MaxSentences:   cli.Int("max-sentences"),
		MonoligualMode: cli.Bool("monolingual"),
		InputPhrase:    cli.String("input"),
		SilentMode:     cli.Bool("silent"),
	}
}

type productionResourceClient struct {
}

func (rc productionResourceClient) Get(address string) (string, error) {
	res, err := http.Get(address)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	return string(bodyBytes), nil
}

func runKyoro(options Options) bool {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	resourceClient := productionResourceClient{}
	anki := anki.NewAnkiConnect(http.DefaultClient, "http://localhost", 8765, logger)
	sentences := acquisition.NewJibikiSentenceretriever(resourceClient, logger)
	meaning := acquisition.NewJdictMeaningRetriever(resourceClient, logger)
	verifier := verification.NewConsoleSentenceVerifier()
	mao := KyoroProduction{anki, sentences, meaning, verifier, logger}
	return mao.Run(options)
}

func main() {
	app := cli.NewApp()
	app.Name = "Kyoro, a Japanese sentence card generator"
	app.Version = "0.6.0"
	app.Authors = []*cli.Author{
		&cli.Author{
			Name:  "James Waugh",
			Email: "james.waugh.r@gmail.com",
		},
	}
	app.Description = `
Kyoro builds bulk Japanese sentence cards faster by pulling sentences from online
sources and importing them into Aki with sentence and vocabulary-focused modes.`
	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "deck-name, d", Usage: "The Anki deck name to place generated cards in"},
		&cli.StringFlag{Name: "model-name, m", Usage: "The model to create cards with. *Required Fields*: japanese,english,reading,sentences"},
		&cli.StringFlag{Name: "input, i", Usage: "The Japanese phrase to create cards around. Leave blank to use the clipboard."},
		&cli.IntFlag{Name: "max-sentences, n", Value: 5, Usage: "Maximum number of example sentences to build"},
		&cli.BoolFlag{Name: "monolingual, 1", Usage: "Create Japanese-only cards with no English text"},
		&cli.BoolFlag{Name: "silent, s", Usage: "Don't create any confirmation dialogs and add cards unequivocally"},
	}
	app.Commands = []*cli.Command{
		{
			Name:  "sentences",
			Usage: "Creates --max-sentences sentence cards, one sentence per card",
			Action: func(c *cli.Context) error {
				options := makeKyoroOptions(c)
				options.SentencesOnFrontMode = true
				runKyoro(options)
				return nil
			},
		},
		{
			Name:  "vocab",
			Usage: "Create a single card with --max-sentences sentences in the card.",
			Action: func(c *cli.Context) error {
				options := makeKyoroOptions(c)
				options.SentencesOnFrontMode = false
				runKyoro(options)
				return nil
			},
		},
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)
}
