const DEFAULT_ROW_COUNT = 8;
const DEFAULT_COLUMN_COUNT = 4;
const MIN_COLUMN_LEN = 26;

const buildSheetStyle = () => ({
  bgcolor: '#ffffff',
  color: '#0a0a0a',
  align: 'left',
  valign: 'middle',
  textwrap: false,
  strike: false,
  underline: false,
  font: {
    name: 'Arial',
    size: 10,
    bold: false,
    italic: false
  },
  format: 'normal'
});

const createEmptyRows = (rowCount = DEFAULT_ROW_COUNT, colCount = DEFAULT_COLUMN_COUNT) =>
  Array.from({ length: rowCount }, () => Array.from({ length: colCount }, () => ''));

const normalizeRows = (sourceRows) => {
  if (!Array.isArray(sourceRows) || sourceRows.length === 0) {
    return createEmptyRows();
  }
  const width = Math.max(
    DEFAULT_COLUMN_COUNT,
    ...sourceRows.map((row) => (Array.isArray(row) ? row.length : 0))
  );
  return sourceRows.map((row) => {
    const values = Array.isArray(row) ? row : [];
    return Array.from({ length: width }, (_, index) => {
      const value = values[index];
      return value == null ? '' : String(value);
    });
  });
};

const parseRows = (rawContent) => {
  if (!rawContent || !String(rawContent).trim()) {
    return createEmptyRows();
  }
  try {
    const parsed = JSON.parse(rawContent);
    if (Array.isArray(parsed)) {
      return normalizeRows(parsed);
    }
    if (Array.isArray(parsed?.rows)) {
      return normalizeRows(parsed.rows);
    }
  } catch (error) {
    // ignore parse error and fallback
  }
  return createEmptyRows();
};

const serializeRows = (rows) =>
  JSON.stringify({
    type: 'table',
    version: 1,
    rows
  });

const rowsToSheetData = (rows) => {
  const rowCount = Math.max(DEFAULT_ROW_COUNT, rows.length);
  const colCount = Math.max(
    MIN_COLUMN_LEN,
    ...rows.map((row) => (Array.isArray(row) ? row.length : 0))
  );
  const sheetRows = { len: rowCount };
  rows.forEach((row, rowIndex) => {
    const cells = {};
    row.forEach((value, colIndex) => {
      cells[colIndex] = { text: value == null ? '' : String(value) };
    });
    sheetRows[rowIndex] = { cells };
  });
  return {
    name: 'Sheet1',
    freeze: 'A1',
    cols: { len: colCount },
    rows: sheetRows
  };
};

const firstSheet = (data) => {
  if (Array.isArray(data)) {
    return data[0] || {};
  }
  if (data && typeof data === 'object') {
    if (data.rows) {
      return data;
    }
    const candidates = Object.values(data);
    const found = candidates.find((item) => item && typeof item === 'object' && item.rows);
    return found || {};
  }
  return {};
};

const sheetDataToRows = (data) => {
  const sheet = firstSheet(data);
  const rowMap = sheet?.rows && typeof sheet.rows === 'object' ? sheet.rows : {};
  const rowIndexes = Object.keys(rowMap)
    .filter((key) => /^\d+$/.test(key))
    .map((key) => Number(key));

  let maxRow = DEFAULT_ROW_COUNT - 1;
  let maxCol = DEFAULT_COLUMN_COUNT - 1;

  rowIndexes.forEach((rowIndex) => {
    if (rowIndex > maxRow) {
      maxRow = rowIndex;
    }
    const cellMap = rowMap[rowIndex]?.cells || {};
    Object.keys(cellMap)
      .filter((key) => /^\d+$/.test(key))
      .forEach((key) => {
        const colIndex = Number(key);
        if (colIndex > maxCol) {
          maxCol = colIndex;
        }
      });
  });

  const rows = [];
  for (let r = 0; r <= maxRow; r += 1) {
    const current = [];
    const cellMap = rowMap[r]?.cells || {};
    for (let c = 0; c <= maxCol; c += 1) {
      const value = cellMap[c]?.text;
      current.push(value == null ? '' : String(value));
    }
    rows.push(current);
  }
  return rows;
};

export {
  DEFAULT_ROW_COUNT,
  DEFAULT_COLUMN_COUNT,
  MIN_COLUMN_LEN,
  buildSheetStyle,
  createEmptyRows,
  normalizeRows,
  parseRows,
  serializeRows,
  rowsToSheetData,
  firstSheet,
  sheetDataToRows
};
