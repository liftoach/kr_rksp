FROM node:20-alpine

WORKDIR /app

COPY finance-tracker/package.json ./
COPY finance-tracker/package-lock.json ./

RUN npm install

COPY finance-tracker ./

EXPOSE 3000

CMD ["npm", "start"]