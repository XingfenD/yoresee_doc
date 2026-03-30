const WebSocket = require('ws');

function bindWebSocketGateway({ server, redis, setupWSConnection }) {
  const wss = new WebSocket.Server({ server });

  wss.on('connection', async (conn, req) => {
    const url = new URL(req.url, `http://${req.headers.host}`);
    const roomName = url.pathname.slice(1);

    if (roomName) {
      redis.updateRoomActiveTime(roomName);
    }

    try {
      await setupWSConnection(conn, req, {
        gc: true
      });
    } catch (err) {
      console.error('Failed to setup WS connection:', err);
      conn.close();
    }
  });

  return wss;
}

module.exports = {
  bindWebSocketGateway
};
