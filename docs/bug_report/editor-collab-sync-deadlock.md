# 编辑器中同步遮罩死锁 / Editor collaboration-sync overlay deadlock

| 字段 | 值 |
| --- | --- |
| 状态 / Status | 已修复 / Fixed（未提交 / uncommitted） |
| 修复版本 / Fixed in | see `docs/CHANGELOG` |
| 影响范围 / Impact | 富文本 (`yoresee_rich_text`, docType `4`) 与 markdown (`1`) 文档；表格 (`2`) / 幻灯片 (`3`) 不受影响 |
| 数据后果 / Data consequence | 无数据损坏；纯前端死锁，文档内容始终无法渲染，编辑/协作能力整体不可用 |
| 复现路径 / Reproduction | 打开任意富文本或 markdown 文档 → 一直显示「正在同步文档」，`MarkdownEditor` / `YoreseeRichTextEditor` 组件从不挂载 |

---

## 现象 / Symptom

打开富文本或 markdown 文档时，编辑器区域 permanently 显示加载文案「正在同步文档」，排版/编辑区不出现。

浏览器控制台可见以下关键日志（`targetDocType=4` 即富文本文档）：

```
[STATE] isReady= true collabEnabled= true collabSynced= false collabAttempted= true targetDocType= 4
```

`collabAttempted=true` 且 `collabSynced=false` 说明：同步就绪标记被提前置为「已尝试、未同步」，但 `handleCollabSync` 从未被调用（编辑器从未挂载，WebSocket 从未建立）。

## 排查过程 / Investigation

### Phase 1 — 锚定死锁点

状态机的三态联合决定遮罩是否显示。`DocumentEditor.vue` 模板结构当时是：

```html
<div v-if="isReady && collabEnabled && !collabSynced && collabAttempted" class="editor-loading">正在同步…</div>
<div v-else-if="!isReady && !isError" class="editor-loading">加载中…</div>
<div v-else-if="isError" class="editor-loading">错误</div>
<MarkdownEditor        v-else-if="type==='1'" … />
<YoreseeRichTextEditor v-else-if="type==='4'" … />
```

关键问题：**加载遮罩是 `v-if/v-else-if` 互斥链的第一分支**。当 `v-if` 为真时，它「替换」掉编辑器，而不是「覆盖」在编辑器之上。遮罩条件一旦成立，编辑器分支（`v-else-if`）永远不进入渲染 → 编辑器从不挂载。

### Phase 2 — 为什么遮罩条件恒成立

`collabAttempted` / `collabSynced` 由 `syncCollabReadyFlag()` 在 `resolveTargetDocument` 内、`resolveContent`（置 `isReady=true`）之后立即调用：

```js
const syncCollabReadyFlag = () => {
  const shouldBeSynced = !collabEnabled.value;        // collab 开启时为 false
  docMachine.markCollabSynced(shouldBeSynced, docId.value);
};
// markCollabSynced(false, docId) → collabAttempted=true, collabSynced=false
```

于是导航到协作文档时：

1. `startNavigation` → `collabSynced=true`（初始值，防陈旧）
2. `resolveContent` → `isReady=true`，此时渲染 flush 触发；遮罩条件 `!collabSynced` 仍为真 → 遮罩分支胜出
3. `syncCollabReadyFlag` → `collabAttempted=true, collabSynced=false`
4. 后续渲染 flush：遮罩条件 `isReady && collabEnabled && !collabSynced && collabAttempted` 恒为真 → 编辑器分支永不进入
5. 编辑器不挂载 → `onMounted` 不跑 → 不调 `setupCollaboration()` → WebSocket 不连接 → `collab-sync` 事件永不触发 → `handleCollabSync` 永不调用 → `collabSynced` 卡在 `false` → **死锁**

`[MOUNT] MarkdownEditor` / `[RTC mounted]` 这类挂载日志从头到尾从未出现，印证编辑器确实没挂载。

### Phase 3 — 还原原始设计确认根因

查看重构前（commit `66131ef`）的实现：

```html
<div v-if="collabEnabled && !collabReady" class="editor-loading">加载中</div>
<MarkdownEditor v-if="isMarkdownDocument" … />
```

加载遮罩是**独立 `v-if`**，`YoreseeRichTextEditor` / `MarkdownEditor` 也是独立 `v-if`——两者**并行渲染**。遮罩依赖 `.editor-loading { position:absolute; inset:0; z-index:2 }` 浮在已挂载的编辑器之上，而非互斥替换。

本次状态机重构（`0fad5c6` → `d82be2a` → `ee9abb5` → `c045ddc`）错误地把遮罩并入互斥 `v-else-if` 链，并用布尔对 `collabSynced` + `collabAttempted` 取代原本的 `lastSyncedDocId`（它记录「哪个文档已确认同步」的身份信息）。丢失身份语义 + 互斥链，两者共同造成死锁。

### Phase 4 — 为什么表格/幻灯片不受影响

`collabEnabled` 仅对 markdown / rich-text 为 `true`（见 `useDocumentEditorPolicy`）。表格/幻灯片 `collabEnabled=false`，遮罩条件第一项就为假，编辑器直接渲染，不参与死锁。

## 修复 / Fix

### 方案：恢复「遮罩覆盖编辑器」的原始语义，并用 `syncedDocId` 重建身份语义

**1. 状态机 `useDocumentStateMachine.js`**

- 删除 `collabSynced` 布尔与 `collabAttempted`，改为 `syncedDocId` ref + 计算属性：

  ```js
  const collabSynced = computed(
    () => targetDocId.value !== '' && syncedDocId.value === targetDocId.value
  );
  ```

- `startNavigation` 不再重置协作态（`syncedDocId` 跨导航保留）：导航到新文档时 `syncedDocId !== targetDocId` 自动变为「未同步」，无需显式布尔；导航回已同步的同文档则保持同步态、不闪屏。
- `markCollabSynced(isSynced, docId)`：`docId` 守卫不变；同步成功置 `syncedDocId=docId`，断开同步则清空（仅当当前即该文档）。

**2. 模板 `DocumentEditor.vue`**

- 遮罩恢复为 ready 分支内的**独立绝对定位覆盖层**，编辑器链 `v-if/v-else-if` 不再与遮罩互斥：

  ```html
  <template v-else>
    <div v-if="collabEnabled && !docMachine.collabSynced.value" class="editor-loading">加载中</div>
    <MarkdownEditor        v-if="targetDocType==='1'" … />
    <YoreseeRichTextEditor v-else-if="targetDocType==='4'" … />
  </template>
  ```

  编辑器始终挂载；遮罩只在 `collabEnabled && !collabSynced` 时浮于其上。WebSocket 建立 → `collab-sync` → `collabSynced=true` → 遮罩消失。

**3. 删除并行残留状态**

- 移除 `useDocumentRouteContext` 中遗留的 `collabReady`，以及生命周期里对其的赋值与全部调试日志（`[STATE]`/`[COLLAB_SYNC]`/`[DOC_LIFECYCLE]`/`[MOUNT]`/`[RTC mounted]` 等）。
- `useDocumentEditorLifecycle` 的 `syncCollabReadyFlag()` 在 kbId watcher 中原本会无条件 `markCollabSynced(false)`，对已同步同文档会清空 `syncedDocId` 造成卡死；该调用已被移除（其原判断已由 `collabSynced` 计算属性自动表达）。

**4. 关联清理：Tiptap v3 警告**

编辑器挂载后另见两条 `[tiptap warn]`：`Duplicate extension names: ['link','underline']` 与 `Collaboration ... not compatible with extension-undo-redo`。根因是 Tiptap v3 的 `StarterKit` 现已内置 `link`、`underline`、`undoRedo`（自带 history），与显式引入的 `Link`/`Underline` 及协作所需的 `Collaboration` 冲突。已在 `useRichTextEditorRuntime.js` 的 `createEditor` 中：

- `StarterKit.configure({ link:false, underline:false })`，保留显式带自定义配置的 `Link`/`Underline`；
- 检测附加扩展是否含 `Collaboration`（`ext.name === 'collaboration'`），仅在该情况下关闭 StarterKit 的 `undoRedo`，使协作自带 history 生效。

## 验证 / Verification

- 构建通过（`npm run build`）。
- 富文本/ markdown 文档：遮罩短暂出现 → 编辑器挂载 → WebSocket 同步 → 遮罩消失。
- 表格/幻灯片文档：遮罩条件恒假，编辑器直接渲染。
- `[tiptap warn]` 两条警告消失。

## 附录 A — 关键日志原文

```
[STATE] isReady= true collabEnabled= true collabSynced= false collabAttempted= true targetDocType= 4
```

（注：`collabAttempted` / `[STATE]` 仅为排查期临时埋点，修复后已从源码删除。）
