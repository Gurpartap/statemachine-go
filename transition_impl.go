package statemachine

type transitionImpl struct {
	from string
	to   string
}

func newTransitionImpl(from, to string) *transitionImpl {
	return &transitionImpl{
		from: from,
		to:   to,
	}
}

func (t *transitionImpl) From() string {
	return t.from
}

func (t *transitionImpl) To() string {
	return t.to
}

// // callback1(next: {
// //   callback2(next: {
// //     callback3(next: {
// //       applyTransition()
// //     })
// //   })
// // })
// func (t *transitionImpl) applyTransitionAroundCallbacks(callbacks []TransitionCallbackFunc, applyTransition func()) {
// 	if len(callbacks) == 0 {
// 		applyTransition()
// 		return
// 	}
//
// 	calledBackNext := false
// 	t.exec(callbacks[0], func() {
// 		calledBackNext = true
// 		t.applyTransitionAroundCallbacks(callbacks[1:], applyTransition)
// 	})
// 	if !calledBackNext && len(callbacks) != 1 {
// 		fmt.Printf("len(callbacks): %d\n", len(callbacks))
// 		panic("non-last around callbacks must call next()")
// 	}
//
// 	return
// }
//
// func (t *transitionImpl) exec(callback TransitionCallbackFunc, aroundExecFunc func()) {
// 	args := make(map[reflect.Type]interface{})
// 	args[reflect.TypeOf(new(Transition))] = t
// 	if aroundExecFunc != nil {
// 		// we use ptr reference because of reflect.PtrTo(t) in dynamicFunc{}.Call()
// 		args[reflect.TypeOf(new(func()))] = aroundExecFunc
// 	}
// 	fn := dynafunc.NewDynamicFunc(callback, args)
// 	if err := fn.Call(); err != nil {
// 		panic(err)
// 	}
// }
