/**
 * Extract all unique mentioned user external IDs from the Tiptap editor state.
 *
 * @param {import('@tiptap/vue-3').Editor} editor
 * @returns {string[]}
 */
export function extractMentionIds(editor) {
  if (!editor) return [];
  const ids = [];
  editor.state.doc.descendants((node) => {
    if (node.type.name === 'mention' && node.attrs.id) {
      ids.push(node.attrs.id);
    }
  });
  return [...new Set(ids)];
}
