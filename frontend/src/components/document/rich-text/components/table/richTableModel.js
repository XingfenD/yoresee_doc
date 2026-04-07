const DEFAULT_ROWS = 3;
const DEFAULT_COLS = 3;
const MAX_SIDE = 30;

const clampSize = (value, fallback) => {
  const number = Number(value);
  if (!Number.isFinite(number)) {
    return fallback;
  }
  return Math.max(1, Math.min(MAX_SIDE, Math.floor(number)));
};

const toCellText = (value) => String(value ?? '');

const encodeUtf8Base64 = (value) => {
  const source = String(value || '');
  if (!source) {
    return '';
  }
  try {
    if (typeof TextEncoder !== 'undefined') {
      const bytes = new TextEncoder().encode(source);
      let binary = '';
      const chunkSize = 0x8000;
      for (let index = 0; index < bytes.length; index += chunkSize) {
        const chunk = bytes.subarray(index, index + chunkSize);
        binary += String.fromCharCode(...chunk);
      }
      return btoa(binary);
    }
    return btoa(unescape(encodeURIComponent(source)));
  } catch (_) {
    return '';
  }
};

const decodeUtf8Base64 = (value) => {
  const source = String(value || '').trim();
  if (!source) {
    return '';
  }
  try {
    const binary = atob(source);
    if (typeof TextDecoder !== 'undefined') {
      const bytes = Uint8Array.from(binary, (char) => char.charCodeAt(0));
      return new TextDecoder().decode(bytes);
    }
    return decodeURIComponent(escape(binary));
  } catch (_) {
    return '';
  }
};

export const createRichTableModel = ({ rows = DEFAULT_ROWS, cols = DEFAULT_COLS } = {}) => {
  const rowCount = clampSize(rows, DEFAULT_ROWS);
  const colCount = clampSize(cols, DEFAULT_COLS);
  return {
    version: 1,
    rows: Array.from({ length: rowCount }, () =>
      Array.from({ length: colCount }, () => '')
    )
  };
};

export const normalizeRichTableModel = (value, fallback = {}) => {
  const fallbackModel = createRichTableModel(fallback);
  if (!value || typeof value !== 'object') {
    return fallbackModel;
  }

  const sourceRows = Array.isArray(value.rows) ? value.rows : [];
  if (sourceRows.length === 0) {
    return fallbackModel;
  }

  const normalizedRowCount = clampSize(sourceRows.length, fallbackModel.rows.length);
  const rawMaxCols = sourceRows.reduce((max, row) => {
    const size = Array.isArray(row) ? row.length : 0;
    return Math.max(max, size);
  }, 0);
  const normalizedColCount = clampSize(
    rawMaxCols || fallbackModel.rows[0]?.length || DEFAULT_COLS,
    fallbackModel.rows[0]?.length || DEFAULT_COLS
  );

  const rows = Array.from({ length: normalizedRowCount }, (_, rowIndex) => {
    const sourceRow = Array.isArray(sourceRows[rowIndex]) ? sourceRows[rowIndex] : [];
    return Array.from({ length: normalizedColCount }, (_, colIndex) =>
      toCellText(sourceRow[colIndex] || '')
    );
  });

  return {
    version: 1,
    rows
  };
};

export const cloneRichTableModel = (value) => normalizeRichTableModel(value);

export const encodeRichTableModelForAttr = (value) =>
  encodeURIComponent(JSON.stringify(normalizeRichTableModel(value)));

export const decodeRichTableModelFromAttr = (value) => {
  const source = String(value || '').trim();
  if (!source) {
    return createRichTableModel();
  }
  try {
    return normalizeRichTableModel(JSON.parse(decodeURIComponent(source)));
  } catch (_) {
    return createRichTableModel();
  }
};

export const encodeRichTableModelToMarkdownBlock = (value) => {
  const payload = JSON.stringify(normalizeRichTableModel(value));
  const encoded = encodeUtf8Base64(payload);
  return encoded ? `base64:${encoded}` : '';
};

export const decodeRichTableModelFromMarkdownBlock = (value) => {
  const source = String(value || '').trim();
  if (!source) {
    return null;
  }

  const jsonSource = source.startsWith('base64:')
    ? decodeUtf8Base64(source.slice('base64:'.length))
    : source;

  if (!jsonSource) {
    return null;
  }

  try {
    return normalizeRichTableModel(JSON.parse(jsonSource));
  } catch (_) {
    return null;
  }
};
