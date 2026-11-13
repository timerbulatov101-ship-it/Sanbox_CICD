FROM node:16.20.0-alpine3.18 AS builder

ARG VERSION="0.0.1"
ARG VUE_APP_API_URL="http://localhost:8081"

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

ENV NODE_ENV=production
ENV VUE_APP_API_URL=${VUE_APP_API_URL}

RUN npm run build -- --dest "/app/build/frontend-${VERSION}/site"

FROM node:16.20.0-alpine3.18 AS runtime

ARG VERSION="0.0.1"

WORKDIR /app

RUN npm install -g http-server

COPY --from=builder /app/build/frontend-${VERSION}/site ./site

EXPOSE 8080

CMD ["http-server", "./site", "-p", "8080", "-a", "0.0.0.0", "-c-1"]