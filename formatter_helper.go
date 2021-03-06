package multilogger

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	multicolor "github.com/coveooss/multilogger/color"
	"github.com/coveooss/multilogger/errors"
	"github.com/coveooss/multilogger/reutils"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

func (f *Formatter) doFormat(entry *logrus.Entry) (string, error) {
	f.replacerLock.Lock()
	if f.replacer == nil {
		if err := f.presetFormatString(); err != nil {
			return "", err
		}
	}
	defer f.replacerLock.Unlock()

	output := f.replacer.format + "\n"

	usedFields := make(map[string]uint, len(entry.Data))
	var printFields []*fieldReplacer
	for _, replacer := range f.replacer.fields {
		result, delayed := replacer.replace(entry, usedFields)
		if delayed != nil {
			printFields = append([]*fieldReplacer{delayed}, printFields...)
		}
		output = output[:replacer.position] + result + output[replacer.position:]
	}

	if len(printFields) > 0 {
		for _, replacer := range printFields {
			result := make([]string, 0, len(entry.Data))
			for key := range entry.Data {
				if usedFields[key] > 0 || entry.Data[key] == nil && replacer.ignoreEmpty {
					continue
				}
				result = append(result, key)
			}
			sort.Strings(result)
			for i, key := range result {
				if replacer.noKeyFieldFormat {
					result[i] = fmt.Sprintf("%s=%v", key, entry.Data[key])
				} else {
					result[i] = f.replacer.keyReplacer.formatValue(key+"=", entry.Level) +
						f.replacer.fieldReplacer.formatValue(entry.Data[key], entry.Level)
				}
			}
			fields := strings.Join(result, " ")
			printKey := replacer.printKey && (!replacer.ignoreEmpty || fields != "")
			fields = replacer.format("Fields", fields, printKey, entry.Level)
			output = strings.Replace(output, replacementToken, fields, 1)
		}
	}

	return output, nil
}

func (f *Formatter) presetFormatString() error {
	var errors errors.Array
	f.replacer = &replacer{Formatter: f}
	r := f.replacer
	r.keyReplacer = r.newField(true)
	r.fieldReplacer = r.newField(false)
	r.format = reFormat.ReplaceAllStringFunc(f.format, func(match string) (result string) {
		matches, _ := reutils.MultiMatch(match, reFormat)

		result = replacementToken
		fieldReplacer := r.newField(false)
		if field := matches["field"]; field != "" {
			fieldReplacer.tt = fieldTokenType
			fieldReplacer.fieldName = field
			f.replacer.fields = append(r.fields, fieldReplacer)
		} else if token := matches["token"]; token != "" {
			fieldReplacer.tt = reverseTokens[token]
			if fieldReplacer.tt == unsetTokenType {
				switch token {
				case "field":
					fieldReplacer.tt = fieldWrapperTokenType
					r.fieldReplacer = fieldReplacer
					result = ""
				case "key":
					fieldReplacer.tt = keyWrapperTokenType
					r.keyReplacer = fieldReplacer
					result = ""
				}
			} else {
				r.fields = append(r.fields, fieldReplacer)
			}
		} else {
			result = ""
		}

		if limit, err := strconv.ParseUint(matches["limit"], 10, 0); err == nil {
			limit := uint(limit)
			fieldReplacer.limit = &limit
		}
		if width, err := strconv.Atoi(matches["width"]); err == nil {
			fieldReplacer.width = &width
		}

		attributes := strings.Split(matches["attributes"], ",")
		colors := make([]interface{}, 0, len(attributes))
		for _, attribute := range attributes {
			attribute = strings.ToLower(strings.TrimSpace(attribute))
			if attribute == "" {
				continue
			}
			switch attribute {
			case "color":
				fieldReplacer.color = true
			case "upper":
				fieldReplacer.transform = uppercaseTransform
			case "lower":
				fieldReplacer.transform = lowercaseTransform
			case "title":
				fieldReplacer.transform = titleTransform
			case "key":
				fieldReplacer.printKey = true
			case "curly", "curlybrackets":
				fieldReplacer.wrapper = curlyBrackets
			case "square", "squarebrackets":
				fieldReplacer.wrapper = squareBrackets
			case "round", "roundbrackets", "parens", "parenthesis":
				fieldReplacer.wrapper = roundBrackets
			case "angle", "anglebrackets":
				fieldReplacer.wrapper = angleBrackets
			case "space":
				fieldReplacer.addSpace = true
			case "none":
				fieldReplacer.noKeyFieldFormat = true
			case "ignore", "ignoreempty":
				fieldReplacer.ignoreEmpty = true
			default:
				if f.color {
					colors = append(colors, attribute)
				}
			}
		}
		if len(colors) > 0 {
			var err error
			if fieldReplacer.attributes, err = multicolor.TryConvertAttributes(colors); err != nil {
				errors = append(errors, err)
			}
			if fieldReplacer.tt == unsetTokenType {
				result = color.New(fieldReplacer.attributes...).Sprint()
				if result != reset {
					// There is no token or field specified, in that case, we do not reset the color attributes.
					result = strings.TrimSuffix(result, reset)
				}
			}
		}

		return
	})

	rePlaceHolder, l := regexp.MustCompile(replacementToken), len(replacementToken)
	fields := r.fields
	r.fields = make([]*fieldReplacer, len(fields))
	for i, position := range rePlaceHolder.FindAllStringIndex(r.format, -1) {
		begin, end := position[0]-i*l, position[1]-i*l
		r.format = r.format[:begin] + r.format[end:]
		index := len(fields) - i - 1
		r.fields[index] = fields[i]
		r.fields[index].position = uint(begin)
	}

	return errors.AsError()
}

// https://regex101.com/r/SPI8hT/1
var (
	reFormat = regexp.MustCompile(`%(?:(?P<width>-?\d+)?(?:\.(?P<limit>\d+))?(?:(?P<token>(?:time|(?:global)?delay|delta|message|msg|level|lvl|module|func|file|line|caller|fields|key|field))|(?P<field>\w+)))?(?i)(?::(?P<attributes>(?:[,+\-\s]*(?:color|upper|lower|title|none|key|ignore(?:empty)?|space|parens|parenthesis|(?:square|curly|round|angle)(?:brackets)?|(?:bg)?(?:hi)?(?:black|red|green|yellow|blue|magenta|cyan|white)|(?:bold|faint|italic|underline|blinkslow|blinkrapid|reversevideo|concealed|crossedout|reset))\s*)+))?%`)
	reset    = string([]byte{27, 91, 48, 109})
)

type transformType uint8

const (
	noTransform transformType = iota
	uppercaseTransform
	lowercaseTransform
	titleTransform
)

type bracketType uint8

const (
	noBracket bracketType = iota
	curlyBrackets
	angleBrackets
	roundBrackets
	squareBrackets
)

type tokenType uint8

const (
	unsetTokenType tokenType = iota
	tokenMessage
	tokenLevel
	tokenTime
	tokenDelta
	tokenDelay
	tokenGlobalDelay
	tokenModule
	tokenFunc
	tokenFile
	tokenLine
	tokenCaller
	tokenFields
	fieldTokenType
	fieldWrapperTokenType
	keyWrapperTokenType
)

func init() {
	reverseTokens = make(map[string]tokenType)
	for t := unsetTokenType + 1; t < fieldTokenType; t++ {
		reverseTokens[strings.ToLower(t.String())] = t
	}
	reverseTokens["lvl"] = tokenLevel
	reverseTokens["msg"] = tokenMessage
	reverseTokens["global"] = tokenGlobalDelay
}

var reverseTokens map[string]tokenType

//go:generate stringer -type=tokenType -trimprefix token -output formatter_generated.go
