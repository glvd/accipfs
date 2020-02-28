package service

import (
	"fmt"
	"github.com/glvd/accipfs/aws"
	"github.com/glvd/accipfs/config"
	"strings"
	"sync"
)

// Service ...
type Service struct {
	cfg        *config.Config
	serveMutex sync.RWMutex
	serve      []Node
	i          *nodeClientIPFS
	e          *nodeClientETH
	nodes      map[string][]byte
}

// New ...
func New(config config.Config) (s *Service, e error) {
	s = &Service{
		cfg:   &config,
		nodes: make(map[string][]byte),
	}
	s.i, e = newNodeIPFS(config)
	if e != nil {
		return nil, e
	}
	s.e, e = newETH(config)
	if e != nil {
		return nil, e
	}
	return s, e
}

// RegisterServer ...
func (s *Service) RegisterServer(node Node) {
	s.serveMutex.Lock()
	defer s.serveMutex.Unlock()
	s.serve = append(s.serve, node)
}

// Run ...
func (s *Service) Run() {
	s.serveMutex.RLock()
	defer s.serveMutex.RUnlock()
	for _, s := range s.serve {
		s.Start()
	}
}

func (s *Service) syncDNS() {
	//defer fmt.Println("<更新网关数据完成...>")
	var records []string
	// build node records
	for node := range s.nodes {
		if !strings.Contains(node, "enode") {
			continue
		}
		// get ip address
		uri := strings.Split(node, "@")[1]
		ip := strings.Split(uri, ":")[0]
		records = append(records, ip)
	}

	if len(records) == 0 {
		return
	}
	fmt.Println("<正在更新网关数据...>", records)

	dnsService := aws.NewRoute53(s.cfg)

	// get remote dns record
	var remoteIPs []string
	remoteRecordSets, err := dnsService.GetRecordSets()
	if err != nil {
		fmt.Println("<访问远端网关失败> ", err.Error())
		return
	}
	if len(remoteRecordSets) != 0 {
		for _, recordSet := range remoteRecordSets {
			remoteIPs = append(remoteIPs, *recordSet.ResourceRecords[0].Value)
		}
	}
	//// add new record
	//ipAdd := removeDuplicateElement(general.DiffStrArray(records, remoteIPs))
	//fmt.Println("[resource to be added]", ipAdd)
	//setsAdd := dnsService.BuildMultiValueRecordSets(ipAdd)
	//if len(setsAdd) > 0 {
	//	res, err := dnsService.ChangeSets(setsAdd, "UPSERT")
	//	if err != nil {
	//		fmt.Println("[add resource record fail]", err.Error())
	//	} else {
	//		fmt.Println("[add resource record success]", res.String())
	//	}
	//}
	//
	//// delete record out of date
	//failedSets := dnsService.FilterFailedRecords(remoteRecordSets)
	//fmt.Println("[resource to be deleted]", len(failedSets))
	//if len(failedSets) > 0 {
	//	res, err := dnsService.ChangeSets(failedSets, "DELETE")
	//	if err != nil {
	//		fmt.Println("[delete resource record fail]", err.Error())
	//	} else {
	//		fmt.Println("[delete resource record success]", res.String())
	//	}
	//}

	return
}
