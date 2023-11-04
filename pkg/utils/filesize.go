package utils

var ToBytesSize = map[string]int{
	"bytes": 1,
	"kb":    1024,
	"mb":    1024 * 1024,
	"gb":    1024 * 1024 * 1024,
	"tb":    1024 * 1024 * 1024 * 1024,
}

// var FromBytesSize = map[string]int{
// 	"kb": 1024,
// 	"mb": 1024 * 1024,
// 	"gb": 1024 * 1024 * 1024,
// 	"tb": 1024 * 1024 * 1024 * 1024,
// }

func ConvertFileSize(size float64, currentUnit, resultUnit string) float64 {
	toByteMultiplier := ToBytesSize[currentUnit]
	result := size * float64(toByteMultiplier)
	resultUnitFactor := ToBytesSize[resultUnit]
	result = result / float64(resultUnitFactor)
	return result
}
