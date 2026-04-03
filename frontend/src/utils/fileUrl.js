function joinPath(basePath, targetPath) {
  const left = (basePath || '').replace(/\/+$/, '');
  const right = (targetPath || '').replace(/^\/+/, '');
  if (!left) return `/${right}`;
  if (!right) return left;
  return `${left}/${right}`;
}

function isInternalMinioHost(hostname, currentHostname) {
  const host = (hostname || '').toLowerCase();
  if (!host) return false;
  if (host === 'minio') return true;
  if (host === 'localhost' || host === '127.0.0.1') {
    return host !== (currentHostname || '').toLowerCase();
  }
  return false;
}

export function resolveFileUrl(rawUrl) {
  if (!rawUrl) return '';

  try {
    const parsed = new URL(rawUrl, window.location.origin);
    const publicBase = (import.meta.env.VITE_ATTACHMENT_PUBLIC_BASE_URL || '').trim();

    if (publicBase) {
      const base = new URL(publicBase, window.location.origin);
      const basePath = (base.pathname || '').replace(/\/+$/, '');
      if (parsed.origin === base.origin && basePath && parsed.pathname.startsWith(`${basePath}/`)) {
        return parsed.toString();
      }
      const pathname = joinPath(base.pathname, parsed.pathname);
      return `${base.origin}${pathname}${parsed.search}${parsed.hash}`;
    }

    if (isInternalMinioHost(parsed.hostname, window.location.hostname)) {
      return `${window.location.origin}${parsed.pathname}${parsed.search}${parsed.hash}`;
    }

    return parsed.toString();
  } catch (_) {
    return rawUrl;
  }
}
