import Graph from "graphology";

type LayoutMapping = { [key: string]: { x: number; y: number } };

type AnimationSettings = {
  duration?: number;
  fps?: number;
};

const defaultSettings = {
  duration: 500,
  fps: 30,
};

async function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export async function layoutAnimate(
  graph: Graph,
  targetPosition: LayoutMapping,
  settings?: AnimationSettings
) {
  const originalPosition: LayoutMapping = {};

  graph.forEachNode(function (node, attr) {
    originalPosition[node] = {
      x: attr.x,
      y: attr.y,
    };
  });

  const s = settings ? settings : defaultSettings;
  const fps = s.fps || defaultSettings.fps;
  const duration = s.duration || defaultSettings.duration;
  const totalFrame = (fps * duration) / 1000;

  for (let frameIndex = 0; frameIndex < totalFrame; ++frameIndex) {
    const spf = sleep(1000 / fps);

    const timeProgress = frameIndex / totalFrame;
    const linearProgress = timeProgress;
    const bezierProgress =
      linearProgress * linearProgress * (3 - 2 * linearProgress);
    const conjugateBezierProgress = 1 - bezierProgress;

    graph.updateEachNodeAttributes(function (node, attr) {
      // Updating node's positon
      attr.x =
        originalPosition[node].x * conjugateBezierProgress +
        targetPosition[node].x * bezierProgress;

      attr.y =
        originalPosition[node].y * conjugateBezierProgress +
        targetPosition[node].y * bezierProgress;

      return attr;
    });
    await spf;
  }
}
