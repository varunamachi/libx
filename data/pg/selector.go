package pg

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/varunamachi/libx/data"
)

type buffer struct {
	strings.Builder
}

func (buf *buffer) write(str string) *buffer {
	buf.WriteString(str)
	return buf
}

func (buf *buffer) writeInt(val int64) *buffer {
	buf.WriteString(strconv.FormatInt(val, 10))
	return buf
}

func GenQueryX(
	params *data.QueryParams,
	fmtStr string,
	fargs ...string) (string, []interface{}) {
	selector, args := NewSelectorGenerator().SelectorX(params)
	return fmt.Sprintf(fmtStr, fargs) + " WHERE " + selector, args
}

func GenQuery(
	filter *data.Filter,
	fmtStr string,
	fargs ...string) (string, []interface{}) {
	selector, args := NewSelectorGenerator().Selector(filter)
	return fmt.Sprintf(fmtStr, fargs) + " WHERE " + selector, args
}

type SelectorGenerator struct {
	// filter      *data.Filter
	dollerIndex int64
	args        []interface{}
	fragments   []string
}

func NewSelectorGenerator() *SelectorGenerator {
	return &SelectorGenerator{}
}

func (gen *SelectorGenerator) Reset() *SelectorGenerator {
	gen.dollerIndex = 0
	gen.args = make([]interface{}, 0, 100)
	gen.fragments = make([]string, 0, 30)
	return gen
}

func (gen *SelectorGenerator) Selector(
	filter *data.Filter) (string, []interface{}) {
	gen.Reset().
		matchers(filter.Props).
		bools(filter.Bools).
		dateRanges(filter.Dates).
		ranges(filter.Ranges).
		searches(filter.Searches)

	buf := buffer{}
	for idx, frag := range gen.fragments {
		buf.write(frag)
		if idx < len(gen.fragments)-1 {
			buf.write(" AND ")
		}
	}
	return buf.String(), gen.args

}

func (gen *SelectorGenerator) SelectorX(
	cmnParam *data.QueryParams) (string, []interface{}) {

	filter := cmnParam.Filter
	gen.Reset().
		matchers(filter.Props).
		bools(filter.Bools).
		dateRanges(filter.Dates).
		ranges(filter.Ranges).
		searches(filter.Searches)

	buf := buffer{}
	for idx, frag := range gen.fragments {
		buf.write(frag)
		if idx < len(gen.fragments)-1 {
			buf.write(" AND ")
		}
	}

	buf.write(" OFFSET = $").writeInt(gen.dollerIndex)
	gen.addArg(cmnParam.Offset()).dollerIndex++
	buf.write(" LIMIT = $").writeInt(gen.dollerIndex)
	gen.addArg(cmnParam.Limit()).dollerIndex++
	buf.write(" ORDER BY $").writeInt(gen.dollerIndex)
	gen.addArg(cmnParam.Sort).dollerIndex++
	if cmnParam.SortDesc {
		buf.write(" DESC")
	}
	return buf.String(), gen.args
}

func (gen *SelectorGenerator) matchers(
	pol map[string]*data.Matcher) *SelectorGenerator {
	buf, idx := buffer{}, 0

	for key, prop := range pol {
		if len(prop.Fields) != 0 {
			continue
		}
		buf.write(key)
		if prop.Invert {
			buf.write(" NOT ")
		}
		buf.write(" IN (")
		for jdx, p := range prop.Fields {
			buf.write("$").writeInt(gen.dollerIndex)
			gen.addArg(p).dollerIndex++
			if jdx < len(prop.Fields)-1 {
				buf.write(", ")
			}
		}
		buf.write(") ")

		if idx < len(pol)-1 {
			buf.write(" AND ")
		}
		idx++
	}
	gen.fragments = append(gen.fragments, buf.String())
	return gen
}

func (gen *SelectorGenerator) bools(
	bools map[string]interface{}) *SelectorGenerator {
	buf, idx := buffer{}, 0

	for key, boolVal := range bools {
		buf.write(key).write(" = $").writeInt(gen.dollerIndex)
		gen.addArg(boolVal).dollerIndex++ // :P

		if idx < len(bools)-1 {
			buf.write(" AND ")
		}
		idx++
	}
	gen.fragments = append(gen.fragments, buf.String())
	return gen
}

func (gen *SelectorGenerator) dateRanges(
	dates map[string]*data.DateRangeMatcher) *SelectorGenerator {
	buf, idx := buffer{}, 0

	for key, dt := range dates {
		buf.write(key)
		if dt.Invert {
			buf.write(" NOT")
		}
		buf.write(" BETWEEN $").writeInt(gen.dollerIndex)
		gen.addArg(dt.From).dollerIndex++
		buf.write(" AND $").writeInt(gen.dollerIndex)
		gen.addArg(dt.To).dollerIndex++

		if idx < len(dates)-1 {
			buf.write(" AND ")
		}
		idx++
	}
	gen.fragments = append(gen.fragments, buf.String())
	return gen
}

func (gen *SelectorGenerator) ranges(
	ranges map[string]*data.RangeMatcher) *SelectorGenerator {
	buf, idx := buffer{}, 0

	for key, dt := range ranges {
		buf.write(key)
		if dt.Invert {
			buf.write(" NOT")
		}
		buf.write(" BETWEEN $").writeInt(gen.dollerIndex)
		gen.addArg(dt.From).dollerIndex++
		buf.write(" AND $").writeInt(gen.dollerIndex)
		gen.addArg(dt.To).dollerIndex++

		if idx < len(ranges)-1 {
			buf.write(" AND ")
		}
		idx++
	}
	gen.fragments = append(gen.fragments, buf.String())
	return gen
}

func (gen *SelectorGenerator) searches(
	searches map[string]*data.Matcher) *SelectorGenerator {
	buf, idx := buffer{}, 0

	for key, prop := range searches {
		if len(prop.Fields) != 0 {
			continue
		}
		for jdx, p := range prop.Fields {
			buf.write(key)
			if prop.Invert {
				buf.write(" NOT ")
			}
			buf.write("SIMILAR TO $").writeInt(gen.dollerIndex)
			gen.addArg(p).dollerIndex++
			if jdx < len(prop.Fields)-1 {
				buf.write(" OR ")
			}
		}

		if idx < len(searches)-1 {
			buf.write(" AND ")
		}
		idx++
	}
	gen.fragments = append(gen.fragments, buf.String())
	return gen
}

func (gen *SelectorGenerator) addArg(arg interface{}) *SelectorGenerator {
	gen.args = append(gen.args, arg)
	return gen
}
