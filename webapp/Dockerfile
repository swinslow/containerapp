FROM node:11.1.0

RUN mkdir /webapp
WORKDIR /webapp

ENV PATH /webapp/node_modules/.bin:$PATH

COPY package.json /webapp/package.json
RUN npm install

CMD ["npm", "start"]
