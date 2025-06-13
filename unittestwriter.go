package unittestwriter

import (
	"fmt"
	"os"
	"time"

	// "os/exec"
	// "path/filepath"
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
	// fmt.Println("input: ", input)
	parts := strings.Split(input, "_")
	// fmt.Println("parts just after split: ")
	for i := 0; i < len(parts); i++ {
		word := strings.ToLower(parts[i])
		// fmt.Println("word: ", word)
		if word == "id" {
			parts[i] = "ID"
		} else if len(word) > 0 {
			parts[i] = strings.ToUpper(word[:1]) + word[1:]
			// fmt.Println("parts: ", parts)
		}
	}
	// Ensure first letter of the final result is capitalized
	result := strings.Join(parts, "")
	if len(result) > 0 {
		result = strings.ToUpper(result[:1]) + result[1:]
	}
	// fmt.Println("result: ", result)
	return result
}
// Add these helper functions to unittestwriter.go

// generateValidValueForCheckConstraint generates a valid value based on check constraint
func generateValidValueForCheckConstraint(tableName, columnName string, columnType string, checkConstraints []dbschemareader.CheckConstraint) string {
	// fmt.Printf("DEBUG: Table=%s, Column=%s, Type=%s\n", tableName, columnName, columnType)
	// if tableName == "resources" && columnName == "publication_year"{
	// 	fmt.Printf("DEBUG: Table=%s, Column=%s, Type=%s\n", tableName, columnName, columnType)
	// }
	// Find check constraint for this column
	for _, constraint := range checkConstraints {
		if constraint.CheckConstraintColumnName == columnName {
			// if tableName == "resources" && columnName == "publication_year"{
			// 	fmt.Printf("DEBUG: Found constraint for %s: %s\n", columnName, constraint.CheckConstraintValue)
			// }
			result := generateValueFromConstraint(constraint.CheckConstraintValue, columnType)
			// fmt.Printf("DEBUG: Generated value for %s: %s\n", columnName, result)  // ← This line should show the result
			return result
		}
	}	
	// No constraint found, use default random generation
	// fmt.Printf("DEBUG: No constraint found for %s, using default\n", columnName)
	return generateDefaultRandomValue(columnType)
}

// Enhanced generateValueFromConstraint to handle all constraint types
func generateValueFromConstraint(constraintValue, columnType string) string {
	originalValue := constraintValue
	constraintValue = strings.ToLower(strings.TrimSpace(constraintValue))
	
	// fmt.Printf("DEBUG: generateValueFromConstraint called with: '%s', columnType: '%s'\n", constraintValue, columnType)
	
	// Handle complex constraints with OR conditions
	if strings.Contains(constraintValue, " or ") {
		return handleOrConstraint(constraintValue, columnType)
	}
	
	// Handle complex constraints with AND conditions
	if strings.Contains(constraintValue, " and ") {
		return handleAndConstraint(constraintValue, columnType)
	}
	
	// ✅ NEW: Handle arithmetic expressions (addition, subtraction, multiplication, division)
	if strings.Contains(constraintValue, "+") || 
	   strings.Contains(constraintValue, "-") || 
	   strings.Contains(constraintValue, "*") || 
	   strings.Contains(constraintValue, "/") {
		return handleArithmeticConstraint(constraintValue, columnType)
	}

	// Handle REGEX constraints (various regex operators)
	if strings.Contains(constraintValue, "~*") || 
	   strings.Contains(constraintValue, "~") || 
	   strings.Contains(constraintValue, "regexp") ||
	   strings.Contains(constraintValue, "rlike") {
		return handleRegexConstraint(originalValue, columnType)
	}
	
	// Handle IN constraints with various formats
	if strings.Contains(constraintValue, " in ") || 
	   strings.Contains(constraintValue, " any ") {
		return handleInConstraint(constraintValue, columnType)
	}
	
	// Handle range constraints
	if strings.Contains(constraintValue, ">=") || strings.Contains(constraintValue, ">") {
		return handleRangeConstraint(constraintValue, columnType, "min")
	}
	
	if strings.Contains(constraintValue, "<=") || strings.Contains(constraintValue, "<") {
		return handleRangeConstraint(constraintValue, columnType, "max")
	}
	
	// Handle equality constraints
	if strings.Contains(constraintValue, "=") && 
	   !strings.Contains(constraintValue, ">=") && 
	   !strings.Contains(constraintValue, "<=") &&
	   !strings.Contains(constraintValue, "!=") &&
	   !strings.Contains(constraintValue, "<>") {
		return handleEqualityConstraint(constraintValue, columnType)
	}
	
	// Handle LIKE/ILIKE constraints
	if strings.Contains(constraintValue, "like") || strings.Contains(constraintValue, "ilike") {
		return handleLikeConstraint(constraintValue, columnType)
	}
	
	// Handle BETWEEN constraints
	if strings.Contains(constraintValue, "between") {
		return handleBetweenConstraint(constraintValue, columnType)
	}
	
	// Handle NOT NULL constraints
	if strings.Contains(constraintValue, "not null") || strings.Contains(constraintValue, "is not null") {
		return generateDefaultRandomValue(columnType)
	}
	
	// Handle length constraints
	if strings.Contains(constraintValue, "length") || strings.Contains(constraintValue, "char_length") {
		return handleLengthConstraint(constraintValue, columnType)
	}
	
	// Default fallback
	// fmt.Printf("DEBUG: No specific constraint handler found, using default for columnType: %s\n", columnType)
	return generateDefaultRandomValue(columnType)
}
////////////////??????????????????//////////

// Handle arithmetic expressions in check constraints
func handleArithmeticConstraint(constraintValue, columnType string) string {
	// Handle different patterns of arithmetic constraints
	
	// Pattern: columnA = columnB + columnC + columnD
	if strings.Contains(constraintValue, "=") {
		return handleArithmeticEquality(constraintValue, columnType)
	}
	
	// Pattern: columnA + columnB > value
	if strings.Contains(constraintValue, ">") {
		return handleArithmeticComparison(constraintValue, columnType, ">")
	}
	
	// Pattern: columnA + columnB < value
	if strings.Contains(constraintValue, "<") {
		return handleArithmeticComparison(constraintValue, columnType, "<")
	}
	
	// Pattern: columnA + columnB >= value
	if strings.Contains(constraintValue, ">=") {
		return handleArithmeticComparison(constraintValue, columnType, ">=")
	}
	
	// Pattern: columnA + columnB <= value
	if strings.Contains(constraintValue, "<=") {
		return handleArithmeticComparison(constraintValue, columnType, "<=")
	}
	
	// Default fallback
	return generateDefaultRandomValue(columnType)
}

// Handle arithmetic equality constraints like: totalamount = subtotal + taxamount + shippingfee
func handleArithmeticEquality(constraintValue, columnType string) string {
	// Split on = to get left and right sides
	parts := strings.Split(constraintValue, "=")
	if len(parts) != 2 {
		return generateDefaultRandomValue(columnType)
	}
	
	//leftSide := strings.TrimSpace(parts[0])
	rightSide := strings.TrimSpace(parts[1])
	
	// Parse the arithmetic expression on the right side
	expression := parseArithmeticExpression(rightSide)
	
	if len(expression.Variables) > 0 {
		// Generate code that calculates the sum of referenced variables
		return generateArithmeticExpression(expression, columnType)
	}
	
	return generateDefaultRandomValue(columnType)
}

// Handle arithmetic comparison constraints like: subtotal + taxamount > 100
func handleArithmeticComparison(constraintValue, columnType, operator string) string {
	// Split on the operator
	parts := strings.Split(constraintValue, operator)
	if len(parts) != 2 {
		return generateDefaultRandomValue(columnType)
	}
	
	leftSide := strings.TrimSpace(parts[0])
	rightSide := strings.TrimSpace(parts[1])
	
	// Parse the arithmetic expression
	expression := parseArithmeticExpression(leftSide)
	
	// Parse the comparison value
	comparisonValue, err := strconv.ParseFloat(rightSide, 64)
	if err != nil {
		return generateDefaultRandomValue(columnType)
	}
	
	// Generate a value that satisfies the constraint
	return generateValueForArithmeticComparison(expression, operator, comparisonValue, columnType)
}

// Arithmetic expression structure
type ArithmeticExpression struct {
	Variables []string
	Operators []string
	Constants []float64
}

// Parse arithmetic expression into components
func parseArithmeticExpression(expression string) ArithmeticExpression {
	result := ArithmeticExpression{
		Variables: []string{},
		Operators: []string{},
		Constants: []float64{},
	}
	
	// Remove spaces and parentheses
	cleaned := strings.ReplaceAll(expression, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")
	
	// Split on operators while preserving them
	tokens := splitWithOperators(cleaned, []string{"+", "-", "*", "/"})
	
	for _, token := range tokens {
		if isOperator(token) {
			result.Operators = append(result.Operators, token)
		} else {
			// Check if it's a number or variable
			if value, err := strconv.ParseFloat(token, 64); err == nil {
				result.Constants = append(result.Constants, value)
			} else {
				result.Variables = append(result.Variables, token)
			}
		}
	}
	
	return result
}

// Generate arithmetic expression code
func generateArithmeticExpression(expr ArithmeticExpression, columnType string) string {
	if len(expr.Variables) == 0 {
		return generateDefaultRandomValue(columnType)
	}
	
	// Build the expression using variable references
	var expressionParts []string
	
	for _, variable := range expr.Variables {
		// Convert database column names to Go variable names
		goVarName := convertColumnNameToGoVariable(variable)
		expressionParts = append(expressionParts, goVarName)
	}
	
	// Join with + operator (most common case)
	if len(expressionParts) > 1 {
		return strings.Join(expressionParts, " + ")
	} else if len(expressionParts) == 1 {
		return expressionParts[0]
	}
	
	return generateDefaultRandomValue(columnType)
}

// Convert database column name to Go variable name
func convertColumnNameToGoVariable(columnName string) string {
	// Common conversions for your schema
	switch strings.ToLower(columnName) {
	case "subtotal":
		return "arg.Subtotal"
	case "taxamount":
		return "arg.Taxamount"
	case "shippingfee":
		return "arg.Shippingfee"
	case "totalamount":
		return "arg.Totalamount"
	case "discount":
		return "arg.Discount"
	case "amount":
		return "arg.Amount"
	default:
		// Generic conversion: convert to camelCase and add arg. prefix
		camelCase := toCamelCase(columnName)
		return "arg." + camelCase
	}
}

// Generate value for arithmetic comparison
func generateValueForArithmeticComparison(expr ArithmeticExpression, operator string, comparisonValue float64, columnType string) string {
	// For simplicity, generate a value that would satisfy the constraint
	switch operator {
	case ">":
		return fmt.Sprintf("util.RandomReal(%.2f, %.2f)", comparisonValue+1, comparisonValue+100)
	case ">=":
		return fmt.Sprintf("util.RandomReal(%.2f, %.2f)", comparisonValue, comparisonValue+100)
	case "<":
		if comparisonValue > 1 {
			return fmt.Sprintf("util.RandomReal(0.00, %.2f)", comparisonValue-1)
		}
		return "0.00"
	case "<=":
		return fmt.Sprintf("util.RandomReal(0.00, %.2f)", comparisonValue)
	default:
		return generateDefaultRandomValue(columnType)
	}
}

// Helper functions
func splitWithOperators(s string, operators []string) []string {
	result := []string{s}
	
	for _, op := range operators {
		var newResult []string
		for _, part := range result {
			subParts := strings.Split(part, op)
			for i, subPart := range subParts {
				if i > 0 {
					newResult = append(newResult, op)
				}
				if subPart != "" {
					newResult = append(newResult, subPart)
				}
			}
		}
		result = newResult
	}
	
	return result
}

func isOperator(token string) bool {
	operators := []string{"+", "-", "*", "/"}
	for _, op := range operators {
		if token == op {
			return true
		}
	}
	return false
}

func toCamelCase(s string) string {
	words := strings.FieldsFunc(s, func(c rune) bool {
		return c == '_' || c == '-' || c == ' '
	})
	
	if len(words) == 0 {
		return s
	}
	
	result := strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		result += strings.Title(strings.ToLower(words[i]))
	}
	
	return result
}


// Check if a column has an arithmetic constraint
func checkForArithmeticConstraint(tableName, columnName string, checkConstraints []dbschemareader.CheckConstraint) bool {
    for _, constraint := range checkConstraints {
        if constraint.CheckConstraintColumnName == columnName {
            constraintValue := strings.ToLower(strings.TrimSpace(constraint.CheckConstraintValue))
            // Check for arithmetic operators
            if strings.Contains(constraintValue, "+") || 
               strings.Contains(constraintValue, "-") || 
               strings.Contains(constraintValue, "*") || 
               strings.Contains(constraintValue, "/") {
                return true
            }
        }
    }
    return false
}

// Get all arithmetic constraints for a table
func getArithmeticConstraints(tableName string, checkConstraints []dbschemareader.CheckConstraint) []dbschemareader.CheckConstraint {
    var arithmeticConstraints []dbschemareader.CheckConstraint
    
    for _, constraint := range checkConstraints {
        constraintValue := strings.ToLower(strings.TrimSpace(constraint.CheckConstraintValue))
        if strings.Contains(constraintValue, "+") || 
           strings.Contains(constraintValue, "-") || 
           strings.Contains(constraintValue, "*") || 
           strings.Contains(constraintValue, "/") {
            arithmeticConstraints = append(arithmeticConstraints, constraint)
        }
    }
    
    return arithmeticConstraints
}

// Generate arithmetic assignment after struct creation
func generateArithmeticAssignment(constraint dbschemareader.CheckConstraint, columns []dbschemareader.Table_columns) string {
    constraintValue := strings.ToLower(strings.TrimSpace(constraint.CheckConstraintValue))
    
    // Handle pattern: columnA = columnB + columnC + columnD
    if strings.Contains(constraintValue, "=") {
        parts := strings.Split(constraintValue, "=")
        if len(parts) != 2 {
            return ""
        }
        
        leftSide := strings.TrimSpace(parts[0])
        rightSide := strings.TrimSpace(parts[1])
        
        // Find the target column
        targetColumn := findColumnByName(leftSide, columns)
        if targetColumn == nil {
            return ""
        }
        
        // Parse the arithmetic expression
        expression := parseArithmeticExpressionForAssignment(rightSide, columns)
        if expression == "" {
            return ""
        }
        
        // Generate the assignment
        return fmt.Sprintf("arg.%s = %s", targetColumn.ColumnNameParams, expression)
    }
    
    return ""
}

// Parse arithmetic expression for assignment
func parseArithmeticExpressionForAssignment(expression string, columns []dbschemareader.Table_columns) string {
    // Remove spaces and parentheses
    cleaned := strings.ReplaceAll(expression, " ", "")
    cleaned = strings.ReplaceAll(cleaned, "(", "")
    cleaned = strings.ReplaceAll(cleaned, ")", "")
    
    // Split on + operator (most common case)
    if strings.Contains(cleaned, "+") {
        parts := strings.Split(cleaned, "+")
        var goExpression []string
        
        for _, part := range parts {
            part = strings.TrimSpace(part)
            if part != "" {
                // Find the corresponding column
                column := findColumnByName(part, columns)
                if column != nil {
                    goExpression = append(goExpression, "arg."+column.ColumnNameParams)
                } else {
                    // It might be a constant
                    if _, err := strconv.ParseFloat(part, 64); err == nil {
                        goExpression = append(goExpression, part)
                    }
                }
            }
        }
        
        if len(goExpression) > 0 {
            return strings.Join(goExpression, " + ")
        }
    }
    
    return ""
}

// Find column by name
func findColumnByName(name string, columns []dbschemareader.Table_columns) *dbschemareader.Table_columns {
    for i := range columns {
        if strings.EqualFold(columns[i].Column_name, name) {
            return &columns[i]
        }
    }
    return nil
}








////////????????????????????????????/////////
// Comprehensive regex constraint handler
func handleRegexConstraint(constraintValue, columnType string) string {
	// constraintLower := strings.ToLower(constraintValue)
	// fmt.Printf("DEBUG: handleRegexConstraint called with: %s\n", constraintValue)

	// In handleRegexConstraint function, replace the email handling:
	if isEmailPattern(constraintValue) {
		// fmt.Printf("DEBUG: Detected email pattern, returning util.RandomEmail()\n")
		return "util.RandomEmail()"
	}	

	// Phone number patterns
	if isPhonePattern(constraintValue) {
		return generateValidPhone(constraintValue)
	}
	
	// URL patterns
	if isUrlPattern(constraintValue) {
		return generateValidUrl(constraintValue)
	}
	
	// Postal code patterns
	if isPostalCodePattern(constraintValue) {
		return generateValidPostalCode(constraintValue)
	}
	
	// Date format patterns
	if isDateFormatPattern(constraintValue) {
		return generateValidDateFormat(constraintValue)
	}
	
	// Username patterns
	if isUsernamePattern(constraintValue) {
		return generateValidUsername(constraintValue)
	}
	
	// Password patterns
	if isPasswordPattern(constraintValue) {
		return generateValidPassword(constraintValue)
	}
	
	// Generic alphanumeric patterns
	if isAlphaNumericPattern(constraintValue) {
		return generateAlphaNumeric(constraintValue)
	}
	
	// Numeric only patterns
	if isNumericPattern(constraintValue) {
		return generateNumericString(constraintValue)
	}
	
	// Default fallback for unrecognized patterns
	return generateDefaultRandomValue(columnType)
}

// Email pattern detection - comprehensive
// Email pattern detection - enhanced
func isEmailPattern(pattern string) bool {
	pattern = strings.ToLower(pattern)
	
	// Common email indicators - check for regex pattern specifically
	emailIndicators := []string{
		"@",
		"email",
		"mail",
		"\\.com",
		"\\.org", 
		"\\.net",
		"[a-z].*@.*[a-z]",
		"^.*@.*\\.",
		"[a-za-z0-9._%-]+@[a-za-z0-9-]+",
	}
	
	// Must contain @ and have domain-like structure for email regex
	hasAt := strings.Contains(pattern, "@")
	hasDomainStructure := strings.Contains(pattern, ".") || strings.Contains(pattern, "\\.")
	
	if hasAt && hasDomainStructure {
		return true
	}
	
	// Check for any email-specific indicators
	for _, indicator := range emailIndicators {
		if strings.Contains(pattern, indicator) {
			return true
		}
	}
	
	return false
}

// Phone pattern detection
func isPhonePattern(pattern string) bool {
	pattern = strings.ToLower(pattern)
	phoneIndicators := []string{
		"phone", "tel", "mobile", "cell",
		"\\d{3}", "\\d{10}", "\\d{11}",
		"[0-9].*[0-9].*[0-9]", // At least 3 digits
		"\\+\\d", // International format
		"\\(\\d{3}\\)", // US format with parentheses
	}
	
	for _, indicator := range phoneIndicators {
		if strings.Contains(pattern, indicator) {
			return true
		}
	}
	return false
}

func generateValidPhone(pattern string) string {
	if strings.Contains(pattern, "\\+") {
		return `"+1234567890"`
	}
	if strings.Contains(pattern, "\\(") && strings.Contains(pattern, "\\)") {
		return `"(555) 123-4567"`
	}
	if strings.Contains(pattern, "-") {
		return `"555-123-4567"`
	}
	return `"5551234567"`
}

// URL pattern detection
func isUrlPattern(pattern string) bool {
	pattern = strings.ToLower(pattern)
	urlIndicators := []string{
		"http", "https", "www", "url", "uri",
		"://", "\\.com", "\\.org", "\\.net",
	}
	
	for _, indicator := range urlIndicators {
		if strings.Contains(pattern, indicator) {
			return true
		}
	}
	return false
}

func generateValidUrl(pattern string) string {
	if strings.Contains(pattern, "https") {
		return `"https://example.com"`
	}
	if strings.Contains(pattern, "http") {
		return `"http://example.com"`
	}
	return `"https://www.example.com"`
}

// Postal code pattern detection
func isPostalCodePattern(pattern string) bool {
	pattern = strings.ToLower(pattern)
	postalIndicators := []string{
		"zip", "postal", "postcode",
		"\\d{5}", "\\d{4}", // US/International
		"[a-z]\\d[a-z]", // Canadian format
	}
	
	for _, indicator := range postalIndicators {
		if strings.Contains(pattern, indicator) {
			return true
		}
	}
	return false
}

func generateValidPostalCode(pattern string) string {
	if strings.Contains(pattern, "[a-z].*\\d.*[a-z]") {
		return `"A1B2C3"` // Canadian format
	}
	if strings.Contains(pattern, "\\d{5}") {
		return `"12345"`
	}
	return `"12345"`
}

// Date format pattern detection
func isDateFormatPattern(pattern string) bool {
	pattern = strings.ToLower(pattern)
	dateIndicators := []string{
		"yyyy", "mm", "dd", "date",
		"\\d{4}", "\\d{2}",
		"/", "-", ".",
	}
	
	hasDateStructure := false
	for _, indicator := range dateIndicators {
		if strings.Contains(pattern, indicator) {
			hasDateStructure = true
			break
		}
	}
	
	return hasDateStructure && (strings.Contains(pattern, "/") || 
		strings.Contains(pattern, "-") || strings.Contains(pattern, "\\."))
}

func generateValidDateFormat(pattern string) string {
	if strings.Contains(pattern, "/") {
		return `"12/31/2023"`
	}
	if strings.Contains(pattern, "-") {
		return `"2023-12-31"`
	}
	return `"2023.12.31"`
}

// Username pattern detection
func isUsernamePattern(pattern string) bool {
	pattern = strings.ToLower(pattern)
	usernameIndicators := []string{
		"username", "user", "login",
		"^[a-z]", "^[a-z0-9]",
		"[a-z0-9_]", "[a-z0-9\\.]",
	}
	
	for _, indicator := range usernameIndicators {
		if strings.Contains(pattern, indicator) {
			return true
		}
	}
	return false
}

func generateValidUsername(pattern string) string {
	base := "user123"
	if strings.Contains(pattern, "_") {
		base = "user_123"
	}
	if strings.Contains(pattern, "\\.") {
		base = "user.123"
	}
	return fmt.Sprintf(`"%s"`, base)
}

// Password pattern detection
func isPasswordPattern(pattern string) bool {
	pattern = strings.ToLower(pattern)
	passwordIndicators := []string{
		"password", "pwd", "pass",
		"(?=.*[a-z])", "(?=.*[A-Z])", "(?=.*\\d)", "(?=.*[@$!%*?&])",
		"[a-z].*[A-Z].*\\d", // Mixed case with numbers
	}
	
	for _, indicator := range passwordIndicators {
		if strings.Contains(pattern, indicator) {
			return true
		}
	}
	return false
}

func generateValidPassword(pattern string) string {
	base := "Password123"
	if strings.Contains(pattern, "[@$!%*?&]") || strings.Contains(pattern, "special") {
		base = "Password123!"
	}
	return fmt.Sprintf(`"%s"`, base)
}

// Generic alphanumeric detection
func isAlphaNumericPattern(pattern string) bool {
	return strings.Contains(pattern, "[a-z]") && strings.Contains(pattern, "[0-9]") ||
		   strings.Contains(pattern, "a-z") && strings.Contains(pattern, "0-9")
}

func generateAlphaNumeric(pattern string) string {
	return `"abc123"`
}

// Numeric pattern detection  
func isNumericPattern(pattern string) bool {
	return (strings.Contains(pattern, "\\d") || strings.Contains(pattern, "[0-9]")) &&
		   !strings.Contains(pattern, "[a-z]") && !strings.Contains(pattern, "a-z")
}

func generateNumericString(pattern string) string {
	return `"123456"`
}

// Enhanced IN constraint handler to handle ANY and other formats
func handleInConstraint(constraintValue, columnType string) string {
	// fmt.Printf("DEBUG: handleInConstraint called with: '%s'\n", constraintValue)
	
	constraintValue = strings.ToLower(constraintValue)
	
	// Handle ANY(ARRAY[...]) format (PostgreSQL specific)
	if strings.Contains(constraintValue, "any") && strings.Contains(constraintValue, "array") {
		return handleAnyArrayConstraint(constraintValue, columnType)
	}
	
	// Handle standard IN (...) format
	inIndex := strings.Index(constraintValue, " in ")
	if inIndex == -1 {
		// fmt.Printf("DEBUG: No ' in ' found in constraint\n")
		return generateDefaultRandomValue(columnType)
	}
	
	// Get the part after "IN"
	inPart := constraintValue[inIndex+4:]
	// fmt.Printf("DEBUG: Part after 'in ': '%s'\n", inPart)
	
	// Extract values between parentheses
	start := strings.Index(inPart, "(")
	end := strings.LastIndex(inPart, ")")
	if start == -1 || end == -1 {
		// fmt.Printf("DEBUG: No parentheses found in constraint\n")
		return generateDefaultRandomValue(columnType)
	}
	
	values := inPart[start+1 : end]
	// fmt.Printf("DEBUG: Values extracted: '%s'\n", values)
	
	// Split by comma and clean up
	validValues := strings.Split(values, ",")
	var cleanValues []string
	for _, val := range validValues {
		cleaned := strings.TrimSpace(val)
		cleaned = strings.Trim(cleaned, "'\"") // Remove quotes
		if cleaned != "" {
			cleanValues = append(cleanValues, cleaned)
		}
	}
	
	// fmt.Printf("DEBUG: Clean values: %v\n", cleanValues)
	
	if len(cleanValues) > 0 {
		// Pick the first valid value for consistency in tests
		selectedValue := cleanValues[0]
		// fmt.Printf("DEBUG: Selected value: '%s' for columnType: '%s'\n", selectedValue, columnType)
		
		if columnType == "varchar" {
			return fmt.Sprintf(`"%s"`, selectedValue)
		}
		return selectedValue
	}
	
	// fmt.Printf("DEBUG: No clean values found, using default\n")
	return generateDefaultRandomValue(columnType)
}

// Handle PostgreSQL ANY(ARRAY[...]) constraints
func handleAnyArrayConstraint(constraintValue, columnType string) string {
	// Extract values from ANY(ARRAY['val1', 'val2']) format
	arrayStart := strings.Index(constraintValue, "array[")
	if arrayStart == -1 {
		return generateDefaultRandomValue(columnType)
	}
	
	// Find the matching closing bracket
	arrayPart := constraintValue[arrayStart+6:] // Skip "array["
	bracketEnd := strings.Index(arrayPart, "]")
	if bracketEnd == -1 {
		return generateDefaultRandomValue(columnType)
	}
	
	valuesStr := arrayPart[:bracketEnd]
	
	// Split and clean values
	values := strings.Split(valuesStr, ",")
	var cleanValues []string
	for _, val := range values {
		cleaned := strings.TrimSpace(val)
		cleaned = strings.Trim(cleaned, "'\"") // Remove quotes
		// Remove type casting like ::character varying
		if idx := strings.Index(cleaned, "::"); idx != -1 {
			cleaned = cleaned[:idx]
		}
		if cleaned != "" {
			cleanValues = append(cleanValues, cleaned)
		}
	}
	
	if len(cleanValues) > 0 {
		// Pick the first valid value
		return fmt.Sprintf(`"%s"`, cleanValues[0])
	}
	
	return generateDefaultRandomValue(columnType)
}

// Length constraint handler
func handleLengthConstraint(constraintValue, columnType string) string {
	// Extract length requirements
	if strings.Contains(constraintValue, "length") {
		// Try to extract specific length requirements
		if strings.Contains(constraintValue, ">=") {
			// minimum length
			return "util.RandomString(10)" // Safe minimum
		}
		if strings.Contains(constraintValue, "<=") {
			// maximum length
			return "util.RandomString(5)" // Safe maximum
		}
		if strings.Contains(constraintValue, "=") {
			// exact length - try to extract number
			return "util.RandomString(8)" // Default safe length
		}
	}
	
	return generateDefaultRandomValue(columnType)
}

// handleRangeConstraint handles >= and <= constraints
func handleRangeConstraint(constraintValue, columnType, rangeType string) string {
	// var operator, value string
	var value string
	var isStrictInequality bool // true for > and <, false for >= and <=
	
	if strings.Contains(constraintValue, ">=") {
		parts := strings.Split(constraintValue, ">=")
		if len(parts) == 2 {
			value = strings.TrimSpace(parts[1])
			isStrictInequality = false
		}
	} else if strings.Contains(constraintValue, ">") {
		parts := strings.Split(constraintValue, ">")
		if len(parts) == 2 {
			value = strings.TrimSpace(parts[1])
			isStrictInequality = true
		}
	} else if strings.Contains(constraintValue, "<=") {
		parts := strings.Split(constraintValue, "<=")
		if len(parts) == 2 {
			value = strings.TrimSpace(parts[1])
			isStrictInequality = false
		}
	} else if strings.Contains(constraintValue, "<") {
		parts := strings.Split(constraintValue, "<")
		if len(parts) == 2 {
			value = strings.TrimSpace(parts[1])
			isStrictInequality = true
		}
	}
	
	if value == "" {
		return generateDefaultRandomValue(columnType)
	}
	
	// Parse the constraint value
	
	if columnType == "bigint" || columnType == "int" {
		if constraintVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			if rangeType == "min" {
				// For > constraints, use value + 1 as minimum
				// For >= constraints, use value as minimum
				minVal := constraintVal
				if isStrictInequality {
					minVal = constraintVal + 1
				}
				//return fmt.Sprintf("util.RandomInteger(%d, %d)", constraintVal, constraintVal+100)
				return fmt.Sprintf("util.RandomInteger(%d, %d)", minVal, minVal+100)
			} else {
				// For < constraints, use value - 1 as maximum
				// For <= constraints, use value as maximum
				maxVal := constraintVal
				if isStrictInequality {
					maxVal = constraintVal - 1
				}
				return fmt.Sprintf("util.RandomInteger(%d, %d)", maxVal-100, maxVal)
			}
		}
	} else if columnType == "real" || columnType == "float" {
		if constraintVal, err := strconv.ParseFloat(value, 64); err == nil {
			if rangeType == "min" {
				// For > constraints, use a slightly higher value
				// For >= constraints, use the exact value
				minVal := constraintVal
				if isStrictInequality {
					minVal = constraintVal + 0.01
				}
				return fmt.Sprintf("util.RandomReal(%.2f, %.2f)", minVal, minVal+100)
			} else {
				// For < constraints, use a slightly lower value
				// For <= constraints, use the exact value
				maxVal := constraintVal
				if isStrictInequality {
					maxVal = constraintVal - 0.01
				}
				return fmt.Sprintf("util.RandomReal(%.2f, %.2f)", maxVal-100, maxVal)
			}
		}
	}
	
	return generateDefaultRandomValue(columnType)
}

// handleEqualityConstraint handles = constraints
func handleEqualityConstraint(constraintValue, columnType string) string {
	parts := strings.Split(constraintValue, "=")
	if len(parts) != 2 {
		return generateDefaultRandomValue(columnType)
	}
	
	value := strings.TrimSpace(parts[1])
	value = strings.Trim(value, "'\"") // Remove quotes
	
	if columnType == "varchar" {
		return fmt.Sprintf(`"%s"`, value)
	}
	
	return value
}

// handleLikeConstraint handles LIKE constraints
func handleLikeConstraint(constraintValue, columnType string) string {
	// For email LIKE '%@%.%', generate a valid email
	if strings.Contains(constraintValue, "@") && strings.Contains(constraintValue, ".") {
		// Check if util.RandomEmail() exists, otherwise create a simple email format
		return "util.RandomEmail()"
	}
	
	// For other LIKE patterns, use default random generation
	return generateDefaultRandomValue(columnType)
}

// handleBetweenConstraint handles BETWEEN constraints
func handleBetweenConstraint(constraintValue, columnType string) string {
	// Extract values from BETWEEN clause: amount BETWEEN 0 AND 1000
	betweenIndex := strings.Index(constraintValue, "between")
	if betweenIndex == -1 {
		return generateDefaultRandomValue(columnType)
	}
	
	betweenPart := constraintValue[betweenIndex+7:] // Skip "between"
	andIndex := strings.Index(betweenPart, "and")
	if andIndex == -1 {
		return generateDefaultRandomValue(columnType)
	}
	
	minVal := strings.TrimSpace(betweenPart[:andIndex])
	maxVal := strings.TrimSpace(betweenPart[andIndex+3:]) // Skip "and"
	
	if columnType == "bigint" || columnType == "int" {
		if min, err := strconv.ParseInt(minVal, 10, 64); err == nil {
			if max, err := strconv.ParseInt(maxVal, 10, 64); err == nil {
				return fmt.Sprintf("util.RandomInteger(%d, %d)", min, max)
			}
		}
	} else if columnType == "real" || columnType == "float" {
		if min, err := strconv.ParseFloat(minVal, 64); err == nil {
			if max, err := strconv.ParseFloat(maxVal, 64); err == nil {
				return fmt.Sprintf("util.RandomReal(%.2f, %.2f)", min, max)
			}
		}
	}
	
	return generateDefaultRandomValue(columnType)
}

// handleOrConstraint handles OR conditions like "page_count is null or page_count >= 0"
func handleOrConstraint(constraintValue, columnType string) string {
	// Split by " or " to get individual conditions
	// if constraintValue == "publication_year is null or publication_year > 0"{
	// 	fmt.Printf("DEBUG: handleOrConstraint called with: '%s', columnType: '%s'\n", constraintValue, columnType)
	// }
	conditions := strings.Split(constraintValue, " or ")	
	for _, condition := range conditions {
		condition = strings.TrimSpace(condition)
		// if constraintValue == "publication_year is null or publication_year > 0"{
		// 	fmt.Printf("DEBUG: handleOrConstraint inside for loop, condition: '%s'\n", condition)
		// }

		// Check if this condition allows NULL
		if strings.Contains(condition, "is null") {

			// For OR constraints with NULL option, we can choose to either:
			// 1. Use NULL (for nullable columns)
			// 2. Use a valid non-NULL value
			// Let's choose a valid non-NULL value for test consistency
			if strings.Contains(condition, "email") && columnType == "varchar" {
				return "util.RandomEmail()"
			} else {
				// if constraintValue == "publication_year is null or publication_year > 0"{
				// 	fmt.Printf("DEBUG: handleOrConstraint inside else, condition: '%s'\n", condition)
				// }
				continue
			}
		}
		
		// Process the non-NULL condition
		if strings.Contains(condition, ">=") || strings.Contains(condition, ">") {
			// if constraintValue == "publication_year is null or publication_year > 0"{
			// 	fmt.Printf("DEBUG: handleOrConstraint inside non-NULL condition, condition: '%s'\n", condition)
			// }
			//fmt.Println(handleRangeConstraint(condition, columnType, "min"))
			return handleRangeConstraint(condition, columnType, "min")
		}
		
		if strings.Contains(condition, "<=") || strings.Contains(condition, "<") {
			return handleRangeConstraint(condition, columnType, "max")
		}
		
		if strings.Contains(condition, " in ") {
			return handleInConstraint(condition, columnType)
		}
		
		if strings.Contains(condition, "=") && !strings.Contains(condition, ">=") && !strings.Contains(condition, "<=") {
			return handleEqualityConstraint(condition, columnType)
		}
		
		if strings.Contains(condition, "like") {
			return handleLikeConstraint(condition, columnType)
		}
		
		if strings.Contains(condition, "between") {
			return handleBetweenConstraint(condition, columnType)
		}
	}
	
	// If no valid condition found, return default
	return generateDefaultRandomValue(columnType)
}

// handleAndConstraint handles AND conditions like "age >= 18 and age <= 65"
func handleAndConstraint(constraintValue, columnType string) string {
	// Split by " and " to get individual conditions
	conditions := strings.Split(constraintValue, " and ")
	
	var minVal, maxVal *float64
	// var validValues []string
	
	// Process each condition to build combined constraints
	for _, condition := range conditions {
		condition = strings.TrimSpace(condition)
		
		// Extract range constraints
		if strings.Contains(condition, ">=") {
			parts := strings.Split(condition, ">=")
			if len(parts) == 2 {
				if val, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64); err == nil {
					minVal = &val
				}
			}
		} else if strings.Contains(condition, ">") {
			parts := strings.Split(condition, ">")
			if len(parts) == 2 {
				if val, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64); err == nil {
					adjusted := val + 1 // Add 1 for > (not >=)
					minVal = &adjusted
				}
			}
		}
		
		if strings.Contains(condition, "<=") {
			parts := strings.Split(condition, "<=")
			if len(parts) == 2 {
				if val, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64); err == nil {
					maxVal = &val
				}
			}
		} else if strings.Contains(condition, "<") {
			parts := strings.Split(condition, "<")
			if len(parts) == 2 {
				if val, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64); err == nil {
					adjusted := val - 1 // Subtract 1 for < (not <=)
					maxVal = &adjusted
				}
			}
		}
		
		// Handle IN constraints within AND
		if strings.Contains(condition, " in ") {
			return handleInConstraint(condition, columnType)
		}
	}
	
	// Generate value based on combined min/max constraints
	if minVal != nil && maxVal != nil {
		// Both min and max constraints
		if columnType == "bigint" || columnType == "int" {
			return fmt.Sprintf("util.RandomInteger(%d, %d)", int64(*minVal), int64(*maxVal))
		} else if columnType == "real" || columnType == "float" {
			return fmt.Sprintf("util.RandomReal(%.2f, %.2f)", *minVal, *maxVal)
		}
	} else if minVal != nil {
		// Only min constraint
		if columnType == "bigint" || columnType == "int" {
			return fmt.Sprintf("util.RandomInteger(%d, %d)", int64(*minVal), int64(*minVal)+100)
		} else if columnType == "real" || columnType == "float" {
			return fmt.Sprintf("util.RandomReal(%.2f, %.2f)", *minVal, *minVal+100)
		}
	} else if maxVal != nil {
		// Only max constraint
		if columnType == "bigint" || columnType == "int" {
			return fmt.Sprintf("util.RandomInteger(%d, %d)", int64(*maxVal)-100, int64(*maxVal))
		} else if columnType == "real" || columnType == "float" {
			return fmt.Sprintf("util.RandomReal(%.2f, %.2f)", *maxVal-100, *maxVal)
		}
	}
	
	// Default fallback
	return generateDefaultRandomValue(columnType)
}

func generateDefaultRandomValue(columnType string) string {
	switch columnType {
	case "varchar":
		return "util.RandomName(8)"
	case "bigint":
		return "util.RandomInteger(1, 100)"
	case "real":
		return "util.RandomReal(1, 100)"
	case "timestamptz":
		return "time.Now().UTC()"
	case "date":
		return "time.Date(2025, 5, 29, 0, 0, 0, 0, time.UTC)"  // Fixed date for consistency
	case "uuid":
		return "uuid.New()"
	case "bool":
		return "true"
	default:
		return "util.RandomName(8)"
	}
}

// checkIfConstraintAllowsNull checks if a column's constraint explicitly allows NULL
func checkIfConstraintAllowsNull(columnName string, checkConstraints []dbschemareader.CheckConstraint) bool {
	for _, constraint := range checkConstraints {
		if constraint.CheckConstraintColumnName == columnName {
			constraintValue := strings.ToLower(strings.TrimSpace(constraint.CheckConstraintValue))
			// Check if constraint has "IS NULL" condition
			if strings.Contains(constraintValue, "is null") {
				return true
			}
		}
	}
	return false
}

// Enhanced CreateRandomFunction with check constraint support
func CreateRandomFunctionWithConstraints(tableX []dbschemareader.Table_Struct, i int, outputFile *os.File, userTableName string) {
	funcSig := "func createRandom" + tableX[i].FunctionSignature + "(t *testing.T"
	
	// Track used parameter names to avoid duplicates
	usedParams := make(map[string]int)
	var paramList []string
	
	for k := 0; k < len(tableX[i].ForeignKeys); k++ {
		// CHECK 1: SKIP self-referencing foreign keys in function signature
		if tableX[i].Table_name == tableX[i].ForeignKeys[k].FK_Related_TableName {
			continue
		}
		
		CamelCase := ToCamelCase(tableX[i].ForeignKeys[k].FK_Related_TableName_Singular_Object)
		baseParamName := tableX[i].ForeignKeys[k].FK_Related_SingularTableName
		
		// Handle duplicate table references
		paramName := baseParamName
		if count, exists := usedParams[baseParamName]; exists {
			usedParams[baseParamName] = count + 1
			paramName = baseParamName + strconv.Itoa(count + 1)
		} else {
			usedParams[baseParamName] = 1
		}
		
		paramList = append(paramList, paramName + " " + CamelCase)
	}
	
	// Build function signature
	for _, param := range paramList {
		funcSig = funcSig + ", " + param
	}
	
	CamelCase := ToCamelCase(tableX[i].FunctionSignature)
	funcSig = funcSig + ") " + CamelCase
	_, _ = outputFile.WriteString(funcSig + " {" + "\n")

	if tableX[i].Table_name == userTableName {
		_, _ = outputFile.WriteString(`	`+tableX[i].UserTableSpecs.AuthColumnName+`, err := util.HashPassword(util.RandomString(6))` + "\n")
		_, _ = outputFile.WriteString("	require.NoError(t, err)" + "\n")
	}
	_, _ = outputFile.WriteString("	arg := Create" + tableX[i].FunctionSignature + "Params{" + "\n")
	
	// Reset for parameter usage tracking
	usedParams = make(map[string]int)
	
	for j := 0; j < len(tableX[i].Table_Columns); j++ {
		// Skip primary keys and auto-generated columns
		if !tableX[i].IsSessionsTable {
			if tableX[i].Table_Columns[j].PrimaryFlag && strings.ToLower(tableX[i].Table_Columns[j].ColumnType) != "varchar" {
				continue
			}	
		}
		if (tableX[i].Table_Columns[j].ColumnType == "timestamptz" &&
		tableX[i].Table_Columns[j].DefaultValue == "now()") {
			continue
		}
		if (tableX[i].Table_Columns[j].ColumnType == "date" &&
		tableX[i].Table_Columns[j].DefaultValue == "CURRENT_DATE") {
			continue
		}
		
		if tableX[i].Table_Columns[j].ForeignFlag {
			// Handle foreign keys (existing logic)
			for k := 0; k < len(tableX[i].ForeignKeys); k++ {
				if tableX[i].ForeignKeys[k].FK_Column == tableX[i].Table_Columns[j].Column_name {
					// Handle self-referencing foreign keys
					if tableX[i].Table_name == tableX[i].ForeignKeys[k].FK_Related_TableName {
						// For self-reference, set to NULL
						if strings.Contains(tableX[i].Table_Columns[j].ColumnType, "uuid") {
							_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    	pgtype.UUID{Valid: false},"+"\n")
						} else {
							_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    	pgtype.Text{Valid: false},"+"\n")
						}
					} else {
						// Normal foreign key handling (existing logic)
						baseParamName := tableX[i].ForeignKeys[k].FK_Related_SingularTableName
						paramName := baseParamName
						if count, exists := usedParams[baseParamName]; exists {
							usedParams[baseParamName] = count + 1
							paramName = baseParamName + strconv.Itoa(count + 1)
						} else {
							usedParams[baseParamName] = 1
						}
						
						if tableX[i].Table_Columns[j].ColumnType == "uuid" && tableX[i].Table_Columns[j].Not_Null {
							FormatedFieldName := FormatFieldName(tableX[i].ForeignKeys[k].FK_Related_Table_Column)
							_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    	" + paramName + "." + FormatedFieldName+","+"\n")
						} else if tableX[i].Table_Columns[j].ColumnType == "uuid" && !tableX[i].Table_Columns[j].Not_Null {
							FormatedFieldName := FormatFieldName(tableX[i].ForeignKeys[k].FK_Related_Table_Column)
							_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    	pgtype.UUID{Bytes: " + paramName + "." + FormatedFieldName + ", Valid: true},"+"\n")
						}
					}
				}
			}
		} else {
			// Enhanced column value generation with check constraint support
				// if tableX[i].Table_name == "orders" && tableX[i].Table_Columns[j].Column_name == "totalamount"{
				// 	fmt.Println("tableX[i].Table_name", tableX[i].Table_name)
				// 	fmt.Println("tableX[i].Table_Columns[j].Column_name)", tableX[i].Table_Columns[j].Column_name)
				// 	fmt.Println("tableX[i].Table_Columns[j].ColumnType)", tableX[i].Table_Columns[j].ColumnType)
				// 	fmt.Println("tableX[i].Table_Columns[j].Not_Null)", tableX[i].Table_Columns[j].Not_Null)
				// 	time.Sleep(30 * time.Second)
				// }
			if tableX[i].Table_Columns[j].ColumnType == "varchar" && tableX[i].Table_Columns[j].Not_Null {
				if tableX[i].Table_name == userTableName && tableX[i].Table_Columns[j].Column_name == tableX[i].UserTableSpecs.AuthColumnName{
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    "+tableX[i].UserTableSpecs.AuthColumnName+"," + "\n")
				} else {
					// FIXED: Use check constraint aware value generation
					validValue := generateValidValueForCheckConstraint(tableX[i].Table_name, tableX[i].Table_Columns[j].Column_name, tableX[i].Table_Columns[j].ColumnType, tableX[i].CheckConstraints)
					// fmt.Printf("FINAL DEBUG: Using value %s for %s.%s\n", validValue, tableX[i].Table_name, tableX[i].Table_Columns[j].Column_name)
					
					// CRITICAL FIX: Check if the validValue contains util. functions and handle accordingly
					if strings.Contains(validValue, "util.") {
						// This is a function call, use as-is
						_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    " + validValue + "," + "\n")
					} else {
						// This is a literal string value, keep the quotes
						_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    " + validValue + "," + "\n")
					}
				}
			} else if tableX[i].Table_Columns[j].ColumnType == "varchar" && !tableX[i].Table_Columns[j].Not_Null {
				validValue := generateValidValueForCheckConstraint(tableX[i].Table_name, tableX[i].Table_Columns[j].Column_name, tableX[i].Table_Columns[j].ColumnType, tableX[i].CheckConstraints)
				// For nullable varchar with constraints, wrap in pgtype
				if strings.Contains(validValue, "util.") {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    pgtype.Text{String: " + validValue + ", Valid: true}," + "\n")
				} else {
					// Remove quotes from validValue since we're wrapping it in pgtype.Text
					cleanValue := strings.Trim(validValue, `"`)
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    pgtype.Text{String: \"" + cleanValue + "\", Valid: true}," + "\n")
				}
			} else if tableX[i].Table_Columns[j].ColumnType == "bigint" && tableX[i].Table_Columns[j].Not_Null {
				validValue := generateValidValueForCheckConstraint(tableX[i].Table_name, tableX[i].Table_Columns[j].Column_name, tableX[i].Table_Columns[j].ColumnType, tableX[i].CheckConstraints)
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    " + validValue + "," + "\n")
			} else if tableX[i].Table_Columns[j].ColumnType == "bigint" && !tableX[i].Table_Columns[j].Not_Null {
				validValue := generateValidValueForCheckConstraint(tableX[i].Table_name, tableX[i].Table_Columns[j].Column_name, tableX[i].Table_Columns[j].ColumnType, tableX[i].CheckConstraints)
				
				// Special handling for OR constraints with NULL option
				constraintAllowsNull := checkIfConstraintAllowsNull(tableX[i].Table_Columns[j].Column_name, tableX[i].CheckConstraints)
				
				if constraintAllowsNull {
					if strings.Contains(validValue, "util.") {
						_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    pgtype.Int8{Int64: " + validValue + ", Valid: true}," + "\n")
					} else {
						if numVal, err := strconv.ParseInt(validValue, 10, 64); err == nil {
							_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    pgtype.Int8{Int64: " + strconv.FormatInt(numVal, 10) + ", Valid: true}," + "\n")
						} else {
							_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    pgtype.Int8{Int64: " + validValue + ", Valid: true}," + "\n")
						}
					}
				} else {
					if strings.Contains(validValue, "util.") {
						_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    pgtype.Int8{Int64: " + validValue + ", Valid: true}," + "\n")
					} else {
						_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    pgtype.Int8{Int64: " + validValue + ", Valid: true}," + "\n")
					}
				}
			} else if tableX[i].Table_Columns[j].ColumnType == "real" && tableX[i].Table_Columns[j].Not_Null {
				// Check if this column has an arithmetic constraint
				hasArithmeticConstraint := checkForArithmeticConstraint(tableX[i].Table_name, tableX[i].Table_Columns[j].Column_name, tableX[i].CheckConstraints)
				if hasArithmeticConstraint {
					// For arithmetic constraints, we'll set this after struct creation
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    0, // Will be calculated after struct creation" + "\n")
				} else {
					// Normal constraint handling
					validValue := generateValidValueForCheckConstraint(tableX[i].Table_name, tableX[i].Table_Columns[j].Column_name, tableX[i].Table_Columns[j].ColumnType, tableX[i].CheckConstraints)
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    " + validValue + "," + "\n")
				}
			} else if tableX[i].Table_Columns[j].ColumnType == "real" && !tableX[i].Table_Columns[j].Not_Null {
				validValue := generateValidValueForCheckConstraint(tableX[i].Table_name, tableX[i].Table_Columns[j].Column_name, tableX[i].Table_Columns[j].ColumnType, tableX[i].CheckConstraints)
				
				constraintAllowsNull := checkIfConstraintAllowsNull(tableX[i].Table_Columns[j].Column_name, tableX[i].CheckConstraints)
				
				if constraintAllowsNull {
					if strings.Contains(validValue, "util.") {
						_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    pgtype.Float4{Float32: " + validValue + ", Valid: true}," + "\n")
					} else {
						if numVal, err := strconv.ParseFloat(validValue, 32); err == nil {
							_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    pgtype.Float4{Float32: " + strconv.FormatFloat(numVal, 'f', 2, 32) + ", Valid: true}," + "\n")
						} else {
							_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    pgtype.Float4{Float32: " + validValue + ", Valid: true}," + "\n")
						}
					}
				} else {
					if strings.Contains(validValue, "util.") {
						_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    pgtype.Float4{Float32: " + validValue + ", Valid: true}," + "\n")
					} else {
						_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    pgtype.Float4{Float32: " + validValue + ", Valid: true}," + "\n")
					}
				}
			} else {
				// For other column types, use existing logic
				if tableX[i].Table_Columns[j].ColumnType == "timestamptz" && tableX[i].Table_Columns[j].Not_Null {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    time.Now().UTC()," + "\n")
				} else if tableX[i].Table_Columns[j].ColumnType == "timestamptz" && !tableX[i].Table_Columns[j].Not_Null {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true}," + "\n")
				} else if tableX[i].Table_Columns[j].ColumnType == "date" && tableX[i].Table_Columns[j].Not_Null {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    time.Date(2025, 5, 29, 0, 0, 0, 0, time.UTC)," + "\n")
				} else if tableX[i].Table_Columns[j].ColumnType == "date" && !tableX[i].Table_Columns[j].Not_Null {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    pgtype.Date{Time: time.Date(2025, 5, 29, 0, 0, 0, 0, time.UTC), Valid: true}," + "\n")
				} else if tableX[i].Table_Columns[j].ColumnType == "uuid" && tableX[i].Table_Columns[j].Not_Null {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    uuid.New()," + "\n")
				} else if tableX[i].Table_Columns[j].ColumnType == "uuid" && !tableX[i].Table_Columns[j].Not_Null {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    pgtype.UUID{Bytes: uuid.New(), Valid: true}," + "\n")
				} else if tableX[i].Table_Columns[j].ColumnType == "bool" && tableX[i].Table_Columns[j].Not_Null {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    true," + "\n")
				} else if tableX[i].Table_Columns[j].ColumnType == "bool" && !tableX[i].Table_Columns[j].Not_Null {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    pgtype.Bool{Bool: true, Valid: true}," + "\n")
				}
			}
		}
	}
	
	// Rest of the function remains the same...
	_, _ = outputFile.WriteString("	}" + "\n")

	// After the struct creation (after the closing brace)
	// ✅ NEW: Handle arithmetic constraints after struct creation
	arithmeticConstraints := getArithmeticConstraints(tableX[i].Table_name, tableX[i].CheckConstraints)
	for _, constraint := range arithmeticConstraints {
		expression := generateArithmeticAssignment(constraint, tableX[i].Table_Columns)
		if expression != "" {
			_, _ = outputFile.WriteString("	" + expression + "\n")
		}
	}

	_, _ = outputFile.WriteString("	" + tableX[i].OutputFileName + ", err := testStore.Create" + tableX[i].FunctionSignature + "(context.Background(), arg)" + "\n")
	_, _ = outputFile.WriteString("	require.NoError(t, err)" + "\n")
	_, _ = outputFile.WriteString("	require.NotEmpty(t, " + tableX[i].OutputFileName + ")" + "\n")
	
	// Validation logic remains the same...
	for j := 0; j < len(tableX[i].Table_Columns); j++ {
		if !tableX[i].IsSessionsTable {
			if tableX[i].Table_Columns[j].PrimaryFlag && strings.ToLower(tableX[i].Table_Columns[j].ColumnType) != "varchar" {
				continue
			}	
		}
		if (tableX[i].Table_Columns[j].ColumnType == "timestamptz" &&
		tableX[i].Table_Columns[j].DefaultValue == "now()") {
			continue
		}
		if (tableX[i].Table_Columns[j].ColumnType == "date" &&
		tableX[i].Table_Columns[j].DefaultValue == "CURRENT_DATE") {
			continue
		}

		if tableX[i].Table_Columns[j].ColumnType == "timestamptz" {
			if tableX[i].Table_Columns[j].Not_Null {
				_, _ = outputFile.WriteString("	require.WithinDuration(t, arg." + tableX[i].Table_Columns[j].ColumnNameParams + ", " + tableX[i].OutputFileName + "." + tableX[i].Table_Columns[j].ColumnNameParams + ", time.Second" + ")" + "\n")
			} else {
				_, _ = outputFile.WriteString("	require.WithinDuration(t, arg." + tableX[i].Table_Columns[j].ColumnNameParams + ".Time, " + tableX[i].OutputFileName + "." + tableX[i].Table_Columns[j].ColumnNameParams + ".Time, time.Second" + ")" + "\n")
			}
		} else {
			_, _ = outputFile.WriteString("	require.Equal(t, arg." + tableX[i].Table_Columns[j].ColumnNameParams + ", " + tableX[i].OutputFileName + "." + tableX[i].Table_Columns[j].ColumnNameParams + ")" + "\n")
		}
	}
	_, _ = outputFile.WriteString("	return " + tableX[i].OutputFileName + "\n")
	_, _ = outputFile.WriteString("}" + "\n")
	_, _ = outputFile.WriteString("\n")
}

func CreateRandomFunction(tableX []dbschemareader.Table_Struct, i int, outputFile *os.File, userTableName string) {
	funcSig := "func createRandom" + tableX[i].FunctionSignature + "(t *testing.T"
	
	// Track used parameter names to avoid duplicates
	usedParams := make(map[string]int)
	var paramList []string
	
	for k := 0; k < len(tableX[i].ForeignKeys); k++ {
		// CHECK 1: SKIP self-referencing foreign keys in function signature
		if tableX[i].Table_name == tableX[i].ForeignKeys[k].FK_Related_TableName {
			continue
		}
		
		CamelCase := ToCamelCase(tableX[i].ForeignKeys[k].FK_Related_TableName_Singular_Object)
		baseParamName := tableX[i].ForeignKeys[k].FK_Related_SingularTableName
		
		// Handle duplicate table references
		paramName := baseParamName
		if count, exists := usedParams[baseParamName]; exists {
			usedParams[baseParamName] = count + 1
			paramName = baseParamName + strconv.Itoa(count + 1)
		} else {
			usedParams[baseParamName] = 1
		}
		
		paramList = append(paramList, paramName + " " + CamelCase)
	}
	
	// Build function signature
	for _, param := range paramList {
		funcSig = funcSig + ", " + param
	}
	
	CamelCase := ToCamelCase(tableX[i].FunctionSignature)
	funcSig = funcSig + ") " + CamelCase
	_, _ = outputFile.WriteString(funcSig + " {" + "\n")

	if tableX[i].Table_name == userTableName {
		_, _ = outputFile.WriteString(`	`+tableX[i].UserTableSpecs.AuthColumnName+`, err := util.HashPassword(util.RandomString(6))` + "\n")
		_, _ = outputFile.WriteString("	require.NoError(t, err)" + "\n")
	}
	_, _ = outputFile.WriteString("	arg := Create" + tableX[i].FunctionSignature + "Params{" + "\n")
	
	// Reset for parameter usage tracking
	usedParams = make(map[string]int)
	
	for j := 0; j < len(tableX[i].Table_Columns); j++ {
		////////////////////This filter must match with sqlcq filters//////////////////////////
		////Since this is NOT the session management table we are excluding the Primary Key////
		///////////////////////////////////////////////////////////////////////////////////////
		if tableX[i].Table_Columns[j].PrimaryFlag{
			continue
		}
		if (tableX[i].Table_Columns[j].ColumnType == "timestamptz" &&
		tableX[i].Table_Columns[j].DefaultValue == "now()") {
			continue
		}
		if (tableX[i].Table_Columns[j].ColumnType == "date" &&
		tableX[i].Table_Columns[j].DefaultValue == "CURRENT_DATE") {
			continue
		}
		if tableX[i].Table_Columns[j].ForeignFlag {
			for k := 0; k < len(tableX[i].ForeignKeys); k++ {
				if tableX[i].ForeignKeys[k].FK_Column == tableX[i].Table_Columns[j].Column_name {
					// HANDLE self-referencing foreign keys
					if tableX[i].Table_name == tableX[i].ForeignKeys[k].FK_Related_TableName {
						// For self-reference, set to NULL
						if strings.Contains(tableX[i].Table_Columns[j].ColumnType, "uuid") {
							_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    	pgtype.UUID{Valid: false},"+"\n")
						} else {
							_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    	pgtype.Text{Valid: false},"+"\n")
						}
					} else {
						// Determine parameter name (handle duplicates)
						baseParamName := tableX[i].ForeignKeys[k].FK_Related_SingularTableName
						paramName := baseParamName
						if count, exists := usedParams[baseParamName]; exists {
							usedParams[baseParamName] = count + 1
							paramName = baseParamName + strconv.Itoa(count + 1)
						} else {
							usedParams[baseParamName] = 1
						}
						
						if tableX[i].Table_Columns[j].ColumnType == "uuid" && tableX[i].Table_Columns[j].Not_Null {
							// Normal foreign key handling
							FormatedFieldName := FormatFieldName(tableX[i].ForeignKeys[k].FK_Related_Table_Column)
							_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    	" + paramName + "." + FormatedFieldName+","+"\n")
						} else if tableX[i].Table_Columns[j].ColumnType == "uuid" && !tableX[i].Table_Columns[j].Not_Null {
							// Nullable foreign key handling
							FormatedFieldName := FormatFieldName(tableX[i].ForeignKeys[k].FK_Related_Table_Column)
							_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    	pgtype.UUID{Bytes: " + paramName + "." + FormatedFieldName + ", Valid: true},"+"\n")
						}
					}
				}
			}
		} else {
			if tableX[i].Table_Columns[j].ColumnType == "varchar" && tableX[i].Table_Columns[j].Not_Null {
				if tableX[i].Table_name == userTableName && tableX[i].Table_Columns[j].Column_name == tableX[i].UserTableSpecs.AuthColumnName{
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    "+tableX[i].UserTableSpecs.AuthColumnName+"," + "\n")
				}else{
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    util.RandomName(8)," + "\n")
				}
			} else if tableX[i].Table_Columns[j].ColumnType == "varchar" && !tableX[i].Table_Columns[j].Not_Null {
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    pgtype.Text{String: util.RandomName(8), Valid: true}," + "\n")
			} else if tableX[i].Table_Columns[j].ColumnType == "bigint" && tableX[i].Table_Columns[j].Not_Null {
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    util.RandomInteger(1, 100)," + "\n")
			} else if tableX[i].Table_Columns[j].ColumnType == "bigint" && !tableX[i].Table_Columns[j].Not_Null {
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    pgtype.Int8{Int64: util.RandomInteger(1, 100), Valid: true}," + "\n")
			} else if tableX[i].Table_Columns[j].ColumnType == "real" && tableX[i].Table_Columns[j].Not_Null {
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    util.RandomReal(1, 100)," + "\n")
			} else if tableX[i].Table_Columns[j].ColumnType == "real" && !tableX[i].Table_Columns[j].Not_Null {
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    pgtype.Float4{Float32: util.RandomReal(1, 100), Valid: true}," + "\n")
			} else if tableX[i].Table_Columns[j].ColumnType == "timestamptz" && tableX[i].Table_Columns[j].Not_Null {
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    time.Now().UTC()," + "\n")
			} else if tableX[i].Table_Columns[j].ColumnType == "timestamptz" && !tableX[i].Table_Columns[j].Not_Null {
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true}," + "\n")
			} else if tableX[i].Table_Columns[j].ColumnType == "date" && tableX[i].Table_Columns[j].Not_Null {
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    time.Now().UTC()," + "\n")
			} else if tableX[i].Table_Columns[j].ColumnType == "date" && !tableX[i].Table_Columns[j].Not_Null {
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    pgtype.Date{Time: time.Now().UTC(), Valid: true}," + "\n")
			} else if tableX[i].Table_Columns[j].ColumnType == "uuid" && tableX[i].Table_Columns[j].Not_Null {
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    uuid.New()," + "\n")
			} else if tableX[i].Table_Columns[j].ColumnType == "uuid" && !tableX[i].Table_Columns[j].Not_Null {
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    pgtype.UUID{Bytes: uuid.New(), Valid: true}," + "\n")
			} else if tableX[i].Table_Columns[j].ColumnType == "bool" && tableX[i].Table_Columns[j].Not_Null {
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    true," + "\n")
			} else if tableX[i].Table_Columns[j].ColumnType == "bool" && !tableX[i].Table_Columns[j].Not_Null {
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[j].ColumnNameParams + ":    pgtype.Bool{Bool: true, Valid: true}," + "\n")
			}
		}
	}
	_, _ = outputFile.WriteString("	}" + "\n")
	_, _ = outputFile.WriteString("	" + tableX[i].OutputFileName + ", err := testStore.Create" + tableX[i].FunctionSignature + "(context.Background(), arg)" + "\n")
	_, _ = outputFile.WriteString("	require.NoError(t, err)" + "\n")
	_, _ = outputFile.WriteString("	require.NotEmpty(t, " + tableX[i].OutputFileName + ")" + "\n")
	for j := 0; j < len(tableX[i].Table_Columns); j++ {
		if tableX[i].Table_Columns[j].PrimaryFlag{
			continue
		}
		if (tableX[i].Table_Columns[j].ColumnType == "timestamptz" &&
		tableX[i].Table_Columns[j].DefaultValue == "now()") {
			continue
		}
		if (tableX[i].Table_Columns[j].ColumnType == "date" &&
		tableX[i].Table_Columns[j].DefaultValue == "CURRENT_DATE") {
			continue
		}

		// In CreateRandomFunction, around line where timestamp validation happens:
		if tableX[i].Table_Columns[j].ColumnType == "timestamptz" {
			if tableX[i].Table_Columns[j].Not_Null {
				// Non-nullable: compare time.Time with time.Time
				_, _ = outputFile.WriteString("	require.WithinDuration(t, arg." + tableX[i].Table_Columns[j].ColumnNameParams + ", " + tableX[i].OutputFileName + "." + tableX[i].Table_Columns[j].ColumnNameParams + ", time.Second" + ")" + "\n")
			} else {
				// Nullable: compare pgtype.Timestamptz.Time with pgtype.Timestamptz.Time
				_, _ = outputFile.WriteString("	require.WithinDuration(t, arg." + tableX[i].Table_Columns[j].ColumnNameParams + ".Time, " + tableX[i].OutputFileName + "." + tableX[i].Table_Columns[j].ColumnNameParams + ".Time, time.Second" + ")" + "\n")
			}
		} else {
			_, _ = outputFile.WriteString("	require.Equal(t, arg." + tableX[i].Table_Columns[j].ColumnNameParams + ", " + tableX[i].OutputFileName + "." + tableX[i].Table_Columns[j].ColumnNameParams + ")" + "\n")
		}
	}
	_, _ = outputFile.WriteString("	return " + tableX[i].OutputFileName + "\n")
	_, _ = outputFile.WriteString("}" + "\n")
	_, _ = outputFile.WriteString("\n")
}

func printTestFuncForCreate(tableX []dbschemareader.Table_Struct, i int, fk_HierarchyX []dbschemareader.FK_Hierarchy, outputFile *os.File) {
	var fkVarMap = make(map[string]string)
	_, _ = outputFile.WriteString("func TestCreate" + tableX[i].FunctionSignature + "(t *testing.T) {" + "\n")
	for k := 0; k < len(fk_HierarchyX); k++ {
		if fk_HierarchyX[k].TableName == tableX[i].Table_name {
			for l := len(fk_HierarchyX[k].RelatedTablesLevels) - 1; l >= 0; l-- {
				for m := 0; m < len(fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList); m++ {
					// CHECK 1: Skip self-referencing foreign keys in dependency creation
					if fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {
						continue
					}

					varName := fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName+"_fk_"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName+strconv.Itoa(k) + strconv.Itoa(l) + strconv.Itoa(m)
					key := fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName+"_fk_"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName
					// Only store if key doesn't exist
					if _, exists := fkVarMap[key]; !exists {
						fkVarMap[key] = varName
						_, _ = outputFile.WriteString("	" + varName+ " := createRandom" + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName_Singular_Object + "(t")
					} else{
						continue
					}
					for g := 0; g < len(fk_HierarchyX); g++ {
						if fk_HierarchyX[g].TableName == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {
							for h := 0; h < len(fk_HierarchyX[g].RelatedTablesLevels); h++ {
								for z := 0; z < len(fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList); z++ {
									// CHECK 2: Skip self-referencing in nested dependencies
									if fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName == fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_TableName {
										continue
									}
				
									key := fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName
									if val, ok := fkVarMap[key]; ok {
										_, _ = outputFile.WriteString(", " + val)
									}
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
					// CHECK 3: Skip self-referencing in main function parameters
					if fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName == fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_TableName {
						continue
					}

					key := fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName
					if val, ok := fkVarMap[key]; ok {
						_, _ = outputFile.WriteString(", " + val)
					}
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
							// CHECK 1: Skip self-referencing foreign keys in dependency creation
							if fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {
								continue
							}
							
							varName := fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName+"_fk_"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName+strconv.Itoa(k) + strconv.Itoa(l) + strconv.Itoa(m)
							key := fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName+"_fk_"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName
							// Only store if key doesn't exist
							if _, exists := fkVarMap[key]; !exists {
								fkVarMap[key] = varName
								_, _ = outputFile.WriteString("	" + varName+ " := createRandom" + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName_Singular_Object + "(t")		
							} else{
								continue
							}
							for g := 0; g < len(fk_HierarchyX); g++ {
								if fk_HierarchyX[g].TableName == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {
									for h := 0; h < len(fk_HierarchyX[g].RelatedTablesLevels); h++ {
										for z := 0; z < len(fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList); z++ {
											// CHECK 2: Skip self-referencing in nested dependencies
											if fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName == fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_TableName {
												continue
											}

											key := fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName
											if val, ok := fkVarMap[key]; ok {
												_, _ = outputFile.WriteString(", " + val)
											}		
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
								// CHECK 3: Skip self-referencing in main function parameters
								if fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName == fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_TableName {
									continue
								}

								key := fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName
								if val, ok := fkVarMap[key]; ok {
									_, _ = outputFile.WriteString(", " + val)
								}
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
			
			// Fixed validation logic with proper pgtype handling
			for h := 0; h < len(tableX[i].Table_Columns); h++ {
				if tableX[i].Table_Columns[h].ColumnType == "timestamptz" {
					if tableX[i].Table_Columns[h].Not_Null {
						// Non-nullable timestamptz: both are time.Time
						_, _ = outputFile.WriteString("	require.WithinDuration(t, " + tableX[i].OutputFileName + "1." + tableX[i].Table_Columns[h].ColumnNameParams + ", " + tableX[i].OutputFileName + "2." + tableX[i].Table_Columns[h].ColumnNameParams + ", time.Second)" + "\n")
					} else {
						// Nullable timestamptz: both are pgtype.Timestamptz
						_, _ = outputFile.WriteString("	require.WithinDuration(t, " + tableX[i].OutputFileName + "1." + tableX[i].Table_Columns[h].ColumnNameParams + ".Time, " + tableX[i].OutputFileName + "2." + tableX[i].Table_Columns[h].ColumnNameParams + ".Time, time.Second)" + "\n")
					}
				}  else if tableX[i].Table_Columns[h].ColumnType == "date" {
					if tableX[i].Table_Columns[h].Not_Null {
						// Non-nullable date: both are time.Time
						_, _ = outputFile.WriteString("	require.Equal(t, " + tableX[i].OutputFileName + "1." + tableX[i].Table_Columns[h].ColumnNameParams + ", " + tableX[i].OutputFileName + "2." + tableX[i].Table_Columns[h].ColumnNameParams + ")" + "\n")
					} else {
						// Nullable date: both are pgtype.Date
						_, _ = outputFile.WriteString("	require.Equal(t, " + tableX[i].OutputFileName + "1." + tableX[i].Table_Columns[h].ColumnNameParams + ", " + tableX[i].OutputFileName + "2." + tableX[i].Table_Columns[h].ColumnNameParams + ")" + "\n")
					}
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

// Enhanced printTestFuncForReadList to handle composite unique constraints
func printTestFuncForReadList(tableX []dbschemareader.Table_Struct, i int, fk_HierarchyX []dbschemareader.FK_Hierarchy, outputFile *os.File) {
	var fkVarMap = make(map[string]string)
	_, _ = outputFile.WriteString("func TestList" + tableX[i].FunctionSignature2 + "(t *testing.T) {" + "\n")
	
	// Check if table has composite unique constraints
	hasCompositeUniqueConstraints := len(tableX[i].CompositeUniqueConstraints) > 0
	
	if hasCompositeUniqueConstraints {
		// For tables with composite unique constraints, create dependencies inside the loop
		// to ensure different foreign key combinations
		_, _ = outputFile.WriteString("	for i := 0; i < 10; i++ {" + "\n")
		
		// Create fresh dependencies for each iteration
		for k := 0; k < len(fk_HierarchyX); k++ {
			if fk_HierarchyX[k].TableName == tableX[i].Table_name {
				for l := len(fk_HierarchyX[k].RelatedTablesLevels) - 1; l >= 0; l-- {
					for m := 0; m < len(fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList); m++ {
						if fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {
							continue
						}

						varName := fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName+"_fk_"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName+strconv.Itoa(k) + strconv.Itoa(l) + strconv.Itoa(m)
						key := fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName+"_fk_"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName
						if _, exists := fkVarMap[key]; !exists {
							fkVarMap[key] = varName
							_, _ = outputFile.WriteString("	" + varName+ " := createRandom" + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName_Singular_Object + "(t")
						} else{
							continue
						}
						for g := 0; g < len(fk_HierarchyX); g++ {
							if fk_HierarchyX[g].TableName == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {
								for h := 0; h < len(fk_HierarchyX[g].RelatedTablesLevels); h++ {
									for z := 0; z < len(fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList); z++ {
										if fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName == fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_TableName {
											continue
										}
										key := fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName
										if val, ok := fkVarMap[key]; ok {
											_, _ = outputFile.WriteString(", " + val)
										}		
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
		
		// Create the main record with fresh dependencies
		_, _ = outputFile.WriteString("		createRandom" + tableX[i].FunctionSignature + "(t")
		for g := 0; g < len(fk_HierarchyX); g++ {
			if fk_HierarchyX[g].TableName == tableX[i].Table_name {
				for h := 0; h < len(fk_HierarchyX[g].RelatedTablesLevels); h++ {
					for z := 0; z < len(fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList); z++ {
						if fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName == fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_TableName {
							continue
						}
						key := fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName
						if val, ok := fkVarMap[key]; ok {
							_, _ = outputFile.WriteString(", " + val)
						}		
					}
					if h == 0 {
						break
					}
				}
			}
		}
		_, _ = outputFile.WriteString(")" + "\n")
		_, _ = outputFile.WriteString("	}" + "\n")
		
	} else {
		// Original logic for tables WITHOUT composite unique constraints
		// Create dependencies once outside the loop
		for k := 0; k < len(fk_HierarchyX); k++ {
			if fk_HierarchyX[k].TableName == tableX[i].Table_name {
				for l := len(fk_HierarchyX[k].RelatedTablesLevels) - 1; l >= 0; l-- {
					for m := 0; m < len(fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList); m++ {
						if fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {
							continue
						}

						varName := fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName+"_fk_"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName+strconv.Itoa(k) + strconv.Itoa(l) + strconv.Itoa(m)
						key := fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName+"_fk_"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName
						if _, exists := fkVarMap[key]; !exists {
							fkVarMap[key] = varName
							_, _ = outputFile.WriteString("	" + varName+ " := createRandom" + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName_Singular_Object + "(t")
						} else{
							continue
						}
						for g := 0; g < len(fk_HierarchyX); g++ {
							if fk_HierarchyX[g].TableName == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {
								for h := 0; h < len(fk_HierarchyX[g].RelatedTablesLevels); h++ {
									for z := 0; z < len(fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList); z++ {
										if fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName == fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_TableName {
											continue
										}
										key := fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName
										if val, ok := fkVarMap[key]; ok {
											_, _ = outputFile.WriteString(", " + val)
										}		
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
		
		// Create 10 test records with same dependencies
		_, _ = outputFile.WriteString("	for i := 0; i < 10; i++ {" + "\n")
		_, _ = outputFile.WriteString("		createRandom" + tableX[i].FunctionSignature + "(t")
		for g := 0; g < len(fk_HierarchyX); g++ {
			if fk_HierarchyX[g].TableName == tableX[i].Table_name {
				for h := 0; h < len(fk_HierarchyX[g].RelatedTablesLevels); h++ {
					for z := 0; z < len(fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList); z++ {
						if fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName == fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_TableName {
							continue
						}
						key := fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName
						if val, ok := fkVarMap[key]; ok {
							_, _ = outputFile.WriteString(", " + val)
						}		
					}
					if h == 0 {
						break
					}
				}
			}
		}
		_, _ = outputFile.WriteString(")" + "\n")
		_, _ = outputFile.WriteString("	}" + "\n")
	}

	// Rest of the function remains the same (List parameters and validation)
	_, _ = outputFile.WriteString("	arg := List" + tableX[i].FunctionSignature2 + "Params{" + "\n")
	for g := 0; g < len(tableX[i].Table_Columns); g++ {
		// if tableX[i].Table_Columns[g].ForeignFlag && !tableX[i].Table_Columns[g].Not_Null {
		// 	for r := 0; r < len(tableX[i].ForeignKeys); r++ {
		// 		if tableX[i].ForeignKeys[r].FK_Column == tableX[i].Table_Columns[g].Column_name {
		// 			if tableX[i].Table_name == tableX[i].ForeignKeys[r].FK_Related_TableName {
		// 				if strings.Contains(tableX[i].Table_Columns[g].ColumnType, "uuid") {
		// 					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[g].ColumnNameParams + ": pgtype.UUID{Valid: false},"+"\n")
		// 				} else {
		// 					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[g].ColumnNameParams + ": pgtype.Text{Valid: false},"+"\n")
		// 				}
		// 			} else {
		// 				FormatedFieldName := FormatFieldName(tableX[i].ForeignKeys[r].FK_Related_Table_Column)
		// 				if hasCompositeUniqueConstraints {
		// 					// For composite unique tables, we can't use shared FK vars since they're created in loop
		// 					// Use NULL for list queries instead
		// 					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[g].ColumnNameParams + ": pgtype.UUID{Valid: false},"+"\n")
		// 				} else {
		// 					key := tableX[i].Table_name+"_fk_"+tableX[i].ForeignKeys[r].FK_Related_SingularTableName
		// 					if val, ok := fkVarMap[key]; ok {
		// 						_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[g].ColumnNameParams + ": pgtype.UUID{Bytes: " + val + "." + FormatedFieldName + ", Valid: true},"+"\n")
		// 					}						
		// 				}
		// 			}
		// 		}
		// 	}
		// }
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
	_, _ = outputFile.WriteString("	}" + "\n")
	_, _ = outputFile.WriteString("}" + "\n")
	_, _ = outputFile.WriteString("\n")
}

// printTestFuncForUpdateWithConstraints - Enhanced version with check constraint support
func printTestFuncForUpdateWithConstraints(tableX []dbschemareader.Table_Struct, i int, fk_HierarchyX []dbschemareader.FK_Hierarchy, outputFile *os.File, userTableName string) {
	var fkVarMap = make(map[string]string)
	_, _ = outputFile.WriteString("func TestUpdate" + tableX[i].FunctionSignature + "(t *testing.T) {" + "\n")
	
	// Create dependencies - skip self-referencing (existing logic)
	for k := 0; k < len(fk_HierarchyX); k++ {
		if fk_HierarchyX[k].TableName == tableX[i].Table_name {
			for l := len(fk_HierarchyX[k].RelatedTablesLevels) - 1; l >= 0; l-- {
				for m := 0; m < len(fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList); m++ {
					// CHECK 1: Skip self-referencing foreign keys in dependency creation
					if fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {
						continue
					}
					varName := fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName+"_fk_"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName+strconv.Itoa(k) + strconv.Itoa(l) + strconv.Itoa(m)
					key := fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName+"_fk_"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName
					// Only store if key doesn't exist
					if _, exists := fkVarMap[key]; !exists {
						fkVarMap[key] = varName
						_, _ = outputFile.WriteString("	" + varName+ " := createRandom" + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName_Singular_Object + "(t")
					} else{
						continue
					}
					for g := 0; g < len(fk_HierarchyX); g++ {
						if fk_HierarchyX[g].TableName == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {
							for h := 0; h < len(fk_HierarchyX[g].RelatedTablesLevels); h++ {
								for z := 0; z < len(fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList); z++ {
									// CHECK 2: Skip self-referencing in nested dependencies
									if fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName == fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_TableName {
										continue
									}
									key := fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName
									if val, ok := fkVarMap[key]; ok {
										_, _ = outputFile.WriteString(", " + val)
									}		
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
	
	// Create record to update (existing logic)
	_, _ = outputFile.WriteString("	" + tableX[i].OutputFileName + "1 := createRandom" + tableX[i].FunctionSignature + "(t")
	for g := 0; g < len(fk_HierarchyX); g++ {
		if fk_HierarchyX[g].TableName == tableX[i].Table_name {
			if len(fk_HierarchyX[g].RelatedTablesLevels) > 0 {
				for h := 0; h < len(fk_HierarchyX[g].RelatedTablesLevels); h++ {
					for z := 0; z < len(fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList); z++ {
						// CHECK 3: Skip self-referencing in main function parameters
						if fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName == fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_TableName {
							continue
						}
						key := fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName
						if val, ok := fkVarMap[key]; ok {
							_, _ = outputFile.WriteString(", " + val)
						}		
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
	
	// Get primary key column name (existing logic)
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
	
	// Special handling for users table (existing logic)
	if tableX[i].Table_name == userTableName {
		_, _ = outputFile.WriteString(`	`+tableX[i].UserTableSpecs.AuthColumnName+`, err := util.HashPassword(util.RandomString(6))` + "\n")
		_, _ = outputFile.WriteString("	require.NoError(t, err)" + "\n")
	}
	
	// Build update parameters - ENHANCED WITH CONSTRAINT SUPPORT
	_, _ = outputFile.WriteString("	arg := Update" + tableX[i].FunctionSignature + "Params{" + "\n")
	for p := 0; p < len(tableX[i].Table_Columns); p++ {
		// Skip fields that are excluded from update (existing logic)
		if tableX[i].Table_Columns[p].UniqueFlag {
			continue
		}
		if (tableX[i].Table_Columns[p].ColumnType == "timestamptz" && tableX[i].Table_Columns[p].DefaultValue == "now()") {
			continue
		}
		if (tableX[i].Table_Columns[p].ColumnType == "date" && tableX[i].Table_Columns[p].DefaultValue == "CURRENT_DATE") {
			continue
		}
		
		// Handle different column types
		if tableX[i].Table_Columns[p].ForeignFlag {
			// Foreign key handling (existing logic)
			for k := 0; k < len(tableX[i].ForeignKeys); k++ {
				if tableX[i].ForeignKeys[k].FK_Column == tableX[i].Table_Columns[p].Column_name {
					// CHECK 4: Handle self-referencing foreign keys
					if tableX[i].Table_name == tableX[i].ForeignKeys[k].FK_Related_TableName {
						// For self-reference, set to NULL
						if strings.Contains(tableX[i].Table_Columns[p].ColumnType, "uuid") {
							_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.UUID{Valid: false},"+"\n")
						} else {
							_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.Text{Valid: false},"+"\n")
						}
					} else {
						// Normal foreign key handling 
						FormatedFieldName := FormatFieldName(tableX[i].ForeignKeys[k].FK_Related_Table_Column)
						key := tableX[i].Table_name+"_fk_"+tableX[i].ForeignKeys[k].FK_Related_SingularTableName
						if val, ok := fkVarMap[key]; ok {
							if tableX[i].Table_Columns[p].Not_Null {
								_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": " + val + "." + FormatedFieldName+","+"\n")
							} else {
								// For nullable foreign keys, wrap in pgtype
								if strings.Contains(tableX[i].Table_Columns[p].ColumnType, "uuid") {
									_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.UUID{Bytes: " + val + "." + FormatedFieldName + ", Valid: true},"+"\n")
								} else if strings.Contains(tableX[i].Table_Columns[p].ColumnType, "varchar") {
									_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.Text{String: " + val + "." + FormatedFieldName + ", Valid: true},"+"\n")
								} else if strings.Contains(tableX[i].Table_Columns[p].ColumnType, "bigint") {
									_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.Int8{Int64: " + val + "." + FormatedFieldName + ", Valid: true},"+"\n")
								}
							}
						}						
					}
				}
			}
		} else {
			// Primary key handling (existing logic)
			if tableX[i].Table_Columns[p].PrimaryFlag {
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": " + tableX[i].OutputFileName + "1." + getByColumnName + "," + "\n")
				continue
			}
			
			// ENHANCED: Column type specific handling with constraint support
			if tableX[i].Table_Columns[p].ColumnType == "varchar" {
				if tableX[i].Table_Columns[p].Not_Null {
					if tableX[i].Table_name == userTableName && tableX[i].Table_Columns[p].Column_name == tableX[i].UserTableSpecs.AuthColumnName{
						_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ":    "+tableX[i].UserTableSpecs.AuthColumnName+"," + "\n")
					} else {
						// ENHANCED: Use constraint-aware value generation
						validValue := generateValidValueForCheckConstraint(tableX[i].Table_name, tableX[i].Table_Columns[p].Column_name, tableX[i].Table_Columns[p].ColumnType, tableX[i].CheckConstraints)
						_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": " + validValue + "," + "\n")
					}
				} else {
					// ENHANCED: Use constraint-aware value generation for nullable varchar
					validValue := generateValidValueForCheckConstraint(tableX[i].Table_name, tableX[i].Table_Columns[p].Column_name, tableX[i].Table_Columns[p].ColumnType, tableX[i].CheckConstraints)
					if strings.Contains(validValue, "util.") {
						_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.Text{String: " + validValue + ", Valid: true}," + "\n")
					} else {
						_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.Text{String: " + validValue + ", Valid: true}," + "\n")
					}
				}
			} else if tableX[i].Table_Columns[p].ColumnType == "bigint" {
				if tableX[i].Table_Columns[p].Not_Null {
					// ENHANCED: Use constraint-aware value generation
					validValue := generateValidValueForCheckConstraint(tableX[i].Table_name, tableX[i].Table_Columns[p].Column_name, tableX[i].Table_Columns[p].ColumnType, tableX[i].CheckConstraints)
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": " + validValue + "," + "\n")
				} else {
					// ENHANCED: Use constraint-aware value generation for nullable bigint
					validValue := generateValidValueForCheckConstraint(tableX[i].Table_name, tableX[i].Table_Columns[p].Column_name, tableX[i].Table_Columns[p].ColumnType, tableX[i].CheckConstraints)
					
					constraintAllowsNull := checkIfConstraintAllowsNull(tableX[i].Table_Columns[p].Column_name, tableX[i].CheckConstraints)
					
					if constraintAllowsNull {
						if strings.Contains(validValue, "util.") {
							_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.Int8{Int64: " + validValue + ", Valid: true}," + "\n")
						} else {
							if numVal, err := strconv.ParseInt(validValue, 10, 64); err == nil {
								_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.Int8{Int64: " + strconv.FormatInt(numVal, 10) + ", Valid: true}," + "\n")
							} else {
								_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.Int8{Int64: " + validValue + ", Valid: true}," + "\n")
							}
						}
					} else {
						if strings.Contains(validValue, "util.") {
							_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.Int8{Int64: " + validValue + ", Valid: true}," + "\n")
						} else {
							_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.Int8{Int64: " + validValue + ", Valid: true}," + "\n")
						}
					}
				}
			} else if tableX[i].Table_Columns[p].ColumnType == "real" {
				if tableX[i].Table_Columns[p].Not_Null {
					// ENHANCED: Use constraint-aware value generation
					
					
					// Check if this column has an arithmetic constraint
					hasArithmeticConstraint := checkForArithmeticConstraint(tableX[i].Table_name, tableX[i].Table_Columns[p].Column_name, tableX[i].CheckConstraints)
					if hasArithmeticConstraint {
						// For arithmetic constraints, we'll set this after struct creation
						_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ":    0, // Will be calculated after struct creation" + "\n")
					} else {
						// Normal constraint handling
						validValue := generateValidValueForCheckConstraint(tableX[i].Table_name, tableX[i].Table_Columns[p].Column_name, tableX[i].Table_Columns[p].ColumnType, tableX[i].CheckConstraints)
						_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ":    " + validValue + "," + "\n")
					}
					
					
					
					
					
					// validValue := generateValidValueForCheckConstraint(tableX[i].Table_name, tableX[i].Table_Columns[p].Column_name, tableX[i].Table_Columns[p].ColumnType, tableX[i].CheckConstraints)
					// _, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": " + validValue + "," + "\n")
				} else {
					// ENHANCED: Use constraint-aware value generation for nullable real
					validValue := generateValidValueForCheckConstraint(tableX[i].Table_name, tableX[i].Table_Columns[p].Column_name, tableX[i].Table_Columns[p].ColumnType, tableX[i].CheckConstraints)
					
					constraintAllowsNull := checkIfConstraintAllowsNull(tableX[i].Table_Columns[p].Column_name, tableX[i].CheckConstraints)
					
					if constraintAllowsNull {
						if strings.Contains(validValue, "util.") {
							_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.Float4{Float32: " + validValue + ", Valid: true}," + "\n")
						} else {
							if numVal, err := strconv.ParseFloat(validValue, 32); err == nil {
								_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.Float4{Float32: " + strconv.FormatFloat(numVal, 'f', 2, 32) + ", Valid: true}," + "\n")
							} else {
								_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.Float4{Float32: " + validValue + ", Valid: true}," + "\n")
							}
						}
					} else {
						if strings.Contains(validValue, "util.") {
							_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.Float4{Float32: " + validValue + ", Valid: true}," + "\n")
						} else {
							_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.Float4{Float32: " + validValue + ", Valid: true}," + "\n")
						}
					}
				}
			} else {
				// For other column types, use existing logic (timestamptz, date, uuid, bool)
				if tableX[i].Table_Columns[p].ColumnType == "timestamptz" {
					if tableX[i].Table_Columns[p].Not_Null {
						_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": time.Now().UTC()," + "\n")
					} else {
						_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true}," + "\n")
					}
				} else if tableX[i].Table_Columns[p].ColumnType == "date" {
					if tableX[i].Table_Columns[p].Not_Null {
						_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": time.Date(2025, 5, 29, 0, 0, 0, 0, time.UTC)," + "\n")
					} else {
						_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.Date{Time: time.Date(2025, 5, 29, 0, 0, 0, 0, time.UTC), Valid: true}," + "\n")
					}
				} else if tableX[i].Table_Columns[p].ColumnType == "uuid" {
					if tableX[i].Table_Columns[p].Not_Null {
						_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": uuid.New()," + "\n")
					} else {
						_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.UUID{Bytes: uuid.New(), Valid: true}," + "\n")
					}
				} else if tableX[i].Table_Columns[p].ColumnType == "bool" {
					if tableX[i].Table_Columns[p].Not_Null {
						_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": true," + "\n")
					} else {
						_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.Bool{Bool: true, Valid: true}," + "\n")
					}
				}
			}
		}
	}
	_, _ = outputFile.WriteString("	}" + "\n")
	

	// After the struct creation (after the closing brace)
	// ✅ NEW: Handle arithmetic constraints after struct creation
	arithmeticConstraints := getArithmeticConstraints(tableX[i].Table_name, tableX[i].CheckConstraints)
	for _, constraint := range arithmeticConstraints {
		expression := generateArithmeticAssignment(constraint, tableX[i].Table_Columns)
		if expression != "" {
			_, _ = outputFile.WriteString("	" + expression + "\n")
		}
	}


	// Execute update (existing logic)
	_, _ = outputFile.WriteString("	" + tableX[i].OutputFileName + "2, err := testStore.Update" + tableX[i].FunctionSignature + "(context.Background(), " + "arg" + ")" + "\n")
	_, _ = outputFile.WriteString("	" + "require.NoError(t, err)" + "\n")
	_, _ = outputFile.WriteString("	" + "require.NotEmpty(t, " + tableX[i].OutputFileName + "2)" + "\n")
	_, _ = outputFile.WriteString("\n")
	
	// Validate updated values (existing logic)
	for h := 0; h < len(tableX[i].Table_Columns); h++ {
		// Skip excluded fields
		if tableX[i].Table_Columns[h].UniqueFlag {
			continue
		}
		if (tableX[i].Table_Columns[h].ColumnType == "timestamptz" && tableX[i].Table_Columns[h].DefaultValue == "now()") {
			continue
		}
		if (tableX[i].Table_Columns[h].ColumnType == "date" && tableX[i].Table_Columns[h].DefaultValue == "CURRENT_DATE") {
			continue
		}
		
		// Validation logic with proper pgtype handling
		if tableX[i].Table_Columns[h].PrimaryFlag || tableX[i].Table_Columns[h].ForeignFlag {
			_, _ = outputFile.WriteString("	require.Equal(t, " + tableX[i].OutputFileName + "1." + tableX[i].Table_Columns[h].ColumnNameParams + ", " + tableX[i].OutputFileName + "2." + tableX[i].Table_Columns[h].ColumnNameParams + ")" + "\n")
		} else if tableX[i].Table_Columns[h].ColumnType == "timestamptz" {
			if tableX[i].Table_Columns[h].Not_Null {
				// Non-nullable timestamptz: both arg and result are time.Time
				_, _ = outputFile.WriteString("	require.WithinDuration(t, arg." + tableX[i].Table_Columns[h].ColumnNameParams + ", " + tableX[i].OutputFileName + "2." + tableX[i].Table_Columns[h].ColumnNameParams + ", time.Second)" + "\n")
			} else {
				// Nullable timestamptz: both arg and result are pgtype.Timestamptz
				_, _ = outputFile.WriteString("	require.WithinDuration(t, arg." + tableX[i].Table_Columns[h].ColumnNameParams + ".Time, " + tableX[i].OutputFileName + "2." + tableX[i].Table_Columns[h].ColumnNameParams + ".Time, time.Second)" + "\n")
			}
		} else if tableX[i].Table_Columns[h].ColumnType == "date" {
			if tableX[i].Table_Columns[h].Not_Null {
				// Non-nullable date: both arg and result are time.Time
				_, _ = outputFile.WriteString("	require.Equal(t, arg." + tableX[i].Table_Columns[h].ColumnNameParams + ", " + tableX[i].OutputFileName + "2." + tableX[i].Table_Columns[h].ColumnNameParams + ")" + "\n")
			} else {
				// Nullable date: both arg and result are pgtype.Date
				_, _ = outputFile.WriteString("	require.Equal(t, arg." + tableX[i].Table_Columns[h].ColumnNameParams + ", " + tableX[i].OutputFileName + "2." + tableX[i].Table_Columns[h].ColumnNameParams + ")" + "\n")
			}
		} else {
			_, _ = outputFile.WriteString("	require.Equal(t, arg." + tableX[i].Table_Columns[h].ColumnNameParams + ", " + tableX[i].OutputFileName + "2." + tableX[i].Table_Columns[h].ColumnNameParams + ")" + "\n")
		}
	}
	_, _ = outputFile.WriteString("\n")
	_, _ = outputFile.WriteString("}" + "\n")
	_, _ = outputFile.WriteString("\n")
}

func printTestFuncForUpdate(tableX []dbschemareader.Table_Struct, i int, fk_HierarchyX []dbschemareader.FK_Hierarchy, outputFile *os.File, userTableName string) {
	var fkVarMap = make(map[string]string)
	_, _ = outputFile.WriteString("func TestUpdate" + tableX[i].FunctionSignature + "(t *testing.T) {" + "\n")
	
	// Create dependencies - skip self-referencing
	for k := 0; k < len(fk_HierarchyX); k++ {
		if fk_HierarchyX[k].TableName == tableX[i].Table_name {
			for l := len(fk_HierarchyX[k].RelatedTablesLevels) - 1; l >= 0; l-- {
				for m := 0; m < len(fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList); m++ {
					// CHECK 1: Skip self-referencing foreign keys in dependency creation
					if fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {
						continue
					}
					varName := fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName+"_fk_"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName+strconv.Itoa(k) + strconv.Itoa(l) + strconv.Itoa(m)
					key := fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName+"_fk_"+fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_SingularTableName
					// Only store if key doesn't exist
					if _, exists := fkVarMap[key]; !exists {
						fkVarMap[key] = varName
						_, _ = outputFile.WriteString("	" + varName+ " := createRandom" + fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName_Singular_Object + "(t")
					} else{
						continue
					}
					for g := 0; g < len(fk_HierarchyX); g++ {
						if fk_HierarchyX[g].TableName == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {
							for h := 0; h < len(fk_HierarchyX[g].RelatedTablesLevels); h++ {
								for z := 0; z < len(fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList); z++ {
									// CHECK 2: Skip self-referencing in nested dependencies
									if fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName == fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_TableName {
										continue
									}
									key := fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName
									if val, ok := fkVarMap[key]; ok {
										_, _ = outputFile.WriteString(", " + val)
									}		
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
	
	// Create record to update
	_, _ = outputFile.WriteString("	" + tableX[i].OutputFileName + "1 := createRandom" + tableX[i].FunctionSignature + "(t")
	for g := 0; g < len(fk_HierarchyX); g++ {
		if fk_HierarchyX[g].TableName == tableX[i].Table_name {
			if len(fk_HierarchyX[g].RelatedTablesLevels) > 0 {
				for h := 0; h < len(fk_HierarchyX[g].RelatedTablesLevels); h++ {
					for z := 0; z < len(fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList); z++ {
						// CHECK 3: Skip self-referencing in main function parameters
						if fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName == fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_TableName {
							continue
						}
						key := fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName
						if val, ok := fkVarMap[key]; ok {
							_, _ = outputFile.WriteString(", " + val)
						}		
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
	
	// Get primary key column name
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
	
	// Special handling for users table
	if tableX[i].Table_name == userTableName {
		_, _ = outputFile.WriteString(`	`+tableX[i].UserTableSpecs.AuthColumnName+`, err := util.HashPassword(util.RandomString(6))` + "\n")
		_, _ = outputFile.WriteString("	require.NoError(t, err)" + "\n")
	}
	
	// Build update parameters
	_, _ = outputFile.WriteString("	arg := Update" + tableX[i].FunctionSignature + "Params{" + "\n")
	for p := 0; p < len(tableX[i].Table_Columns); p++ {
		// Skip fields that are excluded from update (same logic as API handler)
		if tableX[i].Table_Columns[p].UniqueFlag {
			continue
		}
		if (tableX[i].Table_Columns[p].ColumnType == "timestamptz" && tableX[i].Table_Columns[p].DefaultValue == "now()") {
			continue
		}
		if (tableX[i].Table_Columns[p].ColumnType == "date" && tableX[i].Table_Columns[p].DefaultValue == "CURRENT_DATE") {
			continue
		}
		
		// Handle different column types
		if tableX[i].Table_Columns[p].ForeignFlag {
			for k := 0; k < len(tableX[i].ForeignKeys); k++ {
				if tableX[i].ForeignKeys[k].FK_Column == tableX[i].Table_Columns[p].Column_name {
					// CHECK 4: Handle self-referencing foreign keys
					if tableX[i].Table_name == tableX[i].ForeignKeys[k].FK_Related_TableName {
						// For self-reference, set to NULL
						if strings.Contains(tableX[i].Table_Columns[p].ColumnType, "uuid") {
							_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.UUID{Valid: false},"+"\n")
						} else {
							_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.Text{Valid: false},"+"\n")
						}
					} else {
						// Normal foreign key handling - use pgtype for nullable, regular for non-nullable
						FormatedFieldName := FormatFieldName(tableX[i].ForeignKeys[k].FK_Related_Table_Column)
						key := tableX[i].Table_name+"_fk_"+tableX[i].ForeignKeys[k].FK_Related_SingularTableName
						if val, ok := fkVarMap[key]; ok {
							if tableX[i].Table_Columns[p].Not_Null {
								_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": " + val + "." + FormatedFieldName+","+"\n")
							} else {
								// For nullable foreign keys, wrap in pgtype
								if strings.Contains(tableX[i].Table_Columns[p].ColumnType, "uuid") {
									_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.UUID{Bytes: " + val + "." + FormatedFieldName + ", Valid: true},"+"\n")
								} else if strings.Contains(tableX[i].Table_Columns[p].ColumnType, "varchar") {
									_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.Text{String: " + val + "." + FormatedFieldName + ", Valid: true},"+"\n")
								} else if strings.Contains(tableX[i].Table_Columns[p].ColumnType, "bigint") {
									_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.Int8{Int64: " + val + "." + FormatedFieldName + ", Valid: true},"+"\n")
								}
							}
						}						
					}
				}
			}
		} else {
			// Primary key handling
			if tableX[i].Table_Columns[p].PrimaryFlag {
				_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": " + tableX[i].OutputFileName + "1." + getByColumnName + "," + "\n")
				continue
			}			
			// Column type specific handling with nullable support
			if tableX[i].Table_Columns[p].ColumnType == "varchar" {
				if tableX[i].Table_Columns[p].Not_Null {
					if tableX[i].Table_name == userTableName && tableX[i].Table_Columns[p].Column_name == tableX[i].UserTableSpecs.AuthColumnName{
						_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ":    "+tableX[i].UserTableSpecs.AuthColumnName+"," + "\n")
					}else{
						_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ":    util.RandomName(8)," + "\n")
					}
				} else {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.Text{String: util.RandomName(8), Valid: true}," + "\n")
				}
			} else if tableX[i].Table_Columns[p].ColumnType == "bigint" {
				if tableX[i].Table_Columns[p].Not_Null {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": util.RandomInteger(1, 100)," + "\n")
				} else {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.Int8{Int64: util.RandomInteger(1, 100), Valid: true}," + "\n")
				}
			} else if tableX[i].Table_Columns[p].ColumnType == "real" {
				if tableX[i].Table_Columns[p].Not_Null {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": util.RandomReal(1, 100)," + "\n")
				} else {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.Float4{Float32: util.RandomReal(1, 100), Valid: true}," + "\n")
				}
			} else if tableX[i].Table_Columns[p].ColumnType == "timestamptz" {
				if tableX[i].Table_Columns[p].Not_Null {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": time.Now().UTC()," + "\n")
				} else {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true}," + "\n")
				}
			} else if tableX[i].Table_Columns[p].ColumnType == "date" {
				if tableX[i].Table_Columns[p].Not_Null {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": time.Now().UTC()," + "\n")
				} else {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.Date{Time: time.Now().UTC(), Valid: true}," + "\n")
				}
			} else if tableX[i].Table_Columns[p].ColumnType == "uuid" {
				if tableX[i].Table_Columns[p].Not_Null {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": uuid.New()," + "\n")
				} else {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.UUID{Bytes: uuid.New(), Valid: true}," + "\n")
				}
			} else if tableX[i].Table_Columns[p].ColumnType == "bool" {
				if tableX[i].Table_Columns[p].Not_Null {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": true," + "\n")
				} else {
					_, _ = outputFile.WriteString("		" + tableX[i].Table_Columns[p].ColumnNameParams + ": pgtype.Bool{Bool: true, Valid: true}," + "\n")
				}
			}
		}
	}
	_, _ = outputFile.WriteString("	}" + "\n")
	
	// Execute update
	_, _ = outputFile.WriteString("	" + tableX[i].OutputFileName + "2, err := testStore.Update" + tableX[i].FunctionSignature + "(context.Background(), " + "arg" + ")" + "\n")
	_, _ = outputFile.WriteString("	" + "require.NoError(t, err)" + "\n")
	_, _ = outputFile.WriteString("	" + "require.NotEmpty(t, " + tableX[i].OutputFileName + "2)" + "\n")
	_, _ = outputFile.WriteString("\n")
	
	// Validate updated values
	for h := 0; h < len(tableX[i].Table_Columns); h++ {
		// Skip excluded fields
		if tableX[i].Table_Columns[h].UniqueFlag {
			continue
		}
		if (tableX[i].Table_Columns[h].ColumnType == "timestamptz" && tableX[i].Table_Columns[h].DefaultValue == "now()") {
			continue
		}
		if (tableX[i].Table_Columns[h].ColumnType == "date" && tableX[i].Table_Columns[h].DefaultValue == "CURRENT_DATE") {
			continue
		}
		
		// Validation logic with proper pgtype handling
		if tableX[i].Table_Columns[h].PrimaryFlag || tableX[i].Table_Columns[h].ForeignFlag {
			_, _ = outputFile.WriteString("	require.Equal(t, " + tableX[i].OutputFileName + "1." + tableX[i].Table_Columns[h].ColumnNameParams + ", " + tableX[i].OutputFileName + "2." + tableX[i].Table_Columns[h].ColumnNameParams + ")" + "\n")
		} else if tableX[i].Table_Columns[h].ColumnType == "timestamptz" {
			if tableX[i].Table_Columns[h].Not_Null {
				// Non-nullable timestamptz: both arg and result are time.Time
				_, _ = outputFile.WriteString("	require.WithinDuration(t, arg." + tableX[i].Table_Columns[h].ColumnNameParams + ", " + tableX[i].OutputFileName + "2." + tableX[i].Table_Columns[h].ColumnNameParams + ", time.Second)" + "\n")
			} else {
				// Nullable timestamptz: both arg and result are pgtype.Timestamptz
				_, _ = outputFile.WriteString("	require.WithinDuration(t, arg." + tableX[i].Table_Columns[h].ColumnNameParams + ".Time, " + tableX[i].OutputFileName + "2." + tableX[i].Table_Columns[h].ColumnNameParams + ".Time, time.Second)" + "\n")
			}
		} else if tableX[i].Table_Columns[h].ColumnType == "date" {
			if tableX[i].Table_Columns[h].Not_Null {
				// Non-nullable date: both arg and result are time.Time
				_, _ = outputFile.WriteString("	require.WithinDuration(t, arg." + tableX[i].Table_Columns[h].ColumnNameParams + ", " + tableX[i].OutputFileName + "2." + tableX[i].Table_Columns[h].ColumnNameParams + ", time.Second)" + "\n")
			} else {
				// Nullable date: both arg and result are pgtype.Date
				_, _ = outputFile.WriteString("	require.WithinDuration(t, arg." + tableX[i].Table_Columns[h].ColumnNameParams + ".Time, " + tableX[i].OutputFileName + "2." + tableX[i].Table_Columns[h].ColumnNameParams + ".Time, time.Second)" + "\n")
			}
		} else {
			_, _ = outputFile.WriteString("	require.Equal(t, arg." + tableX[i].Table_Columns[h].ColumnNameParams + ", " + tableX[i].OutputFileName + "2." + tableX[i].Table_Columns[h].ColumnNameParams + ")" + "\n")
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
					// CHECK 1: Skip self-referencing foreign keys in dependency creation
					if fk_HierarchyX[k].RelatedTablesLevels[l].Hierarchy_TableName == fk_HierarchyX[k].RelatedTablesLevels[l].RelatedTableList[m].FK_Related_TableName {
						continue
					}
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
									// CHECK 2: Skip self-referencing in nested dependencies
									if fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName == fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_TableName {
										continue
									}
									key := fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName
									if val, ok := fkVarMap[key]; ok {
										_, _ = outputFile.WriteString(", " + val)
									}		
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
						// CHECK 3: Skip self-referencing in main function parameters
						if fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName == fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_TableName {
							continue
						}
						key := fk_HierarchyX[g].RelatedTablesLevels[h].Hierarchy_TableName+"_fk_"+fk_HierarchyX[g].RelatedTablesLevels[h].RelatedTableList[z].FK_Related_SingularTableName
						if val, ok := fkVarMap[key]; ok {
							_, _ = outputFile.WriteString(", " + val)
						}		
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

func TestWriter(projectFolderName string, gitHubAccountName string, dirPath string, tableX []dbschemareader.Table_Struct, fk_HierarchyX []dbschemareader.FK_Hierarchy) {
	//generating main_test.go
	fmt.Println("time.Now(): ",time.Now())
	main_testFunc(projectFolderName, gitHubAccountName, dirPath)
	//generate unit tests for db query files in sqlc folder
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
		var userTableName string //functionSignature, sessionStatusColumn, sessionTokenColumn, fkUserColumn, userTablePrimaryColumnName
		for x, table := range tableX {
			if table.IsSessionsTable {
				// functionSignature = table.FunctionSignature
				// sessionStatusColumn = table.SessionTableSpecs.StatusColumnName
				// sessionTokenColumn = table.SessionTableSpecs.TokenColumnName
				// fkUserColumn = table.SessionTableSpecs.FkUserColumn
			}
			if table.IsUserTable {
				userTableName = table.UserTableSpecs.TableName
				for _, col := range tableX[x].Table_Columns {
					if col.PrimaryFlag && len(col.DefaultValue) > 0 {
						// userTablePrimaryColumnName = col.ColumnNameParams
					}
				}
			}
		}	
		// CreateRandomFunction(tableX[:], i, outputFile, userTableName)
		CreateRandomFunctionWithConstraints(tableX[:], i, outputFile, userTableName)
		printTestFuncForCreate(tableX[:], i, fk_HierarchyX[:], outputFile)
		printTestFuncForReadGet(tableX[:], i, fk_HierarchyX[:], outputFile)
		printTestFuncForReadList(tableX[:], i, fk_HierarchyX[:], outputFile)
		// printTestFuncForUpdate(tableX[:], i, fk_HierarchyX[:], outputFile, userTableName)
		printTestFuncForUpdateWithConstraints(tableX[:], i, fk_HierarchyX[:], outputFile, userTableName)
		printTestFuncForDelete(tableX[:], i, fk_HierarchyX[:], outputFile)
		outputFile.Close()
	}
	fmt.Println("unit tests have been written successfully for db queries")
}
	// //writing jwt_maker_test.go
	// outputFileName = dirPath + "/token/jwt_maker_test.go"
	// outputFile, errs = os.Create(outputFileName)
	// if errs != nil {
	// 	fmt.Println("Failed to create file:", errs)
	// 	return
	// }
	// defer outputFile.Close()
	// _, _ = outputFile.WriteString("package token" + "\n")
	// _, _ = outputFile.WriteString("\n")
	// _, _ = outputFile.WriteString("import (" + "\n")
	// _, _ = outputFile.WriteString(`	"testing"` + "\n")
	// _, _ = outputFile.WriteString(`	"time"` + "\n")
	// _, _ = outputFile.WriteString("\n")
	// _, _ = outputFile.WriteString(`	"github.com/golang-jwt/jwt"` + "\n")
	// _, _ = outputFile.WriteString(`	"github.com/stretchr/testify/require"` + "\n")
	// _, _ = outputFile.WriteString(`	"github.com/`+gitHubAccountName+`/` + projectFolderName + `/util"` + "\n")
	// _, _ = outputFile.WriteString("\n")
	// _, _ = outputFile.WriteString(`)` + "\n")
	// _, _ = outputFile.WriteString("\n")

	// _, _ = outputFile.WriteString(`func TestJWTMaker(t *testing.T) {` + "\n")
	// _, _ = outputFile.WriteString(`	maker, err := NewJWTMaker(util.RandomString(32))` + "\n")
	// _, _ = outputFile.WriteString(`	require.NoError(t, err)` + "\n")
	// _, _ = outputFile.WriteString("\n")
	// _, _ = outputFile.WriteString(`	username := util.RandomName(8)` + "\n")
	// _, _ = outputFile.WriteString(`	role := util.UserLevel_1_Role` + "\n")
	// _, _ = outputFile.WriteString(`	duration := time.Minute` + "\n")
	// _, _ = outputFile.WriteString("\n")
	// _, _ = outputFile.WriteString(`	issuedAt := time.Now()` + "\n")
	// _, _ = outputFile.WriteString(`	expiredAt := issuedAt.Add(duration)` + "\n")
	// _, _ = outputFile.WriteString("\n")
	// _, _ = outputFile.WriteString(`	token, payload, err := maker.CreateToken(username, role, duration)` + "\n")
	// _, _ = outputFile.WriteString(`	//token, err := maker.CreateToken(username, role, duration)` + "\n")
	// _, _ = outputFile.WriteString(`	require.NoError(t, err)` + "\n")
	// _, _ = outputFile.WriteString(`	require.NotEmpty(t, token)` + "\n")
	// _, _ = outputFile.WriteString(`	require.NotEmpty(t, payload)` + "\n")
	// _, _ = outputFile.WriteString("\n")
	// _, _ = outputFile.WriteString(`	payload, err = maker.VerifyToken(token)` + "\n")
	// _, _ = outputFile.WriteString(`	require.NoError(t, err)` + "\n")
	// _, _ = outputFile.WriteString(`	require.NotEmpty(t, payload)` + "\n")
	// _, _ = outputFile.WriteString("\n")
	// _, _ = outputFile.WriteString(`	require.NotZero(t, payload.ID)` + "\n")
	// _, _ = outputFile.WriteString(`	require.Equal(t, username, payload.Username)` + "\n")
	// _, _ = outputFile.WriteString(`	//require.Equal(t, role, payload.Role)` + "\n")
	// _, _ = outputFile.WriteString(`	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)` + "\n")
	// _, _ = outputFile.WriteString(`	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)` + "\n")
	// _, _ = outputFile.WriteString(`}` + "\n")

	// _, _ = outputFile.WriteString(`func TestExpiredJWTToken(t *testing.T) {` + "\n")
	// _, _ = outputFile.WriteString(`	maker, err := NewJWTMaker(util.RandomString(32))` + "\n")
	// _, _ = outputFile.WriteString(`	require.NoError(t, err)` + "\n")
	// _, _ = outputFile.WriteString("\n")
	// _, _ = outputFile.WriteString(`	//token, payload, err := maker.CreateToken(util.RandomName(8), util.DepositorRole, -time.Minute)` + "\n")
	// _, _ = outputFile.WriteString(`	token, payload, err := maker.CreateToken(util.RandomName(8), util.UserLevel_1_Role, -time.Minute)` + "\n")
	// _, _ = outputFile.WriteString(`	require.NoError(t, err)` + "\n")
	// _, _ = outputFile.WriteString(`	require.NotEmpty(t, token)` + "\n")
	// _, _ = outputFile.WriteString(`	require.NotEmpty(t, payload)` + "\n")
	// _, _ = outputFile.WriteString("\n")
	// _, _ = outputFile.WriteString(`	payload, err = maker.VerifyToken(token)` + "\n")
	// _, _ = outputFile.WriteString(`	require.Error(t, err)` + "\n")
	// _, _ = outputFile.WriteString(`	require.EqualError(t, err, ErrExpiredToken.Error())` + "\n")
	// _, _ = outputFile.WriteString(`	require.Nil(t, payload)` + "\n")
	// _, _ = outputFile.WriteString(`}` + "\n")
	// _, _ = outputFile.WriteString("\n")

	// _, _ = outputFile.WriteString(`func TestInvalidJWTTokenAlgNone(t *testing.T) {` + "\n")
	// _, _ = outputFile.WriteString(`	//payload, err := NewPayload(util.RandomOwner(), util.DepositorRole, time.Minute)` + "\n")
	// _, _ = outputFile.WriteString(`	payload, err := NewPayload(util.RandomName(8), util.UserLevel_1_Role, time.Minute)` + "\n")
	// _, _ = outputFile.WriteString(`	require.NoError(t, err)` + "\n")
	// _, _ = outputFile.WriteString("\n")
	// _, _ = outputFile.WriteString(`	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)` + "\n")
	// _, _ = outputFile.WriteString(`	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)` + "\n")
	// _, _ = outputFile.WriteString(`	require.NoError(t, err)` + "\n")
	// _, _ = outputFile.WriteString("\n")
	// _, _ = outputFile.WriteString(`	maker, err := NewJWTMaker(util.RandomString(32))` + "\n")
	// _, _ = outputFile.WriteString(`	require.NoError(t, err)` + "\n")
	// _, _ = outputFile.WriteString("\n")
	// _, _ = outputFile.WriteString(`	payload, err = maker.VerifyToken(token)` + "\n")
	// _, _ = outputFile.WriteString(`	require.Error(t, err)` + "\n")
	// _, _ = outputFile.WriteString(`	require.EqualError(t, err, ErrInvalidToken.Error())` + "\n")
	// _, _ = outputFile.WriteString(`	require.Nil(t, payload)` + "\n")
	// _, _ = outputFile.WriteString(`}` + "\n")
	// _, _ = outputFile.WriteString("\n")
	// fmt.Println("jwt_maker_test.go file has been generated successfully")
	// outputFile.Close()

	// //writing paseto_maker_test.go
	// outputFileName = dirPath + "/token/paseto_maker_test.go"
	// outputFile, errs = os.Create(outputFileName)
	// if errs != nil {
	// 	fmt.Println("Failed to create file:", errs)
	// 	return
	// }
	// defer outputFile.Close()
	// _, _ = outputFile.WriteString("package token" + "\n")
	// _, _ = outputFile.WriteString("\n")
	// _, _ = outputFile.WriteString("import (" + "\n")
	// _, _ = outputFile.WriteString(`	"testing"` + "\n")
	// _, _ = outputFile.WriteString(`	"time"` + "\n")
	// _, _ = outputFile.WriteString("\n")
	// _, _ = outputFile.WriteString(`	"github.com/stretchr/testify/require"` + "\n")
	// _, _ = outputFile.WriteString(`	"github.com/`+gitHubAccountName+`/` + projectFolderName + `/util"` + "\n")
	// _, _ = outputFile.WriteString("\n")
	// _, _ = outputFile.WriteString(`)` + "\n")
	// _, _ = outputFile.WriteString("\n")

	// _, _ = outputFile.WriteString(`func TestPasetoMaker(t *testing.T) {` + "\n")
	// _, _ = outputFile.WriteString(`	maker, err := NewPasetoMaker(util.RandomString(32))` + "\n")
	// _, _ = outputFile.WriteString(`	require.NoError(t, err)` + "\n")
	// _, _ = outputFile.WriteString("\n")
	// _, _ = outputFile.WriteString(`	username := util.RandomName(8)` + "\n")
	// _, _ = outputFile.WriteString(`	role := util.UserLevel_1_Role` + "\n")
	// _, _ = outputFile.WriteString(`	duration := time.Minute` + "\n")
	// _, _ = outputFile.WriteString("\n")
	// _, _ = outputFile.WriteString(`	issuedAt := time.Now()` + "\n")
	// _, _ = outputFile.WriteString(`	expiredAt := issuedAt.Add(duration)` + "\n")
	// _, _ = outputFile.WriteString("\n")
	// _, _ = outputFile.WriteString(`	token, payload, err := maker.CreateToken(username, role, duration)` + "\n")
	// _, _ = outputFile.WriteString(`	//token, err := maker.CreateToken(username, role, duration)` + "\n")
	// _, _ = outputFile.WriteString(`	require.NoError(t, err)` + "\n")
	// _, _ = outputFile.WriteString(`	require.NotEmpty(t, token)` + "\n")
	// _, _ = outputFile.WriteString(`	require.NotEmpty(t, payload)` + "\n")
	// _, _ = outputFile.WriteString("\n")
	// _, _ = outputFile.WriteString(`	payload, err = maker.VerifyToken(token)` + "\n")
	// _, _ = outputFile.WriteString(`	require.NoError(t, err)` + "\n")
	// _, _ = outputFile.WriteString(`	require.NotEmpty(t, payload)` + "\n")
	// _, _ = outputFile.WriteString("\n")
	// _, _ = outputFile.WriteString(`	require.NotZero(t, payload.ID)` + "\n")
	// _, _ = outputFile.WriteString(`	require.Equal(t, username, payload.Username)` + "\n")
	// _, _ = outputFile.WriteString(`	//require.Equal(t, role, payload.Role)` + "\n")
	// _, _ = outputFile.WriteString(`	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)` + "\n")
	// _, _ = outputFile.WriteString(`	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)` + "\n")
	// _, _ = outputFile.WriteString(`}` + "\n")

	// _, _ = outputFile.WriteString(`func TestExpiredPasetoToken(t *testing.T) {` + "\n")
	// _, _ = outputFile.WriteString(`	maker, err := NewPasetoMaker(util.RandomString(32))` + "\n")
	// _, _ = outputFile.WriteString(`	require.NoError(t, err)` + "\n")
	// _, _ = outputFile.WriteString("\n")
	// _, _ = outputFile.WriteString(`	//token, payload, err := maker.CreateToken(util.RandomName(8), util.DepositorRole, -time.Minute)` + "\n")
	// _, _ = outputFile.WriteString(`	token, payload, err := maker.CreateToken(util.RandomName(8), util.UserLevel_1_Role, -time.Minute)` + "\n")
	// _, _ = outputFile.WriteString(`	require.NoError(t, err)` + "\n")
	// _, _ = outputFile.WriteString(`	require.NotEmpty(t, token)` + "\n")
	// _, _ = outputFile.WriteString(`	require.NotEmpty(t, payload)` + "\n")
	// _, _ = outputFile.WriteString("\n")
	// _, _ = outputFile.WriteString(`	payload, err = maker.VerifyToken(token)` + "\n")
	// _, _ = outputFile.WriteString(`	require.Error(t, err)` + "\n")
	// _, _ = outputFile.WriteString(`	require.EqualError(t, err, ErrExpiredToken.Error())` + "\n")
	// _, _ = outputFile.WriteString(`	require.Nil(t, payload)` + "\n")
	// _, _ = outputFile.WriteString(`}` + "\n")
	// _, _ = outputFile.WriteString("\n")
	// fmt.Println("paseto_maker_test.go file has been generated successfully")
	// outputFile.Close()

	// //Executing goimports
	// cmd := exec.Command("goimports", "-w", ".")
	// cmd.Dir = dirPath + "/db/sqlc"
	// cmd.Run()
	// println("goimports executed successfully")
	// // var stderr bytes.Buffer

	// //Executing go mod tidy
	// cmd = exec.Command("go", "mod", "tidy")
	// cmd.Dir = dirPath
	// cmd.Run()
	// println("go mod tidy executed successfully")

