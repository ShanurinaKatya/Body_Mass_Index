# Лабораторная работа №6. Завершение интерфейса пользователя в React

Веб-приложение «ИМТ Калькулятор»: SPA на React с менеджером состояний Redux Toolkit,
авторизацией по JWT и взаимодействием с Go-веб-сервисом через кодогенерацию Swagger и Axios.

## Структура репозитория

```
.
├── backend/    # Go-сервис услуг (Gin, GORM, PostgreSQL, Redis, MinIO, Swagger)
├── frontend/   # React SPA (Vite, TypeScript, React-Bootstrap, Redux Toolkit)
└── readme.md
```

- `frontend/readme.md` — документация лабораторной работы №5 (FSD, mock-данные).
- `backend/readme.md` — документация веб-сервиса (методы, авторизация, Swagger).

## Стек

- React 19, TypeScript, Vite 8
- React-Bootstrap
- Redux Toolkit, React Redux, Redux Thunk (встроен в Toolkit)
- Axios
- swagger-typescript-api (кодогенерация клиента)
- Go (Gin, GORM, JWT, Redis, MinIO, PostgreSQL)

## Запуск

### 1. Бэкенд

```bash
cd backend

# PostgreSQL, Redis, Adminer (:8081)
docker compose -f db/docker-compose.yml up -d

# MinIO (:9000, консоль :9001)
docker compose -f minio/docker-compose.yml up -d

# Создать в MinIO публичный на чтение бакет body-mass-index
# (консоль http://localhost:9001, логин root / rootpassword)

go run ./cmd/migrate
go run ./cmd/patients_categories    # :8080
```

Swagger UI: http://localhost:8080/swagger/index.html

### 2. Фронтенд

```bash
cd frontend
npm install
npm run generate-api   # генерация src/shared/api/Api.ts из Swagger
npm run dev            # http://localhost:3000
```

Кодогенерация берёт спецификацию из `http://localhost:8080/swagger/doc.json`,
а если сервис не запущен — из `backend/docs/swagger.json`.

### Тестовые данные

После запуска сервиса создайте пользователя и несколько услуг:

```bash
# Регистрация
curl -X POST http://localhost:8080/api/users/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"123456"}'

# Токен
TOKEN=$(curl -s -X POST http://localhost:8080/api/users/auth \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"123456"}' | python3 -c 'import sys,json;print(json.load(sys.stdin)["access_token"])')

# Черновик
curl -X POST http://localhost:8080/api/patients_categories \
  -H "Authorization: Bearer $TOKEN" -F 'name=ИМТ для школьников'

# Публикация (id из ответа предыдущего запроса)
curl -X PUT http://localhost:8080/api/patients_categories/1/publish \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"description":"Расчёт ИМТ для школьников","age":15,"weight":55}'
```

## Обратный прокси (CORS)

В `frontend/vite.config.ts` настроены два маршрута:

| Маршрут | Куда | Назначение |
|---|---|---|
| `/api` | `http://localhost:8080` | методы веб-сервиса |
| `/s3` | `http://localhost:9000/body-mass-index` | изображения и видео из MinIO |

Браузер обращается только к `localhost:3000`, поэтому ограничение CORS не возникает.
Адреса объектов MinIO (`http://localhost:9000/body-mass-index/...`) на фронтенде
переписываются в `/s3/...` функцией `toMediaUrl`.

## Redux Toolkit

`frontend/src/app/store/patients_categories.ts`:

```
store
├── services  (entities/service/model)
│   ├── tiles, feed, draft
│   ├── loading, error, usingMock
│   └── thunks: fetchTiles, fetchFeed, fetchDraft,
│              createService, publishService, deleteService, likeService
├── filter    (features/filter/model)
│   └── age (100 — все услуги, 0–99 — точное совпадение)
└── user      (entities/user/model)
    └── username, isAuthenticated, loading, error
```

- Для методов услуг используются `createAsyncThunk` и сгенерированный клиент `api.*`.
- Авторизация и деавторизация выполняются **вручную через axios**, без кодогенерации
  и thunk: после ответа состояние в Redux устанавливается обычными action.
- Redux DevTools включён (`devTools: true`).

## Авторизация и хранение JWT

- JWT сохраняется в `localStorage` под ключом `patients_categories.jwt`.
- Axios-interceptor подставляет `Authorization: Bearer <token>` во все запросы.
- Redux в память не сохраняется: после F5 интерфейс возвращается в состояние гостя,
  а JWT остаётся в `Application → Local Storage`.
- При выходе JWT удаляется, состояние пользователя очищается, фильтр сбрасывается.

## Mock-данные

Если бэкенд недоступен **и** пользователь не авторизован, thunk'и подставляют
коллекцию `SERVICES_MOCK` (`frontend/src/shared/model/patients_categories.ts`).
Для авторизованного пользователя ошибка выводится на страницу.

## Страницы

| Маршрут | Назначение | Доступ |
|---|---|---|
| `/feed` | лента услуг, «Далее», лайк | все |
| `/cards` | плитка, фильтр по возрасту, удаление своих услуг | все |
| `/add` | создание черновика и публикация услуги | авторизованные |
| `/login` | вход по логину и паролю | гости |

Панель навигации: «Лента», «Добавление», «Плитка», «Вход» / «Выход».

## Порядок показа (скриншоты)

1. **1–2** — `/feed` и `/cards` без запущенного сервиса (mock, режим гостя).
2. **3–6** — `/feed`, `/cards`, `/add`, `/login` с запущенным бэкендом.
3. **7–8** — изменение записи в БД через Adminer и отображение изменений во фронтенде.
4. **9** — вход, F5: Redux снова гость, JWT остаётся в Local Storage.
5. **10** — JWT из Application используется в Insomnia/Postman:
   `GET /api/patients_categories` с `Authorization: Bearer <token>`;
   у своих услуг `is_creator = 1`.
6. **11–12** — состояния Redux в DevTools и axios в сгенерированном `Api.ts`.
7. **13–14** — thunk'и на сгенерированном клиенте с тематическими именами.
8. **15–16** — Network: запросы `/api/...` через прокси и `/s3/...` к MinIO.

## Проверки

```bash
cd frontend
npm run generate-api
npm run lint
npm run build
```

## Диаграмма прецедентов

Исходник PlantUML: `docs/patients_categories.puml`.
