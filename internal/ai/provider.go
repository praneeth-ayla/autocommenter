package ai

import (
	"fmt"

	"github.com/praneeth-ayla/AutoCommenter/internal/ai/gemini"
	"github.com/praneeth-ayla/AutoCommenter/internal/contextstore"
	"github.com/praneeth-ayla/AutoCommenter/internal/scanner"
)

type Provider interface {
	Validate() error
	GenerateComments(content string, contexts []contextstore.FileDetails) (string, error)
	GenerateContextBatch(files []scanner.Data) ([]contextstore.FileDetails, error)
	GenerateReadme(contexts []contextstore.FileDetails, existingReadme string) (string, error)
}

func NewProvider(name string) (Provider, error) {
	var p Provider

	switch name {
	case "gemini":
		p = gemini.New()
	default:
		return nil, fmt.Errorf("unknown provider: %s", name)
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}

	return p, nil
}
