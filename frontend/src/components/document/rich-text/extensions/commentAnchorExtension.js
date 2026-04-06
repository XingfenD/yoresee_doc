import { Mark, mergeAttributes } from '@tiptap/core';

export const COMMENT_ANCHOR_ATTR = 'data-comment-anchor-id';

export const CommentAnchorExtension = Mark.create({
  name: 'commentAnchor',
  inclusive: false,

  addAttributes() {
    return {
      anchorId: {
        default: '',
        parseHTML: (element) => String(element.getAttribute(COMMENT_ANCHOR_ATTR) || '').trim(),
        renderHTML: (attributes) => {
          const anchorId = String(attributes.anchorId || '').trim();
          if (!anchorId) {
            return {};
          }
          return {
            [COMMENT_ANCHOR_ATTR]: anchorId
          };
        }
      }
    };
  },

  parseHTML() {
    return [
      {
        tag: `span[${COMMENT_ANCHOR_ATTR}]`
      }
    ];
  },

  renderHTML({ HTMLAttributes }) {
    return [
      'span',
      mergeAttributes(HTMLAttributes, {
        class: 'yoresee-comment-anchor'
      }),
      0
    ];
  },

  addCommands() {
    return {
      setCommentAnchor:
        (anchorId) =>
        ({ commands }) =>
          commands.setMark(this.name, { anchorId }),
      unsetCommentAnchor:
        () =>
        ({ commands }) =>
          commands.unsetMark(this.name)
    };
  }
});
