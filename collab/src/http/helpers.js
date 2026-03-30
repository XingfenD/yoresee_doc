function writeJson(res, status, payload) {
  res.writeHead(status, { 'Content-Type': 'application/json' });
  res.end(JSON.stringify(payload));
}

function writeText(res, status, body) {
  res.writeHead(status);
  res.end(body);
}

function getRequestURL(req) {
  const host = req.headers && req.headers.host ? req.headers.host : 'localhost';
  return new URL(req.url || '/', `http://${host}`);
}

module.exports = {
  writeJson,
  writeText,
  getRequestURL
};
