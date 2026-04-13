const Y = require('yjs');
const redis = require('./redis');
const mq = require('./mq');
const config = require('./config');
const { keyCollabDocUpdates } = require('./key');

const compactThreshold = parseInt(process.env.DOC_UPDATE_COMPACT_THRESHOLD, 10) || 1000;
const dirtyNotifyThreshold = config.dirtyDocNotifyThreshold;

const updateCounts = new Map();
const notifyCounts = new Map();
const dirtyLogged = new Set();

const persistUpdate = async (docId, doc, update) => {
  if (!docId) {
    return;
  }
  const redisKey = keyCollabDocUpdates(docId);
  await redis.appendListBuffer(redisKey, update);
  await redis.addDirtyDoc(docId);
  await redis.updateRoomActiveTime(`doc-${docId}`);
  if (!dirtyLogged.has(docId)) {
    dirtyLogged.add(docId);
    console.log(`[collab] dirty marked docId=${docId} redisKey=${redisKey}`);
  }

  const nextNotify = (notifyCounts.get(docId) || 0) + 1;
  if (nextNotify >= dirtyNotifyThreshold) {
    await mq.publishDirtyDoc(docId);
    notifyCounts.set(docId, 0);
  } else {
    notifyCounts.set(docId, nextNotify);
  }

  const nextCount = (updateCounts.get(docId) || 0) + 1;
  if (nextCount >= compactThreshold) {
    const merged = Y.encodeStateAsUpdate(doc);
    await redis.replaceListBuffer(redisKey, merged);
    updateCounts.set(docId, 0);
  } else {
    updateCounts.set(docId, nextCount);
  }
};

const attachPersistence = (doc, docId) => {
  if (doc._persistenceBound) {
    return;
  }
  doc._persistenceBound = true;
  doc.on('update', (update) => {
    void persistUpdate(docId, doc, update);
  });
};

module.exports = { persistUpdate, attachPersistence };
