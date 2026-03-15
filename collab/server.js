const http = require('http');
const WebSocket = require('ws');
const { setupWSConnection } = require('y-websocket/bin/utils');
const config = require('./src/config');
const redis = require('./src/redis');
const grpc = require('./src/grpc');

const server = http.createServer();
const wss = new WebSocket.Server({ server });

server.on('request', (req, res) => {
  if (req.url === '/health') {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({
      status: 'ok',
      redis: redis.redisClient ? 'connected' : 'disconnected',
      persistence: redis.redisPersistence ? 'initialized' : 'not_initialized'
    }));
    return;
  }

  if (req.url === '/api/active-rooms') {
    redis.getActiveRooms().then(rooms => {
      res.writeHead(200, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ rooms, count: rooms.length }));
    });
    return;
  }

  res.writeHead(404);
  res.end();
});

wss.on('connection', (conn, req) => {
  const url = new URL(req.url, `http://${req.headers.host}`);
  const roomName = url.pathname.slice(1);

  if (roomName) {
    redis.updateRoomActiveTime(roomName);
  }

  setupWSConnection(conn, req, {
    gc: true
  });
});

server.listen(config.port, async () => {
  await redis.initRedis();
  redis.initYRedis();
  grpc.initGrpcClient();

  // Perform health check on startup
  await grpc.healthCheck();

  console.log(`Collab server listening on ${config.port}`);
  console.log(`Redis config: ${config.redisHost}:${config.redisPort}/${config.redisDb}`);
  console.log(`HTTP API: /api/active-rooms`);
});

process.on('SIGTERM', async () => {
  await redis.closeRedis();
  process.exit(0);
});