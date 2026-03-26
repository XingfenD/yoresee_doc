import { computed, useSlots } from 'vue';

export function useForwardedSlotNames(excludeNames = []) {
  const slots = useSlots();
  const excludeSet = new Set(excludeNames);

  return computed(() =>
    Object.keys(slots).filter((name) => !excludeSet.has(name))
  );
}
