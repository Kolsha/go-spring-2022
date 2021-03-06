# Тестирование задач

 - Каждая задача должна лежать в отдельном пакете или пакетах.
 - Зависимостей между пакетами на верхнем уровне быть не должно.
 - Автоматические тесты на задачу должны запускаться через go test.
 - Авторское решение задачи должно лежать в репозитории и включаться
   билд тегом solution.
 - Некоторые тесты в задаче могут быть приватными. Такие тесты должны
   включаться билд тегом private. Эти тесты не будут доступны
   студентам, но будут запускаться в момент проверки решения в
   тестирующей системе.
 - При посылке решения, на сервер отправляются все файлы внутри пакета.
 - При тестировании, используются изменённые файлы пакета и
   оригинальные файлы тестов.
 - Файл пакета можно защитить от изменения, добавив `//go:build !change` в начало файла.
   В этом случае при тестировании посылки всегда будет использоваться оригинальная версия файла.
 - Тесты могут использовать файлы из директории testdata. Менять testdata нельзя.

```sh
# Запуск тестов студентом
go test ./...

# Запуск тестов при разработке задачи
go test -tags solution,private ./...

# Запуск тестов на сервере
## 1. Скопировать файлы пакета из посылки.
## 2. Скопировать все файлы тестов из приватного репозитория.
## 3. Скопировать !change файлы из приватного репозитория.
## 4. Скопировать testdata из приватного репозитория.
go test -tags private ./...
```

Для проверки submission'ов есть testtool (см. docs/testtool.md).
