import { watch } from 'vue';

const APP_NAME = 'Yoresee';

/**
 * Sets document.title reactively based on a page label and an entity name.
 *
 * @param {import('vue').Ref<string>|import('vue').ComputedRef<string>} pageLabel
 *   Translated page type, e.g. t('pageTitle.documentEditor')
 * @param {import('vue').Ref<string>|import('vue').ComputedRef<string>} entityName
 *   Dynamic name of the current entity (document title, KB name, etc.)
 */
export function usePageTitle(pageLabel, entityName) {
  const update = () => {
    const label = pageLabel.value;
    const name = entityName.value;
    if (label && name) {
      document.title = `${label} - ${name} - ${APP_NAME}`;
    } else if (label) {
      document.title = `${label} - ${APP_NAME}`;
    } else {
      document.title = APP_NAME;
    }
  };

  watch([pageLabel, entityName], update, { immediate: true });
}
