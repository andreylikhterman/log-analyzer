package input

import (
	"fmt"
	"log"
	"os"

	"github.com/vorduin/slices"
)

type RequestTemplate struct {
	RequeredFlags []string
	OptionalFlags []string
}

func checkFlags(pattern RequestTemplate, parts []string) error {
	for _, word := range parts {
		if word[0:2] == "--" {
			if !slices.Contains(pattern.RequeredFlags, word[2:]) && !slices.Contains(pattern.OptionalFlags, word[2:]) {
				return fmt.Errorf("неизвестный флаг: %s", word)
			}
		}
	}

	return nil
}

func checkCountFlags(pattern RequestTemplate, parts []string) error {
	for _, flag := range pattern.RequeredFlags {
		if slices.Count(parts, "--"+flag) != 1 {
			return fmt.Errorf("флаг %s не указан или указан более одного раза", flag)
		}
	}

	for _, flag := range pattern.OptionalFlags {
		if slices.Count(parts, "--"+flag) > 1 {
			return fmt.Errorf("флаг %s не может быть указан более одного раза", flag)
		}
	}

	for i := 0; i < len(parts)-1; i += 2 {
		if !(parts[i][0:2] == "--" && parts[i+1][0:2] != "--") {
			return fmt.Errorf("неверный запрос")
		}
	}

	return nil
}

func getFlags(pattern RequestTemplate, parts []string) map[string]string {
	flags := make(map[string]string)

	for _, flag := range pattern.RequeredFlags {
		flags[flag] = parts[slices.Index(parts, "--"+flag)+1]
	}

	for _, flag := range pattern.OptionalFlags {
		if slices.Contains(parts, "--"+flag) {
			flags[flag] = parts[slices.Index(parts, "--"+flag)+1]
		} else {
			flags[flag] = ""
		}
	}

	return flags
}

func Request(pattern RequestTemplate) (configMap map[string]string) {
	parts := os.Args[1:]

	err := checkFlags(pattern, parts)
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	err = checkCountFlags(pattern, parts)
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	return getFlags(pattern, parts)
}
