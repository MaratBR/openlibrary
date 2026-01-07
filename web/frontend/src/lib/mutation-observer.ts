const callbacks = new WeakMap<Node, (record: MutationRecord) => void>()

const observer = new MutationObserver((records) => {
  for (const record of records) {
    const cb = callbacks.get(record.target)
    if (cb) cb(record)
  }
})

export function addGlobalMutationObserverCallback(
  node: Node,
  callback: (record: MutationRecord) => void,
) {
  callbacks.set(node, callback)
  observer.observe(node)
}

export function removeNodeFromGlobalMutationObserver(node: Node) {
  callbacks.delete(node)
}
