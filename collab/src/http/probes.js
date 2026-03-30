const { writeJson } = require('./helpers');

function buildReadinessState(isDraining, redisClient) {
  const ready = !isDraining && !!redisClient;
  return {
    ready,
    redisState: redisClient ? 'connected' : 'disconnected',
    detail: isDraining ? 'server is draining' : undefined
  };
}

function handleProbeRequest(pathname, res, { isDraining, redisClient }) {
  if (pathname === '/livez') {
    writeJson(res, 200, { status: 'ok' });
    return true;
  }

  const state = buildReadinessState(isDraining, redisClient);

  if (pathname === '/readyz') {
    writeJson(res, state.ready ? 200 : 503, {
      status: state.ready ? 'ok' : 'not_ready',
      redis: state.redisState,
      detail: state.detail
    });
    return true;
  }

  if (pathname === '/health') {
    writeJson(res, 200, {
      status: 'ok',
      redis: state.redisState,
      persistence: 'custom',
      ready: state.ready ? 'true' : 'false',
      detail: state.detail
    });
    return true;
  }

  return false;
}

module.exports = {
  handleProbeRequest
};
