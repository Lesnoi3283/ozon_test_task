# ozon_test_task
Test task solution for OZON internship

### Запуск

Запуск с in-memory хранилищем (Redis):
1. Изменить в `.env` файле значение `IN_MEMORY_STORAGE` на `true`
2. Команда ```docker compose --profile app --profile redis up --build```

Запуск с PostgreSQL:
1. Изменить в .env файле значение `IN_MEMORY_STORAGE` на `false`
2. Команда ```docker compose --profile app --profile postgres up --build```

### Окружение

Для удобства проверки и запуска требуемые переменные окружения были вынесены в `.env` файл.