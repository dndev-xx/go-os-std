package cmd

type Word struct {
	val string
	cnt int32
}

func NewWord(val string, cnt int32) Word {
	return Word{
		val: val,
		cnt: cnt,
	}
}
