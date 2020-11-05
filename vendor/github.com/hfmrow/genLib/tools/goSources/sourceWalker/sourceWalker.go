// sourceWalker.go

package sourceWalker

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	glfs "github.com/hfmrow/genLib/files"
	glsg "github.com/hfmrow/genLib/strings"
)

// GoSourceFileStructure: contain AST file informations
type GoSourceFileStructure struct {
	Filename string
	Package  string

	Imports []imported
	Func    []function
	Struct  []structure
	Var     []variable

	Eol string // End of line of the input file

	// Unexported
	data          []byte        // File content
	astOut        *bytes.Buffer // AST representation of the input file
	linesIndexes  [][]int
	offset        int // Define if we start at 0 or 1  when counting lines and offsets positions.
	astFile       *ast.File
	fset          *token.FileSet // Positions are relative to fset.
	tmpMethods    []function
	varInsideFunc bool
}

// type ShortDescr struct{ imported } // Reserved for enduser usage ...

type imported struct {
	Name        string
	NameFromSrc string
	Content     content
	File        string
}

type function struct {
	Ident    identObj
	Content  content
	File     string
	Exported bool
}

type structure struct {
	Ident    identObj
	Content  content
	Fields   []field
	Methods  []function
	File     string
	Exported bool
}

type variable struct {
	Objects  field
	Content  content
	File     string
	Found    identObj
	Exported bool
}

type field struct {
	List []identObj
	Type string
	Name string
}

type identObj struct {
	Name     string
	Kind     string
	Type     string
	Value    string
	Idx      int
	Exported bool
}

type content struct {
	OfstStart int
	OfstEnd   int
	LineStart int
	LineEnd   int
	Content   []byte
	Comment   string
}

// GetFuncByName: Optional "unExported": empty = both.
func (gsfs *GoSourceFileStructure) GetFuncByName(fName string, unExported ...bool) (funct *function) {

	for idx, fnc := range gsfs.Func {
		if fnc.Ident.Name == fName {
			return &gsfs.Func[idx]
		}
	}

	return
}

// GetStructByName:
func (gsfs *GoSourceFileStructure) GetStructByName(sName string) (stru *structure) {

	for idx, stc := range gsfs.Struct {
		if stc.Ident.Name == sName {
			return &gsfs.Struct[idx]
		}
	}

	return
}

// GetVarByName: "Position" contain the position in "list" and "values" fields.
func (gsfs *GoSourceFileStructure) GetVarByName(vName string) (vari *variable) {

	for idx, vr := range gsfs.Var {
		for idn, vn := range vr.Objects.List {
			if vn.Name == vName {
				tmpV := gsfs.Var[idx]
				tmpV.Found.Name = vn.Name
				tmpV.Found.Type = vr.Objects.Type
				tmpV.Found.Value = vn.Value
				tmpV.Found.Idx = idn
				return &tmpV
			}
		}
	}

	return
}

// varWalker: called from the walker or others to check for variables tokens
func varWalker(t token.Token) (kind string, ok bool) {
	switch t {
	case token.VAR:
		return "var", true
	case token.CONST:
		return "const", true
	case token.ASSIGN:
		return "=", true
	case token.DEFINE:
		return ":=", true
	}
	return
}

// getImportOnly:
func (gsfs *GoSourceFileStructure) getImportOnly() (err error) {
	if filename, err := filepath.Rel(filepath.Join(os.Getenv("GOPATH"), "src"), gsfs.Filename); err == nil {
		// getImports
		for _, val := range gsfs.astFile.Imports {
			content := gsfs.getContentFromPos(val.Pos(), val.End())
			gsfs.Imports = append(gsfs.Imports,
				imported{
					Name:        val.Path.Value,
					NameFromSrc: string(gsfs.data[val.Pos()-1 : val.End()-1]),
					Content:     content,
					File:        filename})
		}
	}
	return
}

// goInspect: parse go file and retrieve into structure that was found.
func (gsfs *GoSourceFileStructure) goInspect() {
	filename, _ := filepath.Rel(filepath.Join(os.Getenv("GOPATH"), "src"), gsfs.Filename)
	gsfs.getImportOnly()
	// Get nodes infos
	ast.Inspect(gsfs.astFile, func(node ast.Node) bool {
		switch val := node.(type) {
		case *ast.FuncDecl: // Functions
			exported := val.Name.IsExported()
			content := gsfs.getContentFromPos(val.Pos(), val.End(), val.Doc.Text())
			if val.Recv == nil { // Functions
				var funct function
				obj := gsfs.getIdent(val.Name)
				funct.Ident.Name = obj.Name
				funct.Ident.Kind = obj.Kind
				funct.Content = content
				funct.File = filename
				funct.Exported = exported
				gsfs.Func = append(gsfs.Func, funct)
			} else { // Methods
				var method function
				method.Ident.Name = gsfs.getIdent(val.Name).Name
				method.Ident.Type = gsfs.getFields(val.Recv.List).Ident.Type
				method.Content = content
				method.File = filename
				method.Exported = exported
				gsfs.tmpMethods = append(gsfs.tmpMethods, method)
			}
		case *ast.GenDecl:
			for _, spec := range val.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					exported := s.Name.IsExported()
					stru := gsfs.getStruct(s)
					stru.Content = gsfs.getContentFromPos(val.Pos(), val.End(), val.Doc.Text())
					stru.Ident.Type = ""
					stru.File = filename
					stru.Exported = exported
					gsfs.Struct = append(gsfs.Struct, stru)

				case *ast.ValueSpec:
					if _, ok := varWalker(val.Tok); ok && !gsfs.varInsideFunc {
						if fld := gsfs.getSpecs([]ast.Spec{s}); fld != nil {
							gsfs.Var = append(gsfs.Var, variable{
								Objects: *fld,
								File:    filename,
								Content: gsfs.getContentFromPos(val.Pos(), val.End(), s.Doc.Text()),
							})
						}
					}
				}
			}
		}
		return true
	})
}

// getIdent:
func (gsfs *GoSourceFileStructure) getIdent(ident *ast.Ident) (obj identObj) {
	obj.Name = ident.Name
	if ident.Obj != nil {
		obj.Kind = ident.Obj.Kind.String()
		if ident.Obj.Type != nil {
			obj.Type = ident.Obj.Type.(*ast.Ident).Name
		}
		if ident.Obj.Data != nil {
			switch iod := ident.Obj.Data.(type) {
			case *ast.BasicLit:
				obj.Value = iod.Value
			}
		}
	}
	return
}

// getBasicLit
func (gsfs *GoSourceFileStructure) getBasicLit(bl *ast.BasicLit) (outValue, outType string) {
	return bl.Value, strings.ToLower(bl.Kind.String())
}

// getAssignStmt:
func (gsfs *GoSourceFileStructure) getAssignStmt(aStmt *ast.AssignStmt) (fld *field) {

	var obj identObj
	fld = new(field)
	// fld.Type = aStmt.Tok.String()
	for _, lhs := range aStmt.Lhs {
		switch lhst := lhs.(type) {
		case *ast.Ident:
			obj = gsfs.getIdent(lhst)
			obj.Kind = aStmt.Tok.String()
			obj.Type = "func"
		}
		fld.List = append(fld.List, obj)
	}
	for idx, rhs := range aStmt.Rhs {
		switch rhst := rhs.(type) {
		case *ast.BasicLit:
			fld.List[idx].Value, fld.Type = gsfs.getBasicLit(rhst)

		case *ast.Ident:
			fld.List[idx].Value = gsfs.getIdent(rhst).Name
			if fld.List[idx].Value == "true" || fld.List[idx].Value == "false" {
				fld.Type = "bool"
			}
		}
	}

	return
}

// getSpecs:
func (gsfs *GoSourceFileStructure) getSpecs(specs []ast.Spec) (fld *field) {
	for _, spec := range specs {
		var tmpStr string
		switch s := spec.(type) {
		case *ast.ValueSpec:
			fld = new(field)
			for _, idnt := range s.Names {
				obj := gsfs.getIdent(idnt)
				obj.Exported = idnt.IsExported()
				fld.List = append(fld.List, obj)
			}
			if s.Type != nil {
				switch st := s.Type.(type) {
				case *ast.Ident:
					fld.Type = gsfs.getIdent(st).Name
				case *ast.StarExpr:
					if st.X != nil {
						fld.Type = gsfs.getStarX(st.X)
					}
				case *ast.ArrayType:
					fld.Type = gsfs.getArray(st)
				}
			}
			var values []string
			if s.Values != nil {
				for _, value := range s.Values {
					switch v := value.(type) {
					case *ast.FuncLit: // Right arg is a function
						fld.Type = "func"
					case *ast.Ident:
						fld.Type = gsfs.getIdent(v).Kind
						fld.Name = gsfs.getIdent(v).Name
					case *ast.SelectorExpr:
						if v.X != nil {
							switch vXt := v.X.(type) {
							case *ast.Ident:
								tmpStr = gsfs.getIdent(vXt).Name
							}
						}
						if v.Sel != nil {
							tmpStr += "." + gsfs.getIdent(v.Sel).Name
						}
						fld.Type = tmpStr
					case *ast.BasicLit:
						fld.Type = strings.ToLower(v.Kind.String())
						values = append(values, v.Value)
					case *ast.CallExpr:
						if v.Fun != nil {
							switch vt := v.Fun.(type) {
							case *ast.Ident:
								fld.Type = gsfs.getIdent(vt).Name
							}
						}
						if v.Args != nil {
							for _, arg := range v.Args {
								switch va := arg.(type) {
								case *ast.BasicLit:
									values = append(values, va.Value)
								}
							}
						}
					}
				}
			}
			// Fill with values
			for idx := len(fld.List) - 1; idx >= 0; idx-- {
				if len(values) > idx {
					fld.List[idx].Value = values[idx]
				}
			}
			if len(fld.List) == 0 {
				fld = nil
			}
		case *ast.TypeSpec: // Struct
			// (not implemented here. it was done above)
		}
	}
	return
}

// getArray:
func (gsfs *GoSourceFileStructure) getArray(ary *ast.ArrayType) (fld string) {
	switch fvv := ary.Elt.(type) {
	case *ast.ArrayType:
		switch fvvE := fvv.Elt.(type) {
		case *ast.Ident:
			fld = "[][]" + fvvE.Name
		case *ast.StarExpr:
			fld = "[][]" + gsfs.getStarX(fvvE.X)
		}
	case *ast.Ident:
		fld = "[]" + fvv.Name
	case *ast.StarExpr:
		fld = "[]" + gsfs.getStarX(fvv.X)
	}
	return
}

// getStarX:
func (gsfs *GoSourceFileStructure) getStarX(sX ast.Expr) (fld string) {
	switch fvX := sX.(type) {
	case *ast.Ident:
		fld = "*" + fvX.Name
	case *ast.ArrayType:
		fld = "*" + gsfs.getArray(fvX)
	}
	return
}

// getField:
func (gsfs *GoSourceFileStructure) getFields(fields []*ast.Field) (stru structure) {
	for _, fList := range fields {
		var fld field
		for _, ident := range fList.Names {
			fld.List = append(fld.List, gsfs.getIdent(ident))
		}
		switch fv := fList.Type.(type) {
		case *ast.Ident:
			fld.Type = fv.Name
		case *ast.StarExpr:
			fld.Type = gsfs.getStarX(fv.X)
		case *ast.ArrayType:
			fld.Type = gsfs.getArray(fv)
		}
		stru.Ident.Type = fld.Type
		stru.Fields = append(stru.Fields, fld)
	}
	return
}

// getStruct:
func (gsfs *GoSourceFileStructure) getStruct(s *ast.TypeSpec) (stru structure) {
	switch st := s.Type.(type) {
	case *ast.StructType:
		stru = gsfs.getFields(st.Fields.List)
	}
	// stru = gsfs.getFields(s.Type.(*ast.StructType).Fields.List)
	obj := gsfs.getIdent(s.Name)
	stru.Ident.Name = obj.Name
	stru.Ident.Kind = obj.Kind
	return
}

func (gsfs *GoSourceFileStructure) GetImportsOnly(filename string) (err error) {
	gsfs.Filename = filename
	if err = gsfs.loadDataFile(); err == nil {
		err = gsfs.getImportOnly()
	}
	return
}

// GoSourceFileStructureSetup: setup and retieve information for designed file.
// Notice: the lines numbers and offsets start at 0. Set "zero" at false to start at 1.
func (gsfs *GoSourceFileStructure) GoSourceFileStructureSetup(filename string, zero ...bool) (err error) {
	gsfs.offset = 1 // lines start at 0 (substract 1 for each offsets position)
	if len(zero) > 0 {
		if zero[0] {
			gsfs.offset = 0
		}
	}
	gsfs.Filename = filename
	if err = gsfs.loadDataFile(); err == nil {
		if err = gsfs.fillDeclaration(); err == nil {
			gsfs.filteringMethods()
		}
	}
	return
}

// AppendFile:
func (gsfs *GoSourceFileStructure) AppendFile(filename string) (err error) {
	gsfs.tmpMethods = []function{}
	return gsfs.GoSourceFileStructureSetup(filename)
}

func (gsfs *GoSourceFileStructure) loadDataFile() (err error) {
	// Loading data (file)
	gsfs.fset = token.NewFileSet()
	if gsfs.astFile, err = parser.ParseFile(gsfs.fset, gsfs.Filename, nil, parser.ParseComments); err == nil {
		gsfs.data, err = ioutil.ReadFile(gsfs.Filename)
	}
	return
}

// fillDeclaration:
func (gsfs *GoSourceFileStructure) fillDeclaration() (err error) {
	// Setting internal variables
	gsfs.Eol = glsg.GetTextEOL(gsfs.data)
	gsfs.data = append(gsfs.data, []byte(gsfs.Eol)...) // Add an eol to avoid a f..k..g issue where the last line wasn't analysed
	eolRegx := regexp.MustCompile(gsfs.Eol)
	eolPositions := eolRegx.FindAllIndex(gsfs.data, -1)

	// Define and prepare slice of line indexes
	gsfs.linesIndexes = make([][]int, len(eolPositions)+1)
	gsfs.linesIndexes[0] = []int{0, eolPositions[0][0]}
	// Creating lines indexes
	for idx := 1; idx < len(eolPositions); idx++ {
		gsfs.linesIndexes[idx] = []int{eolPositions[idx-1][1], eolPositions[idx][0]}
	}

	gsfs.Package = gsfs.astFile.Name.String() // get package name
	gsfs.goInspect()
	return
}

// getLineFromOffsets: get the line number corresponding to offsets. Notice, line number start at 0.
func (gsfs *GoSourceFileStructure) getLineFromOffsets(sOfst, eOfst int) (lStart, lEnd int) {
	for lineNb, lineIdxs := range gsfs.linesIndexes {
		switch {
		case sOfst >= lineIdxs[0] && sOfst <= lineIdxs[1]:
			lStart = lineNb
			if eOfst <= lineIdxs[1] { // only one line
				lEnd = lineNb
				return
			}
		case eOfst >= lineIdxs[0] && eOfst <= lineIdxs[1]:
			lEnd = lineNb
			return
		}
	}
	return
}

// getContentFromPos: fill content structure
func (gsfs *GoSourceFileStructure) getContentFromPos(pos, end token.Pos, comment ...string) (cnt content) {
	// Set to relative offset
	sOfst := gsfs.fset.PositionFor(pos, true).Offset
	eOfst := gsfs.fset.PositionFor(end, true).Offset
	// Make content structure
	cnt.OfstStart, cnt.OfstEnd = sOfst-gsfs.offset, eOfst-gsfs.offset
	cnt.LineStart, cnt.LineEnd = gsfs.getLineFromOffsets(sOfst-1, eOfst-1)
	cnt.Content = gsfs.data[sOfst-1 : eOfst-1]
	if len(comment) > 0 {
		cnt.Comment = comment[0]
	}
	return
}

// astPrintToBuf: Simply display ast content for an overview of declarations. DEBUG purpose ...
func (gsfs *GoSourceFileStructure) astPrintToBuf(saveToFilename ...string) (bytesBuf *bytes.Buffer, err error) {
	var writer io.Writer
	bytesBuf = new(bytes.Buffer)
	writer = bytesBuf

	err = ast.Fprint(writer, gsfs.fset, gsfs.astFile, ast.NotNilFilter)
	if len(saveToFilename) > 0 {
		if len(saveToFilename[0]) != 0 {
			OS := glfs.OsPermsStructNew()
			if err = ioutil.WriteFile(saveToFilename[0], bytesBuf.Bytes(), OS.GROUP_RW|OS.USER_RW|OS.ALL_R); err == nil {
				fmt.Printf("AST file saved successfully: %s\n", saveToFilename[0])
			} else {
				fmt.Printf("Unable to save AST file: %s\n", err.Error())
			}
		}
	}
	return bytesBuf, err
}

// filteringMethods: put methods with their respective structures
func (gsfs *GoSourceFileStructure) filteringMethods() {
	for idx, stru := range gsfs.Struct {
		for _, mtd := range gsfs.tmpMethods {

			if stru.Ident.Name == strings.Trim(mtd.Ident.Type, "*") {
				gsfs.Struct[idx].Methods = append(gsfs.Struct[idx].Methods, mtd)
			}
		}
	}
}
