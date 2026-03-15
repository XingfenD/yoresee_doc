const port = Number(process.env.PORT || 1234);

const redisHost = process.env.REDIS_HOST || 'redis';
const redisPort = Number(process.env.REDIS_PORT || 6379);
const redisPassword = process.env.REDIS_PASSWORD || '';
const redisDb = Number(process.env.REDIS_DB || 0);

const backendAddr = process.env.BACKEND_ADDR || 'backend:9090';
const dirtyDocTopic = process.env.DIRTY_DOC_TOPIC || 'collab.dirty_docs';
const dirtyDocMqType = process.env.DIRTY_DOC_MQ || 'redis';
const rabbitmqUrl = process.env.RABBITMQ_URL || '';
const internalRpcKey = process.env.INTERNAL_RPC_KEY || '';

module.exports = {
  port,
  redisHost,
  redisPort,
  redisPassword,
  redisDb,
  backendAddr,
  dirtyDocTopic,
  dirtyDocMqType,
  rabbitmqUrl,
  internalRpcKey
};
