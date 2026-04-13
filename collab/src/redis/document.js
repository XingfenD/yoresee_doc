const { commandOptions } = require('redis');
const { getClient } = require('./client');
const config = require('../config');

async function getListBuffers(key) {
  const client = getClient();
  if (!client) return [];
  try {
    const items = await client.lRange(commandOptions({ returnBuffers: true }), key, 0, -1);
    return items || [];
  } catch (err) {
    console.error(`Failed to lrange buffer for ${key}:`, err.message);
    return [];
  }
}

async function appendListBuffer(key, buffer) {
  const client = getClient();
  if (!client) return;
  try {
    const payload = Buffer.isBuffer(buffer) ? buffer : Buffer.from(buffer);
    await client.rPush(key, payload);
    if (config.docUpdatesTTL > 0) {
      await client.expire(key, config.docUpdatesTTL);
    }
  } catch (err) {
    console.error(`Failed to rpush buffer for ${key}:`, err.message);
  }
}

async function replaceListBuffer(key, buffer) {
  const client = getClient();
  if (!client) return;
  try {
    const payload = Buffer.isBuffer(buffer) ? buffer : Buffer.from(buffer);
    const pipeline = client.multi();
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
  const client = getClient();
  if (!client) return;
  try {
    await client.sAdd(config.dirtyDocSetKey, docId);
  } catch (err) {
    console.error(`Failed to add dirty doc ${docId}:`, err.message);
  }
}

async function removeDirtyDoc(docId) {
  const client = getClient();
  if (!client) return;
  try {
    await client.sRem(config.dirtyDocSetKey, docId);
  } catch (err) {
    console.error(`Failed to remove dirty doc ${docId}:`, err.message);
  }
}

async function isDirtyDoc(docId) {
  const client = getClient();
  if (!client) return false;
  try {
    const exists = await client.sIsMember(config.dirtyDocSetKey, docId);
    return !!exists;
  } catch (err) {
    console.error(`Failed to check dirty doc ${docId}:`, err.message);
    return false;
  }
}

module.exports = {
  getListBuffers,
  appendListBuffer,
  replaceListBuffer,
  addDirtyDoc,
  removeDirtyDoc,
  isDirtyDoc,
};
