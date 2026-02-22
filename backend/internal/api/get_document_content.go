package api

import (
	"context"
	"net/http"

	"github.com/XingfenD/yoresee_doc/internal/service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/gin-gonic/gin"
)

func (h *GetDocumentContentHandler) handle(ctx context.Context, req Request) (resp Response, err error) {
	getDocReq := req.(*GetDocumentContentRequest)

	// 从上下文中获取用户信息和命名域
	// 这里简化处理，实际应该从JWT中获取用户信息
	// 暂时使用默认值
	userExternalID := "admin"
	namespace := "default"

	// 调用服务获取文档内容
	document, content, err := service.DocumentSvc.GetDocumentWithContent(getDocReq.DocumentExternalID, userExternalID, namespace)
	if err != nil {
		return nil, err
	}

	// 构造响应
	return &GetDocumentContentResponse{
		BaseResponse: GenBaseRespWithErr(status.StatusSuccess),
		Document:     service.DocumentSvc.ConvertToDocumentResponse(document),
		Content:      content,
	}, nil
}

func (h *GetDocumentContentHandler) GinHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// 从路径参数中获取documentExternalID
		documentExternalID := c.Param("documentExternalID")

		req := &GetDocumentContentRequest{
			DocumentExternalID: documentExternalID,
		}

		resp, err := h.handle(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, GenBaseRespWithErrAndCtx(c, err))
			return
		}

		if getDocResp, ok := resp.(*GetDocumentContentResponse); ok {
			getDocResp.Message = ""
		}
		c.JSON(http.StatusOK, resp)
	}
}
