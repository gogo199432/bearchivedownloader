FROM node:18-alpine as build

WORKDIR /app
RUN npm install -g pnpm

# pnpm fetch does require only lockfile
COPY pnpm-lock.yaml ./
RUN pnpm fetch

ADD . ./
RUN pnpm install -r --offline

RUN pnpm build

FROM node:18-alpine
WORKDIR /app
RUN npm install -g pnpm
COPY --from=build /app/package.json ./
COPY pnpm-lock.yaml ./
RUN pnpm fetch --prod
RUN pnpm install -r --offline --prod
COPY --from=build /app/build ./build
EXPOSE 3000
CMD ["pnpm","startProd"]
