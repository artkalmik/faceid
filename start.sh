#!/bin/bash

# Остановка и удаление существующих контейнеров
docker-compose down

# Удаление старых образов
docker-compose rm -f

# Сборка и запуск контейнеров
docker-compose up --build 