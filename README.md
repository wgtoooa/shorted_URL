# URL Shortener

Привет, дорогой друг, ты попал на мой первый проект на golang: укоротитель ссылок.

![Логотип](/static/URLSHORT.png)

## 🚀 Возможности

*   **Профиль:** В начале ты регистрируешься и имеешь свой личный аккаунт.
*   **Сокращение ссылок:** Сокращаешь длинную ссылку и получаешь короткую, также ты можешь сам обозначить ее название
*   **Дополнительно:** Можно удалять ссылки, редактировать короткие названия. В профиле хранятся 5 ссылок, при нажатии на которых ты сразу переместишься на закрепленную страницу.

## 📦 Установка

Пошаговое руководство по установке.

1.  Клонируйте репозиторий:
    ```bash
    git clone https://github.com/wgtoooa/shorted_URL.git
    ```
2.  Перейдите в директорию проекта:
    ```bash
    cd ваш-репозиторий
    ```
3.  Заполнить зависимости:

#### Пример файла `.env`

```env
# Database
POSTGRES_USER=user
POSTGRES_PASSWORD=12345
POSTGRES_HOST=postgres
POSTGRES_NAME=myDB
POSTGRES_PORT=5432

# Redis
REDIS_ADDRESS=redis:6379
REDIS_PASSWORD=12345
REDIS_DB=0

# Server
SERVER_PORT=80
SERVER_HOST=0.0.0.0

# API & Security
API_KEY=abc123xyz
JWT_SECRET=your-super-secret-jwt-key

# Environment
PRODUCTION=true 
```
4. Запуск проекта(нужен Docker)
```bash
#собираем образ
docker compose build

#запускаем  проект в фоновом режиме
docker compose up -d
```
5. Тестирование
```bash
#проверим работу
curl http://localhost:(you-port)/health
#Ожидаем status OK
```
## Поддержка
Если возникнут вопросы писать в телеграм **@wgtoooa**