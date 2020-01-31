# Курс по Го в ШАД

## Программа курса

1. Введение. Программа курса. Отчётность по курсу, критерии
   оценки. Философия дизайна. if, switch, for. Hello, world. Command
   line arguments. Word count. Animated gif. Fetching URL. Fetching
   URL concurrently. Web server. Tour of go. Local IDE
   setup. Submitting solutions to automated
   grading. gofmt. goimports. linting. Submitting PR's with bug fixes.

2. Базовые конструкции языка. names, declarations, variables,
   assignments. type declarations. packages and files. scope. Zero
   value. Выделение памяти. Стек vs куча. Basic data
   types. Constants. Composite data types. Arrays. Slices. Maps. Structs.
   JSON. text/template. string и []byte. Работа с unicode. Unicode
   replacement character.
   Функции. Функции с переменным числом аргументов. Анонимные функции. Ошибки.

3. Методы. Value receiver vs pointer receiver. Embedding. Method
   value. Encapsulation. Интерфейсы. Интерфейсы как
   контракты. io.Writer, io.Reader и их
   реализации. sort.Interface. error. http.Handler. Интерфейсы как
   перечисления. Type assertion. Type switch. The bigger the
   interface, the weaker the abstraction. Обработка ошибок. panic,
   defer, recover. errors.{Unwrap,Is,As}. fmt.Errorf. %w.

4. Горутины и каналы. clock server. echo server. Размер
   канала. Блокирующее и неблокирующее чтение. select
   statement. Channel axioms. `time.After`. `time.NewTicker`. Pipeline
   pattern. Cancellation. Parallel loop. sync.WaitGroup. Обработка
   ошибок в параллельном коде. errgroup.Group. Concurrent web
   crawler. Concurrent directory traversal.

5. Продвинутое тестирование. Subtests. *testing.B. (*T).Logf. (*T).Skipf. (*T).FailNow.
   testing.Short(), testing flags. Генерация моков. testify/{require,assert}. testify/suite. Test fixture.
   Интеграционные тесты. Goroutine leak detector. TestingMain. Coverage. Сравнение бенчмарков.

6. Concurrency with shared memory. sync.Mutex. sync.RWMutex. sync.Cond. atomic. sync.Once.
   Race detector. Async cache. Работа с базой данных. database/sql. sqlx.

7. Package context. Passing request-scoped data. http middleware. chi.Router. Request cancellation.
   Advanced concurrency patterns. Async cache. Graceful server shutdown. context.WithTimeout.
   Batching and cancellation.

8. Reflection. reflect.Type and reflect.Value. struct tags. net/rpc. encoding/gob.
   sync.Map. reflect.DeepEqual.

9. Low-level programming. unsafe. Package binary. bytes.Buffer. cgo,
   syscall.

10. Архитектура GC. Write barrier. Stack growth. GC pause. GOGC. sync.Pool. Шедулер
    горутин. GOMACPROCS. Утечка тредов.

11. Go tooling. pprof. CPU and Memory profiling. Кросс-компиляция. GOOS, GOARCH. CGO_ENABLED=0.
    Build tags. go modules. godoc. x/analysis. Code generation.

12. Полезные библиотеки. CLI applications with cobra. Protobuf and
    GRPC. zap logging.

13. Запасная леция #1. Работа с crypto/* и x/crypto.

14. Запасная леция #2.

## TODO

1. Check go blog.
2. Check gopher puzzlers.
3. Архитектура и паттерны.
4. Core net/http examples.
5. Project layout.
6. Go proverbs. https://go-proverbs.github.io/
7. All stdlib packages.

## Тестирование задач

 - Каждая задача должна лежать в отдельном пакете или пакетах.
 - Зависимостей между пакетами на верхнем уровне быть не должно.
 - Автоматические тесты на задачу должны запускаться через go test.
 - Авторское решение задачи должно лежать в репозитори и включаться
   билд тегом solution.
 - Некоторые тесты в задаче могут быть приватными. Такие тесты должны
   включаться билд тегом private. Эти тесты не будут доступны
   студентам, но будут запускаться в момент проверки решения в
   тестирующей системе.
 - При посылке решения, на сервер отправляются все файлы внутри пакета.
 - При тестировании, используются изменённые файлы пакета, и
   оригинальные файлы тестов.
 - Файл пакета можно защитить от изменения, добавив `// +build !change` в начало файла.
   В этом случае, при тестировании посылки всегда будет использоваться оригинальная версия файла.
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
