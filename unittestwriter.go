package unittestwriter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	// "time"
	// "bytes"

	"github.com/naviscom/dbschemareader"
)

func ToCamelCase(input string) string {
	parts := strings.Split(input, "_")
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			runes := []rune(parts[i])
			runes[0] = unicode.ToUpper(runes[0])
			parts[i] = string(runes)
		}
	}
	return strings.Join(parts, "")
}

func FormatFieldName(input string) string {
	fmt.Println("input: ", input)
	parts := strings.Split(input, "_")
	fmt.Println("parts just after split: ")
	for i := 0; i < len(parts); i++ {
		word := strings.ToLower(parts[i])
		fmt.Println("word: ", word)
		if word == "id" {
			parts[i] = "ID"
		} else if len(word) > 0 {
			parts[i] = strings.ToUpper(word[:1]) + word[1:]
			fmt.Println("parts: ", parts)
		}
	}
	// Ensure first letter of the final result is capitalized
	result := strings.Join(parts, "")
	if len(result) > 0 {
		result = strings.ToUpper(result[:1]) + result[1:]
	}
	fmt.Println("result: ", result)
	return result
}

func main_testFunc(projectFolderName string, gitHubAccountName string, dirPath string) {
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
	_, _ = outputFile.WriteString(` "github.com/jackc/pgx/v5/pgxpool"` + "\n")
	_, _ = outputFile.WriteString(` "github.com/`+gitHubAccountName+`/`+projectFolderName+`/util"` + "\n")
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
	_, _ = outputFile.WriteString(`	config, err := util.LoadConfig("../..")` + "\n")
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

func CreateRandomFunction(tableX []dbschemareader.Table_Struct, i int, outputFile *os.File) {
	funcSig := "func createRandom" + tableX[i].FunctionSignature + "(t *testing.T"
	for k := 0; k < len(tableX[i].ForeignKeys); k++ {
		CamelCase := ToCamelCase(tableX[i].ForeignKeys[k].FK_Related_TableName_Singular_Object)
		funcSig = funcSig + ", " + tableX[i].ForeignKeys[k].FK_Related_SingularTableName + " " + CamelCase
		// funcSig = funcSig + ", " + tableX[i].ForeignKeys[k].FK_Related_SingularTableName + " " + tableX[i].ForeignKeys[k].FK_Related_TableName_Singular_Object
	}
	CamelCase := ToCamelCase(tableX[i].FunctionSignature)
	funcSig = funcSig + ") " + CamelCase
	// funcSig = funcSig + ") " + tableX[i].FunctionSignature
	_, _ = outputFile.WriteString(funcSig + " {" + "\n")

	if tableX[i].Table_name == "users" {
		_, _ = outputFile.WriteString("	hashedPassword, err := util.HashPassword(util.RandomString(6))" + "\n")
		_, _ = outputFile.WriteString("	require.NoError(t, err)" + "\n")

	}
	_, _ = outputFile.WriteString("	arg := Create" + tableX[i].FunctionSignature + "Params{" + "\n")
	var z int
	if tableX[i].Table_Columns[0].PrimaryFlag && (tableX[i].Table_Columns[0].ColumnType == "bigserial" || tableX[i].Table_Columns[0].ColumnType == "uuid") {
		z = 1
	}
	if tableX[i].Table_Columns[0].PrimaryFlag && tableX[i].Table_Columns[0].ColumnType != "bigserial" && tableX[i].Table_Columns[0].ColumnType != "uuid" {
		z = 0
	}
	for j := z; j < len(tableX[i].Table_Columns); j++ {
		if tableX[i].Table_Columns[j].ForeignFlag {
			for k := 0; k < len(tableX[i].ForeignKeys); k++ {
				if tableX[i].ForeignKeys[k].FK_Column == tableX[i].Table_Columns[j].Column_name {
					fmt.Println("tableX[i].Table_Columns[j].Column_name : ", tableX[i].Table_Columns[j].Column_name)					
					FormatedFieldName := FormatFieldName(tableX[i].ForeignKeys[k].FK_Related_Table_Column)
					fmt.Println("FormatedFieldName : ", FormatedFieldName)
					fmt.Println(FormatedFieldName)
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    	" + tableX[i].ForeignKeys[k].FK_Related_SingularTableName + "." + FormatedFieldName+","+"\n")
					// _, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    	" + tableX[i].ForeignKeys[k].FK_Related_SingularTableName + "." + strings.ToUpper((tableX[i].ForeignKeys[k].FK_Related_Table_Column)+","+"\n"))
				}
			}
		} else {
			if tableX[i].Table_name == "users" && (tableX[i].Table_Columns[j].Column_name == "password_changed_at" || tableX[i].Table_Columns[j].Column_name == "password_created_at") {
				continue
			}
			if tableX[i].Table_Columns[j].ColumnType == "varchar" {
				// fmt.Println("tableX[i].Table_name , tableX[i].Table_Columns[j].Column_name", tableX[i].Table_name, tableX[i].Table_Columns[j].Column_name)
				if tableX[i].Table_name == "users" && tableX[i].Table_Columns[j].Column_name == "email" {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    util.RandomEmail()," + "\n")
					continue
				} else if tableX[i].Table_name == "users" && tableX[i].Table_Columns[j].Column_name == "hashed_password" {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    " + "hashedPassword" + "," + "\n")
				} else {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    util.RandomName(8)," + "\n")
				}
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
	for j := z; j < len(tableX[i].Table_Columns); j++ {
		if tableX[i].Table_name == "users" && (tableX[i].Table_Columns[j].Column_name == "password_changed_at" || tableX[i].Table_Columns[j].Column_name == "password_created_at") {
			if tableX[i].Table_name == "users" && tableX[i].Table_Columns[j].Column_name == "password_changed_at" {
				_, _ = outputFile.WriteString("	require.True(t, " + tableX[i].OutputFileName + "." + tableX[i].Table_Columns[j].ColumnNameParams + ".IsZero())" + "\n")
			}
			if tableX[i].Table_name == "users" && tableX[i].Table_Columns[j].Column_name == "password_created_at" {
				_, _ = outputFile.WriteString("	require.NotZero(t, " + tableX[i].OutputFileName + "." + tableX[i].Table_Columns[j].ColumnNameParams + ")" + "\n")
			}
			continue
		}
		if tableX[i].Table_Columns[j].ColumnType == "timestamptz" {
			_, _ = outputFile.WriteString("	require.WithinDuration(t, arg." + tableX[i].Table_Columns[j].ColumnNameParams + ", " + tableX[i].OutputFileName + "." + tableX[i].Table_Columns[j].ColumnNameParams + ", time.Second" + ")" + "\n")
		} else {
			_, _ = outputFile.WriteString("	require.Equal(t, arg." + tableX[i].Table_Columns[j].ColumnNameParams + ", " + tableX[i].OutputFileName + "." + tableX[i].Table_Columns[j].ColumnNameParams + ")" + "\n")
		}
	}

	_, _ = outputFile.WriteString("	return " + tableX[i].OutputFileName + "\n")
	_, _ = outputFile.WriteString("}" + "\n")
	// fmt.Println("	", funcSig+" has been generated successfully")
	_, _ = outputFile.WriteString("\n")
}

func printTestFuncForCreate(tableX []dbschemareader.Table_Struct, i int, fk_HierarchyX []dbschemareader.FK_Hierarchy, outputFile *os.File) {
	var fkVarMap = make(map[string]string)
	_, _ = outputFile.WriteString("func TestCreate" + tableX[i].FunctionSignature + "(t *testing.T) {" + "\n")
	for k := 0; k < len(fk_HierarchyX); k++ {
		if fk_HierarchyX[k].TableName == tableX[i].Table_name {
			for l := len(fk_HierarchyX[k].RelatedTablesLevels) - 1; l >= 0; l-- {
				for m := 0; m < len(fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList); m++ {
					varName := fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName+"_fk_"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName+strconv.Itoa(k) + strconv.Itoa(l) + strconv.Itoa(m)
					key := fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName+"_fk_"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName
					// Only store if key doesn't exist
					if _, exists := fkVarMap[key]; !exists {
						fkVarMap[key] = varName
						_, _ = outputFile.WriteString("	" + varName+ " := createRandom" + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName_Singular_Object + "(t")
					} else{
						continue
					}
					// _, _ = outputFile.WriteString("	" + fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName+"_fk_"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName+ " := createRandom" + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName_Singular_Object + "(t")
					for g := 0; g < len(fk_HierarchyX); g++ {
						if fk_HierarchyX[g].TableName == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {
							for h := 0; h < len(fk_HierarchyX[g].RelatedTablesLevels); h++ {
								for z := 0; z < len(fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList); z++ {				
									key := fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName
									if val, ok := fkVarMap[key]; ok {
										_, _ = outputFile.WriteString(", " + val)
									}
									// _, _ = outputFile.WriteString(", " + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName)
								}
								if h == 0 {
									break
								}
							}
						}
					}
					_, _ = outputFile.WriteString(")" + "\n")
				}
			}
		}
	}
	_, _ = outputFile.WriteString("	createRandom" + tableX[i].FunctionSignature + "(t")
	for g := 0; g < len(fk_HierarchyX); g++ {
		if fk_HierarchyX[g].TableName == tableX[i].Table_name {
			for h := 0; h < len(fk_HierarchyX[g].RelatedTablesLevels); h++ {
				for z := 0; z < len(fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList); z++ {
					key := fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName
					if val, ok := fkVarMap[key]; ok {
						_, _ = outputFile.WriteString(", " + val)
					}
					// _, _ = outputFile.WriteString(", " + tableX[i].Table_name+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName)
				}
				if h == 0 {
					break
				}
			}
		}
	}
	_, _ = outputFile.WriteString(")" + "\n")
	_, _ = outputFile.WriteString("}" + "\n")
	_, _ = outputFile.WriteString("\n")
}

func printTestFuncForReadGet(tableX []dbschemareader.Table_Struct, i int, fk_HierarchyX []dbschemareader.FK_Hierarchy, outputFile *os.File) {
	var s int = 0
	for j := 0; j < len(tableX[i].Table_Columns); j++ {
		var fkVarMap = make(map[string]string)
		if tableX[i].Table_Columns[j].PrimaryFlag || tableX[i].Table_Columns[j].UniqueFlag {
			var getByColumnName string = tableX[i].Table_Columns[j].ColumnNameParams
			_, _ = outputFile.WriteString("func TestGet" + tableX[i].FunctionSignature + strconv.Itoa(s) + "(t *testing.T) {" + "\n")
			for k := 0; k < len(fk_HierarchyX); k++ {
				if fk_HierarchyX[k].TableName == tableX[i].Table_name {
					for l := len(fk_HierarchyX[k].RelatedTablesLevels) - 1; l >= 0; l-- {
						for m := 0; m < len(fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList); m++ {
							varName := fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName+"_fk_"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName+strconv.Itoa(k) + strconv.Itoa(l) + strconv.Itoa(m)
							key := fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName+"_fk_"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName
							// Only store if key doesn't exist
							if _, exists := fkVarMap[key]; !exists {
								fkVarMap[key] = varName
								_, _ = outputFile.WriteString("	" + varName+ " := createRandom" + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName_Singular_Object + "(t")		
							} else{
								continue
							}
							// _, _ = outputFile.WriteString("	" + fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName+"_fk_"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName + " := createRandom" + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName_Singular_Object + "(t")
							for g := 0; g < len(fk_HierarchyX); g++ {
								if fk_HierarchyX[g].TableName == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {
									for h := 0; h < len(fk_HierarchyX[g].RelatedTablesLevels); h++ {
										for z := 0; z < len(fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList); z++ {
											key := fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName
											if val, ok := fkVarMap[key]; ok {
												_, _ = outputFile.WriteString(", " + val)
											}		
											// _, _ = outputFile.WriteString(", " + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName)
										}
										if h == 0 {
											break
										}
									}
								}
							}
							_, _ = outputFile.WriteString(")" + "\n")
						}
					}
				}
			}
			_, _ = outputFile.WriteString("	" + tableX[i].OutputFileName + "1 := createRandom" + tableX[i].FunctionSignature + "(t")
			for g := 0; g < len(fk_HierarchyX); g++ {
				if fk_HierarchyX[g].TableName == tableX[i].Table_name {
					if len(fk_HierarchyX[g].RelatedTablesLevels) > 0 {
						for h := 0; h < len(fk_HierarchyX[g].RelatedTablesLevels); h++ {
							for z := 0; z < len(fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList); z++ {
								key := fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName
								if val, ok := fkVarMap[key]; ok {
									_, _ = outputFile.WriteString(", " + val)
								}
								// _, _ = outputFile.WriteString(", " + tableX[i].Table_name+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName)
							}
							_, _ = outputFile.WriteString(")" + "\n")
							if h == 0 {
								break
							}
						}
					} else {
						_, _ = outputFile.WriteString(")" + "\n")
					}
				}
			}
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

func printTestFuncForReadList(tableX []dbschemareader.Table_Struct, i int, fk_HierarchyX []dbschemareader.FK_Hierarchy, outputFile *os.File) {
	var fkVarMap = make(map[string]string)
	_, _ = outputFile.WriteString("func TestList" + tableX[i].FunctionSignature2 + "(t *testing.T) {" + "\n")
	for k := 0; k < len(fk_HierarchyX); k++ {
		if fk_HierarchyX[k].TableName == tableX[i].Table_name {
			for l := len(fk_HierarchyX[k].RelatedTablesLevels) - 1; l >= 0; l-- {
				for m := 0; m < len(fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList); m++ {
					varName := fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName+"_fk_"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName+strconv.Itoa(k) + strconv.Itoa(l) + strconv.Itoa(m)
					key := fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName+"_fk_"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName
					// Only store if key doesn't exist
					if _, exists := fkVarMap[key]; !exists {
						fkVarMap[key] = varName
						_, _ = outputFile.WriteString("	" + varName+ " := createRandom" + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName_Singular_Object + "(t")
					} else{
						continue
					}
					// _, _ = outputFile.WriteString("	" + fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName+"_fk_"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName + " := createRandom" + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName_Singular_Object + "(t")
					for g := 0; g < len(fk_HierarchyX); g++ {
						if fk_HierarchyX[g].TableName == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {
							for h := 0; h < len(fk_HierarchyX[g].RelatedTablesLevels); h++ {
								for z := 0; z < len(fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList); z++ {
									key := fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName
									if val, ok := fkVarMap[key]; ok {
										_, _ = outputFile.WriteString(", " + val)
									}		
									// _, _ = outputFile.WriteString(", " + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName)
								}
								if h == 0 {
									break
								}
							}
						}
					}
					_, _ = outputFile.WriteString(")" + "\n")
				}
			}
		}
	}
	_, _ = outputFile.WriteString("	for i := 0; i < 10; i++ {" + "\n")
	_, _ = outputFile.WriteString("		createRandom" + tableX[i].FunctionSignature + "(t")
	for g := 0; g < len(fk_HierarchyX); g++ {
		if fk_HierarchyX[g].TableName == tableX[i].Table_name {
			for h := 0; h < len(fk_HierarchyX[g].RelatedTablesLevels); h++ {
				for z := 0; z < len(fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList); z++ {
					key := fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName
					if val, ok := fkVarMap[key]; ok {
						_, _ = outputFile.WriteString(", " + val)
					}		
					// _, _ = outputFile.WriteString(", " + tableX[i].Table_name+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName)
				}
				if h == 0 {
					break
				}
			}
		}
	}
	_, _ = outputFile.WriteString(")" + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString("	" + "}" + "\n")

	_, _ = outputFile.WriteString("	arg := List" + tableX[i].FunctionSignature2 + "Params{" + "\n")
	for g := 0; g < len(tableX[i].Table_Columns); g++ {
		if tableX[i].Table_Columns[g].ForeignFlag {
			for r := 0; r < len(tableX[i].ForeignKeys); r++ {
				if tableX[i].ForeignKeys[r].FK_Column == tableX[i].Table_Columns[g].Column_name {
					// fmt.Println(tableX[i].Table_Columns[g].Column_name)
					FormatedFieldName := FormatFieldName(tableX[i].ForeignKeys[r].FK_Related_Table_Column)
					// if tableX[i].Table_name == "bands" {
					// 	fmt.Println(FormatedFieldName)
					// }
					key := tableX[i].Table_name+"_fk_"+tableX[i].ForeignKeys[r].FK_Related_SingularTableName
					if val, ok := fkVarMap[key]; ok {
						_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[g].ColumnNameParams + ": " + val + "." + FormatedFieldName+","+"\n")
					}						
					// _, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[g].ColumnNameParams + ": " + tableX[i].Table_name+"_fk_"+tableX[i].ForeignKeys[r].FK_Related_SingularTableName + "." + FormatedFieldName + "," + "\n")
					// _, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[g].ColumnNameParams + ": " + tableX[i].Table_name+"_fk_"+tableX[i].ForeignKeys[r].FK_Related_SingularTableName + "." + strings.ToUpper(tableX[i].ForeignKeys[r].FK_Related_Table_Column) + "," + "\n")
				}
			}
		}
	}
	_, _ = outputFile.WriteString("		Limit:         5," + "\n")
	_, _ = outputFile.WriteString("		Offset:        5," + "\n")
	_, _ = outputFile.WriteString("	}" + "\n")

	_, _ = outputFile.WriteString("	" + tableX[i].Table_name + ", err := testStore.List" + tableX[i].FunctionSignature2 + "(context.Background(), " + "arg" + ")" + "\n")
	_, _ = outputFile.WriteString("	" + "require.NoError(t, err)" + "\n")
	_, _ = outputFile.WriteString("	" + "require.Len(t, " + tableX[i].Table_name + ", 5)" + "\n")
	_, _ = outputFile.WriteString("\n")

	_, _ = outputFile.WriteString("	for _, " + tableX[i].OutputFileName + " := range " + tableX[i].Table_name + " {" + "\n")
	_, _ = outputFile.WriteString("		require.NotEmpty(t, " + tableX[i].OutputFileName + ")" + "\n")
	str := "		require.True(t, "
	pipeflag := false
	if len(tableX[i].ForeignKeys) == 1 {
		for g := 0; g < len(tableX[i].Table_Columns); g++ {
			if tableX[i].Table_Columns[g].ForeignFlag {
				for r := 0; r < len(tableX[i].ForeignKeys); r++ {
					if tableX[i].ForeignKeys[r].FK_Column == tableX[i].Table_Columns[g].Column_name {
						str = str + "arg." + tableX[i].Table_Columns[g].ColumnNameParams + " == " + tableX[i].OutputFileName + "." + tableX[i].Table_Columns[g].ColumnNameParams
						pipeflag = true
					}
				}
			}
		}
		str = str + ")"
		_, _ = outputFile.WriteString(str + "\n")
	}
	if len(tableX[i].ForeignKeys) > 1 {
		for g := 0; g < len(tableX[i].Table_Columns); g++ {
			if tableX[i].Table_Columns[g].ForeignFlag {
				if pipeflag {
					str = str + " || "
				}
				for r := 0; r < len(tableX[i].ForeignKeys); r++ {
					if tableX[i].ForeignKeys[r].FK_Column == tableX[i].Table_Columns[g].Column_name {
						key := tableX[i].Table_name+"_fk_"+tableX[i].ForeignKeys[r].FK_Related_SingularTableName
						if val, ok := fkVarMap[key]; ok {
							FormatedFieldName := FormatFieldName(tableX[i].ForeignKeys[r].FK_Related_Table_Column)
							str = str + tableX[i].OutputFileName + "." + tableX[i].Table_Columns[g].ColumnNameParams + " == " + val + "." + FormatedFieldName
						}							
						// str = str + tableX[i].OutputFileName + "." + tableX[i].Table_Columns[g].ColumnNameParams + " == " + tableX[i].Table_name+"_fk_"+tableX[i].ForeignKeys[r].FK_Related_SingularTableName + "." + strings.ToUpper(tableX[i].ForeignKeys[r].FK_Related_Table_Column)
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

func printTestFuncForUpdate(tableX []dbschemareader.Table_Struct, i int, fk_HierarchyX []dbschemareader.FK_Hierarchy, outputFile *os.File) {
	var fkVarMap = make(map[string]string)
	_, _ = outputFile.WriteString("func TestUpdate" + tableX[i].FunctionSignature + "(t *testing.T) {" + "\n")
	for k := 0; k < len(fk_HierarchyX); k++ {
		if fk_HierarchyX[k].TableName == tableX[i].Table_name {
			for l := len(fk_HierarchyX[k].RelatedTablesLevels) - 1; l >= 0; l-- {
				for m := 0; m < len(fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList); m++ {
					varName := fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName+"_fk_"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName+strconv.Itoa(k) + strconv.Itoa(l) + strconv.Itoa(m)
					key := fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName+"_fk_"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName
					// Only store if key doesn't exist
					if _, exists := fkVarMap[key]; !exists {
						fkVarMap[key] = varName
						_, _ = outputFile.WriteString("	" + varName+ " := createRandom" + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName_Singular_Object + "(t")
					} else{
						continue
					}
					// _, _ = outputFile.WriteString("	" + fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName+"_fk_"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName + " := createRandom" + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName_Singular_Object + "(t")
					for g := 0; g < len(fk_HierarchyX); g++ {
						if fk_HierarchyX[g].TableName == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {
							for h := 0; h < len(fk_HierarchyX[g].RelatedTablesLevels); h++ {
								for z := 0; z < len(fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList); z++ {
									key := fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName
									if val, ok := fkVarMap[key]; ok {
										_, _ = outputFile.WriteString(", " + val)
									}		
									// _, _ = outputFile.WriteString(", " + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName)
								}
								if h == 0 {
									break
								}
							}
						}
					}
					_, _ = outputFile.WriteString(")" + "\n")
				}
			}
		}
	}
	_, _ = outputFile.WriteString("	" + tableX[i].OutputFileName + "1 := createRandom" + tableX[i].FunctionSignature + "(t")
	for g := 0; g < len(fk_HierarchyX); g++ {
		if fk_HierarchyX[g].TableName == tableX[i].Table_name {
			if len(fk_HierarchyX[g].RelatedTablesLevels) > 0 {
				for h := 0; h < len(fk_HierarchyX[g].RelatedTablesLevels); h++ {
					for z := 0; z < len(fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList); z++ {
						key := fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName
						if val, ok := fkVarMap[key]; ok {
							_, _ = outputFile.WriteString(", " + val)
						}		
						// _, _ = outputFile.WriteString(", " + tableX[i].Table_name+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName)
					}
					_, _ = outputFile.WriteString(")" + "\n")
					if h == 0 {
						break
					}
				}
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
	if tableX[i].Table_name == "users" {
		_, _ = outputFile.WriteString("	hashedPassword, err := util.HashPassword(util.RandomString(6))" + "\n")
		_, _ = outputFile.WriteString("	require.NoError(t, err)" + "\n")
	}
	_, _ = outputFile.WriteString("	arg := Update" + tableX[i].FunctionSignature + "Params{" + "\n")
	for p := 0; p < len(tableX[i].Table_Columns); p++ {
		if tableX[i].Table_Columns[p].ForeignFlag {
			for k := 0; k < len(tableX[i].ForeignKeys); k++ {
				if tableX[i].ForeignKeys[k].FK_Column == tableX[i].Table_Columns[p].Column_name {
					FormatedFieldName := FormatFieldName(tableX[i].ForeignKeys[k].FK_Related_Table_Column)
					key := tableX[i].Table_name+"_fk_"+tableX[i].ForeignKeys[k].FK_Related_SingularTableName
					if val, ok := fkVarMap[key]; ok {
						_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ":    	" + val + "." + FormatedFieldName+","+"\n")
					}						
					// _, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ":    	" + tableX[i].Table_name+"_fk_"+tableX[i].ForeignKeys[k].FK_Related_SingularTableName + "." + FormatedFieldName+","+"\n")
					// _, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ":    	" + tableX[i].Table_name+"_fk_"+tableX[i].ForeignKeys[k].FK_Related_SingularTableName + "." + strings.ToUpper((tableX[i].ForeignKeys[k].FK_Related_Table_Column)+","+"\n"))
				}
			}
		} else {
			if tableX[i].Table_Columns[p].ColumnType == "bigserial" || tableX[i].Table_Columns[p].ColumnType == "uuid" {
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ":    " + tableX[i].OutputFileName + "1." + getByColumnName + "," + "\n")
				continue
			}
			if tableX[i].Table_Columns[p].ColumnType == "varchar" && tableX[i].Table_Columns[p].PrimaryFlag {
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ":    " + tableX[i].OutputFileName + "1." + getByColumnName + "," + "\n")
				continue
			}
			if tableX[i].Table_Columns[p].ColumnType == "varchar" && !tableX[i].Table_Columns[p].PrimaryFlag {
				if tableX[i].Table_name == "users" && tableX[i].Table_Columns[p].Column_name == "email" {
					continue
				}
				if tableX[i].Table_name == "users" && tableX[i].Table_Columns[p].Column_name == "hashed_password" {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ":    hashedPassword," + "\n")
					continue
				} else {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ":    util.RandomName(8)," + "\n")
					continue
				}
			}
			if tableX[i].Table_Columns[p].ColumnType == "bigint" || tableX[i].Table_Columns[p].ColumnType == "bigserial" {
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ":    util.RandomInteger(1, 100)," + "\n")
				continue
			}
			if tableX[i].Table_Columns[p].ColumnType == "real" {
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ":    util.RandomReal(1, 100)," + "\n")
				continue
			}
			if tableX[i].Table_Columns[p].ColumnType == "timestamptz" {
				if tableX[i].Table_name == "users" && tableX[i].Table_Columns[p].Column_name == "password_created_at" {
					continue
				} else {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ":    time.Now().UTC()," + "\n")
					continue
				}
			}
		}
	}
	_, _ = outputFile.WriteString("	}" + "\n")
	_, _ = outputFile.WriteString("	" + tableX[i].OutputFileName + "2, err := testStore.Update" + tableX[i].FunctionSignature + "(context.Background(), " + "arg" + ")" + "\n")
	_, _ = outputFile.WriteString("	" + "require.NoError(t, err)" + "\n")
	_, _ = outputFile.WriteString("	" + "require.NotEmpty(t, " + tableX[i].OutputFileName + "2)" + "\n")
	_, _ = outputFile.WriteString("\n")
	for h := 0; h < len(tableX[i].Table_Columns); h++ {
		if tableX[i].Table_name == "users" && tableX[i].Table_Columns[h].Column_name == "email" {
			continue
		}
		if tableX[i].Table_name == "users" && tableX[i].Table_Columns[h].Column_name == "password_created_at" {
			continue
		}
		if tableX[i].Table_Columns[h].PrimaryFlag || tableX[i].Table_Columns[h].ForeignFlag {
			_, _ = outputFile.WriteString("	require.Equal(t, " + tableX[i].OutputFileName + "1." + tableX[i].Table_Columns[h].ColumnNameParams + ", " + tableX[i].OutputFileName + "2." + tableX[i].Table_Columns[h].ColumnNameParams + ")" + "\n")
			continue
		}
		if tableX[i].Table_Columns[h].ColumnType == "timestamptz" {
			_, _ = outputFile.WriteString("	require.WithinDuration(t, " + "arg." + tableX[i].Table_Columns[h].ColumnNameParams + ", " + tableX[i].OutputFileName + "2." + tableX[i].Table_Columns[h].ColumnNameParams + ", time.Second)" + "\n")
			continue
		}else{
			_, _ = outputFile.WriteString("	require.Equal(t, " + "arg." + tableX[i].Table_Columns[h].ColumnNameParams + ", " + tableX[i].OutputFileName + "2." + tableX[i].Table_Columns[h].ColumnNameParams + ")" + "\n")
			continue
		}
		
	}
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString("}" + "\n")
	_, _ = outputFile.WriteString("\n")
}

func printTestFuncForDelete(tableX []dbschemareader.Table_Struct, i int, fk_HierarchyX []dbschemareader.FK_Hierarchy, outputFile *os.File) {
	var fkVarMap = make(map[string]string)
	_, _ = outputFile.WriteString("func TestDelete" + tableX[i].FunctionSignature + "(t *testing.T) {" + "\n")
	for k := 0; k < len(fk_HierarchyX); k++ {
		if fk_HierarchyX[k].TableName == tableX[i].Table_name {
			for l := len(fk_HierarchyX[k].RelatedTablesLevels) - 1; l >= 0; l-- {
				for m := 0; m < len(fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList); m++ {
					varName := fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName+"_fk_"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName+strconv.Itoa(k) + strconv.Itoa(l) + strconv.Itoa(m)
					key := fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName+"_fk_"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName
					// Only store if key doesn't exist
					if _, exists := fkVarMap[key]; !exists {
						fkVarMap[key] = varName
						_, _ = outputFile.WriteString("	" + varName+ " := createRandom" + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName_Singular_Object + "(t")
					} else{
						continue
					}
					// _, _ = outputFile.WriteString("	" + fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName+"_fk_"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName + " := createRandom" + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName_Singular_Object + "(t")
					for g := 0; g < len(fk_HierarchyX); g++ {
						if fk_HierarchyX[g].TableName == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {
							for h := 0; h < len(fk_HierarchyX[g].RelatedTablesLevels); h++ {
								for z := 0; z < len(fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList); z++ {
									key := fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName
									if val, ok := fkVarMap[key]; ok {
										_, _ = outputFile.WriteString(", " + val)
									}		
									// _, _ = outputFile.WriteString(", " + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName)
								}
								if h == 0 {
									break
								}
							}
						}
					}
					_, _ = outputFile.WriteString(")" + "\n")
				}
			}
		}
	}
	_, _ = outputFile.WriteString("	" + tableX[i].OutputFileName + "1 := createRandom" + tableX[i].FunctionSignature + "(t")
	for g := 0; g < len(fk_HierarchyX); g++ {
		if fk_HierarchyX[g].TableName == tableX[i].Table_name {
			if len(fk_HierarchyX[g].RelatedTablesLevels) > 0 {
				for h := 0; h < len(fk_HierarchyX[g].RelatedTablesLevels); h++ {
					for z := 0; z < len(fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList); z++ {
						key := fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName
						if val, ok := fkVarMap[key]; ok {
							_, _ = outputFile.WriteString(", " + val)
						}		
						// _, _ = outputFile.WriteString(", " + tableX[i].Table_name+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName)
					}
					_, _ = outputFile.WriteString(")" + "\n")
					if h == 0 {
						break
					}
				}
			} else {
				_, _ = outputFile.WriteString(")" + "\n")
			}
		}
	}
	// for g := 0; g < len(tableX); g++ {
	// 	if tableX[g].Table_name == tableX[i].Table_name {
	// 		if len(tableX[g].ForeignKeys) > 0 {
	// 			for h := 0; h < len(tableX[g].ForeignKeys); h++ {
	// 				_, _ = outputFile.WriteString(", " + tableX[g].ForeignKeys[h].FK_Related_SingularTableName)
	// 			}
	// 			_, _ = outputFile.WriteString(")" + "\n")
	// 			break
	// 		} else {
	// 			_, _ = outputFile.WriteString(")" + "\n")
	// 		}
	// 	}
	// }
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
	_, _ = outputFile.WriteString("	" + "require.Empty(t, " + tableX[i].OutputFileName + "2)" + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString("}" + "\n")
	_, _ = outputFile.WriteString("\n")
}

func TestWriter(projectFolderName string, gitHubAccountName string, dirPath string) {
	//generating main_test.go
	/////////////////////////////
	main_testFunc(projectFolderName, gitHubAccountName, dirPath)
	/////////////////////////////////////////////////
	//generate unit tests for go file in sqlc folder
	/////////////////////////////////////////////////
	var files []string
	pathToSearch := dirPath + "/db/migration"
	err := filepath.Walk(pathToSearch, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	var tableX []dbschemareader.Table_Struct
	var fk_HierarchyX []dbschemareader.FK_Hierarchy
	fmt.Println("generating unit tests.....")
	for _, element := range files {
		if element[len(element)-6:] == `up.sql` {
			tableX, fk_HierarchyX = dbschemareader.ReadSchema(element, tableX)
			for i := 0; i < len(tableX); i++ {
				if tableX[i].Table_name == "sessions" ||  tableX[i].Table_name == "activities"{
					continue
				}
				outputFile, errs := os.Create(dirPath + "/db/sqlc/" + tableX[i].OutputFileName + "_test.go")
				if errs != nil {
					fmt.Println("Failed to create file:", errs)
					return
				}
				defer outputFile.Close()
				_, _ = outputFile.WriteString("package db" + "\n")
				_, _ = outputFile.WriteString("\n")
				_, _ = outputFile.WriteString("import (" + "\n")
				_, _ = outputFile.WriteString(`	"context"` + "\n")
				_, _ = outputFile.WriteString(`	"time"` + "\n")
				_, _ = outputFile.WriteString(`	"testing"` + "\n")
				_, _ = outputFile.WriteString(`	"github.com/stretchr/testify/require"` + "\n")
				_, _ = outputFile.WriteString(`	"github.com/`+gitHubAccountName+`/` + projectFolderName + `/util"` + "\n")
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
	_, _ = outputFile.WriteString(`	//if errors.As(err, &pgErr) {` + "\n")
	_, _ = outputFile.WriteString(`		//return pgErr.Code` + "\n")
	_, _ = outputFile.WriteString(`	//}` + "\n")
	_, _ = outputFile.WriteString(`	//return ""` + "\n")
	_, _ = outputFile.WriteString(`//}` + "\n")

	//writing jwt_maker_test.go
	outputFileName = dirPath + "/token/jwt_maker_test.go"
	outputFile, errs = os.Create(outputFileName)
	if errs != nil {
		fmt.Println("Failed to create file:", errs)
		return
	}
	defer outputFile.Close()
	_, _ = outputFile.WriteString("package token" + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString("import (" + "\n")
	_, _ = outputFile.WriteString(`	"testing"` + "\n")
	_, _ = outputFile.WriteString(`	"time"` + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString(`	"github.com/golang-jwt/jwt"` + "\n")
	_, _ = outputFile.WriteString(`	"github.com/stretchr/testify/require"` + "\n")
	_, _ = outputFile.WriteString(`	"github.com/`+gitHubAccountName+`/` + projectFolderName + `/util"` + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString(`)` + "\n")
	_, _ = outputFile.WriteString("\n")

	_, _ = outputFile.WriteString(`func TestJWTMaker(t *testing.T) {` + "\n")
	_, _ = outputFile.WriteString(`	maker, err := NewJWTMaker(util.RandomString(32))` + "\n")
	_, _ = outputFile.WriteString(`	require.NoError(t, err)` + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString(`	username := util.RandomName(8)` + "\n")
	_, _ = outputFile.WriteString(`	role := util.UserLevel_1_Role` + "\n")
	_, _ = outputFile.WriteString(`	duration := time.Minute` + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString(`	issuedAt := time.Now()` + "\n")
	_, _ = outputFile.WriteString(`	expiredAt := issuedAt.Add(duration)` + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString(`	token, payload, err := maker.CreateToken(username, role, duration)` + "\n")
	_, _ = outputFile.WriteString(`	//token, err := maker.CreateToken(username, role, duration)` + "\n")
	_, _ = outputFile.WriteString(`	require.NoError(t, err)` + "\n")
	_, _ = outputFile.WriteString(`	require.NotEmpty(t, token)` + "\n")
	_, _ = outputFile.WriteString(`	require.NotEmpty(t, payload)` + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString(`	payload, err = maker.VerifyToken(token)` + "\n")
	_, _ = outputFile.WriteString(`	require.NoError(t, err)` + "\n")
	_, _ = outputFile.WriteString(`	require.NotEmpty(t, payload)` + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString(`	require.NotZero(t, payload.ID)` + "\n")
	_, _ = outputFile.WriteString(`	require.Equal(t, username, payload.Username)` + "\n")
	_, _ = outputFile.WriteString(`	//require.Equal(t, role, payload.Role)` + "\n")
	_, _ = outputFile.WriteString(`	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)` + "\n")
	_, _ = outputFile.WriteString(`	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)` + "\n")
	_, _ = outputFile.WriteString(`}` + "\n")

	_, _ = outputFile.WriteString(`func TestExpiredJWTToken(t *testing.T) {` + "\n")
	_, _ = outputFile.WriteString(`	maker, err := NewJWTMaker(util.RandomString(32))` + "\n")
	_, _ = outputFile.WriteString(`	require.NoError(t, err)` + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString(`	//token, payload, err := maker.CreateToken(util.RandomName(8), util.DepositorRole, -time.Minute)` + "\n")
	_, _ = outputFile.WriteString(`	token, payload, err := maker.CreateToken(util.RandomName(8), util.UserLevel_1_Role, -time.Minute)` + "\n")
	_, _ = outputFile.WriteString(`	require.NoError(t, err)` + "\n")
	_, _ = outputFile.WriteString(`	require.NotEmpty(t, token)` + "\n")
	_, _ = outputFile.WriteString(`	require.NotEmpty(t, payload)` + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString(`	payload, err = maker.VerifyToken(token)` + "\n")
	_, _ = outputFile.WriteString(`	require.Error(t, err)` + "\n")
	_, _ = outputFile.WriteString(`	require.EqualError(t, err, ErrExpiredToken.Error())` + "\n")
	_, _ = outputFile.WriteString(`	require.Nil(t, payload)` + "\n")
	_, _ = outputFile.WriteString(`}` + "\n")
	_, _ = outputFile.WriteString("\n")

	_, _ = outputFile.WriteString(`func TestInvalidJWTTokenAlgNone(t *testing.T) {` + "\n")
	_, _ = outputFile.WriteString(`	//payload, err := NewPayload(util.RandomOwner(), util.DepositorRole, time.Minute)` + "\n")
	_, _ = outputFile.WriteString(`	payload, err := NewPayload(util.RandomName(8), util.UserLevel_1_Role, time.Minute)` + "\n")
	_, _ = outputFile.WriteString(`	require.NoError(t, err)` + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString(`	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)` + "\n")
	_, _ = outputFile.WriteString(`	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)` + "\n")
	_, _ = outputFile.WriteString(`	require.NoError(t, err)` + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString(`	maker, err := NewJWTMaker(util.RandomString(32))` + "\n")
	_, _ = outputFile.WriteString(`	require.NoError(t, err)` + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString(`	payload, err = maker.VerifyToken(token)` + "\n")
	_, _ = outputFile.WriteString(`	require.Error(t, err)` + "\n")
	_, _ = outputFile.WriteString(`	require.EqualError(t, err, ErrInvalidToken.Error())` + "\n")
	_, _ = outputFile.WriteString(`	require.Nil(t, payload)` + "\n")
	_, _ = outputFile.WriteString(`}` + "\n")
	_, _ = outputFile.WriteString("\n")
	fmt.Println("jwt_maker_test.go file has been generated successfully")
	outputFile.Close()

	//writing paseto_maker_test.go
	outputFileName = dirPath + "/token/paseto_maker_test.go"
	outputFile, errs = os.Create(outputFileName)
	if errs != nil {
		fmt.Println("Failed to create file:", errs)
		return
	}
	defer outputFile.Close()
	_, _ = outputFile.WriteString("package token" + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString("import (" + "\n")
	_, _ = outputFile.WriteString(`	"testing"` + "\n")
	_, _ = outputFile.WriteString(`	"time"` + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString(`	"github.com/stretchr/testify/require"` + "\n")
	_, _ = outputFile.WriteString(`	"github.com/`+gitHubAccountName+`/` + projectFolderName + `/util"` + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString(`)` + "\n")
	_, _ = outputFile.WriteString("\n")

	_, _ = outputFile.WriteString(`func TestPasetoMaker(t *testing.T) {` + "\n")
	_, _ = outputFile.WriteString(`	maker, err := NewPasetoMaker(util.RandomString(32))` + "\n")
	_, _ = outputFile.WriteString(`	require.NoError(t, err)` + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString(`	username := util.RandomName(8)` + "\n")
	_, _ = outputFile.WriteString(`	role := util.UserLevel_1_Role` + "\n")
	_, _ = outputFile.WriteString(`	duration := time.Minute` + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString(`	issuedAt := time.Now()` + "\n")
	_, _ = outputFile.WriteString(`	expiredAt := issuedAt.Add(duration)` + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString(`	token, payload, err := maker.CreateToken(username, role, duration)` + "\n")
	_, _ = outputFile.WriteString(`	//token, err := maker.CreateToken(username, role, duration)` + "\n")
	_, _ = outputFile.WriteString(`	require.NoError(t, err)` + "\n")
	_, _ = outputFile.WriteString(`	require.NotEmpty(t, token)` + "\n")
	_, _ = outputFile.WriteString(`	require.NotEmpty(t, payload)` + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString(`	payload, err = maker.VerifyToken(token)` + "\n")
	_, _ = outputFile.WriteString(`	require.NoError(t, err)` + "\n")
	_, _ = outputFile.WriteString(`	require.NotEmpty(t, payload)` + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString(`	require.NotZero(t, payload.ID)` + "\n")
	_, _ = outputFile.WriteString(`	require.Equal(t, username, payload.Username)` + "\n")
	_, _ = outputFile.WriteString(`	//require.Equal(t, role, payload.Role)` + "\n")
	_, _ = outputFile.WriteString(`	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)` + "\n")
	_, _ = outputFile.WriteString(`	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)` + "\n")
	_, _ = outputFile.WriteString(`}` + "\n")

	_, _ = outputFile.WriteString(`func TestExpiredPasetoToken(t *testing.T) {` + "\n")
	_, _ = outputFile.WriteString(`	maker, err := NewPasetoMaker(util.RandomString(32))` + "\n")
	_, _ = outputFile.WriteString(`	require.NoError(t, err)` + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString(`	//token, payload, err := maker.CreateToken(util.RandomName(8), util.DepositorRole, -time.Minute)` + "\n")
	_, _ = outputFile.WriteString(`	token, payload, err := maker.CreateToken(util.RandomName(8), util.UserLevel_1_Role, -time.Minute)` + "\n")
	_, _ = outputFile.WriteString(`	require.NoError(t, err)` + "\n")
	_, _ = outputFile.WriteString(`	require.NotEmpty(t, token)` + "\n")
	_, _ = outputFile.WriteString(`	require.NotEmpty(t, payload)` + "\n")
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString(`	payload, err = maker.VerifyToken(token)` + "\n")
	_, _ = outputFile.WriteString(`	require.Error(t, err)` + "\n")
	_, _ = outputFile.WriteString(`	require.EqualError(t, err, ErrExpiredToken.Error())` + "\n")
	_, _ = outputFile.WriteString(`	require.Nil(t, payload)` + "\n")
	_, _ = outputFile.WriteString(`}` + "\n")
	_, _ = outputFile.WriteString("\n")
	fmt.Println("paseto_maker_test.go file has been generated successfully")
	outputFile.Close()

	//Executing goimports
	cmd := exec.Command("goimports", "-w", ".")
	cmd.Dir = dirPath + "/db/sqlc"
	cmd.Run()
	println("goimports executed successfully")
	// var stderr bytes.Buffer

	//Executing go mod tidy
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = dirPath
	cmd.Run()
	println("go mod tidy executed successfully")

}
