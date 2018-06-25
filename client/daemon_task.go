package main

import (
	"fmt"
	"jiacrontab/libs/proto"
	"jiacrontab/model"
	"os"
	"path/filepath"
	"time"
)

type DaemonTask struct {
}

func (t *DaemonTask) CreateDaemonTask(args model.DaemonTask, reply *int64) error {

	ret := model.DB().Create(&args)

	*reply = ret.RowsAffected
	return ret.Error
}

func (t *DaemonTask) ListDaemonTask(args struct{ Page, Pagesize int }, reply *[]model.DaemonTask) error {
	ret := model.DB().Find(reply).Offset((args.Page - 1) * args.Pagesize).Limit(args.Pagesize).Order("update_at desc")

	return ret.Error
}

func (t *DaemonTask) ActionDaemonTask(args proto.ActionDaemonTaskArgs, reply *bool) error {

	var task model.DaemonTask

	*reply = false

	ret := model.DB().Find(&task, "id=?", args.TaskId)

	if (task == model.DaemonTask{}) {

		return ret.Error
	}

	globalDaemon.add(&daemonTask{
		task:   &task,
		action: args.Action,
	})
	return nil
}

func (t *DaemonTask) GetDaemonTask(args int, reply *model.DaemonTask) error {
	ret := model.DB().Find(reply, "task_id", args)
	if (*reply == model.DaemonTask{}) {
		return ret.Error
	}
	return nil
}

func (t *DaemonTask) Log(args int, ret *[]byte) error {
	fp := filepath.Join(globalConfig.logPath, "daemon_task", time.Now().Format("2006/01"), fmt.Sprint(args, ".log"))
	f, err := os.Open(fp)
	defer f.Close()
	if err != nil {
		return err
	}
	fStat, err := f.Stat()
	if err != nil {
		return err
	}
	limit := int64(1024 * 1024)
	var offset int64
	var buffer []byte
	if fStat.Size() > limit {
		buffer = make([]byte, limit)
		offset = fStat.Size() - limit
	} else {
		offset = 0
		buffer = make([]byte, fStat.Size())
	}
	f.Seek(offset, os.SEEK_CUR)

	_, err = f.Read(buffer)
	*ret = buffer

	return err
}
