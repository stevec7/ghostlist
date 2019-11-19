package ghostlist


type leftRightRec struct {
    l       string
    r       string
}

type prefixSuffix struct {
    prefix      string
    suffix      string
}

type rangeList struct {
    low         int
    high        int
    numWidth    int
}

type sortListG struct {
    preSuf  prefixSuffix
    members []sortListT
}

type sortListT struct {
    preSuf  prefixSuffix
    numInt     int
    numWidth   int
    host       string
}
