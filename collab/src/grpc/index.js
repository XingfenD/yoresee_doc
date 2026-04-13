const { initGrpcClient, getDocumentServiceClient, getSystemServiceClient } = require('./client');
const { getDocumentContent, getDocumentYjsSnapshot, saveDocumentYjsSnapshot } = require('./document');
const { healthCheck } = require('./system');

module.exports = {
  initGrpcClient,
  getDocumentContent,
  getDocumentYjsSnapshot,
  saveDocumentYjsSnapshot,
  healthCheck,
  get documentServiceClient() {
    return getDocumentServiceClient();
  },
  get systemServiceClient() {
    return getSystemServiceClient();
  },
};
