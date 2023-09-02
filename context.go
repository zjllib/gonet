package gonet

import (
	"github.com/zjllib/gonet/v3/codec"
	"reflect"
	"sync"
)

type Context struct {
	//会话管理
	sessionMgr *SessionManager
	//message types
	mMsgTypes map[MessageID]reflect.Type
	//message ids
	mMsgIDs map[reflect.Type]MessageID
	//server types
	sessionType reflect.Type
	//消息编码器
	defaultCodec codec.Codec
	//传输端
	server IServer
	client IClient
	//bee worker pool
	workers BeeWorkerPool
	//消息钩子
	mMsgHooks map[MessageID]Hook

	name string

	globalLock sync.Mutex
}

func (c *Context) Name() string {
	return c.name
}
func (c *Context) Server() IServer {
	return c.server
}
func (c *Context) Client() IClient {
	return c.client
}

// 会话管理
func (c *Context) GetSession(id uint64) (ISession, bool) {
	return c.sessionMgr.getAliveSession(id)
}
func (c *Context) CreateSession() ISession {
	idleSession := c.sessionMgr.getIdleSession()
	session := idleSession.(ISession)
	c.sessionMgr.addAliveSession(idleSession)
	return session
}
func (c *Context) RecycleSession(session ISession, err error) {
	c.HandingMessage(session, &Message{
		ID:   SessionClose,
		Body: err,
	})
	//关闭
	session.Close()
	c.sessionMgr.recycleIdleSession(session)
}
func (c *Context) SessionCount() int {
	return c.sessionMgr.CountAliveSession()
}

// 广播会话
func (c *Context) Broadcast(msg interface{}) {
	c.sessionMgr.alive.Range(func(_, item interface{}) bool {
		session, ok := item.(ISession)
		if ok {
			session.Send(msg)
		}
		return true
	})
}

// 映射消息体
func (c *Context) Route(msgID MessageID, msg any, callback Hook) {
	c.globalLock.Lock()
	defer c.globalLock.Unlock()
	msgType := reflect.TypeOf(msg)
	if _, ok := c.mMsgTypes[msgID]; ok {
		panic("error:Duplicate message id")
	}
	if msgType != nil {
		c.mMsgIDs[msgType] = msgID
		c.mMsgTypes[msgID] = msgType
	}
	if callback != nil {
		c.mMsgHooks[msgID] = callback
	}
}

// 获取消息ID
func (c *Context) GetMsgID(msg interface{}) (MessageID, bool) {
	msgID, ok := c.mMsgIDs[reflect.TypeOf(msg)]
	return msgID, ok
}

// 通消息id创建消息体
func (c *Context) CreateMsg(msgID MessageID) interface{} {
	if msg, ok := c.mMsgTypes[msgID]; ok {
		return reflect.New(msg).Interface()
	}
	return nil
}

// 编码消息
func (c *Context) EncodeMessage(msg interface{}) ([]byte, error) {
	return c.defaultCodec.Encode(msg)
}

// 解码消息
func (c *Context) DecodeMessage(msg interface{}, data []byte) error {
	return c.defaultCodec.Decode(data, msg)
}

// 缓存消息
func (c *Context) HandingMessage(session ISession, msg *Message) {
	msg.Head.setSession(session)
	c.workers.rcvMsgCh <- msg
}