vec4 effect(vec4 vcolor, Image texture, vec2 texcoord, vec2 pixel_coords) {	
  vec4 pixel = Texel(texture, texcoord );//This is the current pixel color
  number average = (pixel.r+pixel.b+pixel.g)/3.0;
  pixel.r = average;
  pixel.g = average;
  pixel.b = average;
  return pixel;
}
