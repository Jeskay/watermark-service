package internal

type Position string

const (
	LeftTop     Position = "left_top"
	LeftBottom  Position = "left_bottom"
	RightTop    Position = "right_top"
	RightBottom Position = "right_bottom"
)

func PositionFromString(text string) Position {
	switch text {
	case "left_top":
		return LeftTop
	case "left_bottom":
		return LeftBottom
	case "right_top":
		return RightTop
	case "right_bottom":
		return RightBottom
	}
	return LeftTop
}
