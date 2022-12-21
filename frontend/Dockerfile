FROM node:14.20.1-alpine3.16 as build-stage

WORKDIR /app
COPY . .

RUN cd /app/rabbit_gather

RUN yarn install && yarn build

FROM nginx:1.23.2
COPY --from=build-stage /app/rabbit_gather/dist /usr/share/nginx/html

COPY --from=build-stage /app/nginx.conf /etc/nginx/nginx.conf

