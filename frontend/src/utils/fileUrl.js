import { resolveWithApiBase } from '@/config/baseUrl';

export function resolveFileUrl(rawUrl) {
  const value = String(rawUrl || '').trim();
  if (!value) {
    return '';
  }
  return resolveWithApiBase(value);
}
