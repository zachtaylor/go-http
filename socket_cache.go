package http

import (
	"sync"
)

var socketCache = make(map[string]*Socket)
var socketCacheLock sync.Mutex

func SocketCount() int {
	return len(socketCache)
}

func SocketCache(id string) *Socket {
	return socketCache[id]
}

func GetSockets(username string) SocketSlice {
	sockets := SocketSlice{}
	socketCacheLock.Lock()
	for _, socket := range socketCache {
		if socket.Session == nil {
		} else if username == socket.Username {
			sockets = append(sockets, socket)
		}
	}
	socketCacheLock.Unlock()
	return sockets
}

func StoreSocket(socket *Socket) {
	socketCacheLock.Lock()
	socketCache[socket.Name()] = socket
	socketCacheLock.Unlock()
}

func RemoveSocket(socket *Socket) {
	socketCacheLock.Lock()
	delete(socketCache, socket.Name())
	socketCacheLock.Unlock()
}
