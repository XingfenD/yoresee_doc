/**
 * Resolve a collab URL to a full ws:// or wss:// URL.
 * @param {string} rawUrl
 * @returns {string}
 */
export const resolveCollabUrl = (rawUrl) => {
  const input = String(rawUrl || '').trim();
  if (!input) return '';
  if (input.startsWith('ws://') || input.startsWith('wss://')) return input;
  const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
  const host = window.location.host;
  const path = input.startsWith('/') ? input : `/${input}`;
  return `${protocol}://${host}${path}`;
};

/**
 * Get connected peer count from a Yjs WebSocket provider.
 * @param {object} provider
 * @returns {number}
 */
export const getAwarenessPeerCount = (provider) => {
  if (!provider?.awareness) return 0;
  try {
    return provider.awareness.getStates().size;
  } catch (_) {
    return 0;
  }
};
