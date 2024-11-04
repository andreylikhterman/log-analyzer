package domain

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"
)

const (
	MarkdownFormat = "markdown"
	AdocFormat     = "adoc"
	DefaultFormat  = MarkdownFormat
)

type Config struct {
	Path        string
	TypePath    string
	From        time.Time
	To          time.Time
	Format      string
	FilterField string
	FilterValue *regexp.Regexp
}

func (config *Config) AddPath(path string) error {
	localPath, valid := ValidLocalPath(path)

	switch {
	case valid:
		config.Path = localPath
		config.TypePath = "local"
	case isValidURL(path):
		config.Path = path
		config.TypePath = "url"
	default:
		return fmt.Errorf("неверный локальный путь к логам или URL")
	}

	return nil
}

func (config *Config) AddFrom(from string) error {
	var err error

	if from == "" {
		config.From = time.Time{}
	} else {
		config.From, err = time.Parse("2006-01-02", from)
		if err != nil {
			return fmt.Errorf("неверный формат даты для --from: %v", err)
		}
	}

	return nil
}

func (config *Config) AddTo(to string) error {
	var err error

	if to == "" {
		config.To = time.Time{}
	} else {
		config.To, err = time.Parse("2006-01-02", to)
		if err != nil {
			return fmt.Errorf("неверный формат даты для --to: %v", err)
		}
	}

	return nil
}

func (config *Config) AddFormat(format string) error {
	switch format {
	case AdocFormat:
		config.Format = AdocFormat
	case MarkdownFormat:
		config.Format = MarkdownFormat
	case "":
		config.Format = DefaultFormat
	default:
		return fmt.Errorf("неподдеживаемый формат: %s", format)
	}

	return nil
}

func (config *Config) AddFilterField(field string) error {
	if slices.Contains(config.getFilterFields(), field) || field == "" {
		config.FilterField = field
	} else {
		return fmt.Errorf("неверное поле фильтрации: %s", field)
	}

	return nil
}

func (config *Config) AddFilterValue(value string) error {
	if config.FilterField == "" && value != "" {
		return fmt.Errorf("необходимо указать поле фильтрации")
	}

	pattern, err := regexp.Compile(value)
	if err != nil {
		return fmt.Errorf("не удалась компиляция регулярного выражения: %v", err)
	}

	config.FilterValue = pattern

	return nil
}

func (config *Config) getFilterFields() []string {
	return []string{"agent", "address", "user", "method", "url", "protocol", "status", "referer"}
}

func isValidURL(path string) bool {
	parsedURL, err := url.Parse(path)
	if err != nil {
		return false
	}

	return parsedURL.Scheme != "" && parsedURL.Host != ""
}

func ValidLocalPath(path string) (string, bool) {
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Ошибка при получении домашней директории: %s", err)
		}

		path = filepath.Join(homeDir, path[1:])
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", false
	}

	matches, err := filepath.Glob(absPath)
	if err == nil && len(matches) > 0 {
		return absPath, true
	}

	if _, err := os.Stat(absPath); err == nil {
		return absPath, true
	}

	return "", false
}
