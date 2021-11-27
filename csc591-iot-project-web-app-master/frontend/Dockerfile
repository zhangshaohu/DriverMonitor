FROM node:9.10

ENV NPM_CONFIG_LOGLEVEL warn
ARG app_env
ENV APP_ENV $app_env

RUN mkdir -p /frontend
WORKDIR /frontend

COPY ./package.json ./yarn.lock ./
RUN yarn install

COPY ./src/ ./src/
COPY ./public/ ./public/

CMD npm run start

EXPOSE 3000
