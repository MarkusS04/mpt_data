// Package apimodel contains all models used by api to exchange data with client
package apimodel

// OrderTask is type for client to send only a taskId and the OrderNumber that this task should get
type OrderTask struct {
	TaskID      uint
	OrderNumber uint
}

// OrderTaskDetail is type for client to send only a taskDetailId and the OrderNumber that this taskDetail should get
type OrderTaskDetail struct {
	TaskDetailID uint
	OrderNumber  uint
}
