package ai

import (
	"fmt"

	"github.com/praneeth-ayla/autocommenter/internal/ai/gemini"
	"github.com/praneeth-ayla/autocommenter/internal/contextstore"
	"github.com/praneeth-ayla/autocommenter/internal/scanner"
)

// Provider defines the interface for AI comment generation services.
type Provider interface {
	Validate() error                                                                                    // Validate checks if the provider is configured correctly.
	GenerateComments(content string, contexts []contextstore.FileDetails, style string) (string, error) // GenerateComments creates comments for the given content and contexts.
	GenerateContextBatch(files []scanner.Data) ([]contextstore.FileDetails, error)                      // GenerateContextBatch generates context details for multiple files.
	GenerateReadme(contexts []contextstore.FileDetails, existingReadme string) (string, error)          // GenerateReadme generates a README file based on the provided contexts.
}

// SupportedProviders lists the names of AI providers that the application supports.
var SupportedProviders = []string{
	"gemini", // Gemini is currently the only supported provider.
}

// NewProvider creates and returns a new AI provider based on the given name.
func NewProvider(name string) (Provider, error) {
	var p Provider

	switch name {
	case "gemini":
		p = gemini.New() // Instantiate the Gemini AI provider.
	default:
		// Return an error if the provider name is not recognized.
		return nil, fmt.Errorf("unknown provider: %s", name)
	}

	// Validate the newly created provider before returning it.
	if err := p.Validate(); err != nil {
		return nil, err
	}

	return p, nil
}
