package input

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/vorduin/slices"
)

type RequestTemplate struct {
	Name          string
	RequeredFlags []string
	OptionalFlags []string
}

func checkName(pattern RequestTemplate, parts []string) error {
	if parts[0] != pattern.Name {
		return fmt.Errorf("неверная команда: %s", parts[0])
	}

	return nil
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

	for i := 1; i < len(parts)-1; i += 2 {
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
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	input = strings.TrimSpace(input)
	parts := strings.Fields(input)

	err = checkName(pattern, parts)
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	err = checkFlags(pattern, parts)
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	err = checkCountFlags(pattern, parts)
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	return getFlags(pattern, parts)
}
