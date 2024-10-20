FROM node:latest

WORKDIR /app

RUN npm install -g @redocly/cli

EXPOSE 80