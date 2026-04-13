const { initRedis, initYRedis, closeRedis, getClient } = require('./client');
const { updateRoomActiveTime, getActiveRooms, checkRoomExists } = require('./room');
const { getListBuffers, appendListBuffer, replaceListBuffer, addDirtyDoc, removeDirtyDoc, isDirtyDoc } = require('./document');

module.exports = {
  initRedis,
  initYRedis,
  closeRedis,
  updateRoomActiveTime,
  getActiveRooms,
  checkRoomExists,
  getListBuffers,
  appendListBuffer,
  replaceListBuffer,
  addDirtyDoc,
  removeDirtyDoc,
  isDirtyDoc,
  get redisClient() {
    return getClient();
  },
};
