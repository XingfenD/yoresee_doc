import {
  unaryCall,
  messages,
  documentClient,
  baseToObject,
  handleResponse,
  mapDocument,
  mapAttachment,
  toTimeRange
} from './shared';

const {
  UpdateDocumentMetaRequest,
  UpdateDocumentSettingsRequest,
  CreateDocumentRequest,
  ListRecentDocumentsRequest,
  RecordRecentDocumentRequest,
  GetDocumentContentRequest,
  GetDocumentSettingsRequest,
  GetOwnDocumentsRequest,
  ListDocumentsRequest,
  UploadDocumentAttachmentRequest,
  ListDocumentAttachmentsRequest,
  DeleteDocumentAttachmentRequest,
  RecursiveOptions,
  CreateDocumentContainerType,
  DocumentType
} = messages;

export const updateDocumentMeta = async (documentExternalID, data = {}) => {
  const req = new UpdateDocumentMetaRequest({
    externalId: documentExternalID,
    title: data.title ?? undefined,
    summary: data.summary ?? undefined,
    tags: Array.isArray(data.tags) ? data.tags : undefined
  });

  const resp = await unaryCall(documentClient, 'updateDocumentMeta', req);
  const base = baseToObject(resp);
  return handleResponse(base, {});
};

export const updateDocumentSettings = async (documentExternalID, data = {}) => {
  const req = new UpdateDocumentSettingsRequest({
    externalId: documentExternalID,
    isPublic: typeof data.is_public === 'boolean' ? data.is_public : undefined
  });

  const resp = await unaryCall(documentClient, 'updateDocumentSettings', req);
  const base = baseToObject(resp);
  return handleResponse(base, {
    is_public: typeof resp.isPublic === 'boolean' ? resp.isPublic : undefined
  });
};

export const createDocument = async (data) => {
  const req = new CreateDocumentRequest({
    title: data.title || '',
    type: DocumentType.MARKDOWN,
    isPublic: typeof data.is_public === 'boolean' ? data.is_public : undefined
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
  return handleResponse(base, {
    external_id: resp.externalId
  });
};

export const getDocumentContent = async (documentExternalID) => {
  const req = new GetDocumentContentRequest({
    documentExternalId: documentExternalID
  });

  const resp = await unaryCall(documentClient, 'getDocumentContent', req);
  const base = baseToObject(resp);
  return handleResponse(base, {
    document: mapDocument(resp.document),
    content: resp.content
  });
};

export const getDocumentSettings = async (documentExternalID) => {
  const req = new GetDocumentSettingsRequest({
    documentExternalId: documentExternalID
  });

  const resp = await unaryCall(documentClient, 'getDocumentSettings', req);
  const base = baseToObject(resp);
  return handleResponse(base, {
    is_public: typeof resp.isPublic === 'boolean' ? resp.isPublic : undefined
  });
};

export const getRecentDocuments = async (params = {}) => {
  const req = new ListRecentDocumentsRequest({
    startTime: params.start_time || undefined,
    endTime: params.end_time || undefined,
    page: params.page || undefined,
    pageSize: params.page_size || undefined
  });

  const resp = await unaryCall(documentClient, 'listRecentDocuments', req);
  const base = baseToObject(resp);
  return handleResponse(base, {
    documents: (resp.documents || []).map(mapDocument),
    total: resp.total ?? 0
  });
};

export const recordRecentDocument = async (documentExternalID) => {
  const req = new RecordRecentDocumentRequest({
    documentExternalId: documentExternalID
  });

  const resp = await unaryCall(documentClient, 'recordRecentDocument', req);
  const base = baseToObject(resp);
  return handleResponse(base, {});
};

export const getMyDocuments = async (params = {}) => {
  const req = new GetOwnDocumentsRequest({
    page: params.page || undefined,
    pageSize: params.page_size || undefined,
    directoryOnly: params.directory_only || false
  });

  const resp = await unaryCall(documentClient, 'getOwnDocuments', req);
  const base = baseToObject(resp);
  return handleResponse(base, {
    documents: (resp.documents || []).map(mapDocument),
    total_count: resp.totalCount ?? 0
  });
};

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
    knowledgeBaseExternalId: params.knowledge_base_external_id || undefined,
    titleKeyword: params.title_keyword || undefined,
    type: params.type || undefined,
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
  return handleResponse(base, {
    documents: (resp.documents || []).map(mapDocument),
    total_count: resp.totalCount ?? 0
  });
};

export const uploadDocumentAttachment = async (params = {}) => {
  const fileBytes = params.file_bytes instanceof Uint8Array ? params.file_bytes : new Uint8Array([]);
  const req = new UploadDocumentAttachmentRequest({
    documentExternalId: params.document_external_id || '',
    fileContent: fileBytes,
    fileName: params.file_name || '',
    contentType: params.content_type || undefined
  });

  const resp = await unaryCall(documentClient, 'uploadDocumentAttachment', req);
  const base = baseToObject(resp);
  return handleResponse(base, {
    attachment: mapAttachment(resp.attachment)
  });
};

export const listDocumentAttachments = async (documentExternalID) => {
  const req = new ListDocumentAttachmentsRequest({
    documentExternalId: documentExternalID || ''
  });

  const resp = await unaryCall(documentClient, 'listDocumentAttachments', req);
  const base = baseToObject(resp);
  return handleResponse(base, {
    attachments: (resp.attachments || []).map(mapAttachment)
  });
};

export const deleteDocumentAttachment = async (documentExternalID, attachmentExternalID) => {
  const req = new DeleteDocumentAttachmentRequest({
    documentExternalId: documentExternalID || '',
    attachmentExternalId: attachmentExternalID || ''
  });

  const resp = await unaryCall(documentClient, 'deleteDocumentAttachment', req);
  const base = baseToObject(resp);
  return handleResponse(base, {});
};
