FROM node:22-alpine AS builder

WORKDIR /app
COPY . .
ARG VITE_BACKEND_ADDRESS
ENV VITE_BACKEND_ADDRESS $VITE_BACKEND_ADDRESS
RUN npm install
RUN npm run build

# ---

FROM nginx:1.27.4-alpine-slim AS frontend

COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf  /etc/nginx/conf.d
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
