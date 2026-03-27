import {
  unaryCall,
  messages,
  documentClient,
  baseToObject,
  handleResponse,
  mapTemplate
} from './shared';

const {
  CreateTemplateRequest,
  ListTemplatesRequest,
  GetTemplateRequest,
  ListRecentTemplatesRequest,
  CreateTemplateContainer
} = messages;

const containerMap = {
  own: CreateTemplateContainer.OWN_TEMPLATE,
  knowledge_base: CreateTemplateContainer.KNOWLEDGEBASE_TEMPLATE,
  public: CreateTemplateContainer.PUBLIC_TEMPLATE
};

export const createTemplate = async (data = {}) => {
  const req = new CreateTemplateRequest({
    targetContainer: containerMap[data.target_container] ?? CreateTemplateContainer.OWN_TEMPLATE,
    knowledgeBaseId: data.knowledge_base_id || '',
    templateContent: data.template_content || ''
  });

  const resp = await unaryCall(documentClient, 'createTemplate', req);
  const base = baseToObject(resp);
  return handleResponse(base, {});
};

export const listTemplates = async (params = {}) => {
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
  return handleResponse(base, {
    templates: (resp.templates || []).map(mapTemplate),
    total: resp.total ?? 0
  });
};

export const getTemplate = async (templateId, options = {}) => {
  const req = new GetTemplateRequest({
    templateId: Number(templateId) || 0,
    recordRecentLog: Boolean(options.record_recent_log)
  });

  const resp = await unaryCall(documentClient, 'getTemplate', req);
  const base = baseToObject(resp);
  return handleResponse(base, {
    template: resp.template ? mapTemplate(resp.template) : null
  });
};

export const listRecentTemplates = async (params = {}) => {
  const req = new ListRecentTemplatesRequest({
    startTime: params.start_time || undefined,
    endTime: params.end_time || undefined,
    page: params.page || undefined,
    pageSize: params.page_size || undefined
  });

  const resp = await unaryCall(documentClient, 'listRecentTemplates', req);
  const base = baseToObject(resp);
  return handleResponse(base, {
    templates: (resp.templates || []).map(mapTemplate),
    total: resp.total ?? 0
  });
};
