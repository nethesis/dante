/*
 * Copyright (C) 2019 Nethesis S.r.l.
 * http://www.nethesis.it - info@nethesis.it
 *
 * This file is part of Dante project.
 *
 * Dante is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * Dante is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dante.  If not, see COPYING.
 */

package utils

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type HttpError struct {
	Code        int
	ErrorString string
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("%d - %s", e.Code, e.ErrorString)
}

func ContainsString(stringSlice []string, searchString string) bool {
	for _, value := range stringSlice {
		if value == searchString {
			return true
		}
	}
	return false
}

func ReadJson(filePath string) (map[string]interface{}, error) {
	var mapData map[string]interface{}

	jsonFile, err := os.Open(filePath)
	defer jsonFile.Close()
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(bytes, &mapData)
	if err != nil {
		return nil, err
	}
	return mapData, nil
}

func ReadJsonIgnoreOpenError(filePath string) (map[string]interface{}, bool, error) {
	var mapData map[string]interface{}

	jsonFile, err := os.Open(filePath)
	defer jsonFile.Close()
	if err != nil {
		return nil, true, err
	}
	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, false, err
	}
	json.Unmarshal(bytes, &mapData)
	if err != nil {
		return nil, false, err
	}
	return mapData, false, nil
}

func ReverseSliceFloat(slice []float64) []float64 {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

func ReverseSliceString(slice []string) []string {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

func GetDateStringFromWidgetFilePath(filePath string) string {
	filePathTokens := strings.Split(filePath, "/")
	day := filePathTokens[len(filePathTokens)-2]
	month := filePathTokens[len(filePathTokens)-3]
	year := filePathTokens[len(filePathTokens)-4]
	return year + "-" + month + "-" + day
}

func Anonymize(value string) string {
	h := md5.New()
	h.Write([]byte(value))
	bs := h.Sum(nil)
	anonymizedString := fmt.Sprintf("%x", bs)
	return anonymizedString[:8]
}

func IndexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found
}

func IsDashOrDot(r rune) bool {
	switch r {
	case '-', '.':
		return true
	}
	return false
}

func moveStringAtIndex(val string, indexWhereToInsert int, slice []string) {
	indexToRemove := IndexOf(val, slice)
	slice = append(slice[:indexToRemove], slice[indexToRemove+1:]...)
	newSlice := make([]string, indexWhereToInsert+1)
	copy(newSlice, slice[:indexWhereToInsert])
	newSlice[indexWhereToInsert] = val
	slice = append(newSlice, slice[indexWhereToInsert:]...)
}
