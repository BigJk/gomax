package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/boltdb/bolt"
)

type config struct {
	Mars    marsSettings `json:"mars"`
	EocTest bool         `json:"eocTest"`
	Suite   string       `json:"suite"`
	Phases  []phase      `json:"phases"`
	Running bool         `json:"-"`
}

func (c *config) decode(b []byte) {
	buffer := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buffer)
	dec.Decode(&c)
}

var currentConfig config
var currentWarrior warrior
var flatOperatorRegex *regexp.Regexp
var deepOperatorRegex *regexp.Regexp
var calcChan chan []int
var workerWaitGroup sync.WaitGroup

func main() {
	rand.Seed(time.Now().Unix())

	gob.Register(deepOperator{})
	gob.Register(flatOperator{})

	flatOperatorRegex = regexp.MustCompile("!\\((\\d*)-(\\d*)\\)")
	deepOperatorRegex = regexp.MustCompile("!{(\\d*)-(\\d*)}")

	if len(os.Args) == 1 {
		fmt.Println("gomax commands:")
		fmt.Println()
		fmt.Println("gomax create OUTPUT CONFIG WARRIOR")
		fmt.Println("Builds the database for optimizing.")
		fmt.Println("(eg. 'gomax create warrior.db config.json warrior.red')")
		fmt.Println("     OUTPUT  - Output file. (eg. 'warrior.db')")
		fmt.Println("     CONFIG  - Config file. (eg. 'config.json')")
		fmt.Println("     WARRIOR - Warrior file. (eg. 'warrior.red')")
		fmt.Println()
		fmt.Println("gomax run FILE WORKER")
		fmt.Println("Starts optimizing.")
		fmt.Println("(eg. 'gomax run warrior.db 6')")
		fmt.Println("     FILE    - Path to via 'create' generated database.")
		fmt.Println("     WORKER  - Count of worker / threads")
		fmt.Println()
		fmt.Println("gomax config")
		fmt.Println("Creates a sample config.")
		fmt.Println("(eg. 'gomax config')")
		return
	}

	command := os.Args[1]

	if command == "create" {
		createDatabse()
	} else if command == "run" {
		run()
	} else if command == "config" {
		createConfig()
	}

}

func createConfig() {
	c := config{marsSettings{8000, 80000, 8000, 200, 4000}, true, "suite_folder", make([]phase, 0), false}
	c.Phases = append(c.Phases, phase{"Non-Starter Selection", 25, 50, 85, false, "path_to_warrior", make([]string, 0), make([]string, 0), 0, 0, 0})
	c.Phases = append(c.Phases, phase{"Preselection", 25, 25, 90, false, "", []string{"scn", "clr"}, make([]string, 0), 0, 0, 0})
	c.Phases = append(c.Phases, phase{"Grand Prix", 25, 25, 100, false, "", []string{"cds", "clr", "pap", "pwi", "pws", "sbi", "scn", "stn"}, make([]string, 0), 0, 0, 0})
	b, _ := json.MarshalIndent(c, "", "	")
	ioutil.WriteFile("sampleconfig.json", b, 0777)
	fmt.Println("Sample config 'sampleconfig.json' generated!")
}

func run() {
	dbOutput := os.Args[2]
	workerCount, _ := strconv.Atoi(os.Args[3])

	db, _ = bolt.Open(dbOutput, 0600, nil)

	var states []int

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("settings"))
		currentConfig.decode(b.Get([]byte("config")))
		currentWarrior.decode(b.Get([]byte("warrior")))
		states = decodeInts(b.Get([]byte("states")))
		return nil
	})

	changeChannel = make(chan change, 100)
	calcChan = make(chan []int)

	log.Println("Starting to optimize...")
	log.Println("Current states", states)

	currentConfig.Running = true

	loadTop()

	initMars()

	go boltWorker()
	pushChange("settings", "closed", encodeInterface(false))

	workerWaitGroup.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go worker()
	}

	go calcAll(states)

	initRest()
}

func createDatabse() {
	dbOutput := os.Args[2]
	configPath := os.Args[3]
	warriorPath := os.Args[4]

	db, _ = bolt.Open(dbOutput, 0600, nil)
	defer db.Close()

	b, _ := ioutil.ReadFile(configPath)
	json.Unmarshal(b, &currentConfig)

	initMars()

	for i := 0; i < len(currentConfig.Phases); i++ {
		currentConfig.Phases[i].loadOponents(currentConfig.Suite)
	}

	bytes, _ := ioutil.ReadFile(warriorPath)
	w := parseWarrior(string(bytes))

	db.Update(func(tx *bolt.Tx) error {
		settings, _ := tx.CreateBucket([]byte("settings"))
		settings.Put([]byte("config"), encodeInterface(currentConfig))
		settings.Put([]byte("warriorcode"), bytes)
		settings.Put([]byte("warrior"), encodeInterface(w))
		settings.Put([]byte("states"), encodeInterface(make([]int, w.operatorSize())))
		settings.Put([]byte("closed"), encodeInterface(false))

		for i := 0; i < len(currentConfig.Phases); i++ {
			tx.CreateBucket([]byte("phase" + strconv.Itoa(i)))
		}

		return nil
	})

	fmt.Println("Successfully created GoMax Database!")
	fmt.Println("Use 'gomax run " + dbOutput + " 6' to start optimizing.")
}

func calcAll(states []int) {
	ft := reflect.ValueOf(flatOperator{}).String()
	//dt := reflect.ValueOf(deepOperator{})

	var types []int

	for i := 0; i < currentWarrior.operatorSize(); i++ {
		if reflect.ValueOf(currentWarrior.getOperator(i)).String() == ft {
			types = append(types, 0)
		} else {
			types = append(types, 1)
		}
	}

	for currentConfig.Running {
		for i := 0; i < currentWarrior.operatorSize(); i++ {
			o := currentWarrior.getOperator(i)
			if o.isMax(states[i]) {
				if types[i] == 0 {
					states[i] = rand.Intn(o.max()-o.min()) + o.min()
				} else {
					states[i] = 0
					if i+1 < currentWarrior.operatorSize() {
						if types[i+1] == 1 {
							states[i+1]++
						}
					}
				}
			}
		}
		copyStats := make([]int, len(states))
		copy(copyStats, states)
		copyStats = append([]int(nil), states...)
		calcChan <- copyStats
		pushChange("settings", "states", encodeInterface(copyStats))
		states[0]++
	}
	close(calcChan)
	log.Println("Closing worker...")
	workerWaitGroup.Wait()

	time.Sleep(time.Second * 3)

	for i := 0; i < len(currentConfig.Phases); i++ {
		pushChange("phase"+strconv.Itoa(i), "top", encodeInterface(phaseResultSorter(phaseTop[i])))
	}

	pushChange("settings", "closed", encodeInterface(true))

	log.Println("Closing DB...")
	for len(changeChannel) > 0 {
		time.Sleep(time.Second)
	}
	db.Close()

	log.Println("Everything closed...")

	os.Exit(0)
}

func worker() {
	for {
		s, ok := <-calcChan
		if ok {
			currentWarrior.push(s)
		} else {
			workerWaitGroup.Done()
			return
		}
	}
}
