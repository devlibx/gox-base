Gox-Base project provide utilities which is used commonly in all applications.
1. Serialization utils
2. Json file to object
3. Yaml file to object
4. XML file to object
5. ... 

#Utility


#### Convert anything to string 
This utility will convert int, bool, interface{} to string. Object will output a json string.
```golang
out, _ := Stringify(10) 
// Output = "10"

boolOut, _ := Stringify(true) 
// Output = "true"


type utilTestStruct struct {
	IntValue    int    `json:"int"`
	BoolValue   bool   `json:"bool"`
	StringValue string `json:"string"`
}

objectOut, _ := Stringify(utilTestStruct{
		IntValue:    10,
		BoolValue:   false,
		StringValue: "some value",
	})
// Output = {"int":10,"bool":false,"string":"some value"}
```
###### Stringify with error suppressed
If you don't want to handle error and have default value on error then you
can use "suppress error" version of this method
```golang
intOut1 := StringifySuppressError(10, "0")
// Output = "10"

If there is a error when input is bad then you will get the default 
value "0"
```