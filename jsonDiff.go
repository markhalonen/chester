package main

import (
	"encoding/json"
	"fmt"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getJSONObject(s []byte) map[string]interface{} {
	var f interface{}
	err := json.Unmarshal(s, &f)
	if err != nil {
		fmt.Println(string(s), " is not valid JSON")
		return make(map[string]interface{})
	}
	m := f.(map[string]interface{})
	return m
}

func getDifferingPaths(obA, obB map[string]interface{}, currentPath []interface{}) [][]interface{} {
	paths := make([][]interface{}, 0)
	bothKeys := make(map[string]bool) // Really a set...
	for k := range obA {
		bothKeys[k] = true
	}
	for k := range obB {
		bothKeys[k] = true
	}

	for k := range bothKeys {
		obAVal, aValPresent := obA[k]
		obBVal, bValPresent := obB[k]

		if !aValPresent || !bValPresent {
			paths = append(paths, append(currentPath, k))
			continue
		}

		switch obAVal.(type) {
		case map[string]interface{}:
			switch obBVal.(type) {
			case map[string]interface{}:
				// they are both objects, recur.
				obAValMap := obAVal.(map[string]interface{})
				obBValMap := obBVal.(map[string]interface{})
				paths = append(paths, getDifferingPaths(obAValMap, obBValMap, append(currentPath, k))...)
			default:
				// types are not equal.
				paths = append(paths, append(currentPath, k))
			}
		case []interface{}:
			switch obBVal.(type) {
			case []interface{}:
				ObASlice := obAVal.([]interface{})
				ObBSlice := obBVal.([]interface{})
				for i := 0; i < min(len(ObASlice), len(ObBSlice)); i++ {

				}
			default:
				// types are not equal.
				paths = append(paths, append(currentPath, k))
			}
		case string:
			if obAVal != obBVal {
				paths = append(paths, append(currentPath, k))
			}
		case float64:
			if obAVal != obBVal {
				paths = append(paths, append(currentPath, k))
			}
		default:
			fmt.Println(obAVal, "is of a type I don't know how to handle")
		}

	}

	return paths
}

func printJSONObject(ob map[string]interface{}) {
	for k, v := range ob {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case float64:
			fmt.Println(k, "is float64", vv)
		case []interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		case map[string]interface{}:
			m := v.(map[string]interface{})
			printJSONObject(m)
		default:
			fmt.Println(v, "is of a type I don't know how to handle")
		}
	}
}

func getIgnores(a, b []byte) [][]interface{} {
	aMap := getJSONObject(a)
	bMap := getJSONObject(b)

	diffs := getDifferingPaths(aMap, bMap, make([]interface{}, 0))
	return diffs
}
