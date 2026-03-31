package document_service

import (
	"bytes"
	"context"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
)

const maxAttachmentSize = 5 * 1024 * 1024

func (s *DocumentService) UploadAttachment(
	ctx context.Context,
	userExternalID string,
	documentExternalID string,
	file []byte,
	fileName string,
	contentType string,
) (*dto.AttachmentResponse, error) {
	if s == nil || strings.TrimSpace(userExternalID) == "" || strings.TrimSpace(documentExternalID) == "" {
		return nil, status.GenErrWithCustomMsg(status.StatusParamError, "invalid upload attachment request")
	}
	if len(file) == 0 || len(file) > maxAttachmentSize {
		return nil, status.GenErrWithCustomMsg(status.StatusParamError, "invalid attachment file")
	}
	if config.GlobalConfig == nil || strings.TrimSpace(config.GlobalConfig.Minio.Bucket) == "" {
		return nil, status.GenErrWithCustomMsg(status.StatusServiceInternalError, "attachment storage is not configured")
	}

	documentID, err := s.documentRepo.GetIDByExternalID(documentExternalID).Exec(ctx)
	if err != nil {
		logrus.Errorf("[Service layer: DocumentService] query document id failed, document_external_id=%s, err=%+v", documentExternalID, err)
		return nil, status.GenErrWithCustomMsg(status.StatusReadDBError, "query document failed")
	}
	if documentID <= 0 {
		return nil, status.GenErrWithCustomMsg(status.StatusDocumentNotFound, "document not found")
	}

	userID, err := s.userRepo.GetIDByExternalID(userExternalID).Exec()
	if err != nil || userID <= 0 {
		logrus.Errorf("[Service layer: DocumentService] query user id failed, user_external_id=%s, err=%+v", userExternalID, err)
		return nil, status.GenErrWithCustomMsg(status.StatusUserNotFound, "user not found")
	}

	safeName := strings.TrimSpace(fileName)
	if safeName == "" {
		safeName = "attachment.bin"
	}
	contentType = strings.TrimSpace(contentType)
	if contentType == "" {
		contentType = http.DetectContentType(file)
	}

	attachmentExternalID := utils.GenerateExternalID(utils.ExternalIDContextAttachment)
	objectName := fmt.Sprintf(
		"attachments/%s/%s-%d%s",
		documentExternalID,
		attachmentExternalID,
		time.Now().UnixNano(),
		resolveAttachmentExt(safeName, contentType),
	)

	url, err := storage.UploadFile(
		config.GlobalConfig.Minio.Bucket,
		objectName,
		bytes.NewReader(file),
		int64(len(file)),
		contentType,
	)
	if err != nil {
		logrus.Errorf("[Service layer: DocumentService] upload attachment to minio failed, document_external_id=%s, object_name=%s, err=%+v", documentExternalID, objectName, err)
		return nil, status.GenErrWithCustomMsg(status.StatusServiceInternalError, "upload attachment failed")
	}

	attachment := &model.Attachment{
		ExternalID: attachmentExternalID,
		DocumentID: documentID,
		Name:       safeName,
		Path:       objectName,
		Size:       int64(len(file)),
		MimeType:   contentType,
		UserID:     userID,
	}
	if err = s.attachmentRepo.Create(attachment); err != nil {
		logrus.Errorf("[Service layer: DocumentService] save attachment meta failed, document_external_id=%s, attachment_external_id=%s, err=%+v", documentExternalID, attachmentExternalID, err)
		return nil, status.GenErrWithCustomMsg(status.StatusWriteDBError, "save attachment failed")
	}

	return &dto.AttachmentResponse{
		ExternalID:         attachment.ExternalID,
		DocumentExternalID: documentExternalID,
		Name:               attachment.Name,
		Size:               attachment.Size,
		MimeType:           attachment.MimeType,
		Path:               attachment.Path,
		URL:                url,
		CreatedAt:          attachment.CreatedAt,
		UpdatedAt:          attachment.UpdatedAt,
	}, nil
}

func (s *DocumentService) ListAttachments(ctx context.Context, documentExternalID string) ([]*dto.AttachmentResponse, error) {
	if s == nil || strings.TrimSpace(documentExternalID) == "" {
		return nil, status.GenErrWithCustomMsg(status.StatusParamError, "invalid list attachments request")
	}
	if config.GlobalConfig == nil || strings.TrimSpace(config.GlobalConfig.Minio.Bucket) == "" {
		return nil, status.GenErrWithCustomMsg(status.StatusServiceInternalError, "attachment storage is not configured")
	}

	documentID, err := s.documentRepo.GetIDByExternalID(documentExternalID).Exec(ctx)
	if err != nil {
		logrus.Errorf("[Service layer: DocumentService] query document id failed, document_external_id=%s, err=%+v", documentExternalID, err)
		return nil, status.GenErrWithCustomMsg(status.StatusReadDBError, "query document failed")
	}
	if documentID <= 0 {
		return nil, status.GenErrWithCustomMsg(status.StatusDocumentNotFound, "document not found")
	}

	attachments, err := s.attachmentRepo.ListByDocumentID(documentID)
	if err != nil {
		logrus.Errorf("[Service layer: DocumentService] list attachments failed, document_external_id=%s, document_id=%d, err=%+v", documentExternalID, documentID, err)
		return nil, status.GenErrWithCustomMsg(status.StatusReadDBError, "list attachments failed")
	}

	resp := make([]*dto.AttachmentResponse, 0, len(attachments))
	for _, attachment := range attachments {
		url, presignErr := storage.GetPresignedURL(config.GlobalConfig.Minio.Bucket, attachment.Path, 7*24*time.Hour)
		if presignErr != nil {
			logrus.Warnf("[Service layer: DocumentService] generate attachment presigned url failed, attachment_external_id=%s, err=%+v", attachment.ExternalID, presignErr)
		}
		resp = append(resp, &dto.AttachmentResponse{
			ExternalID:         attachment.ExternalID,
			DocumentExternalID: documentExternalID,
			Name:               attachment.Name,
			Size:               attachment.Size,
			MimeType:           attachment.MimeType,
			Path:               attachment.Path,
			URL:                url,
			CreatedAt:          attachment.CreatedAt,
			UpdatedAt:          attachment.UpdatedAt,
		})
	}
	return resp, nil
}

func (s *DocumentService) DeleteAttachment(ctx context.Context, documentExternalID, attachmentExternalID string) error {
	if s == nil || strings.TrimSpace(documentExternalID) == "" || strings.TrimSpace(attachmentExternalID) == "" {
		return status.GenErrWithCustomMsg(status.StatusParamError, "invalid delete attachment request")
	}
	if config.GlobalConfig == nil || strings.TrimSpace(config.GlobalConfig.Minio.Bucket) == "" {
		return status.GenErrWithCustomMsg(status.StatusServiceInternalError, "attachment storage is not configured")
	}

	documentID, err := s.documentRepo.GetIDByExternalID(documentExternalID).Exec(ctx)
	if err != nil {
		logrus.Errorf("[Service layer: DocumentService] query document id failed, document_external_id=%s, err=%+v", documentExternalID, err)
		return status.GenErrWithCustomMsg(status.StatusReadDBError, "query document failed")
	}
	if documentID <= 0 {
		return status.GenErrWithCustomMsg(status.StatusDocumentNotFound, "document not found")
	}

	attachment, err := s.attachmentRepo.GetByExternalIDAndDocumentID(attachmentExternalID, documentID)
	if err != nil {
		logrus.Errorf("[Service layer: DocumentService] query attachment failed, document_external_id=%s, attachment_external_id=%s, err=%+v", documentExternalID, attachmentExternalID, err)
		return status.GenErrWithCustomMsg(status.StatusDocumentNotFound, "attachment not found")
	}

	if err = storage.DeleteFile(config.GlobalConfig.Minio.Bucket, attachment.Path); err != nil {
		logrus.Errorf("[Service layer: DocumentService] delete attachment object failed, attachment_external_id=%s, path=%s, err=%+v", attachmentExternalID, attachment.Path, err)
		return status.GenErrWithCustomMsg(status.StatusServiceInternalError, "delete attachment failed")
	}
	if err = s.attachmentRepo.DeleteByID(attachment.ID); err != nil {
		logrus.Errorf("[Service layer: DocumentService] delete attachment meta failed, attachment_external_id=%s, err=%+v", attachmentExternalID, err)
		return status.GenErrWithCustomMsg(status.StatusWriteDBError, "delete attachment failed")
	}
	return nil
}

func resolveAttachmentExt(filename, contentType string) string {
	ext := strings.TrimSpace(strings.ToLower(filepath.Ext(filename)))
	if ext != "" {
		return ext
	}
	exts, err := mime.ExtensionsByType(contentType)
	if err == nil && len(exts) > 0 {
		return exts[0]
	}
	return ".bin"
}
