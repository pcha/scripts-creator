package main

import "os"

const genDirVarName = "gendir"

type Config struct {
	GenDir string
}

func buildConfig() Config {
	return Config{
		GenDir: os.Getenv(genDirVarName),
	}
}
