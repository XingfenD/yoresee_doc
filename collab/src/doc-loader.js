const Y = require('yjs');
const redis = require('./redis');
const backend = require('./grpc');

const snapshotTimeoutMs = Number(process.env.BACKEND_SNAPSHOT_TIMEOUT_MS || 3000);

async function loadDoc(docId) {
  const start = Date.now();
  console.log(`[doc-loader] start docId=${docId}`);
  const redisKey = `yjs:doc:${docId}`;
  console.log(`[doc-loader] redis key: ${redisKey}`);

  try {
    const update = await redis.getBuffer(redisKey);
    console.log(`[doc-loader] redis ${docId} ${update && update.length ? 'hit' : 'miss'}`);

    const doc = new Y.Doc();
    if (update && update.length > 0) {
      console.log(`[doc-loader] applying redis update docId=${docId} bytes=${update.length}`);
      Y.applyUpdate(doc, update);
      console.log(`[doc-loader] apply redis update docId=${docId} bytes=${update.length} took=${Date.now() - start}ms`);
      return doc;
    }

    console.log(`[doc-loader] no redis cache, trying snapshot docId=${docId}`);
    try {
      console.log(`[doc-loader] requesting snapshot docId=${docId}`);
      const snapshot = await Promise.race([
        backend.getDocumentYjsSnapshot(docId),
        new Promise((_, reject) => setTimeout(() => reject(new Error('snapshot timeout')), snapshotTimeoutMs))
      ]);
      console.log(`[doc-loader] received snapshot docId=${docId} length=${snapshot ? snapshot.length : 0}`);

      if (snapshot && snapshot.length > 0) {
        console.log(`[doc-loader] applying snapshot docId=${docId} bytes=${snapshot.length}`);
        Y.applyUpdate(doc, snapshot);
        console.log(`[doc-loader] saving snapshot to redis docId=${docId}`);
        await redis.setBuffer(redisKey, snapshot);
        console.log(`[doc-loader] apply snapshot docId=${docId} bytes=${snapshot.length} took=${Date.now() - start}ms`);
      } else {
        console.log(`[doc-loader] snapshot empty docId=${docId} took=${Date.now() - start}ms`);
        console.log(`[doc-loader] trying to get document content docId=${docId}`);
        const content = await backend.getDocumentContent(docId);
        console.log(`[doc-loader] received content docId=${docId} length=${content ? content.length : 0}`);

        if (content) {
          console.log(`[doc-loader] seeding content docId=${docId} length=${content.length}`);
          const ytext = doc.getText('content');
          ytext.insert(0, content);
          const initialUpdate = Y.encodeStateAsUpdate(doc);
          console.log(`[doc-loader] saving initial update to redis docId=${docId}`);
          await redis.setBuffer(redisKey, initialUpdate);
          console.log(`[doc-loader] seeded content docId=${docId} bytes=${initialUpdate.length} took=${Date.now() - start}ms`);
        } else {
          console.log(`[doc-loader] no content to seed docId=${docId} took=${Date.now() - start}ms`);
        }
      }
    } catch (err) {
      console.error(`[doc-loader] snapshot error docId=${docId} err=${err.message} stack=${err.stack} took=${Date.now() - start}ms`);
      // Fallback to getting document content when snapshot fails
      try {
        console.log(`[doc-loader] fallback: trying to get document content docId=${docId}`);
        const content = await backend.getDocumentContent(docId);
        console.log(`[doc-loader] fallback: received content docId=${docId} length=${content ? content.length : 0}`);

        if (content) {
          console.log(`[doc-loader] fallback: seeding content docId=${docId} length=${content.length}`);
          const ytext = doc.getText('content');
          ytext.insert(0, content);
          const initialUpdate = Y.encodeStateAsUpdate(doc);
          console.log(`[doc-loader] fallback: saving initial update to redis docId=${docId}`);
          await redis.setBuffer(redisKey, initialUpdate);
          console.log(`[doc-loader] fallback seeded content docId=${docId} bytes=${initialUpdate.length} took=${Date.now() - start}ms`);
        } else {
          console.log(`[doc-loader] fallback no content to seed docId=${docId} took=${Date.now() - start}ms`);
        }
      } catch (fallbackErr) {
        console.error(`[doc-loader] fallback error docId=${docId} err=${fallbackErr.message} stack=${fallbackErr.stack} took=${Date.now() - start}ms`);
      }
    }

    console.log(`[doc-loader] returning doc docId=${docId} took=${Date.now() - start}ms`);
    return doc;
  } catch (err) {
    console.error(`[doc-loader] fatal error docId=${docId} err=${err.message} stack=${err.stack} took=${Date.now() - start}ms`);
    // Return empty doc as last resort
    return new Y.Doc();
  }
}

module.exports = {
  loadDoc
};