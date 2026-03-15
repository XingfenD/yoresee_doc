const amqp = require('amqplib');
const config = require('./config');
const redis = require('./redis');

let rabbitConn = null;
let rabbitChannel = null;

const shouldPublishRedis = () => {
  const mode = (config.dirtyDocMqType || 'redis').toLowerCase();
  return mode === 'redis' || mode === 'both';
};

const shouldPublishRabbit = () => {
  const mode = (config.dirtyDocMqType || 'redis').toLowerCase();
  return mode === 'rabbitmq' || mode === 'rabbit' || mode === 'both';
};

const getRabbitChannel = async () => {
  if (!config.rabbitmqUrl) {
    return null;
  }
  if (rabbitChannel) {
    return rabbitChannel;
  }
  rabbitConn = await amqp.connect(config.rabbitmqUrl);
  rabbitChannel = await rabbitConn.createChannel();
  return rabbitChannel;
};

const publishRabbit = async (topic, payload) => {
  const channel = await getRabbitChannel();
  if (!channel) {
    return;
  }
  await channel.assertExchange(topic, 'topic', { durable: true });
  channel.publish(topic, '', Buffer.from(payload), {
    contentType: 'application/json'
  });
};

const publishRedis = async (topic, payload) => {
  if (!redis.redisClient) {
    return;
  }
  await redis.redisClient.publish(topic, payload);
};

const publishDirtyDoc = async (docId) => {
  if (!docId) {
    return;
  }
  const payload = JSON.stringify({ doc_id: docId, ts: Date.now() });
  try {
    if (shouldPublishRedis()) {
      await publishRedis(config.dirtyDocTopic, payload);
    }
  } catch (err) {
    console.error(`Failed to publish dirty doc ${docId} to redis:`, err.message);
  }
  try {
    if (shouldPublishRabbit()) {
      await publishRabbit(config.dirtyDocTopic, payload);
    }
  } catch (err) {
    console.error(`Failed to publish dirty doc ${docId} to rabbitmq:`, err.message);
  }
};

const close = async () => {
  if (rabbitChannel) {
    await rabbitChannel.close();
    rabbitChannel = null;
  }
  if (rabbitConn) {
    await rabbitConn.close();
    rabbitConn = null;
  }
};

module.exports = {
  publishDirtyDoc,
  close
};
