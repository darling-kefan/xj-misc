package main

import (
	"database/sql"
	"encoding/xml"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	//"reflect"
	"strings"
)

type Env struct {
	Db *sql.DB
}

type District struct {
	Did       int64
	Pid       int64
	Name      string
	Zipcode   int64
	Leveltype int64
}

func (env *Env) allDistricts() ([]*District, map[*District][]*District) {
	rows, err := env.Db.Query("select `did`, `pid`, `name`, `zipcode`, `leveltype` from districts")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	dts := make([]*District, 0)
	for rows.Next() {
		dt := new(District)
		err := rows.Scan(&dt.Did, &dt.Pid, &dt.Name, &dt.Zipcode, &dt.Leveltype)
		if err != nil {
			log.Fatal(err)
		}
		dts = append(dts, dt)
	}

	provinces := make([]*District, 0)
	dtmaps := make(map[*District][]*District, len(dts))
	for _, dt := range dts {
		if dt.Leveltype == 1 && dt.Pid == 100000 {
			provinces = append(provinces, dt)
		}
		cdts := make([]*District, 0)
		for _, dt2 := range dts {
			if dt2.Pid == dt.Did {
				cdts = append(cdts, dt2)
			}
		}
		if len(cdts) > 0 {
			dtmaps[dt] = cdts
		}
	}

	return provinces, dtmaps
}

type Data struct {
	Dicts []Dict `plist:"address"`
}

type Dict struct {
	Did     string `plist:"id"`
	Name    string `plist:"name"`
	Sub     []Dict `plist:"sub"`
	Zipcode string `plist:"zipcode"`
}

func main() {
	// mysql -uxjroot -pDgyE1ZHdMT5x -h 192.168.0.190 -P 13306 oauth
	db, err := sql.Open("mysql", "xjroot:DgyE1ZHdMT5x@tcp(192.168.0.190:13306)/oauth")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var env *Env = &Env{
		Db: db,
	}

	provinces, dtmaps := env.allDistricts()

	dicts := make([]Dict, 0)
	for _, pr := range provinces {
		province := Dict{
			Did:  pr.Did,
			Name: pr.Name,
		}

		cities := make([]Dict, 0)
		if _, ok := dtmaps[pr]; ok {
			for _, ci := range dtmaps[pr] {
				city := Dict{
					Did:  ci.Did,
					Name: ci.Name,
				}

				countries := make([]Dict, 0)
				if _, ok = dtmaps[ci]; ok {
					for _, co := range dtmaps[ci] {
						country := Dict{
							Did:  co.Did,
							Name: co.Name,
						}
						countries = append(countries, country)
					}
				}
				city.Zipcode = ci.Zipcode
				if len(countries) > 0 {
					city.Sub = countries
				}
				cities = append(cities, city)
			}
		}

		if len(cities) > 0 {
			province.Sub = cities
		}
		dicts = append(dicts, province)
	}

	data := &Data{dicts}

	plistdata, err := plist.MarshalIndent(data, XMLFormat, "\t")

	os.Stdout.Write(plistdata)
}
