FROM node:lts
RUN curl -fsSL https://bun.sh/install | bash

COPY ./bun /bun

WORKDIR /bun/usedcar

CMD [ "npm", "start" ]