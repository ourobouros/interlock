// INTERLOCK | https://github.com/inversepath/interlock
// Copyright (c) 2015 Inverse Path S.r.l.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type jsonObject map[string]interface{}

var censorPattern = regexp.MustCompile("password")

func parseRequest(r *http.Request) (j jsonObject, err error) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return
	}

	if conf.Debug {
		if conf.testMode || !censorPattern.Match(body) {
			log.Printf("%s", body)
		}
	}

	d := json.NewDecoder(strings.NewReader(string(body[:])))
	d.UseNumber()

	err = d.Decode(&j)

	if err != nil {
		return
	}

	return
}

func (j jsonObject) String() (s string) {
	b, err := json.Marshal(j)

	if err != nil {
		log.Print(err)
		return
	}

	s = string(b)

	return
}

func validateRequest(req jsonObject, reqAttrs []string) error {
	for i := 0; i < len(reqAttrs); i++ {
		ok := false

		args := strings.Split(reqAttrs[i], ":")

		if len(args) != 2 {
			return errors.New("unknown validation argument")
		}

		key := args[0]
		kind := args[1]

		if _, ok = req[key]; !ok {
			return fmt.Errorf("missing attribute %s", key)
		}

		switch kind {
		case "s":
			_, ok = req[key].(string)
		case "b":
			_, ok = req[key].(bool)
		case "n":
			_, ok = req[key].(json.Number)
		case "a":
			_, ok = req[key].([]interface{})
		default:
			return errors.New("unknown validation kind")
		}

		if !ok {
			return fmt.Errorf("attribute %s is not a %s", key, kind)
		}
	}

	return nil
}
