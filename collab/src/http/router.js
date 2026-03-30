const { writeJson, writeText, getRequestURL } = require('./helpers');
const { handleProbeRequest } = require('./probes');
const { loadDocForRead, encodeSnapshotResponse, encodeBinaryUpdate } = require('./docs');

function decodeDocID(pathname, prefix) {
  try {
    return decodeURIComponent(pathname.slice(prefix.length));
  } catch (_err) {
    return '';
  }
}

function createRequestHandler({ redis, wsUtils, getIsDraining }) {
  const docs = wsUtils.docs;

  return async function handleRequest(req, res) {
    try {
      const url = getRequestURL(req);
      const pathname = url.pathname;

      if (handleProbeRequest(pathname, res, { isDraining: getIsDraining(), redisClient: redis.redisClient })) {
        return;
      }

      if (pathname === '/api/active-rooms') {
        const rooms = await redis.getActiveRooms();
        writeJson(res, 200, { rooms, count: rooms.length });
        return;
      }

      if (pathname.startsWith('/internal/yjs/doc-snapshot/')) {
        const docId = decodeDocID(pathname, '/internal/yjs/doc-snapshot/');
        if (!docId) {
          writeText(res, 400, 'doc id required');
          return;
        }

        const doc = await loadDocForRead(docs, redis, docId);
        if (!doc) {
          writeText(res, 404, 'doc not loaded');
          return;
        }

        writeJson(res, 200, encodeSnapshotResponse(doc));
        return;
      }

      if (pathname.startsWith('/internal/yjs/doc/')) {
        const docId = decodeDocID(pathname, '/internal/yjs/doc/');
        if (!docId) {
          writeText(res, 400, 'doc id required');
          return;
        }

        const doc = await loadDocForRead(docs, redis, docId);
        if (!doc) {
          writeText(res, 404, 'doc not loaded');
          return;
        }

        res.writeHead(200, { 'Content-Type': 'application/octet-stream' });
        res.end(encodeBinaryUpdate(doc));
        return;
      }

      writeText(res, 404, '');
    } catch (err) {
      console.error('Request handling failed:', err);
      writeText(res, 500, 'internal error');
    }
  };
}

module.exports = {
  createRequestHandler
};
