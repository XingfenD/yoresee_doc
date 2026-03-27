import { unaryCall, messages, settingClient, baseToObject, handleResponse } from './shared';

const { GetSettingsRequest, UpdateSettingsRequest } = messages;

export const getSettings = async (scene = 'system') => {
  const req = new GetSettingsRequest({ scene });
  const resp = await unaryCall(settingClient, 'getSettings', req);
  const base = baseToObject(resp);
  return handleResponse(base, {
    groups: (resp.groups || []).map((group) => ({
      key: group.key,
      title: group.title,
      title_key: group.titleKey,
      items: (group.items || []).map((item) => ({
        key: item.key,
        label: item.label,
        label_key: item.labelKey,
        description: item.description,
        description_key: item.descriptionKey,
        type: item.type,
        ui: {
          component: item.ui?.component || '',
          options: (item.ui?.options || []).map((opt) => ({
            label: opt.label,
            label_key: opt.labelKey,
            value: opt.value
          })),
          placeholder: item.ui?.placeholder || '',
          placeholder_key: item.ui?.placeholderKey || ''
        },
        value: item.value,
        default_value: item.defaultValue,
        required: item.required,
        readonly: item.readonly
      }))
    }))
  });
};

export const updateSettings = async (updates = []) => {
  const req = new UpdateSettingsRequest({
    updates: updates.map((item) => ({
      key: item.key,
      value: item.value
    }))
  });
  const resp = await unaryCall(settingClient, 'updateSettings', req);
  const base = baseToObject(resp);
  return handleResponse(base, {});
};
