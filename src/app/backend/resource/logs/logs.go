package logs

import "strings"

// LINE_INDEX_NOT_FOUND is returned if requested line could not be found
var LINE_INDEX_NOT_FOUND = -1

// Default number of lines that should be returned in case of invalid request
var DefaultDisplayNumLogLines = 100

// MaxLogLines is a number that will be certainly bigger than any number of logs
// Here 2*10^9 logs
var MaxLogLines int = 2000000000

const (
	NewestTimestamp		= "newest"
	OldestTimestamp		= "oldest"
)

// Load logs from begin or the end of the log file, because some log files are too large
const (
	Beginning 	= "beginning"
	End 				= "end"
)

// Representation of log lines
type LogDetails struct {
	// Additional information of the logs
	Info LogInfo `json:"info"`

	// Reference point to keep track of the position of all the logs
	Selector `json:"selection"`

	// Actual log lines of this page
	LogLines `json:"logs"`
}

// Meta information about the selected log lines
type LogInfo struct {
	// Name of pod
	PodName string `json:"podName"`
	// Name of container the logs are for
	ContainerName string `json:"containerName"`
	// Name of init container the logs are for
	InitContainerName string `json:"initContainerName"`
	// Date of the first log line
	FromDate LogTimestamp `json:"fromDate"`
	// Date of the last log line
	ToDate LogTimestamp `json:"toDate"`
	// Some logs lines in the middle of the log file could not be loaded, because the log file is too large
	Truncated bool `json:"truncated"`
}

// Default log view selector that is used in case of invalid request
// Downloads newest DefaultDisplayNumLogLines lines.
var DefaultSelector = &Selector{
	OffsetFrom:      1 - DefaultDisplayNumLogLines,
	OffsetTo:        1,
	ReferencePoint:  NewestLogLineId,
	LogFilePosition: End,
}

// Returns all logs
var AllSelector = &Selector{
	OffsetFrom:     -MaxLogLines,
	OffsetTo:       MaxLogLines,
	ReferencePoint: NewestLogLineId,
}

// Selector is a slice of log
type Selector struct {
	// ReferencePoint is the ID of a line which should serve as a reference point for this selector
	ReferencePoint LogLineId `json:"referencePoint"`
	// First index of the slice relatively to the reference line (this one will be included)
	OffsetFrom int `json:"offsetFrom"`
	// Last index of the slice relatively to the reference line (this one will not be included)
	OffsetTo int `json:"offsetTo"`
	// The log file is loaded either from begin and end. This matters only if the log file is too
	// large to be loaded (avoid oom)
	LogFilePosition string `json:"logFilePosition"`
}

// NewestLogLineId is the reference Id of the newest line
var NewestLogLineId = LogLineId{
	LogTimestamp: NewestTimestamp,
}

// OldestLogLineId is the reference Id of the oldest line
var OldestLogLineId = LogLineId{
	LogTimestamp: OldestTimestamp,
}

// LogLineId uniquely identifies a line in logs
type LogLineId struct {
	// timestamp of this line
	LogTimestamp `json:"timestamp"`
	// line number of the log
	LineNum int `json:"lineNum"`
}

type LogLines []LogLine

// A single log line. Split into timestamp and the actual content
type LogLine struct {
	Timestamp LogTimestamp `json:"timestamp"`
	Content string `json:"content"`
}

// LogTimestamp is a timestamp that appears on the beginning of each log line
type LogTimestamp string

func (self LogLines) SelectLogs(logSelector *Selector) (LogLines, LogTimestamp, LogTimestamp, Selector, bool) {
	requestedNumItems := logSelector.OffsetTo - logSelector.OffsetFrom
	referenceLineIndex := self.getLineIndex(&logSelector.ReferencePoint)
	if referenceLineIndex == LINE_INDEX_NOT_FOUND || requestedNumItems <= 0 || len(self) == 0 {
		return LogLines{}, "", "", Selector{}, false
	}
	fromIndex := referenceLineIndex + logSelector.OffsetFrom
	toIndex := referenceLineIndex + logSelector.OffsetTo
	lastPage := false
	if requestedNumItems > len(self) {
		fromIndex = 0
		toIndex = len(self)
		lastPage = true
	} else if toIndex > len(self) {
		toIndex += -fromIndex
		fromIndex = 0
		lastPage = logSelector.LogFilePosition == Beginning
	} else if fromIndex < 0 {
		toIndex += -fromIndex
		fromIndex = 0
		lastPage = logSelector.LogFilePosition == End
	}

	// set the middle of log array as a reference point, this part of array should not be affected by log deletion/addition.
	newSelection := Selector{
		ReferencePoint:  *self.createLogLineId(len(self) / 2),
		OffsetFrom:      fromIndex - len(self)/2,
		OffsetTo:        toIndex - len(self)/2,
		LogFilePosition: logSelector.LogFilePosition,
	}
	return self[fromIndex:toIndex], self[fromIndex].Timestamp, self[toIndex-1].Timestamp, newSelection, lastPage
}

// GetLineIndex returns the index of the line (referenced from beginning of log array) with provided logLineId
func (self LogLines) getLineIndex(logLineId *LogLineId) int {
	if logLineId == nil || logLineId.LogTimestamp == NewestTimestamp || len(self) == 0 || logLineId.LogTimestamp == "" {
		return len(self) - 1
	} else if logLineId.LogTimestamp == OldestTimestamp {
		return 0
	}

	logTimestamp := logLineId.LogTimestamp
	linesMatched := 0
	matchingStartedAt := 0
	for idx := range self {
		if self[idx].Timestamp == logTimestamp {
			if linesMatched == 0 {
				matchingStartedAt = idx
			}
			linesMatched += 1
		} else if linesMatched > 0 {
			break
		}
	}

	var offset int

	//         lineNume = -3
	// st  id /  |
	// |   | |   |
	// 2 2 2 2 2 1
	if logLineId.LineNum < 0 {
		offset = linesMatched + logLineId.LineNum
	} else {
		offset = logLineId.LineNum - 1
	}
	if 0 <= offset && offset < linesMatched {
		return matchingStartedAt + offset
	} else {
		return LINE_INDEX_NOT_FOUND
	}
}

func (self LogLines) createLogLineId(lineIndex int) *LogLineId {
	logTimestamp := self[lineIndex].Timestamp
	// determine whether use positive or negative indexing
	// check whether last line has the same index(logTimestamp) as requested line.
	// because if same, the unique is upper (less than) the given lineIndex
	var step int
	if self[len(self) - 1].Timestamp == logTimestamp {
		step = 1
	} else {
		step = -1
	}

	offset := step
	for ; 0 <= lineIndex - offset && lineIndex - offset < len(self); offset += step {
		if self[lineIndex - offset].Timestamp != logTimestamp {
			break
		}
	}

	// if last line has the same index. offset is positive, if logTimestamp appears 3 times
	// offset is 3.
	return &LogLineId{
		LogTimestamp: logTimestamp,
		LineNum: offset,
	}
}

// ToLogLines converts rawlogs (string) to LogLines. Proper log lines start with a timestamp which is chopped off
func ToLogLines(rawLogs string) LogLines {
	logLines := LogLines{}
	for _, line := range strings.Split(rawLogs, "\n") {
		if line != "" {
			startsWithDate := '0' <= line[0] && line[0] <= '9'
			idx := strings.Index(line, " ")
			if idx > 0 && startsWithDate {
				timestamp := LogTimestamp(line[0:idx])
				content := line[idx+1:]
				logLines = append(logLines, LogLine{Timestamp: timestamp, Content: content})
			} else {
				logLines = append(logLines, LogLine{Timestamp: LogTimestamp("0"), Content: line})
			}
		}
	}
	return logLines
}
