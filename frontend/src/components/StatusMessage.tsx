interface StatusMessageProps {
  type: 'idle' | 'success' | 'error' | 'loading'
  message: string
  testId?: string
}

export function StatusMessage({ type, message, testId }: StatusMessageProps) {
  if (type === 'idle') return null

  return (
    <div className={`status-message ${type}`} data-testid={testId}>
      {type === 'loading' && <div className="spinner" />}
      {message}
    </div>
  )
}
