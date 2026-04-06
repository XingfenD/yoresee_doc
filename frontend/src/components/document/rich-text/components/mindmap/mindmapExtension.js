import { mergeAttributes, Node } from '@tiptap/core';
import { VueNodeViewRenderer } from '@tiptap/vue-3';
import MindmapNodeView from './MindmapNodeView.vue';

export const DEFAULT_MINDMAP_SOURCE = `# Mindmap
## Topic A
### Detail A1
## Topic B`;

const encodeSource = (value) => encodeURIComponent(String(value || ''));
const decodeSource = (value) => {
  if (!value) {
    return DEFAULT_MINDMAP_SOURCE;
  }
  try {
    return decodeURIComponent(String(value));
  } catch (_) {
    return String(value);
  }
};

export const MindmapExtension = Node.create({
  name: 'mindmapBlock',
  group: 'block',
  atom: true,
  draggable: true,
  selectable: true,

  addAttributes() {
    return {
      source: {
        default: DEFAULT_MINDMAP_SOURCE,
        parseHTML: (element) => decodeSource(
          element.getAttribute('data-source') || element.getAttribute('source')
        ),
        renderHTML: (attributes) => ({
          'data-source': encodeSource(attributes.source || DEFAULT_MINDMAP_SOURCE)
        })
      }
    };
  },

  parseHTML() {
    return [{ tag: 'yoresee-mindmap' }];
  },

  renderHTML({ HTMLAttributes }) {
    return ['yoresee-mindmap', mergeAttributes(HTMLAttributes)];
  },

  addNodeView() {
    return VueNodeViewRenderer(MindmapNodeView);
  },

  addCommands() {
    return {
      insertMindmapBlock:
        (attrs = {}) =>
        ({ commands }) =>
          commands.insertContent({
            type: this.name,
            attrs: {
              source: attrs.source || DEFAULT_MINDMAP_SOURCE
            }
          })
    };
  }
});

export const mindmapCommandMeta = {
  key: 'mindmap',
  label: 'Mindmap',
  command: 'insertMindmapBlock'
};
