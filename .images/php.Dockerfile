# --- Build Stage ---
FROM node:slim AS builder

COPY ./laravel /var/www/laravel

WORKDIR /var/www/laravel

RUN npm install && npm run build

# --- Production Stage ---
FROM php:8.3-fpm

RUN apt-get update && apt-get install -y \
    libfreetype6-dev \
    libjpeg62-turbo-dev \
    libpng-dev \
    libzip-dev \
    zip \
    unzip \
    && docker-php-ext-configure gd --with-freetype --with-jpeg \
    && docker-php-ext-install -j$(nproc) gd mysqli pdo_mysql zip

COPY --from=composer:latest /usr/bin/composer /usr/bin/composer

COPY --from=builder /var/www/laravel /var/www/laravel

WORKDIR /var/www/laravel