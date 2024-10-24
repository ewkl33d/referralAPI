# API для системы реферальных ссылок.

### Необходимо наличие PostgreSQL и Redis

### Для начала работы редактируйте .env файл под свою систему, выполните миграции бд и запустите Redis

# Endpoint'ы
По умолчанию, следующие endpoint'ы доступны:

POST /auth/register - Регистрация нового пользователя.

POST /auth/login - Вход пользователя.

DELETE /referral/delete - Удаление реферального кода.

GET /referral/get/:email - Получение реферального кода по email.

POST /referral/register - Регистрация нового пользователя с использованием реферального кода.

GET /referral/referrals/:id - Получение списка пользователей, приглашенных определенным реферером.

# Тестирование

## Swagger
Для тестирования API вы можете использовать Swagger UI, который доступен по адресу /swagger/index.html

## Postman
Также вы можете использовать Postman для тестирования API. Примеры запросов:

### 1. Регистрация пользователя:

Метод: POST

URL: /auth/register

Пример тела запроса:

```json
{
  "email": "test@example.com",
  "password": "password"
}
```

### 2. Вход пользователя:

Метод: POST

URL:/auth/login

Пример тела запроса:

```json
{
"email": "test@example.com",
"password": "password"
}
```

### В ответе придет jwt токен, который в последствии нужно использовать в последующих запросах в header в формате {Authorization: token}

### 3. Создание реферального кода:

Метод: POST

URL: /referral/create

Пример тела запроса:

```json
{
"expiry": "2024-12-31T23:59:59Z"
}
```

### 4. Удаление реферального кода:

Метод: DELETE

URL: /referral/delete

Метод удаляет код пользователя, чей jwt token пришел

### 5. Получение реферального кода по email:

Метод: GET

пример URL:/referral/get/test@example.com

### 6. Регистрация нового пользователя по реферальному коду:

Метод: POST

URL: /referral/register

Пример тела запроса:

```json
{
"email": "newuser@example.com",
"password": "newpassword",
"referral_code": "ABC123"
}
```

### 7. Получение списка пользователей, приглашенных реферером:

Метод: GET

пример URL: /referral/referrals/1