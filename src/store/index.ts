import Graph from "graphology";
import { createStore } from "vuex";

export default createStore({
  state: {
    graphNodeCount: 0,
    graphEdgeCount: 0,
    graphLayout: "",
    graph: new Graph(),
    graphNodeSelected: -1,
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
