package xorm

import (
	"bufio"
	"regexp"
	"strings"
)

type XSqlMap struct {
	sqlMapRootDir string
	extension     string
}

func XSql(directory, extension string) *XSqlMap {
	return &XSqlMap{
		sqlMapRootDir: directory,
		extension:     extension,
	}
}

func (sqlMap *XSqlMap) RootDir() string {
	return sqlMap.sqlMapRootDir
}

func (sqlMap *XSqlMap) Extension() string {
	return sqlMap.extension
}

type Scanner struct {
	line    string
	queries map[string]string
	current string
}

type stateFn func(*Scanner) stateFn

func getTag(line string) string {
	re := regexp.MustCompile("^\\s*--\\s*id:\\s*(\\S+)")
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		return ""
	}
	return matches[1]
}

func initialState(s *Scanner) stateFn {
	if tag := getTag(s.line); len(tag) > 0 {
		s.current = tag
		return queryState
	}
	return initialState
}

func queryState(s *Scanner) stateFn {
	if tag := getTag(s.line); len(tag) > 0 {
		s.current = tag
	} else {
		s.appendQueryLine()
	}
	return queryState
}

func (s *Scanner) appendQueryLine() {
	current := s.queries[s.current]
	line := strings.Trim(s.line, " \t")
	if len(line) == 0 {
		return
	}

	if len(current) > 0 {
		current = current + "\n"
	}

	current = current + line
	s.queries[s.current] = current
}

func (s *Scanner) Run(io *bufio.Scanner) map[string]string {
	s.queries = make(map[string]string)

	for state := initialState; io.Scan(); {
		s.line = io.Text()
		state = state(s)
	}

	return s.queries
}
