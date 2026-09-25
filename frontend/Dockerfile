FROM node:24-alpine AS dependencies

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm ci


FROM node:24-alpine AS builder

WORKDIR /app

COPY --from=dependencies \
    /app/node_modules \
    ./node_modules

COPY . .

RUN npm run build


FROM node:24-alpine AS production-dependencies

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm ci --omit=dev \
    && npm cache clean --force


FROM node:24-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production

COPY --from=builder \
    /app/build \
    ./build

COPY --from=production-dependencies \
    /app/node_modules \
    ./node_modules

COPY package.json ./

USER node

EXPOSE 3000

CMD ["node", "build"]