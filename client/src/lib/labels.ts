export function labelText(label: string): { group: string; value: string } {
  const [group, ...rest] = label.split(':')
  return rest.length ? { group, value: rest.join(':') } : { group: '', value: group }
}
