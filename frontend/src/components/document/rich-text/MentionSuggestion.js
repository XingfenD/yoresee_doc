import { mentionPopupState, hideMentionPopup } from './mentionPopupState.js';

/**
 * Tiptap suggestion configuration for @mention.
 * UI is rendered by MentionUserList inside the main Vue app (via YoreseeRichTextEditor).
 * This file only drives the reactive state — no createApp, no DOM manipulation.
 */
export const mentionSuggestion = {
  char: '@',
  allowSpaces: false,
  startOfLine: false,

  // items() is not used here — fetching is done inside MentionUserList via watch on keyword
  items: () => [],

  render: () => {
    function getPosition(clientRect) {
      const rect = typeof clientRect === 'function' ? clientRect() : clientRect;
      if (!rect) return { x: 0, y: 0 };
      return { x: rect.left, y: rect.bottom + 4 };
    }

    return {
      onStart(props) {
        mentionPopupState.visible = true;
        mentionPopupState.keyword = props.query || '';
        mentionPopupState.position = getPosition(props.clientRect);
        mentionPopupState.activeIndex = 0;
        mentionPopupState.onSelect = (user) =>
          props.command({ id: user.external_id, label: user.nickname || user.username });
      },

      onUpdate(props) {
        mentionPopupState.keyword = props.query || '';
        mentionPopupState.position = getPosition(props.clientRect);
        mentionPopupState.onSelect = (user) =>
          props.command({ id: user.external_id, label: user.nickname || user.username });
      },

      onKeyDown({ event }) {
        if (!mentionPopupState.visible) return false;
        if (event.key === 'Escape') {
          hideMentionPopup();
          return true;
        }
        if (event.key === 'ArrowDown') {
          mentionPopupState.activeIndex = (mentionPopupState.activeIndex + 1) % 8;
          return true;
        }
        if (event.key === 'ArrowUp') {
          mentionPopupState.activeIndex = Math.max(0, mentionPopupState.activeIndex - 1);
          return true;
        }
        return false;
      },

      onExit() {
        hideMentionPopup();
      },
    };
  },
};
