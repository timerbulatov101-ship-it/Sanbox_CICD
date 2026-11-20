
FROM node:16.20.0-alpine3.18 AS builder

ARG VERSION="0.0.1"
ARG VUE_APP_API_URL="/"

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

ENV NODE_ENV=production
ENV VUE_APP_API_URL=${VUE_APP_API_URL}

RUN npm run build -- --dest "/app/build/frontend-${VERSION}/site"

FROM nginx:1.27-alpine AS runtime

ARG VERSION="0.0.1"

COPY --from=builder /app/build/frontend-${VERSION}/site /usr/share/nginx/html
COPY default.conf /etc/nginx/conf.d/default.conf

EXPOSE 8080

CMD ["nginx", "-g", "daemon off;"]
