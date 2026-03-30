const http = require('http');
const { createRequestHandler } = require('./router');

function createHTTPServer({ redis, wsUtils, getIsDraining }) {
  const requestHandler = createRequestHandler({
    redis,
    wsUtils,
    getIsDraining
  });

  return http.createServer((req, res) => {
    void requestHandler(req, res);
  });
}

module.exports = {
  createHTTPServer
};
