import TurndownService from 'turndown';
import { marked } from 'marked';
import { COMMENT_ANCHOR_ATTR } from '@/components/document/rich-text/extensions/commentAnchorExtension';
import {
  decodeRichTableModelFromAttr,
  decodeRichTableModelFromMarkdownBlock,
  encodeRichTableModelForAttr,
  encodeRichTableModelToMarkdownBlock,
  encodeDrawioSourceToBase64,
  extractDrawioBlocksFromHtml,
  extractLanguageFromCodeClass,
  normalizeDrawioSource,
  resolveDrawioSource,
  resolveMindmapSource
} from './componentCodec';

marked.setOptions({
  gfm: true,
  breaks: true
});

export const createRichTextTurndown = () => {
  const turndown = new TurndownService({
    headingStyle: 'atx',
    bulletListMarker: '-',
    codeBlockStyle: 'fenced',
    emDelimiter: '*'
  });

  turndown.addRule('fencedCodeBlockWithLanguage', {
    filter: (node) => {
      if (!node || node.nodeName !== 'PRE' || !node.firstChild) {
        return false;
      }
      return node.firstChild.nodeName === 'CODE';
    },
    replacement: (_, node) => {
      const codeNode = node.firstChild;
      const className = typeof codeNode.getAttribute === 'function' ? codeNode.getAttribute('class') || '' : '';
      const language = extractLanguageFromCodeClass(className);
      const rawText = String(codeNode.textContent || '').replace(/\n+$/, '');
      const languageSuffix = language ? language : '';
      return `\n\n\`\`\`${languageSuffix}\n${rawText}\n\`\`\`\n\n`;
    }
  });

  turndown.addRule('yoreseeMindmap', {
    filter: (node) => {
      if (!node) {
        return false;
      }
      const nodeName = String(node.nodeName || '').toLowerCase();
      if (nodeName === 'yoresee-mindmap') {
        return true;
      }
      if (typeof node.getAttribute !== 'function') {
        return false;
      }
      const dataType = String(node.getAttribute('data-type') || '').toLowerCase();
      const hasDataSource = Boolean(node.getAttribute('data-source') || node.getAttribute('source'));
      return hasDataSource && dataType.includes('mindmap');
    },
    replacement: (_, node) => {
      const source = resolveMindmapSource(node).trim();
      if (!source) {
        return '';
      }
      return `\n\n\`\`\`mindmap\n${source}\n\`\`\`\n\n`;
    }
  });

  turndown.addRule('yoreseeDrawio', {
    filter: (node) => {
      if (!node) {
        return false;
      }
      const nodeName = String(node.nodeName || '').toLowerCase();
      if (nodeName === 'yoresee-drawio') {
        return true;
      }
      if (typeof node.getAttribute !== 'function') {
        return false;
      }
      const dataType = String(node.getAttribute('data-type') || '').toLowerCase();
      const hasData = Boolean(node.getAttribute('data-diagram') || node.getAttribute('diagram'));
      return hasData && dataType.includes('drawio');
    },
    replacement: (_, node) => {
      const source = resolveDrawioSource(node).trim();
      if (!source) {
        return '';
      }
      const encoded = encodeDrawioSourceToBase64(source);
      if (!encoded) {
        return '';
      }
      return `\n\n\`\`\`drawio\nbase64:${encoded}\n\`\`\`\n\n`;
    }
  });

  turndown.addRule('yoreseeTable', {
    filter: (node) => {
      if (!node) {
        return false;
      }
      const nodeName = String(node.nodeName || '').toLowerCase();
      if (nodeName === 'yoresee-table') {
        return true;
      }
      if (typeof node.getAttribute !== 'function') {
        return false;
      }
      const dataType = String(node.getAttribute('data-type') || '').toLowerCase();
      const hasData = Boolean(node.getAttribute('data-table') || node.getAttribute('table'));
      return hasData && dataType.includes('table');
    },
    replacement: (_, node) => {
      const tableModel = decodeRichTableModelFromAttr(
        node.getAttribute('data-table') || node.getAttribute('table')
      );
      const encoded = encodeRichTableModelToMarkdownBlock(tableModel);
      if (!encoded) {
        return '';
      }
      return `\n\n\`\`\`table\n${encoded}\n\`\`\`\n\n`;
    }
  });

  turndown.addRule('yoreseeCommentAnchor', {
    filter: (node) => node.nodeName === 'SPAN' && node.getAttribute(COMMENT_ANCHOR_ATTR),
    replacement: (content, node) => {
      const anchorId = String(node.getAttribute(COMMENT_ANCHOR_ATTR) || '').trim();
      if (!anchorId) {
        return content || '';
      }
      return `<span ${COMMENT_ANCHOR_ATTR}="${anchorId}">${content || ''}</span>`;
    }
  });

  return turndown;
};

export const markdownToHtml = (value) => {
  const markdownSource = String(value || '');
  if (!markdownSource.trim()) {
    return '<p></p>';
  }

  const sourceWithComponentBlocks = markdownSource
    .replace(/```drawio\s*\n([\s\S]*?)```/gi, (_, drawioSource) => {
      const normalized = normalizeDrawioSource(drawioSource);
      if (!normalized) {
        return '';
      }
      const encoded = encodeURIComponent(normalized);
      return `\n<yoresee-drawio data-diagram="${encoded}" data-type="drawio"></yoresee-drawio>\n`;
    })
    .replace(/```mindmap\s*\n([\s\S]*?)```/gi, (_, mindmapSource) => {
      const encoded = encodeURIComponent(String(mindmapSource || '').trim());
      return `\n<yoresee-mindmap data-source="${encoded}"></yoresee-mindmap>\n`;
    })
    .replace(/```table\s*\n([\s\S]*?)```/gi, (_, tableSource) => {
      const tableModel = decodeRichTableModelFromMarkdownBlock(tableSource);
      if (!tableModel) {
        return '';
      }
      const encoded = encodeRichTableModelForAttr(tableModel);
      return `\n<yoresee-table data-table="${encoded}" data-type="table"></yoresee-table>\n`;
    });

  const parsed = marked.parse(sourceWithComponentBlocks, { async: false });
  return typeof parsed === 'string' ? parsed : markdownSource;
};

export const htmlToMarkdown = (value, turndown) => {
  const source = String(value || '');
  if (!source.trim()) {
    return '';
  }
  const drawioFallbackBlocks = extractDrawioBlocksFromHtml(source);
  let markdown = turndown.turndown(source);
  if (drawioFallbackBlocks.length > 0 && !/```drawio[\s\S]*?```/i.test(markdown)) {
    markdown = `${markdown.trimEnd()}\n\n${drawioFallbackBlocks.join('\n\n')}`.trim();
  }
  return markdown.replace(/\n{3,}/g, '\n\n').trimEnd();
};
