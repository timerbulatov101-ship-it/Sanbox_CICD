# syntax=docker/dockerfile:1.7

FROM node:16.20.2-alpine3.18@sha256:a1f9d027912b58a7c75be7716c97cfbc6d3099f3a97ed84aa490be9dee20e787 AS dependencies

WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci --legacy-peer-deps

FROM node:16.20.2-alpine3.18@sha256:a1f9d027912b58a7c75be7716c97cfbc6d3099f3a97ed84aa490be9dee20e787 AS builder

ARG VERSION="0.0.1"
ARG VUE_APP_API_URL="/"

WORKDIR /app
COPY --from=dependencies /app/node_modules ./node_modules
COPY . .

ENV NODE_ENV=production
ENV VUE_APP_API_URL=${VUE_APP_API_URL}

RUN npm run build -- --dest "/app/build/frontend-${VERSION}/site"

FROM nginx:1.27-alpine@sha256:65645c7bb6a0661892a8b03b89d0743208a18dd2f3f17a54ef4b76fb8e2f2a10 AS runtime

ARG VERSION="0.0.1"

COPY --from=builder /app/build/frontend-${VERSION}/site /usr/share/nginx/html
COPY default.conf /etc/nginx/conf.d/default.conf

EXPOSE 8080
CMD ["nginx", "-g", "daemon off;"]
