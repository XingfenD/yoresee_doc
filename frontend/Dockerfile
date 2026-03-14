# dev
FROM node:20-alpine AS dev
WORKDIR /app

COPY package*.json ./
RUN npm install

EXPOSE 80

CMD ["sh", "-c", "npm install && npm run dev -- --host --port 80"]
