# Kyoro

Kyoro is a configurable command line tool that helps builds Japanese sentence cards in Anki quickly, by taking a term on your clipboard (or by the CLI) and generating a set of sentence cards the term appears in.

## Features

* Multiple sentence or vocabulary card generation
* Minimizes time spent fetching sentences, vocab, readings, audio, and clicking Anki menus when immersing
* Automatically includes audio for each card imported
* Native where possible, text-to-speech otherwise.
* Sentence cards are a candidate for large-scale i+1 optimization with other Anki tools such as [MorphMan](https://github.com/kaegi/MorphMan)

## Usage

TBD

## Use Cases

1. While reading Japanese in digital medium, such as a VN, I want to highlight a word, press a hotkey, and SRS the word for later.
1. After coming across a new word, I want to spend minimal time doing manual labor to get Anki cards created for the word.
1. I have a large corpus of text and I want to SRS it, but the setup time is infeasible. 

## Kyoro vs. Yomichan

[Yomichan](https://foosoft.net/projects/yomichan/) is a popular tool that also supports importing Anki cards for Japanese words, phrases, and kanji.

So, are the two projects different? First, Yomichan is not a tool built for Anki import--it's a pop-up dictionary which *also* supports Anki import.

### Advantages of Kyoro

1. Multiple sentence card import: Import N different cards using a piece of vocabulary
1. Vocabulary card import: Import a single card with N different sentences on it (Yomichan can only include one)
1. Works on text outside of the browser
1. Automatable: Being CLI, it can be used to convert bulk sets of Japanese into sentence cards.
1. Includes audio for the entire sentence, not just a single word.
1. Optionally allows user verification of sentences before importing, i.e, to not import vocabulary above the user's level.

### Advantages of Yomichan

1. Ability to customize data in imported card fields
1. Customizable dictionaries
