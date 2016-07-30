package mat

// Normalized an array of floats into these params if they exist
// if they are not present then thier default values are returned
// x The position of the object along the x-axis.
// y The position of the object along the y-axis.
// angle The angle of the object (in radians).
// sx The scale factor along the x-axis.
// sy The scale factor along the y-axis.
// ox The origin offset along the x-axis.
// oy The origin offset along the y-axis.
// kx Shear along the x-axis.
// ky Shear along the y-axis.
func normalizeDrawCallArgs(args []float32) (float32, float32, float32, float32, float32, float32, float32, float32, float32) {
	var x, y, angle, sx, sy, ox, oy, kx, ky float32
	sx = 1
	sy = 1

	if args == nil || len(args) < 2 {
		return x, y, angle, sx, sy, ox, oy, kx, ky
	}

	args_length := len(args)

	switch args_length {
	case 9:
		ky = args[8]
		fallthrough
	case 8:
		kx = args[7]
		if args_length == 8 {
			ky = kx
		}
		fallthrough
	case 7:
		oy = args[6]
		fallthrough
	case 6:
		ox = args[5]
		if args_length == 6 {
			oy = ox
		}
		fallthrough
	case 5:
		sy = args[4]
		fallthrough
	case 4:
		sx = args[3]
		if args_length == 4 {
			sy = sx
		}
		fallthrough
	case 3:
		angle = args[2]
		fallthrough
	case 2:
		x = args[0]
		y = args[1]
	}

	return x, y, angle, sx, sy, ox, oy, kx, ky
}
