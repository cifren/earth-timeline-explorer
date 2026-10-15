# Development: Vite dev server with hot reload; the source is bind-mounted.
FROM node:24-alpine AS dev
RUN apk add --no-cache make
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]

# Build: static bundle.
FROM node:24-alpine AS build
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
ARG VITE_API_URL
ENV VITE_API_URL=${VITE_API_URL}
RUN npm run build

# Production: static files served by Caddy.
FROM caddy:2-alpine AS prod
COPY Caddyfile /etc/caddy/Caddyfile
COPY --from=build /app/dist /srv
