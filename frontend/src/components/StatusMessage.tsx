import { useEffect, useState } from 'react'

interface StatusMessageProps {
  type: 'idle' | 'success' | 'error' | 'loading'
  message: string
  testId?: string
}

const LINE_STAGGER = 55 // ms between line reveals

export function StatusMessage({ type, message, testId }: StatusMessageProps) {
  const [done, setDone] = useState(false)

  const lines = message.split('\n')

  useEffect(() => {
    if (type !== 'success') return
    const reduce = window.matchMedia?.('(prefers-reduced-motion: reduce)').matches
    if (reduce) {
      setDone(true)
      return
    }
    setDone(false)
    const total = lines.length * LINE_STAGGER + 220
    const id = window.setTimeout(() => setDone(true), total)
    return () => window.clearTimeout(id)
  }, [type, message]) // eslint-disable-line react-hooks/exhaustive-deps

  if (type === 'idle' || type === 'loading') return null

  if (type === 'error') {
    return (
      <div className="stamp" role="alert" data-testid={testId}>
        {message}
      </div>
    )
  }

  // success → printed ticket
  return (
    <div className="ticket-block" aria-live="polite">
      <div className="tear">{done ? 'order' : 'printing'}</div>
      <pre className="ticket-out" data-testid={testId}>
        {lines.map((line, i) => (
          <span
            key={i}
            className="ticket-line"
            style={{ animationDelay: `${i * LINE_STAGGER}ms` }}
          >
            {line}
            {'\n'}
          </span>
        ))}
        {!done && <span className="caret">▍</span>}
      </pre>
      <div className="tear">{done ? 'tear' : '· · ·'}</div>
      {done && (
        <div className="ticket-caption">torn off · copied · ⌘V into WhatsApp</div>
      )}
    </div>
  )
}
