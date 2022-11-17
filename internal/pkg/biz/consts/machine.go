package consts

const (
	MachineStatusIdle = 1
	MachineStatusBusy = 2

	MachineListRedisKey   = "RunnerMachineList"
	MachineUseStatePrefix = "MachineUseState:"
	MachineAliveTime      = 3

	MachineMonitorPrefix = "MachineMonitor:" // 压力机监控数据
)
