<template>
  <div id="sigma-container" style="width: 100vw; height: 100vh"></div>
  <div id="controls">
    <div class="input">
      <label for="zoom-in">Zoom in</label><button id="zoom-in">+</button>
    </div>
    <div class="input">
      <label for="zoom-out">Zoom out</label><button id="zoom-out">-</button>
    </div>
    <div class="input">
      <label for="zoom-reset">Reset zoom</label
      ><button id="zoom-reset">⊙</button>
    </div>
    <div class="input">
      <label for="labels-threshold">Labels threshold</label>
      <input id="labels-threshold" type="range" min="0" max="15" step="0.5" />
    </div>
  </div>
</template>

<script lang="ts">
import Sigma from "sigma";
import Graph from "graphology";
import { parse } from "graphology-gexf/browser";
import { Vue } from "vue-class-component";
import chroma from "chroma-js";
import { v4 as uuid } from "uuid";
import NodeProgramFull from "./webgl/programs/node.full";
import getNodeProgramImage from "sigma/rendering/webgl/programs/node.image";
import FA2Layout from "graphology-layout-forceatlas2/worker";
import forceAtlas2 from "graphology-layout-forceatlas2";

fetch("./arctic.gexf")
  .then((res) => res.text())
  .then((gexf) => {
    // Parse GEXF string:
    const graph = parse(Graph, gexf);

    // Retrieve some useful DOM elements:
    const container = document.getElementById("sigma-container") as HTMLElement;
    const zoomInBtn = document.getElementById("zoom-in") as HTMLButtonElement;
    const zoomOutBtn = document.getElementById("zoom-out") as HTMLButtonElement;
    const zoomResetBtn = document.getElementById(
      "zoom-reset"
    ) as HTMLButtonElement;
    const labelsThresholdRange = document.getElementById(
      "labels-threshold"
    ) as HTMLInputElement;

    // Instanciate sigma:
    const renderer = new Sigma(graph, container, {
      minCameraRatio: 0.1,
      maxCameraRatio: 10,
      nodeProgramClasses: {
        image: getNodeProgramImage(),
        circle: NodeProgramFull,
      },
      renderEdgeLabels: true,
    });
    const camera = renderer.getCamera();

    // Bind zoom manipulation buttons
    zoomInBtn.addEventListener("click", () => {
      camera.animatedZoom({ duration: 600 });
    });
    zoomOutBtn.addEventListener("click", () => {
      camera.animatedUnzoom({ duration: 600 });
    });
    zoomResetBtn.addEventListener("click", () => {
      camera.animatedReset({ duration: 600 });
    });

    // Bind labels threshold to range input
    labelsThresholdRange.addEventListener("input", () => {
      renderer.setSetting(
        "labelRenderedSizeThreshold",
        +labelsThresholdRange.value
      );
    });

    // Set proper range initial value:
    labelsThresholdRange.value =
      renderer.getSetting("labelRenderedSizeThreshold") + "";

    //
    // Drag'n'drop feature
    // ~~~~~~~~~~~~~~~~~~~
    //

    // State for drag'n'drop
    let draggedNode: string | null = null;
    let isDragging = false;

    // On mouse down on a node
    //  - we enable the drag mode
    //  - save in the dragged node in the state
    //  - highlight the node
    //  - disable the camera so its state is not updated
    renderer.on("downNode", (e) => {
      isDragging = true;
      draggedNode = e.node;
      graph.setNodeAttribute(draggedNode, "highlighted", true);
    });

    // On mouse move, if the drag mode is enabled, we change the position of the draggedNode
    renderer.getMouseCaptor().on("mousemovebody", (e) => {
      if (!isDragging || !draggedNode) return;

      // Get new position of node
      const pos = renderer.viewportToGraph(e);

      graph.setNodeAttribute(draggedNode, "x", pos.x);
      graph.setNodeAttribute(draggedNode, "y", pos.y);

      // Prevent sigma to move camera:
      e.preventSigmaDefault();
      e.original.preventDefault();
      e.original.stopPropagation();
    });

    // On mouse up, we reset the autoscale and the dragging mode
    renderer.getMouseCaptor().on("mouseup", () => {
      if (draggedNode) {
        graph.removeNodeAttribute(draggedNode, "highlighted");
      }
      isDragging = false;
      draggedNode = null;
    });

    // Disable the autoscale at the first down interaction
    renderer.getMouseCaptor().on("mousedown", () => {
      if (!renderer.getCustomBBox()) renderer.setCustomBBox(renderer.getBBox());
    });

    //
    // Create node (and edge) by click
    // ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
    //

    // When clicking on the stage, we add a new node and connect it to the closest node
    renderer.on(
      "clickStage",
      ({ event }: { event: { x: number; y: number } }) => {
        // Sigma (ie. graph) and screen (viewport) coordinates are not the same.
        // So we need to translate the screen x & y coordinates to the graph one by calling the sigma helper `viewportToGraph`
        const coordForGraph = renderer.viewportToGraph({
          x: event.x,
          y: event.y,
        });

        // We create a new node
        const node = {
          ...coordForGraph,
          size: 4,
          color: chroma.random().hex(),
          // type: "border",
        };

        // Searching the two closest nodes to auto-create an edge to it
        const closestNodes = graph
          .nodes()
          .map((nodeId) => {
            const attrs = graph.getNodeAttributes(nodeId);
            const distance =
              Math.pow(node.x - attrs.x, 2) + Math.pow(node.y - attrs.y, 2);
            return { nodeId, distance };
          })
          .sort((a, b) => a.distance - b.distance)
          .slice(0, 2);

        // We register the new node into graphology instance
        const id = uuid();
        graph.addNode(id, node);

        // We create the edges
        closestNodes.forEach((e) => graph.addEdge(id, e.nodeId));
      }
    );

    const sensibleSettings = forceAtlas2.inferSettings(graph);
    const layout = new FA2Layout(graph, {
      settings: sensibleSettings,
    });
    layout.start();
  });
export default class GraphComponent extends Vue {}
</script>

<style>
body {
  font-family: sans-serif;
}

html,
body,
#sigma-container {
  width: 100%;
  height: 100%;
  margin: 0;
  padding: 0;
  overflow: hidden;
}

#controls {
  position: absolute;
  right: 1em;
  top: 1em;
  text-align: right;
}

.input {
  position: relative;
  display: inline-block;
  vertical-align: middle;
}

.input:not(:hover) label {
  display: none;
}

.input label {
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  background: black;
  color: white;
  padding: 0.2em;
  border-radius: 2px;
  margin-top: 0.3em;
  font-size: 0.8em;
  white-space: nowrap;
}

.input button {
  width: 2.5em;
  height: 2.5em;
  display: inline-block;
  text-align: center;
  background: white;
  outline: none;
  border: 1px solid dimgrey;
  border-radius: 2px;
  cursor: pointer;
}
</style>
