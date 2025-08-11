package greeter

import (
	cryptornd "crypto/rand"
	"fmt"
	"math/rand"
)

type GrtRepo interface{
    Predefined() string
    Random() string
    Answer(a string) string
}

type GrtRepoFake struct { }

func NewGrtRepo() *GrtRepoFake {
    return &GrtRepoFake{}
}

func (r *GrtRepoFake) Predefined() string {
    return "Hello world from Greeter Repo"
}

func (r *GrtRepoFake) Random() string {
    return cryptornd.Text()
}

func (r *GrtRepoFake) Answer(a string) string {
    return fmt.Sprintf("with ans: %s and random %v", a,rand.Int())
}
