package utils

func CreateIDGenerator() func() int {
    var f func() int
    var counter int
    f = func() int {
        counter++
        return counter
    }
    return f
}
