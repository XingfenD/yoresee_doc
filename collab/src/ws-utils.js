const Y = require('yjs');
const syncProtocol = require('y-protocols/dist/sync.cjs');
const awarenessProtocol = require('y-protocols/dist/awareness.cjs');

const encoding = require('lib0/dist/encoding.cjs');
const decoding = require('lib0/dist/decoding.cjs');
const { loadDoc } = require('./doc-loader');
const redis = require('./redis');
const mq = require('./mq');
const config = require('./config');

const wsReadyStateConnecting = 0;
const wsReadyStateOpen = 1;
const wsReadyStateClosing = 2;
const wsReadyStateClosed = 3;

const gcEnabled = process.env.GC !== 'false' && process.env.GC !== '0';
const compactThreshold = parseInt(process.env.DOC_UPDATE_COMPACT_THRESHOLD, 10) || 1000;
const dirtyNotifyThreshold = config.dirtyDocNotifyThreshold;

const debounce = (fn, wait, options = {}) => {
  let timeout = null;
  let lastInvoke = 0;
  let lastArgs = null;
  const maxWait = options.maxWait || 0;

  const invoke = () => {
    timeout = null;
    lastInvoke = Date.now();
    const args = lastArgs;
    lastArgs = null;
    if (args) {
      fn(...args);
    }
  };

  return (...args) => {
    lastArgs = args;
    const now = Date.now();
    if (maxWait > 0 && now - lastInvoke >= maxWait) {
      if (timeout) {
        clearTimeout(timeout);
        timeout = null;
      }
      invoke();
      return;
    }
    if (!timeout) {
      timeout = setTimeout(invoke, wait);
    }
  };
};

const docPromises = new Map();
const updateCounts = new Map();
const notifyCounts = new Map();



const messageSync = 0;
const messageAwareness = 1;

const extractDocId = (docName) => {
  if (!docName) {
    return '';
  }
  if (docName.startsWith('doc-')) {
    return docName.slice(4);
  }
  return docName;
};

const persistUpdate = async (docId, doc, update) => {
  if (!docId) {
    return;
  }
  const redisKey = `yjs:doc:${docId}`;
  await redis.appendBuffer(redisKey, update);
  await redis.addDirtyDoc(docId);
  await redis.updateRoomActiveTime(`doc-${docId}`);

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
    await redis.setBuffer(redisKey, merged);
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

const updateHandler = (update, origin, doc) => {
  const encoder = encoding.createEncoder();
  encoding.writeVarUint(encoder, messageSync);
  syncProtocol.writeUpdate(encoder, update);
  const message = encoding.toUint8Array(encoder);
  doc.conns.forEach((_, conn) => send(doc, conn, message));
};

class WSSharedDoc extends Y.Doc {
  constructor(name) {
    super({ gc: gcEnabled });
    this.name = name;
    this.conns = new Map();
    this.awareness = new awarenessProtocol.Awareness(this);
    this.awareness.setLocalState(null);

    const awarenessChangeHandler = ({ added, updated, removed }, conn) => {
      const changedClients = added.concat(updated, removed);
      if (conn !== null) {
        const connControlledIDs = this.conns.get(conn);
        if (connControlledIDs !== undefined) {
          added.forEach(clientID => {
            connControlledIDs.add(clientID);
          });
          removed.forEach(clientID => {
            connControlledIDs.delete(clientID);
          });
        }
      }
      const encoder = encoding.createEncoder();
      encoding.writeVarUint(encoder, messageAwareness);
      encoding.writeVarUint8Array(encoder, awarenessProtocol.encodeAwarenessUpdate(this.awareness, changedClients));
      const buff = encoding.toUint8Array(encoder);
      this.conns.forEach((_, c) => {
        send(this, c, buff);
      });
    };
    this.awareness.on('update', awarenessChangeHandler);
    this.on('update', updateHandler);
  }
}

const createYDoc = async (docname, gc = true) => {
  const doc = new WSSharedDoc(docname);
  doc.gc = gc;

  const docId = extractDocId(docname);
  try {
    const loadedDoc = await loadDoc(docId);
    const update = Y.encodeStateAsUpdate(loadedDoc);
    if (update && update.length > 0) {
      Y.applyUpdate(doc, update);
    }
  } catch (err) {
    console.error(`Failed to load doc ${docId}:`, err.message);
  }

  attachPersistence(doc, docId);
  return doc;
};

const getYDoc = async (docname, gc = true) => {
  if (docPromises.has(docname)) {
    console.log(`[ws-utils] returning existing promise docname=${docname}`);
    return docPromises.get(docname);
  }
  console.log(`[ws-utils] creating new doc docname=${docname}`);
  const createPromise = createYDoc(docname, gc)
    .finally(() => {
      docPromises.delete(docname);
    });
  docPromises.set(docname, createPromise);
  return createPromise;
};

const messageListener = (conn, doc, message) => {
  try {
    const encoder = encoding.createEncoder();
    const decoder = decoding.createDecoder(message);
    const messageType = decoding.readVarUint(decoder);
    switch (messageType) {
      case messageSync:
        encoding.writeVarUint(encoder, messageSync);
        syncProtocol.readSyncMessage(decoder, encoder, doc, conn);
        if (encoding.length(encoder) > 1) {
          send(doc, conn, encoding.toUint8Array(encoder));
        }
        break;
      case messageAwareness: {
        awarenessProtocol.applyAwarenessUpdate(doc.awareness, decoding.readVarUint8Array(decoder), conn);
        break;
      }
      default:
        break;
    }
  } catch (err) {
    console.error(err);
    doc.emit('error', [err]);
  }
};

const closeConn = (doc, conn) => {
  if (doc.conns.has(conn)) {
    const controlledIds = doc.conns.get(conn);
    doc.conns.delete(conn);
    awarenessProtocol.removeAwarenessStates(doc.awareness, Array.from(controlledIds), null);
    if (doc.conns.size === 0) {
      docs.delete(doc.name);
      doc.destroy();
    }
  }
  conn.close();
};

const send = (doc, conn, m) => {
  if (conn.readyState !== wsReadyStateConnecting && conn.readyState !== wsReadyStateOpen) {
    closeConn(doc, conn);
  }
  try {
    conn.send(m, err => {
      if (err != null) {
        closeConn(doc, conn);
      }
    });
  } catch (err) {
    closeConn(doc, conn);
  }
};

exports.setupWSConnection = async (conn, req, { docName = req.url.slice(1).split('?')[0], gc = true } = {}) => {
  const pendingMessages = [];
  let ready = false;
  let closed = false;

  let doc = null;

  conn.on('message', message => {
    if (!ready) {
      pendingMessages.push(message);
      return;
    }
    if (doc) {
      messageListener(conn, doc, new Uint8Array(message));
    }
  });
  conn.on('close', () => {
    closed = true;
    if (ready) {
      closeConn(doc, conn);
    }
  });
  conn.on('error', () => {
    closed = true;
    if (ready) {
      closeConn(doc, conn);
    }
  });

  doc = await getYDoc(docName, gc);
  if (closed || conn.readyState === wsReadyStateClosed) {
    return;
  }

  doc.conns.set(conn, new Set());
  ready = true;

  for (const message of pendingMessages) {
    messageListener(conn, doc, new Uint8Array(message));
  }
  pendingMessages.length = 0;

  const encoder = encoding.createEncoder();
  encoding.writeVarUint(encoder, messageSync);
  syncProtocol.writeSyncStep1(encoder, doc);
  send(doc, conn, encoding.toUint8Array(encoder));

  const awarenessStates = doc.awareness.getStates();
  if (awarenessStates.size > 0) {
    const awarenessEncoder = encoding.createEncoder();
    encoding.writeVarUint(awarenessEncoder, messageAwareness);
    encoding.writeVarUint8Array(
      awarenessEncoder,
      awarenessProtocol.encodeAwarenessUpdate(doc.awareness, Array.from(awarenessStates.keys()))
    );
    send(doc, conn, encoding.toUint8Array(awarenessEncoder));
  }
};
