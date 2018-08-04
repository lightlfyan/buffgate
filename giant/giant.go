package giant

import (
	"github.com/globalsign/mgo"
	"github.com/golang/glog"
	"github.com/lightlfyan/buffgate/config"
	"github.com/lightlfyan/buffgate/model"
	"sync"
	"time"
)

// -----------------------------------------------------------------------------
type Giant struct {
	stomach chan *model.ClientEvent
	close   chan int
	session *mgo.Session

	r_buff []*model.ClientEvent
	w_buff []*model.ClientEvent

	rw *sync.RWMutex

	pool *sync.Pool
}

func (g *Giant) Live() {
	g.stomach = make(chan *model.ClientEvent, 10240)

	g.r_buff = make([]*model.ClientEvent, 0, 128)
	g.w_buff = make([]*model.ClientEvent, 0, 128)

	g.rw = &sync.RWMutex{}
	g.pool = &sync.Pool{
		New: func() interface{} {
			return new(model.ClientEvent)
		},
	}

	mgo.SetDebug(true)

	// do some reconnect
	session, err := mgo.Dial(config.Config.MgoUrl)
	if err != nil {
		panic(err)
	}

	g.session = session
	t := time.Tick(time.Second * 10)

	c := g.session.DB("gsensor").C("clientlog")
	q_size := 0

	for {
		q_size = len(g.stomach)

		select {
		case l := <-g.stomach:
			if q_size >= 10 {
				g.rw.Lock()
				g.w_buff = append(g.w_buff, l)
				g.rw.Unlock()
			} else {
				err := c.Insert(l)
				if err != nil {
					glog.Error(*l)
				}
				g.pool.Put(l)
			}
		case <-t:
			glog.Infoln("queue len: ", q_size)
		}
	}
}

func (g *Giant) flush() {
	for {
		if len(g.r_buff) <= 0 {
			if len(g.w_buff) <= 0 {
				time.Sleep(time.Millisecond * 100)
				continue
			} else {
				g.swap()
			}
		}

		bulk := g.session.DB("gsensor").C("clientlog").Bulk()
		glog.Info("flush: ", len(g.r_buff), len(g.stomach))

		idx_end := 0
		offset := 1024

		for i := 0; i < len(g.r_buff); i += offset {
			idx_end = i + offset
			if len(g.r_buff)-i < offset {
				idx_end = len(g.r_buff)
			}
			for _, l := range g.r_buff[i:idx_end] {
				bulk.Insert(l)
			}
			_, err := bulk.Run()
			if err != nil {
				for _, l := range g.r_buff[i:idx_end] {
					glog.Error(*l)
				}
			}
		}

		for _, l := range g.r_buff {
			g.pool.Put(l)
		}
		g.r_buff = g.r_buff[0:0]
	}
}

func (g *Giant) swap() {
	g.rw.RLock()
	g.r_buff, g.w_buff = g.w_buff, g.r_buff
	g.rw.RUnlock()
}

var g *Giant

func Start() {
	g = &Giant{}
	go g.flush()
	g.Live()
}

func GetEvent() *model.ClientEvent {
	l := g.pool.Get().(*model.ClientEvent)
	l.Reset()
	return l
}

func Send(l *model.ClientEvent) {
	g.stomach <- l
}
