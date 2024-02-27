# build frontend
FROM node:18.18.0-alpine AS frontend-stage
WORKDIR /
COPY ./frontend/package*.json ./
RUN npm install
RUN npm install -g vite
COPY ./frontend ./
RUN vite build

# build backend
FROM golang:1.21 AS build-stage
WORKDIR /
COPY ./backend/go.mod ./backend/go.sum ./
RUN go mod download
COPY ./backend /
RUN CGO_ENABLED=0 GOOS=linux go build -o /entrypoint

# Deploy.
FROM gcr.io/distroless/static-debian11 AS release-stage
WORKDIR /
COPY --from=build-stage /entrypoint /entrypoint
COPY --from=build-stage /*.key /
COPY --from=build-stage /*.crt /
COPY --from=frontend-stage /dist ./dist
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/entrypoint"]