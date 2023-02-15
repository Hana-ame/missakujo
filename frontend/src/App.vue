<script setup lang="ts">
import { ref, computed } from "vue"
import { NSpace, NInput, NDatePicker, NCheckbox } from "naive-ui"
// import Date from "date"

const acct = ref("")
const userId = ref("")
const host = ref("")
const token = ref("")
// const acct = ref("")

const range = ref<[number,number]>([Date.now()-1000*3600*24*7,Date.now()])

const deleteRenotes = ref(true)
const deleteReplies = ref(true)


let renoteLessThan : number = 0
let renoteLessThanStr : string = ""


const renotesLessThan = computed<string>({
  get() : string{
    return renoteLessThanStr
  },
  set(newValue : string){
    renoteLessThanStr = newValue
    renoteLessThan = Number(newValue)
  }
})

</script>

<template>
  
  <n-space vertical>
    <n-input
      v-model:value="acct"
      type="text"
      placeholder="username@example.com"
    />
    {{ acct }}
    <n-space>
      <n-input
        v-model:value="userId"
        type="text"
        placeholder="userId"
      />
      <n-input
        v-model:value="host"
        type="text"
        placeholder="host"
      />
    </n-space>
    {{ userId +"@"+ host }}
    <n-input
      v-model:value="token"
      type="text"
      placeholder="token"
    />
    {{ token }}
    <n-date-picker v-model:value="range" type="datetimerange" clearable />
    <pre>{{ JSON.stringify(range) }}</pre>

    <n-space item-style="display: flex; margin:auto;" >
      <n-checkbox size="large" v-model:checked="deleteRenotes" label="DeleteRenotes" />
      <n-checkbox size="large" v-model:checked="deleteReplies" label="DeleteReplies" />
    </n-space>
    {{ deleteRenotes }}
    {{ deleteReplies }}
    <n-input
      v-model:value="renotesLessThan"
      type="text"
      placeholder="999"
    />
    {{ renotesLessThan }}
    {{ renoteLessThan }}
  </n-space>


</template>

<style scoped>
.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}
.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}
.logo.vue:hover {
  filter: drop-shadow(0 0 2em #42b883aa);
}
</style>
