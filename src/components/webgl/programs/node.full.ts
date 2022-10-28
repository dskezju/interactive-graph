/**
 * This class copies sigma/rendering/webgl/programs/node.fast
 */
import { floatColor } from "sigma/utils";
import { NodeDisplayData } from "sigma/types";
import { AbstractNodeProgram } from "sigma/rendering/webgl/programs/common/node";
import { RenderParams } from "sigma/rendering/webgl/programs/common/program";

import vertexShaderSource from "!raw-loader!../shaders/node.full.vert.glsl";
import fragmentShaderSource from "!raw-loader!../shaders/node.full.frag.glsl";

const POINTS = 1,
  ATTRIBUTES = 5;

export default class NodeProgramFull extends AbstractNodeProgram {
  constructor(gl: WebGLRenderingContext) {
    super(gl, vertexShaderSource, fragmentShaderSource, POINTS, ATTRIBUTES);
    this.bind();
  }

  process(data: NodeDisplayData, hidden: boolean, offset: number): void {
    const array = this.array;
    let i = offset * POINTS * ATTRIBUTES;

    if (hidden) {
      array[i++] = 0;
      array[i++] = 0;
      array[i++] = 0;
      array[i++] = 0;
      array[i++] = 0;
      return;
    }

    const color = floatColor(data.color);

    array[i++] = data.x;
    array[i++] = data.y;
    array[i++] = data.size;
    array[i++] = color;
    array[i++] = 0.05;
  }

  bind(): void {
    super.bind();
    const gl = this.gl;

    const ring1WidthLocation = gl.getAttribLocation(
      this.program,
      "a_ring1_width"
    );
    gl.enableVertexAttribArray(ring1WidthLocation);
    gl.vertexAttribPointer(
      ring1WidthLocation,
      1,
      gl.FLOAT,
      false,
      this.attributes * Float32Array.BYTES_PER_ELEMENT,
      16
    );
  }

  render(params: RenderParams): void {
    const gl = this.gl;

    const program = this.program;
    gl.useProgram(program);

    gl.uniform1f(this.ratioLocation, 1 / Math.sqrt(params.ratio));
    gl.uniform1f(this.scaleLocation, params.scalingRatio);
    gl.uniformMatrix3fv(this.matrixLocation, false, params.matrix);

    // gl.uniform1f(gl.getUniformLocation(program, "u_ring1_width"), 0.05);

    gl.drawArrays(gl.POINTS, 0, this.array.length / ATTRIBUTES);
  }
}
