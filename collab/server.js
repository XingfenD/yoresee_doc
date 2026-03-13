const http = require('http');
const WebSocket = require('ws');
const jwt = require('jsonwebtoken');
const { parse } = require('url');
const { setupWSConnection } = require('y-websocket/bin/utils');

const port = Number(process.env.PORT || 1234);
const server = http.createServer();
const wss = new WebSocket.Server({ server });
const jwtSecret = process.env.JWT_SECRET || '';

wss.on('connection', (conn, req) => {
  if (jwtSecret) {
    const { query } = parse(req.url || '', true);
    const token = query?.token;
    if (!token) {
      conn.close(4401, 'token required');
      return;
    }
    try {
      jwt.verify(token, jwtSecret);
    } catch (error) {
      conn.close(4401, 'invalid token');
      return;
    }
  }

  setupWSConnection(conn, req);
});

server.listen(port, () => {
  console.log(`Collab server listening on ${port}`);
});
