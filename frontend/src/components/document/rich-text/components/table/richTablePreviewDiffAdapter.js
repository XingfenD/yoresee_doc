import { normalizeRichTableModel } from './richTableModel';

const renderTableSummary = (node) => {
  const model = normalizeRichTableModel(node?.attrs?.table);
  const rows = model.rows || [];
  const rowCount = rows.length;
  const colCount = rows[0]?.length || 0;

  const previewRows = rows.slice(0, 4).map((row) => {
    const cells = row.slice(0, 4).map((cell) => String(cell || '').trim() || ' ');
    return `| ${cells.join(' | ')} |`;
  });

  const lines = [`[Table] ${rowCount}x${colCount}`];
  if (previewRows.length > 0) {
    lines.push(...previewRows);
  }
  return lines.join('\n');
};

export const richTablePreviewDiffAdapter = {
  toPreview: renderTableSummary,
  toDiff: renderTableSummary
};
