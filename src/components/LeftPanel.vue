<template>
  <el-collapse v-model="activeNames">
    <el-collapse-item title=" " name="1">
      <el-tabs
        type="border-card"
        v-model="activeName"
        class="demo-tabs"
        style="margin-top: 10px; margin-left: 10px"
      >
        <el-tab-pane label="Data" name="Data">
          <div class="titleText">North Wind</div>
          <div class="informtext">Node: {{ getGraphNodeCount() }}</div>
          <div class="informtext">Edges: {{ getGraphEdgeCount() }}</div>

          <el-divider> </el-divider>

          <el-scrollbar height="500px">
            <div
              v-for="(item, i) in getSelectedNodeAttributes()"
              :key="i"
              class="chooseform"
            >
              <el-row>
                <el-col :span="12">{{ i }} </el-col>
                <el-col :span="12">{{ item }} </el-col>
              </el-row>
            </div>
          </el-scrollbar>
        </el-tab-pane>
        <el-tab-pane label="Layout" name="Layout">
          <div class="titleText">Layout</div>
          <div class="informtext">Different layouts to vsualize the graph</div>
          <el-select-v2
            v-model="layoutSelected"
            :options="layoutToSelect"
            placeholder="Please select"
            style="width: 240px"
            width="300px;"
            @change="handelLayoutChange"
          />
        </el-tab-pane>
      </el-tabs>
    </el-collapse-item>
  </el-collapse>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";
import store from "@/store";

const layoutToSelect = [
  { value: "0", label: "original" },
  { value: "1", label: "circle" },
  { value: "2", label: "cluster" },
  { value: "3", label: "force" },
];

export default defineComponent({
  name: "LeftPanel",
  component: {},
  props: {
    attributes: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      activeName: "Data",
      layoutSelected: ref(),
      layoutToSelect: layoutToSelect,
      activeNames: "1",
    };
  },
  methods: {
    getGraphNodeCount() {
      return store.state.graphNodeCount;
    },
    getGraphEdgeCount() {
      return store.state.graphEdgeCount;
    },
    getSelectedNodeAttributes() {
      if (store.state.graph && store.state.graphNodeSelected) {
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        const attrs = store.state.graph.getNodeAttributes(
          store.state.graphNodeSelected
        );
        return attrs;
      }
    },
    handelLayoutChange(layoutId) {
      store.dispatch("set", {
        key: "graphLayout",
        value: layoutToSelect[layoutId].label,
      });
    },
  },
});
</script>

<style>
.titleText {
  font-size: 20px;
  color: #000000;
  margin-top: 20px;
  margin-bottom: 20px;
}
.normaltext {
  font-size: 14px;
  color: #303133;
}
.informtext {
  font-size: 14px;
  color: #303133;
  margin-top: 4px;
  margin-bottom: 4px;
}
.chooseform {
  font-size: 14px;
  color: #303133;
  margin-top: 20px;
}
</style>
