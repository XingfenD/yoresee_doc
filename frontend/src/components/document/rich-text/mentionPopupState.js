import { reactive } from 'vue';

/**
 * Shared reactive state for the mention suggestion popup.
 * Updated by MentionSuggestion.js, read by MentionUserListPortal (rendered inside the main app).
 */
export const mentionPopupState = reactive({
  visible: false,
  keyword: '',
  position: { x: 0, y: 0 },
  activeIndex: 0,
  onSelect: null,  // (user: {external_id, nickname, username}) => void
});

export function hideMentionPopup() {
  mentionPopupState.visible = false;
  mentionPopupState.onSelect = null;
}
