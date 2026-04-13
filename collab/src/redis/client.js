const { createClient } = require('redis');
const config = require('../config');

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

async function closeRedis() {
  if (redisClient) {
    await redisClient.quit();
  }
}

function getClient() {
  return redisClient;
}

module.exports = { initRedis, initYRedis, closeRedis, getClient };
