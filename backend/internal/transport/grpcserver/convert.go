package grpcserver

import (
	"time"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/status"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
)

func baseResponseFromErr(err error) *pb.BaseResponse {
	if err == nil {
		return baseResponseFromStatus(status.StatusSuccess)
	}
	return baseResponseFromStatus(err)
}

func baseResponceFromStatusWithCustomMsg(err error, msg string) *pb.BaseResponse {
	return baseResponseFromStatus(status.GenErrWithCustomMsg(err, msg))
}

func baseResponseFromStatus(err error) *pb.BaseResponse {
	st, ok := err.(*status.Status)
	if !ok {
		st = status.StatusServiceInternalError.(*status.Status)
	}
	return &pb.BaseResponse{
		Code:    int32(st.Code),
		Message: st.Message,
	}
}

func timeToString(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

func toDocumentType(t dto.DocumentType) pb.DocumentType {
	switch t {
	case dto.DocumentType_Markdown:
		return pb.DocumentType_DOCUMENT_TYPE_MARKDOWN
	default:
		return pb.DocumentType_DOCUMENT_TYPE_UNSPECIFIED
	}
}

func fromDocumentType(t pb.DocumentType) dto.DocumentType {
	switch t {
	case pb.DocumentType_DOCUMENT_TYPE_MARKDOWN:
		return dto.DocumentType_Markdown
	default:
		return dto.DocumentType_Markdown
	}
}

func fromCreateTemplateContainer(t pb.CreateTemplateContainer) dto.TemplateContainer {
	switch t {
	case pb.CreateTemplateContainer_KNOWLEDGEBASE_TEMPLATE:
		return dto.TemplateContainerKnowledgeBase
	case pb.CreateTemplateContainer_PUBLIC_TEMPLATE:
		return dto.TemplateContainerPublic
	default:
		return dto.TemplateContainerOwn
	}
}

func toDocumentResponse(doc *dto.DocumentMetaResponse) *pb.DocumentResponse {
	if doc == nil {
		return nil
	}
	resp := &pb.DocumentResponse{
		ExternalId:  doc.ExternalID,
		Title:       doc.Title,
		Type:        toDocumentType(doc.Type),
		Summary:     doc.Summary,
		Status:      int32(doc.Status),
		Tags:        doc.Tags,
		ViewCount:   int32(doc.ViewCount),
		EditCount:   int32(doc.EditCount),
		CreatedAt:   timeToString(doc.CreatedAt),
		UpdatedAt:   timeToString(doc.UpdatedAt),
		HasChildren: doc.HasChildren,
	}
	if len(doc.Children) > 0 {
		resp.Children = make([]*pb.DocumentResponse, 0, len(doc.Children))
		for _, child := range doc.Children {
			resp.Children = append(resp.Children, toDocumentResponse(child))
		}
	}
	return resp
}

func toKnowledgeBaseResponse(kb *dto.KnowledgeBaseResponse) *pb.KnowledgeBaseResponse {
	if kb == nil {
		return nil
	}
	return &pb.KnowledgeBaseResponse{
		ExternalId:            kb.ExternalID,
		Name:                  kb.Name,
		Description:           kb.Description,
		Cover:                 kb.Cover,
		IsPublic:              kb.IsPublic,
		CreatedAt:             timeToString(kb.CreatedAt),
		UpdatedAt:             timeToString(kb.UpdatedAt),
		DeletedAt:             timeToString(kb.DeletedAt),
		CreatorUserExternalId: kb.CreatorUserExternalID,
		CreatorName:           kb.CreatorName,
		DocumentsCount:        kb.DocumentsCount,
	}
}

func toUserResponse(user *dto.UserResponse) *pb.UserResponse {
	if user == nil {
		return nil
	}
	resp := &pb.UserResponse{
		ExternalId: user.ExternalID,
		Username:   user.Username,
		Email:      user.Email,
		Nickname:   user.Nickname,
		Avatar:     user.Avatar,
		Status:     int32(user.Status),
		CreatedAt:  timeToString(user.CreatedAt),
		UpdatedAt:  timeToString(user.UpdatedAt),
	}
	if user.InvitationCode != nil {
		resp.InvitationCode = user.InvitationCode
	}
	return resp
}

func toTemplateResponse(tpl *dto.TemplateResponse) *pb.TemplateResponse {
	if tpl == nil {
		return nil
	}
	return &pb.TemplateResponse{
		Id:                      tpl.ID,
		Name:                    tpl.Name,
		Description:             tpl.Description,
		Content:                 tpl.Content,
		Scope:                   tpl.Scope,
		KnowledgeBaseExternalId: tpl.KnowledgeBaseExternalID,
		Tags:                    tpl.Tags,
		CreatedAt:               timeToString(tpl.CreatedAt),
		UpdatedAt:               timeToString(tpl.UpdatedAt),
	}
}

func toUserGroupResponse(group *dto.UserGroupResponse) *pb.UserGroupResponse {
	if group == nil {
		return nil
	}
	resp := &pb.UserGroupResponse{
		ExternalId:            group.ExternalID,
		Name:                  group.Name,
		Description:           group.Description,
		CreatorUserExternalId: group.CreatorUserExternalID,
		MemberCount:           int32(group.MemberCount),
	}
	return resp
}
