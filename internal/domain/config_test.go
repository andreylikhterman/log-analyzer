package domain

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTempFile(t *testing.T) string {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "testfile")
	require.NoError(t, err, "Не удалось создать временный файл")

	defer tmpFile.Close()

	return tmpFile.Name()
}

func startTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		r.Host = "localhost"
	}))
}

func TestAddPath(t *testing.T) {
	t.Run("ValidLocalPath", func(t *testing.T) {
		config := &Config{}
		tmpFile := createTempFile(t)

		defer os.Remove(tmpFile)

		err := config.AddPath(tmpFile)

		assert.NoError(t, err)
		assert.Equal(t, "local", config.TypePath)
		assert.Equal(t, tmpFile, config.Path)
	})

	t.Run("ValidURLPath", func(t *testing.T) {
		config := &Config{}
		server := startTestServer(t)

		defer server.Close()

		err := config.AddPath(server.URL)

		assert.NoError(t, err)
		assert.Equal(t, "url", config.TypePath)
		assert.Equal(t, server.URL, config.Path)
	})

	t.Run("InvalidPath", func(t *testing.T) {
		config := &Config{}
		err := config.AddPath("/invalid/path")

		assert.Error(t, err, "Ожидалось, что выкинется ошибка для некорректного пути")
	})
}

func TestDateParsing(t *testing.T) {
	config := &Config{}

	t.Run("ValidFromDate", func(t *testing.T) {
		err := config.AddFrom("2023-10-15")

		assert.NoError(t, err)

		expectedDate, _ := time.Parse("2006-01-02", "2023-10-15")
		assert.True(t, config.From.Equal(expectedDate), "Ожидалась дата From %v, но получено %v", expectedDate, config.From)
	})

	t.Run("InvalidFromDate", func(t *testing.T) {
		err := config.AddFrom("15-10-2023")
		assert.Error(t, err, "Ожидалось, что выкинется ошибка для некорректной даты 'from'")
	})

	t.Run("ValidToDate", func(t *testing.T) {
		err := config.AddTo("2023-10-20")

		assert.NoError(t, err)

		expectedDate, _ := time.Parse("2006-01-02", "2023-10-20")
		assert.True(t, config.To.Equal(expectedDate), "Ожидалась дата To %v, но получено %v", expectedDate, config.To)
	})
}

func TestAddFormat(t *testing.T) {
	config := &Config{}

	t.Run("ValidFormat", func(t *testing.T) {
		err := config.AddFormat(AdocFormat)

		assert.NoError(t, err)
		assert.Equal(t, AdocFormat, config.Format, "Ожидался формат %s, но получено %s", AdocFormat, config.Format)
	})

	t.Run("UnsupportedFormat", func(t *testing.T) {
		err := config.AddFormat("unsupported")
		assert.Error(t, err, "Ожидалось, что выкинется ошибка для неподдерживаемого формата")
	})

	t.Run("DefaultFormat", func(t *testing.T) {
		err := config.AddFormat("")
		assert.NoError(t, err)
		assert.Equal(t, DefaultFormat, config.Format, "Ожидался формат по умолчанию %s, но получено %s", DefaultFormat, config.Format)
	})
}

func TestFilterHandling(t *testing.T) {
	config := &Config{}

	t.Run("ValidFilterField", func(t *testing.T) {
		err := config.AddFilterField("agent")

		assert.NoError(t, err)
		assert.Equal(t, "agent", config.FilterField, "Ожидалось поле фильтрации 'agent', но получено %s", config.FilterField)
	})

	t.Run("InvalidFilterField", func(t *testing.T) {
		err := config.AddFilterField("invalidField")
		assert.Error(t, err, "Ожидалось, что выкинется ошибка для некорректного поля фильтрации")
	})

	t.Run("ValidFilterRegex", func(t *testing.T) {
		err1 := config.AddFilterField("agent")
		err2 := config.AddFilterValue("Mozilla.*")

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotNil(t, config.FilterValue, "Ожидалось, что FilterValue не будет nil")
		assert.Equal(t, "Mozilla.*", config.FilterValue.String(),
			"Ожидалось регулярное выражение 'Mozilla.*', но получено %v", config.FilterValue)
	})

	t.Run("InvalidFilterRegex", func(t *testing.T) {
		err := config.AddFilterValue("[invalid")
		assert.Error(t, err, "Ожидалось, что выкинется ошибка для некорректного регулярного выражения")
	})
}
