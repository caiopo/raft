package main

import (
	"io/ioutil"
	"net/http"
	"sync"

	raft "github.com/caiopo/pontoon"
)

const RESPONSE_OK = 200

type SMEntry struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Body []byte `json:"body"`
}

type StateMachine struct {
	Log     []*SMEntry
	LogLock sync.RWMutex

	LastEntryCompleted int
	LCELock            sync.Mutex

	Answers     []map[string]string
	AnswersLock sync.RWMutex

	Replicas     []string
	ReplicasLock sync.RWMutex

	started bool
}

func (s *StateMachine) Apply(cr *raft.CommandRequest) error {
	if !s.started {
		s.start()
	}

	entry := &SMEntry{ID: cr.ID, Name: cr.Name, Body: cr.Body}

	s.AnswersLock.Lock()
	s.Answers = append(s.Answers, make(map[string]string))
	s.AnswersLock.Unlock()

	s.LogLock.Lock()
	s.Log = append(s.Log, entry)
	s.LogLock.Unlock()

	return nil
}

func (s *StateMachine) AddReplica(ip string) {

	s.ReplicasLock.Lock()
	s.Replicas = append(s.Replicas, ip)

	if s.started {
		go s.replica(len(s.Replicas) - 1)
	}

	s.ReplicasLock.Unlock()
}

func (s *StateMachine) replica(replicaNumber int) {

	s.ReplicasLock.RLock()
	myReplica := s.Replicas[replicaNumber]
	s.ReplicasLock.RUnlock()

	curIndex := 0

	for {
		// check if there is a new request in the log
		s.LogLock.RLock()
		next := len(s.Log) > curIndex
		s.LogLock.RUnlock()

		// if yes, send it to the replica
		if !next {
			continue
		}

		s.LogLock.RLock()
		curRequest := s.Log[curIndex]
		s.LogLock.RUnlock()

		body, err := sendTo(myReplica, curRequest)

		if err != nil {
			continue
		}

		// if the response is ok, save it and go for the next request
		s.AnswersLock.Lock()
		s.Answers[curIndex][myReplica] = body
		answerCount := len(s.Answers[curIndex])
		s.AnswersLock.Unlock()

		s.ReplicasLock.RLock()
		majority := (len(s.Replicas) / 2) + 1
		s.ReplicasLock.RUnlock()

		s.LCELock.Lock()
		if answerCount >= majority && s.LastEntryCompleted < curIndex {
			s.LastEntryCompleted = curIndex
		}
		s.LCELock.Unlock()

		curIndex++

	}
}

func (s *StateMachine) start() {
	s.ReplicasLock.RLock()
	for i, _ := range s.Replicas {
		go s.replica(i)
	}

	s.started = true
	s.ReplicasLock.RUnlock()
}

func sendTo(ip string, entry *SMEntry) (string, error) {

	response, err := http.Get("http://" + ip + "/" + string(entry.Body))

	if err != nil || response.StatusCode != RESPONSE_OK {
		return "", err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}
