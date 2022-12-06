package data

import "time"

// 인터페이스를 만들기 위한 구조체 정의
type memHandler struct {
	pageMap map[int]*Page
}

// type memRealTimePerfHandler struct {
// 	realtimeperfMap map[int]*RealTimePerf
// }

func (m *memHandler) GetPages() []*Page {
	list := []*Page{}
	for _, v := range m.pageMap {
		list = append(list, v)
	}
	return list
}

func (m *memHandler) GetPage(index int) *Page {
	page, ok := m.pageMap[index]
	if ok {
		return page
	}
	return nil
}

func (m *memHandler) AddPage(page *Page) bool {
	//page := &Page{index, contents, time.Now()}
	page.UpdatedAt = time.Now()
	m.pageMap[page.Index] = page
	return true
}

func (m *memHandler) UpdatePage(page *Page) bool {

	orgPage, ok := m.pageMap[page.Index]
	if ok {
		orgPage.Contents = page.Contents
		orgPage.UpdatedAt = time.Now()
		return true
	}
	return false
}

func (m *memHandler) DeletePage() bool {
	return true
}

// func (m2 *memRealTimePerfHandler) GetRealTimePerf(index int, stime int, etime int) *RealTimePerf {
// 	real, ok := m2.realtimeperfMap[index]
// 	if ok {
// 		return real
// 	}
// 	return nil
// }

func (m *memHandler) Close() {

}

// func newMemHandler() DBHandler {
// 	m := &memHandler{}
// 	m.pageMap = make(map[int]*Page)

// 	return m
// }
