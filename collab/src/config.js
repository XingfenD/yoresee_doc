const port = Number(process.env.PORT || 1234);

const redisHost = process.env.REDIS_HOST || 'redis';
const redisPort = Number(process.env.REDIS_PORT || 6379);
const redisPassword = process.env.REDIS_PASSWORD || '';
const redisDb = Number(process.env.REDIS_DB || 0);

const backendAddr = process.env.BACKEND_ADDR || 'backend:9090';

module.exports = {
  port,
  redisHost,
  redisPort,
  redisPassword,
  redisDb,
  backendAddr
};