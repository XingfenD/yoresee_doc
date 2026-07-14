/**
 * Normalize a 2D array to uniform width, coercing all values to strings.
 * @param {Array<Array<*>>} rows
 * @param {number} [minWidth=0] - minimum column count
 * @returns {Array<Array<string>>}
 */
export const normalizeRows = (rows, minWidth = 0) => {
  if (!Array.isArray(rows) || rows.length === 0) return [];
  const width = Math.max(minWidth, 1, ...rows.map((row) => (Array.isArray(row) ? row.length : 0)));
  return rows.map((row) => {
    const values = Array.isArray(row) ? row : [];
    return Array.from({ length: width }, (_, i) => {
      const value = values[i];
      return value === null || value === undefined ? '' : String(value);
    });
  });
};

/**
 * Excel-style column label: 0→A, 25→Z, 26→AA, etc.
 * @param {number} index - 0-based column index
 * @returns {string}
 */
export const columnLabelAt = (index) => {
  let value = Number(index) + 1;
  let label = '';
  while (value > 0) {
    const remainder = (value - 1) % 26;
    label = String.fromCharCode(65 + remainder) + label;
    value = Math.floor((value - 1) / 26);
  }
  return label || 'A';
};
