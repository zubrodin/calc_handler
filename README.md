# Калькулятор на Go

Это сервис для выполнения арифметических операций, который позволяет добавлять выражения, получать результаты и обрабатывать задачи асинхронно.---

# Описание

Calculator Service предоставляет REST API для выполнения арифметических операций. Он обрабатывает запросы на добавление выражений, возвращает список выражений и их результаты, а также управляет задачами с использованием горутин для асинхронной обработки.---

# Технологии

+ Go (Golang)
+ REST API
+ JSON для обмена данными
+ PostgreSQL (если понадобится для хранения данных, не включено в текущую версию)

## Установка

1. Убедитесь, что у вас установлен Go (версия 1.16 или выше).
2. Клонируйте репозиторий:

   ```bash
   git clone https://github.com/zubrodin/calc_handler
   cd ./calc_handler
   
3. Установите зависимости :

   ```bash
   go mod tidy

4. Настройте переменные окружения (опционально):

   Для настройки времени выполнения операций, вы можете установить переменные окружения:
   
   ```bash
   export TIME_ADDITION_MS=100
   export TIME_SUBTRACTION_MS=100
   export TIME_MULTIPLICATIONS_MS=100
   export TIME_DIVISIONS_MS=100
   export COMPUTING_POWER=4

---

# Использование

1. Запустите сервер:

    ```bash
   // запуск оркестра
   go run ./cmd/calc_service/main.go
   // запуск агента
   go run ./internal/agent/main.go

2. Откройте ваш браузер или используйте инструмент, такой как curl или Postman, для тестирования API.
---

# Api 

## 1. Добавление выражения

POST /api/v1/calculate

Тело запроса: 
 ```
 {
  "expression": "3 + 4 * 2"
}
```
Ответ: 
```
{
  "id": "some-unique-id"
}
```
## 2. Получение всех выражений
GET /api/v1/expressions
Ответ:
```
{
  "expressions": [
    {
      "id": "some-unique-id",
      "status": "completed",
      "result": 11
    }
  ]
}
```
## 3. Получение выражения по ID
GET /api/v1/expressions/{id}
Ответ: 
```
{
  "expression": {
    "id": "some-unique-id",
    "status": "completed",
    "result": 11
  }
}
```

## 4. Получение задачи
GET /internal/task
Ответ: 
```
{
  "id": "task-id",
  "arg1": 3,
  "arg2": 4,
  "operation": "+",
  "operation_time": 100
}
```

## 5. Отправка результата задачи
POST /internal/task/result
Тело запроса: 
```
{
  "id": "task-id",
  "result": 7
}
```
Ответ:
```
{
  "status": "result recorded"
}
```
