# Constanta
 test task for Constanta


## 1. Создание платежа
POST запрос   
```sh
http://localhost:8001/transactions
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
http://localhost:8001/transactions/{userId}
http://localhost:8001/transactions/{email}
```

## 3. Проверка статуса платежа по ID
GET запрос
```sh
http://localhost:8001/transactions/status/{id}
```

## 4. Изменение статуса платежа
PUT запрос
```sh 
http://localhost:8001/transactions/status/{id}
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
http://localhost:8001/transactions/{id}
```