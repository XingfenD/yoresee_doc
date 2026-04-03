import { clients, messages, unaryCall } from '../grpc_client';
import { resolveFileUrl } from '@/utils/fileUrl';

const {
  documentClient,
  commentClient,
  knowledgeBaseClient,
  membershipClient,
  notificationClient,
  invitationClient,
  settingClient
} = clients;

export {
  unaryCall,
  messages,
  documentClient,
  commentClient,
  knowledgeBaseClient,
  membershipClient,
  notificationClient,
  invitationClient,
  settingClient
};

const { DocumentType } = messages;

export function baseToObject(resp) {
  const base = resp.base;
  return {
    code: base?.code ?? 50000,
    message: base?.message ?? 'unknown error'
  };
}

export function handleResponse(base, data) {
  if (base.code === 0) {
    return { ...base, ...data };
  }
  throw new Error('request failed');
}

export function toTimeRange(input) {
  if (!input) return null;
  const { TimeRange } = messages;
  return new TimeRange({
    start: input.start || '',
    end: input.end || ''
  });
}

export function mapDocument(doc) {
  if (!doc) return null;
  return {
    external_id: doc.externalId,
    title: doc.title,
    type: doc.type === DocumentType.MARKDOWN ? 'markdown' : '',
    summary: doc.summary,
    is_public: doc.isPublic,
    tags: doc.tags,
    view_count: doc.viewCount,
    edit_count: doc.editCount,
    created_at: doc.createdAt,
    updated_at: doc.updatedAt,
    has_children: doc.hasChildren,
    children: (doc.children || []).map(mapDocument)
  };
}

export function mapKnowledgeBase(kb) {
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

export function mapTemplate(tpl) {
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

export function mapUser(user) {
  if (!user) return null;
  return {
    external_id: user.externalId,
    username: user.username,
    email: user.email,
    nickname: user.nickname,
    avatar: resolveFileUrl(user.avatar),
    status: user.status,
    created_at: user.createdAt,
    updated_at: user.updatedAt,
    invitation_code: user.invitationCode ?? null
  };
}

export function mapUserGroup(group) {
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

export function mapOrgNode(node) {
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

export function mapInvitation(inv) {
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

export function mapInvitationRecord(rec) {
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

export function mapNotification(item) {
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

export function mapComment(item) {
  if (!item) return null;
  return {
    external_id: item.externalId,
    document_external_id: item.documentExternalId,
    parent_external_id: item.parentExternalId || null,
    content: item.content,
    anchor_id: item.anchorId || '',
    creator_user_external_id: item.creatorUserExternalId,
    creator_name: item.creatorName,
    creator_avatar: resolveFileUrl(item.creatorAvatar),
    created_at: item.createdAt
  };
}

export function mapAttachment(item) {
  if (!item) return null;
  return {
    external_id: item.externalId,
    document_external_id: item.documentExternalId,
    name: item.name,
    size: item.size,
    mime_type: item.mimeType,
    path: item.path,
    url: resolveFileUrl(item.url),
    created_at: item.createdAt,
    updated_at: item.updatedAt
  };
}

export function mapDocumentVersion(item) {
  if (!item) return null;
  return {
    version: item.version ?? 0,
    title: item.title || '',
    content: item.content || '',
    change_summary: item.changeSummary || '',
    created_at: item.createdAt || ''
  };
}
