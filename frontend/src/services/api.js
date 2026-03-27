export * from './api/knowledgeBase';
export * from './api/document';
export * from './api/template';
export * from './api/membership';
export * from './api/invitation';
export * from './api/setting';
export * from './api/notification';
export * from './api/comment';

import * as knowledgeBaseApi from './api/knowledgeBase';
import * as documentApi from './api/document';
import * as templateApi from './api/template';
import * as membershipApi from './api/membership';
import * as invitationApi from './api/invitation';
import * as settingApi from './api/setting';
import * as notificationApi from './api/notification';
import * as commentApi from './api/comment';

export default {
  ...knowledgeBaseApi,
  ...documentApi,
  ...templateApi,
  ...membershipApi,
  ...invitationApi,
  ...settingApi,
  ...notificationApi,
  ...commentApi
};
