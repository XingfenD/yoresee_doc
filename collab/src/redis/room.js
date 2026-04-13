const { getClient } = require('./client');
const { keyCollabRoomRaw, keyCollabRoomPattern, parseRoomName, keyCollabYjs } = require('../key');

async function updateRoomActiveTime(roomName) {
  const client = getClient();
  if (!client) return;

  try {
    const key = keyCollabRoomRaw(roomName);
    await client.set(key, Date.now().toString(), { EX: 3600 * 24 });
  } catch (err) {
    console.error('Failed to update room active time:', err);
  }
}

async function getActiveRooms() {
  const client = getClient();
  if (!client) return [];

  try {
    const keys = await client.keys(keyCollabRoomPattern());
    return keys.map(k => parseRoomName(k));
  } catch (err) {
    console.error('Failed to get active rooms:', err);
    return [];
  }
}

async function checkRoomExists(roomName) {
  const client = getClient();
  if (!client) return false;

  try {
    const yjsKey = keyCollabYjs(roomName);
    const exists = await client.exists(yjsKey);
    return exists;
  } catch (err) {
    console.error(`Failed to check room ${roomName}:`, err.message);
    return false;
  }
}

module.exports = { updateRoomActiveTime, getActiveRooms, checkRoomExists };
