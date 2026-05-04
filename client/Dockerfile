FROM node:22-alpine AS build-stage
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL
RUN npm run build

FROM node:22-alpine AS release-stage
WORKDIR /app
COPY --from=build-stage /app/.output ./.output
EXPOSE 3000
CMD ["node", ".output/server/index.mjs"]
