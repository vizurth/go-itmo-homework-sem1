# Library

## Задание
В этом домашнем задании вам предстоит реализовать свой собственный сервис **library**.
В последующих домашних заданиях вы будете его развивать.

Ваш сервис должен поддерживать следующее [API](./library.proto):

```protobuf
syntax = "proto3";

package library;

option go_package = "github.com/Go-CT-ITMO/library-yourname;library";

service Library {
  rpc AddBook(AddBookRequest) returns (AddBookResponse) {}
  
  rpc UpdateBook(UpdateBookRequest) returns (UpdateBookResponse) {}

  rpc GetBookInfo(GetBookInfoRequest) returns (GetBookInfoResponse) {}

  rpc RegisterAuthor(RegisterAuthorRequest) returns (RegisterAuthorResponse) {}

  rpc ChangeAuthorInfo(ChangeAuthorInfoRequest) returns (ChangeAuthorInfoResponse) {}

  rpc GetAuthorInfo(GetAuthorInfoRequest) returns (GetAuthorInfoResponse) {}

  rpc GetAuthorBooks(GetAuthorBooksRequest) returns (stream Book) {}
}

message Book {
  string id = 1;
  string name = 2;
  repeated string author_id = 3;
}

message AddBookRequest {
  string name = 1;
  repeated string author_id = 2;
}

message AddBookResponse {
  Book book = 1;
}

message UpdateBookRequest {
  string id = 1;
  string name = 2;
}

message UpdateBookResponse {}

message GetBookInfoRequest {
  string id = 1;
}

message GetBookInfoResponse {
  Book book = 1;
}

message RegisterAuthorRequest {
  string name = 1;
}

message RegisterAuthorResponse {
  string id = 1;
}

message ChangeAuthorInfoRequest {
  string id = 1;
  string name = 2;
}

message ChangeAuthorInfoResponse {}

message GetAuthorInfoRequest {
  string id = 1;
}

message GetAuthorInfoResponse {
  string id = 1;
  string name = 2;
}

message GetAuthorBooksRequest {
  string author_id = 1;
}
```

## Тестирование
* Код тестов можно посмотреть в файле [integration_test.go](./integration/integration_test.go).
* В рамках CI вам **необходимо** реализовать метку **run** для запуска вашего сервиса.
* В рамках CI вы **можете** реализовать метку **generate**, чтобы не пушить сгенерированный код.
* gRPC сервис и gRPC gateway должны быть подняты на портах, указанных в соответствующих переменных окружения **GRPC_PORT** и  **GRPC_GATEWAY_PORT**.
* Для gRPC сервиса и gRPC gateway необходимо реализовать health checks.

```yaml
// library.yaml
make generate
make run &
echo $! > service_pid.txt

// Makefile
run:
  echo "OK"
# TODO: not implemented

generate:
  echo "OK"
# TODO: not implemented
```

```go
mux.HandlePath("GET", "/health", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "OK")
})

s := grpc.NewServer()
grpc_health_v1.RegisterHealthServer(s, health.NewServer())
reflection.Register(s)
```

## Требования
* Необходимо сгенерировать моки и написать свои тесты, степень покрытия будет проверяться в CI.
* Некоторые пути для [gRPC gateway](https://github.com/grpc-ecosystem/grpc-gateway) явно указаны в тестах, это необходимо учитывать.
* Должна быть настроена генерация с использованием **Makefile**.
* Должна быть реализована поддержка валидации.

## Рекомендации
* [Примеры реализаций](https://github.com/Go-CT-ITMO/lectures).
* [Презентация с примерами](https://docs.google.com/presentation/d/1OrPgiktRmxj6fzMUjSM8y3jmQZTMhV3fKSTHIf5WRLc/edit?usp=sharing).
* Для генерации рекомендуется использовать [buf](https://buf.build/explore) или [protoc](https://github.com/protocolbuffers/protobuf/releases).
* Для валидации рекомендуется использовать [ozzo-validation](https://github.com/go-ozzo/ozzo-validation) или [protoc-gen-validate](https://github.com/bufbuild/protoc-gen-validate).
* Рекомендуется использовать **Makefile** из кода учебных практик.
* Не забывайте про логирование.

## Особенности реализации
- Используйте [тесты](./integration/integration_test.go), чтобы осознать недосказанности.
- В этом домашнем задании вы сами организовываете структуру проекта, что будет оцениваться во время ревью.
- В данном домашнем задании необходимо реализовать in-memory хранилище, которое потом будет заменено на базу данных. 

## Сдача
* Открыть pull request из ветки `hw` в ветку `main` **вашего репозитория**.
* В описании PR заполнить количество часов, которые вы потратили на это задание.
* Отправить заявку на ревью в соответствующей форме.
* Время дедлайна фиксируется отправкой формы.
* Изменять файлы в ветке main без PR запрещено.
* Изменять файл [CI workflow](./.github/workflows/library.yaml) запрещено.

## Makefile

Для удобств локальной разработки сделан [`Makefile`](Makefile). Имеются следующие команды:

Запустить полный цикл (линтер, тесты):

```bash 
make all
```

Запустить только тесты:

```bash
make test
``` 

Запустить линтер:

```bash
make lint
```

Подтянуть новые тесты:

```bash
make update
```

При разработке на Windows рекомендуется использовать [WSL](https://learn.microsoft.com/en-us/windows/wsl/install), чтобы
была возможность пользоваться вспомогательными скриптами.
