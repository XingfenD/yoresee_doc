const trimTrailingSlash = (value = '') => String(value).trim().replace(/\/+$/, '');
const trimLeadingSlash = (value = '') => String(value).trim().replace(/^\/+/, '');

const defaultBaseURL = (() => {
  if (typeof window !== 'undefined' && window.location?.origin) {
    return window.location.origin;
  }
  return '';
})();

export const API_BASE_URL = trimTrailingSlash(import.meta.env.VITE_API_BASE_URL || defaultBaseURL);

export const resolveWithApiBase = (target = '') => {
  const value = String(target || '').trim();
  if (!value) {
    return API_BASE_URL;
  }
  try {
    return new URL(value).toString();
  } catch (_) {
    if (!API_BASE_URL) {
      return value;
    }
    const normalizedPath = value.startsWith('/') ? value : `/${trimLeadingSlash(value)}`;
    return `${API_BASE_URL}${normalizedPath}`;
  }
};

export const GRPC_WEB_ENDPOINT = resolveWithApiBase(import.meta.env.VITE_GRPC_WEB_ENDPOINT || '/grpc');
