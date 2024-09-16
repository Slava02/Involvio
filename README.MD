# Involvio

Платформа для увелечения вовлеченности участников организации в ее жизнь

## Содержание
- [Технологии](#технологии)
- [Начало работы](#начало-работы)
- [Тестирование](#тестирование)
- [To do](#to-do)
- [Команда проекта](#команда-проекта)

## Технологии

- GO
- Postgresql

### Архитектура
![img.png](docs/img/img.png)

### База данных
```mermaid
erDiagram
    SPACE {
        INTEGER id PK
        VARCHAR name
        VARCHAR description
    }

    CATEGORY {
        INTEGER id PK
        VARCHAR name
    }

    CATEGORY_SPACE {
        INTEGER category_id PK
        INTEGER space_id PK
    }

    USER {
        INTEGER id PK
        VARCHAR login
    }

    USER_SPACE {
        INTEGER user_id PK
        INTEGER space_id PK
    }

    USER_EVENT {
        INTEGER user_id PK
        INTEGER event_id PK
    }

    EVENT {
        INTEGER id PK
        VARCHAR name
        VARCHAR description
        TIMESTAMP begin_date
        TIMESTAMP end_date
    }

    CATEGORY_EVENT {
        INTEGER event_id PK
        INTEGER category_id PK
        BOOLEAN parameter
    }

    USER_SPACE }|--|| USER : "user_id > id"
    USER_SPACE }|--|| SPACE : "space_id > id"
    USER_EVENT }|--|| USER : "user_id > id"
    USER_EVENT }|--|| EVENT : "event_id > id"
    CATEGORY_SPACE }|--|| SPACE : "space_id > id"
    CATEGORY_SPACE }|--|| CATEGORY : "category_id > id"
    CATEGORY_EVENT }|--|| EVENT : "event_id > id"
    CATEGORY_EVENT }|--|| CATEGORY : "category_id > id"

```

## Начало работы

### Предварительные требования
#### Необходимые компоненты
- [docker](https://docs.docker.com/engine/install/)
- [docker-compose](https://docs.docker.com/compose/install/)
- [migrate cli](https://pkg.go.dev/github.com/golang-migrate/migrate/v4#readme-cli-usage)
#### Переменные среды
- PG_URL - ссылка для подключения к Postgres (по умолчанию - postgres://user:pass@localhost:5432/involvio)
### Запуск
```shell
export DB_PASSWORD=admin ENV_NAME=dev && make dc
```

## Тестирование

## Команда проекта
