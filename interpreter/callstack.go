package interpreter

import "container/list"

type CallStack struct {
	list *list.List
}

func NewCallStack() *CallStack {
	list := list.New()
	return &CallStack{
		list: list,
	}
}

func (s *CallStack) Push(x *CallStackFrame) {
	if x == nil {
		panic("Pushed nill call stack frame")
	}
	s.list.PushBack(x)
}

func (s *CallStack) Pop() *CallStackFrame {
	back := s.list.Back()
	if back == nil {
		return nil
	}
	s.list.Remove(back)
	return back.Value.(*CallStackFrame)
}

func (s *CallStack) Peek() *CallStackFrame {
	back := s.list.Back()
	if back == nil {
		return nil
	}
	return back.Value.(*CallStackFrame)
}

func (s *CallStack) ResolveVariable(variableName string) (*BeccaValue, bool) {
	back := s.list.Back()
	if back == nil {
		return nil, false
	}
	stackFrame := back.Value.(*CallStackFrame)
	value, found := stackFrame.Scope[variableName]
	if found {
		return value, true
	}
	for prev := back.Prev(); prev != nil; {
		stackFrame := back.Value.(*CallStackFrame)
		value, found := stackFrame.Scope[variableName]
		if found {
			return value, true
		}
		back = prev
	}
	return nil, false
}

func (s *CallStack) AssignVariable(variableName string, value *BeccaValue) error {
	s.Peek().Scope[variableName] = value
	return nil
}
