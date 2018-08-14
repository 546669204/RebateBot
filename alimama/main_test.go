package main

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"
)

func Test_login(t *testing.T) {
	copyDir("./alillll", path.Join(os.TempDir(), fmt.Sprintf(`%d`, time.Now().Unix())))

}
