FROM node:13-alpine

RUN addgroup lambda
RUN adduser -D -G lambda lambda

WORKDIR /home/lambda
COPY core handler
COPY lambda handler/lambda
RUN chown -R lambda:lambda handler

USER lambda
WORKDIR /home/lambda/handler
RUN ls -la
RUN npm ci

CMD ["npm", "start"]
