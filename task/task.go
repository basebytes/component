package task

import (
	"fmt"
	"sync"

	"github.com/basebytes/scheduler"
	"github.com/basebytes/types"
)

var (
	taskMap   map[string]scheduler.Task
	configMap map[string]*scheduler.TaskConfig
	lock      sync.RWMutex
	once      sync.Once
)

func Init(defaults map[string]any, configs []*scheduler.TaskConfig, tasks ...scheduler.Task) {
	once.Do(func() {
		lock.Lock()
		defer lock.Unlock()
		if _taskMap, _configMap, err := reload(defaults, configs, tasks...); err == nil {
			taskMap = _taskMap
			configMap = _configMap
		} else {
			panic(err)
		}
	})
}

func Reload(defaults map[string]any, configs []*scheduler.TaskConfig, tasks ...scheduler.Task) (err error) {
	var (
		_taskMap   map[string]scheduler.Task
		_configMap map[string]*scheduler.TaskConfig
	)
	if _taskMap, _configMap, err = reload(defaults, configs, tasks...); err == nil {
		lock.Lock()
		defer lock.Unlock()
		taskMap = _taskMap
		configMap = _configMap
	}
	return
}

func reload(defaults map[string]any, configs []*scheduler.TaskConfig, tasks ...scheduler.Task) (
	taskMap map[string]scheduler.Task, configMap map[string]*scheduler.TaskConfig, err error) {
	taskMap = make(map[string]scheduler.Task, len(tasks))
	configMap = make(map[string]*scheduler.TaskConfig, len(tasks))
	for _, task := range tasks {
		code := task.Code()
		if _, OK := taskMap[code]; OK {
			err = fmt.Errorf("duplicate task code %s", code)
			return
		}
		taskMap[code] = task
	}
	for _, taskConfig := range configs {
		if _, OK := taskMap[taskConfig.Code]; OK {
			taskConfig.SetDefaults(defaults)
			configMap[taskConfig.Code] = taskConfig
		}
	}
	return
}

func Start() (err error) {
	lock.RLock()
	defer lock.RUnlock()
	for code, taskConfig := range configMap {
		if err = registerTask(taskMap[code], taskConfig); err != nil {
			break
		}
	}
	if err == nil {
		scheduler.Start()
	}
	return
}

func Stop(wait bool) {
	ctx := scheduler.Stop()
	if wait {
		<-ctx.Done()
	} else {
		go func() { <-ctx.Done() }()
	}
}

func AddTask(code string) (err error) {
	lock.RLock()
	defer lock.RUnlock()
	if task, ok := taskMap[code]; !ok {
		err = scheduler.TaskNotFoundErr
	} else if scheduler.ExistTask(code) {
		err = scheduler.TaskExistErr
	} else {
		err = registerTask(task, configMap[code])
	}
	return
}

func CancelTask(code string, wait bool) {
	ctx := scheduler.CancelTask(code)
	if wait {
		<-ctx.Done()
	} else {
		go func() { <-ctx.Done() }()
	}
}

func PauseTask(code string) {
	scheduler.PauseTask(code)
}

func ResumeTask(code string) {
	scheduler.ResumeTask(code)
}

func Status() map[string]int32 {
	taskStatuses := scheduler.TaskStatus()
	lock.RLock()
	defer lock.RUnlock()
	if len(taskMap) != len(taskStatuses) {
		for code := range taskMap {
			if _, ok := taskStatuses[code]; !ok {
				taskStatuses[code] = scheduler.Stopped
			}
		}
	}
	return taskStatuses
}

func AddPlan(code, plan string) (err error) {
	lock.RLock()
	defer lock.RUnlock()
	if taskConfig, ok := configMap[code]; ok {
		for _, cfg := range taskConfig.Plans {
			if cfg.Key() == plan {
				return scheduler.AddPlan(code, cfg)
			}
		}
		err = scheduler.PlanNotFoundErr
	} else {
		err = scheduler.TaskNotFoundErr
	}
	return
}

func PausePlan(code, plan string) error {
	return scheduler.PausePlan(code, plan)
}

func ResumePlan(code, plan string) error {
	return scheduler.ResumePlan(code, plan)
}

func CancelPlan(code, plan string) error {
	return scheduler.CancelPlan(code, plan)
}

func ExecuteTask(code string, param *runParam) (string, error) {
	return scheduler.ExecuteTask(code, param.Base, param.MaxDelay)
}

func ExistTask(code string) (ok bool) {
	lock.RLock()
	defer lock.RUnlock()
	_, ok = taskMap[code]
	return
}

func Snapshot(code string) (status *taskStatus, err error) {
	lock.RLock()
	defer lock.RUnlock()
	if cfg, ok := configMap[code]; ok {
		status = &taskStatus{Plans: make([]*PlanStatus, 0, len(cfg.Plans))}
		taskStatuses := scheduler.TaskStatus(code)
		if status.Status, ok = taskStatuses[code]; ok {
			planStatus := scheduler.TaskPlanStatus(code)
			jobStatus := scheduler.TaskJobStatus(code)
			for _, planCfg := range cfg.Plans {
				ps := &PlanStatus{Key: planCfg.Key()}
				if ps.Status, ok = planStatus[ps.Key]; !ok {
					ps.Status = scheduler.Stopped
				}
				ps.Jobs = jobStatus[ps.Key]
				status.Plans = append(status.Plans, ps)
			}
		} else {
			status.Status = scheduler.Stopped
		}
	} else if _, ok = taskMap[code]; ok {
		status = stoppedTaskStatus
	} else {
		err = scheduler.TaskNotFoundErr
	}
	return
}

type taskStatus struct {
	Status int32         `json:"status"`
	Plans  []*PlanStatus `json:"plans,omitempty"`
}

type PlanStatus struct {
	Key    string           `json:"key"`
	Status int32            `json:"status"`
	Jobs   map[string]int32 `json:"jobs,omitempty"`
}

var stoppedTaskStatus = &taskStatus{Status: scheduler.Stopped}

func NewRunParam() *runParam {
	return &runParam{MaxDelay: &types.Duration{}}
}

type runParam struct {
	Base     *types.Time     `json:"base,omitempty"`
	MaxDelay *types.Duration `json:"maxDelay,omitempty"`
}

func registerTask(task scheduler.Task, taskConfig *scheduler.TaskConfig) (err error) {
	return scheduler.RegisterTask(task, taskConfig)
}
