FROM node:18-alpine

WORKDIR /app

COPY ./front .

RUN npm init -y && npm install express http-proxy-middleware

EXPOSE 3000

CMD ["node", "server.js"]