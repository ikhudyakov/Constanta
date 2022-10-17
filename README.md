# Эмулятор платежного сервиса

> host: localhost
> port: 8001

## 1. Создание платежа
POST запрос   
```sh
http://{host}:{port}/transactions
```
body:
```sh
{
        "userid": 1,
        "useremail": "no-reply1@test.ru",
        "amount": 2451.21,
        "currency": "EUR"
}
```
## 2. Получение списка платежей пользователя по его ID или email
GET запрос 
```sh  
http://{host}:{port}/transactions/{userId}
http://{host}:{port}/transactions/{email}
```

## 3. Проверка статуса платежа по ID
GET запрос
```sh
http://{host}:{port}/transactions/status/{id}
```

## 4. Изменение статуса платежа
PUT запрос
```sh 
http://{host}:{port}/transactions/status/{id}
```
body:
```sh
{
    "status": "SUCCESS"
}
```
## 5. Удаление платежа
DELETE запрос 
```sh
http://{host}:{port}/transactions/{id}
```