import Mention from '@tiptap/extension-mention';
import { mentionSuggestion } from '../MentionSuggestion.js';

/**
 * Tiptap Mention extension configured with:
 * - .mention-chip CSS class for unified styling
 * - data-mention-id attribute for extraction
 * - Plain-text renderText for search indexing
 */
export const MentionExtension = Mention.configure({
  HTMLAttributes: {
    class: 'mention-chip',
  },
  renderLabel({ options, node }) {
    return `${options.suggestion.char}${node.attrs.label ?? node.attrs.id}`;
  },
  suggestion: mentionSuggestion,
}).extend({
  // Override renderHTML to include data-mention-id
  renderHTML({ node, HTMLAttributes }) {
    return [
      'span',
      {
        ...HTMLAttributes,
        'data-mention-id': node.attrs.id,
        'data-mention-label': node.attrs.label,
      },
      `@${node.attrs.label ?? node.attrs.id}`,
    ];
  },
});
