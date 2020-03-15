package service

import (
	"fmt"
	"github.com/glvd/accipfs/aws"
	"github.com/glvd/accipfs/config"
	"github.com/gocacher/badger-cache"
	"github.com/gocacher/cacher"
	"github.com/goextension/log"
	"github.com/robfig/cron/v3"
	"strings"
	"sync"
)

const outputHead = "<Service>"

// Service ...
type Service struct {
	cfg        *config.Config
	cache      cacher.Cacher
	cron       *cron.Cron
	serveMutex sync.RWMutex
	serve      []NodeServer
	i          *nodeClientIPFS
	e          *nodeClientETH
	nodes      map[string]bool
}

// New ...
func New(cfg config.Config) (s *Service, e error) {
	s = &Service{
		cfg:   &cfg,
		nodes: make(map[string]bool),
	}
	s.serve = append(s.serve, NewNodeServerIPFS(cfg), NewNodeServerETH(cfg))

	s.i, e = newNodeIPFS(cfg)
	if e != nil {
		return nil, e
	}
	s.e, e = newNodeETH(cfg)
	if e != nil {
		return nil, e
	}

	cache.DefaultCachePath = config.DataDirCache()
	s.cache = cache.New()

	s.cron = cron.New(cron.WithSeconds())
	return s, e
}

// RegisterServer ...
func (s *Service) RegisterServer(node NodeServer) {
	s.serve = append(s.serve, node)
}

// Run ...
func (s *Service) Run() {
	for _, serv := range s.serve {
		if err := serv.Start(); err != nil {
			panic(err)
		}
	}
	jobETH, err := s.cron.AddJob("0 * * * * *", s.e)
	if err != nil {
		panic(err)
	}
	fmt.Println(outputHead, "ETH", "run id", jobETH)

	jobIPFS, err := s.cron.AddJob("0 * * * * *", s.i)
	if err != nil {
		panic(err)
	}
	fmt.Println(outputHead, "IPFS", "run id", jobIPFS)

	s.cron.Run()
}

// Stop ...
func (s *Service) Stop() {
	for _, serv := range s.serve {
		if err := serv.Stop(); err != nil {
			log.Errorw("stop error", "tag", outputHead, "error", err)
		}
	}
}

func syncDNS(cfg config.Config, nodes map[string]bool) {
	//defer fmt.Println("<更新网关数据完成...>")
	var records []string
	// build serviceNode records
	for node := range nodes {
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

	dnsService := aws.NewRoute53(&cfg)

	// get remote dns record
	remoteIPs := make(map[string]bool)
	remoteRecordSets, err := dnsService.GetRecordSets()
	if err != nil {
		log.Infow("visit remote record failed", "tag", outputHead, "error", err.Error())
		return
	}
	if len(remoteRecordSets) != 0 {
		for _, recordSet := range remoteRecordSets {
			remoteIPs[*recordSet.ResourceRecords[0].Value] = true
		}
	}
	// add new record
	ipAdd := DiffStrArray(records, remoteIPs)
	setsAdd := dnsService.BuildMultiValueRecordSets(ipAdd)
	log.Infow("resource adding", "tag", outputHead, "list", ipAdd, "count", len(setsAdd))
	if len(setsAdd) > 0 {
		res, err := dnsService.ChangeSets(setsAdd, "UPSERT")
		if err != nil {
			log.Infow("add resource record fail", "tag", outputHead, "error", err)
		} else {
			log.Infow("add resource record success", "tag", outputHead, "error", "result", res.String())
		}
	}

	// delete record out of date
	failedSets := dnsService.FilterFailedRecords(remoteRecordSets)
	log.Infow("resource deleting", "tag", outputHead, "list", remoteRecordSets, "count", len(failedSets))
	if len(failedSets) > 0 {
		res, err := dnsService.ChangeSets(failedSets, "DELETE")
		if err != nil {
			log.Infow("delete resource record fail", "tag", outputHead, "error", err)
		} else {
			log.Infow("delete resource record success", "tag", outputHead, "error", "result", res.String())
		}
	}

	return
}

// DiffStrArray return the elements in `a` that aren't in `b`.
func DiffStrArray(a []string, b map[string]bool) []string {
	var diff []string
	for _, x := range a {
		if _, found := b[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
