#!/bin/bash
for i in `cat Words`; do go run . --deck-name "Sentence Cards" --model-name "Migaku Japanese" --max-sentences 2 --input "$i"  sentences; done