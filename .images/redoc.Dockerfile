FROM node:latest-slim

WORKDIR /app

RUN npm install -g @redocly/cli

EXPOSE 80