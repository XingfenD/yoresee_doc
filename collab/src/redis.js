const { createClient, commandOptions } = require('redis');
const config = require('./config');

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

async function checkRoomExists(roomName) {
  if (!redisClient) return false;

  try {
    const yjsKey = `yjs:${roomName}`;
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

async function getBuffer(key) {
  if (!redisClient) return null;
  try {
    return await redisClient.get(commandOptions({ returnBuffers: true }), key);
  } catch (err) {
    console.error(`Failed to get buffer for ${key}:`, err.message);
    return null;
  }
}

async function setBuffer(key, buffer) {
  if (!redisClient) return;
  try {
    const payload = Buffer.isBuffer(buffer) ? buffer : Buffer.from(buffer);
    await redisClient.set(key, payload);
  } catch (err) {
    console.error(`Failed to set buffer for ${key}:`, err.message);
  }
}

async function appendBuffer(key, buffer) {
  if (!redisClient) return;
  try {
    const payload = Buffer.isBuffer(buffer) ? buffer : Buffer.from(buffer);
    await redisClient.append(key, payload);
  } catch (err) {
    console.error(`Failed to append buffer for ${key}:`, err.message);
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
  getBuffer,
  setBuffer,
  appendBuffer,
  addDirtyDoc,
  removeDirtyDoc,
  isDirtyDoc,
  closeRedis,
  get redisClient() {
    return redisClient;
  }
};
