<script setup lang="ts">
import { ref, computed, onMounted } from "vue"
import { NSpace, NInput, NDatePicker, NCheckbox, NButton } from "naive-ui"
// import Date from "date"
const BASE = "/api-missakujo"

const acct = ref("")
const userId = ref("")
const host = ref("")
const token = ref("")
// const acct = ref("")

const range = ref<[number,number]>([Date.now()-1000*3600*24*7,Date.now()])

const deleteRenotes = ref(true)
const deleteReplies = ref(true)


// let renoteLessThan : number = 0
// let renoteLessThanStr : string = ""
const renotesLessThan = ref("999")

// const renotesLessThan = computed<string>({
//   get() : string{
//     return renoteLessThanStr
//   },
//   set(newValue : string){
//     renoteLessThanStr = newValue
//     renoteLessThan = Number(newValue)
//   }
// })

const logHRef = ref("")

function handleDelete(){
  let since = range.value[0]
  let until = range.value[1]
  let rnLessThan = Number(renotesLessThan.value)
  let obj = {    
    "host": host.value.trim(),
    "user": userId.value.trim(),
    "token": token.value.trim(),
    "since": since,
    "until": until,

    "renoteLessThan":rnLessThan,
    "timeOffset":28800,
    "deleteReply":deleteReplies.value,
    "deleteRenote":deleteRenotes.value,

  }
  console.log(JSON.stringify(obj))

  fetch(BASE + "/delete",{
    method:"POST",
    headers:{
      "Content-Type":"application/json"
    },
    body:JSON.stringify(obj)
  }).then(res => res.json()).then(data => {
    if (data.sessionID == "") { return }
    // setInterval(() => {      
      fetch("/log/" + data.sessionID)
        .then(res => res.text())
        .then(data => {
          // logTxt.value = data
        })
      logHRef.value = "/log/" + data.sessionID
    // }, 2000)
  })  
}

const loading = ref(false)

function handleChange() {
  loading.value = true
  userId.value = "Loading..."
  fetch(BASE + "/webfinger/" + acct.value).then(res => res.json()).then(data => {
    userId.value = data.userId
    loading.value = false
  }).catch(err => {
    userId.value = "Failed"
  })

  let arr = acct.value.split("@")
  if (arr.length == 2) {
    host.value = arr.pop()!
  }
  // host.value = acct.value.split("@").pop()!
  if (loading.value) {
    userId.value = ""
  }
}

const uuid = ref("")
const logTxt = ref("")

// onMounted(() => {
//   setInterval(() => {
//     if (uuid.value == "") { return }
//       fetch("/log/" + uuid)
//         .then(res => res.text())
//         .then(data => {
//           logTxt.value = data
//         })
//   }, 2000)
// })


</script>

<template>
  <h1>Missakujo</h1>
  <h2>Delete your misskey notes within a period of time</h2>

  <n-space vertical>
    <n-input
      v-model:value="acct"
      type="text"
      placeholder="username@example.com"
      @update:value="handleChange"
    />
    <!-- <div v-if="loading">loading...</div> -->
    <!-- <n-space justify="space-around" > -->
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
    <!-- </n-space> -->
    <!-- {{ userId +"@"+ host }} -->
    
    <div v-if="host.includes('.')">
    create api token <a :href="'https://'+host+'/settings/api'">here</a>.
    choose 'ノートを作成、削除する'
    </div>
    <n-input
      v-model:value="token"
      type="text"
      placeholder="token"
    />
    
    notes in this perid of time will be deleted
    <n-date-picker v-model:value="range" type="datetimerange" clearable />
    <!-- <pre>{{ JSON.stringify(range) }}</pre> -->

    {{ deleteRenotes ? "will delete renotes" : "renotes will not be deleted" }} |
    {{ deleteReplies ? "will delete replies" : "replies will not be deleted" }}
    <n-space item-style="display: flex; margin:auto;" >
      <n-checkbox size="large" v-model:checked="deleteRenotes" label="DeleteRenotes" />
      <n-checkbox size="large" v-model:checked="deleteReplies" label="DeleteReplies" />
    </n-space>
    renotes less than {{ Number(renotesLessThan) }} will be deleted
    <n-input
      v-model:value="renotesLessThan"
      type="text"
      placeholder="999"
    />
    
    <!-- {{ renoteLessThan }} -->
    <n-button type="primary" style="width: 100%;" @click="handleDelete()">
      Sakujo!
    </n-button>
  </n-space>
  <a v-if="logHRef" :href="logHRef+'.txt'">{{ logHRef }}</a>
  <pre>{{ logTxt }}</pre>

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
