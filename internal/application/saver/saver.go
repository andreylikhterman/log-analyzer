package saver

import (
	"fmt"
	"os"
)

const (
	formatADOC     = "adoc"
	formatMarkdown = "markdown"
)

type Saver struct{}

func NewSaver() *Saver {
	return &Saver{}
}

func (saver *Saver) Save(output, name, format string) error {
	var filename string

	switch format {
	case formatADOC:
		filename = name + ".adoc"
	case formatMarkdown:
		filename = name + ".md"
	default:
		return fmt.Errorf("неподдеживаемый формат: %s", format)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("невозможно создать файл %s: %v", filename, err)
	}

	defer file.Close()

	_, err = file.WriteString(output)
	if err != nil {
		return fmt.Errorf("невозможно записать в файл %s: %v", filename, err)
	}

	return nil
}
