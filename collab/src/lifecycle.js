function closeServerGracefully(server) {
  return new Promise((resolve) => {
    server.close(() => resolve());
  });
}

function closeWssGracefully(wss) {
  return new Promise((resolve) => {
    wss.close(() => resolve());
  });
}

function createGracefulShutdown({
  server,
  wss,
  mq,
  redis,
  onDraining,
  timeoutMs = 10000
}) {
  let shuttingDown = false;

  return async function shutdown(signal) {
    if (shuttingDown) {
      return;
    }
    shuttingDown = true;
    onDraining();
    console.log(`Received ${signal}, start graceful shutdown`);

    const forceExitTimer = setTimeout(() => {
      console.error('Graceful shutdown timeout, force exiting');
      process.exit(1);
    }, timeoutMs);
    forceExitTimer.unref();

    await Promise.all([
      closeWssGracefully(wss),
      closeServerGracefully(server)
    ]);

    try {
      await mq.close();
    } catch (err) {
      console.error('Failed to close MQ:', err);
    }

    try {
      await redis.closeRedis();
    } catch (err) {
      console.error('Failed to close Redis:', err);
    }

    clearTimeout(forceExitTimer);
    process.exit(0);
  };
}

function registerShutdownSignals(shutdown) {
  process.on('SIGTERM', () => {
    void shutdown('SIGTERM');
  });

  process.on('SIGINT', () => {
    void shutdown('SIGINT');
  });
}

module.exports = {
  createGracefulShutdown,
  registerShutdownSignals
};
