# Этап 5: PostgreSQL и ручной Helm deployment

## Локальный запуск через Docker Compose

1. Создайте локальный файл окружения:

   ```bash
   cp .env.example .env
   ```

2. Задайте собственный `DB_PASSWORD` в `.env`. Значение должно совпадать с паролем, который передаётся PostgreSQL и backend.

3. Запустите приложение:

   ```bash
   docker compose up -d --build
   docker compose ps
   ```

4. Проверьте подключение PostgreSQL:

   ```bash
   docker compose exec postgres pg_isready -U postgres -d sandbox
   docker compose exec postgres psql -U postgres -d sandbox -c '\dt'
   ```

5. Создайте заказ через API и проверьте, что строка сохранилась:

   ```bash
   curl -i -X POST http://localhost:8081/orders
   docker compose exec postgres psql -U postgres -d sandbox -c 'SELECT id, order_id, created_at FROM orders ORDER BY id DESC;'
   ```

Для установки клиента `psql` в Ubuntu/WSL выполните `sudo apt-get update && sudo apt-get install -y postgresql-client`. Для подключения с хостовой машины используйте `localhost:5432`, пользователя `postgres`, базу `sandbox` и пароль из `.env`.

## PostgreSQL и приложение в Kubernetes через Helm

Создайте Secret на GitHub с именем `DB_PASSWORD` и тем же паролем, который хотите использовать в кластере. В workflow `Deploy to Kubernetes` передаётся пароль в Helm как `database.password`; в Kubernetes создаётся Secret `sandbox-postgres-auth`.

Ручное обновление релиза из корня проекта (переменные читаются из локального `.env`):

```bash
set -a
. ./.env
set +a
IMAGE_TAG=PUT_VERSION_FROM_THE_CI_VERSION_JOB_HERE
helm upgrade --install sandbox ./helm/sandbox \
  --namespace default \
  --set-string image.tag="$IMAGE_TAG" \
  --set-string database.password="$DB_PASSWORD" \
  --wait \
  --timeout 10m
```

Проверка:

```bash
kubectl get pods,pvc,svc -n default
kubectl rollout status statefulset/sandbox-postgres -n default
kubectl rollout status deployment/sandbox-backend -n default
kubectl rollout status deployment/sandbox-frontend -n default
```

Подключение к PostgreSQL из локального `psql` через port-forward:

```bash
kubectl port-forward svc/sandbox-postgres 5433:5432
psql "host=localhost port=5433 dbname=sandbox user=postgres sslmode=disable"
```

В Helm по умолчанию frontend остаётся NodePort для сохранения текущего способа доступа. Chart уже поддерживает переключение на ClusterIP и HTTP Ingress после установки ingress-контроллера: задайте `ingress.enabled=true` и `ingress.host`.

## Важное о workflow

`ci-cd.yml` собирает и публикует образы после push в любую ветку. Для ветки `lab5` используй тег, сгенерированный job `Generate version` в соответствующем запуске CI; тег `latest` обновляется только для `main`. Kubernetes deployment остаётся ручным: используй Helm-команду выше или workflow `Deploy to Kubernetes` через `workflow_dispatch`. Так push не будет неожиданно менять рабочий кластер.
