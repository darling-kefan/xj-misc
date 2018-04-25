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

type Plist struct {
	XMLName xml.Name   `xml:"plist"`
	Version string     `xml:"version,attr"`
	DictKey string     `xml:"dict>key"`
	Dicts   []Province `xml:"dict>array>dict"`
}

type Province struct {
	IdPart
	NamePart
	SubKey string `xml:"key_3,omitempty"`
	SubVal []City `xml:"array>dict,omitempty"`
}

type City struct {
	IdPart
	NamePart
	SubKey     string    `xml:"key_3,omitempty"`
	SubVal     []Country `xml:"array>dict,omitempty"`
	ZipcodeKey string    `xml:"key_4,omitempty"`
	ZipcodeVal int64     `xml:"string_4,omitempty"`
}

type KeyString []interface{}

type Country struct {
	IdPart
	NamePart
}

type IdPart struct {
	IdKey string `xml:"key_1"`
	IdVal int64  `xml:"string_1"`
}

type NamePart struct {
	NameKey string `xml:"key_2"`
	NameVal string `xml:"string_2"`
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

	dicts := make([]Province, 0)
	for _, pr := range provinces {
		province := Province{
			IdPart:   IdPart{"id", pr.Did},
			NamePart: NamePart{"name", pr.Name},
		}
		cities := make([]City, 0)
		if _, ok := dtmaps[pr]; ok {
			for _, ci := range dtmaps[pr] {
				city := City{
					IdPart:   IdPart{"id", ci.Did},
					NamePart: NamePart{"name", ci.Name},
				}
				countries := make([]Country, 0)
				if _, ok = dtmaps[ci]; ok {
					for _, co := range dtmaps[ci] {
						country := Country{
							IdPart:   IdPart{"id", co.Did},
							NamePart: NamePart{"name", co.Name},
						}
						countries = append(countries, country)
					}
				}
				city.ZipcodeKey = "zipcode"
				city.ZipcodeVal = ci.Zipcode
				if len(countries) > 0 {
					city.SubKey = "sub"
					city.SubVal = countries
				}
				cities = append(cities, city)
			}
		}
		if len(cities) > 0 {
			province.SubKey = "sub"
			province.SubVal = cities
		}

		dicts = append(dicts, province)
	}

	plist := &Plist{Version: "1.0", DictKey: "address", Dicts: dicts}

	output, err := xml.MarshalIndent(plist, "", "    ")
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	replace1 := `
                        </array>
                        <array>`
	replace2 := `
                <array></array>`
	r := strings.NewReplacer("key_1", "key", "key_2", "key", "key_3", "key", "key_4", "key", "string_1", "string", "string_2", "string", "string_4", "string", replace1, "", replace2, "")
	newOutput := r.Replace(string(output))
	header := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
`
	newOutput = header + newOutput

	os.Stdout.Write([]byte(newOutput))
}
