package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"golang.org/x/net/websocket"

	_ "net/http/pprof"

	"github.com/boltdb/bolt"
	"github.com/julienschmidt/httprouter"
)

type phaseInformation struct {
	Top       []phaseResult
	Total     int
	Passed    int
	Failed    int
	Bestscore float32
}

type indexPhase struct {
	Phase        phase
	PhaseResults []phaseResult
}

type indexTemplate struct {
	Config config
	Phases []indexPhase
}

var phaseTop [][]phaseResult

func loadTop() {
	phaseTop = make([][]phaseResult, len(currentConfig.Phases))
	c := false

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("settings"))
		c = decodeBool(b.Get([]byte("closed")))
		return nil
	})

	if c {
		log.Println("Loading Phases from cache...")
	} else {
		log.Println("Program wasn't correctly closed... Rebuilding phases...")
	}

	for i := 0; i < len(currentConfig.Phases); i++ {
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("phase" + strconv.Itoa(i)))

			if c {
				p := phaseResultSorter{}
				p.decode(b.Get([]byte("top")))
				phaseTop[i] = []phaseResult(p)
			} else {
				c := b.Cursor()

				for k, v := c.First(); k != nil; k, v = c.Next() {
					p := phaseResult{}
					p.decode(v)
					phaseTop[i] = append(phaseTop[i], p)
					currentConfig.Phases[i].Total++
					if p.Passed {
						currentConfig.Phases[i].Passed++
					}
				}
			}

			return nil
		})

		if !c {
			sort.Sort(phaseResultSorter(phaseTop[i]))
			if len(phaseTop[i]) >= currentConfig.Phases[i].Top {
				phaseTop[i] = phaseTop[i][:currentConfig.Phases[i].Top]
			}
		}

	}
}

func pushWarriorToTop(p int, r phaseResult) {
	if len(phaseTop[p]) >= currentConfig.Phases[p].Top {
		if r.Score > phaseTop[p][currentConfig.Phases[p].Top-1].Score {
			phaseTop[p][currentConfig.Phases[p].Top-1] = r
			sort.Sort(phaseResultSorter(phaseTop[p]))
		}
	} else {
		phaseTop[p] = append(phaseTop[p], r)
		sort.Sort(phaseResultSorter(phaseTop[p]))
	}
}

func constructRest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s := strings.Split(ps.ByName("ints"), ",")
	is := make([]int, len(s))
	for i := 0; i < len(s); i++ {
		c, _ := strconv.Atoi(s[i])
		is[i] = c
	}
	w.Write([]byte(currentWarrior.construct(is)))
}

func configRest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	b, _ := json.MarshalIndent(currentConfig, "", "	")
	w.Write(b)
}

func phasesCount() int {
	return len(currentConfig.Phases)
}

func indexRest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	funcs := template.FuncMap{"pcount": phasesCount}
	t, _ := template.New("Index").Funcs(funcs).ParseFiles("views/index.html")
	index := indexTemplate{currentConfig, make([]indexPhase, len(currentConfig.Phases))}
	for i := 0; i < len(currentConfig.Phases); i++ {
		index.Phases[i] = indexPhase{currentConfig.Phases[i], phaseTop[i]}
	}
	t.ExecuteTemplate(w, "Index", index)
}

func stopRest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	currentConfig.Running = false
	w.Write([]byte("Ok"))
}

func initRest() {
	router := httprouter.New()

	router.GET("/", indexRest)
	router.GET("/api/stop", stopRest)
	router.GET("/api/config", configRest)
	router.GET("/api/warrior/:ints", constructRest)

	router.Handler("GET", "/ws", websocket.Handler(wsHandler))

	router.NotFound = http.FileServer(http.Dir("public"))

	go wsWorker()

	log.Fatal(http.ListenAndServe(":8080", router))
}
