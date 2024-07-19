package discovery

import (
	"context"
	"etcd-client/etcd"
	"fmt"
	"log"
	"sync"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Service struct {
	Name     string
	IP       string
	Port     string
	Protocol string
}

func ServiceRegister(s *Service) {
	var grantLesase bool
	var leaseID clientv3.LeaseID
	ctx := context.Background()
	// 租约
	// 查看是否有租约
	getRes, err := etcd.CLI.Get(ctx, s.Name, clientv3.WithCountOnly())
	if err != nil {
		log.Fatalln(err)
	}
	// 判断key是否存在
	if getRes.Count == 0 {
		grantLesase = true
	}
	// 租约声明
	if grantLesase {
		// 设置租约单次时间是10秒
		// 不进行续费的话只会存在10秒
		leaseRes, err := etcd.CLI.Grant(ctx, 10)
		if err != nil {
			log.Fatalln(err)
		}
		leaseID = leaseRes.ID
	}

	// 事务和数据操作
	kv := clientv3.NewKV(etcd.CLI)
	txn := kv.Txn(ctx)
	// 这一步是检查键存在库中的reversion，如果是0那就说明没有创建
	_, err = txn.If(clientv3.Compare(clientv3.CreateRevision(s.Name), "=", 0)).
		// 成立执行
		Then(
			// 没有创建或者数据过时的话指定租约进行创建
			clientv3.OpPut(s.Name, s.Name, clientv3.WithLease(leaseID)),
			clientv3.OpPut(s.Name+".ip", s.IP, clientv3.WithLease(leaseID)),
			clientv3.OpPut(s.Name+".port", s.Port, clientv3.WithLease(leaseID)),
			clientv3.OpPut(s.Name+".protocol", s.Protocol, clientv3.WithLease(leaseID)),
		).
		// 否则
		Else(
			// 创建过还没过期的话那就进行更新租约进行out
			// WithIgnoreLease参数的意义是更新有租约的数据，如果不加这个参数则新put的数据将不带有租约
			clientv3.OpPut(s.Name, s.Name, clientv3.WithIgnoreLease()),
			clientv3.OpPut(s.Name+".ip", s.IP, clientv3.WithIgnoreLease()),
			clientv3.OpPut(s.Name+".prot", s.Port, clientv3.WithIgnoreLease()),
			clientv3.OpPut(s.Name+".protocol", s.Protocol, clientv3.WithIgnoreLease()),
		).
		// 事务提交
		Commit()
	if err != nil {
		log.Fatalln(err)
	}
	// 更新好数据之后进行续约
	if grantLesase {
		leaseKeepalive, err := etcd.CLI.KeepAlive(ctx, leaseID)
		if err != nil {
			log.Fatalln(err)
		}
		for lease := range leaseKeepalive {
			fmt.Printf("leaseID:%x,ttl:%d\n", lease.ID, lease.TTL)
		}
	}
}

type Services struct {
	services map[string]*Service
	sync.RWMutex
}

var myService = &Services{
	services: map[string]*Service{
		"Hello.Greeter": &Service{},
	},
}

func ServiceDiscovery(svcName string) *Service {
	var s *Service = nil
	myService.Lock()
	s, _ = myService.services[svcName]
	myService.Unlock()
	return s
}

// 监听服务
func WatchServiceName(svcName string) {
	getRes, err := etcd.CLI.Get(context.Background(), svcName, clientv3.WithPrefix())
	if err != nil {
		log.Fatalln(err)
	}
	if getRes.Count > 0 {
		mp := sliceToMap(getRes.Kvs)
		s := &Service{}
		if kv, ok := mp[svcName]; ok {
			s.Name = string(kv.Value)
		}
		if kv, ok := mp[svcName+".ip"]; ok {
			s.IP = string(kv.Value)
		}
		if kv, ok := mp[svcName+".port"]; ok {
			s.Port = string(kv.Value)
		}
		if kv, ok := mp[svcName+".protocol"]; ok {
			s.Protocol = string(kv.Value)
		}
		myService.Lock()
		myService.services[svcName] = s
		myService.Unlock()
	}

	rch := etcd.CLI.Watch(context.Background(), svcName, clientv3.WithPrefix())
	for wres := range rch {
		for _, ev := range wres.Events {
			if ev.Type == clientv3.EventTypeDelete {
				myService.Lock()
				// delete(myService.services, svcName)
				myService.services[svcName] = &Service{}
				myService.Unlock()
			}
			if ev.Type == clientv3.EventTypePut {
				myService.Lock()
				if _, ok := myService.services[svcName]; !ok {
					myService.services[svcName] = &Service{}
				}
				switch string(ev.Kv.Key) {
				case svcName:
					myService.services[svcName].Name = string(ev.Kv.Value)
				case svcName + ".ip":
					myService.services[svcName].IP = string(ev.Kv.Value)
				case svcName + ".port":
					myService.services[svcName].Port = string(ev.Kv.Value)
				case svcName + ".protocol":
					myService.services[svcName].Protocol = string(ev.Kv.Value)
				}
				myService.Unlock()
			}
		}
	}
}

func sliceToMap(list []*mvccpb.KeyValue) map[string]*mvccpb.KeyValue {
	mp := make(map[string]*mvccpb.KeyValue)
	for _, item := range list {
		mp[string(item.Key)] = item
	}
	return mp
}
