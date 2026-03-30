package errs

import (
	"errors"
	"fmt"
)

var (
	ErrTopicEmpty                = errors.New("topic is empty")
	ErrRedisClientNotInitialized = errors.New("redis client not initialized")
	ErrInvalidConsumeMode        = errors.New("invalid consume mode")
	ErrInvalidErrorAction        = errors.New("invalid error action")

	ErrRabbitMQConnectFailed      = errors.New("failed to connect to RabbitMQ")
	ErrRabbitMQOpenChannelFailed  = errors.New("failed to open a channel")
	ErrRabbitMQDeclareExchange    = errors.New("failed to declare an exchange")
	ErrRabbitMQPublishMessage     = errors.New("failed to publish a message")
	ErrRabbitMQDeclareQueue       = errors.New("failed to declare a queue")
	ErrRabbitMQBindQueue          = errors.New("failed to bind a queue")
	ErrRabbitMQRegisterConsumer   = errors.New("failed to register a consumer")
	ErrRabbitMQConsumerChanClosed = errors.New("consumer channel closed")

	ErrInitRedisClient = errors.New("init redis client failed")

	ErrElasticAddressesEmpty = errors.New("elasticsearch addresses are empty")
	ErrElasticPingFailed     = errors.New("ping elasticsearch failed")
	ErrElasticClientNil      = errors.New("elasticsearch client is nil")
	ErrElasticIndexOrDocID   = errors.New("index or docID is empty")
	ErrElasticIndexEmpty     = errors.New("index is empty")
	ErrElasticStatus         = errors.New("elasticsearch status error")
	ErrElasticStatusDecode   = errors.New("elasticsearch status decode error")
	ErrElasticStatusWithBody = errors.New("elasticsearch status body error")

	ErrMinioInitClient        = errors.New("init minio client failed")
	ErrMinioCheckBucketExists = errors.New("check bucket exists failed")
	ErrMinioCreateBucket      = errors.New("create bucket failed")
	ErrMinioUploadFile        = errors.New("upload file failed")
	ErrMinioGeneratePresigned = errors.New("generate presigned url failed")

	ErrConsulKVGet             = errors.New("consul kv get failed")
	ErrConsulKVSet             = errors.New("consul kv set failed")
	ErrConsulBindTargetNil     = errors.New("consul config bind target is nil")
	ErrConsulClientNil         = errors.New("consul client is nil")
	ErrConsulBindTargetInvalid = errors.New("consul config bind target must be pointer to struct")
	ErrConsulFieldMustBeFunc   = errors.New("field must be func()T")
	ErrConsulFieldTagInvalid   = errors.New("field tag invalid")
	ErrConsulTagEmpty          = errors.New("empty tag")
	ErrConsulTagMissingKey     = errors.New("missing key")
	ErrConsulUnsupportedType   = errors.New("unsupported type")

	ErrInvalidGormLogLevel = errors.New("invalid gorm log level")

	ErrLockAcquireFailed = errors.New("failed to acquire lock")
	ErrLockNotHeld       = errors.New("lock not held or expired")
)

func Wrap(base error, cause error) error {
	if cause == nil {
		return base
	}
	return fmt.Errorf("%w: %v", base, cause)
}

func Detail(base error, detail string) error {
	if detail == "" {
		return base
	}
	return fmt.Errorf("%w: %s", base, detail)
}

func Detailf(base error, format string, args ...any) error {
	return Detail(base, fmt.Sprintf(format, args...))
}

func DetailWrap(base error, detail string, cause error) error {
	if cause == nil {
		return Detail(base, detail)
	}
	if detail == "" {
		return Wrap(base, cause)
	}
	return fmt.Errorf("%w: %s: %v", base, detail, cause)
}
