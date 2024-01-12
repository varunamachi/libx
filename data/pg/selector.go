package pg

import (
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

// func GenQueryX(
// 	params *data.CommonParams,
// 	fmtStr string,
// 	fargs ...string) Selector {
// 	selector := NewSelectorGenerator().SelectorX(params)
// 	return NewSel(
// 		fmt.Sprintf(fmtStr, fargs)+" WHERE "+selector.QueryFragment,
// 		selector.Args,
// 	)
// }

// func GenQuery(
// 	filter *data.Filter,
// 	fmtStr string,
// 	fargs ...string) Selector {
// 	selector := NewSelectorGenerator().Selector(filter)
// 	return NewSel(
// 		fmt.Sprintf(fmtStr, fargs)+" WHERE "+selector.QueryFragment,
// 		selector.Args,
// 	)
// }

type SelectorGenerator struct {
	// filter      *data.Filter
	dollerIndex int64
	_args       []interface{}
	fragments   []string
}

func NewSelectorGenerator() *SelectorGenerator {
	return &SelectorGenerator{}
}

func (gen *SelectorGenerator) Reset() *SelectorGenerator {
	gen.dollerIndex = 1
	gen._args = make([]interface{}, 0, 100)
	gen.fragments = make([]string, 0, 30)
	return gen
}

func (gen *SelectorGenerator) Selector(filter *data.Filter) Selector {
	if filter == nil {
		return Selector{}
	}

	gen.Reset().
		matchers(filter.Props).
		bools(filter.Bools).
		dateRanges(filter.Dates).
		ranges(filter.Ranges).
		searches(filter.Searches)

	buf := buffer{}
	for idx, frag := range gen.fragments {
		buf.write(frag)
		if idx != 0 && idx < len(gen.fragments)-1 {
			buf.write(" AND ")
		}
	}
	return NewSel(buf.String(), gen._args)

}

func (gen *SelectorGenerator) SelectorX(cmnParam *data.CommonParams) Selector {

	// TODO -fix ANDs

	gen.Reset().
		matchers(cmnParam.Filter.Props).
		matchers(cmnParam.Filter.Lists).
		bools(cmnParam.Filter.Bools).
		dateRanges(cmnParam.Filter.Dates).
		ranges(cmnParam.Filter.Ranges).
		searches(cmnParam.Filter.Searches)

	buf := buffer{}
	for idx, frag := range gen.fragments {
		if frag == "" {
			continue
		}

		buf.write(frag)
		if idx < len(gen.fragments)-2 {
			buf.write(" AND ")
		}
	}

	if cmnParam.Limit() != 0 {
		buf.write(" OFFSET = \"$").writeInt(gen.dollerIndex).WriteString("\"")
		gen.addArg(cmnParam.Offset())
		buf.write(" LIMIT = \"$").writeInt(gen.dollerIndex).WriteString("\"")
		gen.addArg(cmnParam.Limit())
	}

	if cmnParam.Sort != "" {
		buf.write(" ORDER BY \"$").writeInt(gen.dollerIndex).WriteString("\"")
		gen.addArg(cmnParam.Sort)
		buf.write(data.Qop(cmnParam.SortDescending, " DESC", " ASC"))
	}

	return NewSel(buf.String(), gen._args)
}

func (gen *SelectorGenerator) matchers(
	pol map[string]*data.Matcher) *SelectorGenerator {
	if len(pol) == 0 {
		return gen
	}

	buf, idx := buffer{}, 0
	buf.Grow(100)

	for key, prop := range pol {
		if idx != 0 && idx <= len(prop.Fields) {
			buf.write(" AND ")
		}

		buf.write(key)
		if prop.Invert {
			buf.write(" NOT ")
		}
		buf.write(" IN (")
		for jdx, p := range prop.Fields {
			buf.write("\"$").writeInt(gen.dollerIndex).WriteString("\"")
			gen.addArg(p)
			if jdx < len(prop.Fields)-1 {
				buf.write(", ")
			}
		}
		buf.write(") ")
		idx++
	}
	gen.fragments = append(gen.fragments, buf.String())
	return gen
}

func (gen *SelectorGenerator) bools(
	bools map[string]interface{}) *SelectorGenerator {

	if len(bools) == 0 {
		return gen
	}

	buf, idx := buffer{}, 0
	buf.Grow(100)

	for key, boolVal := range bools {

		if idx != 0 && idx <= len(bools)-1 {
			buf.write(" AND ")
		}

		buf.write(key).
			write(" = \"$").
			writeInt(gen.dollerIndex).
			WriteString("\"")
		gen.addArg(boolVal) // :P
		idx++
	}
	gen.fragments = append(gen.fragments, buf.String())
	return gen
}

func (gen *SelectorGenerator) dateRanges(
	dates map[string]*data.DateRangeMatcher) *SelectorGenerator {
	if len(dates) == 0 {
		return gen
	}
	buf, idx := buffer{}, 0
	buf.Grow(100)

	for key, dt := range dates {

		if idx != 0 && idx <= len(dates)-1 {
			buf.write(" AND ")
		}

		buf.write("(")
		buf.write(key)
		if dt.Invert {
			buf.write(" NOT")
		}
		buf.write(" BETWEEN \"$").writeInt(gen.dollerIndex).WriteString("\"")
		gen.addArg(dt.From)
		buf.write(" AND \"$").writeInt(gen.dollerIndex).WriteString("\"")
		gen.addArg(dt.To)
		buf.write(")")

		idx++
	}
	gen.fragments = append(gen.fragments, buf.String())
	return gen
}

func (gen *SelectorGenerator) ranges(
	ranges map[string]*data.RangeMatcher) *SelectorGenerator {
	if len(ranges) == 0 {
		return gen
	}

	buf, idx := buffer{}, 0
	buf.Grow(100)

	for key, dt := range ranges {

		if idx != 0 && idx <= len(ranges)-1 {
			buf.write(" AND ")
		}

		buf.write("(")
		buf.write(key)
		if dt.Invert {
			buf.write(" NOT")
		}
		buf.write(" BETWEEN \"$").writeInt(gen.dollerIndex).WriteString("\"")
		gen.addArg(dt.From)
		buf.write(" AND \"$").writeInt(gen.dollerIndex).WriteString("\"")
		gen.addArg(dt.To)
		buf.write(")")

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
		if len(prop.Fields) == 0 {
			continue
		}

		if idx != 0 && idx <= len(searches)-1 {
			buf.write(" AND ")
		}

		buf.write("(")
		for jdx, p := range prop.Fields {
			buf.write(key)
			if prop.Invert {
				buf.write(" NOT")
			}
			buf.
				write(" SIMILAR TO \"$").
				writeInt(gen.dollerIndex).
				WriteString("\"")
			gen.addArg(p)
			if jdx < len(prop.Fields)-1 {
				buf.write(" OR ")
			}
		}
		buf.write(")")

		idx++
	}
	fragment := buf.String()
	gen.fragments = append(gen.fragments, fragment)
	return gen
}

func (gen *SelectorGenerator) addArg(arg interface{}) *SelectorGenerator {
	gen._args = append(gen._args, arg)
	gen.dollerIndex++
	return gen
}
