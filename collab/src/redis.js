const { createClient, commandOptions } = require('redis');
const config = require('./config');
const { keyCollabRoomRaw, keyCollabRoomPattern, parseRoomName, keyCollabYjs } = require('./key');

let redisClient = null;

async function initRedis() {
  redisClient = createClient({
    socket: {
      host: config.redisHost,
      port: config.redisPort
    },
    password: config.redisPassword || undefined,
    database: config.redisDb
  });

  redisClient.on('error', (err) => console.error('Redis Client Error', err));

  try {
    await redisClient.connect();
    console.log(`Connected to Redis at ${config.redisHost}:${config.redisPort}`);
  } catch (err) {
    console.error('Failed to connect to Redis:', err);
  }
}

function initYRedis() {
  console.log('Yjs Redis persistence initialized (custom loader)');
}

async function updateRoomActiveTime(roomName) {
  if (!redisClient) return;

  try {
    const key = keyCollabRoomRaw(roomName);
    await redisClient.set(key, Date.now().toString(), { EX: 3600 * 24 });
  } catch (err) {
    console.error('Failed to update room active time:', err);
  }
}

async function getActiveRooms() {
  if (!redisClient) return [];

  try {
    const keys = await redisClient.keys(keyCollabRoomPattern());
    return keys.map(k => parseRoomName(k));
  } catch (err) {
    console.error('Failed to get active rooms:', err);
    return [];
  }
}

async function checkRoomExists(roomName) {
  if (!redisClient) return false;

  try {
    const yjsKey = keyCollabYjs(roomName);
    const exists = await redisClient.exists(yjsKey);
    return exists;
  } catch (err) {
    console.error(`Failed to check room ${roomName}:`, err.message);
    return false;
  }
}

async function closeRedis() {
  if (redisClient) {
    await redisClient.quit();
  }
}

async function getListBuffers(key) {
  if (!redisClient) return [];
  try {
    const items = await redisClient.lRange(commandOptions({ returnBuffers: true }), key, 0, -1);
    return items || [];
  } catch (err) {
    console.error(`Failed to lrange buffer for ${key}:`, err.message);
    return [];
  }
}

async function appendListBuffer(key, buffer) {
  if (!redisClient) return;
  try {
    const payload = Buffer.isBuffer(buffer) ? buffer : Buffer.from(buffer);
    await redisClient.rPush(key, payload);
    if (config.docUpdatesTTL > 0) {
      await redisClient.expire(key, config.docUpdatesTTL);
    }
  } catch (err) {
    console.error(`Failed to rpush buffer for ${key}:`, err.message);
  }
}

async function replaceListBuffer(key, buffer) {
  if (!redisClient) return;
  try {
    const payload = Buffer.isBuffer(buffer) ? buffer : Buffer.from(buffer);
    const pipeline = redisClient.multi();
    pipeline.del(key);
    pipeline.rPush(key, payload);
    if (config.docUpdatesTTL > 0) {
      pipeline.expire(key, config.docUpdatesTTL);
    }
    await pipeline.exec();
  } catch (err) {
    console.error(`Failed to replace list buffer for ${key}:`, err.message);
  }
}

async function addDirtyDoc(docId) {
  if (!redisClient) return;
  try {
    await redisClient.sAdd(config.dirtyDocSetKey, docId);
  } catch (err) {
    console.error(`Failed to add dirty doc ${docId}:`, err.message);
  }
}

async function removeDirtyDoc(docId) {
  if (!redisClient) return;
  try {
    await redisClient.sRem(config.dirtyDocSetKey, docId);
  } catch (err) {
    console.error(`Failed to remove dirty doc ${docId}:`, err.message);
  }
}

async function isDirtyDoc(docId) {
  if (!redisClient) return false;
  try {
    const exists = await redisClient.sIsMember(config.dirtyDocSetKey, docId);
    return !!exists;
  } catch (err) {
    console.error(`Failed to check dirty doc ${docId}:`, err.message);
    return false;
  }
}

module.exports = {
  initRedis,
  initYRedis,
  updateRoomActiveTime,
  getActiveRooms,
  checkRoomExists,
  getListBuffers,
  appendListBuffer,
  replaceListBuffer,
  addDirtyDoc,
  removeDirtyDoc,
  isDirtyDoc,
  closeRedis,
  get redisClient() {
    return redisClient;
  }
};
