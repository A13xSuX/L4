package main

import (
	"errors"
	"strconv"
	"strings"
	"sync"
)

func parseFields(fieldsStr string) ([]int, error) {
	result := make([]int, 0)
	if fieldsStr == "" { // check empty fields
		return result, errors.New("Fields is empty")
	}
	split := strings.Split(fieldsStr, ",")
	for _, field := range split {
		if strings.Contains(field, "-") { //range
			manyNum := strings.Split(field, "-")
			if len(manyNum) != 2 {
				return result, errors.New("Fields is invalid, len!= 2")
			}
			start, err := strconv.Atoi(strings.TrimSpace(manyNum[0]))
			if err != nil {
				return result, errors.New("Start is invalid")
			}
			end, err := strconv.Atoi(strings.TrimSpace(manyNum[len(manyNum)-1]))
			if err != nil {
				return result, errors.New("End is invalid")
			}
			if start > end {
				return result, errors.New("Start cannot be greater than End")
			}
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
		} else { //solo element
			num, err := strconv.Atoi(field)
			if err != nil {
				return result, err
			}
			result = append(result, num)
		}
	}
	return result, nil
}

func processLine(line string, fields []int, delimiter string, separated bool) string {
	if !strings.Contains(line, delimiter) && separated {
		return ""
	}

	column := strings.Split(line, delimiter)
	var resultFields []string

	for _, fieldNum := range fields { //fieldNum from user begin from 1
		if fieldNum-1 < len(column) && fieldNum > 0 {
			resultFields = append(resultFields, column[fieldNum-1])
		} //ignore fields for out
	}
	return strings.Join(resultFields, delimiter)
}
func processLinesConcurrent(lines []string, fields []int, delimiter string, separated bool) []string {
	jobs := make(chan lineJob)
	result := make(chan lineResult, len(lines))

	workers := 4
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				out := processLine(job.text, fields, delimiter, separated)
				result <- lineResult{
					index: job.index,
					text:  out,
				}
			}
		}()
	}

	go func() {
		for i, line := range lines {
			jobs <- lineJob{
				index: i,
				text:  line,
			}
		}
		close(jobs)
		wg.Wait()
		close(result)
	}()

	temp := make([]string, len(lines))
	for res := range result {
		temp[res.index] = res.text
	}

	final := make([]string, 0, len(lines))
	for _, line := range temp {
		if line == "" {
			continue
		}
		final = append(final, line)
	}
	return final
}
