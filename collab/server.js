const http = require('http');
const WebSocket = require('ws');
const Y = require('yjs');
const { setupWSConnection, docs } = require('./src/ws-utils');
const config = require('./src/config');
const redis = require('./src/redis');
const grpc = require('./src/grpc');
const mq = require('./src/mq');

const server = http.createServer();
const wss = new WebSocket.Server({ server });

server.on('request', (req, res) => {
  if (req.url === '/health') {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({
      status: 'ok',
      redis: redis.redisClient ? 'connected' : 'disconnected',
      persistence: 'custom'
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

  if (req.url.startsWith('/internal/yjs/doc/')) {
    const url = new URL(req.url, `http://${req.headers.host}`);
    const docId = decodeURIComponent(url.pathname.replace('/internal/yjs/doc/', ''));
    if (!docId) {
      res.writeHead(400);
      res.end('doc id required');
      return;
    }
    const doc = docs.get(`doc-${docId}`) || docs.get(docId);
    if (!doc) {
      res.writeHead(404);
      res.end('doc not loaded');
      return;
    }
    const update = Y.encodeStateAsUpdate(doc);
    res.writeHead(200, { 'Content-Type': 'application/octet-stream' });
    res.end(Buffer.from(update));
    return;
  }

  res.writeHead(404);
  res.end();
});

wss.on('connection', async (conn, req) => {
  const url = new URL(req.url, `http://${req.headers.host}`);
  const roomName = url.pathname.slice(1);

  if (roomName) {
    redis.updateRoomActiveTime(roomName);
  }

  try {
    await setupWSConnection(conn, req, {
      gc: true
    });
  } catch (err) {
    console.error('Failed to setup WS connection:', err);
    conn.close();
  }
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
  await mq.close();
  await redis.closeRedis();
  process.exit(0);
});
