import { mergeAttributes, Node } from '@tiptap/core';
import { VueNodeViewRenderer } from '@tiptap/vue-3';
import RichTableNodeView from './RichTableNodeView.vue';
import {
  createRichTableModel,
  decodeRichTableModelFromAttr,
  encodeRichTableModelForAttr,
  normalizeRichTableModel
} from './richTableModel';

export const RichTableExtension = Node.create({
  name: 'tableBlock',
  group: 'block',
  atom: true,
  draggable: true,
  selectable: true,

  addAttributes() {
    return {
      table: {
        default: createRichTableModel(),
        parseHTML: (element) =>
          decodeRichTableModelFromAttr(
            element.getAttribute('data-table') || element.getAttribute('table')
          ),
        renderHTML: (attributes) => ({
          'data-table': encodeRichTableModelForAttr(attributes.table),
          'data-type': 'table'
        })
      }
    };
  },

  parseHTML() {
    return [{ tag: 'yoresee-table' }];
  },

  renderHTML({ HTMLAttributes }) {
    return ['yoresee-table', mergeAttributes(HTMLAttributes)];
  },

  addNodeView() {
    return VueNodeViewRenderer(RichTableNodeView);
  },

  addCommands() {
    return {
      insertTableBlock:
        (attrs = {}) =>
        ({ commands }) =>
          commands.insertContent({
            type: this.name,
            attrs: {
              table: normalizeRichTableModel(attrs.table || createRichTableModel(attrs))
            }
          })
    };
  }
});

export const richTableCommandMeta = {
  key: 'table',
  label: '表格',
  command: 'insertTableBlock'
};
