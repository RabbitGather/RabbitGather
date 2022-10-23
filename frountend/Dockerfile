FROM node:14.20.1-alpine3.16 as build-stage

WORKDIR /app
COPY ./rabbit_gather .

RUN yarn install && yarn build

FROM nginx:1.23.2
COPY --from=build-stage /app/dist /usr/share/nginx/html

COPY nginx.conf /etc/nginx/conf.d/nginx.conf

