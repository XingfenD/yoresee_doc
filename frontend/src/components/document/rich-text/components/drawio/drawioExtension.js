import { mergeAttributes, Node } from '@tiptap/core';
import { VueNodeViewRenderer } from '@tiptap/vue-3';
import DrawioNodeView from './DrawioNodeView.vue';

export const DEFAULT_DRAWIO_XML = `<mxfile host="app.diagrams.net" modified="2026-01-01T00:00:00.000Z" agent="yoresee-rich-text" version="24.7.17">
  <diagram id="page-1" name="Page-1">
    <mxGraphModel dx="1142" dy="658" grid="1" gridSize="10" guides="1" tooltips="1" connect="1" arrows="1" fold="1" page="1" pageScale="1" pageWidth="1280" pageHeight="720" background="#ffffff" math="0" shadow="0">
      <root>
        <mxCell id="0" />
        <mxCell id="1" parent="0" />
      </root>
    </mxGraphModel>
  </diagram>
</mxfile>`;

const encodeDiagram = (value) => encodeURIComponent(String(value || ''));
const decodeDiagram = (value) => {
  if (!value) {
    return DEFAULT_DRAWIO_XML;
  }
  try {
    return decodeURIComponent(String(value));
  } catch (_) {
    return String(value) || DEFAULT_DRAWIO_XML;
  }
};

export const DrawioExtension = Node.create({
  name: 'drawioBlock',
  group: 'block',
  atom: true,
  draggable: true,
  selectable: true,

  addAttributes() {
    return {
      diagram: {
        default: DEFAULT_DRAWIO_XML,
        parseHTML: (element) => decodeDiagram(
          element.getAttribute('data-diagram') ||
          element.getAttribute('diagram')
        ),
        renderHTML: (attributes) => ({
          'data-diagram': encodeDiagram(attributes.diagram || '')
        })
      }
    };
  },

  parseHTML() {
    return [{ tag: 'yoresee-drawio' }];
  },

  renderHTML({ HTMLAttributes }) {
    return ['yoresee-drawio', mergeAttributes(HTMLAttributes, { 'data-type': 'drawio' })];
  },

  addNodeView() {
    return VueNodeViewRenderer(DrawioNodeView);
  },

  addCommands() {
    return {
      insertDrawioBlock:
        (attrs = {}) =>
        ({ commands }) =>
          commands.insertContent({
            type: this.name,
            attrs: {
              diagram: attrs.diagram || DEFAULT_DRAWIO_XML
            }
          })
    };
  }
});

export const drawioCommandMeta = {
  key: 'drawio',
  label: 'Draw.io',
  command: 'insertDrawioBlock'
};
