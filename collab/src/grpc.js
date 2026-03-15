const grpc = require('@grpc/grpc-js');
const { DocumentServiceClient, SystemServiceClient } = require('./gen/yoresee_doc/v1/yoresee_doc_grpc_pb');
const pb = require('./gen/yoresee_doc/v1/yoresee_doc_pb');
const config = require('./config');

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
    documentServiceClient = new DocumentServiceClient(
      config.backendAddr,
      grpc.credentials.createInsecure(),
      {
        'grpc.default_service_config': JSON.stringify({
          loadBalancingConfig: [{ round_robin: {} }]
        })
      }
    );
    systemServiceClient = new SystemServiceClient(
      config.backendAddr,
      grpc.credentials.createInsecure(),
      {
        'grpc.default_service_config': JSON.stringify({
          loadBalancingConfig: [{ round_robin: {} }]
        })
      }
    );
    console.log(`gRPC clients initialized, connecting to ${config.backendAddr}`);
  } catch (err) {
    console.error('Failed to initialize gRPC client:', err.message);
  }
}

async function getDocumentContent(documentExternalId) {
  return new Promise((resolve, reject) => {
    if (!documentServiceClient) {
      return reject(new Error('gRPC client not initialized'));
    }
    const request = new pb.GetDocumentContentRequest();
    request.setDocumentExternalId(documentExternalId);

    documentServiceClient.getDocumentContent(request, buildMetadata(), (err, response) => {
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
    if (!documentServiceClient) {
      return reject(new Error('gRPC client not initialized'));
    }
    const request = new pb.GetDocumentYjsSnapshotRequest();
    request.setDocumentExternalId(documentExternalId);

    documentServiceClient.getDocumentYjsSnapshot(request, buildMetadata(), (err, response) => {
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
    if (!documentServiceClient) {
      return reject(new Error('gRPC client not initialized'));
    }
    const request = new pb.SaveDocumentYjsSnapshotRequest();
    request.setDocumentExternalId(documentExternalId);
    request.setState(state);

    documentServiceClient.saveDocumentYjsSnapshot(request, buildMetadata(), (err, response) => {
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

async function healthCheck() {
  if (!systemServiceClient) {
    console.log('System service client not initialized, skipping health check');
    return;
  }

  try {
    const request = new pb.HealthRequest();
    systemServiceClient.health(request, buildMetadata(), (err, response) => {
      if (err) {
        console.error('Health check failed:', err);
      } else {
        const base = response?.getBase?.();
        if (base && base.getCode && base.getCode() === 0) {
          console.log('Health check successful');
        } else {
          console.error('Health check failed:', base ? base.getMessage() : 'Unknown error');
        }
      }
    });
  } catch (err) {
    console.error('Failed to execute health check:', err.message);
  }
}

module.exports = {
  initGrpcClient,
  getDocumentContent,
  getDocumentYjsSnapshot,
  saveDocumentYjsSnapshot,
  healthCheck,
  get documentServiceClient() {
    return documentServiceClient;
  },
  get systemServiceClient() {
    return systemServiceClient;
  }
};
