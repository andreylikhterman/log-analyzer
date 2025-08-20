# 📊 Анализатор логов NGINX — консольная утилита на Go

Программа для анализа логов веб-сервера **NGINX**. Позволяет быстро получить сводную статистику по запросам, ресурсам, кодам ответов и размерам ответов сервера. Поддерживает фильтрацию по временному диапазону и вывод отчётов в формате **Markdown** или **AsciiDoc**.

---

## 🧩 Возможности
- 📂 Поддержка локальных лог-файлов (с шаблонами `glob`) и загрузки по URL
- ⏳ Фильтрация записей по временному диапазону (`from` / `to` в формате ISO8601)
- 🔍 Фильтрация логов по значению поля (`--filter-field` и `--filter-value`)
- 📊 Подсчёт общего количества запросов
- 🔝 Определение самых популярных ресурсов
- 📡 Анализ распределения кодов ответа HTTP
- 📉 Расчёт среднего размера ответа сервера
- 📐 Определение **95-го перцентиля** размера ответа
- 📝 Генерация отчётов в форматах **Markdown** и **AsciiDoc**

---

## 🛠 Установка и запуск

### Способ 1: запуск напрямую
```bash
git clone https://github.com/andreylikhterman/log-analyzer.git
cd log-analyzer
go run cmd/analyzer/main.go --path logs/access.log --format markdown
```

### Способ 2: сборка и запуск
```bash
git clone https://github.com/andreylikhterman/log-analyzer.git
cd log-analyzer
make build
./bin/analyzer --path logs/2024* --from 2024-08-31 --format adoc
```

## 📸 Примеры запуска
```bash
analyzer --path logs/2024* --from 2024-08-31 --filter-field method --filter-value "GET"
```
```bash
analyzer --path https://raw.githubusercontent.com/elastic/examples/master/Common%20Data%20Formats/nginx_logs/nginx_logs --format adoc
```
```bash
analyzer --path logs/**/2024-08-31.txt
```

## 📑 Пример отчёта

[📄 Пример отчёта (Markdown)](assets/analyze.md)

[📄 Пример отчёта (AsciiDoc)](assets/analyze.adoc)
