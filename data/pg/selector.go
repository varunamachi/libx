package pg

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/varunamachi/libx/data"
)

type Selector struct {
	QueryFragment string
	Args          []interface{}
}

func (sel Selector) IsEmpty() bool {
	return sel.QueryFragment == "" && len(sel.Args) == 0
}

func NewSel(qf string, args []interface{}) Selector {
	return Selector{
		QueryFragment: qf,
		Args:          args,
	}
}

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
	params *data.CommonParams,
	fmtStr string,
	fargs ...string) Selector {
	selector := NewSelectorGenerator().SelectorX(params)
	return NewSel(
		fmt.Sprintf(fmtStr, fargs)+" WHERE "+selector.QueryFragment,
		selector.Args,
	)
}

func GenQuery(
	filter *data.Filter,
	fmtStr string,
	fargs ...string) Selector {
	selector := NewSelectorGenerator().Selector(filter)
	return NewSel(
		fmt.Sprintf(fmtStr, fargs)+" WHERE "+selector.QueryFragment,
		selector.Args,
	)
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

func (gen *SelectorGenerator) Selector(filter *data.Filter) Selector {
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
	return NewSel(buf.String(), gen.args)

}

func (gen *SelectorGenerator) SelectorX(cmnParam *data.CommonParams) Selector {

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

	if cmnParam.Limit() != 0 {
		buf.write(" OFFSET = $").writeInt(gen.dollerIndex)
		gen.addArg(cmnParam.Offset()).dollerIndex++
		buf.write(" LIMIT = $").writeInt(gen.dollerIndex)
		gen.addArg(cmnParam.Limit()).dollerIndex++
	}

	if cmnParam.Sort != "" {
		buf.write(" ORDER BY $").writeInt(gen.dollerIndex)
		gen.addArg(cmnParam.Sort).dollerIndex++
		buf.write(data.Qop(cmnParam.SortDesc, " DESC", " ASC"))
	}

	return NewSel(buf.String(), gen.args)
}

func (gen *SelectorGenerator) matchers(
	pol map[string]*data.Matcher) *SelectorGenerator {
	buf, idx := buffer{}, 0
	buf.Grow(100)

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
	buf.Grow(100)

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
	buf.Grow(100)

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
	buf.Grow(100)

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
	buf.Grow(100)

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
