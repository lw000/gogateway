// tuyue_gateway project main.go
package main

import (
	"github.com/judwhite/go-svc/svc"
	"github.com/radovskyb/watcher"
	log "github.com/sirupsen/logrus"
	"gogateway/backend"
	"gogateway/frontend"
	"gogateway/global"
	"time"
)

type Program struct {
	w        *watcher.Watcher
	htpServe GinServer
}

func (p *Program) Init(env svc.Environment) error {
	if env.IsWindowsService() {

	} else {

	}

	// 加载全局配置
	if err := global.LoadGlobalConfig(); err != nil {
		log.Error(err)
		return err
	}

	p.w = watcher.New()
	p.w.SetMaxEvents(1)

	// Only notify Write events.
	p.w.FilterOps(watcher.Write)

	return nil
}

// Start is called after Init. This method must be non-blocking.
func (p *Program) Start() error {
	var err error
	if err = backend.Instance().Start(); err != nil {
		log.Error(err)
		return err
	}

	// TODO: 启动网关服务
	if err = frontend.Instance().Start(); err != nil {
		log.Error(err)
		return err
	}

	frontend.Instance().AddMessageHook(frontend.CheckMessageCode())

	// TODO: 启动http服务
	err = p.htpServe.Start(global.ProjectConfig.Debug, global.ProjectConfig.Servers.Server)
	if err != nil {
		log.Error(err)
		return err
	}

	go func() {
		for {
			select {
			case event := <-p.w.Event:
				log.Println(event) // Print the event's info.
			case err := <-p.w.Error:
				log.Fatalln(err)
			case <-p.w.Closed:
				return
			}
		}
	}()

	// Watch this folder for changes.
	if err := p.w.Add("."); err != nil {
		log.Fatalln(err)
	}

	// Print a list of all of the files and folders currently
	// being watched and their paths.
	for path, f := range p.w.WatchedFiles() {
		log.Printf("%s: %s\n", path, f.Name())
	}

	// Watch test_folder recursively for changes.
	if err := p.w.AddRecursive("./conf/conf.json"); err != nil {
		log.Println(err)
		return err
	}

	// Trigger 2 events after watcher started.
	go func() {
		p.w.Wait()
		p.w.TriggerEvent(watcher.Create, nil)
		p.w.TriggerEvent(watcher.Remove, nil)
	}()

	go func() {
		// Start the watching process - it'll check for changes every 100ms.
		if err := p.w.Start(time.Millisecond * 100); err != nil {
			log.Println(err)
		}
	}()

	return nil
}

// Stop is called in response to syscall.SIGINT, syscall.SIGTERM, or when a
// Windows Service is stopped.
func (p *Program) Stop() error {
	p.w.Close()
	p.htpServe.Stop()
	frontend.Instance().Stop()
	backend.Instance().Stop()
	log.Error("GATEWAY·服务退出")
	return nil
}

func main() {
	err := svc.Run(&Program{})
	if err != nil {
		log.Error(err)
	}
}
