package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/bootstrap"
	"github.com/XingfenD/yoresee_doc/internal/service/mq_service"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/constant"
	"github.com/XingfenD/yoresee_doc/pkg/key"
	"github.com/XingfenD/yoresee_doc/pkg/mq"
	"github.com/sirupsen/logrus"
)

func main() {
	initializer := bootstrap.NewInitializer().
		InitConfig().
		InitPostgres().
		InitRedis().
		InitElasticsearchAllowFail().
		InitMQ().
		InitRepository()
	if err := initializer.Err(); err != nil {
		logrus.Fatalf("Init snapshot-worker failed: %v", err)
	}

	w := &Worker{
		collabCoreHTTP: utils.GetEnvVar("COLLAB_CORE_HTTP", "http://collab-core:1234"),
		mqGroup:        utils.GetEnvVar("SNAPSHOT_MQ_GROUP", "snapshot-worker"),
		topic:          constant.DirtyDocTopicDefault,
		mqBackend:      mq.BackendRabbitMQ,
		client:         &http.Client{Timeout: 10 * time.Second},
		dirtySetKey:    key.KeyCollabDirtyDocSet(),
		inFlight:       &sync.Map{},
	}

	logrus.Infof("Snapshot worker started: topic=%s group=%s collabCore=%s", w.topic, w.mqGroup, w.collabCoreHTTP)

	go w.runMQConsumer(mq_service.MQSvc)
	go w.runScanLoop()

	initializer.ShutdownOnSignal(0)
}
