function startServer({ server, config, redis, grpc }) {
  server.listen(config.port, async () => {
    await redis.initRedis();
    redis.initYRedis();
    grpc.initGrpcClient();

    await grpc.healthCheck();

    console.log(`Collab server listening on ${config.port}`);
    console.log(`Redis config: ${config.redisHost}:${config.redisPort}/${config.redisDb}`);
    console.log('HTTP API: /api/active-rooms');
  });
}

module.exports = {
  startServer
};
