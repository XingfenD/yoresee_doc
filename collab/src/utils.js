const Y = require('yjs');
const redis = require('./redis');
const grpc = require('./grpc');

async function checkAndInitRoom(roomName) {
  if (!redis.redisClient || !grpc.documentServiceClient) {
    return null;
  }

  try {
    const exists = await redis.checkRoomExists(roomName);
    if (exists) {
      console.log(`Room ${roomName} already has data in Redis, skipping init`);
      return null;
    }

    const docId = roomName.replace(/^doc-/, '');
    console.log(`Room ${roomName} is empty, fetching content from backend for doc ${docId}...`);

    const content = await grpc.getDocumentContent(docId);
    if (content) {
      console.log(`Got content from backend, length: ${content.length}`);
      const ydoc = new Y.Doc();
      const ytext = ydoc.getText('content');
      ytext.insert(0, content);
      return ydoc;
    }

    return null;
  } catch (err) {
    console.error(`Failed to check/init room ${roomName}:`, err.message);
    return null;
  }
}

module.exports = {
  checkAndInitRoom
};