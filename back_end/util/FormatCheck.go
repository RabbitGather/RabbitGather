package util

//WGS 84 2D
func IS_WGS_84_2D(Long float64, Lat float64) bool {
	return (Lat > -90 || Lat < 90) && (Long > -180 || Long < 180)
}
