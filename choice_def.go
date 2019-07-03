package statemachine

import (
	"reflect"

	"github.com/Gurpartap/statemachine-go/internal/dynafunc"
)

type ChoiceConditionDef struct {
	RegisteredFunc string          `json:",omitempty" hcl:"registered_func" hcle:"omitempty"`
	Condition      ChoiceCondition `json:"-" hcle:"omit"`
}

type ChoiceDef struct {
	Condition   *ChoiceConditionDef `json:",omitempty" hcl:"condition" hcle:"omitempty"`
	UnlessGuard *TransitionGuardDef `json:",omitempty" hcl:"unless_condition" hcle:"omitempty"`
	OnTrue      *EventDef           `json:",omitempty" hcl:"on_true" hcle:"omitempty"`
	OnFalse     *EventDef           `json:",omitempty" hcl:"on_false" hcle:"omitempty"`
}

func (def *ChoiceDef) SetCondition(condition ChoiceCondition) {
	def.Condition = &ChoiceConditionDef{Condition: condition}
}

func (def *ChoiceDef) SetUnlessGuard(guard TransitionGuard) {
	def.UnlessGuard = &TransitionGuardDef{Guard: guard}
}

func (def *ChoiceDef) SetOnTrue(eventBuilderFn func(eventBuilder EventBuilder)) {
	eventBuilder := NewEventBuilder("")
	eventBuilderFn(eventBuilder)
	e := newEventImpl()
	eventBuilder.Build(e)

	def.OnTrue = e.def
}

func (def *ChoiceDef) SetOnFalse(eventBuilderFn func(eventBuilder EventBuilder)) {
	eventBuilder := NewEventBuilder("")
	eventBuilderFn(eventBuilder)
	e := newEventImpl()
	eventBuilder.Build(e)

	def.OnFalse = e.def
}

func execChoice(condition ChoiceCondition, args map[reflect.Type]interface{}) bool {
	switch reflect.TypeOf(condition).Kind() {
	case reflect.Func:
		fn := dynafunc.NewDynamicFunc(condition, args)
		if err := fn.Call(); err != nil {
			panic(err)
		}
		// condition func must return a bool
		return fn.Out[0].Bool()
	case reflect.Ptr:
		if reflect.ValueOf(condition).Elem().Kind() == reflect.Bool {
			return reflect.ValueOf(condition).Elem().Bool()
		}
		fallthrough
	default:
		panic("choice must either be a compatible func or pointer to a bool variable")
	}
}
