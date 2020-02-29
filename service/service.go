package service

import (
	"fmt"
	"github.com/glvd/accipfs/aws"
	"github.com/glvd/accipfs/config"
	"strings"
	"sync"
)

const outputHead = "<Service>"

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
	fmt.Println(outputHead, "<正在更新网关数据...>", records)

	dnsService := aws.NewRoute53(s.cfg)

	// get remote dns record
	remoteIPs := make(map[string]bool)
	remoteRecordSets, err := dnsService.GetRecordSets()
	if err != nil {
		fmt.Println(outputHead, "<访问远端网关失败> ", err.Error())
		return
	}
	if len(remoteRecordSets) != 0 {
		for _, recordSet := range remoteRecordSets {
			remoteIPs[*recordSet.ResourceRecords[0].Value] = true
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

func removeDuplicateElement(addrs []string) []string {
	result := make([]string, 0, len(addrs))
	temp := map[string]struct{}{}
	for _, item := range addrs {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// DiffStrArray return the elements in `a` that aren't in `b`.
func DiffStrArray(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
