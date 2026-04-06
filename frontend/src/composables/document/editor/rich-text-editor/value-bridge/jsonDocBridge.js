const EMPTY_DOC = Object.freeze({ type: 'doc', content: [] });
const EMPTY_DOC_JSON = '{"type":"doc","content":[]}';

const safeClone = (value) => {
  try {
    return JSON.parse(JSON.stringify(value));
  } catch (_) {
    return { ...EMPTY_DOC };
  }
};

export const normalizeJsonDoc = (value) => {
  if (value && typeof value === 'object' && value.type === 'doc') {
    return safeClone(value);
  }
  return { ...EMPTY_DOC };
};

export const serializeJsonDoc = (value) => {
  try {
    return JSON.stringify(normalizeJsonDoc(value));
  } catch (_) {
    return EMPTY_DOC_JSON;
  }
};

export const cloneJsonDoc = safeClone;
