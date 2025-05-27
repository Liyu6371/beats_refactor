package beats

import (
	"beats_refactor/config"
	"beats_refactor/logger"
	"context"
	"sync"
)

type Beats struct {
	wg sync.WaitGroup

	ctx    context.Context
	cancel context.CancelFunc
}

func New(c context.Context) (*Beats, error) {
	// 配置解析
	if _, err := config.InitConfig(); err != nil {
		return nil, err
	}
	// 初始化日志
	if err := logger.InitLogger(*config.GetLoggerConfig()); err != nil {
		return nil, err
	}
	// 尝试按照配置拉起 Task 任务
	ctx, cancel := context.WithCancel(c)
	return &Beats{
		wg:     sync.WaitGroup{},
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (b *Beats) Start() {
	return
}
