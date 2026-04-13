const grpc = require('@grpc/grpc-js');
const { DocumentServiceClient, SystemServiceClient } = require('../gen/yoresee_doc/v1/yoresee_doc_grpc_pb');
const config = require('../config');

let documentServiceClient = null;
let systemServiceClient = null;

const buildMetadata = () => {
  const md = new grpc.Metadata();
  if (config.internalRpcKey) {
    md.set('x-internal-key', config.internalRpcKey);
  }
  return md;
};

function initGrpcClient() {
  try {
    const opts = {
      'grpc.default_service_config': JSON.stringify({
        loadBalancingConfig: [{ round_robin: {} }]
      })
    };
    const creds = grpc.credentials.createInsecure();
    documentServiceClient = new DocumentServiceClient(config.backendAddr, creds, opts);
    systemServiceClient = new SystemServiceClient(config.backendAddr, creds, opts);
    console.log(`gRPC clients initialized, connecting to ${config.backendAddr}`);
  } catch (err) {
    console.error('Failed to initialize gRPC client:', err.message);
  }
}

module.exports = {
  initGrpcClient,
  buildMetadata,
  getDocumentServiceClient: () => documentServiceClient,
  getSystemServiceClient: () => systemServiceClient,
};
