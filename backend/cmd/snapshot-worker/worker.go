package main

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/service/mq_service"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/key"
	"github.com/XingfenD/yoresee_doc/pkg/mq"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
)

type Worker struct {
	collabCoreHTTP string
	mqGroup        string
	topic          string
	mqBackend      mq.Backend
	client         *http.Client
	dirtySetKey    string
	inFlight       *sync.Map
}

func (w *Worker) runMQConsumer(mqSvc *mq_service.MQService) {
	err := mqSvc.Consume(
		context.Background(),
		w.mqBackend,
		mq.ConsumeOptions{
			Topic:   w.topic,
			Mode:    mq.ConsumeModeGroup,
			Group:   w.mqGroup,
			AutoAck: false,
			OnError: mq.ErrorActionRequeue,
		},
		func(ctx context.Context, message mq.Message) error {
			docID := parseDocID(message.Body)
			if docID == "" {
				return nil
			}
			return w.snapshotDoc(ctx, docID, true)
		},
	)
	if err != nil {
		logrus.Fatalf("MQ consumer failed: %v", err)
	}
}

func (w *Worker) runScanLoop() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		w.scanAndSnapshot()
	}
}

func (w *Worker) scanAndSnapshot() {
	ctx := context.Background()
	rdb := storage.GetRedis()
	if rdb == nil {
		return
	}

	docIDs, err := rdb.SMembers(ctx, w.dirtySetKey).Result()
	if err != nil {
		logrus.Errorf("Scan dirty docs failed: %v", err)
		return
	}

	candidates := filterCandidates(ctx, docIDs)
	logrus.Infof("Scan dirty docs: candidates=%d", len(candidates))
	if len(candidates) > 0 {
		logrus.Infof("Dirty doc candidates: %v", candidates)
	}
	for _, docID := range candidates {
		_ = w.snapshotDoc(ctx, docID, false)
	}
}

func filterCandidates(ctx context.Context, docIDs []string) []string {
	now := time.Now().UnixMilli()
	candidates := make([]string, 0, len(docIDs))
	for _, docID := range docIDs {
		if docID == "" {
			continue
		}
		lastStr, err := storage.GetRedis().Get(ctx, key.KeyCollabRoom(docID)).Result()
		if err != nil {
			logrus.Infof("Dirty doc %s skipped: room key missing", docID)
			continue
		}
		last, err := utils.ParseInt64(lastStr)
		if err != nil {
			logrus.Infof("Dirty doc %s skipped: invalid room timestamp %s", docID, lastStr)
			continue
		}
		if now-last < 10_000 {
			logrus.Infof("Dirty doc %s skipped: lastEditAgo=%dms", docID, now-last)
			continue
		}
		candidates = append(candidates, docID)
	}
	return candidates
}
