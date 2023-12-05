package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/naviscom/dbSchemaReader"
)

type FK_Hierarchy struct{
	TableName				string
	RelatedTablesLevels		[]RelatedTables
}

type RelatedTables struct {
	RelatedTableList	[]RelatedTable
}

type RelatedTable struct {
	FK_Related_TableName					string
	FK_Related_SingularTableName			string
	FK_Related_Table_Column					string
	FK_Related_TableName_Plural_Object		string
	FK_Related_TableName_Singular_Object	string
}


func main_testFunc(dirPath string){
	outputFileName := dirPath + "/db/sqlc/main_test.go"
	outputFile, errs := os.Create(outputFileName)
	if errs != nil {
		fmt.Println("Failed to create file:", errs)
		return
	}
	defer outputFile.Close()
	_, _ = outputFile.WriteString("package db" + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString("import (" + "\n")
	_, _ = outputFile.WriteString(` "database/sql"` + "\n")
	_, _ = outputFile.WriteString(` "log"` + "\n")
	_, _ = outputFile.WriteString(` "os"` + "\n")
	_, _ = outputFile.WriteString(` "testing"` + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString(` _ "github.com/lib/pq"` + "\n")
	_, _ = outputFile.WriteString(")" + "\n")
	_, _ = outputFile.WriteString("\n")

	_, _ = outputFile.WriteString("var testQueries *Queries"+ "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString("const ("+ "\n")
	_, _ = outputFile.WriteString(` dbDriver = "postgres"`+ "\n")
	_, _ = outputFile.WriteString(` dbSource = "postgresql://root:secret@localhost:5432/catalyst?sslmode=disable"`+ "\n")
	_, _ = outputFile.WriteString(")"+"\n")
	_, _ = outputFile.WriteString("\n")

	_, _ = outputFile.WriteString("func TestMain(m *testing.M) {" + "\n")
	_, _ = outputFile.WriteString("	conn, err := sql.Open(dbDriver, dbSource)" + "\n")
	_, _ = outputFile.WriteString("	if err != nil {" + "\n")
	_, _ = outputFile.WriteString(`		log.Fatal("cannot connect to db:", err)` + "\n")
	_, _ = outputFile.WriteString("	}" + "\n")
	_, _ = outputFile.WriteString("\n")

	_, _ = outputFile.WriteString("	testQueries = New(conn)" + "\n")
	_, _ = outputFile.WriteString("	os.Exit(m.Run())" + "\n")
	_, _ = outputFile.WriteString("}" + "\n")
	outputFile.Close()
	fmt.Println("main_test.go file has been generated successfully")
}

func printTestFuncForCreate(tableX[]dbSchemaReader.Table_Struct, i int, outputFile *os.File){
	var relatedTable RelatedTable
	var relatedTables RelatedTables
	var fk_Hierarchy FK_Hierarchy
	var fk_HierarchyX []FK_Hierarchy
	if len(tableX[i].ForeignKeys) > 0 {
		for j :=0; j < len(tableX[i].ForeignKeys); j++{
			relatedTable.FK_Related_TableName = tableX[i].ForeignKeys[j].FK_Related_TableName
			relatedTable.FK_Related_SingularTableName = tableX[i].ForeignKeys[j].FK_Related_SingularTableName
			relatedTable.FK_Related_Table_Column = tableX[i].ForeignKeys[j].FK_Related_Table_Column
			relatedTable.FK_Related_TableName_Singular_Object = tableX[i].ForeignKeys[j].FK_Related_TableName_Singular_Object
			relatedTable.FK_Related_TableName_Plural_Object = tableX[i].ForeignKeys[j].FK_Related_TableName_Plural_Object
			relatedTables.RelatedTableList = append(relatedTables.RelatedTableList, relatedTable)
		}
		fk_Hierarchy.TableName = tableX[i].Table_name
		fk_Hierarchy.RelatedTablesLevels = append(fk_Hierarchy.RelatedTablesLevels, relatedTables)
		fk_HierarchyX = append(fk_HierarchyX, fk_Hierarchy)
	}else{
		fk_Hierarchy.TableName = tableX[i].Table_name
		fk_Hierarchy.RelatedTablesLevels = append(fk_Hierarchy.RelatedTablesLevels, relatedTables)
		fk_HierarchyX = append(fk_HierarchyX, fk_Hierarchy)
	}
	var c int
	var d int
	var e int
	c = 0
	for k :=0; k < len(fk_HierarchyX); k++{
		d = len(fk_HierarchyX[k].RelatedTablesLevels)
		e = d - c
		if fk_HierarchyX[k].TableName == tableX[i].Table_name {
			for l :=len(fk_HierarchyX[k].RelatedTablesLevels)-e; l < len(fk_HierarchyX[k].RelatedTablesLevels); l++{
				for m:=0; m < len(fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList); m++{
					for z :=0; z < len(tableX); z++{
						if fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName == tableX[z].Table_name {	
							if len(tableX[z].ForeignKeys) > 0 {
								relatedTables.RelatedTableList = nil
								for y :=0; y < len(tableX[z].ForeignKeys); y++{
									relatedTable.FK_Related_TableName = tableX[z].ForeignKeys[y].FK_Related_TableName
									relatedTable.FK_Related_SingularTableName = tableX[z].ForeignKeys[y].FK_Related_SingularTableName
									relatedTable.FK_Related_Table_Column = tableX[z].ForeignKeys[y].FK_Related_Table_Column
									relatedTable.FK_Related_TableName_Singular_Object = tableX[z].ForeignKeys[y].FK_Related_TableName_Singular_Object
									relatedTable.FK_Related_TableName_Plural_Object = tableX[z].ForeignKeys[y].FK_Related_TableName_Plural_Object
						
									relatedTables.RelatedTableList = append(relatedTables.RelatedTableList, relatedTable)							
								}
								fk_HierarchyX[k].RelatedTablesLevels = append(fk_HierarchyX[k].RelatedTablesLevels, relatedTables)
							}
						}
					}
				}

			}
			c = d
		}
	}
	_, _ = outputFile.WriteString("func TestCreate"+tableX[i].FunctionSignature+"(t *testing.T) {"+"\n")
	for k := 0; k < len(fk_HierarchyX); k++{
		if fk_HierarchyX[k].TableName == tableX[i].Table_name {
			for l := len(fk_HierarchyX[k].RelatedTablesLevels)-1; l >=0; l--{
				for m := 0; m < len(fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList); m++{
					_, _ = outputFile.WriteString("	"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName+" := createRandom"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName_Singular_Object+"(t")					
					for g := 0; g < len(tableX); g++ {
						if 	tableX[g].Table_name == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {
							for h := 0; h < len(tableX[g].ForeignKeys); h++ {
								_, _ = outputFile.WriteString(", "+ tableX[g].ForeignKeys[h].FK_Related_SingularTableName)					
							}
							break
						}
					}
					_, _ = outputFile.WriteString(")"+"\n")
				}
			}
		}
	}
	_, _ = outputFile.WriteString("	createRandom"+tableX[i].FunctionSignature+"(t")
	for g := 0; g < len(tableX); g++ {
		if 	tableX[g].Table_name == tableX[i].Table_name {
			for h := 0; h < len(tableX[g].ForeignKeys); h++ {
				_, _ = outputFile.WriteString(", "+ tableX[g].ForeignKeys[h].FK_Related_SingularTableName)					
			}
			break
		}
	}
	_, _ = outputFile.WriteString(")"+"\n")
	_, _ = outputFile.WriteString("}"+"\n")
	_, _ = outputFile.WriteString("\n")
	fmt.Println("	","func TestCreate"+tableX[i].FunctionSignature+"(t *testing.T) has been generated successfully")
}

func printTestFuncForReadGet(table []dbSchemaReader.Table_Struct, i int, outputFile *os.File){
	for j := 0; j < len(table[i].Table_Columns); j++ {
		if table[i].Table_Columns[j].PrimaryFlag || table[i].Table_Columns[j].UniqueFlag {
			_, _ = outputFile.WriteString("func TestGet"+table[i].FunctionSignature+strconv.Itoa(j)+"(t *testing.T) {"+"\n")
			_, _ = outputFile.WriteString("\n")
			_, _ = outputFile.WriteString("}"+"\n")
			_, _ = outputFile.WriteString("\n")			
		}
	}	
}

func printTestFuncForReadList(table []dbSchemaReader.Table_Struct, i int, outputFile *os.File){
	_, _ = outputFile.WriteString("func TestList"+table[i].FunctionSignature2+"(t *testing.T) {"+"\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString("}"+"\n")
	_, _ = outputFile.WriteString("\n")
}

func printTestFuncForUpdate(table []dbSchemaReader.Table_Struct, i int, outputFile *os.File){
	_, _ = outputFile.WriteString("func TestUpdate"+table[i].FunctionSignature+"(t *testing.T) {"+"\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString("}"+"\n")
	_, _ = outputFile.WriteString("\n")
}

func printTestFuncForDelete(table []dbSchemaReader.Table_Struct, i int, outputFile *os.File){
	_, _ = outputFile.WriteString("func TestDelete"+table[i].FunctionSignature+"(t *testing.T) {"+"\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString("}"+"\n")
	_, _ = outputFile.WriteString("\n")
}

func main() {
	//generating main_test.go
	dirPath := os.Args[1]
	var files []string
	var x int
	var projectFolderName string //, projectFolderPath string
	var namewithpath bool
	namewithpath = false
	temp := strings.Split(dirPath, "")
	for x = len(temp) - 1; x >= 0; x-- {
		if temp[x] == `/` {
			namewithpath = true
			break
		}
	}
	if namewithpath {
		projectFolderName = strings.Join(temp[x+1:], "")
		// projectFolderPath = strings.Join(temp[:x], "")
	} else {
		projectFolderName = dirPath
	}
	/////////////////////////////
	//generate main_test.go file
	/////////////////////////////
	main_testFunc(dirPath)

	
	//generate unit tests for go file in sqlc folder
	pathToSearch := dirPath+"/db/migration"
	err := filepath.Walk(pathToSearch, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	var tableX []dbSchemaReader.Table_Struct
	for _, element := range files {
		if element[len(element)-6:] == `up.sql` {
			tableX = dbSchemaReader.ReadSchema(element)				
			for i:=0; i<len(tableX); i++{
				outputFile, errs := os.Create(dirPath+"/db/sqlc/"+tableX[i].OutputFileName+"_test.go")
				if errs != nil {
				  fmt.Println("Failed to create file:", errs)
				  return
				}
				defer outputFile.Close()
				fmt.Println("generating ",tableX[i].OutputFileName+"_test.go")

				_, _ = outputFile.WriteString("package db" + "\n")
				_, _ = outputFile.WriteString("\n")
				_, _ = outputFile.WriteString("import ("+"\n")
				_, _ = outputFile.WriteString(`	"context"` + "\n")
				_, _ = outputFile.WriteString(`	// "time"` + "\n")
				_, _ = outputFile.WriteString(`	"testing"` + "\n")
				_, _ = outputFile.WriteString(`	"github.com/stretchr/testify/require"` + "\n")
				_, _ = outputFile.WriteString(`	"github.com/naviscom/`+projectFolderName+`/util"` + "\n")
				_, _ = outputFile.WriteString(")" + "\n")
				_, _ = outputFile.WriteString("\n")
				//////////////////////////////////////////////
				//generating  createRandom func for the table
				//////////////////////////////////////////////
				funcSig := "func createRandom"+tableX[i].FunctionSignature+"(t *testing.T"
				// _, _ = outputFile.WriteString("func createRandom"+tableX[i].FunctionSignature+"(t *testing.T")
				for k := 0; k < len(tableX[i].ForeignKeys); k++ {
					funcSig = funcSig + ", "+ tableX[i].ForeignKeys[k].FK_Related_SingularTableName+" "+ tableX[i].ForeignKeys[k].FK_Related_TableName_Singular_Object
					// _, _ = outputFile.WriteString(", "+ tableX[i].ForeignKeys[k].FK_Related_SingularTableName+" "+ tableX[i].ForeignKeys[k].FK_Related_TableName_Singular_Object)					
				}
				funcSig = funcSig + ") "+tableX[i].FunctionSignature
				_, _ = outputFile.WriteString(funcSig+" {" + "\n")
				_, _ = outputFile.WriteString("	arg := Create"+tableX[i].FunctionSignature+"Params{" + "\n")
				for j := 1; j < len(tableX[i].Table_Columns); j++ {
					if 	tableX[i].Table_Columns[j].ForeignFlag {
						for k := 0; k < len(tableX[i].ForeignKeys); k++ {
							if 	tableX[i].ForeignKeys[k].FK_Column == tableX[i].Table_Columns[j].Column_name {
								_, _ = outputFile.WriteString("		"+tableX[i].Table_Columns[j].ColumnNameParams+":    	"+ tableX[i].ForeignKeys[k].FK_Related_SingularTableName+"."+strings.ToUpper((tableX[i].ForeignKeys[k].FK_Related_Table_Column)+","+ "\n"))
							}
						}
					} else {
						if tableX[i].Table_Columns[j].ColumnType == "varchar" {
							_, _ = outputFile.WriteString("		"+tableX[i].Table_Columns[j].ColumnNameParams+":    util.RandomName(8)," + "\n")
						}
						if tableX[i].Table_Columns[j].ColumnType == "bigint" {
							_, _ = outputFile.WriteString("		"+tableX[i].Table_Columns[j].ColumnNameParams+":    util.RandomInteger(1, 100)," + "\n")
						}
						if tableX[i].Table_Columns[j].ColumnType == "real" {
							_, _ = outputFile.WriteString("		"+tableX[i].Table_Columns[j].ColumnNameParams+":    util.RandomReal(1, 100)," + "\n")
						}
						if tableX[i].Table_Columns[j].ColumnType == "timestamptz" {
							_, _ = outputFile.WriteString("		"+tableX[i].Table_Columns[j].ColumnNameParams+":    time.Now().UTC(),"+ "\n")
						}	
					}		
				}
				_, _ = outputFile.WriteString("	}" + "\n")
				_, _ = outputFile.WriteString("	"+tableX[i].OutputFileName+", err := testQueries.Create"+tableX[i].FunctionSignature+"(context.Background(), arg)" + "\n")
				_, _ = outputFile.WriteString("	require.NoError(t, err)" + "\n")
				_, _ = outputFile.WriteString("	require.NotEmpty(t, "+tableX[i].OutputFileName+")" + "\n")
				for j := 1; j < len(tableX[i].Table_Columns); j++ {
					_, _ = outputFile.WriteString("	require.Equal(t, arg."+tableX[i].Table_Columns[j].ColumnNameParams+", "+ tableX[i].OutputFileName+"."+tableX[i].Table_Columns[j].ColumnNameParams+")"+ "\n")				
				}
				_, _ = outputFile.WriteString("	return "+tableX[i].OutputFileName + "\n")
				_, _ = outputFile.WriteString("}" + "\n")
				fmt.Println("	",funcSig+ " has been generated successfully")

				_, _ = outputFile.WriteString("\n")
				printTestFuncForCreate(tableX[:], i, outputFile)
				printTestFuncForReadGet(tableX[:], i, outputFile)
				printTestFuncForReadList(tableX[:], i, outputFile)
				printTestFuncForUpdate(tableX[:], i, outputFile)
				printTestFuncForDelete(tableX[:], i, outputFile)
				outputFile.Close()
			  }
		}
	}
}