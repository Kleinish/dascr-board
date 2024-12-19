# Build frontend from source
FROM node:23-alpine AS build

# Set environment variables
ENV VITE_API_BASE=http://localhost:8000/
ENV VITE_API_URL=http://localhost:8000/api
ENV VITE_WS_URL=ws://localhost:8000/ws

WORKDIR /usr/src/app

# Copy project files
COPY . .

# Upgrade npm to the latest version, install dependencies, and build
RUN npm install -g npm@11.0.0 && \
    npm install && \
    npm run build

# Create runtime image
FROM caddy:alpine

# Copy Caddyfile and build output
COPY ./caddy/Caddyfile /etc/caddy/Caddyfile
COPY --from=build /usr/src/app/build /usr/share/caddy

# Expose the application port
EXPOSE 8080
