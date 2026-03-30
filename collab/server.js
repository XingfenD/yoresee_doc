const http = require('http');
const WebSocket = require('ws');
const Y = require('yjs');
const wsUtils = require('./src/ws-utils');
const { setupWSConnection } = wsUtils;
const config = require('./src/config');
const redis = require('./src/redis');
const grpc = require('./src/grpc');
const mq = require('./src/mq');

const server = http.createServer();
const wss = new WebSocket.Server({ server });
let isDraining = false;
let shuttingDown = false;

function writeJson(res, status, payload) {
  res.writeHead(status, { 'Content-Type': 'application/json' });
  res.end(JSON.stringify(payload));
}

function closeServerGracefully() {
  return new Promise((resolve) => {
    server.close(() => resolve());
  });
}

function closeWssGracefully() {
  return new Promise((resolve) => {
    wss.close(() => resolve());
  });
}

server.on('request', async (req, res) => {
  if (req.url === '/livez') {
    writeJson(res, 200, { status: 'ok' });
    return;
  }

  if (req.url === '/readyz') {
    const ready = !isDraining && !!redis.redisClient;
    writeJson(res, ready ? 200 : 503, {
      status: ready ? 'ok' : 'not_ready',
      redis: redis.redisClient ? 'connected' : 'disconnected',
      detail: isDraining ? 'server is draining' : undefined
    });
    return;
  }

  if (req.url === '/health') {
    const ready = !isDraining && !!redis.redisClient;
    writeJson(res, 200, {
      status: 'ok',
      redis: redis.redisClient ? 'connected' : 'disconnected',
      persistence: 'custom',
      ready: ready ? 'true' : 'false',
      detail: isDraining ? 'server is draining' : undefined
    });
    return;
  }

  if (req.url === '/api/active-rooms') {
    redis.getActiveRooms().then(rooms => {
      res.writeHead(200, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ rooms, count: rooms.length }));
    });
    return;
  }

  if (req.url.startsWith('/internal/yjs/doc-snapshot/')) {
    const url = new URL(req.url, `http://${req.headers.host}`);
    const docId = decodeURIComponent(url.pathname.replace('/internal/yjs/doc-snapshot/', ''));
    if (!docId) {
      res.writeHead(400);
      res.end('doc id required');
      return;
    }
    const docs = wsUtils.docs;
    let doc = docs ? (docs.get(`doc-${docId}`) || docs.get(docId)) : null;
    if (!doc) {
      const redisKey = `collab:yjs:doc:updates:${docId}`;
      const updates = await redis.getListBuffers(redisKey);
      if (!updates || updates.length === 0) {
        res.writeHead(404);
        res.end('doc not loaded');
        return;
      }
      doc = new Y.Doc();
      for (const update of updates) {
        if (update && update.length > 0) {
          Y.applyUpdate(doc, update);
        }
      }
    }
    const update = Y.encodeStateAsUpdate(doc);
    const ytext = doc.getText('content');
    const text = ytext ? ytext.toString() : '';
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ state: Buffer.from(update).toString('base64'), content: text }));
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
    const docs = wsUtils.docs;
    let doc = docs ? (docs.get(`doc-${docId}`) || docs.get(docId)) : null;
    if (!doc) {
      const redisKey = `collab:yjs:doc:updates:${docId}`;
      const updates = await redis.getListBuffers(redisKey);
      if (!updates || updates.length === 0) {
        res.writeHead(404);
        res.end('doc not loaded');
        return;
      }
      doc = new Y.Doc();
      for (const update of updates) {
        if (update && update.length > 0) {
          Y.applyUpdate(doc, update);
        }
      }
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

async function shutdown(signal) {
  if (shuttingDown) {
    return;
  }
  shuttingDown = true;
  isDraining = true;
  console.log(`Received ${signal}, start graceful shutdown`);

  const forceExitTimer = setTimeout(() => {
    console.error('Graceful shutdown timeout, force exiting');
    process.exit(1);
  }, 10000);
  forceExitTimer.unref();

  await Promise.all([
    closeWssGracefully(),
    closeServerGracefully()
  ]);

  try {
    await mq.close();
  } catch (err) {
    console.error('Failed to close MQ:', err);
  }
  try {
    await redis.closeRedis();
  } catch (err) {
    console.error('Failed to close Redis:', err);
  }

  clearTimeout(forceExitTimer);
  process.exit(0);
}

process.on('SIGTERM', () => {
  void shutdown('SIGTERM');
});

process.on('SIGINT', () => {
  void shutdown('SIGINT');
});
