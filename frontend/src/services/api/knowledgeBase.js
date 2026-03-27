import {
  unaryCall,
  messages,
  knowledgeBaseClient,
  baseToObject,
  handleResponse,
  toTimeRange,
  mapKnowledgeBase,
  mapDocument
} from './shared';

const {
  ListKnowledgeBasesRequest,
  ListRecentKnowledgeBasesRequest,
  GetKnowledgeBaseRequest,
  CreateKnowledgeBaseRequest
} = messages;

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

export const createKnowledgeBase = async (data = {}) => {
  const req = new CreateKnowledgeBaseRequest({
    name: data.name || '',
    description: data.description || '',
    cover: data.cover || '',
    isPublic: Boolean(data.is_public)
  });

  const resp = await unaryCall(knowledgeBaseClient, 'createKnowledgeBase', req);
  const base = baseToObject(resp);
  return handleResponse(base, {
    external_id: resp.externalId
  });
};

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

export const getKnowledgeBaseDocuments = async (knowledgeBaseExternalID, params = {}) => {
  return getKnowledgeBaseDetail(knowledgeBaseExternalID, {
    record_recent_log: false,
    page: 1,
    page_size: 1000,
    directory_only: params.directory_only || false,
    ...params
  });
};
