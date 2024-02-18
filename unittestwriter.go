package unittestwriter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	// "time"
	// "bytes"


	"github.com/naviscom/dbSchemaReader"
)

func main_testFunc(dirPath string) {
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
	_, _ = outputFile.WriteString(` //_ "github.com/lib/pq"` + "\n")
	_, _ = outputFile.WriteString(")" + "\n")
	_, _ = outputFile.WriteString("\n")

	_, _ = outputFile.WriteString("//const (" + "\n")
	_, _ = outputFile.WriteString(` //dbDriver = "postgres"` + "\n")
	_, _ = outputFile.WriteString(` //dbSource = "postgresql://root:secret@localhost:5432/catalyst?sslmode=disable"` + "\n")
	_, _ = outputFile.WriteString("//)" + "\n")
	_, _ = outputFile.WriteString("\n")

	_, _ = outputFile.WriteString("//var testQueries *Queries" + "\n")
	_, _ = outputFile.WriteString("//var testDB *sql.DB" + "\n")
	_, _ = outputFile.WriteString("var testStore *Store" + "\n")
	_, _ = outputFile.WriteString("\n")

	_, _ = outputFile.WriteString("func TestMain(m *testing.M ) {" + "\n")
	_, _ = outputFile.WriteString("	//var err error" + "\n")
	_, _ = outputFile.WriteString(`	config, err := util.LoadConfig("../..")`+ "\n")
	_, _ = outputFile.WriteString(`	if err != nil {` + "\n")
	_, _ = outputFile.WriteString(`		log.Fatal("cannot load config:", err)` + "\n")
	_, _ = outputFile.WriteString(`	}` + "\n")
	_, _ = outputFile.WriteString("\n")

	_, _ = outputFile.WriteString(`connPool, err := pgxpool.New(context.Background(), config.DBSource)` + "\n")
	_, _ = outputFile.WriteString("	if err != nil {" + "\n")
	_, _ = outputFile.WriteString(`		log.Fatal("cannot connect to db:", err)` + "\n")
	_, _ = outputFile.WriteString("	}" + "\n")
	_, _ = outputFile.WriteString("\n")

	_, _ = outputFile.WriteString(`	testStore = NewStore(connPool)` + "\n")
	_, _ = outputFile.WriteString("	os.Exit(m.Run())" + "\n")
	_, _ = outputFile.WriteString("}" + "\n")
	outputFile.Close()
	fmt.Println("main_test.go file has been generated successfully")
}

func CreateRandomFunction(tableX []dbSchemaReader.Table_Struct, i int, outputFile *os.File) {
	funcSig := "func createRandom" + tableX[i].FunctionSignature + "(t *testing.T"
	// _, _ = outputFile.WriteString("func createRandom"+tableX[i].FunctionSignature+"(t *testing.T")
	for k := 0; k < len(tableX[i].ForeignKeys); k++ {
		funcSig = funcSig + ", " + tableX[i].ForeignKeys[k].FK_Related_SingularTableName + " " + tableX[i].ForeignKeys[k].FK_Related_TableName_Singular_Object
		// _, _ = outputFile.WriteString(", "+ tableX[i].ForeignKeys[k].FK_Related_SingularTableName+" "+ tableX[i].ForeignKeys[k].FK_Related_TableName_Singular_Object)
	}
	funcSig = funcSig + ") " + tableX[i].FunctionSignature
	_, _ = outputFile.WriteString(funcSig + " {" + "\n")
	_, _ = outputFile.WriteString("	arg := Create" + tableX[i].FunctionSignature + "Params{" + "\n")
	for j := 1; j < len(tableX[i].Table_Columns); j++ {
		if tableX[i].Table_Columns[j].ForeignFlag {
			for k := 0; k < len(tableX[i].ForeignKeys); k++ {
				if tableX[i].ForeignKeys[k].FK_Column == tableX[i].Table_Columns[j].Column_name {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    	" + tableX[i].ForeignKeys[k].FK_Related_SingularTableName + "." + strings.ToUpper((tableX[i].ForeignKeys[k].FK_Related_Table_Column)+","+"\n"))
				}
			}
		} else {
			if tableX[i].Table_Columns[j].ColumnType == "varchar" {
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    util.RandomName(8)," + "\n")
			}
			if tableX[i].Table_Columns[j].ColumnType == "bigint" {
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    util.RandomInteger(1, 100)," + "\n")
			}
			if tableX[i].Table_Columns[j].ColumnType == "real" {
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    util.RandomReal(1, 100)," + "\n")
			}
			if tableX[i].Table_Columns[j].ColumnType == "timestamptz" {
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    time.Now().UTC()," + "\n")
			}
		}
	}
	_, _ = outputFile.WriteString("	}" + "\n")
//	_, _ = outputFile.WriteString("	" + tableX[i].OutputFileName + ", err := testQueries.Create" + tableX[i].FunctionSignature + "(context.Background(), arg)" + "\n")
	_, _ = outputFile.WriteString("	" + tableX[i].OutputFileName + ", err := testStore.Create" + tableX[i].FunctionSignature + "(context.Background(), arg)" + "\n")
_, _ = outputFile.WriteString("	require.NoError(t, err)" + "\n")
	_, _ = outputFile.WriteString("	require.NotEmpty(t, " + tableX[i].OutputFileName + ")" + "\n")
	for j := 1; j < len(tableX[i].Table_Columns); j++ {
		if tableX[i].Table_Columns[j].ColumnType == "timestamptz" {
			_, _ = outputFile.WriteString("	require.WithinDuration(t, arg." + tableX[i].Table_Columns[j].ColumnNameParams + ", " + tableX[i].OutputFileName + "." + tableX[i].Table_Columns[j].ColumnNameParams +", time.Second" +")" + "\n")
		}else{
			_, _ = outputFile.WriteString("	require.Equal(t, arg." + tableX[i].Table_Columns[j].ColumnNameParams + ", " + tableX[i].OutputFileName + "." + tableX[i].Table_Columns[j].ColumnNameParams + ")" + "\n")
		}
	}
	_, _ = outputFile.WriteString("	return " + tableX[i].OutputFileName + "\n")
	_, _ = outputFile.WriteString("}" + "\n")
	// fmt.Println("	", funcSig+" has been generated successfully")
	_, _ = outputFile.WriteString("\n")
}

func printTestFuncForCreate(tableX []dbSchemaReader.Table_Struct, i int, fk_HierarchyX []dbSchemaReader.FK_Hierarchy, outputFile *os.File) {
	_, _ = outputFile.WriteString("func TestCreate" + tableX[i].FunctionSignature + "(t *testing.T) {" + "\n")
	for k := 0; k < len(fk_HierarchyX); k++ {
		if fk_HierarchyX[k].TableName == tableX[i].Table_name {
			for l := len(fk_HierarchyX[k].RelatedTablesLevels) - 1; l >= 0; l-- {
				for m := 0; m < len(fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList); m++ {
					_, _ = outputFile.WriteString("	" + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName + " := createRandom" + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName_Singular_Object + "(t")
					for g := 0; g < len(tableX); g++ {
						if tableX[g].Table_name == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {
							for h := 0; h < len(tableX[g].ForeignKeys); h++ {
								_, _ = outputFile.WriteString(", " + tableX[g].ForeignKeys[h].FK_Related_SingularTableName)
							}
							break
						}
					}
					_, _ = outputFile.WriteString(")" + "\n")
				}
			}
		}
	}
	_, _ = outputFile.WriteString("	createRandom" + tableX[i].FunctionSignature + "(t")
	for g := 0; g < len(tableX); g++ {
		if tableX[g].Table_name == tableX[i].Table_name {
			for h := 0; h < len(tableX[g].ForeignKeys); h++ {
				_, _ = outputFile.WriteString(", " + tableX[g].ForeignKeys[h].FK_Related_SingularTableName)
			}
			break
		}
	}
	_, _ = outputFile.WriteString(")" + "\n")
	_, _ = outputFile.WriteString("}" + "\n")
	_, _ = outputFile.WriteString("\n")
}

func printTestFuncForReadGet(tableX []dbSchemaReader.Table_Struct, i int, fk_HierarchyX []dbSchemaReader.FK_Hierarchy, outputFile *os.File) {
	var s int = 0
	for j := 0; j < len(tableX[i].Table_Columns); j++ {
		if tableX[i].Table_Columns[j].PrimaryFlag || tableX[i].Table_Columns[j].UniqueFlag {
			var getByColumnName string = tableX[i].Table_Columns[j].ColumnNameParams
			_, _ = outputFile.WriteString("func TestGet" + tableX[i].FunctionSignature + strconv.Itoa(s) + "(t *testing.T) {" + "\n")
			for k := 0; k < len(fk_HierarchyX); k++ {
				if fk_HierarchyX[k].TableName == tableX[i].Table_name {
					for l := len(fk_HierarchyX[k].RelatedTablesLevels) - 1; l >= 0; l-- {
						for m := 0; m < len(fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList); m++ {
							_, _ = outputFile.WriteString("	" + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName + " := createRandom" + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName_Singular_Object + "(t")
							for g := 0; g < len(tableX); g++ {
								if tableX[g].Table_name == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {
									for h := 0; h < len(tableX[g].ForeignKeys); h++ {
										_, _ = outputFile.WriteString(", " + tableX[g].ForeignKeys[h].FK_Related_SingularTableName)
									}
									break
								}
							}
							_, _ = outputFile.WriteString(")" + "\n")
						}
					}
				}
			}
			_, _ = outputFile.WriteString("	" + tableX[i].OutputFileName + "1 := createRandom" + tableX[i].FunctionSignature + "(t")
			for g := 0; g < len(tableX); g++ {
				if tableX[g].Table_name == tableX[i].Table_name {
					if len(tableX[g].ForeignKeys) > 0 {
						for h := 0; h < len(tableX[g].ForeignKeys); h++ {
							_, _ = outputFile.WriteString(", " + tableX[g].ForeignKeys[h].FK_Related_SingularTableName)
						}
						_, _ = outputFile.WriteString(")" + "\n")
						break
					} else {
						_, _ = outputFile.WriteString(")" + "\n")
					}
				}
			}
			// if j == 0 {
			// 	for g := 0; g < len(tableX); g++ {
			// 		if tableX[g].Table_name == tableX[i].Table_name {
			// 			for h := 0; h < len(tableX[g].Table_Columns); h++ {
			// 				if tableX[g].Table_Columns[h].PrimaryFlag {
			// 					getByColumnName = tableX[g].Table_Columns[h].ColumnNameParams
			// 					break
			// 				}
			// 			}
			// 		}
			// 	}
			// }
			// if j == 1 {
			// 	for g := 0; g < len(tableX); g++ {
			// 		if tableX[g].Table_name == tableX[i].Table_name {
			// 			for h := 0; h < len(tableX[g].Table_Columns); h++ {
			// 				if tableX[g].Table_Columns[h].UniqueFlag {
			// 					getByColumnName = tableX[g].Table_Columns[h].ColumnNameParams
			// 					break
			// 				}
			// 			}
			// 		}
			// 	}
			// }
			_, _ = outputFile.WriteString("	" + tableX[i].OutputFileName + "2, err := testStore.Get" + tableX[i].FunctionSignature + strconv.Itoa(s) + "(context.Background(), " + tableX[i].OutputFileName + "1." + getByColumnName + ")" + "\n")
			_, _ = outputFile.WriteString("	" + "require.NoError(t, err)" + "\n")
			_, _ = outputFile.WriteString("	" + "require.NotEmpty(t, " + tableX[i].OutputFileName + "2)" + "\n")
			_, _ = outputFile.WriteString("\n")
			for h := 0; h < len(tableX[i].Table_Columns); h++ {
				if tableX[i].Table_Columns[h].ColumnType == "timestamptz" {
					_, _ = outputFile.WriteString("	require.WithinDuration(t, " + tableX[i].OutputFileName + "1." + tableX[i].Table_Columns[h].ColumnNameParams + ", " + tableX[i].OutputFileName + "2." + tableX[i].Table_Columns[h].ColumnNameParams + ", time.Second)" + "\n")
				} else {
					_, _ = outputFile.WriteString("	require.Equal(t, " + tableX[i].OutputFileName + "1." + tableX[i].Table_Columns[h].ColumnNameParams + ", " + tableX[i].OutputFileName + "2." + tableX[i].Table_Columns[h].ColumnNameParams + ")" + "\n")

				}
			}
			_, _ = outputFile.WriteString("}" + "\n")
			_, _ = outputFile.WriteString("\n")
		}
		s++
	}
}

func printTestFuncForReadList(tableX []dbSchemaReader.Table_Struct, i int, fk_HierarchyX []dbSchemaReader.FK_Hierarchy, outputFile *os.File) {
	_, _ = outputFile.WriteString("func TestList" + tableX[i].FunctionSignature2 + "(t *testing.T) {" + "\n")
	for k := 0; k < len(fk_HierarchyX); k++ {
		if fk_HierarchyX[k].TableName == tableX[i].Table_name {
			for l := len(fk_HierarchyX[k].RelatedTablesLevels) - 1; l >= 0; l-- {
				for m := 0; m < len(fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList); m++ {
					_, _ = outputFile.WriteString("	" + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName + " := createRandom" + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName_Singular_Object + "(t")
					for g := 0; g < len(tableX); g++ {
						if tableX[g].Table_name == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {
							for h := 0; h < len(tableX[g].ForeignKeys); h++ {
								_, _ = outputFile.WriteString(", " + tableX[g].ForeignKeys[h].FK_Related_SingularTableName)
							}
							break
						}
					}
					_, _ = outputFile.WriteString(")" + "\n")
				}
			}
		}
	}
	_, _ = outputFile.WriteString("	for i := 0; i < 10; i++ {" + "\n")
	_, _ = outputFile.WriteString("		createRandom" + tableX[i].FunctionSignature + "(t")
	for g := 0; g < len(tableX); g++ {
		if tableX[g].Table_name == tableX[i].Table_name {
			for h := 0; h < len(tableX[g].ForeignKeys); h++ {
				_, _ = outputFile.WriteString(", " + tableX[g].ForeignKeys[h].FK_Related_SingularTableName)
			}
			break
		}
	}
	_, _ = outputFile.WriteString(")" + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString("	"+"}" + "\n")

	_, _ = outputFile.WriteString("	arg := List" + tableX[i].FunctionSignature2 + "Params{" + "\n")
	for g := 0; g < len(tableX[i].Table_Columns); g++ {
		if tableX[i].Table_Columns[g].ForeignFlag {
			for r := 0; r < len(tableX[i].ForeignKeys); r++ {
				if tableX[i].ForeignKeys[r].FK_Column == tableX[i].Table_Columns[g].Column_name {
					_, _ = outputFile.WriteString("		"+tableX[i].Table_Columns[g].ColumnNameParams+": "+tableX[i].ForeignKeys[r].FK_Related_SingularTableName+"."+strings.ToUpper(tableX[i].ForeignKeys[r].FK_Related_Table_Column)+","+"\n")
				}
			}
		}
	}
	_, _ = outputFile.WriteString("		Limit:         5,"+"\n")
	_, _ = outputFile.WriteString("		Offset:        5,"+"\n")
	_, _ = outputFile.WriteString("	}" + "\n")

	_, _ = outputFile.WriteString("	" + tableX[i].Table_name + ", err := testStore.List" + tableX[i].FunctionSignature2 + "(context.Background(), " + "arg" + ")" + "\n")
	_, _ = outputFile.WriteString("	" + "require.NoError(t, err)" + "\n")
	_, _ = outputFile.WriteString("	" + "require.Len(t, " + tableX[i].Table_name + ", 5)" + "\n")
	_, _ = outputFile.WriteString("\n")

	_, _ = outputFile.WriteString("	for _, "+tableX[i].OutputFileName+" := range "+tableX[i].Table_name+" {"+"\n")
	_, _ = outputFile.WriteString("		require.NotEmpty(t, "+tableX[i].OutputFileName+")" + "\n")
	str := "		require.True(t, "
	pipeflag :=false
	if len(tableX[i].ForeignKeys) == 1{
		for g := 0; g < len(tableX[i].Table_Columns); g++ {
			if tableX[i].Table_Columns[g].ForeignFlag {
				for r := 0; r < len(tableX[i].ForeignKeys); r++ {
					if tableX[i].ForeignKeys[r].FK_Column == tableX[i].Table_Columns[g].Column_name {
						str = str +"arg."+tableX[i].Table_Columns[g].ColumnNameParams +" == " +tableX[i].OutputFileName+"."+tableX[i].Table_Columns[g].ColumnNameParams
						pipeflag = true
					}
				}
			}				
		}
		str = str + ")"
		_, _ = outputFile.WriteString(str + "\n")
	}
	if len(tableX[i].ForeignKeys) > 1{
		for g := 0; g < len(tableX[i].Table_Columns); g++ {
			if tableX[i].Table_Columns[g].ForeignFlag {
				if pipeflag { str = str + " || "}
				for r := 0; r < len(tableX[i].ForeignKeys); r++ {
					if tableX[i].ForeignKeys[r].FK_Column == tableX[i].Table_Columns[g].Column_name {
						str = str +tableX[i].OutputFileName +"."+tableX[i].Table_Columns[g].ColumnNameParams +" == " +tableX[i].ForeignKeys[r].FK_Related_SingularTableName+"."+strings.ToUpper(tableX[i].ForeignKeys[r].FK_Related_Table_Column)
						pipeflag = true
					}
				}
			}				
		}
		str = str + ")"
		_, _ = outputFile.WriteString(str + "\n")
	}
	// require.True(t, transfer.FromAccountID == account1.ID || transfer.ToAccountID == account1.ID)
	_, _ = outputFile.WriteString("	}" + "\n")
	_, _ = outputFile.WriteString("}" + "\n")
	_, _ = outputFile.WriteString("\n")
}

func printTestFuncForUpdate(tableX []dbSchemaReader.Table_Struct, i int, fk_HierarchyX []dbSchemaReader.FK_Hierarchy, outputFile *os.File) {
	_, _ = outputFile.WriteString("func TestUpdate" + tableX[i].FunctionSignature + "(t *testing.T) {" + "\n")
	for k := 0; k < len(fk_HierarchyX); k++ {
		if fk_HierarchyX[k].TableName == tableX[i].Table_name {
			for l := len(fk_HierarchyX[k].RelatedTablesLevels) - 1; l >= 0; l-- {
				for m := 0; m < len(fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList); m++ {
					_, _ = outputFile.WriteString("	" + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName + " := createRandom" + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName_Singular_Object + "(t")
					for g := 0; g < len(tableX); g++ {
						if tableX[g].Table_name == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {
							for h := 0; h < len(tableX[g].ForeignKeys); h++ {
								_, _ = outputFile.WriteString(", " + tableX[g].ForeignKeys[h].FK_Related_SingularTableName)
							}
							break
						}
					}
					_, _ = outputFile.WriteString(")" + "\n")
				}
			}
		}
	}
	_, _ = outputFile.WriteString("	" + tableX[i].OutputFileName + "1 := createRandom" + tableX[i].FunctionSignature + "(t")
	for g := 0; g < len(tableX); g++ {
		if tableX[g].Table_name == tableX[i].Table_name {
			if len(tableX[g].ForeignKeys) > 0 {
				for h := 0; h < len(tableX[g].ForeignKeys); h++ {
					_, _ = outputFile.WriteString(", " + tableX[g].ForeignKeys[h].FK_Related_SingularTableName)
				}
				_, _ = outputFile.WriteString(")" + "\n")
				break
			} else {
				_, _ = outputFile.WriteString(")" + "\n")
			}
		}
	}
	var getByColumnName string
	for g := 0; g < len(tableX); g++ {
		if tableX[g].Table_name == tableX[i].Table_name {
			for h := 0; h < len(tableX[g].Table_Columns); h++ {
				if tableX[g].Table_Columns[h].PrimaryFlag {
					getByColumnName = tableX[g].Table_Columns[h].ColumnNameParams
					break
				}
			}
		}
	}
	_, _ = outputFile.WriteString("	arg := Update" + tableX[i].FunctionSignature + "Params{" + "\n")
	for p := 0; p < len(tableX[i].Table_Columns); p++ {
		if tableX[i].Table_Columns[p].ForeignFlag {
			for k := 0; k < len(tableX[i].ForeignKeys); k++ {
				if tableX[i].ForeignKeys[k].FK_Column == tableX[i].Table_Columns[p].Column_name {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ":    	" + tableX[i].ForeignKeys[k].FK_Related_SingularTableName + "." + strings.ToUpper((tableX[i].ForeignKeys[k].FK_Related_Table_Column)+","+"\n"))
				}
			}
		} else {
			if tableX[i].Table_Columns[p].ColumnType == "bigserial" {
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ":    "+ tableX[i].OutputFileName + "1." + getByColumnName +"," + "\n")
				continue
			}
			if tableX[i].Table_Columns[p].ColumnType == "varchar" {
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ":    util.RandomName(8)," + "\n")
				continue
			}
			if tableX[i].Table_Columns[p].ColumnType == "bigint" || tableX[i].Table_Columns[p].ColumnType == "bigserial"{
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ":    util.RandomInteger(1, 100)," + "\n")
				continue
			}
			if tableX[i].Table_Columns[p].ColumnType == "real" {
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ":    util.RandomReal(1, 100)," + "\n")
				continue
			}
			if tableX[i].Table_Columns[p].ColumnType == "timestamptz" {
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ":    time.Now().UTC()," + "\n")
				continue
			}
		}
	}
	_, _ = outputFile.WriteString("	}" + "\n")

	_, _ = outputFile.WriteString("	" + tableX[i].OutputFileName + "2, err := testStore.Update" + tableX[i].FunctionSignature + "(context.Background(), " + "arg" + ")" + "\n")
	_, _ = outputFile.WriteString("	" + "require.NoError(t, err)" + "\n")
	_, _ = outputFile.WriteString("	" + "require.NotEmpty(t, " + tableX[i].OutputFileName + "2)" + "\n")
	_, _ = outputFile.WriteString("\n")
	for h := 0; h < len(tableX[i].Table_Columns); h++ {
		if tableX[i].Table_Columns[h].PrimaryFlag || tableX[i].Table_Columns[h].ForeignFlag {
			_, _ = outputFile.WriteString("	require.Equal(t, " + tableX[i].OutputFileName + "1." + tableX[i].Table_Columns[h].ColumnNameParams + ", " + tableX[i].OutputFileName + "2." + tableX[i].Table_Columns[h].ColumnNameParams + ")" + "\n")
		}else {
			if tableX[i].Table_Columns[h].ColumnType == "timestamptz" {
				_, _ = outputFile.WriteString("	require.WithinDuration(t, " + "arg." + tableX[i].Table_Columns[h].ColumnNameParams + ", " + tableX[i].OutputFileName + "2." + tableX[i].Table_Columns[h].ColumnNameParams + ", time.Second)" + "\n")
			} else {
				_, _ = outputFile.WriteString("	require.Equal(t, " + "arg." + tableX[i].Table_Columns[h].ColumnNameParams + ", " + tableX[i].OutputFileName + "2." + tableX[i].Table_Columns[h].ColumnNameParams + ")" + "\n")
			}	
		}
	}
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString("}" + "\n")
	_, _ = outputFile.WriteString("\n")
}

func printTestFuncForDelete(tableX []dbSchemaReader.Table_Struct, i int, fk_HierarchyX []dbSchemaReader.FK_Hierarchy, outputFile *os.File) {
	_, _ = outputFile.WriteString("func TestDelete" + tableX[i].FunctionSignature + "(t *testing.T) {" + "\n")
	for k := 0; k < len(fk_HierarchyX); k++ {
		if fk_HierarchyX[k].TableName == tableX[i].Table_name {
			for l := len(fk_HierarchyX[k].RelatedTablesLevels) - 1; l >= 0; l-- {
				for m := 0; m < len(fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList); m++ {
					_, _ = outputFile.WriteString("	" + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName + " := createRandom" + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName_Singular_Object + "(t")
					for g := 0; g < len(tableX); g++ {
						if tableX[g].Table_name == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {
							for h := 0; h < len(tableX[g].ForeignKeys); h++ {
								_, _ = outputFile.WriteString(", " + tableX[g].ForeignKeys[h].FK_Related_SingularTableName)
							}
							break
						}
					}
					_, _ = outputFile.WriteString(")" + "\n")
				}
			}
		}
	}
	_, _ = outputFile.WriteString("	" + tableX[i].OutputFileName + "1 := createRandom" + tableX[i].FunctionSignature + "(t")
	for g := 0; g < len(tableX); g++ {
		if tableX[g].Table_name == tableX[i].Table_name {
			if len(tableX[g].ForeignKeys) > 0 {
				for h := 0; h < len(tableX[g].ForeignKeys); h++ {
					_, _ = outputFile.WriteString(", " + tableX[g].ForeignKeys[h].FK_Related_SingularTableName)
				}
				_, _ = outputFile.WriteString(")" + "\n")
				break
			} else {
				_, _ = outputFile.WriteString(")" + "\n")
			}
		}
	}
	var getByColumnName string
	for g := 0; g < len(tableX); g++ {
		if tableX[g].Table_name == tableX[i].Table_name {
			for h := 0; h < len(tableX[g].Table_Columns); h++ {
				if tableX[g].Table_Columns[h].PrimaryFlag {
					getByColumnName = tableX[g].Table_Columns[h].ColumnNameParams
					break
				}
			}
		}
	}
	_, _ = outputFile.WriteString("	" + "err := testStore.Delete" + tableX[i].FunctionSignature + "(context.Background(), " + tableX[i].OutputFileName + "1." + getByColumnName + ")" + "\n")
	_, _ = outputFile.WriteString("	" + "require.NoError(t, err)" + "\n")
	_, _ = outputFile.WriteString("	" + tableX[i].OutputFileName + "2, err := testStore.Get" + tableX[i].FunctionSignature + "0" + "(context.Background(), " + tableX[i].OutputFileName + "1." + getByColumnName + ")" + "\n")
	_, _ = outputFile.WriteString("	" + "require.Error(t, err)" + "\n")
	_, _ = outputFile.WriteString("	" + "require.EqualError(t, err, ErrRecordNotFound.Error())" + "\n")
	_, _ = outputFile.WriteString("	" + "require.Empty(t, "+ tableX[i].OutputFileName +"2)"+ "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString("}" + "\n")
	_, _ = outputFile.WriteString("\n")
}

func TestWriter(projectFolderPath string) {
	//generating main_test.go
	dirPath := projectFolderPath
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
	/////////////////////////////////////////////////
	//generate unit tests for go file in sqlc folder
	/////////////////////////////////////////////////
	pathToSearch := dirPath + "/db/migration"
	err := filepath.Walk(pathToSearch, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	var tableX []dbSchemaReader.Table_Struct
	var fk_HierarchyX []dbSchemaReader.FK_Hierarchy
	for _, element := range files {
		if element[len(element)-6:] == `up.sql` {
			tableX, fk_HierarchyX = dbSchemaReader.ReadSchema(element)
			for i := 0; i < len(tableX); i++ {
				outputFile, errs := os.Create(dirPath + "/db/sqlc/" + tableX[i].OutputFileName + "_test.go")
				if errs != nil {
					fmt.Println("Failed to create file:", errs)
					return
				}
				defer outputFile.Close()
				fmt.Println("generating ", tableX[i].OutputFileName+"_test.go")
				_, _ = outputFile.WriteString("package db" + "\n")
				_, _ = outputFile.WriteString("\n")
				_, _ = outputFile.WriteString("import (" + "\n")
				_, _ = outputFile.WriteString(`	"context"` + "\n")
				_, _ = outputFile.WriteString(`	"time"` + "\n")
				_, _ = outputFile.WriteString(`	"testing"` + "\n")
				_, _ = outputFile.WriteString(`	"github.com/stretchr/testify/require"` + "\n")
				_, _ = outputFile.WriteString(`	"github.com/naviscom/` + projectFolderName + `/util"` + "\n")
				_, _ = outputFile.WriteString(")" + "\n")
				_, _ = outputFile.WriteString("\n")
				CreateRandomFunction(tableX[:], i, outputFile)
				printTestFuncForCreate(tableX[:], i, fk_HierarchyX[:], outputFile)
				printTestFuncForReadGet(tableX[:], i, fk_HierarchyX[:], outputFile)
				printTestFuncForReadList(tableX[:], i, fk_HierarchyX[:], outputFile)
				printTestFuncForUpdate(tableX[:], i, fk_HierarchyX[:], outputFile)
				printTestFuncForDelete(tableX[:], i, fk_HierarchyX[:], outputFile)
				outputFile.Close()
			}
		}
	}

	//writing error.go
	outputFileName := dirPath + "/db/sqlc/error.go"
	outputFile, errs := os.Create(outputFileName)
	if errs != nil {
		fmt.Println("Failed to create file:", errs)
		return
	}
	defer outputFile.Close()
	_, _ = outputFile.WriteString("package db" + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString("import (" + "\n")
	_, _ = outputFile.WriteString(`	//"database/sql"` + "\n")
	_, _ = outputFile.WriteString(`	"github.com/jackc/pgx/v5"` + "\n")
	_, _ = outputFile.WriteString(`	//"github.com/jackc/pgx/v5/pgconn"` + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString(`)` + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString(`//const (` + "\n")
	_, _ = outputFile.WriteString(`	//ForeignKeyViolation = "23503"` + "\n")
	_, _ = outputFile.WriteString(`	//UniqueViolation     = "23505"` + "\n")
	_, _ = outputFile.WriteString(`//)` + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString(`var ErrRecordNotFound = pgx.ErrNoRows` + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString(`//var ErrUniqueViolation = &pgconn.PgError{` + "\n")
	_, _ = outputFile.WriteString(`	//Code: UniqueViolation,` + "\n")
	_, _ = outputFile.WriteString(`//}` + "\n")
	_, _ = outputFile.WriteString(`//func ErrorCode(err error) string {` + "\n")
	_, _ = outputFile.WriteString(`	//var pgErr *pgconn.PgError` + "\n")
	_, _ = outputFile.WriteString(`	//if errors.As(err, &pgErr) {`+"\n")
	_, _ = outputFile.WriteString(`		//return pgErr.Code` + "\n")
	_, _ = outputFile.WriteString(`	//}` + "\n")
	_, _ = outputFile.WriteString(`	//return ""` + "\n")
	_, _ = outputFile.WriteString(`//}` + "\n")

	//Executing goimports
	cmd := exec.Command("goimports", "-w", ".")
	cmd.Dir = dirPath+"/db/sqlc"
	cmd.Run()
	println("goimports executed successfully")
	// var stderr bytes.Buffer

	//Executing go mod tidy
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = dirPath
	cmd.Run()
	println("go mod tidy executed successfully")

}
