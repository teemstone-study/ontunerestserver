package data

import "time"

type Page struct {
	Index     int       `json:"index"`
	Contents  string    `json:"contents"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RealTimePerf struct {
	Ontunetime       int64 `json:"_ontunetime"`
	Agenttime        int64 `json:"_agenttime"`
	Agentid          int   `json:"_agentid"`
	User             int   `json:"_user"`
	Sys              int   `json:"_sys"`
	Wait             int   `json:"_wait"`
	Idle             int   `json:"_idle"`
	Processorcount   int   `json:"_processorcount"`
	Runqueue         int   `json:"_runqueue"`
	Blockqueue       int   `json:"_blockqueue"`
	Waitqueue        int   `json:"_waitqueue"`
	Pqueue           int   `json:"_pqueue"`
	Pcrateuser       int   `json:"_pcrateuser"`
	Pcratesys        int   `json:"_pcratesys"`
	Memorysize       int   `json:"_memorysize"`
	Memoryused       int   `json:"_memoryused"`
	Memorypinned     int   `json:"_memorypinned"`
	Memorysys        int   `json:"_memorysys"`
	Memoryuser       int   `json:"_memoryuser"`
	Memorycache      int   `json:"_memorycache"`
	Avm              int   `json:"_avm"`
	Pagingspacein    int   `json:"_pagingspacein"`
	Pagingspaceout   int   `json:"_pagingspaceout"`
	Filesystemin     int   `json:"_filesystemin"`
	Filesystemout    int   `json:"_filesystemout"`
	Memoryscan       int   `json:"_memoryscan"`
	Memoryfreed      int   `json:"_memoryfreed"`
	Swapsize         int   `json:"_swapsize"`
	Swapused         int   `json:"_swapused"`
	Swapactive       int   `json:"_swapactive"`
	Fork             int   `json:"_fork"`
	Exec             int   `json:"_exec"`
	Interupt         int   `json:"_interupt"`
	Systemcall       int   `json:"_systemcall"`
	Contextswitch    int   `json:"_contextswitch"`
	Semaphore        int   `json:"_semaphore"`
	Msg              int   `json:"_msg"`
	Diskreadwrite    int   `json:"_diskreadwrite"`
	Diskiops         int   `json:"_diskiops"`
	Networkreadwrite int   `json:"_networkreadwrite"`
	Networkiops      int   `json:"_networkiops"`
	Topcommandid     int   `json:"_topcommandid"`
	Topcommandcount  int   `json:"_topcommandcount"`
	Topuserid        int   `json:"_topuserid"`
	Topcpu           int   `json:"_topcpu"`
	Topdiskid        int   `json:"_topdiskid"`
	Topvgid          int   `json:"_topvgid"`
	Topbusy          int   `json:"_topbusy"`
	Maxpid           int   `json:"_maxpid"`
	Threadcount      int   `json:"_threadcount"`
	Pidcount         int   `json:"_pidcount"`
	Linuxbuffer      int   `json:"_linuxbuffer"`
	Linuxcached      int   `json:"_linuxcached"`
	Linuxsrec        int   `json:"_linuxsrec"`
	Memused_Mb       int   `json:"_memused_Mb"`
	Irq              int   `json:"_irq"`
	Softirq          int   `json:"_softirq"`
	Swapused_Mb      int   `json:"_swapused_Mb"`
	Dusm             int   `json:"_dusm"`
}

type RealTimeCpu struct {
	Ontunetime    int64 `json:"_ontunetime"`
	Agenttime     int64 `json:"_agenttime"`
	Agentid       int   `json:"_agentid"`
	Index         int   `json:"_index"`
	User          int   `json:"_user"`
	Sys           int   `json:"_sys"`
	Wait          int   `json:"_wait"`
	Idle          int   `json:"_idle"`
	Runqueue      int   `json:"_runqueue"`
	Fork          int   `json:"_fork"`
	Exec          int   `json:"_exec"`
	Interupt      int   `json:"_interupt"`
	Systemcall    int   `json:"_systemcall"`
	Contextswitch int   `json:"_contextswitch"`
}

type RealTimeNet struct {
	Ontunetime int64 `json:"_ontunetime"`
	Agenttime  int64 `json:"_agenttime"`
	Agentid    int   `json:"_agentid"`
	Ionameid   int   `json:"_ionameid"`
	Readrate   int   `json:"_readrate"`
	Writerate  int   `json:"_writerate"`
	Readiops   int   `json:"_readiops"`
	Writeiops  int   `json:"_writeiops"`
	Errorps    int   `json:"_errorps"`
	Collision  int   `json:"_collision"`
}

type RealTimeDisk struct {
	Ontunetime   int64 `json:"_ontunetime"`
	Agenttime    int64 `json:"_agenttime"`
	Agentid      int   `json:"_agentid"`
	Ionameid     int   `json:"_ionameid"`
	Readrate     int   `json:"_readrate"`
	Writerate    int   `json:"_writerate"`
	Iops         int   `json:"_iops"`
	Busy         int   `json:"_busy"`
	Descid       int   `json:"_descid"`
	Readsvctime  int   `json:"_readsvctime"`
	Writesvctime int   `json:"_writesvctime"`
}

type DBHandler interface {
	GetPages() []*Page
	GetPage(index int) *Page
	AddPage(page *Page) bool
	UpdatePage(page *Page) bool
	DeletePage() bool
	//GetAllRealTimePerf() []*RealTimePerf
	GetRealTimePerf(ids string, stime int, etime int) *RealTimePerf
	GetRealTimeCpu(ids string, stime int, etime int) *RealTimeCpu
	GetRealTimeNet(ids string, stime int, etime int) *RealTimeNet
	GetRealTimeDisk(ids string, stime int, etime int) *RealTimeDisk
	Close()
}

func NewDBHandler(dbConn string) DBHandler {
	//return newMemHandler()
	//return newSqliteHandler(filepath)
	return newPgHandler(dbConn)
}
