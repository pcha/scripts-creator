package scripts

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sortTasks(t *testing.T) {
	rm := Task{
		Name: "rm",
		Dependencies: []string{
			"cat",
		},
	}
	cat := Task{
		Name: "cat",
		Dependencies: []string{
			"chown",
			"chmod",
		},
	}
	touch := Task{
		Name: "touch",
	}
	chmod := Task{
		Name: "chmod",
		Dependencies: []string{
			"touch",
		},
	}
	chown := Task{
		Name: "chown",
		Dependencies: []string{
			"touch",
		},
	}
	tasks := []Task{
		rm,
		cat,
		touch,
		chmod,
		chown,
	}

	sorter := new(TopologicalSorter)
	tasks = sorter.Sort(tasks)
	assert.Len(t, tasks, 5)
	executed := make(map[string]bool, len(tasks))
	for _, task := range tasks {
		for _, d := range task.Dependencies {
			_, f := executed[d]
			assert.True(t, f, "The task %q would be executed before dependency %q", task.Name, d)
		}
		executed[task.Name] = true
	}
	assert.Len(t, executed, 5)

	fmt.Println(executed)
}
