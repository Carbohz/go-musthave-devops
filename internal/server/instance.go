package server

import (
	"encoding/json"
	"github.com/Carbohz/go-musthave-devops/internal/metrics"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"os"
	"time"
)

type Instance struct {
	Cfg Config
	is  *internalStorage // `is` is shortcut for internal storage
	dbs *dbStorage       // `dbs` is shortcut for database storage
}

func CreateInstance(cfg Config) Instance {
	var instance Instance
	instance.Cfg = cfg

	if cfg.Restore && cfg.StoreFile != "" {
		instance.LoadMetrics()
	}

	if cfg.StoreInterval > 0 && cfg.StoreFile != "" {
		go instance.DumpMetrics()
	}

	//// new function
	//handler.PassSecretKey(cfg.Key)

	return instance
}

func (instance Instance) RunInstance() {
	r := chi.NewRouter()
	r.Use(middleware.Compress(5))
	server.SetupRouters(r)
	server := &http.Server{
		Addr:    instance.Cfg.Address,
		Handler: r,
	}
	server.SetKeepAlivesEnabled(false)
	log.Printf("Listening on address %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func (instance Instance) BeforeShutDown() {
	if instance.Cfg.StoreFile != "" {
		log.Println("Dumping metrics and exiting")

		cfg := instance.Cfg

		flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC

		f, err := os.OpenFile(cfg.StoreFile, flag, 0644)
		if err != nil {
			log.Fatal("Can't open file for dumping: ", err)
		}
		defer f.Close()

		encoder := json.NewEncoder(f)

		if err := encoder.Encode(*instance.is); err != nil {
			log.Fatal("Can't encode server's metrics: ", err)
		}
	}
}

// insertion into storage
func (instance Instance) StoreGaugeMetric(name string, value float64) {
	if instance.Cfg.DBPath != "" {
		// working with db
	} else {
		// working with internal storage
		instance.is.StoreGaugeMetric(name, value)
	}
}

func (instance Instance) StoreCounterMetric(name string, value int64) {
	if instance.Cfg.DBPath != "" {
		// working with db
	} else {
		// working with internal storage
		instance.is.StoreCounterMetric(name, value)
	}
}

// storage lookup
func (instance Instance) FindGaugeMetric(name string) (float64, error) {
	if instance.Cfg.DBPath != "" {
		// working with db
		return -1.0, nil
	} else {
		// working with internal storage
		return instance.is.FindGaugeMetric(name)
	}
}

func (instance Instance) FindCounterMetric(name string) (int64, error) {
	if instance.Cfg.DBPath != "" {
		// working with db
		return -1, nil
	} else {
		// working with internal storage
		return instance.is.FindCounterMetric(name)
	}
}

// get data
func (instance Instance) GetGaugeMetrics() map[string]metrics.GaugeMetric {
	if instance.Cfg.DBPath != "" {
		// working with db
		return make(map[string]metrics.GaugeMetric)
	} else {
		// working with internal storage
		return instance.is.GetGaugeMetrics()
	}
}

func (instance Instance) GetCounterMetrics() map[string]metrics.CounterMetric {
	if instance.Cfg.DBPath != "" {
		// working with db
		return make(map[string]metrics.CounterMetric)
	} else {
		// working with internal storage
		return instance.is.GetCounterMetrics()
	}
}

// update data
func (instance Instance) UpdateGaugeMetrics(name string, value float64) {
	if instance.Cfg.DBPath != "" {
		// working with db

	} else {
		// working with internal storage
		instance.is.UpdateGaugeMetrics(name, value)
	}
}

func (instance Instance) UpdateCounterMetrics(name string, value int64) {
	if instance.Cfg.DBPath != "" {
		// working with db

	} else {
		// working with internal storage
		instance.is.UpdateCounterMetrics(name, value)
	}
}

// dump metrics (to file or to db)
func (instance Instance) DumpMetrics() {
	ticker := time.NewTicker(instance.Cfg.StoreInterval)
	for {
		<-ticker.C
		log.Printf("Dumping metrics to file %s", instance.Cfg.StoreFile)
		if instance.Cfg.DBPath != "" {
			// working with db

		} else {
			// working with internal storage
			instance.is.DumpMetricsToFile(instance)
		}
	}

	//if instance.Cfg.DBPath != "" {
	//	// working with db
	//
	//} else {
	//	// working with internal storage
	//	instance.is.DumpMetricsToFile(instance)
	//}
}

// load metrics (from file or from db)
func (instance Instance) LoadMetrics() {
	if instance.Cfg.DBPath != "" {
		// working with db

	} else {
		// working with internal storage
		instance.is.LoadMetricsFromFile(instance)
	}
}
