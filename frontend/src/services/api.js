import { ElMessage } from 'element-plus';
import { clients, messages, unaryCall } from './grpc_client';

const { documentClient, knowledgeBaseClient } = clients;
const {
  ListKnowledgeBasesRequest,
  GetKnowledgeBaseRequest,
  CreateDocumentRequest,
  GetDocumentContentRequest,
  GetOwnDocumentsRequest,
  ListDocumentsRequest,
  RecursiveOptions,
  TimeRange,
  CreateDocumentContainerType,
  DocumentType
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

function handleResponse(base, data) {
  if (base.code === 20000) {
    return { ...base, ...data };
  }
  const errorMessage = base.message || '请求失败';
  ElMessage.error(errorMessage);
  throw new Error(errorMessage);
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

// 获取知识库详情
export const getKnowledgeBaseDetail = async (knowledgeBaseExternalID, params = {}) => {
  const req = new GetKnowledgeBaseRequest({
    knowledgeBaseExternalId: knowledgeBaseExternalID,
    recordRecentLog: Boolean(params.record_recent_log),
    page: params.page || undefined,
    pageSize: params.page_size || undefined
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

// 获取知识库文档列表
export const getKnowledgeBaseDocuments = async (knowledgeBaseExternalID, params = {}) => {
  return getKnowledgeBaseDetail(knowledgeBaseExternalID, {
    record_recent_log: false,
    page: 1,
    page_size: 1000,
    ...params
  });
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

  const resp = await unaryCall(documentClient, 'createDocument', req);
  const base = baseToObject(resp);
  const dataResp = {
    external_id: resp.externalId
  };
  return handleResponse(base, dataResp);
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

// 获取我的文档列表
export const getMyDocuments = async (params = {}) => {
  const req = new GetOwnDocumentsRequest({
    page: params.page || undefined,
    pageSize: params.page_size || undefined
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

export default {
  getKnowledgeBases,
  getKnowledgeBaseDetail,
  getKnowledgeBaseDocuments,
  createDocument,
  getDocumentContent,
  getMyDocuments,
  listDocuments
};
