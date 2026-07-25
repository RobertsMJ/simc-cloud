package simc

// Some expected strings:
// waist=,id=244573,bonus_id=12214/13667/12497/12066/8960/12384/8791/13622/12667,content_tuning=3615,crafted_stats=49/32,crafting_quality=5
// server=illidan
// omnium_talents=136819:1/136822:1

type ValueMarshaler interface{ MarshalSimcValue() (string, error) }
type ValueUnmarshaler interface{ UnmarshalSimcValue(string) error }
