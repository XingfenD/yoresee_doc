const grpc = require('@grpc/grpc-js');
const { DocumentServiceClient, SystemServiceClient } = require('./gen/yoresee_doc_grpc_pb');
const config = require('./config');

let documentServiceClient = null;
let systemServiceClient = null;

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
    const GetDocumentContentRequest = require('./gen/yoresee_doc_pb').GetDocumentContentRequest;
    const request = new GetDocumentContentRequest();
    request.setDocumentExternalId(documentExternalId);

    documentServiceClient.getDocumentContent(
      request,
      (err, response) => {
        if (err) {
          return reject(err);
        }
        const base = response.getBase();
        if (base && base.getCode() !== 0) {
          return reject(new Error(base.getMessage() || 'Failed to get document content'));
        }
        resolve(response.getContent() || '');
      }
    );
  });
}

async function healthCheck() {
  if (!systemServiceClient) {
    console.log('System service client not initialized, skipping health check');
    return;
  }

  try {
    const HealthRequest = require('./gen/yoresee_doc_pb').HealthRequest;
    const request = new HealthRequest();

    systemServiceClient.health(request, (err, response) => {
      if (err) {
        console.error('Health check failed:', err);
      } else {
        const base = response.getBase();
        if (base && base.getCode() === 0) {
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
  healthCheck,
  get documentServiceClient() {
    return documentServiceClient;
  },
  get systemServiceClient() {
    return systemServiceClient;
  }
};