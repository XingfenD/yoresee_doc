import { ref, computed } from 'vue';

/**
 * Composable for @mention detection in plain textarea/input elements.
 *
 * Usage:
 *   const { mentionState, mentionedUsers, onInput, onKeydown, selectMentionUser, clearMentions } = useMentionInput(textareaRef);
 *
 * mentionState: { visible, keyword, position }
 * mentionedUsers: accumulated list of { external_id, nickname } for this input session
 * onInput(event): call on textarea @input
 * onKeydown(event): call on textarea @keydown for arrow/enter/esc handling
 * selectMentionUser(user): call when user picks a suggestion
 * clearMentions(): reset mentionedUsers (call after comment submit)
 */
export function useMentionInput(textareaRef) {
  const mentionState = ref({ visible: false, keyword: '', position: { x: 0, y: 0 }, activeIndex: 0 });
  const mentionedUsers = ref([]);

  // Track the start position of the current @query in the textarea value
  let mentionStart = -1;

  function onInput(event) {
    const el = event.target;
    const value = el.value;
    const cursor = el.selectionStart;

    // Find the nearest @ before cursor on the same line
    let atPos = -1;
    for (let i = cursor - 1; i >= 0; i--) {
      const ch = value[i];
      if (ch === '@') {
        atPos = i;
        break;
      }
      if (ch === '\n' || ch === ' ') {
        break;
      }
    }

    if (atPos === -1) {
      mentionStart = -1;
      mentionState.value.visible = false;
      return;
    }

    const keyword = value.slice(atPos + 1, cursor);
    mentionStart = atPos;
    const rect = getCaretCoords(el, atPos);
    mentionState.value = {
      visible: true,
      keyword,
      position: { x: rect.x, y: rect.y },
      activeIndex: 0,
    };
  }

  function onKeydown(event) {
    if (!mentionState.value.visible) return;
    if (event.key === 'ArrowDown') {
      event.preventDefault();
      mentionState.value.activeIndex = (mentionState.value.activeIndex + 1) % 8;
    } else if (event.key === 'ArrowUp') {
      event.preventDefault();
      mentionState.value.activeIndex = Math.max(0, mentionState.value.activeIndex - 1);
    } else if (event.key === 'Escape') {
      event.preventDefault();
      mentionState.value.visible = false;
    }
    // Enter is handled by the parent via selectMentionUser
  }

  function selectMentionUser(user) {
    const el = textareaRef.value;
    if (!el || mentionStart === -1) {
      mentionState.value.visible = false;
      return;
    }
    const value = el.value;
    const cursor = el.selectionStart;
    const displayName = user.nickname || user.username;
    const before = value.slice(0, mentionStart);
    const after = value.slice(cursor);
    const inserted = `@${displayName} `;
    el.value = before + inserted + after;
    // Trigger Vue reactivity via input event
    el.dispatchEvent(new Event('input'));
    // Move cursor after inserted text
    const newCursor = mentionStart + inserted.length;
    el.setSelectionRange(newCursor, newCursor);

    // Track mentioned user (deduplicate)
    if (!mentionedUsers.value.find(u => u.external_id === user.external_id)) {
      mentionedUsers.value.push({ external_id: user.external_id, nickname: displayName });
    }

    mentionStart = -1;
    mentionState.value.visible = false;
  }

  function clearMentions() {
    mentionedUsers.value = [];
  }

  return {
    mentionState,
    mentionedUsers,
    onInput,
    onKeydown,
    selectMentionUser,
    clearMentions,
  };
}

/**
 * Approximate caret pixel position for a textarea at a given character offset.
 * Returns { x, y } in viewport coordinates, positioned below the caret line.
 */
function getCaretCoords(el, offset) {
  const rect = el.getBoundingClientRect();
  // Use a hidden mirror div to measure position
  const mirror = document.createElement('div');
  const style = window.getComputedStyle(el);
  for (const prop of [
    'font', 'fontSize', 'fontFamily', 'fontWeight', 'letterSpacing',
    'lineHeight', 'padding', 'border', 'boxSizing', 'wordWrap', 'whiteSpace',
  ]) {
    mirror.style[prop] = style[prop];
  }
  mirror.style.position = 'absolute';
  mirror.style.visibility = 'hidden';
  mirror.style.width = el.offsetWidth + 'px';
  mirror.style.height = 'auto';
  mirror.style.top = '0';
  mirror.style.left = '-9999px';
  mirror.style.overflow = 'hidden';
  mirror.style.whiteSpace = 'pre-wrap';

  const textBefore = el.value.slice(0, offset);
  mirror.textContent = textBefore;
  const span = document.createElement('span');
  span.textContent = '|';
  mirror.appendChild(span);
  document.body.appendChild(mirror);

  const spanRect = span.getBoundingClientRect();
  const mirrorRect = mirror.getBoundingClientRect();
  const lineHeight = parseInt(style.lineHeight) || 20;

  document.body.removeChild(mirror);

  return {
    x: rect.left + (spanRect.left - mirrorRect.left),
    y: rect.top - el.scrollTop + (span.offsetTop) + lineHeight + 4,
  };
}
