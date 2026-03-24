import { clients, messages, unaryCall } from './grpc_client';

const { authClient, systemClient } = clients;
const { AuthLoginRequest, AuthRegisterRequest, QuerySideBarDisplayRequest, SystemInfoRequest } = messages;

function baseToObject(resp) {
  const base = resp.base;
  return {
    code: base?.code ?? 50000,
    message: base?.message ?? 'unknown error'
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

function handleResponse(base, data) {
  if (base.code === 0) {
    return { ...base, ...data };
  }
  throw new Error('request failed');
}

// 登录
export const login = async (email, password) => {
  const req = new AuthLoginRequest({ email, password });

  const resp = await unaryCall(authClient, 'login', req);
  const base = baseToObject(resp);
  const data = {
    token: resp.token,
    user: mapUser(resp.user)
  };
  return handleResponse(base, data);
};

// 注册
export const register = async (username, password, email, invitationCode) => {
  const req = new AuthRegisterRequest({
    username,
    password,
    email,
    invitationCode: invitationCode || undefined
  });

  const resp = await unaryCall(authClient, 'register', req);
  const base = baseToObject(resp);
  return handleResponse(base, {});
};

// 获取系统信息
export const getSystemInfo = async () => {
  const req = new SystemInfoRequest();
  const resp = await unaryCall(systemClient, 'systemInfo', req);
  const base = baseToObject(resp);
  const data = {
    system_name: resp.systemName,
    system_register_mode: resp.systemRegisterMode
  };
  return handleResponse(base, data);
};

// 查询侧边栏可见项
export const querySideBarDisplay = async (scene) => {
  const req = new QuerySideBarDisplayRequest({ scene });
  const resp = await unaryCall(authClient, 'querySideBarDisplay', req);
  const base = baseToObject(resp);
  const data = {
    display_tabs: resp.displayTabs || []
  };
  return handleResponse(base, data);
};
