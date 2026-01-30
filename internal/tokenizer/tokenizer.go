package tokenizer

import (
	"log"

	"github.com/pkoukk/tiktoken-go"
)

var encoding *tiktoken.Tiktoken

// Init initializes the tokenizer with cl100k_base encoding
// This should be called at startup to preload the encoding data
func Init() error {
	var err error
	encoding, err = tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		log.Printf("[WARN] Failed to initialize tiktoken: %v, using fallback", err)
		return err
	}
	log.Printf("[INFO] Tiktoken initialized with cl100k_base encoding")
	return nil
}

// CountTokens counts the number of tokens in a text string
func CountTokens(text string) int {
	if encoding == nil {
		// Fallback: estimate ~4 characters per token
		return len(text) / 4
	}
	return len(encoding.Encode(text, nil, nil))
}
