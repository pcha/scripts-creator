package scripts

type Sorter interface {
	Sort(ts []Task) []Task
}

type TopologicalSorter struct {
	pending     map[string]Task
	sortedStack []Task
}

func (s TopologicalSorter) Sort(ts []Task) []Task {
	var ks []string
	s.pending, ks = buildMap(ts)
	s.sortedStack = s.addToStack(ks)
	return s.sortedStack
}

func buildMap(ts []Task) (map[string]Task, []string) {
	m := make(map[string]Task, len(ts))
	ks := make([]string, len(ts))
	for i, t := range ts {
		m[t.Name] = t
		ks[i] = t.Name
	}
	return m, ks
}

func (s TopologicalSorter) addToStack(list []string) []Task {
	for _, n := range list {
		t, f := getKey(n, s.pending)
		if f {
			s.sortedStack = s.addToStack(t.Dependencies)
			s.sortedStack = append(s.sortedStack, t)
		}
	}
	return s.sortedStack
}

func getKey(k string, m map[string]Task) (Task, bool) {
	t, f := m[k]
	delete(m, k)
	return t, f
}
