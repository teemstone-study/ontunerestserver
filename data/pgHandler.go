package data

import (
	"database/sql" //모든 sql 에 공통되는 interface ,
	"time"

	_ "github.com/lib/pq" //sqlite 패키지도 import 해야한다. _ 는 암시적 사용 예시
)

// 인터페이스를 만들기 위한 구조체 정의
type pgHandler struct {
	db *sql.DB
}

func (p *pgHandler) GetPages() []*Page {
	pages := []*Page{}

	rows, err := p.db.Query("SELECT pageindex, contents, updatedat FROM pages")
	if err != nil {
		panic(err)
	}
	defer rows.Close() //function 이 종료 되기전에 rows 를 Close 시켜라

	for rows.Next() {
		var p Page
		rows.Scan(&p.Index, &p.Contents, &p.UpdatedAt)
		pages = append(pages, &p)
	}

	return pages
}

func (p *pgHandler) GetPage(index int) *Page {
	rows, err := p.db.Query("SELECT pageindex, contents, updatedat FROM pages WHERE pageindex=$1", index)
	if err != nil {
		panic(err)
	}
	defer rows.Close() //function 이 종료 되기전에 rows 를 Close 시켜라
	//log.Println(">>>>>>>>>> pageindex >>>>>> " +  strconv.Itoa(index))
	if rows.Next() {
		var pg Page //= new(Page)
		rows.Scan(&pg.Index, &pg.Contents, &pg.UpdatedAt)
		return &pg
	} else {
		return nil
	}
}

func (p *pgHandler) AddPage(page *Page) bool {
	stmt, err := p.db.Prepare("INSERT INTO pages (pageindex, contents, updatedat) VALUES ($1, $2, NOW())")
	if err != nil {
		panic(err)
	}

	rst, err := stmt.Exec(page.Index, page.Contents)
	if err != nil {
		panic(err)
	}

	page.UpdatedAt = time.Now()
	cnt, _ := rst.RowsAffected()
	return (cnt > 0)
}

func (p *pgHandler) UpdatePage(page *Page) bool {

	stmt, err := p.db.Prepare("UPDATE pages SET contents=$1, updatedat=NOW() WHERE pageindex=$2")
	if err != nil {
		panic(err)
	}

	rst, err := stmt.Exec(page.Contents, page.Index)
	if err != nil {
		panic(err)
	}

	page.UpdatedAt = time.Now()
	cnt, _ := rst.RowsAffected()
	return (cnt > 0)
}

func (p *pgHandler) DeletePage() bool {

	stmt, err := p.db.Prepare("DELETE FROM pages")
	if err != nil {
		panic(err)
	}

	rst, err := stmt.Exec()
	if err != nil {
		panic(err)
	}
	cnt, _ := rst.RowsAffected()
	//return true
	return cnt >= 0
}

// To do : 인자로 받아서 테이블 명을 지정해야함.
func (p *pgHandler) GetRealTimePerf(ids string, stime int, etime int) *RealTimePerf {
	rows, err := p.db.Query("SELECT _ontunetime,  _agenttime,     _agentid,         _user,         _sys,             _wait,       _idle,         _processorcount, _runqueue,   _blockqueue,       "+
		"				_waitqueue,   _pqueue,        _pcrateuser,      _pcratesys,    _memorysize,      _memoryused, _memorypinned, _memorysys,      _memoryuser, _memorycache, 			"+
		"				_avm,         _pagingspacein, _pagingspaceout,  _filesystemin, _filesystemout,   _memoryscan, _memoryfreed,  _swapsize,       _swapused,   _swapactive,  			"+
		"				_fork,        _exec,          _interupt,        _systemcall,   _contextswitch,   _semaphore,  _msg,          _diskreadwrite,  _diskiops,   _networkreadwrite, "+
		"				_networkiops, _topcommandid,  _topcommandcount, _topuserid,    _topcpu,          _topdiskid,  _topvgid,      _topbusy,        _maxpid,     _threadcount,      "+
		"				_pidcount,    _linuxbuffer,   _linuxcached,     _linuxsrec,    _memused_Mb,      _irq,        _softirq,      _swapused_Mb,    _dusm                           "+
		"FROM realtimeperf_22120109 WHERE _agentid in ("+ids+") and _ontunetime=$1 and _agenttime=$2", stime, etime)
	if err != nil {
		panic(err)
	}
	defer rows.Close() //function 이 종료 되기전에 rows 를 Close 시켜라
	//log.Println(">>>>>>>>>> pageindex >>>>>> " +  strconv.Itoa(index))
	if rows.Next() {
		var real RealTimePerf //= new(Page)
		//  rows.Scan(&pg.Index, &pg.Contents, &pg.UpdatedAt)
		rows.Scan(&real.Ontunetime, &real.Agenttime, &real.Agentid, &real.User, &real.Sys, &real.Wait, &real.Idle, &real.Processorcount, &real.Runqueue, &real.Blockqueue,
			&real.Waitqueue, &real.Pqueue, &real.Pcrateuser, &real.Pcratesys, &real.Memorysize, &real.Memoryused, &real.Memorypinned, &real.Memorysys, &real.Memoryuser, &real.Memorycache,
			&real.Avm, &real.Pagingspacein, &real.Pagingspaceout, &real.Filesystemin, &real.Filesystemout, &real.Memoryscan, &real.Memoryfreed, &real.Swapsize, &real.Swapused, &real.Swapactive,
			&real.Fork, &real.Exec, &real.Interupt, &real.Systemcall, &real.Contextswitch, &real.Semaphore, &real.Msg, &real.Diskreadwrite, &real.Diskiops, &real.Networkreadwrite,
			&real.Networkiops, &real.Topcommandid, &real.Topcommandcount, &real.Topuserid, &real.Topcpu, &real.Topdiskid, &real.Topvgid, &real.Topbusy, &real.Maxpid, &real.Threadcount,
			&real.Pidcount, &real.Linuxbuffer, &real.Linuxcached, &real.Linuxsrec, &real.Memused_Mb, &real.Irq, &real.Softirq, &real.Swapused_Mb, &real.Dusm)
		return &real
	} else {
		return nil
	}
}

// To do : 인자로 받아서 테이블 명을 지정해야함.
func (p *pgHandler) GetRealTimeCpu(ids string, stime int, etime int) *RealTimeCpu {
	rows, err := p.db.Query("SELECT _ontunetime, _agenttime, _agentid, _index, _user, _sys, _wait, _idle,	_runqueue, _fork,	_exec,	_interupt,	_systemcall,	_contextswitch "+
		"FROM realtimecpu_22120109 WHERE _agentid in ("+ids+") and _ontunetime=$1 and _agenttime=$2", stime, etime)
	if err != nil {
		panic(err)
	}
	defer rows.Close() //function 이 종료 되기전에 rows 를 Close 시켜라
	//log.Println(">>>>>>>>>> pageindex >>>>>> " +  strconv.Itoa(index))
	if rows.Next() {
		var cpu RealTimeCpu //= new(Page)
		//  rows.Scan(&pg.Index, &pg.Contents, &pg.UpdatedAt)
		rows.Scan(&cpu.Ontunetime, &cpu.Agenttime, &cpu.Agentid, &cpu.Index, &cpu.User, &cpu.Sys, &cpu.Wait, &cpu.Idle, &cpu.Runqueue, &cpu.Fork, &cpu.Exec, &cpu.Interupt, &cpu.Systemcall, &cpu.Contextswitch)
		return &cpu
	} else {
		return nil
	}
}

// To do : 인자로 받아서 테이블 명을 지정해야함.
func (p *pgHandler) GetRealTimeNet(ids string, stime int, etime int) *RealTimeNet {
	rows, err := p.db.Query("SELECT _ontunetime, _agenttime, _agentid, _ionameid, _readrate, _writerate, _readiops, _writeiops, _errorps, _collision "+
		"FROM realtimenet_22120109 WHERE _agentid in ("+ids+") and _ontunetime=$1 and _agenttime=$2", stime, etime)
	if err != nil {
		panic(err)
	}
	defer rows.Close() //function 이 종료 되기전에 rows 를 Close 시켜라
	//log.Println(">>>>>>>>>> pageindex >>>>>> " +  strconv.Itoa(index))
	if rows.Next() {
		var net RealTimeNet //= new(Page)
		//  rows.Scan(&pg.Index, &pg.Contents, &pg.UpdatedAt)
		rows.Scan(&net.Ontunetime, &net.Agenttime, &net.Agentid, &net.Ionameid, &net.Readrate, &net.Writerate, &net.Readiops, &net.Writeiops, &net.Errorps, &net.Collision)
		return &net
	} else {
		return nil
	}
}

// To do : 인자로 받아서 테이블 명을 지정해야함.
func (p *pgHandler) GetRealTimeDisk(ids string, stime int, etime int) *RealTimeDisk {
	rows, err := p.db.Query("SELECT _ontunetime, _agenttime, _agentid, _ionameid, _readrate, _writerate, _iops, _busy, _descid, _readsvctime, _writesvctime "+
		"FROM realtimedisk_22120109 WHERE _agentid in ("+ids+") and _ontunetime=$1 and _agenttime=$2", stime, etime)
	if err != nil {
		panic(err)
	}
	defer rows.Close() //function 이 종료 되기전에 rows 를 Close 시켜라
	//log.Println(">>>>>>>>>> pageindex >>>>>> " +  strconv.Itoa(index))
	if rows.Next() {
		var disk RealTimeDisk //= new(Page)
		//  rows.Scan(&pg.Index, &pg.Contents, &pg.UpdatedAt)
		rows.Scan(&disk.Ontunetime, &disk.Agenttime, &disk.Agentid, &disk.Ionameid, &disk.Readrate, &disk.Writerate, &disk.Iops, &disk.Busy, &disk.Descid, &disk.Readsvctime, &disk.Writesvctime)
		return &disk
	} else {
		return nil
	}
}

func (p *pgHandler) Close() {
	p.db.Close()
}

func newPgHandler(dbConn string) DBHandler {
	database, err := sql.Open("postgres", dbConn)
	if err != nil {
		panic(err)
	}

	statement, err := database.Prepare(
		`CREATE TABLE IF NOT EXISTS pages (
			pageindex int PRIMARY KEY,
			contents  TEXT,			
			updatedat TIMESTAMP
		)`)

	if err != nil {
		panic(err)
	}

	_, err = statement.Exec()
	if err != nil {
		panic(err)
	}

	return &pgHandler{db: database}
}
