package grpcserver

import (
	"time"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/mapper/doc_type_mapper"
	"github.com/XingfenD/yoresee_doc/internal/media"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
)

func baseResponseFromErr(err error) *pb.BaseResponse {
	if err == nil {
		return baseResponseFromStatus(status.StatusSuccess)
	}
	return baseResponseFromStatus(err)
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

func toDocumentResponse(doc *dto.DocumentMetaResponse) *pb.DocumentResponse {
	if doc == nil {
		return nil
	}
	resp := &pb.DocumentResponse{
		ExternalId:  doc.ExternalID,
		Title:       doc.Title,
		Type:        doc_type_mapper.ToProtoType(doc.Type),
		Summary:     doc.Summary,
		Tags:        doc.Tags,
		ViewCount:   int32(doc.ViewCount),
		EditCount:   int32(doc.EditCount),
		CreatedAt:   timeToString(doc.CreatedAt),
		UpdatedAt:   timeToString(doc.UpdatedAt),
		HasChildren: doc.HasChildren,
		IsPublic:    doc.IsPublic,
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
		Avatar:     media.BuildAvatarURL(user.ExternalID, user.AvatarObjectKey, user.AvatarVersion),
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
		Type:                    doc_type_mapper.ToProtoType(tpl.Type),
		Scope:                   tpl.Scope,
		KnowledgeBaseExternalId: tpl.KnowledgeBaseExternalID,
		Tags:                    tpl.Tags,
		CreatedAt:               timeToString(tpl.CreatedAt),
		UpdatedAt:               timeToString(tpl.UpdatedAt),
	}
}

func toAttachmentResponse(attachment *dto.AttachmentResponse) *pb.AttachmentResponse {
	if attachment == nil {
		return nil
	}
	return &pb.AttachmentResponse{
		ExternalId:         attachment.ExternalID,
		DocumentExternalId: attachment.DocumentExternalID,
		Name:               attachment.Name,
		Size:               attachment.Size,
		MimeType:           attachment.MimeType,
		Path:               attachment.Path,
		Url:                attachment.URL,
		CreatedAt:          timeToString(attachment.CreatedAt),
		UpdatedAt:          timeToString(attachment.UpdatedAt),
	}
}

func toDocumentVersionResponse(item *dto.DocumentVersionResponse) *pb.DocumentVersionResponse {
	if item == nil {
		return nil
	}
	return &pb.DocumentVersionResponse{
		Version:       int32(item.Version),
		Title:         item.Title,
		Content:       item.Content,
		ChangeSummary: item.ChangeSummary,
		CreatedAt:     timeToString(item.CreatedAt),
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

func toInvitationResponse(inv *dto.InvitationResponse) *pb.InvitationResponse {
	if inv == nil {
		return nil
	}
	resp := &pb.InvitationResponse{
		Code:                inv.Code,
		CreatedByExternalId: inv.CreatedByExternalID,
		CreatedByName:       inv.CreatedByName,
		UsedCnt:             inv.UsedCnt,
		Disabled:            inv.Disabled,
		CreatedAt:           timeToString(inv.CreatedAt),
	}
	if inv.MaxUsedCnt != nil {
		resp.MaxUsedCnt = inv.MaxUsedCnt
	}
	if inv.ExpiresAt != nil {
		resp.ExpiresAt = utils.Of(timeToString(*inv.ExpiresAt))
	}
	if inv.Note != nil {
		resp.Note = inv.Note
	}
	return resp
}

func toInvitationRecordResponse(rec *dto.InvitationRecordResponse) *pb.InvitationRecordResponse {
	if rec == nil {
		return nil
	}
	resp := &pb.InvitationRecordResponse{
		Code:   rec.Code,
		UsedBy: rec.UsedBy,
		UsedAt: timeToString(rec.UsedAt),
		Status: rec.Status,
	}
	if rec.UsedByExternalID != "" {
		resp.UsedByExternalId = utils.Of(rec.UsedByExternalID)
	}
	return resp
}
