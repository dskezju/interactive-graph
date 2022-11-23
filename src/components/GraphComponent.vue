<template>
  <el-container>
    <el-aside style="width: 360px; background-color: #ffffff">
      <LeftPanel v-bind:attributes="attributes" />
    </el-aside>
    <el-main style="width: 100%; height: 100%; padding: 0px">
      <div id="sigma-container" style="width: 100vw; height: 100vh"></div>
      <div id="controls">
        <div class="input">
          <label for="zoom-in">Zoom in</label><button id="zoom-in">+</button>
        </div>
        <div class="input">
          <label for="zoom-out">Zoom out</label><button id="zoom-out">-</button>
        </div>
        <div class="input">
          <label for="zoom-reset">Reset zoom</label>
          <button id="zoom-reset">⊙</button>
        </div>
        <div class="input">
          <label for="labels-threshold">Labels threshold</label>
          <input
            id="labels-threshold"
            type="range"
            min="0"
            max="15"
            step="0.5"
          />
        </div>
      </div>
      <div id="search">
        <input
          type="search"
          id="search-input"
          list="suggestions"
          placeholder="Try searching for a node..."
        />
        <datalist id="suggestions"></datalist>
      </div>
      <div id="nodeContextMenu" class="contextMenu">
        <div>
          <button id="nodeDelete" @click="handleNodeDeleteClick">delete</button>
          <button id="edgeAdd" @click="handleEdgeAddClick">add edge</button>
        </div>
      </div>
      <div id="edgeContextMenu" class="contextMenu">
        <div>
          <button id="edgeDelete" @click="handleEdgeDeleteClick">delete</button>
        </div>
      </div>
      <div id="stageContextMenu" class="contextMenu">
        <div>
          <button id="nodeAdd" @click="handleNodeAddClick">add node</button>
        </div>
      </div>
    </el-main>
  </el-container>
</template>

<script lang="ts">
import Sigma from "sigma";
import Graph, { MultiGraph } from "graphology";
import chroma from "chroma-js";
import NodeProgramFull from "./webgl/programs/node.full";
import getNodeProgramImage from "sigma/rendering/webgl/programs/node.image";
import FA2Layout from "graphology-layout-forceatlas2/worker";
import forceAtlas2 from "graphology-layout-forceatlas2";
import { Coordinates, EdgeDisplayData, NodeDisplayData } from "sigma/types";
import circlepack from "graphology-layout/circlepack";
import circular from "graphology-layout/circular";
import { layoutAnimate } from "@/lib/layoutAnimation";
import { drawHover } from "@/utils/canvas";
import LeftPanel from "@/components/LeftPanel.vue";
import { defineComponent, ref } from "vue";

import store from "@/store";
import axios from "axios";

import { BACKEND } from "@/config";

function colorize(str: string) {
  for (
    var i = 0, hash = 0;
    i < str.length;
    hash = str.charCodeAt(i++) + ((hash << 5) - hash)
  );
  let color = Math.floor(
    Math.abs(((Math.sin(hash) * 10000) % 1) * 16777216)
  ).toString(16);
  return "#" + Array(6 - color.length + 1).join("0") + color;
}

// Type and declare internal state:
interface State {
  hoveredNode?: string;
  hoveredEdge?: string;
  searchQuery: string;

  // State derived from query:
  selectedNode?: string;
  selectedEdge?: string;
  suggestions?: Set<string>;

  // State derived from hovered node:
  hoveredNeighbors?: Set<string>;

  isDragging?: boolean;
  nodeToConnect?: string;
}
const state: State = {
  searchQuery: "",
  isDragging: false,
};

export default defineComponent({
  name: "GraphComponent",
  components: {
    LeftPanel,
  },
  data() {
    return {
      graph: new Graph(),
      attributes: [],
      fa2layout: new FA2Layout(new Graph()),
      stageContextMenu: ref<HTMLElement>(),
      nodeContextMenu: ref<HTMLElement>(),
      edgeContextMenu: ref<HTMLElement>(),
      renderer: ref<Sigma>(),
    };
  },
  created() {
    this.initGraph();
  },
  mounted() {
    this.stageContextMenu = document.getElementById(
      "stageContextMenu"
    ) as HTMLElement;
    this.nodeContextMenu = document.getElementById(
      "nodeContextMenu"
    ) as HTMLElement;
    this.edgeContextMenu = document.getElementById(
      "edgeContextMenu"
    ) as HTMLElement;
  },
  methods: {
    initGraph() {
      /* *****2022.11.9****** */
      axios({
        method: "GET",
        url: BACKEND + "/graph/",
      })
        .then((res) => res.data)
        .then((jsonObj) => {
          const graph = new MultiGraph();
          this.graph = graph;
          store.dispatch("set", {
            key: "graph",
            value: graph,
          });

          graph.import(jsonObj);

          store.dispatch("set", {
            key: "graphNodeCount",
            value: graph.nodes().length,
          });

          store.dispatch("set", {
            key: "graphEdgeCount",
            value: graph.edges().length,
          });

          graph.forEachNode((node, attr) => {
            attr.color = chroma(colorize(attr["labels"][0])).hex();
            attr.label =
              attr["productName"] ||
              attr["companyName"] ||
              attr["shipName"] ||
              attr["categoryName"] ||
              attr["labels"];
            attr.size = attr["reorderLevel"] / 5;
            return attr;
          });

          circular.assign(graph);

          // Retrieve some useful DOM elements:
          const container = document.getElementById(
            "sigma-container"
          ) as HTMLElement;
          const zoomInBtn = document.getElementById(
            "zoom-in"
          ) as HTMLButtonElement;
          const zoomOutBtn = document.getElementById(
            "zoom-out"
          ) as HTMLButtonElement;
          const zoomResetBtn = document.getElementById(
            "zoom-reset"
          ) as HTMLButtonElement;
          const labelsThresholdRange = document.getElementById(
            "labels-threshold"
          ) as HTMLInputElement;

          // Instanciate sigma:
          const renderer = new Sigma(graph, container, {
            minCameraRatio: 0.001,
            maxCameraRatio: 1000,
            nodeProgramClasses: {
              image: getNodeProgramImage(),
              circle: NodeProgramFull,
            },
            renderEdgeLabels: true,
            enableEdgeHoverEvents: "debounce",
            enableEdgeClickEvents: true,
          });
          this.renderer = renderer;

          renderer.on("rightClickNode", (e) => {
            this.handleNodeRightClick(e.event);
          });

          renderer.on("rightClickEdge", (e) => {
            this.handleEdgeRightClick(e.event);
          });

          graph.forEachNode((node, attr) => {
            let subtitles: string[] = [];
            for (const [key, value] of Object.entries(attr)) {
              subtitles.push(
                `${key}: ${
                  typeof value === "number" ? value.toLocaleString() : value
                }`
              );
            }
            attr.subtitles = subtitles;
            return attr;
          });

          renderer.setSetting("hoverRenderer", (context, data, settings) =>
            drawHover(
              context,
              { ...renderer.getNodeDisplayData(data.key), ...data },
              settings
            )
          );

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
          // Create node (and edge) by click
          // ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
          //

          renderer.on("clickStage", () => {
            this.handleStageClick();
          });

          // When clicking on the stage, we add a new node and connect it to the closest node
          renderer.on("rightClickStage", () => {
            // // Sigma (ie. graph) and screen (viewport) coordinates are not the same.
            // // So we need to translate the screen x & y coordinates to the graph one by calling the sigma helper `viewportToGraph`
            // const coordForGraph = renderer.viewportToGraph({
            //   x: event.x,
            //   y: event.y,
            // });
            // // We create a new node
            // const node = {
            //   ...coordForGraph,
            //   size: 4,
            //   color: chroma.random().hex(),
            //   // type: "border",
            // };
            // // Searching the two closest nodes to auto-create an edge to it
            // const closestNodes = graph
            //   .nodes()
            //   .map((nodeId) => {
            //     const attrs = graph.getNodeAttributes(nodeId);
            //     const distance =
            //       Math.pow(node.x - attrs.x, 2) + Math.pow(node.y - attrs.y, 2);
            //     return { nodeId, distance };
            //   })
            //   .sort((a, b) => a.distance - b.distance)
            //   .slice(0, 2);
            // // We register the new node into graphology instance
            // const id = uuid();
            // graph.addNode(id, node);
            // // We create the edges
            // closestNodes.forEach((e) => graph.addEdge(id, e.nodeId));
            /* Input:  */
            /* add node test */
            // axios({
            //   method: "POST",
            //   url: "http://localhost:8083/graph/node/",
            //   data: {
            //     method: "add",
            //     payload: {
            //       // key: 1,
            //       attributes: {
            //         labels: "Label_test",
            //         attribute1: "attribute1",
            //         attribute2: "attribute2",
            //         attribute3: "attribute3",
            //       },
            //     },
            //   },
            // });
            /* delete node and its relationships by key test */
            // axios({
            //   method: "POST",
            //   url: "http://localhost:8083/node",
            //   data: {
            //     method: "delete",
            //     payload: {
            //       key: 1046,
            //     },
            //   },
            // });
            /* update node by key test */
            // axios({
            //   method: "POST",
            //   url: "http://localhost:8083/node",
            //   data: {
            //     method: "update",
            //     payload: {
            //       key: 1045,
            //       attributes: {
            //         labels: "Label_update",
            //         attribute1: "attribute1_update",
            //         attribute2: "attribute2",
            //         // attribute3: "attribute3",
            //       },
            //     },
            //   },
            // });
            /* add relation test */
            axios({
              method: "POST",
              url: "http://localhost:8083/graph/edge/",
              data: {
                method: "add",
                payload: {
                  source: 1048,
                  target: 1049,
                  attributes: {
                    // one edge 'type' must be specified
                    type: "TYPE_TEST",
                    attribute1: "attribute1",
                  },
                },
              },
            });

            /* delete relation by key test */
            // axios({
            //   method: "POST",
            //   url: "http://localhost:8083/graph/edge/",
            //   data: {
            //     method: "delete",
            //     payload: {
            //       key: 3100,
            //     },
            //   },
            // });

            /* update relation by key test */
            // axios({
            //   method: "POST",
            //   url: "http://localhost:8083/graph/edge/",
            //   data: {
            //     method: "update",
            //     payload: {
            //       key: 3103,
            //       // source: 1048,
            //       // target: 1049,
            //       attributes: {
            //         attribute1: "newattribute",
            //       },
            //     },
            //   },
            // });
            /* delete relation by key test */
            // axios({
            //   method: "POST",
            //   url: "http://localhost:8083/graph/edge/",
            //   data: {
            //     method: "delete",
            //     payload: {
            //       key: 3100,
            //     },
            //   },
            // });
          });

          renderer.on("enterEdge", ({ edge }) => {
            state.hoveredEdge = edge;
            renderer.refresh();
          });
          renderer.on("leaveEdge", () => {
            state.hoveredEdge = undefined;
            renderer.refresh();
          });
          renderer.on("clickEdge", ({ edge }) => {
            state.selectedEdge = edge;
            store.dispatch("set", {
              key: "graphItemSelected",
              value: { type: "edge", id: edge },
            });
            renderer.refresh();
          });

          //
          // highlight and search
          // ~~~~~~~~~~~~~~~~~~~~
          //

          const searchInput = document.getElementById(
            "search-input"
          ) as HTMLInputElement;
          const searchSuggestions = document.getElementById(
            "suggestions"
          ) as HTMLDataListElement;

          // Feed the datalist autocomplete values:
          searchSuggestions.innerHTML = graph
            .nodes()
            .map(
              (node) =>
                `<option value="${graph.getNodeAttribute(
                  node,
                  "label"
                )}"></option>`
            )
            .join("\n");

          // Actions:
          function setSearchQuery(query: string) {
            state.searchQuery = query;

            if (searchInput.value !== query) searchInput.value = query;

            if (query) {
              const lcQuery = query.toLowerCase();
              const suggestions = graph
                .nodes()
                .map((n) => ({
                  id: n,
                  label: graph.getNodeAttribute(n, "label") as string,
                }))
                .filter(({ label }) =>
                  ("" + label).toLowerCase().includes(lcQuery)
                );

              // If we have a single perfect match, them we remove the suggestions, and
              // we consider the user has selected a node through the datalist
              // autocomplete:
              if (suggestions.length === 1 && suggestions[0].label === query) {
                state.selectedNode = suggestions[0].id;
                state.suggestions = undefined;

                // Move the camera to center it on the selected node:
                const nodePosition = renderer.getNodeDisplayData(
                  state.selectedNode
                ) as Coordinates;
                renderer.getCamera().animate(nodePosition, {
                  duration: 500,
                });
              }
              // Else, we display the suggestions list:
              else {
                state.selectedNode = undefined;
                state.suggestions = new Set(suggestions.map(({ id }) => id));
              }
            }
            // If the query is empty, then we reset the selectedNode / suggestions state:
            else {
              state.selectedNode = undefined;
              state.suggestions = undefined;
            }

            // Refresh rendering:
            renderer.refresh();
          }
          function setHoveredNode(node?: string) {
            if (node) {
              state.hoveredNode = node;
              state.hoveredNeighbors = new Set(graph.neighbors(node));
            } else {
              state.hoveredNode = undefined;
              state.hoveredNeighbors = undefined;
            }

            // Refresh rendering:
            renderer.refresh();
          }

          // Bind search input interactions:
          searchInput.addEventListener("input", () => {
            setSearchQuery(searchInput.value || "");
          });
          searchInput.addEventListener("blur", () => {
            setSearchQuery("");
          });

          // Bind graph interactions:
          renderer.on("enterNode", ({ node }) => {
            if (!state.isDragging) {
              setHoveredNode(node);
            }
          });
          renderer.on("leaveNode", () => {
            if (!state.isDragging) {
              setHoveredNode(undefined);
            }
          });

          // Render nodes accordingly to the internal state:
          // 1. If a node is selected, it is highlighted
          // 2. If there is query, all non-matching nodes are greyed
          // 3. If there is a hovered node, all non-neighbor nodes are greyed
          renderer.setSetting("nodeReducer", (node, data) => {
            const res: Partial<NodeDisplayData> = { ...data };

            if (
              state.hoveredNeighbors &&
              !state.hoveredNeighbors.has(node) &&
              state.hoveredNode !== node
            ) {
              res.label = "";
              res.color = "#f6f6f6";
            }

            if (state.selectedNode === node) {
              res.highlighted = true;
            } else if (state.suggestions && !state.suggestions.has(node)) {
              res.label = "";
              res.color = "#f6f6f6";
            }

            return res;
          });

          // Render edges accordingly to the internal state:
          // 1. If a node is hovered, the edge is hidden if it is not connected to the
          //    node
          // 2. If there is a query, the edge is only visible if it connects two
          //    suggestions
          renderer.setSetting("edgeReducer", (edge, data) => {
            const res: Partial<EdgeDisplayData> = { ...data };

            if (
              state.hoveredNode &&
              !graph.hasExtremity(edge, state.hoveredNode)
            ) {
              res.hidden = true;
            }

            if (
              state.suggestions &&
              (!state.suggestions.has(graph.source(edge)) ||
                !state.suggestions.has(graph.target(edge)))
            ) {
              res.hidden = true;
            }

            if (edge == state.hoveredEdge || edge == state.selectedEdge) {
              res.color = "#cc0000";
            }

            return res;
          });

          //
          // Drag'n'drop feature
          // ~~~~~~~~~~~~~~~~~~~
          //

          // State for drag'n'drop
          let draggedNode: string | null = null;
          state.isDragging = false;

          // On mouse down on a node
          //  - we enable the drag mode
          //  - save in the dragged node in the state
          //  - highlight the node
          //  - disable the camera so its state is not updated
          renderer.on("downNode", (e) => {
            state.isDragging = true;
            state.selectedNode = e.node;
            store.dispatch("set", {
              key: "graphItemSelected",
              value: { type: "node", id: e.node },
            });
            this.fa2layout.stop();
            draggedNode = e.node;
            graph.setNodeAttribute(draggedNode, "highlighted", true);
          });

          renderer.on("clickNode", (e) => {
            this.handleNodeClick(e.node);
          });

          // On mouse move, if the drag mode is enabled, we change the position of the draggedNode
          renderer.getMouseCaptor().on("mousemovebody", (e) => {
            if (!state.isDragging || !draggedNode) return;

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
            state.isDragging = false;
            // layout.start();
            draggedNode = null;
          });

          // Disable the autoscale at the first down interaction
          renderer.getMouseCaptor().on("mousedown", () => {
            if (!renderer.getCustomBBox())
              renderer.setCustomBBox(renderer.getBBox());
          });

          renderer.on("clickStage", () => {
            state.selectedNode = undefined;
            state.selectedEdge = undefined;
            store.dispatch("set", {
              key: "graphItemSelected",
              value: null,
            });
            renderer.refresh();
          });

          renderer.on("rightClickStage", (e) => {
            this.handleStageRightClick(e.event);
          });
        });
    },
    handleGraphLayoutChange(layout: string) {
      if (layout == "circle") {
        this.fa2layout.kill();
        layoutAnimate(this.graph, circular(this.graph));
      } else if (layout == "cluster") {
        this.fa2layout.kill();
        layoutAnimate(
          this.graph,
          circlepack(this.graph, {
            hierarchyAttributes: ["labels"],
            scale: 0.005,
          })
        );
      } else if (layout == "force") {
        const sensibleSettings = forceAtlas2.inferSettings(this.graph);
        const fa2layout = new FA2Layout(this.graph, {
          settings: sensibleSettings,
        });
        this.fa2layout = fa2layout;
        fa2layout.start();
      }
      return 0;
    },
    handleNodeRightClick(e) {
      e.original.preventDefault();
      if (this.nodeContextMenu) {
        this.nodeContextMenu.style.display = "initial";
        this.nodeContextMenu.style.top = e.original.pageY + "px";
        this.nodeContextMenu.style.left = e.original.pageX + "px";
      }
    },
    handleEdgeRightClick(e) {
      e.original.preventDefault();
      if (this.edgeContextMenu) {
        this.edgeContextMenu.style.display = "initial";
        this.edgeContextMenu.style.top = e.original.pageY + "px";
        this.edgeContextMenu.style.left = e.original.pageX + "px";
      }
    },
    handleStageRightClick(e) {
      e.original.preventDefault();
      store.dispatch("set", {
        key: "graphRightClickPosition",
        value: { x: e.x, y: e.y },
      });
      if (this.stageContextMenu) {
        this.stageContextMenu.style.display = "initial";
        this.stageContextMenu.style.top = e.original.pageY + "px";
        this.stageContextMenu.style.left = e.original.pageX + "px";
      }
    },
    handleStageClick() {
      if (this.nodeContextMenu) {
        this.nodeContextMenu.style.display = "none";
      }
      if (this.stageContextMenu) {
        this.stageContextMenu.style.display = "none";
      }
    },
    handleNodeDeleteClick() {
      const graphItemSelected = store.state.graphItemSelected;
      if (graphItemSelected == null) {
        return;
      }
      if (this.nodeContextMenu) {
        this.nodeContextMenu.style.display = "none";
      }

      axios({
        method: "POST",
        url: BACKEND + "/graph/node/",
        data: {
          method: "delete",
          payload: {
            key: +graphItemSelected["id"],
          },
        },
      }).then(() => {
        /// backend responds with success
        store.dispatch("set", {
          key: "graphItemSelected",
          value: null,
        });
        this.graph.dropNode(graphItemSelected["id"]);

        store.dispatch("decrement", {
          key: "graphNodeCount",
          value: null,
        });
      });
    },
    handleNodeAddClick() {
      // Sigma (ie. graph) and screen (viewport) coordinates are not the same.
      // So we need to translate the screen x & y coordinates to the graph one by calling the sigma helper `viewportToGraph`
      if (this.renderer) {
        const coordForGraph = this.renderer.viewportToGraph(
          store.state.graphRightClickPosition
        );

        const node = {
          ...coordForGraph,
          color: chroma(colorize("new")).hex(),
          labels: ["new"],
        };

        axios({
          method: "POST",
          url: BACKEND + "/graph/node/",
          data: {
            method: "add",
            payload: {
              attributes: {},
            },
          },
        }).then((rsp) => {
          if (rsp.data.success) {
            const id = rsp.data.message as number;
            this.graph.addNode(id, node);
          }
          store.dispatch("increment", {
            key: "graphNodeCount",
            value: null,
          });
        });

        if (this.stageContextMenu) {
          this.stageContextMenu.style.display = "none";
        }
      }
    },
    handleEdgeDeleteClick() {
      const graphItemSelected = store.state.graphItemSelected;
      if (graphItemSelected == null) {
        return;
      }
      if (this.edgeContextMenu) {
        this.edgeContextMenu.style.display = "none";
      }

      axios({
        method: "POST",
        url: BACKEND + "/graph/edge/",
        data: {
          method: "delete",
          payload: {
            key: +graphItemSelected["id"],
          },
        },
      }).then(() => {
        /// backend responds with success
        store.dispatch("set", {
          key: "graphItemSelected",
          value: null,
        });
        this.graph.dropEdge(graphItemSelected["id"]);

        store.dispatch("decrement", {
          key: "graphEdgeCount",
          value: null,
        });
      });
    },
    handleEdgeAddClick() {
      if (this.nodeContextMenu) {
        this.nodeContextMenu.style.display = "none";
      }

      const graphItemSelected = store.state.graphItemSelected;
      if (graphItemSelected == null) {
        return;
      }
      state.nodeToConnect = graphItemSelected.id;
    },
    handleNodeClick(node) {
      if (state.nodeToConnect) {
        axios({
          method: "POST",
          url: BACKEND + "/graph/edge/",
          data: {
            method: "add",
            payload: {
              source: +state.nodeToConnect,
              target: +node,
              attributes: {
                type: "new",
              },
            },
          },
        }).then(() => {
          this.graph.addEdge(state.nodeToConnect, node);
          state.nodeToConnect = undefined;
          store.dispatch("increment", {
            key: "graphEdgeCount",
            value: null,
          });
        });
      }
    },
  },
  computed: {
    getGraphLayout() {
      return store.state.graphLayout;
    },
  },
  watch: {
    getGraphLayout() {
      this.handleGraphLayoutChange(store.state.graphLayout);
      return;
    },
  },
});
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

#search {
  position: absolute;
  right: 1em;
  top: 4em;
}

.contextMenu {
  display: none;
  position: absolute;
  width: 200px;
  background-color: white;
  box-shadow: 0 0 5px grey;
  border-radius: 3px;
}

.contextMenu button {
  width: 100%;
  background-color: white;
  border: none;
  margin: 0;
  padding: 10px;
  text-align: left;
}

.contextMenu button:hover {
  background-color: lightgray;
}

.chooseform {
  margin-top: 5px;
}
</style>
