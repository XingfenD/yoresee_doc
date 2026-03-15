const Y = require('yjs');
const redis = require('./redis');
const backend = require('./grpc');

const snapshotTimeoutMs = Number(process.env.BACKEND_SNAPSHOT_TIMEOUT_MS || 3000);

async function loadDoc(docId) {
  const start = Date.now();
  console.log(`[doc-loader] start docId=${docId}`);
  const redisKey = `yjs:doc:${docId}`;
  const update = await redis.getBuffer(redisKey);
  console.log(`[doc-loader] redis ${docId} ${update && update.length ? 'hit' : 'miss'}`);

  const doc = new Y.Doc();
  if (update && update.length > 0) {
    Y.applyUpdate(doc, update);
    console.log(`[doc-loader] apply redis update docId=${docId} bytes=${update.length} took=${Date.now() - start}ms`);
    return doc;
  }

  try {
    const snapshot = await Promise.race([
      backend.getDocumentYjsSnapshot(docId),
      new Promise((_, reject) => setTimeout(() => reject(new Error('snapshot timeout')), snapshotTimeoutMs))
    ]);
    if (snapshot && snapshot.length > 0) {
      Y.applyUpdate(doc, snapshot);
      await redis.setBuffer(redisKey, snapshot);
      console.log(`[doc-loader] apply snapshot docId=${docId} bytes=${snapshot.length} took=${Date.now() - start}ms`);
    } else {
      console.log(`[doc-loader] snapshot empty docId=${docId} took=${Date.now() - start}ms`);
      const content = await backend.getDocumentContent(docId);
      if (content) {
        const ytext = doc.getText('content');
        ytext.insert(0, content);
        const initialUpdate = Y.encodeStateAsUpdate(doc);
        await redis.setBuffer(redisKey, initialUpdate);
        console.log(`[doc-loader] seeded content docId=${docId} bytes=${initialUpdate.length} took=${Date.now() - start}ms`);
      } else {
        console.log(`[doc-loader] no content to seed docId=${docId} took=${Date.now() - start}ms`);
      }
    }
  } catch (err) {
    console.error(`[doc-loader] snapshot error docId=${docId} err=${err.message} took=${Date.now() - start}ms`);
  }

  return doc;
}

module.exports = {
  loadDoc
};
