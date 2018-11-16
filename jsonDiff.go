package main

import (
	"encoding/json"
	"strconv"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func getJSONObject(s []byte) (map[string]interface{}, error) {
	var f interface{}
	err := json.Unmarshal(s, &f)
	if err != nil {
		return nil, err
	}
	m := f.(map[string]interface{})
	return m, err
}

func compareObs(obAVal, obBVal interface{}, currentPath []interface{}, k interface{}) [][]interface{} {
	paths := make([][]interface{}, 0)
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
				aVal := ObASlice[i]
				bVal := ObBSlice[i]
				paths = append(paths, compareObs(aVal, bVal, append(currentPath, k), i)...)
			}
			for i := min(len(ObASlice), len(ObBSlice)); i < max(len(ObASlice), len(ObBSlice)); i++ {

				paths = append(paths, append(append(currentPath, k), i))
			}
		default:
			// types are not equal.
			paths = append(paths, append(currentPath, k))
		}
	default:
		if obAVal != obBVal {
			paths = append(paths, append(currentPath, k))
		}
	}
	return paths
}

func getDifferingPaths(obA, obB map[string]interface{}, currentPath []interface{}) [][]interface{} {
	paths := make([][]interface{}, 0)
	bothKeys := []string{}

	for k := range obA {
		bothKeys = append(bothKeys, k)
	}

	for k := range obB {
		inA := false

		// Don't add duplicates
		for k2 := range obA {
			if k == k2 {
				inA = true
			}
		}

		if !inA {
			bothKeys = append(bothKeys, k)
		}
	}

	for _, k := range bothKeys {
		obAVal, aValPresent := obA[k]
		obBVal, bValPresent := obB[k]

		if !aValPresent || !bValPresent {
			paths = append(paths, append(currentPath, k))
			continue
		}

		paths = append(paths, compareObs(obAVal, obBVal, currentPath, k)...)

	}

	return paths
}

func getIgnores(a, b string) ([][]interface{}, error) {
	aMap, err1 := getJSONObject([]byte(a))
	if err1 != nil {
		return nil, err1
	}
	bMap, err2 := getJSONObject([]byte(b))
	if err2 != nil {
		return nil, err2
	}
	diffs := getDifferingPaths(aMap, bMap, make([]interface{}, 0))
	return diffs, nil
}

func getJSONPath(path []interface{}) string {
	if len(path) == 0 {
		return ""
	}
	result := ""
	for _, v := range path {
		switch v.(type) {
		case int:
			result = result + "[" + strconv.Itoa(v.(int)) + "]"
		case string:
			result = result + "[\"" + v.(string) + "\"]"
		}
	}
	return result
}

func getMessage(ignores [][]interface{}) string {
	if len(ignores) == 0 {
		return ""
	}
	result := "Chester calculated that the following JSON paths differ:\n"
	for _, ignore := range ignores {
		result = result + "- " + getJSONPath(ignore) + "\n"
	}
	return result
}

func jsonDiffMessage(a, b string) string {
	ignores, err := getIgnores(a, b)
	if err != nil {
		return ""
	}
	return getMessage(ignores)
}
