import Graph from "graphology";
import { createStore } from "vuex";

interface State {
  graphNodeCount: number;
  graphEdgeCount: number;
  graphLayout: string;
  graph: Graph | null;
  graphItemSelected: { type: string; id: string } | null;
  graphRightClickPosition: { x: number; y: number };
}

export default createStore<State>({
  state: {
    graphNodeCount: 0,
    graphEdgeCount: 0,
    graphLayout: "",
    graph: null,
    graphItemSelected: null,
    graphRightClickPosition: { x: 0, y: 0 },
  },
  getters: {},
  mutations: {
    SET(state, payload) {
      state[payload["key"]] = payload["value"];
    },

    INCREMENT(state, payload) {
      ++state[payload["key"]];
    },

    DECREMENT(state, payload) {
      --state[payload["key"]];
    },
  },
  actions: {
    set(context, payload) {
      context.commit("SET", payload);
    },

    increment(context, payload) {
      context.commit("INCREMENT", payload);
    },

    decrement(context, payload) {
      context.commit("DECREMENT", payload);
    },
  },
  modules: {},
});
