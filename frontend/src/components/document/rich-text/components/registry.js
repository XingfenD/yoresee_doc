import { MindmapExtension, mindmapCommandMeta } from './mindmap/mindmapExtension';
import { DrawioExtension, drawioCommandMeta } from './drawio/drawioExtension';
import { mindmapPreviewDiffAdapter } from './mindmap/mindmapPreviewDiffAdapter';
import { drawioPreviewDiffAdapter } from './drawio/drawioPreviewDiffAdapter';

const COMPONENT_REGISTRY = {
  mindmap: {
    extension: MindmapExtension,
    toolbarItem: mindmapCommandMeta,
    nodeType: 'mindmapBlock',
    previewDiffAdapter: mindmapPreviewDiffAdapter
  },
  drawio: {
    extension: DrawioExtension,
    toolbarItem: drawioCommandMeta,
    nodeType: 'drawioBlock',
    previewDiffAdapter: drawioPreviewDiffAdapter
  }
};

export const DEFAULT_RICH_TEXT_COMPONENTS = ['mindmap', 'drawio'];

export const resolveRichTextComponentSystem = (enabledComponents = DEFAULT_RICH_TEXT_COMPONENTS) => {
  const keys = Array.isArray(enabledComponents) && enabledComponents.length > 0
    ? enabledComponents
    : DEFAULT_RICH_TEXT_COMPONENTS;

  const extensions = [];
  const toolbarItems = [];

  keys.forEach((key) => {
    const component = COMPONENT_REGISTRY[key];
    if (!component) {
      return;
    }
    if (component.extension) {
      extensions.push(component.extension);
    }
    if (component.toolbarItem) {
      toolbarItems.push(component.toolbarItem);
    }
  });

  return {
    extensions,
    toolbarItems
  };
};

export const resolveRichTextPreviewDiffAdapterRegistry = (
  enabledComponents = DEFAULT_RICH_TEXT_COMPONENTS
) => {
  const keys = Array.isArray(enabledComponents) && enabledComponents.length > 0
    ? enabledComponents
    : DEFAULT_RICH_TEXT_COMPONENTS;

  return keys.reduce((result, key) => {
    const component = COMPONENT_REGISTRY[key];
    if (!component?.nodeType || !component?.previewDiffAdapter) {
      return result;
    }
    result[component.nodeType] = component.previewDiffAdapter;
    return result;
  }, {});
};
