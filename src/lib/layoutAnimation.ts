import Graph from "graphology";

type LayoutMapping = { [key: string]: { x: number; y: number } };

type AnimationSettings = {
  duration?: number;
};

const defaultSettings = {
  duration: 500,
};

export function layoutAnimate(
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
  const duration = s.duration || defaultSettings.duration;

  let finished = false;
  const startTime = Date.now();

  const animate = function () {
    const timeProgress = Math.min((Date.now() - startTime) / duration, 1);
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

    if (!finished) {
      if (timeProgress == 1) {
        finished = true;
      }
      window.requestAnimationFrame(animate);
    }
  };

  window.requestAnimationFrame(animate);
}
