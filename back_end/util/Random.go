package util

import (
	"github.com/bwmarrin/snowflake"
	"math"
	"math/rand"
	"time"
)

var node *snowflake.Node
var randomInst *rand.Rand

func init() {
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		panic(err.Error())
	}
}

func init() {
	randomInst = rand.New(rand.NewSource(Snowflake().Int64()))
}
func Snowflake() snowflake.ID {
	return node.Generate()
}

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

func GetSnowflakeIntWithLength(lg int64) int64 {
	return CutIntMax(Snowflake().Int64(), lg)
}

func GetRandomInt(min int, max int) int {
	if min < 0 {
		panic("min < 0")
	}
	return rand.Intn(max-min) + min //CutIntBetweenPint(Snowflake().Int64(), int64(math.Log10(float64(min)))+1, int64(math.Log10(float64(max)))+1)
}

func RandomInLength(i int) int {
	return rand.Intn(int(math.Floor(math.Pow(10.0, float64(i)))))
}

//func NewVerificationCodeWithLength(d int) string {
//	return fmt.Sprintf(fmt.Sprintf("%%0%dd", d), RandomInLength(d))
//}

var Random *rand.Rand

func init() {
	Random = rand.New(rand.NewSource(time.Now().Unix()))
}

func randFloats(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}
