package gonet

// /////////////////////////////
// ///    Option Func   ////////
// ////////////////////////////
type Option func(*AppContext) error

func MaxSessions(max int) Option {
	return func(o *AppContext) error {
		o.maxSessionCount = max
		return nil
	}
}

func WorkerPoolMaxSize(max int) Option {
	return func(o *AppContext) error {
		o.maxWorkerPoolSize = max
		return nil
	}
}

// cache for messages
func WithMessageCache(cache IMessageCache) Option {
	return func(o *AppContext) error {
		o.msgCache = cache
		return nil
	}
}

// message codec,default is json codec
func WithMessageCodec(codec ICodec) Option {
	return func(o *AppContext) error {
		o.codec = codec
		return nil
	}
}
