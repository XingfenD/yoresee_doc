import {
  unaryCall,
  messages,
  commentClient,
  baseToObject,
  handleResponse,
  mapComment
} from './shared';

const {
  CreateDocumentCommentRequest,
  ListDocumentCommentsRequest,
  DeleteDocumentCommentRequest,
  UpdateDocumentCommentRequest
} = messages;

export const createDocumentComment = async (data = {}) => {
  const req = new CreateDocumentCommentRequest({
    documentExternalId: data.document_external_id || '',
    content: data.content || '',
    parentExternalId: data.parent_external_id || undefined,
    anchorId: data.anchor_id || undefined
  });
  const resp = await unaryCall(commentClient, 'createDocumentComment', req);
  const base = baseToObject(resp);
  return handleResponse(base, {
    comment: resp.comment ? mapComment(resp.comment) : null
  });
};

export const listDocumentComments = async (params = {}) => {
  const req = new ListDocumentCommentsRequest({
    documentExternalId: params.document_external_id || '',
    page: params.page || 1,
    pageSize: params.page_size || 10
  });
  const resp = await unaryCall(commentClient, 'listDocumentComments', req);
  const base = baseToObject(resp);
  return handleResponse(base, {
    comments: (resp.comments || []).map(mapComment),
    total: resp.total ?? 0
  });
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
  return handleResponse(base, {
    comment: resp.comment ? mapComment(resp.comment) : null
  });
};
