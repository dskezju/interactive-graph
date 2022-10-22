<template>
  <GraphComponent />
</template>

<script lang="ts">
import { defineComponent } from "vue";
import GraphComponent from "@/components/GraphComponent.vue"; // @ is an alias to /src

export default defineComponent({
  name: "GraphView",
  components: {
    GraphComponent,
  },
});
</script>
