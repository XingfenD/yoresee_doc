const pb = require('../gen/yoresee_doc/v1/yoresee_doc_pb');
const { getSystemServiceClient, buildMetadata } = require('./client');

async function healthCheck() {
  const client = getSystemServiceClient();
  if (!client) {
    console.log('System service client not initialized, skipping health check');
    return;
  }

  try {
    const request = new pb.HealthRequest();
    client.health(request, buildMetadata(), (err, response) => {
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
  healthCheck,
};
