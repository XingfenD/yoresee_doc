import {
  unaryCall,
  messages,
  invitationClient,
  baseToObject,
  handleResponse,
  mapInvitation,
  mapInvitationRecord
} from './shared';

const {
  ListInvitationsRequest,
  CreateInvitationRequest,
  UpdateInvitationRequest,
  DeleteInvitationRequest,
  ListInvitationRecordsRequest
} = messages;

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
  return handleResponse(base, {
    invitations: (resp.invitations || []).map(mapInvitation).filter(Boolean),
    total: resp.total ?? 0
  });
};

export const createInvitation = async (params = {}) => {
  const req = new CreateInvitationRequest({
    maxUsedCnt: typeof params.max_used_cnt === 'number' ? params.max_used_cnt : undefined,
    expiresAt: params.expires_at || undefined,
    note: params.note || undefined
  });

  const resp = await unaryCall(invitationClient, 'createInvitation', req);
  const base = baseToObject(resp);
  return handleResponse(base, {
    invitation: resp.invitation ? mapInvitation(resp.invitation) : null
  });
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
  return handleResponse(base, {
    records: (resp.records || []).map(mapInvitationRecord).filter(Boolean),
    total: resp.total ?? 0
  });
};
