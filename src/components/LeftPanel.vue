<template>
  <el-tabs
    type="border-card"
    v-model="activeName"
    class="demo-tabs"
    style="margin-top: 10px; margin-left: 10px"
  >
    <!-- <el-tab-pane label="Edit" name="Edit">
      <div class="titleText">Knowledge Graph Visulization</div>
      <div class="normaltext">
        To change settings in the following part to better illustrate graph
        informatoion.
      </div>

      <el-divider> </el-divider>

      <div class="informtext">Which tag will be visiualized</div>
      <el-row style="margin-top: 10px; font-size: 14px">
        <el-col :span="6"> Tag </el-col>
        <el-col :span="6">Color</el-col>
        <el-col :span="6">Size</el-col>
        <el-col :span="6"></el-col>
      </el-row>
      <div v-for="(item, i) in tests" :key="i" class="chooseform">
        <el-row>
          <el-col :span="6" style="margin-top: 6px">{{ item }} </el-col>
          <el-col :span="6"><el-switch v-model="switchColor[item]" /></el-col>
          <el-col :span="6"><el-switch v-model="switchSize[item]" /></el-col>
          <el-col :span="6"></el-col>
        </el-row>
        <el-divider> </el-divider>
      </div>

      <div class="informtext">Choose information to show in the graph</div>
      <el-select-v2
        v-model="value"
        :options="options"
        placeholder="Please select"
        style="width: 240px"
        multiple
        width="300px;"
      />

      <div class="informtext" style="margin-top: 20px">
        Select how the edges will be colored
      </div>
      <el-select-v2
        v-model="value2"
        :options="options2"
        placeholder="Please select"
        default="Default - All edges are gray"
      />
    </el-tab-pane> -->
    <el-tab-pane label="Data" name="Data">
      <div class="titleText">North Wind</div>
      <div class="informtext">Node: {{ getGraphNodeCount() }}</div>
      <div class="informtext">Edges: {{ getGraphEdgeCount() }}</div>

      <el-divider> </el-divider>

      <div v-if="getIsGraphItemSelected()">
        <el-form
          ref="formRef"
          :model="attrs"
          label-width="120px"
          class="demo-dynamic"
        >
          <el-scrollbar max-height="50vh" :noresize="false">
            <div v-for="(dict, i) in attrs" :key="i" class="chooseform">
              <el-row>
                <el-col :span="8"
                  ><el-input
                    v-model="dict.key"
                    class="attr-input"
                    :placeholder="dict.key"
                  />
                </el-col>
                <el-col :span="12">
                  <el-input
                    v-model="dict.value"
                    class="attr-input"
                    :placeholder="dict.value"
                /></el-col>
                <el-col :span="4">
                  <el-button
                    class="attributes-delete-button"
                    @click.prevent="removeDomain(dict)"
                    type="danger"
                    plain
                  >
                    <el-icon style="vertical-align: middle">
                      <Close />
                    </el-icon>
                  </el-button>
                </el-col>
              </el-row>
            </div>
          </el-scrollbar>
        </el-form>
        <el-button @click="addDomain">Add</el-button>
        <el-button type="primary" @click="onSubmit">Update</el-button>
      </div>
      <!-- <el-row>
        <el-col :span="12">
          <div class="informtext">Color nodes select:</div>
          <el-select-v2
            v-model="value2"
            :options="options2"
            placeholder="Please select"
            default="Default - All edges are gray"
          />
        </el-col>
        <el-col :span="12">
          <div class="informtext">Size nodes select:</div>
          <el-select-v2
            v-model="value2"
            :options="options2"
            placeholder="Please select"
            default="Default - All edges are gray"
          />
        </el-col>

        <div class="demo-collapse" style="width: 100%">
          <el-collapse v-model="activeNames" @change="handleChange">
            <el-collapse-item title="Tags" name="1">
              <div>
                Consistent with real life: in line with the process and logic of
                real life, and comply with languages and habits that the users
                are used to;
              </div>
              <div>
                Consistent within interface: all elements should be consistent,
                such as: design style, icons and texts, position of elements,
                etc.
              </div>
            </el-collapse-item>
            <el-collapse-item title="Tag1" name="2">
              <div>
                Operation feedback: enable the users to clearly perceive their
                operations by style updates and interactive effects;
              </div>
              <div>
                Visual feedback: reflect current state by updating or
                rearranging elements of the page.
              </div>
            </el-collapse-item>
            <el-collapse-item title="Tag2" name="3">
              <div>
                Simplify the process: keep operating process simple and
                intuitive;
              </div>
              <div>
                Definite and clear: enunciate your intentions clearly so that
                the users can quickly understand and make decisions;
              </div>
              <div>
                Easy to identify: the interface should be straightforward, which
                helps the users to identify and frees them from memorizing and
                recalling.
              </div>
            </el-collapse-item>
            <el-collapse-item title="Tag3" name="4">
              <div>
                Decision making: giving advices about operations is acceptable,
                but do not make decisions for the users;
              </div>
              <div>
                Controlled consequences: users should be granted the freedom to
                operate, including canceling, aborting or terminating current
                operation.
              </div>
            </el-collapse-item>
          </el-collapse>
        </div>
      </el-row> -->
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
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";
import store from "@/store";
import type { FormInstance } from "element-plus";
import { Close } from "@element-plus/icons-vue";

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
      type: Array,
      required: true,
    },
  },
  data() {
    return {
      activeName: "Data",
      layoutSelected: ref(),
      layoutToSelect: layoutToSelect,
      formRef: ref<FormInstance>(),
      attrs: ref(),
      isToUpdateScrollBar: false,
    };
  },
  methods: {
    addDomain() {
      this.attrs.push({
        key: "",
        value: "",
      });
      this.isToUpdateScrollBar = !this.isToUpdateScrollBar;
    },
    removeDomain(dict) {
      const index = this.attrs.indexOf(dict);
      if (index !== -1) {
        this.attrs.splice(index, 1);
      }
    },
    getGraphNodeCount() {
      return store.state.graphNodeCount;
    },
    getGraphEdgeCount() {
      return store.state.graphEdgeCount;
    },

    getIsGraphItemSelected() {
      this.attrs = this.getSelectedItemAttributes();
      return store.state.graphItemSelected;
    },
    getSelectedItemAttributes() {
      if (store.state.graph && store.state.graphItemSelected) {
        if (store.state.graphItemSelected["type"] == "node") {
          const attributes = store.state.graph.getNodeAttributes(
            store.state.graphItemSelected["id"]
          );

          this.attrs = Object.entries(attributes).map(([key, value]) => ({
            key: key,
            value: value.toString(),
          }));

          return this.attrs;
        } else if (store.state.graphItemSelected["type"] == "edge") {
          const attributes = store.state.graph.getEdgeAttributes(
            store.state.graphItemSelected["id"]
          );

          this.attrs = Object.entries(attributes).map(([key, value]) => ({
            key: key,
            value: value,
          }));

          return this.attrs;
        }
      }
    },
    handelLayoutChange(layoutId) {
      store.dispatch("set", {
        key: "graphLayout",
        value: layoutToSelect[layoutId].label,
      });
    },
    onSubmit() {
      console.log(this.attrs);
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

.attr-input :first-child {
  box-shadow: none;
}

.el-button.attributes-delete-button {
  border: none;
  color: transparent;
  background-color: transparent;
}

.el-row:hover > div > .el-button.attributes-delete-button {
  color: #ddd;
  background-color: transparent;
}

.el-row:hover > div > .el-button.attributes-delete-button:hover {
  color: #a00;
}

.el-button.attributes-delete-button:focus {
  background-color: transparent;
}
</style>
