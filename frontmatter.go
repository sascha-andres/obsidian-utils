package obsidianutils

import (
	"errors"
	"os"

	"github.com/adrg/frontmatter"
	"gopkg.in/yaml.v2"
)

// FrontmatterProcessor is an interface for managing key-value pairs in frontmatter metadata.
type FrontmatterProcessor interface {

	// GetValue retrieves the value associated with the given key in the frontmatter metadata.
	// Returns the value as `any` and an error if the key does not exist or another issue occurs.
	GetValue(key string) (any, error)
	// SetValue sets the value associated with the given key in the frontmatter metadata.
	SetValue(key string, value any) error

	// GenerateMarkDownDocument generates a Markdown document with the current frontmatter metadata and content.
	// Returns the document as a byte slice and an error if the generation fails.
	GenerateMarkDownDocument() ([]byte, error)
}

// SimpleFrontmatterProcessor processes markdown files containing frontmatter metadata.
// It allows reading and modifying frontmatter key-value pairs.
// The type stores the file path, frontmatter data, and markdown content.
type SimpleFrontmatterProcessor struct {

	// note stores the file path to the markdown file containing frontmatter metadata.
	note string

	// fm stores the frontmatter data.
	fm map[string]any
	// markDownData stores the markdown content.
	markDownData []byte
}

// NewSimpleFrontmatterProcessor initializes a SimpleFrontmatterProcessor with the given path to a markdown file.
// It allows processing of frontmatter metadata for the specified file.
func NewSimpleFrontmatterProcessor(pathToNote string) *SimpleFrontmatterProcessor {
	return &SimpleFrontmatterProcessor{note: pathToNote}
}

// GenerateMarkDownDocument builds a complete Markdown document by combining frontmatter and markdown content.
// Returns the generated document as a byte slice or an error if no data is available or marshalling fails.
func (sfp *SimpleFrontmatterProcessor) GenerateMarkDownDocument() ([]byte, error) {
	if len(sfp.fm) == 0 && len(sfp.markDownData) == 0 {
		return nil, errors.New("no markdown data loaded")
	}
	data, err := yaml.Marshal(sfp.fm)
	if err != nil {
		return nil, err
	}
	d := []byte("---\n")
	d = append(d, data...)
	d = append(d, []byte("---\n")...)
	d = append(d, sfp.markDownData...)
	return d, nil
}

// GetValue retrieves the value associated with the given key in the frontmatter metadata.
// Returns the value as `any` and an error if the key does not exist or another issue occurs.
func (sfp *SimpleFrontmatterProcessor) GetValue(key string) (any, error) {
	if err := sfp.readDataIfRequired(); err != nil {
		return nil, err
	}
	if value, ok := sfp.fm[key]; ok {
		return value, nil
	}
	return nil, errors.New("key not found")
}

// SetValue sets the value associated with the given key in the frontmatter metadata.
// Returns an error if the key does not exist or another issue occurs.
func (sfp *SimpleFrontmatterProcessor) SetValue(key string, value any) error {
	if err := sfp.readDataIfRequired(); err != nil {
		return err
	}
	sfp.fm[key] = value
	return nil
}

// readDataIfRequired reads the frontmatter data and markdown content from the file if they have not already been read.
func (sfp *SimpleFrontmatterProcessor) readDataIfRequired() error {
	if len(sfp.fm) > 0 || len(sfp.markDownData) > 0 {
		return nil
	}

	f, err := os.OpenFile(sfp.note, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	sfp.markDownData, err = frontmatter.Parse(f, &sfp.fm)
	return err
}
