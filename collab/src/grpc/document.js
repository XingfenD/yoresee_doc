const pb = require('../gen/yoresee_doc/v1/yoresee_doc_pb');
const { getDocumentServiceClient, buildMetadata } = require('./client');

async function getDocumentContent(documentExternalId) {
  return new Promise((resolve, reject) => {
    const client = getDocumentServiceClient();
    if (!client) {
      return reject(new Error('gRPC client not initialized'));
    }
    const request = new pb.GetDocumentContentRequest();
    request.setDocumentExternalId(documentExternalId);

    client.getDocumentContent(request, buildMetadata(), (err, response) => {
      if (err) {
        return reject(err);
      }
      const base = response?.getBase?.();
      if (base && base.getCode && base.getCode() !== 0) {
        return reject(new Error(base.getMessage() || 'Failed to get document content'));
      }
      resolve(response?.getContent?.() || '');
    });
  });
}

async function getDocumentYjsSnapshot(documentExternalId) {
  return new Promise((resolve, reject) => {
    const client = getDocumentServiceClient();
    if (!client) {
      return reject(new Error('gRPC client not initialized'));
    }
    const request = new pb.GetDocumentYjsSnapshotRequest();
    request.setDocumentExternalId(documentExternalId);

    client.getDocumentYjsSnapshot(request, buildMetadata(), (err, response) => {
      if (err) {
        return reject(err);
      }
      const base = response?.getBase?.();
      if (base && base.getCode && base.getCode() !== 0) {
        return reject(new Error(base.getMessage() || 'Failed to get document snapshot'));
      }
      const state = response?.getState_asU8?.();
      resolve(state && state.length ? Buffer.from(state) : null);
    });
  });
}

async function saveDocumentYjsSnapshot(documentExternalId, state) {
  return new Promise((resolve, reject) => {
    const client = getDocumentServiceClient();
    if (!client) {
      return reject(new Error('gRPC client not initialized'));
    }
    const request = new pb.SaveDocumentYjsSnapshotRequest();
    request.setDocumentExternalId(documentExternalId);
    request.setState(state);

    client.saveDocumentYjsSnapshot(request, buildMetadata(), (err, response) => {
      if (err) {
        return reject(err);
      }
      const base = response?.getBase?.();
      if (base && base.getCode && base.getCode() !== 0) {
        return reject(new Error(base.getMessage() || 'Failed to save document snapshot'));
      }
      resolve(true);
    });
  });
}

module.exports = {
  getDocumentContent,
  getDocumentYjsSnapshot,
  saveDocumentYjsSnapshot,
};
