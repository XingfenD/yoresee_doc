import { resolveWithApiBase } from '@/config/baseUrl';

const STORAGE_OBJECT_PREFIXES = ['avatars/', 'attachments/'];

const toStorageObjectPath = (rawPath = '') => {
  const normalized = String(rawPath || '')
    .trim()
    .replace(/^\/+/, '');
  if (!normalized) {
    return '';
  }

  if (normalized.startsWith('storage/')) {
    return normalized.slice('storage/'.length);
  }

  for (const prefix of STORAGE_OBJECT_PREFIXES) {
    if (normalized.startsWith(prefix)) {
      return normalized;
    }
  }

  return '';
};

const toStorageUrl = (objectPath = '') => {
  const normalized = String(objectPath || '').trim().replace(/^\/+/, '');
  if (!normalized) {
    return '';
  }
  return resolveWithApiBase(`/storage/${normalized}`);
};

const extractObjectPathFromAbsoluteUrl = (parsedUrl) => {
  const directPath = toStorageObjectPath(parsedUrl?.pathname || '');
  if (directPath) {
    return directPath;
  }

  const segments = String(parsedUrl?.pathname || '')
    .trim()
    .replace(/^\/+/, '')
    .split('/')
    .filter(Boolean);
  if (!segments.length) {
    return '';
  }

  const startIndex = segments.findIndex((segment) => segment === 'avatars' || segment === 'attachments');
  if (startIndex < 0 || startIndex >= segments.length) {
    return '';
  }
  return segments.slice(startIndex).join('/');
};

export function resolveFileUrl(rawUrl) {
  const value = String(rawUrl || '').trim();
  if (!value) {
    return '';
  }

  try {
    const parsed = new URL(value);
    const objectPath = extractObjectPathFromAbsoluteUrl(parsed);
    if (objectPath) {
      return toStorageUrl(objectPath);
    }
    return parsed.toString();
  } catch (_) {
    const objectPath = toStorageObjectPath(value);
    if (objectPath) {
      return toStorageUrl(objectPath);
    }
    return resolveWithApiBase(value);
  }
}
