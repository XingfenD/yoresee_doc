import {
  unaryCall,
  messages,
  membershipClient,
  baseToObject,
  handleResponse,
  mapUserGroup,
  mapUser,
  mapOrgNode
} from './shared';

const {
  CreateUserGroupRequest,
  DeleteUserGroupRequest,
  GetUserGroupRequest,
  UpdateUserGroupRequest,
  ListUsersRequest,
  UpdateUserRequest,
  ListUserGroupMembersRequest,
  ListUserGroupsRequest,
  ListOrgNodesRequest,
  GetOrgNodeRequest,
  CreateOrgNodeRequest,
  UpdateOrgNodeRequest,
  DeleteOrgNodeRequest,
  MoveOrgNodeRequest,
  ListOrgNodeMembersRequest
} = messages;

export const listUserGroups = async (params = {}) => {
  const req = new ListUserGroupsRequest({
    keyword: params.keyword || undefined,
    page: params.page || undefined,
    pageSize: params.page_size || undefined
  });

  const resp = await unaryCall(membershipClient, 'listUserGroups', req);
  const base = baseToObject(resp);
  return handleResponse(base, {
    user_groups: (resp.userGroups || []).map(mapUserGroup),
    total: resp.total ?? 0
  });
};

export const getUserGroup = async (externalId) => {
  const req = new GetUserGroupRequest({
    externalId: externalId || ''
  });

  const resp = await unaryCall(membershipClient, 'getUserGroup', req);
  const base = baseToObject(resp);
  return handleResponse(base, {
    user_group: resp.userGroup ? mapUserGroup(resp.userGroup) : null
  });
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
  return handleResponse(base, {
    external_id: resp.externalId
  });
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
  return handleResponse(base, {
    users: (resp.users || []).map(mapUser),
    total: resp.total ?? 0
  });
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
  return handleResponse(base, {
    users: (resp.users || []).map(mapUser),
    total: resp.total ?? 0
  });
};

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
  return handleResponse(base, {
    org_nodes: (resp.orgNodes || []).map(mapOrgNode),
    total: resp.total ?? 0
  });
};

export const getOrgNode = async (externalId, params = {}) => {
  const req = new GetOrgNodeRequest({
    externalId: externalId || '',
    includeChildren: Boolean(params.include_children)
  });

  const resp = await unaryCall(membershipClient, 'getOrgNode', req);
  const base = baseToObject(resp);
  return handleResponse(base, {
    org_node: resp.orgNode ? mapOrgNode(resp.orgNode) : null
  });
};

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
  return handleResponse(base, {
    external_id: resp.externalId
  });
};

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

export const deleteOrgNode = async (externalId) => {
  const req = new DeleteOrgNodeRequest({
    externalId: externalId || ''
  });

  const resp = await unaryCall(membershipClient, 'deleteOrgNode', req);
  const base = baseToObject(resp);
  return handleResponse(base, {});
};

export const moveOrgNode = async (data = {}) => {
  const req = new MoveOrgNodeRequest({
    externalId: data.external_id || '',
    newParentExternalId: data.new_parent_external_id || ''
  });

  const resp = await unaryCall(membershipClient, 'moveOrgNode', req);
  const base = baseToObject(resp);
  return handleResponse(base, {});
};

export const listOrgNodeMembers = async (params = {}) => {
  const req = new ListOrgNodeMembersRequest({
    externalId: params.external_id || '',
    keyword: params.keyword || undefined,
    page: params.page || undefined,
    pageSize: params.page_size || undefined
  });

  const resp = await unaryCall(membershipClient, 'listOrgNodeMembers', req);
  const base = baseToObject(resp);
  return handleResponse(base, {
    users: (resp.users || []).map(mapUser),
    total: resp.total ?? 0
  });
};
