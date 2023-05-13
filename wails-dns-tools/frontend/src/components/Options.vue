<template>
  <div class="menu-area absolute right-5 top-4" data-wails-no-drag>
    <ul class="menu bg-base-100 p-2 rounded-box text-xl">
      <li>
        <a @click="handleOptions" class="p-2.5">
          <svg-icon name="setting-config" />
        </a>
      </li>
      <li>
        <a @click="handleQuit" class="p-2.5">
          <svg-icon name="power" />
        </a>
      </li>
    </ul>
  </div>

  <div
    class="options-panel absolute bg-white px-5 py-7 w-7/12 h-5/6 rounded-2xl"
    v-show="showOptionsPanel"
    data-wails-no-drag
  >
    <div class="form-control">
      <div class="divider">导出PDF设置</div>
      <div
        class="w-full px-8 py-2 pb-4 rounded-xl hover:bg-slate-50"
        v = config
      >
        <label class="label">
          <span class="label-text">PDF页面高度</span>
        </label>
        <input
          type="text" Height
          placeholder="height"
          v-model="config.pdfHeight"
          class="input input-bordered w-full"
        />
        <label class="label">
          <span class="label-text">PDF页面宽度</span>
        </label>
        <input
          type="text" Width
          placeholder="width"
          v-model="config.pdfWidth"
          class="input input-bordered w-full"
        />
        <label class="label">
          <span class="label-text">PDF保存路径</span>
        </label>
        <div>当前选择路径：{{config.pdfExportPath}}</div>
        <button class="btn btn-wide btn-outline mx-auto mt-8 confirm"
                @click="selectSave">选择路径</button>

        </div>
      </div>
      <button
        class="btn btn-wide btn-outline mx-auto mt-8 confirm"
        @click="confirm"
      >
        确定
      </button>
  </div>
</template>

<script setup>
import { ref, reactive } from "vue";
import { Quit, GetOptions, SelectSavePath, SetOptions, LogInfo } from "../../wailsjs";

const config = reactive({
  pdfHeight: 768,
  pdfWidth: 1024,
  pdfExportPath: ""
});
const showOptionsPanel = ref(false);

const handleQuit = () => Quit();

const handleOptions = async () => {
  const res = await GetOptions();
  LogInfo(`获取到的配置：${JSON.stringify(res)}`);
  console.log("获取到的配置：", res);
  config.pdfHeight = res.pdfHeight;
  config.pdfWidth = res.pdfWidth;
  config.pdfExportPath = res.pdfExportPath;
  showOptionsPanel.value = true;
};

const selectSave = async () => {
  const res = await SelectSavePath();
  LogInfo(`选择的路径：${JSON.stringify(res)}`);
  console.log("选择的路径：", res);
  config.pdfExportPath = res;
};

const confirm = () => {
  showOptionsPanel.value = false;
  const data = JSON.stringify({
    pdfHeight: config.pdfHeight,
    pdfWidth: config.pdfWidth,
    pdfExportPath: config.pdfExportPath
  });
  LogInfo(`存储的配置：${JSON.stringify(data)}`);
  console.log("config", data);
  SetOptions({
    pdfHeight: config.pdfHeight,
    pdfWidth: config.pdfWidth,
    pdfExportPath: config.pdfExportPath,
  });
};
</script>

<style lang="less" scoped>
@import "../animation.less";
.menu-area {
  .text-focus-in();
  .menu {
    background-color: transparent;
    box-shadow: 0 0 12px 1px rgb(0 0 0 / 15%);
  }
}
.options-panel {
  transition: 0.25s;
  &::after {
    position: absolute;
    content: "";
    background-image: url(../assets/images/background-image.png);
    background-repeat: no-repeat;
    background-attachment: fixed;
    background-size: cover;
    background-position: top;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    z-index: 0;
    filter: blur(10px);
    pointer-events: none;
  }
}
.confirm {
  font-size: 1rem;
  font-weight: 700;
  transition: 0.25s;
  &:hover {
    background-color: #fdd9659b;
    border-color: #fdd9659b;
    color: #252525;
  }
}
</style>
