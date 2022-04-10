## coverme

В этой задаче нужно покрыть простой todo-app http сервис unit тестами.

Необходимо покрыть все sub-package'и.
Package main можно не тестировать.

Существующие файлы менять не нужно.
Нужно создавать новые файлы с тестами.

Поскольку целью задачи являются именно unit тесты, для решения не нужны вызовы `app.Start` или `app.run`.

Тестирующая система будет проверять code coverage.
Порог задан в [coverage_test.go](./app/coverage_test.go)

Важно понимать, что coverage 100% - не решение всех проблем.
В коде по-прежнему могут быть ошибки.
Coverage 100% говорит ровно о том, что все строки кода выполнялись.
Хорошие тесты, в первую очередь, тестируют функциональность.

Как посмотреть общий coverage:
```
go test -v -cover ./coverme/...
```

Как посмотреть coverage пакета в html:
```
go test -v -coverprofile=/tmp/coverage.out ./coverme/models/...
go tool cover -html=/tmp/coverage.out
```
Аналогичная функциональность поддерживается в Goland.
Также рекомендуем ознакомиться с рассказом о cover в блоге (см. ссылки).

## Ссылки

1. слайды: https://p.go.manytask.org/04-testing/lecture.slide
2. cover: https://blog.golang.org/cover
3. assertions: https://github.com/stretchr/testify
4. [gomock](https://github.com/golang/mock) для создания мока базы данных при тестировании сервера
5. [httptest.ResponseRecorder](https://golang.org/pkg/net/http/httptest/#ResponseRecorder) для тестирования handler'ов сервера
6. [httptest.Server](https://golang.org/pkg/net/http/httptest/#Server) для тестирования клиента
7. Если вы ждёте, когда же выложат лекцию: https://www.youtube.com/watch?v=ndmB0bj7eyw

## O сервисе

Todo-app с минимальной функциональностью + client.

Запуск:
```
✗ go run ./coverme/main.go -port 6029
```

Health check:
```
✗ curl -i -X GET localhost:6029/
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 19 Mar 2020 21:46:02 GMT
Content-Length: 24

"API is up and working!"
```

Создать новое todo:
```
✗ curl -i localhost:6029/todo/create -d '{"title":"A","content":"a"}'
HTTP/1.1 201 Created
Content-Type: application/json
Date: Thu, 19 Mar 2020 21:41:31 GMT
Content-Length: 51

{"id":0,"title":"A","content":"a","finished":false}
```

Получить todo по id:
```
✗ curl -i localhost:6029/todo/0
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 19 Mar 2020 21:44:17 GMT
Content-Length: 51

{"id":0,"title":"A","content":"a","finished":false}
```

Получить все todo:
```
✗ curl -i -X GET localhost:6029/todo
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 19 Mar 2020 21:44:37 GMT
Content-Length: 53

[{"id":0,"title":"A","content":"a","finished":false}]
```

Завершить todo:
```
✗ curl -i -X POST localhost:6029/todo/0/finish
HTTP/1.1 200 OK
Date: Thu, 24 Mar 2022 15:40:49 GMT
Content-Length: 0

✗ curl -i -X GET localhost:6029/todo
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 24 Mar 2022 15:41:04 GMT
Content-Length: 52

[{"id":0,"title":"A","content":"a","finished":true}]%
```
