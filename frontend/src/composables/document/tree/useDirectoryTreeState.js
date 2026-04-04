import { computed, nextTick, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { normalizeDocumentType } from '@/utils/documentType';

export function useDirectoryTreeState({
  t,
  router,
  kbId,
  docId,
  treeComponentRef,
  getKnowledgeBaseDocuments,
  getMyDocuments
}) {
  const treeLoading = ref(false);
  const directoryTree = ref([]);
  const knowledgeBaseName = ref('示例知识库');
  const currentDocTitle = ref('示例文档');
  const currentDocType = ref('1');
  const isAllExpanded = ref(true);

  const treeRef = computed(() => treeComponentRef.value?.getTreeRef());

  const transformDocumentsToTree = (documents, parentId = null) => {
    const tree = [];
    documents.forEach((doc) => {
      const treeNode = {
        id: doc.external_id,
        label: doc.title,
        isFolder: !!doc.has_children,
        isLeaf: !doc.has_children,
        type: normalizeDocumentType(doc.type, '1'),
        parentId,
        children: []
      };
      if (doc.children && doc.children.length > 0) {
        treeNode.children = transformDocumentsToTree(doc.children, treeNode.id);
      }
      tree.push(treeNode);
    });
    return tree;
  };

  const findNodeById = (nodes, targetId) => {
    for (const node of nodes) {
      if (String(node.id) === String(targetId)) {
        return node;
      }
      if (node.children && node.children.length > 0) {
        const found = findNodeById(node.children, targetId);
        if (found) {
          return found;
        }
      }
    }
    return null;
  };

  const updateCurrentDocTitle = () => {
    if (!docId.value || docId.value === 'example') {
      currentDocType.value = '1';
      return;
    }
    const node = findNodeById(directoryTree.value, docId.value);
    if (node) {
      currentDocTitle.value = node.label;
      currentDocType.value = normalizeDocumentType(node.type, '1');
      return;
    }
    currentDocType.value = '1';
  };

  const updateTreeNodeTitle = (nodes, targetId, title) => {
    for (const node of nodes) {
      if (String(node.id) === String(targetId)) {
        node.label = title;
        return true;
      }
      if (node.children && node.children.length > 0) {
        if (updateTreeNodeTitle(node.children, targetId, title)) {
          return true;
        }
      }
    }
    return false;
  };

  const findPathToNode = (tree, targetId, path = []) => {
    for (const node of tree) {
      if (node.id === targetId) {
        return [...path, node];
      }
      if (node.children && node.children.length > 0) {
        const result = findPathToNode(node.children, targetId, [...path, node]);
        if (result) {
          return result;
        }
      }
    }
    return null;
  };

  const expandToCurrentDoc = async () => {
    if (!treeRef.value || !docId.value || directoryTree.value.length === 0) {
      return;
    }
    const path = findPathToNode(directoryTree.value, docId.value);
    if (path && path.length > 0) {
      await nextTick();
      for (let i = 0; i < path.length - 1; i += 1) {
        const node = treeRef.value.getNode(path[i].id);
        if (node) {
          node.expanded = true;
        }
      }
    }
  };

  const fetchDocuments = async () => {
    if (kbId.value === 'example') {
      return;
    }

    treeLoading.value = true;
    try {
      if (kbId.value === 'personal') {
        const response = await getMyDocuments({ page: 1, page_size: 1000, directory_only: true });
        knowledgeBaseName.value = t('home.myDocuments');
        directoryTree.value = transformDocumentsToTree(response.documents || []);
      } else {
        const response = await getKnowledgeBaseDocuments(kbId.value, { directory_only: true });
        knowledgeBaseName.value = response.knowledge_base.name;
        directoryTree.value = transformDocumentsToTree(response.documents);
      }
      await expandToCurrentDoc();
      updateCurrentDocTitle();
    } catch (error) {
      console.error('获取文档列表失败:', error);
      ElMessage.error(t('knowledgeBase.fetchError'));
    } finally {
      treeLoading.value = false;
    }
  };

  const goBack = () => {
    if (kbId.value === 'personal' || kbId.value === 'example') {
      router.push('/mydocuments');
    } else {
      router.push(`/knowledge-base/${kbId.value}`);
    }
  };

  const handleTreeNodeClick = (data) => {
    if (data?.isCreating) return;
    if (!data?.id) return;
    if (kbId.value === 'personal') {
      router.push(`/mydocument/${data.id}`);
    } else {
      router.push(`/knowledge-base/${kbId.value}/document/${data.id}`);
    }
  };

  const toggleExpandAll = () => {
    isAllExpanded.value = !isAllExpanded.value;
    if (!treeRef.value) {
      return;
    }
    const nodes = treeRef.value.store?.nodesMap;
    if (!nodes) {
      return;
    }
    Object.values(nodes).forEach((node) => {
      node.expanded = isAllExpanded.value;
    });
  };

  const closeContextMenu = () => {
    if (treeComponentRef.value) {
      treeComponentRef.value.closeContextMenu?.();
    }
  };

  return {
    treeLoading,
    directoryTree,
    knowledgeBaseName,
    currentDocTitle,
    currentDocType,
    isAllExpanded,
    fetchDocuments,
    updateCurrentDocTitle,
    updateTreeNodeTitle,
    expandToCurrentDoc,
    goBack,
    handleTreeNodeClick,
    toggleExpandAll,
    closeContextMenu
  };
}
