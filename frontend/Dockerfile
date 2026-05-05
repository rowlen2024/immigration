FROM node:22-alpine AS builder
WORKDIR /build
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npx nuxi build

FROM node:22-alpine
WORKDIR /app
COPY --from=builder /build/.output ./.output
COPY --from=builder /build/node_modules ./node_modules
EXPOSE 3000
ENV NITRO_HOST=0.0.0.0
ENV NITRO_PORT=3000
CMD ["node", ".output/server/index.mjs"]
