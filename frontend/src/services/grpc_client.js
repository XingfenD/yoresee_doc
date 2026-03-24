import { createPromiseClient, Code } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import {
  AuthService,
  DocumentService,
  KnowledgeBaseService,
  MembershipService,
  SystemService
} from '@/gen/yoresee_doc/v1/yoresee_doc_connect.js';
import * as messages from '@/gen/yoresee_doc/v1/yoresee_doc_pb.js';

const CONNECT_ENDPOINT = import.meta.env.VITE_GRPC_WEB_ENDPOINT || '/grpc';

const transport = createConnectTransport({
  baseUrl: CONNECT_ENDPOINT,
  useBinaryFormat: false
});

const authClient = createPromiseClient(AuthService, transport);
const documentClient = createPromiseClient(DocumentService, transport);
const knowledgeBaseClient = createPromiseClient(KnowledgeBaseService, transport);
const membershipClient = createPromiseClient(MembershipService, transport);
const systemClient = createPromiseClient(SystemService, transport);

export function buildHeaders() {
  const token = localStorage.getItem('token');
  const language = localStorage.getItem('language') || 'en';
  const headers = {
    'accept-language': language
  };
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }
  return headers;
}

export async function unaryCall(client, method, request) {
  try {
    return await client[method](request, { headers: buildHeaders() });
  } catch (err) {
    const code = err?.code;
    if (code === Code.Unauthenticated || code === 'unauthenticated') {
      localStorage.removeItem('token');
      localStorage.removeItem('userInfo');
      window.location.href = '/login';
    }
    throw err;
  }
}

export const clients = {
  authClient,
  documentClient,
  knowledgeBaseClient,
  membershipClient,
  systemClient
};

export { messages };
