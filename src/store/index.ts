import { createStore } from "vuex";

export default createStore({
  state: {
    graphNodeCount: 0,
    graphEdgeCount: 0,
    graphLayout: "",
    graph: null,
    graphNodeSelected: null,
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
