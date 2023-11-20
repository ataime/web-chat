FROM node:16
WORKDIR /app
COPY chat-app .
RUN npm install
CMD ["npm", "run", "serve"]
