import { clients, messages, unaryCall } from './grpc_client';

const { documentClient, commentClient, knowledgeBaseClient, membershipClient, notificationClient, systemClient, invitationClient, settingClient } = clients;
const {
  ListKnowledgeBasesRequest,
  ListRecentKnowledgeBasesRequest,
  GetKnowledgeBaseRequest,
  UpdateDocumentMetaRequest,
  CreateKnowledgeBaseRequest,
  CreateDocumentRequest,
  CreateTemplateRequest,
  CreateUserGroupRequest,
  DeleteUserGroupRequest,
  GetUserGroupRequest,
  UpdateUserGroupRequest,
  ListUsersRequest,
  UpdateUserRequest,
  ListUserGroupMembersRequest,
  ListTemplatesRequest,
  GetTemplateRequest,
  ListRecentTemplatesRequest,
  ListRecentDocumentsRequest,
  RecordRecentDocumentRequest,
  GetDocumentContentRequest,
  GetOwnDocumentsRequest,
  ListDocumentsRequest,
  ListUserGroupsRequest,
  RecursiveOptions,
  TimeRange,
  CreateDocumentContainerType,
  CreateTemplateContainer,
  DocumentType,
  ListInvitationsRequest,
  CreateInvitationRequest,
  UpdateInvitationRequest,
  DeleteInvitationRequest,
  ListInvitationRecordsRequest,
  ListOrgNodesRequest,
  GetOrgNodeRequest,
  CreateOrgNodeRequest,
  UpdateOrgNodeRequest,
  DeleteOrgNodeRequest,
  MoveOrgNodeRequest,
  ListOrgNodeMembersRequest,
  GetSettingsRequest,
  UpdateSettingsRequest,
  CreateNotificationRequest,
  ListNotificationsRequest,
  MarkNotificationsReadRequest,
  MarkAllNotificationsReadRequest,
  CreateDocumentCommentRequest,
  ListDocumentCommentsRequest,
  DeleteDocumentCommentRequest,
  UpdateDocumentCommentRequest,
  CommentScope
} = messages;

function baseToObject(resp) {
  const base = resp.base;
  return {
    code: base?.code ?? 50000,
    message: base?.message ?? 'unknown error'
  };
}

function mapDocument(doc) {
  if (!doc) return null;
  return {
    external_id: doc.externalId,
    title: doc.title,
    type: doc.type === DocumentType.MARKDOWN ? 'markdown' : '',
    summary: doc.summary,
    status: doc.status,
    tags: doc.tags,
    view_count: doc.viewCount,
    edit_count: doc.editCount,
    created_at: doc.createdAt,
    updated_at: doc.updatedAt,
    has_children: doc.hasChildren,
    children: (doc.children || []).map(mapDocument)
  };
}

function mapKnowledgeBase(kb) {
  if (!kb) return null;
  return {
    external_id: kb.externalId,
    name: kb.name,
    description: kb.description,
    cover: kb.cover,
    is_public: kb.isPublic,
    created_at: kb.createdAt,
    updated_at: kb.updatedAt,
    deleted_at: kb.deletedAt,
    creator_user_external_id: kb.creatorUserExternalId,
    creator_name: kb.creatorName,
    documents_count: kb.documentsCount
  };
}

function mapTemplate(tpl) {
  if (!tpl) return null;
  return {
    id: tpl.id,
    name: tpl.name,
    description: tpl.description,
    content: tpl.content,
    scope: tpl.scope,
    knowledge_base_external_id: tpl.knowledgeBaseExternalId,
    tags: tpl.tags,
    created_at: tpl.createdAt,
    updated_at: tpl.updatedAt
  };
}

function mapUserGroup(group) {
  if (!group) return null;
  return {
    external_id: group.externalId,
    name: group.name,
    description: group.description,
    creator_user_external_id: group.creatorUserExternalId,
    member_count: group.memberCount,
    members: (group.members || []).map(mapUser)
  };
}

function mapUser(user) {
  if (!user) return null;
  return {
    external_id: user.externalId,
    username: user.username,
    email: user.email,
    nickname: user.nickname,
    avatar: user.avatar,
    status: user.status,
    created_at: user.createdAt,
    updated_at: user.updatedAt,
    invitation_code: user.invitationCode ?? null
  };
}

function mapOrgNode(node) {
  if (!node) return null;
  return {
    external_id: node.externalId,
    parent_external_id: node.parentExternalId,
    name: node.name,
    path: node.path,
    description: node.description,
    creator_user_external_id: node.creatorUserExternalId,
    member_count: node.memberCount,
    children: (node.children || []).map(mapOrgNode)
  };
}

function mapInvitation(inv) {
  if (!inv) return null;
  return {
    code: inv.code,
    created_by_external_id: inv.createdByExternalId,
    created_by_name: inv.createdByName,
    used_cnt: inv.usedCnt,
    max_used_cnt: inv.maxUsedCnt ?? null,
    expires_at: inv.expiresAt ?? null,
    created_at: inv.createdAt,
    disabled: inv.disabled,
    note: inv.note ?? null
  };
}

function mapInvitationRecord(rec) {
  if (!rec) return null;
  return {
    code: rec.code,
    used_by: rec.usedBy,
    used_by_external_id: rec.usedByExternalId ?? null,
    used_at: rec.usedAt,
    status: rec.status,
    row_key: `${rec.code || ''}_${rec.usedAt || ''}_${rec.status || ''}`
  };
}

function mapNotification(item) {
  if (!item) return null;
  return {
    external_id: item.externalId,
    type: item.type,
    status: item.status,
    title: item.title,
    content: item.content,
    payload: item.payload,
    created_at: item.createdAt,
    read: item.status === 'read'
  };
}

function mapComment(item) {
  if (!item) return null;
  return {
    external_id: item.externalId,
    document_external_id: item.documentExternalId,
    parent_external_id: item.parentExternalId || null,
    content: item.content,
    anchor_id: item.anchorId || '',
    quote: item.quote || '',
    creator_user_external_id: item.creatorUserExternalId,
    creator_name: item.creatorName,
    creator_avatar: item.creatorAvatar,
    created_at: item.createdAt
  };
}
function handleResponse(base, data) {
  if (base.code === 0) {
    return { ...base, ...data };
  }
  throw new Error('request failed');
}

function toTimeRange(input) {
  if (!input) return null;
  return new TimeRange({
    start: input.start || '',
    end: input.end || ''
  });
}

// 获取知识库列表
export const getKnowledgeBases = async (params = {}) => {
  const req = new ListKnowledgeBasesRequest({
    onlyMine: typeof params.only_mine === 'boolean' ? params.only_mine : undefined,
    nameKeyword: params.name_keyword || undefined,
    isPublic: typeof params.is_public === 'boolean' ? params.is_public : undefined,
    createTimeRange: params.create_time_range ? toTimeRange(params.create_time_range) : undefined,
    updateTimeRange: params.update_time_range ? toTimeRange(params.update_time_range) : undefined,
    orderBy: params.order_by || undefined,
    orderDesc: typeof params.order_desc === 'boolean' ? params.order_desc : undefined,
    page: params.page || undefined,
    pageSize: params.page_size || undefined
  });

  const resp = await unaryCall(knowledgeBaseClient, 'listKnowledgeBases', req);
  const base = baseToObject(resp);
  const data = {
    knowledge_bases: (resp.knowledgeBases || []).map(mapKnowledgeBase),
    total: resp.total ?? 0
  };
  return handleResponse(base, data);
};

// 创建知识库
export const createKnowledgeBase = async (data = {}) => {
  const req = new CreateKnowledgeBaseRequest({
    name: data.name || '',
    description: data.description || '',
    cover: data.cover || '',
    isPublic: Boolean(data.is_public)
  });

  const resp = await unaryCall(knowledgeBaseClient, 'createKnowledgeBase', req);
  const base = baseToObject(resp);
  const dataResp = {
    external_id: resp.externalId
  };
  return handleResponse(base, dataResp);
};

// 获取知识库详情
export const getKnowledgeBaseDetail = async (knowledgeBaseExternalID, params = {}) => {
  const req = new GetKnowledgeBaseRequest({
    knowledgeBaseExternalId: knowledgeBaseExternalID,
    recordRecentLog: Boolean(params.record_recent_log),
    page: params.page || undefined,
    pageSize: params.page_size || undefined,
    directoryOnly: params.directory_only || false
  });

  const resp = await unaryCall(knowledgeBaseClient, 'getKnowledgeBase', req);
  const base = baseToObject(resp);
  const data = {
    knowledge_base: mapKnowledgeBase(resp.knowledgeBase),
    documents: (resp.documents || []).map(mapDocument),
    total_count: resp.totalCount ?? 0
  };
  return handleResponse(base, data);
};

// 获取最近访问的知识库
export const getRecentKnowledgeBases = async (params = {}) => {
  const req = new ListRecentKnowledgeBasesRequest({
    startTime: params.start_time || undefined,
    endTime: params.end_time || undefined,
    page: params.page || undefined,
    pageSize: params.page_size || undefined
  });

  const resp = await unaryCall(knowledgeBaseClient, 'listRecentKnowledgeBases', req);
  const base = baseToObject(resp);
  const data = {
    knowledge_bases: (resp.knowledgeBases || []).map(mapKnowledgeBase),
    total: resp.total ?? 0
  };
  return handleResponse(base, data);
};

// 获取知识库文档列表
export const getKnowledgeBaseDocuments = async (knowledgeBaseExternalID, params = {}) => {
  return getKnowledgeBaseDetail(knowledgeBaseExternalID, {
    record_recent_log: false,
    page: 1,
    page_size: 1000,
    directory_only: params.directory_only || false,
    ...params
  });
};

// 更新文档元数据
export const updateDocumentMeta = async (documentExternalID, data = {}) => {
  const req = new UpdateDocumentMetaRequest({
    externalId: documentExternalID,
    title: data.title ?? undefined,
    summary: data.summary ?? undefined,
    tags: Array.isArray(data.tags) ? data.tags : undefined,
    status: typeof data.status === 'number' ? data.status : undefined
  });

  const resp = await unaryCall(documentClient, 'updateDocumentMeta', req);
  const base = baseToObject(resp);
  return handleResponse(base, {});
};

// 创建文档
export const createDocument = async (data) => {
  const req = new CreateDocumentRequest({
    title: data.title || '',
    type: DocumentType.MARKDOWN
  });
  if (data.type === 'markdown') {
    req.type = DocumentType.MARKDOWN;
  }
  if (data.container_type === 'knowledge_base') {
    req.containerType = CreateDocumentContainerType.KNOWLEDGE_BASE;
    if (data.knowledge_base_external_id) {
      req.knowledgeBaseExternalId = data.knowledge_base_external_id;
    }
  } else if (data.container_type === 'own') {
    req.containerType = CreateDocumentContainerType.OWN;
  }
  if (data.parent_external_id) {
    req.parentExternalId = data.parent_external_id;
  }
  if (data.template_id) {
    req.templateId = data.template_id;
  }

  const resp = await unaryCall(documentClient, 'createDocument', req);
  const base = baseToObject(resp);
  const dataResp = {
    external_id: resp.externalId
  };
  return handleResponse(base, dataResp);
};

// 创建模板
export const createTemplate = async (data = {}) => {
  const containerMap = {
    own: CreateTemplateContainer.OWN_TEMPLATE,
    knowledge_base: CreateTemplateContainer.KNOWLEDGEBASE_TEMPLATE,
    public: CreateTemplateContainer.PUBLIC_TEMPLATE
  };
  const req = new CreateTemplateRequest({
    targetContainer: containerMap[data.target_container] ?? CreateTemplateContainer.OWN_TEMPLATE,
    knowledgeBaseId: data.knowledge_base_id || '',
    templateContent: data.template_content || ''
  });

  const resp = await unaryCall(documentClient, 'createTemplate', req);
  const base = baseToObject(resp);
  return handleResponse(base, {});
};

// 获取模板列表
export const listTemplates = async (params = {}) => {
  const containerMap = {
    own: CreateTemplateContainer.OWN_TEMPLATE,
    knowledge_base: CreateTemplateContainer.KNOWLEDGEBASE_TEMPLATE,
    public: CreateTemplateContainer.PUBLIC_TEMPLATE
  };

  const req = new ListTemplatesRequest({
    onlyMine: Boolean(params.only_mine),
    targetContainer: params.target_container ? containerMap[params.target_container] : undefined,
    knowledgeBaseId: params.knowledge_base_id || undefined,
    nameKeyword: params.name_keyword || undefined,
    orderBy: params.order_by || undefined,
    orderDesc: typeof params.order_desc === 'boolean' ? params.order_desc : undefined,
    page: params.page || undefined,
    pageSize: params.page_size || undefined
  });

  const resp = await unaryCall(documentClient, 'listTemplates', req);
  const base = baseToObject(resp);
  const data = {
    templates: (resp.templates || []).map(mapTemplate),
    total: resp.total ?? 0
  };
  return handleResponse(base, data);
};

// 获取模板详情
export const getTemplate = async (templateId, options = {}) => {
  const req = new GetTemplateRequest({
    templateId: Number(templateId) || 0,
    recordRecentLog: Boolean(options.record_recent_log)
  });

  const resp = await unaryCall(documentClient, 'getTemplate', req);
  const base = baseToObject(resp);
  const data = {
    template: resp.template ? mapTemplate(resp.template) : null
  };
  return handleResponse(base, data);
};

// 获取最近模板
export const listRecentTemplates = async (params = {}) => {
  const req = new ListRecentTemplatesRequest({
    startTime: params.start_time || undefined,
    endTime: params.end_time || undefined,
    page: params.page || undefined,
    pageSize: params.page_size || undefined
  });

  const resp = await unaryCall(documentClient, 'listRecentTemplates', req);
  const base = baseToObject(resp);
  const data = {
    templates: (resp.templates || []).map(mapTemplate),
    total: resp.total ?? 0
  };
  return handleResponse(base, data);
};

// 获取用户组列表（系统管理）
export const listUserGroups = async (params = {}) => {
  const req = new ListUserGroupsRequest({
    keyword: params.keyword || undefined,
    page: params.page || undefined,
    pageSize: params.page_size || undefined
  });

  const resp = await unaryCall(membershipClient, 'listUserGroups', req);
  const base = baseToObject(resp);
  const data = {
    user_groups: (resp.userGroups || []).map(mapUserGroup),
    total: resp.total ?? 0
  };
  return handleResponse(base, data);
};

export const listInvitations = async (params = {}) => {
  const req = new ListInvitationsRequest({
    maxUsedCnt: typeof params.max_used_cnt === 'number' ? params.max_used_cnt : undefined,
    expiresAtStart: params.expires_at_start || undefined,
    expiresAtEnd: params.expires_at_end || undefined,
    createdAtStart: params.created_at_start || undefined,
    createdAtEnd: params.created_at_end || undefined,
    disabled: typeof params.disabled === 'boolean' ? params.disabled : undefined,
    orderBy: params.order_by || undefined,
    orderDesc: typeof params.order_desc === 'boolean' ? params.order_desc : undefined,
    page: params.page || 1,
    pageSize: params.page_size || 20,
    onlyMine: typeof params.only_mine === 'boolean' ? params.only_mine : undefined,
    keyword: params.keyword || undefined
  });

  const resp = await unaryCall(invitationClient, 'listInvitations', req);
  const base = baseToObject(resp);
  const data = {
    invitations: (resp.invitations || []).map(mapInvitation).filter(Boolean),
    total: resp.total ?? 0
  };
  return handleResponse(base, data);
};

export const createInvitation = async (params = {}) => {
  const req = new CreateInvitationRequest({
    maxUsedCnt: typeof params.max_used_cnt === 'number' ? params.max_used_cnt : undefined,
    expiresAt: params.expires_at || undefined,
    note: params.note || undefined
  });

  const resp = await unaryCall(invitationClient, 'createInvitation', req);
  const base = baseToObject(resp);
  const data = {
    invitation: resp.invitation ? mapInvitation(resp.invitation) : null
  };
  return handleResponse(base, data);
};

export const updateInvitation = async (params = {}) => {
  const req = new UpdateInvitationRequest({
    code: params.code || '',
    maxUsedCnt: typeof params.max_used_cnt === 'number' ? params.max_used_cnt : undefined,
    expiresAt: params.expires_at || undefined,
    disabled: typeof params.disabled === 'boolean' ? params.disabled : undefined,
    note: params.note || undefined
  });

  const resp = await unaryCall(invitationClient, 'updateInvitation', req);
  const base = baseToObject(resp);
  return handleResponse(base, {});
};

export const deleteInvitation = async (code) => {
  const req = new DeleteInvitationRequest({ code: code || '' });

  const resp = await unaryCall(invitationClient, 'deleteInvitation', req);
  const base = baseToObject(resp);
  return handleResponse(base, {});
};

export const listInvitationRecords = async (params = {}) => {
  const req = new ListInvitationRecordsRequest({
    code: params.code || undefined,
    status: params.status || undefined,
    usedAtStart: params.used_at_start || undefined,
    usedAtEnd: params.used_at_end || undefined,
    page: params.page || 1,
    pageSize: params.page_size || 20,
    onlyMine: typeof params.only_mine === 'boolean' ? params.only_mine : undefined,
    keyword: params.keyword || undefined
  });

  const resp = await unaryCall(invitationClient, 'listInvitationRecords', req);
  const base = baseToObject(resp);
  const data = {
    records: (resp.records || []).map(mapInvitationRecord).filter(Boolean),
    total: resp.total ?? 0
  };
  return handleResponse(base, data);
};

export const getUserGroup = async (externalId) => {
  const req = new GetUserGroupRequest({
    externalId: externalId || ''
  });

  const resp = await unaryCall(membershipClient, 'getUserGroup', req);
  const base = baseToObject(resp);
  const data = {
    user_group: resp.userGroup ? mapUserGroup(resp.userGroup) : null
  };
  return handleResponse(base, data);
};

export const createUserGroup = async (data = {}) => {
  const req = new CreateUserGroupRequest({
    name: data.name || '',
    description: data.description || '',
    memberUserExternalIds: Array.isArray(data.member_user_external_ids)
      ? data.member_user_external_ids
      : []
  });

  const resp = await unaryCall(membershipClient, 'createUserGroup', req);
  const base = baseToObject(resp);
  const dataResp = {
    external_id: resp.externalId
  };
  return handleResponse(base, dataResp);
};

export const deleteUserGroup = async (externalId) => {
  const req = new DeleteUserGroupRequest({
    externalId: externalId || ''
  });

  const resp = await unaryCall(membershipClient, 'deleteUserGroup', req);
  const base = baseToObject(resp);
  return handleResponse(base, {});
};

export const updateUserGroup = async (data = {}) => {
  const req = new UpdateUserGroupRequest({
    externalId: data.external_id || '',
    name: data.name ?? undefined,
    description: data.description ?? undefined,
    syncMembers: Boolean(data.sync_members),
    memberUserExternalIds: Array.isArray(data.member_user_external_ids)
      ? data.member_user_external_ids
      : []
  });

  const resp = await unaryCall(membershipClient, 'updateUserGroup', req);
  const base = baseToObject(resp);
  return handleResponse(base, {});
};

export const listUsers = async (params = {}) => {
  const req = new ListUsersRequest({
    keyword: params.keyword || undefined,
    page: params.page || undefined,
    pageSize: params.page_size || undefined
  });

  const resp = await unaryCall(membershipClient, 'listUsers', req);
  const base = baseToObject(resp);
  const data = {
    users: (resp.users || []).map(mapUser),
    total: resp.total ?? 0
  };
  return handleResponse(base, data);
};

export const updateUser = async (data = {}) => {
  const req = new UpdateUserRequest({
    externalId: data.external_id || '',
    username: data.username ?? undefined,
    email: data.email ?? undefined,
    nickname: data.nickname ?? undefined,
    status: typeof data.status === 'number' ? data.status : undefined
  });

  const resp = await unaryCall(membershipClient, 'updateUser', req);
  const base = baseToObject(resp);
  return handleResponse(base, {});
};

export const listUserGroupMembers = async (params = {}) => {
  const req = new ListUserGroupMembersRequest({
    externalId: params.external_id || '',
    keyword: params.keyword || undefined,
    page: params.page || undefined,
    pageSize: params.page_size || undefined
  });

  const resp = await unaryCall(membershipClient, 'listUserGroupMembers', req);
  const base = baseToObject(resp);
  const data = {
    users: (resp.users || []).map(mapUser),
    total: resp.total ?? 0
  };
  return handleResponse(base, data);
};

// 获取文档内容
export const getDocumentContent = async (documentExternalID, params = {}) => {
  const req = new GetDocumentContentRequest({
    documentExternalId: documentExternalID
  });

  const resp = await unaryCall(documentClient, 'getDocumentContent', req);
  const base = baseToObject(resp);
  const data = {
    document: mapDocument(resp.document),
    content: resp.content
  };
  return handleResponse(base, data);
};

// 获取最近文档
export const getRecentDocuments = async (params = {}) => {
  const req = new ListRecentDocumentsRequest({
    startTime: params.start_time || undefined,
    endTime: params.end_time || undefined,
    page: params.page || undefined,
    pageSize: params.page_size || undefined
  });

  const resp = await unaryCall(documentClient, 'listRecentDocuments', req);
  const base = baseToObject(resp);
  const data = {
    documents: (resp.documents || []).map(mapDocument),
    total: resp.total ?? 0
  };
  return handleResponse(base, data);
};

export const recordRecentDocument = async (documentExternalID) => {
  const req = new RecordRecentDocumentRequest({
    documentExternalId: documentExternalID
  });

  const resp = await unaryCall(documentClient, 'recordRecentDocument', req);
  const base = baseToObject(resp);
  return handleResponse(base, {});
};

// 获取我的文档列表
export const getMyDocuments = async (params = {}) => {
  const req = new GetOwnDocumentsRequest({
    page: params.page || undefined,
    pageSize: params.page_size || undefined,
    directoryOnly: params.directory_only || false
  });

  const resp = await unaryCall(documentClient, 'getOwnDocuments', req);
  const base = baseToObject(resp);
  const data = {
    documents: (resp.documents || []).map(mapDocument),
    total_count: resp.totalCount ?? 0
  };
  return handleResponse(base, data);
};

// 额外暴露 ListDocuments 供未来使用
export const listDocuments = async (params = {}) => {
  const options = params.options
    ? new RecursiveOptions({
        includeChildren: typeof params.options.include_children === 'boolean' ? params.options.include_children : undefined,
        recursive: typeof params.options.recursive === 'boolean' ? params.options.recursive : undefined,
        depth: typeof params.options.depth === 'number' ? params.options.depth : undefined
      })
    : undefined;

  const req = new ListDocumentsRequest({
    userExternalId: params.user_external_id || undefined,
    rootDocumentExternalId: params.root_document_external_id || undefined,
    titleKeyword: params.title_keyword || undefined,
    type: params.type || undefined,
    status: typeof params.status === 'number' ? params.status : undefined,
    tags: Array.isArray(params.tags) ? params.tags : undefined,
    createTimeRange: params.create_time_range ? toTimeRange(params.create_time_range) : undefined,
    updateTimeRange: params.update_time_range ? toTimeRange(params.update_time_range) : undefined,
    orderBy: params.order_by || undefined,
    orderDesc: typeof params.order_desc === 'boolean' ? params.order_desc : undefined,
    page: params.page || undefined,
    pageSize: params.page_size || undefined,
    options
  });

  const resp = await unaryCall(documentClient, 'listDocuments', req);
  const base = baseToObject(resp);
  const data = {
    documents: (resp.documents || []).map(mapDocument)
  };
  return handleResponse(base, data);
};

// 获取组织节点列表
export const listOrgNodes = async (params = {}) => {
  const req = new ListOrgNodesRequest({
    parentExternalId: params.parent_external_id || '',
    keyword: params.keyword || undefined,
    page: params.page || undefined,
    pageSize: params.page_size || undefined,
    includeChildren: Boolean(params.include_children)
  });

  const resp = await unaryCall(membershipClient, 'listOrgNodes', req);
  const base = baseToObject(resp);
  const data = {
    org_nodes: (resp.orgNodes || []).map(mapOrgNode),
    total: resp.total ?? 0
  };
  return handleResponse(base, data);
};

// 获取组织节点详情
export const getOrgNode = async (externalId, params = {}) => {
  const req = new GetOrgNodeRequest({
    externalId: externalId || '',
    includeChildren: Boolean(params.include_children)
  });

  const resp = await unaryCall(membershipClient, 'getOrgNode', req);
  const base = baseToObject(resp);
  const data = {
    org_node: resp.orgNode ? mapOrgNode(resp.orgNode) : null
  };
  return handleResponse(base, data);
};

// 创建组织节点
export const createOrgNode = async (data = {}) => {
  const req = new CreateOrgNodeRequest({
    creatorUserExternalId: data.creator_user_external_id || '',
    name: data.name || '',
    description: data.description || '',
    parentExternalId: data.parent_external_id || '',
    memberUserExternalIds: Array.isArray(data.member_user_external_ids)
      ? data.member_user_external_ids
      : []
  });

  const resp = await unaryCall(membershipClient, 'createOrgNode', req);
  const base = baseToObject(resp);
  const dataResp = {
    external_id: resp.externalId
  };
  return handleResponse(base, dataResp);
};

// 更新组织节点
export const updateOrgNode = async (data = {}) => {
  const req = new UpdateOrgNodeRequest({
    externalId: data.external_id || '',
    name: data.name ?? undefined,
    description: data.description ?? undefined,
    syncMembers: Boolean(data.sync_members),
    memberUserExternalIds: Array.isArray(data.member_user_external_ids)
      ? data.member_user_external_ids
      : []
  });

  const resp = await unaryCall(membershipClient, 'updateOrgNode', req);
  const base = baseToObject(resp);
  return handleResponse(base, {});
};

// 删除组织节点
export const deleteOrgNode = async (externalId) => {
  const req = new DeleteOrgNodeRequest({
    externalId: externalId || ''
  });

  const resp = await unaryCall(membershipClient, 'deleteOrgNode', req);
  const base = baseToObject(resp);
  return handleResponse(base, {});
};

// 移动组织节点
export const moveOrgNode = async (data = {}) => {
  const req = new MoveOrgNodeRequest({
    externalId: data.external_id || '',
    newParentExternalId: data.new_parent_external_id || ''
  });

  const resp = await unaryCall(membershipClient, 'moveOrgNode', req);
  const base = baseToObject(resp);
  return handleResponse(base, {});
};

// 获取组织节点成员列表
export const listOrgNodeMembers = async (params = {}) => {
  const req = new ListOrgNodeMembersRequest({
    externalId: params.external_id || '',
    keyword: params.keyword || undefined,
    page: params.page || undefined,
    pageSize: params.page_size || undefined
  });

  const resp = await unaryCall(membershipClient, 'listOrgNodeMembers', req);
  const base = baseToObject(resp);
  const data = {
    users: (resp.users || []).map(mapUser),
    total: resp.total ?? 0
  };
  return handleResponse(base, data);
};

// 获取系统设置
export const getSettings = async (scene = 'system') => {
  const req = new GetSettingsRequest({ scene });
  const resp = await unaryCall(settingClient, 'getSettings', req);
  const base = baseToObject(resp);
  const data = {
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
  };
  return handleResponse(base, data);
};

// 更新系统设置
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

// create notification (manual or system)
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

// list notifications for current user
export const listNotifications = async (params = {}) => {
  const req = new ListNotificationsRequest({
    page: params.page || 1,
    pageSize: params.page_size || 10,
    status: params.status || undefined
  });
  const resp = await unaryCall(notificationClient, 'listNotifications', req);
  const base = baseToObject(resp);
  const data = {
    notifications: (resp.notifications || []).map(mapNotification),
    total: resp.total ?? 0
  };
  return handleResponse(base, data);
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

export const createDocumentComment = async (data = {}) => {
  const req = new CreateDocumentCommentRequest({
    documentExternalId: data.document_external_id || '',
    content: data.content || '',
    parentExternalId: data.parent_external_id || undefined,
    anchorId: data.anchor_id || undefined,
    quote: data.quote || undefined
  });
  const resp = await unaryCall(commentClient, 'createDocumentComment', req);
  const base = baseToObject(resp);
  const dataResp = {
    comment: resp.comment ? mapComment(resp.comment) : null
  };
  return handleResponse(base, dataResp);
};

export const listDocumentComments = async (params = {}) => {
  const req = new ListDocumentCommentsRequest({
    documentExternalId: params.document_external_id || '',
    page: params.page || 1,
    pageSize: params.page_size || 10,
    scope: params.scope ?? CommentScope.COMMENT_SCOPE_ALL
  });
  const resp = await unaryCall(commentClient, 'listDocumentComments', req);
  const base = baseToObject(resp);
  const dataResp = {
    comments: (resp.comments || []).map(mapComment),
    total: resp.total ?? 0
  };
  return handleResponse(base, dataResp);
};

export const deleteDocumentComment = async (externalId) => {
  const req = new DeleteDocumentCommentRequest({
    externalId: externalId || ''
  });
  const resp = await unaryCall(commentClient, 'deleteDocumentComment', req);
  const base = baseToObject(resp);
  return handleResponse(base, {});
};

export const updateDocumentComment = async (data = {}) => {
  const req = new UpdateDocumentCommentRequest({
    externalId: data.external_id || '',
    content: data.content || ''
  });
  const resp = await unaryCall(commentClient, 'updateDocumentComment', req);
  const base = baseToObject(resp);
  const dataResp = {
    comment: resp.comment ? mapComment(resp.comment) : null
  };
  return handleResponse(base, dataResp);
};

export default {
  getKnowledgeBases,
  createKnowledgeBase,
  getRecentKnowledgeBases,
  getKnowledgeBaseDetail,
  getKnowledgeBaseDocuments,
  updateDocumentMeta,
  createDocument,
  createTemplate,
  getDocumentContent,
  getRecentDocuments,
  recordRecentDocument,
  getMyDocuments,
  listDocuments,
  listTemplates,
  getTemplate,
  listRecentTemplates,
  listUserGroups,
  getUserGroup,
  createUserGroup,
  updateUserGroup,
  deleteUserGroup,
  listUsers,
  updateUser,
  listUserGroupMembers,
  listInvitations,
  createInvitation,
  updateInvitation,
  deleteInvitation,
  listInvitationRecords,
  listOrgNodes,
  getOrgNode,
  createOrgNode,
  updateOrgNode,
  deleteOrgNode,
  moveOrgNode,
  listOrgNodeMembers,
  getSettings,
  updateSettings,
  createNotification,
  listNotifications,
  markNotificationsRead,
  markAllNotificationsRead,
  createDocumentComment,
  listDocumentComments,
  deleteDocumentComment,
  updateDocumentComment,
  CommentScope
};

export { CommentScope };
