<template>
  <div class="draggable-tree-wrapper">
    <el-tree
      ref="treeRef"
      class="draggable-tree"
      :data="data"
      :node-key="nodeKey"
      :props="defaultProps"
      :check-strictly="false"
      :show-checkbox="showCheckbox"
      :default-checked-keys="checkedKeys"
      :expand-on-click-node="false"
      :draggable="false"
      @node-click="handleNodeClick"
      @check-change="handleCheckChange"
    >
      <template #default="{ node, data }">
        <div class="custom-tree-node" :class="{ 'is-dragging': draggingKey === node.key }">
          <span v-if="draggable" class="drag-handle" title="拖拽排序">
            <el-icon><Rank /></el-icon>
          </span>
          <span class="node-label" :title="node.label">
            {{ node.label }}
          </span>
          <span class="node-actions" @click.stop>
            <el-button
              v-if="data.actions?.includes('add') || showDefaultActions"
              type="primary"
              link
              size="small"
              @click.stop="handleAction('add', node, data)"
            >
              <el-icon><Plus /></el-icon>
              <span>新增</span>
            </el-button>
            <el-button
              v-if="data.actions?.includes('edit') || showDefaultActions"
              type="primary"
              link
              size="small"
              @click.stop="handleAction('edit', node, data)"
            >
              <el-icon><Edit /></el-icon>
              <span>编辑</span>
            </el-button>
            <el-button
              v-if="data.actions?.includes('delete') || showDefaultActions"
              type="danger"
              link
              size="small"
              @click.stop="handleAction('delete', node, data)"
            >
              <el-icon><Delete /></el-icon>
              <span>删除</span>
            </el-button>
            <slot name="node-actions" :node="node" :data="data" />
          </span>
        </div>
      </template>
    </el-tree>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, onBeforeUnmount, nextTick, computed } from 'vue'
import Sortable from 'sortablejs'

const props = defineProps({
  data: {
    type: Array,
    default: () => []
  },
  draggable: {
    type: Boolean,
    default: false
  },
  showCheckbox: {
    type: Boolean,
    default: false
  },
  checkedKeys: {
    type: Array,
    default: () => []
  },
  defaultProps: {
    type: Object,
    default: () => ({
      label: 'label',
      children: 'children'
    })
  },
  nodeKey: {
    type: String,
    default: 'id'
  },
  showDefaultActions: {
    type: Boolean,
    default: false
  },
  allowDrag: {
    type: Function,
    default: null
  },
  allowDrop: {
    type: Function,
    default: null
  }
})

const emit = defineEmits(['node-drop', 'node-click', 'check-change', 'node-action'])

const treeRef = ref(null)
const sortableInstances = ref([])
const draggingKey = ref(null)

const propsLabel = computed(() => props.defaultProps.label || 'label')
const propsChildren = computed(() => props.defaultProps.children || 'children')

function getNodeByKey(data, key, childrenKey = 'children') {
  for (const node of data) {
    if (node[props.nodeKey] === key) return node
    if (node[childrenKey] && node[childrenKey].length) {
      const found = getNodeByKey(node[childrenKey], key, childrenKey)
      if (found) return found
    }
  }
  return null
}

function getParentAndIndex(data, key, parent = null, childrenKey = 'children') {
  for (let i = 0; i < data.length; i++) {
    if (data[i][props.nodeKey] === key) {
      return { parent, list: data, index: i }
    }
    if (data[i][childrenKey] && data[i][childrenKey].length) {
      const found = getParentAndIndex(data[i][childrenKey], key, data[i], childrenKey)
      if (found) return found
    }
  }
  return null
}

function handleNodeClick(data, node) {
  emit('node-click', {
    data,
    node,
    nodeData: { ...data }
  })
}

function handleCheckChange(data, checked, indeterminate) {
  const checkedNodes = treeRef.value?.getCheckedNodes() || []
  const checkedKeysResult = treeRef.value?.getCheckedKeys() || []
  const halfCheckedNodes = treeRef.value?.getHalfCheckedNodes() || []
  const halfCheckedKeys = treeRef.value?.getHalfCheckedKeys() || []

  emit('check-change', {
    data,
    checked,
    indeterminate,
    checkedNodes,
    checkedKeys: checkedKeysResult,
    halfCheckedNodes,
    halfCheckedKeys
  })
}

function handleAction(action, node, data) {
  emit('node-action', {
    action,
    node,
    data,
    nodeData: { ...data }
  })
}

function initSortable() {
  destroySortable()

  if (!props.draggable || !treeRef.value) return

  nextTick(() => {
    const treeEl = treeRef.value?.$el
    if (!treeEl) return

    const childLists = treeEl.querySelectorAll('.el-tree-node > .el-tree-node__children')
    const rootList = treeEl.querySelector('.el-tree-node.is-hidden ~ .el-tree-node') 
      ? null 
      : treeEl.querySelector(':scope > .el-tree > .el-tree-node__children') || treeEl

    const allLists = [treeEl.querySelector('.el-tree'), ...childLists].filter(Boolean)

    allLists.forEach((listEl) => {
      const sortable = Sortable.create(listEl, {
        group: 'draggable-tree-group',
        animation: 200,
        handle: '.drag-handle',
        ghostClass: 'sortable-ghost',
        chosenClass: 'sortable-chosen',
        dragClass: 'sortable-drag',
        forceFallback: true,
        fallbackClass: 'sortable-fallback',
        fallbackOnBody: true,
        fallbackTolerance: 4,

        onStart: (evt) => {
          const el = evt.item
          const nodeKey = el.getAttribute('data-key')
          draggingKey.value = nodeKey
        },

        onMove: (evt) => {
          if (!props.allowDrag && !props.allowDrop) return true

          const draggedEl = evt.dragged
          const draggedKey = draggedEl.getAttribute('data-key')
          const draggedNode = getNodeByKey(props.data, draggedKey, propsChildren.value)

          const relatedEl = evt.related
          const relatedKey = relatedEl.getAttribute('data-key')
          const relatedNode = getNodeByKey(props.data, relatedKey, propsChildren.value)

          if (props.allowDrag && !props.allowDrag(draggedNode)) {
            return false
          }

          if (props.allowDrop && relatedNode) {
            if (!props.allowDrop(draggedNode, relatedNode, evt.willInsertAfter)) {
              return false
            }
          }

          return true
        },

        onEnd: (evt) => {
          draggingKey.value = null

          const { item, to, from, newIndex, oldIndex } = evt

          if (to === from && newIndex === oldIndex) return

          const movedKey = item.getAttribute('data-key')
          const movedNode = getNodeByKey(props.data, movedKey, propsChildren.value)

          if (!movedNode) return

          const sourceInfo = getParentAndIndex(props.data, movedKey, null, propsChildren.value)
          if (!sourceInfo) return

          const sourceList = sourceInfo.list
          const sourceIndex = sourceInfo.index
          const sourceParent = sourceInfo.parent

          const movedItem = sourceList.splice(sourceIndex, 1)[0]

          let targetInfo = null
          let targetParent = null

          if (to.classList.contains('el-tree')) {
            targetInfo = { parent: null, list: props.data }
          } else {
            const parentNodeEl = to.closest('.el-tree-node')
            if (parentNodeEl) {
              const parentKey = parentNodeEl.getAttribute('data-key')
              targetParent = getNodeByKey(props.data, parentKey, propsChildren.value)
              if (targetParent) {
                if (!targetParent[propsChildren.value]) {
                  targetParent[propsChildren.value] = []
                }
                targetInfo = { parent: targetParent, list: targetParent[propsChildren.value] }
              }
            }
          }

          if (!targetInfo) {
            sourceList.splice(sourceIndex, 0, movedItem)
            return
          }

          const targetList = targetInfo.list
          const actualNewIndex = Math.min(newIndex, targetList.length)
          targetList.splice(actualNewIndex, 0, movedItem)

          emit('node-drop', {
            movedNode: movedItem,
            movedData: { ...movedItem },
            sourceParent: sourceParent ? { ...sourceParent } : null,
            targetParent: targetParent ? { ...targetParent } : null,
            sourceIndex,
            targetIndex: actualNewIndex,
            oldData: JSON.parse(JSON.stringify(props.data))
          })
        }
      })

      sortableInstances.value.push(sortable)
    })

    const treeNodes = treeEl.querySelectorAll('.el-tree-node')
    treeNodes.forEach((nodeEl) => {
      const contentEl = nodeEl.querySelector(':scope > .el-tree-node__content')
      if (contentEl) {
        const key = nodeEl.getAttribute('data-key')
        contentEl.setAttribute('data-key', key)
        nodeEl.setAttribute('data-key', key)
      }
    })
  })
}

function destroySortable() {
  sortableInstances.value.forEach((instance) => {
    instance.destroy()
  })
  sortableInstances.value = []
}

watch(
  () => [props.draggable, props.data],
  () => {
    if (props.draggable) {
      nextTick(() => {
        initSortable()
      })
    } else {
      destroySortable()
    }
  },
  { deep: true, immediate: false }
)

watch(
  () => props.checkedKeys,
  (newKeys) => {
    if (treeRef.value) {
      nextTick(() => {
        treeRef.value.setCheckedKeys(newKeys)
      })
    }
  },
  { deep: true }
)

defineExpose({
  getCheckedNodes: () => treeRef.value?.getCheckedNodes() || [],
  getCheckedKeys: () => treeRef.value?.getCheckedKeys() || [],
  getHalfCheckedNodes: () => treeRef.value?.getHalfCheckedNodes() || [],
  getHalfCheckedKeys: () => treeRef.value?.getHalfCheckedKeys() || [],
  setCheckedKeys: (keys, leafOnly = false) => treeRef.value?.setCheckedKeys(keys, leafOnly),
  setCheckedNodes: (nodes, leafOnly = false) => treeRef.value?.setCheckedNodes(nodes, leafOnly),
  setChecked: (key, checked, deep = true) => treeRef.value?.setChecked(key, checked, deep),
  getNode: (key) => treeRef.value?.getNode(key),
  remove: (data) => treeRef.value?.remove(data),
  append: (data, parentNode) => treeRef.value?.append(data, parentNode),
  insertBefore: (data, refNode) => treeRef.value?.insertBefore(data, refNode),
  insertAfter: (data, refNode) => treeRef.value?.insertAfter(data, refNode),
  updateKeyChildren: (key, children) => treeRef.value?.updateKeyChildren(key, children),
  filter: (value) => treeRef.value?.filter(value),
  setCurrentKey: (key) => treeRef.value?.setCurrentKey(key),
  getCurrentKey: () => treeRef.value?.getCurrentKey(),
  getCurrentNode: () => treeRef.value?.getCurrentNode(),
  setCurrentNode: (node) => treeRef.value?.setCurrentNode(node),
  $el: () => treeRef.value?.$el,
  refresh: () => initSortable()
})

onMounted(() => {
  if (props.draggable) {
    nextTick(() => {
      initSortable()
    })
  }
})

onBeforeUnmount(() => {
  destroySortable()
})
</script>

<style lang="scss" scoped>
.draggable-tree-wrapper {
  width: 100%;
  padding: 8px 0;
}

.draggable-tree {
  background: transparent;
  user-select: none;
}

:deep(.el-tree-node) {
  .el-tree-node__content {
    height: 40px;
    padding-left: 8px !important;
    border-radius: 6px;
    transition: background-color 0.2s ease;

    &:hover {
      background-color: #f0f7ff;

      .node-actions {
        opacity: 1;
      }
    }

    .custom-tree-node.is-dragging {
      opacity: 0.5;
    }
  }

  &.is-current > .el-tree-node__content {
    background-color: #ecf5ff;
    color: #409eff;
  }

  &.is-expanded > .el-tree-node__expand-icon {
    transform: rotate(90deg);
  }
}

:deep(.el-tree-node__children) {
  overflow: hidden;
  transition: max-height 0.3s ease;
}

:deep(.el-tree-node__expand-icon) {
  transition: transform 0.3s ease;
  color: #909399;
  cursor: pointer;

  &:hover {
    color: #409eff;
  }
}

.custom-tree-node {
  display: flex;
  align-items: center;
  width: 100%;
  padding-right: 8px;
  box-sizing: border-box;

  .drag-handle {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    margin-right: 4px;
    color: #c0c4cc;
    cursor: move;
    border-radius: 4px;
    transition: all 0.2s ease;

    &:hover {
      color: #409eff;
      background-color: #ecf5ff;
    }

    .el-icon {
      font-size: 16px;
    }
  }

  .node-label {
    flex: 1;
    font-size: 14px;
    color: #303133;
    line-height: 40px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .node-actions {
    display: flex;
    align-items: center;
    gap: 4px;
    opacity: 0;
    transition: opacity 0.2s ease;

    .el-button {
      padding: 4px 6px;
      font-size: 12px;
      height: 28px;

      .el-icon {
        font-size: 14px;
      }

      span {
        margin-left: 2px;
      }
    }
  }
}

.sortable-ghost {
  opacity: 0.4;
  background-color: #d9ecff !important;
  border: 2px dashed #409eff;
}

.sortable-chosen {
  .drag-handle {
    color: #409eff !important;
  }
}

.sortable-drag {
  background-color: #ffffff;
  box-shadow: 0 8px 24px rgba(64, 158, 255, 0.2);
  border-radius: 8px;
  z-index: 9999;
}

.sortable-fallback {
  opacity: 0.9;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}
</style>
