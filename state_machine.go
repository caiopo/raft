package main

import (
	"net/http"
	"sync"
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

	Answers     []map[int]string
	AnswersLock sync.RWMutex

	Replicas     []string
	ReplicasLock sync.RWMutex

	starter sync.Once
}

func (s *StateMachine) Apply(cr *CommandRequest) error {
	s.starter(s.start)

	entry := &SMEntry{ID: cr.ID, Name: cr.Name, Body: cr.Body}

	s.AnswersLock.Lock()
	s.Answers = append(s.Answers, make(map[int]string))
	s.AnswersLock.Unlock()

	s.LogLock.Lock()
	s.Log = append(s.Log, entry)
	s.LogLock.Unlock()

	return nil
}

func (s *StateMachine) AddReplica(ip string) {

	s.ReplicasLock.Lock()
	s.Replicas = append(s.Replicas, ip)
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
		if next {
			s.LogLock.RLock()
			curRequest := s.Log[curIndex]
			s.LogLock.RUnlock()

			response, err := sendTo(myReplica, curRequest)

			if err != nil {
				continue
			}

			defer response.Body.Close()

			// if the response is ok, save it and go for the next request
			if response.StatusCode == RESPONSE_OK {
				body, err := ioutil.ReadAll(response.Body)

				if err != nil {
					continue
				}

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
	}
}

func (s *StateMachine) start() {

	for i, _ := range s.Replicas {
		go s.replica(i)
	}

}

func sendTo(ip string, entry *SMEntry) {
	http.Get("http://" + ip + string(entry.Body))
}
