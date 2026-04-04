const CACHE_STORAGE_KEY = 'topNavDisplayMenusCacheV1';
const CACHE_TTL_MS = 24 * 60 * 60 * 1000;
const memoryCache = new Map();

const resolveCurrentUserKey = () => {
  if (typeof window === 'undefined') return 'guest';
  try {
    const raw = localStorage.getItem('userInfo');
    if (!raw) return 'guest';
    const userInfo = JSON.parse(raw);
    return String(userInfo?.external_id || userInfo?.externalId || 'guest').trim() || 'guest';
  } catch {
    return 'guest';
  }
};

const buildScopedKey = () => `${resolveCurrentUserKey()}::top-nav`;

const readStorageCache = () => {
  if (typeof window === 'undefined') return {};
  try {
    const raw = localStorage.getItem(CACHE_STORAGE_KEY);
    if (!raw) return {};
    const parsed = JSON.parse(raw);
    if (!parsed || typeof parsed !== 'object') return {};
    return parsed;
  } catch {
    return {};
  }
};

const writeStorageCache = (cacheObj) => {
  if (typeof window === 'undefined') return;
  try {
    localStorage.setItem(CACHE_STORAGE_KEY, JSON.stringify(cacheObj));
  } catch {
    // ignore storage write errors
  }
};

const normalizeMenus = (menus) => {
  if (!Array.isArray(menus)) return [];
  return menus
    .map((item) => String(item || '').trim())
    .filter(Boolean);
};

const isEntryExpired = (entry, now = Date.now()) => {
  const expiresAt = Number(entry?.expiresAt || 0);
  return !Number.isFinite(expiresAt) || expiresAt <= now;
};

const normalizeEntry = (entry) => {
  if (Array.isArray(entry)) {
    // Backward compatibility for old schema: treat as expired.
    return null;
  }
  if (!entry || typeof entry !== 'object') return null;
  const menus = normalizeMenus(entry.menus);
  const expiresAt = Number(entry.expiresAt || 0);
  if (!Number.isFinite(expiresAt)) return null;
  return {
    menus,
    expiresAt
  };
};

const purgeExpiredStorageEntries = (storageCache, now = Date.now()) => {
  let changed = false;
  Object.keys(storageCache || {}).forEach((key) => {
    const normalized = normalizeEntry(storageCache[key]);
    if (!normalized || isEntryExpired(normalized, now)) {
      delete storageCache[key];
      changed = true;
    }
  });
  if (changed) {
    writeStorageCache(storageCache);
  }
};

export const getTopNavDisplayMenusCache = () => {
  const key = buildScopedKey();
  const now = Date.now();

  const memoryValue = normalizeEntry(memoryCache.get(key));
  if (memoryValue && !isEntryExpired(memoryValue, now)) {
    return [...memoryValue.menus];
  }
  if (memoryValue && isEntryExpired(memoryValue, now)) {
    memoryCache.delete(key);
  }

  const storageCache = readStorageCache();
  purgeExpiredStorageEntries(storageCache, now);
  const storageEntry = normalizeEntry(storageCache[key]);
  if (!storageEntry || isEntryExpired(storageEntry, now)) {
    return null;
  }

  memoryCache.set(key, storageEntry);
  return [...storageEntry.menus];
};

export const setTopNavDisplayMenusCache = (menus) => {
  const key = buildScopedKey();
  const normalized = normalizeMenus(menus);
  const entry = {
    menus: normalized,
    expiresAt: Date.now() + CACHE_TTL_MS
  };
  memoryCache.set(key, entry);

  const storageCache = readStorageCache();
  purgeExpiredStorageEntries(storageCache);
  storageCache[key] = entry;
  writeStorageCache(storageCache);
};
