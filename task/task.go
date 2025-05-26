package task

import (
	"beats_refactor/config"
	"beats_refactor/pipeline"
	"beats_refactor/sender"
	"beats_refactor/source"
	"context"
	"sync"
)

// TaskService 配置项里一个任务会对应一个 TaskService 实例
// 每个 TaskService 实例会包含一个 SourceService、SenderService 和 PipelineService
// 这些服务会根据配置文件中的 Source、Sender 和 Pipeline 字段进行初始化
type TaskInstance struct {
	wg sync.WaitGroup

	ctx    context.Context
	cancel context.CancelFunc

	TaskConf        config.TaskConfig
	SourceService   *source.SourceService
	SenderService   *sender.SenderService
	PipelineService pipeline.PipelineService
}

// NewTaskService 实例化 TaskSerive
func NewTaskService(c config.TaskConfig, ctx context.Context) (*TaskInstance, error) {
	if c.IsCloudMonitorTask() {
		return newCloudMonitorTaskInstance(c, ctx)
	}
	return newSourceMonitorTaskInstance(c, ctx)
}

// newCloudMonitorTaskInstance 创建云监控任务实例
func newCloudMonitorTaskInstance(c config.TaskConfig, ctx context.Context) (*TaskInstance, error) {
	return nil, nil
}

// newSourceMonitorTaskInstance 创建源监控任务实例
func newSourceMonitorTaskInstance(c config.TaskConfig, ctx context.Context) (*TaskInstance, error) {
	context, cancel := context.WithCancel(ctx)
	sourceService, err := source.NewSourceService(c.Source, ctx)
	if err != nil {
		return nil, err
	}

	senderService, err := sender.NewSenderService(c.Sender, ctx)
	if err != nil {
		return nil, err
	}
	return &TaskInstance{
		wg:            sync.WaitGroup{},
		ctx:           context,
		cancel:        cancel,
		TaskConf:      c,
		SourceService: sourceService,
		SenderService: senderService,
	}, nil
}

func (t *TaskInstance) Start() {
}
