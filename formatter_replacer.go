package multilogger

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/acarl005/stripansi"
	multicolor "github.com/coveooss/multilogger/color"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type replacer struct {
	*Formatter
	format        string
	fields        []*fieldReplacer
	keyReplacer   *fieldReplacer
	fieldReplacer *fieldReplacer
}

func (r *replacer) newField(color bool) *fieldReplacer {
	return &fieldReplacer{
		replacer: r,
		color:    color,
		out:      make(map[logrus.Level]func(...interface{}) string),
	}
}

type fieldReplacer struct {
	*replacer
	transform        transformType
	wrapper          bracketType
	addSpace         bool
	ignoreEmpty      bool
	printKey         bool
	noKeyFieldFormat bool
	tt               tokenType
	attributes       []multicolor.Attribute
	fieldName        string
	color            bool
	width            *int
	limit            *uint
	position         uint
	out              map[logrus.Level]func(...interface{}) string
	outMutex         sync.RWMutex
}

func (r *fieldReplacer) replace(entry *logrus.Entry, used map[string]uint) (string, *fieldReplacer) {
	var field string
	key, printKey := r.tt.String(), r.printKey

	computeduration := func(begin time.Time) string {
		delay := entry.Time.Sub(begin)
		round := r.RoundDuration
		if round == 0 {
			round = roundDuration
		}
		if delay = delay.Round(round); delay == 0 {
			return fmt.Sprintf("<%s", r.FormatDuration(round))
		}
		if r.FormatDuration != nil {
			return r.FormatDuration(delay)
		}
		return delay.String()
	}

	// Find the right replacement
	switch r.tt {
	case fieldTokenType:
		key = r.fieldName
		value := entry.Data[key]
		if r.ignoreEmpty && value == nil || value == "" {
			field = ""
			printKey = false
		} else {
			field = fmt.Sprint(value)
		}

		used[key]++
	case tokenMessage:
		field = entry.Message
	case tokenLevel:
		field = r.LevelName[entry.Level]
		if field == "" {
			field = fmt.Sprint(entry.Level)
		}
	case tokenTime:
		if globalZone != nil {
			field = entry.Time.In(globalZone).Format(r.TimestampFormat)
		} else {
			field = entry.Time.Format(r.TimestampFormat)
		}
	case tokenDelta:
		field = computeduration(r.last)
	case tokenDelay:
		field = computeduration(r.baseTime)
	case tokenGlobalDelay:
		field = computeduration(globalTime)
	case tokenModule:
		field = fmt.Sprint(entry.Data[moduleFieldName])
		used[moduleFieldName]++
	case tokenFunc:
		if entry.Caller != nil {
			field = entry.Caller.Function
		}
	case tokenFile:
		if entry.Caller != nil {
			field = entry.Caller.File
		}
	case tokenLine:
		if entry.Caller != nil {
			field = fmt.Sprint(entry.Caller.Line)
		}
	case tokenCaller:
		if entry.Caller != nil {
			field = r.FormatCaller(entry.Caller)
		}
	case tokenFields:
		return replacementToken, r
	}

	return r.format(key, field, printKey, entry.Level), nil
}

func (r *fieldReplacer) formatValue(value interface{}, level logrus.Level) string {
	return r.format("", fmt.Sprint(value), false, level)
}

func (r *fieldReplacer) format(key, value string, printKey bool, level logrus.Level) string {
	if r.ignoreEmpty && strings.TrimSpace(value) == "" {
		return ""
	}

	if !r.Formatter.color {
		value = stripansi.Strip(value)
	}
	// Process the field transformation
	switch r.transform {
	case uppercaseTransform:
		value = strings.ToUpper(value)
	case lowercaseTransform:
		value = strings.ToLower(value)
	case titleTransform:
		value = strings.Title(value)
	}

	// Process the spaces
	if r.limit != nil {
		value = fmt.Sprintf("%.*s", *r.limit, value)
	}
	if r.width != nil {
		value = fmt.Sprintf("%*s", *r.width, value)
	}

	switch r.wrapper {
	case squareBrackets:
		value = fmt.Sprintf("[%s]", value)
	case curlyBrackets:
		value = fmt.Sprintf("{%s}", value)
	case roundBrackets:
		value = fmt.Sprintf("(%s)", value)
	case angleBrackets:
		value = fmt.Sprintf("<%s>", value)
	}

	if r.addSpace {
		value += " "
	}

	r.outMutex.RLock()
	sprint := r.out[level]
	r.outMutex.RUnlock()
	if sprint == nil {
		if r.Formatter.color && (r.color || len(r.attributes) != 0) {
			var attributes []multicolor.Attribute
			attributes = append(attributes, r.attributes...)
			if r.color {
				attributes = append(attributes, r.ColorMap[level]...)
			}
			sprint = color.New(attributes...).Sprint
		} else {
			sprint = fmt.Sprint
		}
		r.outMutex.Lock()
		r.out[level] = sprint
		r.outMutex.Unlock()
	}
	value = sprint(value)

	// Process the prefix
	if r != r.fieldReplacer && r != r.keyReplacer {
		if printKey {
			if r.noKeyFieldFormat {
				value = fmt.Sprintf("%s=%s", key, value)
			} else {
				value = r.keyReplacer.formatValue(key+"=", level) + r.fieldReplacer.formatValue(value, level)
			}
		} else if !r.noKeyFieldFormat {
			value = r.fieldReplacer.formatValue(value, level)
		}
	}

	return value
}

const replacementToken = "💚"
