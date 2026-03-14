const http = require('http');
const WebSocket = require('ws');
const { setupWSConnection, setPersistence } = require('y-websocket/bin/utils');
const { createClient } = require('redis');
const { RedisPersistence } = require('y-redis');
const grpc = require('@grpc/grpc-js');
const { DocumentServiceClient } = require('./gen/yoresee_doc_grpc_pb');
const Y = require('yjs');

const port = Number(process.env.PORT || 1234);
const server = http.createServer();
const wss = new WebSocket.Server({ server });

const redisHost = process.env.REDIS_HOST || 'redis';
const redisPort = Number(process.env.REDIS_PORT || 6379);
const redisPassword = process.env.REDIS_PASSWORD || '';
const redisDb = Number(process.env.REDIS_DB || 0);

const backendAddr = process.env.BACKEND_ADDR || 'backend:9090';

let redisClient = null;
let redisPersistence = null;
let documentServiceClient = null;

function initGrpcClient() {
  try {
    documentServiceClient = new DocumentServiceClient(
      backendAddr,
      grpc.credentials.createInsecure(),
      {
        'grpc.default_service_config': JSON.stringify({
          loadBalancingConfig: [{ round_robin: {} }]
        })
      }
    );
    console.log(`gRPC client initialized, connecting to ${backendAddr}`);
  } catch (err) {
    console.error('Failed to initialize gRPC client:', err.message);
  }
}

async function getDocumentContent(documentExternalId) {
  return new Promise((resolve, reject) => {
    if (!documentServiceClient) {
      return reject(new Error('gRPC client not initialized'));
    }
    const GetDocumentContentRequest = require('./gen/yoresee_doc_pb').GetDocumentContentRequest;
    const request = new GetDocumentContentRequest();
    request.setDocumentExternalId(documentExternalId);
    
    documentServiceClient.getDocumentContent(
      request,
      (err, response) => {
        if (err) {
          return reject(err);
        }
        const base = response.getBase();
        if (base && base.getCode() !== 0) {
          return reject(new Error(base.getMessage() || 'Failed to get document content'));
        }
        resolve(response.getContent() || '');
      }
    );
  });
}

async function checkAndInitRoom(roomName) {
  if (!redisClient || !documentServiceClient) {
    return null;
  }
  
  const yjsKey = `yjs:${roomName}`;
  try {
    const exists = await redisClient.exists(yjsKey);
    if (exists) {
      console.log(`Room ${roomName} already has data in Redis, skipping init`);
      return null;
    }
    
    const docId = roomName.replace(/^doc-/, '');
    console.log(`Room ${roomName} is empty, fetching content from backend for doc ${docId}...`);
    
    const content = await getDocumentContent(docId);
    if (content) {
      console.log(`Got content from backend, length: ${content.length}`);
      const ydoc = new Y.Doc();
      const ytext = ydoc.getText('content');
      ytext.insert(0, content);
      return ydoc;
    }
    
    return null;
  } catch (err) {
    console.error(`Failed to check/init room ${roomName}:`, err.message);
    return null;
  }
}

async function initRedis() {
  redisClient = createClient({
    socket: {
      host: redisHost,
      port: redisPort
    },
    password: redisPassword || undefined,
    database: redisDb
  });

  redisClient.on('error', (err) => console.error('Redis Client Error', err));

  try {
    await redisClient.connect();
    console.log(`Connected to Redis at ${redisHost}:${redisPort}`);
  } catch (err) {
    console.error('Failed to connect to Redis:', err);
  }
}

function initYRedis() {
  try {
    const redisOpts = {
      host: redisHost,
      port: redisPort,
      ...(redisPassword && { password: redisPassword }),
      db: redisDb,
      maxRetriesPerRequest: null
    };

    redisPersistence = new RedisPersistence({
      redisOpts
    });

    setPersistence(redisPersistence);
    console.log('Yjs Redis persistence initialized with ioredis');
  } catch (err) {
    console.error('Failed to initialize Yjs Redis persistence:', err);
  }
}

async function updateRoomActiveTime(roomName) {
  if (!redisClient) return;

  try {
    const key = `collab:room:${roomName}`;
    await redisClient.set(key, Date.now().toString(), { EX: 3600 * 24 });
  } catch (err) {
    console.error('Failed to update room active time:', err);
  }
}

async function getActiveRooms() {
  if (!redisClient) return [];

  try {
    const keys = await redisClient.keys('collab:room:*');
    return keys.map(key => key.replace('collab:room:', ''));
  } catch (err) {
    console.error('Failed to get active rooms:', err);
    return [];
  }
}

server.on('request', (req, res) => {
  if (req.url === '/health') {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({
      status: 'ok',
      redis: redisClient ? 'connected' : 'disconnected',
      persistence: redisPersistence ? 'initialized' : 'not_initialized'
    }));
    return;
  }

  if (req.url === '/api/active-rooms') {
    getActiveRooms().then(rooms => {
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
    updateRoomActiveTime(roomName);
  }

  setupWSConnection(conn, req, {
    gc: true
  });
});

server.listen(port, async () => {
  await initRedis();
  initYRedis();
  initGrpcClient();

  console.log(`Collab server listening on ${port}`);
  console.log(`Redis config: ${redisHost}:${redisPort}/${redisDb}`);
  console.log(`HTTP API: /api/active-rooms`);
});

process.on('SIGTERM', async () => {
  if (redisClient) {
    await redisClient.quit();
  }
  if (redisPersistence) {
    await redisPersistence.destroy();
  }
  process.exit(0);
});