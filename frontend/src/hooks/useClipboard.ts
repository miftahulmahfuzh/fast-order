import { useState } from 'react'

export function useClipboard() {
  const [isCopied, setIsCopied] = useState(false)

  const copy = async (text: string): Promise<boolean> => {
    try {
      await navigator.clipboard.writeText(text)
      setIsCopied(true)
      setTimeout(() => setIsCopied(false), 2000)
      return true
    } catch {
      setIsCopied(false)
      return false
    }
  }

  return { copy, isCopied }
}
