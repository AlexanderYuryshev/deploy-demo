# Демо-приложение для деплоя на хостинги и/или VPS

Пример приложения для развертывания на сервисах хостинга и VDS

## Локальный запуск
### Перед запуском
- Получите необходимые данные по образцу файла [.env.example](.env.example) - адрес, пароль, имя пользователя, имя базы данных для подключения, [данные провайдера](https://create.t3.gg/en/usage/next-auth#setting-up-the-default-discordprovider) авторизации - в примере Discord, [токен](https://github.com/marketplace/models-github) для использования OpenAI, [токен](https://core.telegram.org/bots/faq#how-do-i-create-a-bot) для бота Telegram.
- Поместите переменные окружения в файл .env.
- Установите [Docker](https://www.docker.com/products/docker-desktop/).
### Запуск через Docker
- `docker compose up`

## Деплой
### Vercel/Netlify
1. Создайте свой репозиторий Git с проектом.
2. Заведите аккаунт на [Vercel](https://vercel.com/signup) или на [Netlify](https://app.netlify.com/signup).
3. Создайте базу данных для проекта (можно сделать из кабинета хостинга, основные варианты с бесплатным тарифом - [Prisma Postgres](https://www.prisma.io/postgres) и [Neon](https://neon.com/)).
4. Привяжите репозиторий к проекту в кабинете хостинга. Скорее всего, сразу запустится деплой, но упадет, поскольку базы данных и окружения еще нет.
5. Привяжите базу данных к проекту в кабинете хостинга.
6. В настройках проекта укажите переменные окружения для используемой базы данных (есть вариант интеграции через кабинет хостинга) и OpenAI:
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_DATABASE`
- `OPENAI_API_KEY`
6. Установите скрипт сборки в проекте: `npm run db && next build`. Запустите деплой.
7. Сборка и развертывание проекта будут происходить автоматически при пуше в главную ветку. Также можно запустить деплой вручную через кабинет. Адрес, на котором будет развернуто приложение, также будет отображаться в кабинете хостинга.
- ! Скрипт `npm run db` полностью удалит данные, если они существуют. Для последующих сборок измените скрипт на `next build`.
- ! Telegram-бот таким образом не развернется, эти хостинги подразумевают деплой только фронтенд-приложения. Это нужно делать либо на отдельной машине, либо попробовать адаптировать код - [инструкция для Vercel](https://dev.to/jj/create-a-serverless-telegram-bot-using-go-and-vercel-4fdb), [шаблон (не пробовал)](https://pkg.go.dev/github.com/frasnym/go-telegram-bot-vercel-boilerplate#section-readme).

### VDS
1. Создайте свой репозиторий Git с проектом.
2. Настройте удаленный сервер и получите его данные для подключения через SSH. Создайте на нем папку /app и настройте доступ к ней или измените [джобу](.github\workflows\deploy.yml) так, чтобы там была указана другая и доступная для записи папка.
3. Создайте профиль на [Docker Hub](https://hub.docker.com/signup) и получите его данные для загрузки собранных образов.
4. Укажите в [секретах](https://docs.github.com/ru/actions/how-tos/security-for-github-actions/security-guides/using-secrets-in-github-actions) для репозитория следующие переменные:
- `DOCKERHUB_USERNAME`
- `DOCKERHUB_TOKEN`
- `SSH_PRIVATE_KEY`
- `SSH_HOST`
- `SSH_USER`
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_DATABASE`
- `TELEGRAM_BOT_TOKEN`
- `OPENAI_API_KEY`
5. Сборка образов, их загрузка в Docker Hub и развертывание проекта будут происходить автоматически при пуше в главную ветку. Также можно запустить деплой вручную через GitHub Actions.

## Стек

### Приложение (фронт + бэк)
Использован шаблон Create T3 App. Информацию можно найти на [сайте T3 Stack](https://create.t3.gg/) и в [канале Discord](https://t3.gg/discord).

Библиотеки, которые используются в шаблоне:
- [Next.js](https://nextjs.org)
- [NextAuth.js](https://next-auth.js.org)
- [Prisma](https://prisma.io)
- [Tailwind CSS](https://tailwindcss.com)
- [tRPC](https://trpc.io)
- [Langchain](https://js.langchain.com/docs/introduction/)

Дополнительные гайды по деплою приложения на t3 - [Vercel](https://create.t3.gg/en/deployment/vercel), [Netlify](https://create.t3.gg/en/deployment/netlify), [Docker](https://create.t3.gg/en/deployment/docker).

### Telegram-bot и миграция БД - [Go](https://go.dev/)

### БД - [PostgreSQL](https://www.postgresql.org/)