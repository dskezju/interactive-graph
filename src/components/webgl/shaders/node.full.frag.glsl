precision mediump float;

varying vec4 v_color;
varying float v_border;
varying float v_ring1_width;

const vec4 transparent = vec4(0.0, 0.0, 0.0, 0.0);
const vec4 white = vec4(1.0, 1.0, 1.0, 1.0);
const float radius_2 = 0.5;

void main(void) {
  float radius_1 = radius_2 - v_ring1_width;
  vec4 v_color_light = v_color * 1.2;
  v_color_light[3] = v_color[3];
  vec4 v_color_dark = v_color * 0.8;
  v_color_dark[3] = v_color[3];
  float distToCenter = length(gl_PointCoord - vec2(0.5, 0.5));

  float t = 0.0;
  if (distToCenter < radius_1 - v_border)
    gl_FragColor = v_color_light;
  else if (distToCenter < radius_1)
    gl_FragColor = mix(v_color_dark, v_color_light, (radius_1 - distToCenter) / v_border);
  else if (distToCenter < radius_2 - v_border)
    gl_FragColor = v_color_dark;
  else if (distToCenter < radius_2)
    gl_FragColor = mix(transparent, v_color_dark, (radius_2 - distToCenter) / v_border);
  else
    gl_FragColor = transparent;
}
