const { createClient } = require('redis');
const { RedisPersistence } = require('y-redis');
const { setPersistence } = require('y-websocket/bin/utils');
const config = require('./config');

let redisClient = null;
let redisPersistence = null;

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
  try {
    const redisOpts = {
      host: config.redisHost,
      port: config.redisPort,
      ...(config.redisPassword && { password: config.redisPassword }),
      db: config.redisDb,
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
  if (redisPersistence) {
    await redisPersistence.destroy();
  }
}

module.exports = {
  initRedis,
  initYRedis,
  updateRoomActiveTime,
  getActiveRooms,
  checkRoomExists,
  closeRedis,
  get redisClient() {
    return redisClient;
  },
  get redisPersistence() {
    return redisPersistence;
  }
};