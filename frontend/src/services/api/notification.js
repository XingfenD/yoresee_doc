import {
  unaryCall,
  messages,
  notificationClient,
  baseToObject,
  handleResponse,
  mapNotification
} from './shared';

const {
  CreateNotificationRequest,
  ListNotificationsRequest,
  MarkNotificationsReadRequest,
  MarkAllNotificationsReadRequest
} = messages;

export const createNotification = async (data = {}) => {
  const req = new CreateNotificationRequest({
    receiverExternalIds: data.receiver_external_ids || [],
    type: data.type || '',
    title: data.title || '',
    content: data.content || '',
    payloadJson: data.payload_json || ''
  });
  const resp = await unaryCall(notificationClient, 'createNotification', req);
  const base = baseToObject(resp);
  return handleResponse(base, {});
};

export const listNotifications = async (params = {}) => {
  const req = new ListNotificationsRequest({
    page: params.page || 1,
    pageSize: params.page_size || 10,
    status: params.status || undefined
  });
  const resp = await unaryCall(notificationClient, 'listNotifications', req);
  const base = baseToObject(resp);
  return handleResponse(base, {
    notifications: (resp.notifications || []).map(mapNotification),
    total: resp.total ?? 0
  });
};

export const markNotificationsRead = async (externalIds = []) => {
  const req = new MarkNotificationsReadRequest({ externalIds });
  const resp = await unaryCall(notificationClient, 'markNotificationsRead', req);
  const base = baseToObject(resp);
  return handleResponse(base, {});
};

export const markAllNotificationsRead = async () => {
  const req = new MarkAllNotificationsReadRequest({});
  const resp = await unaryCall(notificationClient, 'markAllNotificationsRead', req);
  const base = baseToObject(resp);
  return handleResponse(base, {});
};
