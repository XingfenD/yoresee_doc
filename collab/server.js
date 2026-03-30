const wsUtils = require('./src/ws-utils');
const { setupWSConnection } = wsUtils;
const config = require('./src/config');
const redis = require('./src/redis');
const grpc = require('./src/grpc');
const mq = require('./src/mq');
const { createHTTPServer } = require('./src/http/server');
const { bindWebSocketGateway } = require('./src/ws-gateway');
const { startServer } = require('./src/bootstrap');
const { createGracefulShutdown, registerShutdownSignals } = require('./src/lifecycle');

let isDraining = false;

const server = createHTTPServer({
  redis,
  wsUtils,
  getIsDraining: () => isDraining
});
const wss = bindWebSocketGateway({
  server,
  redis,
  setupWSConnection
});
startServer({ server, config, redis, grpc });

const shutdown = createGracefulShutdown({
  server,
  wss,
  mq,
  redis,
  onDraining: () => {
    isDraining = true;
  }
});

registerShutdownSignals(shutdown);
