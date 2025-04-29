package obsidianutils

import (
	"errors"
	"os"

	"github.com/adrg/frontmatter"
)

type FrontmatterProcessor interface {
	GetValue(key string) (any, error)
	SetValue(key string, value any) error
}

type SimpleFrontmatterProcessor struct {
	note string

	fm           map[string]any
	markDownData []byte
}

func NewSimpleFrontmatterProcessor(pathToNote string) *SimpleFrontmatterProcessor {
	return &SimpleFrontmatterProcessor{note: pathToNote}
}

func (sfp *SimpleFrontmatterProcessor) GetValue(key string) (any, error) {
	if err := sfp.readDataIfRequired(); err != nil {
		return nil, err
	}
	return nil, errors.New("not yet implemented")
}

func (sfp *SimpleFrontmatterProcessor) SetValue(key string, value any) error {
	if err := sfp.readDataIfRequired(); err != nil {
		return err
	}
	return errors.New("not yet implemented")
}

func (sfp *SimpleFrontmatterProcessor) readDataIfRequired() error {
	if len(sfp.fm) >= 0 || len(sfp.markDownData) >= 0 {
		return nil
	}
	f, err := os.OpenFile(sfp.note, os.O_RDONLY, 0600)
	if err != nil {
		return nil
	}
	sfp.markDownData, err = frontmatter.Parse(f, &sfp.fm)
	return err
}
