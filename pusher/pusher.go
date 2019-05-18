package pusher

import (
	"fmt"
)

type PusherI interface {
	Push(msg string) error
	Init(config map[string]interface{}) error
}

// Errors
var (
	ErrNameAlreadyExists = fmt.Errorf("Pusher: name already exists")
	ErrNotFoundPusher    = fmt.Errorf("Pusher: not found pusher by name")
)

var (
	pushersMap = map[string]PusherI{}
)

func AppendPusher(name string, pusher PusherI) (err error) {
	if pushersMap[name] != nil {
		return ErrNameAlreadyExists
	}

	pushersMap[name] = pusher
	return
}

func Push(pusherName, msg string) (err error) {
	pusher := pushersMap[pusherName]
	if pusher == nil {
		return ErrNotFoundPusher
	}

	return pusher.Push(msg)
}

func Run(configs map[string]map[string]interface{}) (err error) {
	for name, config := range configs {
		pusher := pushersMap[name]
		if pusher == nil {
			return ErrNotFoundPusher
		}

		if err = pusher.Init(config); err != nil {
			return
		}
	}

	return
}
