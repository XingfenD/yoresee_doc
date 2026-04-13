const Y = require('yjs');
const { keyCollabDocUpdates } = require('../key');

function findDocInMemory(docs, docId) {
  if (!docs) {
    return null;
  }
  return docs.get(`doc-${docId}`) || docs.get(docId) || null;
}

async function buildDocFromRedis(redis, docId) {
  const redisKey = keyCollabDocUpdates(docId);
  const updates = await redis.getListBuffers(redisKey);
  if (!updates || updates.length === 0) {
    return null;
  }

  const doc = new Y.Doc();
  for (const update of updates) {
    if (update && update.length > 0) {
      Y.applyUpdate(doc, update);
    }
  }
  return doc;
}

async function loadDocForRead(docs, redis, docId) {
  const inMemoryDoc = findDocInMemory(docs, docId);
  if (inMemoryDoc) {
    return inMemoryDoc;
  }
  return buildDocFromRedis(redis, docId);
}

const XML_FRAGMENT_DOC_TYPES = new Set(['yoresee_rich_text']);

function encodeSnapshotResponse(doc, docType) {
  const update = Y.encodeStateAsUpdate(doc);
  let content = '';

  if (XML_FRAGMENT_DOC_TYPES.has(docType)) {
    const fragment = doc.getXmlFragment('content');
    content = fragment ? fragment.toString() : '';
  } else {
    const ytext = doc.getText('content');
    content = ytext ? ytext.toString() : '';
  }

  return {
    state: Buffer.from(update).toString('base64'),
    content
  };
}

function encodeBinaryUpdate(doc) {
  return Buffer.from(Y.encodeStateAsUpdate(doc));
}

module.exports = {
  loadDocForRead,
  encodeSnapshotResponse,
  encodeBinaryUpdate
};
