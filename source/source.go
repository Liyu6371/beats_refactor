package source

import (
	"beats_refactor/config"
	"beats_refactor/logger"
	"context"
	"errors"
	"sync"
)

type SourceService struct {
	wg sync.WaitGroup
	ch chan []byte
	mu sync.RWMutex

	ctx    context.Context
	cancel context.CancelFunc

	SourceConf      []*config.SourceConfig
	runningInstance map[string]SourceInstance
}

// check 校验配置中的 Source 实例是否已注册
func check(c []*config.SourceConfig) error {
	for _, source := range c {
		if _, ok := instanceFactory[source.Name]; !ok {
			return errors.New("source instance not registered: " + source.Name)
		}
	}
	return nil
}

// NewSourceService 创建 SourceService 实例
func NewSourceService(c []*config.SourceConfig, ctx context.Context) (*SourceService, error) {
	if err := check(c); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	return &SourceService{
		wg: sync.WaitGroup{},

		ch: make(chan []byte, 100), // 缓冲通道，大小可根据需要调整
		mu: sync.RWMutex{},

		ctx:             ctx,
		cancel:          cancel,
		SourceConf:      c,
		runningInstance: map[string]SourceInstance{},
	}, nil
}

func (s *SourceService) Start() {
	// 在NewSourceService中已经校验了配置，在此处拉起所有的 sourceInstance 实例
	// 传递 ctx 用于控制拉起的具体的 sourceInstance 实例的生命周期
	for _, i := range s.SourceConf {
		fn := instanceFactory[i.Name]
		instance, err := fn(i)
		if err != nil {
			logger.Errorf("create source instance %s failed: %v", i.Name, err)
			continue
		}
		s.wg.Add(1)
		s.addInstance(i.Name, instance)

		// 拉起一定数量的 worker 实例
		// for j := range i.Worker {
		// 	instance, err := fn(i)
		// 	instName := fmt.Sprintf("%s-%d", i.Name, j)
		// 	if err != nil {
		// 		logger.Errorf("create source instance %s worker %d failed: %v", i.Name, j, err)
		// 		continue
		// 	}
		// 	s.wg.Add(1)
		// 	s.addInstance(instName, instance)
		// 	go func(n string) {
		// 		defer func() {
		// 			s.wg.Done()
		// 			s.deleteInstance(n)
		// 		}()
		// 		logger.Debugf("source instance %s worker %d started", i.Name, j)
		// 		instance.Start(s.ctx, s.ch, n)
		// 		logger.Debugf("source instance %s worker %d stopped", i.Name, j)
		// 	}(instName)
		// }
	}
}

func (s *SourceService) addInstance(name string, instance SourceInstance) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.runningInstance[name]; !exists {
		s.runningInstance[name] = instance
	}
}

func (s *SourceService) deleteInstance(name string) {
	s.mu.RLock()
	_, ok := s.runningInstance[name]
	s.mu.RUnlock()
	if !ok {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	// 直接删除，无需检查是否存在
	delete(s.runningInstance, name)
}

func (s *SourceService) GetChan() chan []byte {
	return s.ch
}

func (s *SourceService) Stop() {
	s.cancel()  // 取消上下文，通知所有实例停止
	s.wg.Wait() // 等待所有 goroutine 完成
	close(s.ch) // 关闭通道，通知所有接收者
	logger.Info("source service stopped")
}
