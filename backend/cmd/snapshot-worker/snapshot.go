package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/domain_event"
	"github.com/XingfenD/yoresee_doc/internal/service/document_service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/key"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/bytedance/sonic"
	"github.com/sirupsen/logrus"
)

func (w *Worker) snapshotDoc(ctx context.Context, docID string, force bool) error {
	if _, loaded := w.inFlight.LoadOrStore(docID, struct{}{}); loaded {
		return nil
	}
	defer w.inFlight.Delete(docID)

	logrus.Infof("Snapshot start docId=%s force=%v", docID, force)

	docType, err := document_service.DocumentSvc.GetDocumentTypeByExternalID(ctx, docID)
	if err != nil {
		if errors.Is(err, status.StatusDocumentNotFound) {
			logrus.Infof("Snapshot doc not found docId=%s, cleaning up Redis", docID)
			cleanupCollabKeys(ctx, docID)
			return nil
		}
		logrus.Errorf("Snapshot get doc type failed docId=%s err=%v", docID, err)
		return err
	}

	state, content, err := w.fetchSnapshot(ctx, docID, string(docType))
	if err != nil {
		logrus.Errorf("Snapshot fetch failed docId=%s err=%v", docID, err)
		return err
	}
	if len(state) == 0 {
		logrus.Infof("Snapshot empty docId=%s", docID)
		return nil
	}
	logrus.Infof("Snapshot fetched docId=%s bytes=%d", docID, len(state))

	contentChanged, err := document_service.DocumentSvc.SaveDocumentSnapshotAndContent(ctx, docID, state, content)
	if err != nil {
		logrus.Errorf("Snapshot save failed docId=%s err=%v", docID, err)
		return err
	}
	if contentChanged {
		if pubErr := domain_event.PublishDocumentUpsertEvent(ctx, docID); pubErr != nil {
			logrus.Warnf("Snapshot publish search sync event failed docId=%s err=%v", docID, pubErr)
		}
	}

	if err := flushCollabKeys(ctx, docID, state); err != nil {
		logrus.Errorf("Snapshot redis flush failed docId=%s err=%v", docID, err)
		return err
	}

	if force {
		logrus.Infof("Snapshot saved (mq) for %s", docID)
	} else {
		logrus.Infof("Snapshot saved (scan) for %s", docID)
	}
	return nil
}

func (w *Worker) fetchSnapshot(ctx context.Context, docID, docType string) (state []byte, content string, err error) {
	url := fmt.Sprintf("%s/internal/yjs/doc-snapshot/%s?type=%s", strings.TrimRight(w.collabCoreHTTP, "/"), docID, docType)
	resp, err := w.client.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, "", nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	var payload struct {
		State   string `json:"state"`
		Content string `json:"content"`
	}
	if err := sonic.Unmarshal(body, &payload); err != nil {
		return nil, "", err
	}
	if payload.State == "" {
		return nil, payload.Content, nil
	}

	state, err = utils.DecodeBase64(payload.State)
	if err != nil {
		return nil, "", err
	}
	return state, payload.Content, nil
}

func flushCollabKeys(ctx context.Context, docID string, state []byte) error {
	rdb := storage.GetRedis()
	if rdb == nil {
		return nil
	}
	updatesKey := key.KeyCollabDocUpdates(docID)
	pipe := rdb.TxPipeline()
	pipe.Del(ctx, updatesKey)
	pipe.RPush(ctx, updatesKey, state)
	pipe.SRem(ctx, key.KeyCollabDirtyDocSet(), docID)
	_, err := pipe.Exec(ctx)
	return err
}

func cleanupCollabKeys(ctx context.Context, docID string) {
	rdb := storage.GetRedis()
	if rdb == nil {
		return
	}
	pipe := rdb.TxPipeline()
	pipe.SRem(ctx, key.KeyCollabDirtyDocSet(), docID)
	pipe.Del(ctx, key.KeyCollabDocUpdates(docID))
	pipe.Del(ctx, key.KeyCollabRoom(docID))
	if _, err := pipe.Exec(ctx); err != nil {
		logrus.Warnf("Snapshot cleanup Redis failed docId=%s err=%v", docID, err)
	}
}
