import { createApp, h } from 'vue';
import MentionUserList from '@/components/shared/MentionUserList.vue';
import { listUsers } from '@/services/api/membership.js';

/**
 * Tiptap suggestion configuration for @mention.
 */
export const mentionSuggestion = {
  char: '@',
  allowSpaces: false,
  startOfLine: false,

  items: async ({ query }) => {
    try {
      const result = await listUsers({ keyword: query, page: 1, page_size: 8 });
      return result.users || [];
    } catch {
      return [];
    }
  },

  render: () => {
    let app = null;
    let mountEl = null;
    let currentProps = {};

    function getPosition(clientRect) {
      const rect = typeof clientRect === 'function' ? clientRect() : clientRect;
      if (!rect) return { x: 0, y: 0 };
      return { x: rect.left, y: rect.bottom + 4 };
    }

    function remount() {
      if (app) {
        app.unmount();
        app = null;
      }
      if (!mountEl) {
        mountEl = document.createElement('div');
        document.body.appendChild(mountEl);
      }

      const props = { ...currentProps };
      app = createApp({
        render: () => h(MentionUserList, props),
      });
      app.mount(mountEl);
    }

    return {
      onStart(props) {
        currentProps = {
          visible: true,
          keyword: props.query || '',
          position: getPosition(props.clientRect),
          activeIndex: 0,
          onSelect: (user) => props.command({ id: user.external_id, label: user.nickname || user.username }),
        };
        remount();
      },

      onUpdate(props) {
        currentProps = {
          visible: true,
          keyword: props.query || '',
          position: getPosition(props.clientRect),
          activeIndex: currentProps.activeIndex || 0,
          onSelect: (user) => props.command({ id: user.external_id, label: user.nickname || user.username }),
        };
        remount();
      },

      onKeyDown({ event }) {
        if (event.key === 'Escape') {
          currentProps = { ...currentProps, visible: false };
          remount();
          return true;
        }
        if (event.key === 'ArrowDown') {
          currentProps = { ...currentProps, activeIndex: (currentProps.activeIndex + 1) % 8 };
          remount();
          return true;
        }
        if (event.key === 'ArrowUp') {
          currentProps = { ...currentProps, activeIndex: Math.max(0, currentProps.activeIndex - 1) };
          remount();
          return true;
        }
        if (event.key === 'Enter') {
          // Handled by MentionUserList click; skip here to avoid double submit
          return false;
        }
        return false;
      },

      onExit() {
        if (app) {
          app.unmount();
          app = null;
        }
        if (mountEl) {
          document.body.removeChild(mountEl);
          mountEl = null;
        }
      },
    };
  },
};
