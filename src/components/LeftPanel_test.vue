<template>
  <el-tabs
    type="border-card"
    v-model="activeName"
    class="demo-tabs"
    @tab-change="handleClick"
    style="margin-top: 10px; margin-left: 10px"
  >
    <el-tab-pane label="Edit" name="Edit">
      <div class="titleText">Knowledge Graph visiualize</div>
      <div class="normaltext">
        To change settings in the following part to better illustrate graph
        informatoion.
      </div>

      <el-divider> </el-divider>

      <div class="informtext">Which tag will be visiualized</div>
      <el-row style="margin-top: 10px; font-size: 14px">
        <el-col :span="8"> Tag </el-col>
        <el-col :span="8">Color</el-col>
        <el-col :span="8">Size</el-col>
      </el-row>
      <div v-for="(label, i) in attributes" :key="i" class="chooseform">
        <div style="font-size: 14px; margin-bottom: 4px; font-weight: bold">
          {{ i }}
        </div>
        <el-row v-for="(tag, index) in label" :key="index">
          <el-col :span="8" style="margin-top: 6px">{{ tag }} </el-col>
          <el-col :span="8"
            ><el-switch v-model="switchColor[i + '-' + tag]"
          /></el-col>
          <el-col :span="8"
            ><el-switch v-model="switchSize[i + '-' + tag]"
          /></el-col>
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
    </el-tab-pane>
    <el-tab-pane label="Visiualize" name="Visiualize">
      <div class="titleText">City Knowledge Graph</div>
      <div class="informtext">Node: 2000</div>
      <div class="informtext">Edges: 6290</div>

      <el-divider> </el-divider>

      <el-row>
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
          <el-collapse>
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
      </el-row>
    </el-tab-pane>
  </el-tabs>
</template>

<script lang="ts">
import { defineComponent } from "vue";

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
      activeName: "Edit",
      value: [],
      tests: ["Tag1", "Tag2", "Tag3", "Tag4", "Tag5"],
      options: [
        { value: "1", label: "Tag1" },
        { value: "2", label: "Tag2" },
        { value: "3", label: "Tag3" },
        { value: "4", label: "Tag4" },
        { value: "5", label: "Tag5" },
      ],
      value2: [],
      options2: [
        { value: "1", label: "Default - All edges are gray" },
        { value: "2", label: "Use original color" },
        { value: "3", label: "Use the source node color" },
        { value: "4", label: "Use the target node color" },
      ],
      switchColor: {},
      switchSize: {},
    };
  },
  created() {
    for (var label in this.attributes) {
      for (var attr in this.attributes[label]) {
        this.switchColor[label + "-" + attr] = false;
        this.switchSize[label + "-" + attr] = false;
      }
    }
  },
  methods: {
    handleClick(name) {
      // console.log(name);
      // if (name == "Visiualize") {
      // }
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
