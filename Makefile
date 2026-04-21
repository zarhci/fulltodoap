include .env
export

export PROJECT_ROOT=$(shell pwd)


up:
	@docker compose up -d todoapp-postgres

down:
	@docker compose down todoapp-postgres
	@docker compose down port-forwarder


cleanup:
	@read -p "Очистить все volume файлы окружения? Опастность утери данных: [y/N] " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down todoapp-postgres && rm -rf out/dbeaver && \
		echo "Все volume файлы окружения удалены!"; \
	else \
		echo "Операция отменена. Volume файлы окружения сохранены."; \
	fi


migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Ошибка: Не указано имя миграции (отсутствует параметр seq). Используйте 'make migrate-create seq=имя_миграции'"; \
		exit 1; \
	fi; \

	docker compose run --rm todoapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"


migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-force:
	@if [ -z "$(version)" ]; then \
		echo "Ошибка: Не указана версия. Используйте make migrate-force version=1"; \
		exit 1; \
	fi; \
	$(MAKE) migrate-action action="force $(version)"

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Ошибка: Не указано действие миграции (отсутствует параметр action). Используйте 'make migrate-action action=up' или 'make migrate-action action=down'"; \
		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate \
		-path /migrations \
		-database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@todoapp-postgres:5432/$(POSTGRES_DB)?sslmode=disable" \
		$(action)


port-forward:
	@docker compose up -d port-forwarder

port-forward-close:
	@docker compose down port-forwarder


todo-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	go mod tidy && \
	go run cmd/main.go