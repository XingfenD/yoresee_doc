function keyCollabDocUpdates(docId) {
  return `collab:yjs:doc:updates:${docId}`;
}

function keyCollabRoom(docId) {
  return `collab:room:doc-${docId}`;
}

function keyCollabRoomRaw(roomName) {
  return `collab:room:${roomName}`;
}

function keyCollabRoomPattern() {
  return 'collab:room:*';
}

function parseRoomName(redisKey) {
  return redisKey.slice('collab:room:'.length);
}

function keyCollabYjs(roomName) {
  return `collab:yjs:${roomName}`;
}

module.exports = {
  keyCollabDocUpdates,
  keyCollabRoom,
  keyCollabRoomRaw,
  keyCollabRoomPattern,
  parseRoomName,
  keyCollabYjs,
};
