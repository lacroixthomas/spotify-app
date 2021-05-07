FROM node:10.16.0-alpine
WORKDIR /app
RUN npm install -g serve
COPY . ./
RUN rm -rf build
RUN npm install --silent
RUN npm run build
EXPOSE 5000
ENTRYPOINT ["serve", "-s", "build"]