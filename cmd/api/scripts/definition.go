package scripts

type Task struct {
	Name         string   `json:"name"`
	Command      string   `json:"command"`
	Dependencies []string `json:"dependencies"`
}

type TasksList []Task

type Definition struct {
	Filename string    `json:"filename"`
	Tasks    TasksList `json:"tasks"`
}
