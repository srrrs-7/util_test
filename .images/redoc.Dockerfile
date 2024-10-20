FROM node:latest

WORKDIR /app

RUN npm install @redocly/cli

EXPOSE 80

CMD ["npx", "redoc-cli", "bundle", "/app/openapi.yaml", "--output", "/app/index.html", "--watch"]