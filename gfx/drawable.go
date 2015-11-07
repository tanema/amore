package gfx

type Drawable interface {
	/**
	 * Draws the object with the specified transformation.
	 *
	 * @param x The position of the object along the x-axis.
	 * @param y The position of the object along the y-axis.
	 * @param angle The angle of the object (in radians).
	 * @param sx The scale factor along the x-axis.
	 * @param sy The scale factor along the y-axis.
	 * @param ox The origin offset along the x-axis.
	 * @param oy The origin offset along the y-axis.
	 * @param kx Shear along the x-axis.
	 * @param ky Shear along the y-axis.
	 **/
	Draw(args ...float32)
}
