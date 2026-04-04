<template>
  <TableDiffViewer
    v-if="previewKind === 'table'"
    :left-rows="tableDiffRows.leftRows"
    :right-rows="tableDiffRows.rightRows"
    :left-title="leftTitle"
    :right-title="rightTitle"
  />
  <TextDiffViewer
    v-else
    :left-text="resolvedTextDiff.leftText"
    :right-text="resolvedTextDiff.rightText"
    :left-title="leftTitle"
    :right-title="rightTitle"
  />
</template>

<script setup>
import { computed } from 'vue';
import TextDiffViewer from '@/components/document/TextDiffViewer.vue';
import TableDiffViewer from '@/components/document/render/TableDiffViewer.vue';
import {
  resolveDocumentDiffContentPair,
  resolveDocumentPreviewKind,
  resolveTableDiffRowsPair
} from '@/composables/document/render/useDocumentRenderBridge';

const props = defineProps({
  leftContent: {
    type: String,
    default: ''
  },
  rightContent: {
    type: String,
    default: ''
  },
  documentType: {
    type: [String, Number],
    default: '1'
  },
  leftTitle: {
    type: String,
    default: ''
  },
  rightTitle: {
    type: String,
    default: ''
  }
});

const previewKind = computed(() => resolveDocumentPreviewKind(props.documentType));
const resolvedTextDiff = computed(() => resolveDocumentDiffContentPair({
  leftContent: props.leftContent,
  rightContent: props.rightContent,
  documentType: props.documentType
}));
const tableDiffRows = computed(() => resolveTableDiffRowsPair({
  leftContent: props.leftContent,
  rightContent: props.rightContent
}));
</script>
