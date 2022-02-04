package funcstat

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"go/ast"
	"go/format"
	"go/types"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/cfg"
)

const doc = "funcstat get information"

var Analyzer = &analysis.Analyzer{
	Name:       "funcstat",
	Doc:        doc,
	Run:        run,
	ResultType: reflect.TypeOf((map[*ast.FuncDecl]*Result)(nil)),
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		ctrlflow.Analyzer,
	},
}

type Result struct {
	Func                 *ast.FuncDecl
	Name                 *ast.Ident
	Lines                int
	Bytes                int
	NumParams            int
	NumResults           int
	CyclomaticComplexity int
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	cfgs := pass.ResultOf[ctrlflow.Analyzer].(*ctrlflow.CFGs)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	cw := csv.NewWriter(os.Stdout)
	header := []string{
		"package",
		"file",
		"line",
		"name",
		"lines",
		"bytes",
		"params",
		"returns",
		"cyclomatic complexity",
	}

	var rerr error
	results := make(map[*ast.FuncDecl]*Result)
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		if rerr != nil {
			return
		}

		fundecl, _ := n.(*ast.FuncDecl)
		if fundecl == nil {
			return
		}

		sig, _ := pass.TypesInfo.TypeOf(fundecl.Name).(*types.Signature)
		if sig == nil {
			return
		}

		result := &Result{
			Func:       fundecl,
			Name:       fundecl.Name,
			NumParams:  sig.Params().Len(),
			NumResults: sig.Results().Len(),
		}

		var buf bytes.Buffer
		format.Node(&buf, pass.Fset, fundecl)

		result.Bytes = buf.Len()
		scanner := bufio.NewScanner(&buf)
		for scanner.Scan() {
			result.Lines++
		}

		result.CyclomaticComplexity = cyclomaticComplexity(cfgs.FuncDecl(fundecl))

		if len(results) == 0 {
			cw.Write(header)
			cw.Flush()
			if err := cw.Error(); err != nil {
				rerr = err
				return
			}
		}

		results[fundecl] = result

		if err := toCSV(pass, cw, result); err != nil {
			rerr = err
			return
		}
	})

	if rerr != nil {
		return (map[*ast.FuncDecl]*Result)(nil), rerr
	}

	return results, nil
}

func cyclomaticComplexity(g *cfg.CFG) int {
	c := 2
	for _, b := range g.Blocks {
		if b.Live {
			c += len(b.Succs) - 1
		}
	}
	return c
}

func toCSV(pass *analysis.Pass, cw *csv.Writer, r *Result) error {
	record := make([]string, 9)

	record[0] = pass.Pkg.Path()
	pos := pass.Fset.Position(r.Func.Pos())
	record[1] = filepath.Base(pos.Filename)
	record[2] = strconv.Itoa(pos.Line)
	record[3] = r.Name.Name
	record[4] = strconv.Itoa(r.Lines)
	record[5] = strconv.Itoa(r.Bytes)
	record[6] = strconv.Itoa(r.NumParams)
	record[7] = strconv.Itoa(r.NumResults)
	record[8] = strconv.Itoa(r.CyclomaticComplexity)

	cw.Write(record)
	cw.Flush()
	return cw.Error()
}
