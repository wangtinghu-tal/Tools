<template>
  <div class="upload-container" data-wails-no-drag>
    <div class="hero-overlay bg-opacity-0"></div>
    <div class="hero-content text-center text-neutral-content">
      <div class="max-w-md">
        <div class="file-info-area mt-4 h-56">
          <div v-show="state === STATE_INIT" class="pt-4">
            <h1 class="mb-5 text-5xl font-bold text-slate-900 select-none">
              网络检测工具
            </h1>
            <p v-show="!filePath" class="mb-5 text-gray-400 select-none">
              请在右侧配置中添加网络检测依赖配置选项，然后点击按钮开始进行网络状态检测
            </p>
          </div>
          <div class="progress-area" v-show="state !== STATE_INIT">
            <div class="stats shadow w-full">
              <div class="stat text-left">
                <div class="stat-title truncate">{{ filename }}</div>
                <div
                  v-if="lineNum && state === STATE_WAIT"
                  class="stat-value mt-2 text-center"
                >
                  <span class="stat-desc font-normal">已处理</span>
                  {{ lineNum }}
                  <span class="stat-desc font-normal">行</span>
                </div>
                <div
                  v-else-if="state === STATE_DONE"
                  class="stat-value mt-2 text-center"
                >
                  <span class="stat-desc text-2xl">处理完成</span>
                </div>
                <div v-else class="stat-value mt-2 text-center">
                  <span class="stat-desc text-2xl">等待处理</span>
                </div>
                <div
                  v-show="state === STATE_DONE"
                  class="stat-desc truncate mt-2"
                >
                  DNS解析列表：{{ dnsRes }}
                </div>
                <div
                  v-show="state === STATE_DONE"
                  class="stat-desc truncate mt-2"
                >
                  网络下行速度：{{ download }} Mb/s
                </div>
                <div
                  v-show="state === STATE_DONE"
                  class="stat-desc truncate mt-2"
                >
                  网络延迟：{{ delay }} ms
                </div>
                <div
                  v-show="state === STATE_DONE"
                  class="stat-desc truncate mt-2"
                >
                  网络丢包率：{{ loss }} %
                </div>
              </div>
            </div>
          </div>
        </div>
        <div
          class="select-area relative transition duration-500 select-none cursor-pointer mx-auto mt-1"
          :class="{
            'state-init': state === STATE_INIT,
            'state-wait': state === STATE_WAIT,
            'state-done': state === STATE_DONE,
          }"
          @click="handleClickButton"
        >
          <div class="add-files-area flex transition duration-300 z-10">
            <div class="add-files button w-full">
              <span class="file-ico ico-add"
                ><svg-icon name="file-medical"
              /></span>
              <span class="file-ico ico-exist"
                ><svg-icon name="file-alt"
              /></span>
              <span class="file-ico ico-wait"><svg-icon name="loading" /></span>
              <span class="file-ico ico-success"
                ><svg-icon name="file-check-alt"
              /></span>
              <span>{{ buttonText }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
    <input type="file" ref="fileEleRef" id="fileEle" class="hidden" />
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from "vue";
import {
  CheckSpeed,
  CheckLatency,
  CheckDNS,
  EventsOn,
  LogInfo,
} from "../../wailsjs";

// 状态 1 等待检测  2 处理中 3 检测完毕
const STATE_INIT = 1;
const STATE_WAIT = 2;
const STATE_DONE = 3;
const state = ref(STATE_INIT);

const buttonText = ref("开始检测");
const filePath = ref([""]);
const filename = ref("");
const newFilePath = ref("");
const lineNum = ref(0);
const download = ref(0); // 下行速度
const delay = ref(0); // 网络延迟
const loss = ref(0); // 丢包率
const dnsRes = ref("") // DNS解析结果

const handleClickButton = () => {
  switch (state.value) {
    case STATE_INIT:
      test();
      break;
    case STATE_DONE:
      state.value = STATE_INIT;
      filePath.value = "";
      newFilePath.value = "";
      lineNum.value = 0;
      download.value = 0;
      loss.value = 0;
      delay.value = 0;
      dnsRes.value = ""
      break;
  }
};

const test = async () => {
  state.value = STATE_WAIT;
  const downloadSpeed = await CheckSpeed();
  const late = await CheckLatency();
  const list = await CheckDNS();
  LogInfo(`处理结果：${list}`);
  state.value = STATE_DONE;
  if (downloadSpeed > 0) {
    download.value = downloadSpeed;
  }
  if (late !== undefined && late != null && late.length > 0 && late[0] > 0) {
    loss.value = late[0];
  }
  if (late !== undefined && late != null && late.length > 1 && late[1] > 0) {
    delay.value = late[1];
  }
  if(typeof list !== 'undefined' && list != null && list !== '') {
    dnsRes.value = list;
  }
};

const onFilterChange = () => {
  EventsOn("filter-change", (num) => {
    lineNum.value = num;
  });
};

watch(state, (val) => {
  if (val === STATE_INIT) buttonText.value = "开始检测";
  else if (val === STATE_WAIT) buttonText.value = "检测中...";
  else if (val === STATE_DONE) {
    buttonText.value = "检测完成";
  }
});

onMounted(() => {
  onFilterChange();
});
</script>

<style lang="less" scoped>
@import "../animation.less";
.upload-container {
  .text-focus-in();
  .file-info-area {
    .progress-area {
      .stats {
        background-color: transparent;
        box-shadow: 0 0 16px 1px rgb(0 0 0 / 20%);
      }
    }
  }
  .select-area {
    width: 240px;
    height: 60px;
    border-radius: 2.5rem;
    color: #252525;
    box-shadow: 0 0 16px 1px rgb(0 0 0 / 25%);

    &.state-init {
      .ico-add {
        display: block;
      }
      .ico-exist,
      .ico-wait,
      .ico-success {
        display: none;
      }
    }
    &.state-wait {
      .ico-wait {
        display: block;
        .ele-rotate360();
      }
      .ico-add,
      .ico-exist,
      .ico-success {
        display: none;
      }
    }
    &.state-done {
      .ico-success {
        display: block;
      }
      .ico-add,
      .ico-wait,
      .ico-exist {
        display: none;
      }
    }

    .add-files-area {
      padding-left: 35px;
      border-radius: 2.5rem;
      &:hover {
        background-color: #fdda65;
      }
    }
    .add-files {
      font-size: 2rem;
      font-weight: bolder;
      display: flex;
      flex-flow: row;
      align-items: center;
      span {
        font-size: 25px;
        line-height: 60px;
      }
    }

    .reupload-area {
      display: none;
      font-size: 0.875rem;
      background: #f1f1f1;
      height: 3rem;
      border-radius: 1.875rem;
      padding: 0.41667rem;
      transition: 0.25s;
      width: 6.75rem;
      top: 0.4rem;
      right: 0.83333rem;
      &:hover {
        background-color: #fdda65;
      }
      .button {
        width: 6.75rem;
        font-size: 1rem;
        font-weight: 700;
        height: 3rem;
        border-radius: 1.875rem;
        line-height: 3rem;
      }
    }
    .add-files .file-ico {
      font-size: 2.2rem;
      line-height: 60px;
      margin-right: 17px;
      transform: translateY(-1px);
    }
  }
}
</style>
