package vm

import "fmt"
import "github.com/robertkrimen/otto"

type VM struct {
	Owner string
	Name  string
	Url   string
	Js    *otto.Otto
	Ruby  *RubyVM
}

func Factory(owner string, lang string) *VM {
	new_vm := VM{Owner: owner}
	if lang == "ruby" {
	  new_vm.Ruby = rubyfactory()
	}
	if lang == "javascript" {
		new_vm.Js = otto.New()
	}
	return &new_vm
}

func (vm *VM) Eval(code string) error {
	var err error
	if vm.Js != nil {
		err = vm.EvalJs(code)
	}
	if vm.Ruby != nil {
		err = vm.EvalRuby(code)
	}
	return err
}

func (vm *VM) EvalJs(js_code string) error {
	fmt.Println(string(js_code)[0:15])
	fmt.Println("--eval begins--")

	src, err := vm.Js.Compile("", js_code)

	if err != nil {
		fmt.Println("compile failed!", err)
		return err
	} else {
		fmt.Println("compile good!")
		setup, err := vm.Js.Run(src)
		if err != nil {
			fmt.Println("eval failed", err, vm.Js.Context().Stacktrace)
			return err
		} else {
			descriptor_value, err := setup.Call(setup)
			if err != nil {
				fmt.Println("js func setup eval fail")
				return err
			} else {
				descriptor_map, _ := descriptor_value.Export()
				descriptor := descriptor_map.(map[string]interface{})
				vm.Name = descriptor["name"].(string)
				return nil
			}
		}
	}
}

