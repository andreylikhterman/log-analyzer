package saver

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaver_Save_AdocFormat(t *testing.T) {
	saver := NewSaver()

	const (
		output = "Some content"
		name   = "testfile"
		format = "adoc"
	)

	err := saver.Save(output, name, format)

	assert.NoError(t, err)

	filename := name + ".adoc"
	defer os.Remove(filename)

	_, err = os.Stat(filename)
	assert.NoError(t, err)

	content, err := os.ReadFile(filename)
	assert.NoError(t, err)
	assert.Equal(t, output, string(content))
}

func TestSaver_Save_MarkdownFormat(t *testing.T) {
	saver := NewSaver()

	const (
		output = "Some content"
		name   = "testfile"
		format = "markdown"
	)

	err := saver.Save(output, name, format)

	assert.NoError(t, err)

	filename := name + ".md"
	defer os.Remove(filename)

	_, err = os.Stat(filename)
	assert.NoError(t, err)

	content, err := os.ReadFile(filename)
	assert.NoError(t, err)
	assert.Equal(t, output, string(content))
}

func TestSaver_Save_UnsupportedFormat(t *testing.T) {
	saver := NewSaver()

	const (
		output = "Some content"
		name   = "testfile"
		format = "unsupported"
	)

	err := saver.Save(output, name, format)

	assert.Error(t, err)
	assert.Equal(t, "неподдеживаемый формат: unsupported", err.Error())
}

func TestSaver_Save_CreateFileError(t *testing.T) {
	saver := NewSaver()

	output := "Some content"
	name := "nonexistent/testfile"
	format := "markdown"

	err := saver.Save(output, name, format)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "невозможно создать файл")
}

func TestSaver_Save_WriteFileError(t *testing.T) {
	saver := NewSaver()

	output := "Some content"
	name := "/root/testfile"
	format := "markdown"

	err := saver.Save(output, name, format)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "невозможно создать файл")
}
