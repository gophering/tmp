package tmp

import (
	"strings"
	"html"
	"fmt"
	"strconv"
	"sort"
	"encoding/json"
)

/* Весь код преобразующий интерфейсы взят из расширения CossackPyra/pyraconv */

var localisations = make(map[string][]string)
var templates = make(map[string]string)


func SetTemplate(name, text string){
	templates[name] = text
}

func SetLocalisation(name string, list []string){
 localisations[name] = list
}

func toString(i1 interface{}) string {
	if i1 == nil {
		return ""
	}
	switch i2 := i1.(type) {
	default:
		return fmt.Sprint(i2)
	case bool:
		if i2 {
			return "true"
		} else {
			return "false"
		}
	case string:
		return i2
	case *bool:
		if i2 == nil {
			return ""
		}
		if *i2 {
			return "true"
		} else {
			return "false"
		}
	case *string:
		if i2 == nil {
			return ""
		}
		return *i2
	case *json.Number:
		return i2.String()
	case json.Number:
		return i2.String()
	}
	return ""
}

func toStringMap(i1 interface{}) map[string]string {
	switch i2 := i1.(type) {
	case map[string]interface{}:
		m1 := map[string]string{}
		for k, v := range i2 {
			m1[k] = toString(v)
		}
		return m1
	case map[string]string:
		return i2
	default:
		return map[string]string{}
	}
}

func toIntMapStringMap(i1 interface{}) map[int]map[string]string {
	switch i2 := i1.(type) {
	case map[int]interface{}:
		m1 := map[int]map[string]string{}
		for k, v := range i2 {
			m1[k] = toStringMap(v)
		}
		return m1
	case map[int]map[string]string:
		return i2
	default:
		return map[int]map[string]string{}
	}
}


func searchpos(body string, tag string, key string) map[int]string{
	exit := make(map[int]string, 3)
	p1 := strings.Index(body, "{{#"+tag)
	p0 := p1
	p2 := 0
	lenbody := len(body)
	lenkey := len(key)
	lentag := len(tag)
	p3 := 0
	p4 := 0
	for ; p0+p2<p1+1; {
		p3 = strings.Index(body[p0+3:lenbody], "{{#"+tag)
		p4 = strings.Index(body[p0+3:lenbody], "{{/"+tag)
	if p4==-1{
		fmt.Println("Ошибка шаблона")
	}else{
		if p3==-1{
			p2 = p4
		}else{
			if p4<p3{
				p2 = p4
			}else{
				p0 = p4
			}
		}
	}
	}
	exit[0] = body[0:p1]
	exit[1] = body[p1+4+lenkey:p1+p2+3]
	exit[2] = body[p1+p2+8+lentag:lenbody]
	return exit
}

func Typ(inter interface{}) int{
	typ := 0
	switch inter.(type){
		case string: 
			typ = 1
		break
		case map[int]string:
			typ = 2
		case map[string]string:
			typ = 3
		break
	}
	return typ
}

func Render(body string, arr map[string]interface{}, languge string) string{
complete := true
for ;complete; {
p1 := strings.Index(body, "{{")
if p1>-1 {
	p2:= strings.Index(body, "}}")
	key := body[p1+2:p2]
	tag := key[0:1]
	key = strings.Replace(key, "{", "", -1)
	key0 := strings.Replace(key, "#not", "", -1)
	key0 = strings.Replace(key0, "#if", "", -1)
	key0 = strings.Replace(key0, "#file", "", -1)
	key0 = strings.Replace(key0, "#array", "", -1)
	key0 = strings.Replace(key0, "%", "", -1)
	key0 = strings.Replace(key0, " ", "", -1)
	keys := strings.Split(key0, ".")
	//fmt.Println(key)
	//fmt.Printf("%+v\n",arr["ppp"])
	switch arr[keys[0]].(type){
//////////////////////////////////////////// nil
		case nil:
			if tag=="{"{
				body = strings.Replace(body, "{{{"+key+"}}}", "", -1)
			}else if tag=="%"{
				i, _ := strconv.ParseInt(keys[0],0,8)
				body = strings.Replace(body, "{{"+key+"}}", localisations[languge][i-1], -1)
			}else if tag == "#"{
				if strings.Index(key, "#file ")>-1{
					body = strings.Replace(body, "{{"+key+"}}", templates[keys[0]], -1)
				}else if strings.Index(key, "#if ")>-1{
					parts := searchpos(body, "if", key)
							body = parts[0]+parts[2]
				}else if strings.Index(key, "#not ")>-1{
					parts := searchpos(body, "not", key)
							body = parts[0]+parts[1]+parts[2]
				}
			}else{
					body = strings.Replace(body, "{{"+key+"}}", "", -1)
			}
/////////////////////////////////////////// string
		break
		case string:
			a := toString(arr[keys[0]])
			if tag=="{"{
					body = strings.Replace(body, "{{{"+key+"}}}",a, -1)
			}else if tag == "#"{
				if strings.Index(key, "#file ")>-1{
					body = strings.Replace(body, "{{"+key+"}}", templates[a], -1)
				}else if strings.Index(key, "#if ")>-1{
					parts := searchpos(body, "if", key)
						if a!="" {
							body = parts[0]+parts[1]+parts[2]
						}else{
							body = parts[0]+parts[2]
						}
				}else if strings.Index(key, "#not ")>-1{
					parts := searchpos(body, "not", key)
						if a!="" {
							body = parts[0]+parts[2]
						}else{
							body = parts[0]+parts[1]+parts[2]
						}
				}
			}else{
					body = strings.Replace(body, "{{"+key+"}}", html.EscapeString(a), -1)				
			}
		break
//////////////////////////////// map[string]string
		case map[string]string:
			a := toStringMap(arr[keys[0]])
			if tag=="{"{
				body = strings.Replace(body, "{{{"+key+"}}}", a[keys[1]], -1)
			}else if tag=="%"{
				i, _ := strconv.ParseInt(keys[0],0,8)
				body = strings.Replace(body, "{{"+key+"}}", localisations[languge][i-1], -1)
			}else if tag == "#"{
				if strings.Index(key, "#file ")>-1{
					body = strings.Replace(body, "{{"+key+"}}", templates[keys[0]], -1)
				}else if strings.Index(key, "#if ")>-1{
					parts := searchpos(body, "if", key)
						if _, ok:=a[keys[1]]; ok && a[keys[1]]!="" {
							body = parts[0]+parts[1]+parts[2]
						}else{
							body = parts[0]+parts[2]
						}
				}else if strings.Index(key, "#not ")>-1{
					parts := searchpos(body, "not", key)
						if _, ok:=a[keys[1]]; ok && a[keys[1]]!="" {
							body = parts[0]+parts[2]
						}else{
							body = parts[0]+parts[1]+parts[2]
						}
				}
			}else{
					body = strings.Replace(body, "{{"+key+"}}", html.EscapeString(a[keys[1]]), -1)
			}
/////////////////////////////////// map[int]map[string]string
		case map[int]map[string]string:
			if tag=="{"{
					body = strings.Replace(body, "{{{"+key+"}}}", "", -1)
			}else if tag=="%"{
				i, _ := strconv.ParseInt(keys[0],0,8)
				body = strings.Replace(body, "{{"+key+"}}", localisations[languge][i-1], -1)
			}else if tag == "#"{
				if strings.Index(key, "#file ")>-1{
					body = strings.Replace(body, "{{"+key+"}}", templates[keys[0]], -1)
				}else if strings.Index(key, "#if ")>-1{
					parts := searchpos(body, "if", key)
						if arr[keys[0]]==nil{
							body = parts[0]+parts[1]+parts[2]
						}else{
							body = parts[0]+parts[2]
						}
				}else if strings.Index(key, "#not ")>-1{
					parts := searchpos(body, "not", key)
						if arr[keys[0]]==nil{
							body = parts[0]+parts[2]
						}else{
							body = parts[0]+parts[1]+parts[2]
						}
				}else if strings.Index(key, "#array ")>-1{
					parts := searchpos(body, "array", key)
					arr2 := arr
					body2 := ""
					arr1 := toIntMapStringMap(arr[keys[0]])
					 var k []int
					for x, _ := range arr1 {
						k = append(k, x)
					}
					sort.Ints(k)
					for _, x := range k {
						arr2[keys[0]] = arr1[x]

						body2 = body2+Render(parts[1], arr2, languge)
					}
					body = parts[0]+body2+parts[2]
				}
			}else{
				body = strings.Replace(body, "{{"+key+"}}", "", -1)
			}
		break
/////////////////////////////////// default
		default:
			if tag=="{"{
				body = strings.Replace(body, "{{{"+key+"}}}", "", -1)
			}else if tag=="%"{
				i, _ := strconv.ParseInt(keys[0],0,8)
				body = strings.Replace(body, "{{"+key+"}}", localisations[languge][i-1], -1)
			}else if tag == "#"{
				if strings.Index(key, "#file ")>-1{
					body = strings.Replace(body, "{{"+key+"}}", templates[keys[0]], -1)
				}else if strings.Index(key, "#if ")>-1{
					parts := searchpos(body, "if", key)
						if arr[keys[0]]==nil{
							body = parts[0]+parts[1]+parts[2]
						}else{
							body = parts[0]+parts[2]
						}
				}else if strings.Index(key, "#not ")>-1{
					parts := searchpos(body, "not", key)
						if arr[keys[0]]==nil{
							body = parts[0]+parts[2]
						}else{
							body = parts[0]+parts[1]+parts[2]
						}
				}
			}else{
					body = strings.Replace(body, "{{"+key+"}}", "", -1)

			}
	}
	//fmt.Println(key)
}else{
complete = false
}
}
return body
}
