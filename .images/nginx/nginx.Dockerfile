FROM nginx:stable-alpine

COPY ./.images/nginx/laravel.conf /etc/nginx/conf.d/default.conf

EXPOSE 80